package writer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResponseError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Err        string `json:"error"`
}

func NewResponseError(statusCode int, message string, err error) *ResponseError {
	re := &ResponseError{
		Message: message,
	}
	if err != nil {
		re.Err = err.Error()
	}
	return re
}

func Write(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data) // insted of `_`, an error var should be used to log it
}

func WriteError(w http.ResponseWriter, re *ResponseError) {
	Write(w, re.StatusCode, re)
}

func ConvertGrpcError(entityName string, err error) *ResponseError {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.NotFound:
			return NewResponseError(http.StatusNotFound, fmt.Sprintf("%s not found", entityName), err)
		}
		// add here additional error codes
	}
	return NewResponseError(http.StatusInternalServerError, "Internal server error", err)
}
