package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jauharfz/GO-CHRPY/internal/auth"
	"github.com/jauharfz/GO-CHRPY/internal/database"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	apiKey, _ := auth.GetApiKey(r.Header)
	if apiKey != cfg.polkaKey {
		cfg.respondWithError(w, http.StatusUnauthorized, "invalid headers")
	}

	var request struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	_ = json.NewDecoder(r.Body).Decode(&request)
	if request.Event != "user.upgraded" {
		cfg.respondWithJSON(w, http.StatusNoContent, "")
		return
	}
	id, _ := uuid.Parse(request.Data.UserID)
	userDB, err := cfg.db.UpdateChirpyRedStatus(ctx, database.UpdateChirpyRedStatusParams{IsChirpyRed: true, ID: id})
	if err != nil {
		cfg.respondWithError(w, http.StatusNotFound, fmt.Sprintf("failed to update chirpy status %v", err))
	}

	user := User(userDB)
	_ = cfg.respondWithJSON(w, http.StatusNoContent, user)

}
