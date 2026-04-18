package endpoint

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/Zaba505/infra/pkg/errorpb"
	"github.com/Zaba505/infra/services/machine/endpoint/endpointpb"
	"github.com/Zaba505/infra/services/machine/service"
	"github.com/go-chi/chi/v5"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

type FirestoreClient interface {
	CreateMachine(ctx context.Context, req *service.CreateMachineRequest) (*service.CreateMachineResponse, error)
	FindMachineByMAC(ctx context.Context, req *service.FindMachineByMACRequest) (*service.FindMachineByMACResponse, error)
	Close() error
}

type registerMachinesHandler struct {
	tracer          trace.Tracer
	log             *slog.Logger
	firestoreClient FirestoreClient
}

func RegisterMachines(mux *chi.Mux, firestoreClient FirestoreClient) {
	handler := &registerMachinesHandler{
		tracer:          otel.Tracer("machine/endpoint"),
		log:             slog.Default(),
		firestoreClient: firestoreClient,
	}

	mux.Method(http.MethodPost, "/api/v1/machines", handler)
}

func (h *registerMachinesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorHandler(ctx, w, errorpb.NewInternalError("/api/v1/machines", fmt.Sprintf("failed to read request body: %v", err)))
		return
	}

	var req endpointpb.RegisterMachineRequest
	if err := proto.Unmarshal(body, &req); err != nil {
		errorHandler(ctx, w, errorpb.NewValidationError("/api/v1/machines", []*errorpb.InvalidField{
			{Field: proto.String("body"), Reason: proto.String(fmt.Sprintf("invalid protobuf: %v", err))},
		}))
		return
	}

	nics := req.GetNics()
	if len(nics) == 0 {
		errorHandler(ctx, w, errorpb.NewValidationError("/api/v1/machines", []*errorpb.InvalidField{
			{Field: proto.String("nics"), Reason: proto.String("at least one NIC is required")},
		}))
		return
	}

	var invalidFields []*errorpb.InvalidField
	for i, nic := range nics {
		if err := validateMACAddress(nic.GetMac()); err != nil {
			invalidFields = append(invalidFields, &errorpb.InvalidField{
				Field:  proto.String(fmt.Sprintf("nics[%d].mac", i)),
				Reason: proto.String(err.Error()),
			})
		}
	}
	if len(invalidFields) > 0 {
		errorHandler(ctx, w, errorpb.NewValidationError("/api/v1/machines", invalidFields))
		return
	}

	for _, nic := range nics {
		resp, err := h.firestoreClient.FindMachineByMAC(ctx, &service.FindMachineByMACRequest{
			MAC: nic.GetMac(),
		})
		if err != nil {
			errorHandler(ctx, w, errorpb.NewInternalError("/api/v1/machines", fmt.Sprintf("failed to check MAC uniqueness: %v", err)))
			return
		}
		if resp.Found {
			errorHandler(ctx, w, errorpb.NewConflictError("/api/v1/machines", resp.MachineID, map[string]string{"mac_address": nic.GetMac()}))
			return
		}
	}

	machineID, err := uuid.NewV7()
	if err != nil {
		errorHandler(ctx, w, errorpb.NewInternalError("/api/v1/machines", fmt.Sprintf("failed to generate machine ID: %v", err)))
		return
	}

	serviceReq := &service.MachineRequest{
		CPUs:          convertCPUs(req.GetCpus()),
		MemoryModules: convertMemoryModules(req.GetMemoryModules()),
		Accelerators:  convertAccelerators(req.GetAccelerators()),
		NICs:          convertNICs(nics),
		Drives:        convertDrives(req.GetDrives()),
	}

	_, err = h.firestoreClient.CreateMachine(ctx, &service.CreateMachineRequest{
		MachineID: machineID.String(),
		Machine:   serviceReq,
	})
	if err != nil {
		errorHandler(ctx, w, errorpb.NewInternalError("/api/v1/machines", fmt.Sprintf("failed to create machine: %v", err)))
		return
	}

	machineIDStr := machineID.String()
	resp := &endpointpb.RegisterMachineResponse{
		MachineId: &machineIDStr,
	}
	respBody, err := proto.Marshal(resp)
	if err != nil {
		errorHandler(ctx, w, errorpb.NewInternalError("/api/v1/machines", fmt.Sprintf("failed to marshal response: %v", err)))
		return
	}

	w.Header().Set("Content-Type", "application/x-protobuf")
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}

var macAddressRegex = regexp.MustCompile(`^([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$`)

func validateMACAddress(mac string) error {
	if mac == "" {
		return fmt.Errorf("MAC address cannot be empty")
	}
	if !macAddressRegex.MatchString(mac) {
		return fmt.Errorf("invalid MAC address format, expected format: aa:bb:cc:dd:ee:ff")
	}
	return nil
}

func convertCPUs(cpus []*endpointpb.CPU) []service.CPU {
	result := make([]service.CPU, len(cpus))
	for i, cpu := range cpus {
		result[i] = service.CPU{
			Manufacturer:   cpu.GetManufacturer(),
			ClockFrequency: cpu.GetClockFrequency(),
			Cores:          cpu.GetCores(),
		}
	}
	return result
}

func convertMemoryModules(modules []*endpointpb.MemoryModule) []service.MemoryModule {
	result := make([]service.MemoryModule, len(modules))
	for i, module := range modules {
		result[i] = service.MemoryModule{
			Size: module.GetSize(),
		}
	}
	return result
}

func convertAccelerators(accelerators []*endpointpb.Accelerator) []service.Accelerator {
	result := make([]service.Accelerator, len(accelerators))
	for i, accelerator := range accelerators {
		result[i] = service.Accelerator{
			Manufacturer: accelerator.GetManufacturer(),
		}
	}
	return result
}

func convertNICs(nics []*endpointpb.NIC) []service.NIC {
	result := make([]service.NIC, len(nics))
	for i, nic := range nics {
		result[i] = service.NIC{
			MAC: nic.GetMac(),
		}
	}
	return result
}

func convertDrives(drives []*endpointpb.Drive) []service.Drive {
	result := make([]service.Drive, len(drives))
	for i, drive := range drives {
		result[i] = service.Drive{
			Capacity: drive.GetCapacity(),
		}
	}
	return result
}

func errorHandler(ctx context.Context, w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *errorpb.ValidationProblem:
		e.WriteHttpResponse(ctx, w)
	case *errorpb.ConflictProblem:
		e.WriteHttpResponse(ctx, w)
	case *errorpb.Problem:
		e.WriteHttpResponse(ctx, w)
	default:
		genericErr := errorpb.NewInternalError("", err.Error())
		genericErr.WriteHttpResponse(ctx, w)
	}
}
