package main

import (
	"context"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if cfg.platform != "dev" {
		cfg.respondWithError(w, http.StatusForbidden, "failed to reset")
		return
	}
	cfg.db.Reset(ctx)
	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
