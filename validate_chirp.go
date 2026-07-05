package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Body string `json:"body"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		cfg.respondWithError(w, 400, fmt.Sprintf("cannot decode json request from client %v", err))
		return
	}

	if len(req.Body) > 100 {
		cfg.respondWithError(w, 400, fmt.Sprintf("request body is too long %v", err))
		return

	}

	response := map[string]bool{
		"valid": true,
	}

	cfg.respondWithJSON(w, 200, response)
}
