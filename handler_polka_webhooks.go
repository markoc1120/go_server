package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/markoc1120/go_server/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get API key from header", err)
		return
	}
	if apiKey != cfg.polkaAPIKey {
		respondWithError(w, http.StatusUnauthorized, "You are not allowed to do this", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user_id to uuid.UUID", err)
		return
	}
	if params.Event == "user.upgraded" {
		err = cfg.db.UpdateUserToChirpyRed(r.Context(), userID)
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "User not found", nil)
			return
		}
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't update user to chirpy_red", err)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
