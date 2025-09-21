package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/markoc1120/go_server/internal/auth"
	"github.com/markoc1120/go_server/internal/database"
	"github.com/markoc1120/go_server/internal/models"
	"github.com/markoc1120/go_server/internal/response"
)

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		response.WithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.config.Secret)
	if err != nil {
		response.WithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := models.CreateChirpRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't decode params", err)
		return
	}

	cleanedBody, err := validateChirp(params.Body)
	if err != nil {
		response.WithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}
	response.WithJSON(w, http.StatusCreated, models.Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func validateChirp(chirp string) (string, error) {
	const maxChirpLength = 140
	if len(chirp) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}
	return cleanBody(chirp), nil
}

func getProfaneWords() []string {
	return []string{"kerfuffle", "sharbert", "fornax"}
}

func cleanBody(body string) string {
	profaneWords := getProfaneWords()
	var result strings.Builder
	for w := range strings.SplitSeq(body, " ") {
		if slices.Contains(profaneWords, strings.ToLower(w)) {
			result.WriteString("**** ")
			continue
		}
		result.WriteString(w + " ")
	}

	return strings.TrimSpace(result.String())
}
