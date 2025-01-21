package main

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, data any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}

func ReadJson(r *http.Request, data any) error {
	return json.NewDecoder(r.Body).Decode(data)
}
