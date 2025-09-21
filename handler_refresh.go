package main

import (
	"net/http"
	"time"

	"github.com/markoc1120/go_server/internal/auth"
	"github.com/markoc1120/go_server/internal/models"
	"github.com/markoc1120/go_server/internal/response"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		response.WithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		response.WithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.config.Secret, time.Hour)
	if err != nil {
		response.WithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}
	response.WithJSON(w, http.StatusOK, models.TokenResponse{Token: accessToken})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		response.WithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}
	err = cfg.db.RevokeRefreshTokenByToken(r.Context(), refreshToken)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't revoke session", err)
		return
	}
	response.WithJSON(w, http.StatusNoContent, nil)
}
