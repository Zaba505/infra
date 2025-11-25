package service

import (
	"context"
	"testing"
)

type mockFirestoreClient struct {
	machines      map[string]*MachineRequest
	createErr     error
	findByMACErr  error
	existingMacID string
}

func (m *mockFirestoreClient) CreateMachine(ctx context.Context, req *CreateMachineRequest) (*CreateMachineResponse, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.machines == nil {
		m.machines = make(map[string]*MachineRequest)
	}
	m.machines[req.MachineID] = req.Machine
	return &CreateMachineResponse{}, nil
}

func (m *mockFirestoreClient) FindMachineByMAC(ctx context.Context, req *FindMachineByMACRequest) (*FindMachineByMACResponse, error) {
	if m.findByMACErr != nil {
		return nil, m.findByMACErr
	}
	if m.existingMacID != "" {
		return &FindMachineByMACResponse{
			MachineID: m.existingMacID,
			Found:     true,
		}, nil
	}
	return &FindMachineByMACResponse{Found: false}, nil
}

func (m *mockFirestoreClient) Close() error {
	return nil
}

func TestMockClient(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateMachine success", func(t *testing.T) {
		mock := &mockFirestoreClient{}
		req := &CreateMachineRequest{
			MachineID: "test-id",
			Machine: &MachineRequest{
				NICs: []NIC{{MAC: "aa:bb:cc:dd:ee:ff"}},
			},
		}

		_, err := mock.CreateMachine(ctx, req)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(mock.machines) != 1 {
			t.Errorf("expected 1 machine, got %d", len(mock.machines))
		}
	})

	t.Run("FindMachineByMAC not found", func(t *testing.T) {
		mock := &mockFirestoreClient{}

		resp, err := mock.FindMachineByMAC(ctx, &FindMachineByMACRequest{
			MAC: "aa:bb:cc:dd:ee:ff",
		})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.Found {
			t.Errorf("expected not found, got found with ID %s", resp.MachineID)
		}
	})

	t.Run("FindMachineByMAC found", func(t *testing.T) {
		mock := &mockFirestoreClient{
			existingMacID: "existing-id",
		}

		resp, err := mock.FindMachineByMAC(ctx, &FindMachineByMACRequest{
			MAC: "aa:bb:cc:dd:ee:ff",
		})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !resp.Found {
			t.Error("expected found, got not found")
		}
		if resp.MachineID != "existing-id" {
			t.Errorf("expected ID 'existing-id', got '%s'", resp.MachineID)
		}
	})
}
