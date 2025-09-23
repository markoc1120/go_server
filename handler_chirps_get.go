package main

import (
	"context"
	"database/sql"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/markoc1120/go_server/internal/database"
	"github.com/markoc1120/go_server/internal/models"
	"github.com/markoc1120/go_server/internal/response"
)

func getChirpInstances(ctx context.Context, userID *uuid.UUID, db *database.Queries) ([]database.Chirp, error) {
	if userID != nil {
		return db.GetChirpsByUserID(ctx, *userID)
	}
	return db.GetChirps(ctx)
}

func getSortedChirps(chirps []database.Chirp, sortType string) []database.Chirp {
	if sortType == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	}
	return chirps
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var userID *uuid.UUID
	if authorID := r.URL.Query().Get("author_id"); authorID != "" {
		id, err := uuid.Parse(authorID)
		if err != nil {
			response.WithError(w, http.StatusBadRequest, "Invalid author_id query parameter", err)
			return
		}
		userID = &id
	}

	chirps, err := getChirpInstances(ctx, userID, cfg.db)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't retrieve all instances from db", err)
		return
	}
	sortType := r.URL.Query().Get("sort")
	sortedChirps := getSortedChirps(chirps, sortType)

	payload := []models.Chirp{}
	for _, chirp := range sortedChirps {
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
		response.WithError(w, http.StatusBadRequest, "Invalid chirpID in the url", err)
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
