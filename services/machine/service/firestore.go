package service

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type CreateMachineRequest struct {
	MachineID string
	Machine   *MachineRequest
}

type CreateMachineResponse struct{}

type FindMachineByMACRequest struct {
	MAC string
}

type FindMachineByMACResponse struct {
	MachineID string
	Found     bool
}

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

func (c *FirestoreClient) CreateMachine(ctx context.Context, req *CreateMachineRequest) (*CreateMachineResponse, error) {
	docRef := c.client.Collection("machines").Doc(req.MachineID)

	data := map[string]interface{}{
		"id":             req.MachineID,
		"cpus":           req.Machine.CPUs,
		"memory_modules": req.Machine.MemoryModules,
		"accelerators":   req.Machine.Accelerators,
		"nics":           req.Machine.NICs,
		"drives":         req.Machine.Drives,
	}

	_, err := docRef.Set(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create machine document: %w", err)
	}

	return &CreateMachineResponse{}, nil
}

func (c *FirestoreClient) FindMachineByMAC(ctx context.Context, req *FindMachineByMACRequest) (*FindMachineByMACResponse, error) {
	normalizedMAC := strings.ToLower(req.MAC)

	iter := c.client.Collection("machines").
		Where("nics", "array-contains", map[string]interface{}{"mac": normalizedMAC}).
		Limit(1).
		Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return &FindMachineByMACResponse{Found: false}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query machines by MAC: %w", err)
	}

	var data struct {
		ID   string `firestore:"id"`
		NICs []NIC  `firestore:"nics"`
	}
	if err := doc.DataTo(&data); err != nil {
		return nil, fmt.Errorf("failed to decode machine document: %w", err)
	}

	for _, nic := range data.NICs {
		if strings.EqualFold(nic.MAC, req.MAC) {
			return &FindMachineByMACResponse{
				MachineID: data.ID,
				Found:     true,
			}, nil
		}
	}

	return &FindMachineByMACResponse{Found: false}, nil
}

func (c *FirestoreClient) Close() error {
	return c.client.Close()
}
