package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/faxter/chirpy/domain"
	"github.com/faxter/chirpy/internal/database"
	"github.com/google/uuid"
)

func ReadinessEndpoint(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.Header().Add(KEY_CONTENT_TYPE, CONTENT_TYPE_PLAIN)
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte("OK"))
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
		respondWithError(responseWriter, 501, logmsg)
		return
	}
	user := domain.User{ID: dbUser.ID, Email: dbUser.Email, CreatedAt: dbUser.CreatedAt, UpdatedAt: dbUser.CreatedAt}
	respondWithJSON(responseWriter, 201, user)
}

func (a *ApiConfig) CreateChirpEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		logmsg := fmt.Sprintf("Error decoding parameters: %s", err)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	if len(params.Body) > MAX_CHIRP_LENGTH {
		respondWithError(responseWriter, 400, "Chirp is too long")
		return
	}

	cleanedBody := censorWords(params.Body)

	dbChirp, err := a.Queries.CreateChirp(context.Background(), database.CreateChirpParams{Body: cleanedBody, UserID: params.UserId})
	if err != nil {
		msg := fmt.Sprint("Could not create database entry:", err)
		respondWithError(responseWriter, 501, msg)
		return
	}

	chirp := domain.Chirp{
		Id: dbChirp.ID, CreatedAt: dbChirp.CreatedAt, UpdatedAt: dbChirp.UpdatedAt, Body: dbChirp.Body, UserId: dbChirp.UserID}
	respondWithJSON(responseWriter, 201, chirp)
}

func censorWords(text string) string {
	bodyWordList := strings.Split(text, " ")
	badSet := make(map[string]struct{})
	for _, w := range BadWords() {
		badSet[strings.ToLower(w)] = struct{}{}
	}

	result := make([]string, len(bodyWordList))

	for i, word := range bodyWordList {
		if _, found := badSet[strings.ToLower(word)]; found {
			result[i] = "****"
		} else {
			result[i] = word
		}
	}

	return strings.Join(result, " ")
}

func (a *ApiConfig) GetChirpsEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	chirpList, err := a.Queries.GetChirps(context.Background())
	if err != nil {
		msg := fmt.Sprint("Could not get chirps from database:", err)
		respondWithError(responseWriter, 501, msg)
		return
	}

	result := []domain.Chirp{}
	for _, dbChirp := range chirpList {
		result = append(result, domain.Chirp{Id: dbChirp.ID, CreatedAt: dbChirp.CreatedAt, UpdatedAt: dbChirp.UpdatedAt, Body: dbChirp.Body, UserId: dbChirp.UserID})
	}

	respondWithJSON(responseWriter, 200, result)
}
