package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zaba505/infra/services/machine/endpoint/endpointpb"
	machineerrors "github.com/Zaba505/infra/services/machine/errors"
	"github.com/Zaba505/infra/services/machine/service"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/protobuf/proto"
)

type mockFirestoreClient struct {
	findResp  *service.FindMachineByMACResponse
	findErr   error
	createErr error
}

func (m *mockFirestoreClient) FindMachineByMAC(_ context.Context, _ *service.FindMachineByMACRequest) (*service.FindMachineByMACResponse, error) {
	return m.findResp, m.findErr
}

func (m *mockFirestoreClient) CreateMachine(_ context.Context, _ *service.CreateMachineRequest) (*service.CreateMachineResponse, error) {
	return &service.CreateMachineResponse{}, m.createErr
}

func (m *mockFirestoreClient) Close() error { return nil }

func TestValidateMACAddress(t *testing.T) {
	tests := []struct {
		name    string
		mac     string
		wantErr bool
	}{
		{
			name:    "valid MAC address lowercase",
			mac:     "52:54:00:12:34:56",
			wantErr: false,
		},
		{
			name:    "valid MAC address uppercase",
			mac:     "AA:BB:CC:DD:EE:FF",
			wantErr: false,
		},
		{
			name:    "valid MAC address mixed case",
			mac:     "aA:bB:cC:dD:eE:fF",
			wantErr: false,
		},
		{
			name:    "empty MAC address",
			mac:     "",
			wantErr: true,
		},
		{
			name:    "invalid format - missing colons",
			mac:     "aabbccddeeff",
			wantErr: true,
		},
		{
			name:    "invalid format - wrong separator",
			mac:     "aa-bb-cc-dd-ee-ff",
			wantErr: true,
		},
		{
			name:    "invalid format - too short",
			mac:     "aa:bb:cc:dd:ee",
			wantErr: true,
		},
		{
			name:    "invalid format - too long",
			mac:     "aa:bb:cc:dd:ee:ff:gg",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			mac:     "zz:yy:xx:ww:vv:uu",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMACAddress(tt.mac)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMACAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterMachinesHandler_ServeHTTP(t *testing.T) {
	mac := "aa:bb:cc:dd:ee:ff"
	validBody, _ := proto.Marshal(&endpointpb.RegisterMachineRequest{
		Nics: []*endpointpb.NIC{{Mac: &mac}},
	})

	tests := []struct {
		name      string
		body      []byte
		client    *mockFirestoreClient
		wantCode  int
		checkBody func(t *testing.T, body []byte)
	}{
		{
			name:     "invalid proto body",
			body:     []byte{0xFF},
			client:   &mockFirestoreClient{},
			wantCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body []byte) {
				var p machineerrors.ValidationProblem
				if err := json.Unmarshal(body, &p); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(p.InvalidFields) == 0 || p.InvalidFields[0].Field != "body" {
					t.Errorf("expected invalid field 'body', got %v", p.InvalidFields)
				}
			},
		},
		{
			name: "no NICs",
			body: func() []byte {
				b, _ := proto.Marshal(&endpointpb.RegisterMachineRequest{})
				return b
			}(),
			client:   &mockFirestoreClient{},
			wantCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body []byte) {
				var p machineerrors.ValidationProblem
				if err := json.Unmarshal(body, &p); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(p.InvalidFields) == 0 || p.InvalidFields[0].Field != "nics" {
					t.Errorf("expected invalid field 'nics', got %v", p.InvalidFields)
				}
			},
		},
		{
			name: "invalid MAC address",
			body: func() []byte {
				badMAC := "not-a-mac"
				b, _ := proto.Marshal(&endpointpb.RegisterMachineRequest{
					Nics: []*endpointpb.NIC{{Mac: &badMAC}},
				})
				return b
			}(),
			client:   &mockFirestoreClient{},
			wantCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body []byte) {
				var p machineerrors.ValidationProblem
				if err := json.Unmarshal(body, &p); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(p.InvalidFields) == 0 || p.InvalidFields[0].Field != "nics[0].mac" {
					t.Errorf("expected invalid field 'nics[0].mac', got %v", p.InvalidFields)
				}
			},
		},
		{
			name:     "FindMachineByMAC error",
			body:     validBody,
			client:   &mockFirestoreClient{findErr: fmt.Errorf("firestore unavailable")},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "duplicate MAC conflict",
			body: validBody,
			client: &mockFirestoreClient{
				findResp: &service.FindMachineByMACResponse{Found: true, MachineID: "existing-id"},
			},
			wantCode: http.StatusConflict,
			checkBody: func(t *testing.T, body []byte) {
				var p machineerrors.ConflictProblem
				if err := json.Unmarshal(body, &p); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if p.ExistingMachineID != "existing-id" {
					t.Errorf("want existing_machine_id 'existing-id', got %q", p.ExistingMachineID)
				}
			},
		},
		{
			name: "CreateMachine error",
			body: validBody,
			client: &mockFirestoreClient{
				findResp:  &service.FindMachineByMACResponse{Found: false},
				createErr: fmt.Errorf("write failed"),
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "success",
			body: validBody,
			client: &mockFirestoreClient{
				findResp: &service.FindMachineByMACResponse{Found: false},
			},
			wantCode: http.StatusCreated,
			checkBody: func(t *testing.T, body []byte) {
				var resp endpointpb.RegisterMachineResponse
				if err := proto.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to decode proto response: %v", err)
				}
				if resp.GetMachineId() == "" {
					t.Error("expected non-empty machine_id")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &registerMachinesHandler{
				tracer:          noop.NewTracerProvider().Tracer(""),
				log:             slog.Default(),
				firestoreClient: tt.client,
			}

			r := httptest.NewRequest(http.MethodPost, "/api/v1/machines", bytes.NewReader(tt.body))
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)

			if w.Code != tt.wantCode {
				t.Errorf("want status %d, got %d (body: %s)", tt.wantCode, w.Code, w.Body.String())
			}
			if tt.checkBody != nil {
				tt.checkBody(t, w.Body.Bytes())
			}
		})
	}
}
