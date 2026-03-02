package response

import (
	"net/http"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func SendMessage(w http.ResponseWriter, message string, statusCode int) {
	SendJSON(w, statusCode, MessageResponse{
		Message: message,
	})
}

func SendError(w http.ResponseWriter, statusCode int, code, message string) {
	SendJSON(w, statusCode, ErrorResponse{
		Code:    code,
		Message: message,
	})
}

func SendValidationError(w http.ResponseWriter, statusCode int, message string, fields map[string]string) {
	SendJSON(w, statusCode, ErrorResponse{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Fields:  fields,
	})
}
