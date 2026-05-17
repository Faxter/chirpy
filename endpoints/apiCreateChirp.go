package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/faxter/chirpy/domain"
	"github.com/faxter/chirpy/internal/auth"
	"github.com/faxter/chirpy/internal/database"
	"github.com/google/uuid"
)

func (a *ApiConfig) CreateChirpEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		logmsg := fmt.Sprintf("Error extracting token from authorization header: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	userId, err := auth.ValidateJWT(token, a.Secret)
	if err != nil {
		logmsg := fmt.Sprintf("Error validating user: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 401, logmsg)
		return
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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

	dbChirp, err := a.Queries.CreateChirp(request.Context(), database.CreateChirpParams{Body: cleanedBody, UserID: userId})
	if err != nil {
		msg := fmt.Sprint("Could not create database entry:", err)
		respondWithError(responseWriter, 501, msg)
		return
	}

	chirp := domain.Chirp{
		Id:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID}
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

func (a *ApiConfig) DeleteSingleChirpEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		logmsg := fmt.Sprintf("Error extracting token from authorization header: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 401, logmsg)
		return
	}

	userIdOfRequestingUser, err := auth.ValidateJWT(token, a.Secret)
	if err != nil {
		logmsg := fmt.Sprintf("Error validating user: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 403, logmsg)
		return
	}

	chirpIdString := request.PathValue("chirpID")
	chirpIdToDelete, err := uuid.Parse(chirpIdString)
	if err != nil {
		logmsg := fmt.Sprintf("Error parsing chirp id: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	chirpData, err := a.Queries.GetChirp(request.Context(), chirpIdToDelete)
	if err != nil {
		logmsg := fmt.Sprintf("Could not find chirp in database: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 404, logmsg)
		return
	}

	if chirpData.UserID != userIdOfRequestingUser {
		logmsg := fmt.Sprintln("You are not allowed to delete chirp of another user!")
		fmt.Println(logmsg)
		respondWithError(responseWriter, 403, logmsg)
		return
	}

	err = a.Queries.DeleteChirp(request.Context(), chirpIdToDelete)
	if err != nil {
		logmsg := fmt.Sprintf("Error deleting chirp from database: %s", err)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	respondWithJSON(responseWriter, 204, "")
}
