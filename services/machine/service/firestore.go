package service

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type FirestoreClient struct {
	client *firestore.Client
}

func NewFirestoreClient(ctx context.Context, projectID string) (*FirestoreClient, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %w", err)
	}
	return &FirestoreClient{client: client}, nil
}

func (c *FirestoreClient) CreateMachine(ctx context.Context, machineID string, machine *MachineRequest) error {
	docRef := c.client.Collection("machines").Doc(machineID)

	data := map[string]interface{}{
		"id":             machineID,
		"cpus":           machine.CPUs,
		"memory_modules": machine.MemoryModules,
		"accelerators":   machine.Accelerators,
		"nics":           machine.NICs,
		"drives":         machine.Drives,
	}

	_, err := docRef.Set(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to create machine document: %w", err)
	}

	return nil
}

func (c *FirestoreClient) FindMachineByMAC(ctx context.Context, mac string) (string, bool, error) {
	normalizedMAC := strings.ToLower(mac)

	iter := c.client.Collection("machines").
		Where("nics", "array-contains", map[string]interface{}{"mac": normalizedMAC}).
		Limit(1).
		Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("failed to query machines by MAC: %w", err)
	}

	var data struct {
		ID   string `firestore:"id"`
		NICs []NIC  `firestore:"nics"`
	}
	if err := doc.DataTo(&data); err != nil {
		return "", false, fmt.Errorf("failed to decode machine document: %w", err)
	}

	for _, nic := range data.NICs {
		if strings.EqualFold(nic.MAC, mac) {
			return data.ID, true, nil
		}
	}

	return "", false, nil
}

func (c *FirestoreClient) Close() error {
	return c.client.Close()
}
