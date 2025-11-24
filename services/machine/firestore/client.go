package firestore

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/Zaba505/infra/services/machine/models"
	"google.golang.org/api/iterator"
)

type Client struct {
	client *firestore.Client
}

func NewClient(ctx context.Context, projectID string) (*Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %w", err)
	}
	return &Client{client: client}, nil
}

func (c *Client) CreateMachine(ctx context.Context, machineID string, machine *models.MachineRequest) error {
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

func (c *Client) FindMachineByMAC(ctx context.Context, mac string) (string, bool, error) {
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
		ID   string       `firestore:"id"`
		NICs []models.NIC `firestore:"nics"`
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

func (c *Client) Close() error {
	return c.client.Close()
}
