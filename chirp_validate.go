package main

import (
	"encoding/json"
	"net/http"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameteres struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameteres{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		// w.WriteHeader(http.StatusInternalServerError)
		// w.Write([]byte(`{"error":"Something went wrong"}`))
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		// w.WriteHeader(http.StatusBadRequest)
		// w.Write([]byte(`{"error":"Chirp is too long"}`))
		return
	}

	type payload struct {
		CleanedBody string `json:"cleaned_body"`
	}
	respondWithJSON(w, http.StatusOK, payload{CleanedBody: cleanBody(params.Body)})
}
