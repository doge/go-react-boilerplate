package response

import (
	"encoding/json"
	"net/http"
)

func SendMessage(w http.ResponseWriter, message string, statusCode int) {

	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(map[string]string{
		"status": message,
	})
}
