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
	"github.com/jauharfz/GO-CHRPY/internal/auth"
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
	userID, err := verifyIds(r.Header, cfg.secret)
	if err != nil {
		cfg.respondWithError(w, 401, fmt.Sprintf("failed to verify token %v", err))
		return
	}

	chirp.Body, err = validateChirp(chirp.Body)
	if err != nil {
		cfg.respondWithError(w, 400, "failed to validate chirp body")
		return
	}
	chirpReq, _ := cfg.db.CreateChirp(ctx, database.CreateChirpParams{Body: chirp.Body, UserID: userID})

	cfg.respondWithJSON(w, 201, Chirp(chirpReq))
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	chirps, _ := cfg.db.GetAllChirps(ctx)
	chirpRes := make([]Chirp, len(chirps))
	for idx := range chirps {
		chirpRes[idx] = Chirp(chirps[idx])
	}

	_ = cfg.respondWithJSON(w, 200, chirpRes)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)

	chirp, err := cfg.db.GetChirpsById(ctx, id)
	if err != nil {
		cfg.respondWithError(w, 404, "failed to get chirps id")
	}

	_ = cfg.respondWithJSON(w, 200, Chirp(chirp))
}

func validateChirp(body string) (string, error) {
	if len(body) > 1000 {
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

func verifyIds(header http.Header, tokenSecret string) (uuid.UUID, error) {
	token, err := auth.GetBearerToken(header)
	if err != nil {
		return uuid.Nil, err
	}
	return auth.VerifyJWT(token, tokenSecret)
}
