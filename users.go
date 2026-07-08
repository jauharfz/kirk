package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jauharfz/GO-CHRPY/internal/auth"
	"github.com/jauharfz/GO-CHRPY/internal/database"
)

type User struct {
	ID               uuid.UUID `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	ExpiresInSeconds int32     `json:"expires_in_seconds,omitempty"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user User

	_ = json.NewDecoder(r.Body).Decode(&user)
	hashedPwd, _ := auth.HashPassword(user.Password)
	user.Password = ""
	userRes, _ := cfg.db.CreateUser(ctx, database.CreateUserParams{Email: user.Email, Password: hashedPwd})
	user = User(userRes)

	_ = cfg.respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var userReq User

	_ = json.NewDecoder(r.Body).Decode(&userReq)
	userDb, _ := cfg.db.GetUsersByEmail(ctx, userReq.Email)
	ok, _ := auth.CheckPasswordHash(userReq.Password, userDb.Password)
	if !ok {
		cfg.respondWithError(w, 401, "email or password didn't match")
	}
	if userDb.ExpiresInSeconds > 3600 || userDb.ExpiresInSeconds <= 0 {
		userDb.ExpiresInSeconds = 3600
	}
	token, _ := auth.MakeJWT(userDb.ID, cfg.secret, time.Duration(userDb.ExpiresInSeconds))
	userReq = User(userDb)
	response := struct {
		ID               uuid.UUID `json:"id"`
		CreatedAt        time.Time `json:"created_at"`
		UpdatedAt        time.Time `json:"updated_at"`
		Email            string    `json:"email"`
		ExpiresInSeconds int32     `json:"expires_in_seconds,omitempty"`
		Token            string    `json:"token"`
	}{
		ID:               userReq.ID,
		CreatedAt:        userReq.CreatedAt,
		UpdatedAt:        userReq.UpdatedAt,
		Email:            userReq.Email,
		ExpiresInSeconds: userReq.ExpiresInSeconds,
		Token:            token,
	}

	_ = cfg.respondWithJSON(w, http.StatusOK, response)
}
