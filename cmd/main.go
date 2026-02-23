package main

import (
	"encoding/json"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	res := map[string]any{
		"success": true,
		"message": "healthy",
	}

	json.NewEncoder(w).Encode(res)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/validate", Validate)
	mux.HandleFunc("POST /api/mutate", Mutate)
	mux.HandleFunc("GET /api/healthz", health)

	err := http.ListenAndServe("0.0.0.0:8000", mux)
	if err != nil {
		panic(err)
	}
}
