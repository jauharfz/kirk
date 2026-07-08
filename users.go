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

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
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
	token, _ := auth.MakeJWT(userDb.ID, cfg.secret, time.Duration(3600))
	refreshToken := auth.MakeRefreshToken()
	userReq = User(userDb)

	_ = cfg.db.InsertRefreshToken(ctx, database.InsertRefreshTokenParams{Token: refreshToken, UserID: userReq.ID, ExpiresAt: time.Now().Add(60 * 24 * time.Second)})

	response := struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}{
		ID:           userReq.ID,
		CreatedAt:    userReq.CreatedAt,
		UpdatedAt:    userReq.UpdatedAt,
		Email:        userReq.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}

	_ = cfg.respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		cfg.respondWithError(w, 400, fmt.Sprintf("cannot decode json request from client %v", err))
		return
	}
	userID, err := verifyIds(r.Header, cfg.secret)
	if err != nil {
		cfg.respondWithError(w, 401, fmt.Sprintf("failed to verify token %v", err))
		return
	}

	password, _ := auth.HashPassword(newUser.Password)

	user, _ := cfg.db.UpdateUserEmailAndPwd(ctx, database.UpdateUserEmailAndPwdParams{Email: newUser.Email, Password: password, ID: userID})

	response := struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	cfg.respondWithJSON(w, http.StatusOK, response)
}
