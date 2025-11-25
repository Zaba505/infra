package endpoint

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zaba505/infra/services/machine/errors"
	"github.com/Zaba505/infra/services/machine/service"
)

type mockFirestoreClient struct {
	machines      map[string]*service.MachineRequest
	createErr     error
	findByMACErr  error
	existingMacID string
}

func (m *mockFirestoreClient) CreateMachine(ctx context.Context, req *service.CreateMachineRequest) (*service.CreateMachineResponse, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.machines == nil {
		m.machines = make(map[string]*service.MachineRequest)
	}
	m.machines[req.MachineID] = req.Machine
	return &service.CreateMachineResponse{}, nil
}

func (m *mockFirestoreClient) FindMachineByMAC(ctx context.Context, req *service.FindMachineByMACRequest) (*service.FindMachineByMACResponse, error) {
	if m.findByMACErr != nil {
		return nil, m.findByMACErr
	}
	if m.existingMacID != "" {
		return &service.FindMachineByMACResponse{
			MachineID: m.existingMacID,
			Found:     true,
		}, nil
	}
	return &service.FindMachineByMACResponse{Found: false}, nil
}

func (m *mockFirestoreClient) Close() error {
	return nil
}

func TestPostMachinesHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("successful machine registration", func(t *testing.T) {
		mock := &mockFirestoreClient{}
		handler := &postMachinesHandler{firestoreClient: mock}

		req := &MachineRequest{
			CPUs: []CPU{
				{Manufacturer: "Intel", ClockFrequency: 2400000000, Cores: 8},
			},
			MemoryModules: []MemoryModule{
				{Size: 17179869184},
			},
			NICs: []NIC{
				{MAC: "52:54:00:12:34:56"},
			},
			Drives: []Drive{
				{Capacity: 500107862016},
			},
		}

		resp, err := handler.Handle(ctx, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.ID == "" {
			t.Error("expected non-empty machine ID")
		}

		if len(mock.machines) != 1 {
			t.Errorf("expected 1 machine in store, got %d", len(mock.machines))
		}
	})

	t.Run("validation error - missing NICs", func(t *testing.T) {
		mock := &mockFirestoreClient{}
		handler := &postMachinesHandler{firestoreClient: mock}

		req := &MachineRequest{
			NICs: []NIC{},
		}

		_, err := handler.Handle(ctx, req)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}

		validationErr, ok := err.(*errors.ValidationProblem)
		if !ok {
			t.Fatalf("expected *errors.ValidationProblem, got %T", err)
		}

		if validationErr.Status != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, validationErr.Status)
		}

		if len(validationErr.InvalidFields) == 0 {
			t.Error("expected invalid fields, got none")
		}
	})

	t.Run("validation error - invalid MAC format", func(t *testing.T) {
		mock := &mockFirestoreClient{}
		handler := &postMachinesHandler{firestoreClient: mock}

		req := &MachineRequest{
			NICs: []NIC{
				{MAC: "invalid-mac"},
			},
		}

		_, err := handler.Handle(ctx, req)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}

		validationErr, ok := err.(*errors.ValidationProblem)
		if !ok {
			t.Fatalf("expected *errors.ValidationProblem, got %T", err)
		}

		if len(validationErr.InvalidFields) == 0 {
			t.Error("expected invalid fields, got none")
		}
	})

	t.Run("conflict error - duplicate MAC address", func(t *testing.T) {
		mock := &mockFirestoreClient{
			existingMacID: "existing-machine-id",
		}
		handler := &postMachinesHandler{firestoreClient: mock}

		req := &MachineRequest{
			NICs: []NIC{
				{MAC: "52:54:00:12:34:56"},
			},
		}

		_, err := handler.Handle(ctx, req)
		if err == nil {
			t.Fatal("expected conflict error, got nil")
		}

		conflictErr, ok := err.(*errors.ConflictProblem)
		if !ok {
			t.Fatalf("expected *errors.ConflictProblem, got %T", err)
		}

		if conflictErr.Status != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, conflictErr.Status)
		}

		if conflictErr.MACAddress != "52:54:00:12:34:56" {
			t.Errorf("expected MAC '52:54:00:12:34:56', got '%s'", conflictErr.MACAddress)
		}

		if conflictErr.ExistingMachineID != "existing-machine-id" {
			t.Errorf("expected existing ID 'existing-machine-id', got '%s'", conflictErr.ExistingMachineID)
		}
	})

	t.Run("internal error - FindMachineByMAC fails", func(t *testing.T) {
		mock := &mockFirestoreClient{
			findByMACErr: fmt.Errorf("firestore error"),
		}
		handler := &postMachinesHandler{firestoreClient: mock}

		req := &MachineRequest{
			NICs: []NIC{
				{MAC: "52:54:00:12:34:56"},
			},
		}

		_, err := handler.Handle(ctx, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		internalErr, ok := err.(*errors.Problem)
		if !ok {
			t.Fatalf("expected *errors.Problem, got %T", err)
		}

		if internalErr.Status != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, internalErr.Status)
		}
	})

	t.Run("internal error - CreateMachine fails", func(t *testing.T) {
		mock := &mockFirestoreClient{
			createErr: fmt.Errorf("firestore create error"),
		}
		handler := &postMachinesHandler{firestoreClient: mock}

		req := &MachineRequest{
			NICs: []NIC{
				{MAC: "52:54:00:12:34:56"},
			},
		}

		_, err := handler.Handle(ctx, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		internalErr, ok := err.(*errors.Problem)
		if !ok {
			t.Fatalf("expected *errors.Problem, got %T", err)
		}

		if internalErr.Status != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, internalErr.Status)
		}
	})
}

func TestErrorHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("handles ValidationProblem", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := errors.NewValidationError("/test", []errors.InvalidField{
			{Field: "test", Reason: "test reason"},
		})

		errorHandler(ctx, w, err)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/problem+json" {
			t.Errorf("expected Content-Type 'application/problem+json', got '%s'", contentType)
		}
	})

	t.Run("handles ConflictProblem", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := errors.NewConflictError("/test", "aa:bb:cc:dd:ee:ff", "test-id")

		errorHandler(ctx, w, err)

		if w.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
		}
	})

	t.Run("handles generic error", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := fmt.Errorf("generic error")

		errorHandler(ctx, w, err)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}
