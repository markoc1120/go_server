package main

import (
	"net/http"

	"github.com/markoc1120/go_server/internal/response"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.config.Platform != "dev" {
		response.WithError(w, http.StatusForbidden, "You can't do this, reset is only allowed in dev environment.", nil)
		return
	}
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't delete all users.", err)
		return
	}
	response.WithJSON(w, http.StatusOK, nil)
}
