package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jauharfz/GO-CHRPY/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var chirp Chirp
	err := json.NewDecoder(r.Body).Decode(&chirp)
	if err != nil {
		cfg.respondWithError(w, 400, fmt.Sprintf("cannot decode json request from client %v", err))
		return
	}

	chirp.Body, _ = validateChirp(chirp.Body)
	cfg.db.CreateChirp(ctx, database.CreateChirpParams{Body: chirp.Body, UserID: chirp.UserID})

	cfg.respondWithJSON(w, 201, chirp)
}

func validateChirp(body string) (string, error) {
	if len(body) > 100 {
		return "", errors.New("body string is too long")
	}
	badWords := map[string]struct{}{
		"sharbert":  {},
		"fornax":    {},
		"kerfuffle": {},
	}
	return getCleanedBody(body, badWords), nil

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
