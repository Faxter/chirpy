package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
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
