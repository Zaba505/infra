package errors

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zaba505/infra/services/machine/models"
)

func TestProblem_WriteHttpResponse(t *testing.T) {
	p := &Problem{
		Type:     "https://api.example.com/errors/test",
		Title:    "Test Error",
		Status:   http.StatusBadRequest,
		Detail:   "This is a test error",
		Instance: "/test",
	}

	w := httptest.NewRecorder()
	p.WriteHttpResponse(context.Background(), w)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/problem+json" {
		t.Errorf("expected Content-Type 'application/problem+json', got '%s'", contentType)
	}

	var response Problem
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Type != p.Type {
		t.Errorf("expected Type '%s', got '%s'", p.Type, response.Type)
	}
	if response.Title != p.Title {
		t.Errorf("expected Title '%s', got '%s'", p.Title, response.Title)
	}
	if response.Status != p.Status {
		t.Errorf("expected Status %d, got %d", p.Status, response.Status)
	}
	if response.Detail != p.Detail {
		t.Errorf("expected Detail '%s', got '%s'", p.Detail, response.Detail)
	}
	if response.Instance != p.Instance {
		t.Errorf("expected Instance '%s', got '%s'", p.Instance, response.Instance)
	}
}

func TestValidationProblem_WriteHttpResponse(t *testing.T) {
	fields := []models.InvalidField{
		{Field: "nics", Reason: "at least one NIC is required"},
	}
	vp := NewValidationError("/api/v1/machines", fields)

	w := httptest.NewRecorder()
	vp.WriteHttpResponse(context.Background(), w)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/problem+json" {
		t.Errorf("expected Content-Type 'application/problem+json', got '%s'", contentType)
	}

	var response ValidationProblem
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response.InvalidFields) != 1 {
		t.Errorf("expected 1 invalid field, got %d", len(response.InvalidFields))
	}
}

func TestConflictProblem_WriteHttpResponse(t *testing.T) {
	cp := NewConflictError("/api/v1/machines", "aa:bb:cc:dd:ee:ff", "018c7dbd-a000-7000-8000-fedcba987650")

	w := httptest.NewRecorder()
	cp.WriteHttpResponse(context.Background(), w)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/problem+json" {
		t.Errorf("expected Content-Type 'application/problem+json', got '%s'", contentType)
	}

	var response ConflictProblem
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.MACAddress != "aa:bb:cc:dd:ee:ff" {
		t.Errorf("expected MACAddress 'aa:bb:cc:dd:ee:ff', got '%s'", response.MACAddress)
	}
	if response.ExistingMachineID != "018c7dbd-a000-7000-8000-fedcba987650" {
		t.Errorf("expected ExistingMachineID '018c7dbd-a000-7000-8000-fedcba987650', got '%s'", response.ExistingMachineID)
	}
}

func TestNewValidationError(t *testing.T) {
	fields := []models.InvalidField{
		{Field: "test", Reason: "test reason"},
	}
	vp := NewValidationError("/test", fields)

	if vp.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, vp.Status)
	}
	if vp.Title != "Validation Error" {
		t.Errorf("expected title 'Validation Error', got '%s'", vp.Title)
	}
	if len(vp.InvalidFields) != 1 {
		t.Errorf("expected 1 invalid field, got %d", len(vp.InvalidFields))
	}
}

func TestNewConflictError(t *testing.T) {
	cp := NewConflictError("/test", "aa:bb:cc:dd:ee:ff", "test-id")

	if cp.Status != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, cp.Status)
	}
	if cp.Title != "Duplicate MAC Address" {
		t.Errorf("expected title 'Duplicate MAC Address', got '%s'", cp.Title)
	}
	if cp.MACAddress != "aa:bb:cc:dd:ee:ff" {
		t.Errorf("expected MACAddress 'aa:bb:cc:dd:ee:ff', got '%s'", cp.MACAddress)
	}
	if cp.ExistingMachineID != "test-id" {
		t.Errorf("expected ExistingMachineID 'test-id', got '%s'", cp.ExistingMachineID)
	}
}

func TestNewInternalError(t *testing.T) {
	p := NewInternalError("/test", "Something went wrong")

	if p.Status != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, p.Status)
	}
	if p.Title != "Internal Server Error" {
		t.Errorf("expected title 'Internal Server Error', got '%s'", p.Title)
	}
	if p.Detail != "Something went wrong" {
		t.Errorf("expected detail 'Something went wrong', got '%s'", p.Detail)
	}
}
