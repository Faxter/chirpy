package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/faxter/chirpy/domain"
)

func ReadinessEndpoint(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.Header().Add(KEY_CONTENT_TYPE, CONTENT_TYPE_PLAIN)
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte("OK"))
}

func ChirpValidatorEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		logmsg := fmt.Sprintf("Error decoding parameters: %s", err)
		log.Println(logmsg)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	if len(params.Body) > MAX_CHIRP_LENGTH {
		respondWithError(responseWriter, 400, "Chirp is too long")
		return
	}

	respondWithCleanedJson(params.Body, responseWriter)
}

func censorWords(wordList []string) []string {
	badSet := make(map[string]struct{})
	for _, w := range BadWords() {
		badSet[strings.ToLower(w)] = struct{}{}
	}

	result := make([]string, len(wordList))

	for i, word := range wordList {
		if _, found := badSet[strings.ToLower(word)]; found {
			result[i] = "****"
		} else {
			result[i] = word
		}
	}

	return result
}

func (a *ApiConfig) CreateUserEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"email"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		logmsg := fmt.Sprintf("Error decoding parameters: %s", err)
		log.Println(logmsg)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	dbUser, err := a.Queries.CreateUser(context.Background(), params.Body)
	if err != nil {
		logmsg := fmt.Sprintf("Error creating user %s in database: %s", params.Body, err)
		log.Println(logmsg)
		respondWithError(responseWriter, 500, logmsg)
		return
	}
	user := domain.User{ID: dbUser.ID, Email: dbUser.Email, CreatedAt: dbUser.CreatedAt, UpdatedAt: dbUser.CreatedAt}
	respondWithJSON(responseWriter, 201, user)
}
