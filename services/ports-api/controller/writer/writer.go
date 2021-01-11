package writer

import (
	"encoding/json"
	"net/http"
)

func Write(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data) // insted of `_`, an error var should be used to log it
}
