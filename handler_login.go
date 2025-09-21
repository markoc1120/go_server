package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/markoc1120/go_server/internal/auth"
	"github.com/markoc1120/go_server/internal/database"
	"github.com/markoc1120/go_server/internal/models"
	"github.com/markoc1120/go_server/internal/response"
	"github.com/markoc1120/go_server/internal/validation"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var params models.LoginRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if err := validation.ValidateEmail(params.Email); err != nil {
		response.WithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	if params.Password == "" {
		response.WithError(w, http.StatusBadRequest, "password is required", nil)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		response.WithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		response.WithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.config.Secret, time.Hour)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Error generating JWT accessToken", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Error generating refreshToken", err)
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	response.WithJSON(w, http.StatusOK, models.LoggedInUser{
		User: models.User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
