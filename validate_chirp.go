package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	badWords := map[string]struct{}{
		"sharbert":  {},
		"fornax":    {},
		"kerfuffle": {},
	}
	resBody := getCleanedBody(req.Body, badWords)

	response := map[string]interface{}{
		"valid":        true,
		"cleaned_body": resBody,
	}

	cfg.respondWithJSON(w, 200, response)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	splitBody := strings.Split(body, " ")
	for i := 0; i < len(splitBody); i++ {
		lowerStr := strings.ToLower(splitBody[i])
		if _, ok := badWords[lowerStr]; ok {
			splitBody[i] = "****"
		}
	}
	return strings.Join(splitBody, " ")
}
