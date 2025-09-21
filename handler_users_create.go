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

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	var params models.CreateUserRequest

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

	if err := validation.ValidatePassword(params.Password); err != nil {
		response.WithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	passwordHash, err := auth.HashPassword(params.Password)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Error during hashing password", err)
		return
	}

	user, err := cfg.db.CreateUser(
		r.Context(),
		database.CreateUserParams{Email: params.Email, HashedPassword: passwordHash},
	)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	response.WithJSON(w, http.StatusCreated, models.User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
