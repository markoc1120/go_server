package main

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/markoc1120/go_server/internal/models"
	"github.com/markoc1120/go_server/internal/response"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't retrieve all instances from db", err)
		return
	}

	payload := []models.Chirp{}
	for _, chirp := range chirps {
		payload = append(payload, models.Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	response.WithJSON(w, http.StatusOK, payload)
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	query := r.PathValue("chirpID")
	id, err := uuid.Parse(query)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Invalid chirp ID in the url", err)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.WithError(w, http.StatusNotFound, "chirp not found", nil)
			return
		}
		response.WithError(w, http.StatusInternalServerError, "Couldn't retrieve the single chirp instance from db", err)
		return
	}
	response.WithJSON(w, http.StatusOK, models.Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
