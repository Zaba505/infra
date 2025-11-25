package endpoint

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Zaba505/infra/services/machine/errors"
	"github.com/Zaba505/infra/services/machine/service"
	"github.com/google/uuid"
	"github.com/z5labs/humus"
	"github.com/z5labs/humus/rest"
	"github.com/z5labs/humus/rest/rpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type FirestoreClient interface {
	CreateMachine(ctx context.Context, req *service.CreateMachineRequest) (*service.CreateMachineResponse, error)
	FindMachineByMAC(ctx context.Context, req *service.FindMachineByMACRequest) (*service.FindMachineByMACResponse, error)
	Close() error
}

type postMachinesHandler struct {
	tracer          trace.Tracer
	log             *slog.Logger
	firestoreClient FirestoreClient
}

func (h *postMachinesHandler) Handle(ctx context.Context, req *MachineRequest) (*MachineResponse, error) {
	invalidFields := req.Validate()
	if len(invalidFields) > 0 {
		errFields := make([]errors.InvalidField, len(invalidFields))
		for i, f := range invalidFields {
			errFields[i] = errors.InvalidField{
				Field:  f.Field,
				Reason: f.Reason,
			}
		}
		return nil, errors.NewValidationError("/api/v1/machines", errFields)
	}

	for _, nic := range req.NICs {
		resp, err := h.firestoreClient.FindMachineByMAC(ctx, &service.FindMachineByMACRequest{
			MAC: nic.MAC,
		})
		if err != nil {
			return nil, errors.NewInternalError("/api/v1/machines", fmt.Sprintf("failed to check MAC uniqueness: %v", err))
		}
		if resp.Found {
			return nil, errors.NewConflictError("/api/v1/machines", nic.MAC, resp.MachineID)
		}
	}

	machineID, err := uuid.NewV7()
	if err != nil {
		return nil, errors.NewInternalError("/api/v1/machines", fmt.Sprintf("failed to generate machine ID: %v", err))
	}

	serviceReq := &service.MachineRequest{
		CPUs:          convertCPUs(req.CPUs),
		MemoryModules: convertMemoryModules(req.MemoryModules),
		Accelerators:  convertAccelerators(req.Accelerators),
		NICs:          convertNICs(req.NICs),
		Drives:        convertDrives(req.Drives),
	}

	_, err = h.firestoreClient.CreateMachine(ctx, &service.CreateMachineRequest{
		MachineID: machineID.String(),
		Machine:   serviceReq,
	})
	if err != nil {
		return nil, errors.NewInternalError("/api/v1/machines", fmt.Sprintf("failed to create machine: %v", err))
	}

	return &MachineResponse{
		ID: machineID.String(),
	}, nil
}

func convertCPUs(cpus []CPU) []service.CPU {
	result := make([]service.CPU, len(cpus))
	for i, cpu := range cpus {
		result[i] = service.CPU{
			Manufacturer:   cpu.Manufacturer,
			ClockFrequency: cpu.ClockFrequency,
			Cores:          cpu.Cores,
		}
	}
	return result
}

func convertMemoryModules(modules []MemoryModule) []service.MemoryModule {
	result := make([]service.MemoryModule, len(modules))
	for i, module := range modules {
		result[i] = service.MemoryModule{
			Size: module.Size,
		}
	}
	return result
}

func convertAccelerators(accelerators []Accelerator) []service.Accelerator {
	result := make([]service.Accelerator, len(accelerators))
	for i, accelerator := range accelerators {
		result[i] = service.Accelerator{
			Manufacturer: accelerator.Manufacturer,
		}
	}
	return result
}

func convertNICs(nics []NIC) []service.NIC {
	result := make([]service.NIC, len(nics))
	for i, nic := range nics {
		result[i] = service.NIC{
			MAC: nic.MAC,
		}
	}
	return result
}

func convertDrives(drives []Drive) []service.Drive {
	result := make([]service.Drive, len(drives))
	for i, drive := range drives {
		result[i] = service.Drive{
			Capacity: drive.Capacity,
		}
	}
	return result
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

func PostMachines(firestoreClient FirestoreClient) rest.ApiOption {
	handler := &postMachinesHandler{
		tracer:          otel.Tracer("machine/endpoint"),
		log:             humus.Logger("machine/endpoint"),
		firestoreClient: firestoreClient,
	}

	return rest.Handle(
		http.MethodPost,
		rest.BasePath("/api/v1/machines"),
		rpc.HandleJson(handler),
		rest.OnError(rest.ErrorHandlerFunc(errorHandler)),
	)
}
