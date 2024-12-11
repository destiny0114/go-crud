package domain

import (
	"encoding/json"
	"net/http"
)

type SimpleResponse struct {
	Message string `json:"message"`
}

func HelloHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := SimpleResponse{Message: "API worked!"}
	json.NewEncoder(w).Encode(res)
}
