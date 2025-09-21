package main

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/markoc1120/go_server/internal/auth"
	"github.com/markoc1120/go_server/internal/response"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		response.WithError(w, http.StatusUnauthorized, "Couldn't find token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.config.Secret)
	if err != nil {
		response.WithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	query := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(query)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Invalid chirp ID in the url", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.WithError(w, http.StatusNotFound, "chirp not found", nil)
			return
		}
		response.WithError(w, http.StatusInternalServerError, "Couldn't retrieve the single chirp instance from db", err)
		return
	}

	if chirp.UserID != userID {
		response.WithError(w, http.StatusForbidden, "You can't do this", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't delete the chirp instance from db", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
