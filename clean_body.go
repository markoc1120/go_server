package main

import (
	"slices"
	"strings"
)

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
