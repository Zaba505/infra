package errorpb

import (
	"context"
	"net/http"

	"google.golang.org/protobuf/proto"
)

func (p *Problem) Error() string {
	return p.GetDetail()
}

func (p *Problem) WriteHttpResponse(_ context.Context, w http.ResponseWriter) {
	writeProtoError(w, int(p.GetStatus()), p)
}

func (vp *ValidationProblem) Error() string {
	return vp.GetProblem().GetDetail()
}

func (vp *ValidationProblem) WriteHttpResponse(_ context.Context, w http.ResponseWriter) {
	writeProtoError(w, int(vp.GetProblem().GetStatus()), vp)
}

func (cp *ConflictProblem) Error() string {
	return cp.GetProblem().GetDetail()
}

func (cp *ConflictProblem) WriteHttpResponse(_ context.Context, w http.ResponseWriter) {
	writeProtoError(w, int(cp.GetProblem().GetStatus()), cp)
}

func writeProtoError(w http.ResponseWriter, status int, msg proto.Message) {
	b, err := proto.Marshal(msg)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/problem+protobuf")
	w.WriteHeader(status)
	w.Write(b)
}

func NewValidationError(instance string, fields []*InvalidField) *ValidationProblem {
	return &ValidationProblem{
		Problem: &Problem{
			Type:     proto.String("https://api.example.com/errors/validation-error"),
			Title:    proto.String("Validation Error"),
			Status:   proto.Int32(http.StatusBadRequest),
			Detail:   proto.String("The request body failed validation"),
			Instance: proto.String(instance),
		},
		InvalidFields: fields,
	}
}

func NewConflictError(instance, existingResourceID string, conflictingFields map[string]string) *ConflictProblem {
	return &ConflictProblem{
		Problem: &Problem{
			Type:     proto.String("https://api.example.com/errors/conflict"),
			Title:    proto.String("Conflict"),
			Status:   proto.Int32(http.StatusConflict),
			Detail:   proto.String("A resource with the given identifier already exists"),
			Instance: proto.String(instance),
		},
		ExistingResourceId: proto.String(existingResourceID),
		ConflictingFields:  conflictingFields,
	}
}

func NewInternalError(instance, detail string) *Problem {
	return &Problem{
		Type:     proto.String("https://api.example.com/errors/internal-error"),
		Title:    proto.String("Internal Server Error"),
		Status:   proto.Int32(http.StatusInternalServerError),
		Detail:   proto.String(detail),
		Instance: proto.String(instance),
	}
}
