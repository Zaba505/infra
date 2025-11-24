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

func (m *mockFirestoreClient) CreateMachine(ctx context.Context, machineID string, machine *MachineRequest) error {
	if m.createErr != nil {
		return m.createErr
	}
	if m.machines == nil {
		m.machines = make(map[string]*MachineRequest)
	}
	m.machines[machineID] = machine
	return nil
}

func (m *mockFirestoreClient) FindMachineByMAC(ctx context.Context, mac string) (string, bool, error) {
	if m.findByMACErr != nil {
		return "", false, m.findByMACErr
	}
	if m.existingMacID != "" {
		return m.existingMacID, true, nil
	}
	return "", false, nil
}

func (m *mockFirestoreClient) Close() error {
	return nil
}

func TestMockClient(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateMachine success", func(t *testing.T) {
		mock := &mockFirestoreClient{}
		req := &MachineRequest{
			NICs: []NIC{{MAC: "aa:bb:cc:dd:ee:ff"}},
		}

		err := mock.CreateMachine(ctx, "test-id", req)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(mock.machines) != 1 {
			t.Errorf("expected 1 machine, got %d", len(mock.machines))
		}
	})

	t.Run("FindMachineByMAC not found", func(t *testing.T) {
		mock := &mockFirestoreClient{}

		id, found, err := mock.FindMachineByMAC(ctx, "aa:bb:cc:dd:ee:ff")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if found {
			t.Errorf("expected not found, got found with ID %s", id)
		}
	})

	t.Run("FindMachineByMAC found", func(t *testing.T) {
		mock := &mockFirestoreClient{
			existingMacID: "existing-id",
		}

		id, found, err := mock.FindMachineByMAC(ctx, "aa:bb:cc:dd:ee:ff")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !found {
			t.Error("expected found, got not found")
		}
		if id != "existing-id" {
			t.Errorf("expected ID 'existing-id', got '%s'", id)
		}
	})
}
