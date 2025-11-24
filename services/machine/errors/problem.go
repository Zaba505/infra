package errors

import (
	"context"
	"encoding/json"
	"net/http"
)

type InvalidField struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func (p *Problem) Error() string {
	return p.Detail
}

func (p *Problem) WriteHttpResponse(ctx context.Context, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)
	json.NewEncoder(w).Encode(p)
}

type ValidationProblem struct {
	Problem
	InvalidFields []InvalidField `json:"invalid_fields"`
}

func (vp *ValidationProblem) WriteHttpResponse(ctx context.Context, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(vp.Status)
	json.NewEncoder(w).Encode(vp)
}

type ConflictProblem struct {
	Problem
	MACAddress        string `json:"mac_address"`
	ExistingMachineID string `json:"existing_machine_id"`
}

func (cp *ConflictProblem) WriteHttpResponse(ctx context.Context, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(cp.Status)
	json.NewEncoder(w).Encode(cp)
}

func NewValidationError(instance string, fields []InvalidField) *ValidationProblem {
	return &ValidationProblem{
		Problem: Problem{
			Type:     "https://api.example.com/errors/validation-error",
			Title:    "Validation Error",
			Status:   http.StatusBadRequest,
			Detail:   "The request body failed validation",
			Instance: instance,
		},
		InvalidFields: fields,
	}
}

func NewConflictError(instance, mac, existingID string) *ConflictProblem {
	return &ConflictProblem{
		Problem: Problem{
			Type:     "https://api.example.com/errors/duplicate-mac-address",
			Title:    "Duplicate MAC Address",
			Status:   http.StatusConflict,
			Detail:   "A machine with MAC address " + mac + " already exists",
			Instance: instance,
		},
		MACAddress:        mac,
		ExistingMachineID: existingID,
	}
}

func NewInternalError(instance, detail string) *Problem {
	return &Problem{
		Type:     "https://api.example.com/errors/internal-error",
		Title:    "Internal Server Error",
		Status:   http.StatusInternalServerError,
		Detail:   detail,
		Instance: instance,
	}
}
