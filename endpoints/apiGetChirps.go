package endpoints

import (
	"fmt"
	"net/http"

	"github.com/faxter/chirpy/domain"
	"github.com/faxter/chirpy/internal/database"
	"github.com/google/uuid"
)

func (a *ApiConfig) GetChirpsEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	author := request.URL.Query().Get("author_id")
	chirpList := []database.Chirp{}
	var err error
	if author == "" {
		chirpList, err = a.Queries.GetChirps(request.Context())
		if err != nil {
			msg := fmt.Sprint("Could not get chirps from database:", err)
			respondWithError(responseWriter, 501, msg)
			return
		}
	} else {
		authorId, err := uuid.Parse(author)
		if err != nil {
			msg := fmt.Sprintf("Could not convert %s into UUID: %s", author, err)
			respondWithError(responseWriter, 502, msg)
			return
		}
		chirpList, err = a.Queries.GetChirpsByAuthor(request.Context(), authorId)
		if err != nil {
			msg := fmt.Sprint("Could not get chirps from database:", err)
			respondWithError(responseWriter, 501, msg)
			return
		}
	}

	result := []domain.Chirp{}
	for _, dbChirp := range chirpList {
		result = append(result, domain.Chirp{
			Id:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserId:    dbChirp.UserID})
	}

	respondWithJSON(responseWriter, 200, result)
}

func (a *ApiConfig) GetSingleChirpEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	requestIdString := request.PathValue("chirpID")
	requestedChirpId, err := uuid.Parse(requestIdString)
	if err != nil {
		msg := fmt.Sprintf("Could not convert %s into UUID: %s", requestIdString, err)
		respondWithError(responseWriter, 502, msg)
		return
	}

	chirp, err := a.Queries.GetChirp(request.Context(), requestedChirpId)
	if err != nil {
		msg := fmt.Sprintf("Could not get chirp from database: %s", err)
		respondWithError(responseWriter, 404, msg)
		return
	}

	result := domain.Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	}
	respondWithJSON(responseWriter, 200, result)
}
