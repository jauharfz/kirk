package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jauharfz/GO-CHRPY/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	token, _ := auth.GetBearerToken(r.Header)
	refreshToken, _ := cfg.db.GetRefreshTokenByToken(ctx, token)
	if refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.RevokeAt.Valid {
		cfg.respondWithError(w, 401, "token already expired or revoked")
		return
	}
	token, _ = auth.MakeJWT(refreshToken.UserID, cfg.secret, 3600)
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	cfg.respondWithJSON(w, 200, response)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	token, _ := auth.GetBearerToken(r.Header)
	refreshToken, _ := cfg.db.GetRefreshTokenByToken(ctx, token)
	err := cfg.db.RevokeRefreshToken(ctx, refreshToken.Token)
	if err != nil {
		cfg.respondWithError(w, 401, fmt.Sprintf("failed to revoke refresh token %v", err))
		return
	}
	refreshToken, _ = cfg.db.GetRefreshTokenByToken(ctx, refreshToken.Token)
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", refreshToken.Token))
	cfg.respondWithJSON(w, 204, "")
}
