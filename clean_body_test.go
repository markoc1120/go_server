package main

import "testing"

func TestCleanBody(t *testing.T) {
	var cases = []struct {
		body        string
		cleanedBody string
	}{
		{
			body:        "I had something interesting for breakfast",
			cleanedBody: "I had something interesting for breakfast",
		},
		{
			body:        "I hear Mastodon is better than Chirpy. sharbert I need to migrate",
			cleanedBody: "I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			body:        "I really need a kerfuffle to go to bed sooner, Fornax !",
			cleanedBody: "I really need a **** to go to bed sooner, **** !",
		},
	}

	for _, test := range cases {
		t.Run(test.body, func(t *testing.T) {
			got := cleanBody(test.body)
			if got != test.cleanedBody {
				t.Errorf("got %q, want %q", got, test.cleanedBody)
			}
		})
	}
}
