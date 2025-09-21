package main

import (
	"encoding/json"
	"net/http"

	"github.com/markoc1120/go_server/internal/auth"
	"github.com/markoc1120/go_server/internal/database"
	"github.com/markoc1120/go_server/internal/models"
	"github.com/markoc1120/go_server/internal/response"
	"github.com/markoc1120/go_server/internal/validation"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
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

	decoder := json.NewDecoder(r.Body)
	params := models.UpdateUserRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if err := validation.ValidateEmail(params.Email); err != nil {
		response.WithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	if err := validation.ValidatePassword(params.Password); err != nil {
		response.WithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	passwordHash, err := auth.HashPassword(params.Password)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Error during hashing password", err)
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		HashedPassword: passwordHash,
		ID:             userID,
		Email:          params.Email,
	})
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't update password", err)
		return
	}

	response.WithJSON(w, http.StatusOK, models.User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
