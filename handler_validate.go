package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"
)

func validateChirpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body string `json:"body"`
		}
		type returnVals struct {
			Valid bool   `json:"valid"`
			Error string `json:"error"`
			CleanedBody string  `json:"cleaned_body"`
		}

		params := parameters{}
		respBody := returnVals{}
		status := http.StatusOK

		err := json.NewDecoder(r.Body).Decode(&params)

		if err != nil {
			status = http.StatusBadRequest
			respBody.Error = "Invalid JSON"

		} else if utf8.RuneCountInString(params.Body) > 140 {
			status = http.StatusBadRequest
			respBody.Valid = false
			respBody.Error = "Chirp is too long"

		} else {
			respBody.Valid = true
			respBody.CleanedBody = cleanBody(params.Body)
			
		}

		dat, err := json.Marshal(respBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return

		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(dat)
	}
}

func cleanBody(Body string) string {
	profaneWords := map[string]struct{} {
		"kerfuffle" :{},
		"sharbert" :{},
		"fornax" :{}, 
	}
	
	
	
	words := strings.Split(Body," ")

	for i,word:= range words {
		if _,exists := profaneWords[strings.ToLower(word)];exists{
			words[i]= "****"
		}
	}

	return strings.Join(words," ")
}