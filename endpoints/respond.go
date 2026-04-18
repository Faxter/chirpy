package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func respondWithError(writer http.ResponseWriter, code int, msg string) {
	type returnError struct {
		Error string `json:"error"`
	}
	dat, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error marshalling JSON for error response (not sending response): %s", err)
		writer.WriteHeader(500)
		return
	}
	writer.Header().Set(KEY_CONTENT_TYPE, CONTENT_TYPE_JSON)
	writer.WriteHeader(code)
	writer.Write(dat)
}

func respondWithJSON(writer http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		respondWithError(writer, 500, "Error marshalling payload")
		return
	}

	writer.Header().Set(KEY_CONTENT_TYPE, CONTENT_TYPE_JSON)
	writer.WriteHeader(code)
	writer.Write(dat)
}

func respondWithCleanedJson(body string, responseWriter http.ResponseWriter) {
	type returnCleaned struct {
		Cleaned string `json:"cleaned_body"`
	}

	bodyWordList := strings.Split(body, " ")
	cleanedBody := strings.Join(censorWords(bodyWordList), " ")
	respBody := returnCleaned{
		Cleaned: cleanedBody,
	}

	respondWithJSON(responseWriter, 200, respBody)
}
