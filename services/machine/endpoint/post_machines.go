package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Zaba505/infra/services/machine/errors"
	"github.com/Zaba505/infra/services/machine/models"
	"github.com/google/uuid"
	"github.com/z5labs/humus/rest"
	"github.com/z5labs/humus/rest/rpc"
)

type FirestoreClient interface {
	CreateMachine(ctx context.Context, machineID string, machine *models.MachineRequest) error
	FindMachineByMAC(ctx context.Context, mac string) (string, bool, error)
	Close() error
}

type postMachinesHandler struct {
	firestoreClient FirestoreClient
}

func (h *postMachinesHandler) Handle(ctx context.Context, req *models.MachineRequest) (*models.MachineResponse, error) {
	invalidFields := req.Validate()
	if len(invalidFields) > 0 {
		return nil, errors.NewValidationError("/api/v1/machines", invalidFields)
	}

	for _, nic := range req.NICs {
		existingID, found, err := h.firestoreClient.FindMachineByMAC(ctx, nic.MAC)
		if err != nil {
			return nil, errors.NewInternalError("/api/v1/machines", fmt.Sprintf("failed to check MAC uniqueness: %v", err))
		}
		if found {
			return nil, errors.NewConflictError("/api/v1/machines", nic.MAC, existingID)
		}
	}

	machineID, err := uuid.NewV7()
	if err != nil {
		return nil, errors.NewInternalError("/api/v1/machines", fmt.Sprintf("failed to generate machine ID: %v", err))
	}

	if err := h.firestoreClient.CreateMachine(ctx, machineID.String(), req); err != nil {
		return nil, errors.NewInternalError("/api/v1/machines", fmt.Sprintf("failed to create machine: %v", err))
	}

	return &models.MachineResponse{
		ID: machineID.String(),
	}, nil
}

func errorHandler(ctx context.Context, w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *errors.ValidationProblem:
		e.WriteHttpResponse(ctx, w)
	case *errors.ConflictProblem:
		e.WriteHttpResponse(ctx, w)
	case *errors.Problem:
		e.WriteHttpResponse(ctx, w)
	default:
		genericErr := errors.NewInternalError("", err.Error())
		genericErr.WriteHttpResponse(ctx, w)
	}
}

type responseWithLocation struct {
	resp     *models.MachineResponse
	location string
}

func (r *responseWithLocation) WriteHttpResponse(ctx context.Context, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", r.location)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(r.resp)
}

type postMachinesHandlerWithLocation struct {
	inner *postMachinesHandler
}

func (h *postMachinesHandlerWithLocation) Handle(ctx context.Context, req *models.MachineRequest) (*responseWithLocation, error) {
	resp, err := h.inner.Handle(ctx, req)
	if err != nil {
		return nil, err
	}

	return &responseWithLocation{
		resp:     resp,
		location: fmt.Sprintf("/api/v1/machines/%s", resp.ID),
	}, nil
}

func PostMachines(firestoreClient FirestoreClient) rest.ApiOption {
	handler := &postMachinesHandlerWithLocation{
		inner: &postMachinesHandler{firestoreClient: firestoreClient},
	}

	return rest.Handle(
		http.MethodPost,
		rest.BasePath("/api/v1/machines"),
		rpc.HandleJson(handler),
		rest.OnError(rest.ErrorHandlerFunc(errorHandler)),
	)
}
