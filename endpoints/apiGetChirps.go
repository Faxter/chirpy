package endpoints

import (
	"fmt"
	"net/http"

	"github.com/faxter/chirpy/domain"
)

func (a *ApiConfig) GetChirpsEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	chirpList, err := a.Queries.GetChirps(request.Context())
	if err != nil {
		msg := fmt.Sprint("Could not get chirps from database:", err)
		respondWithError(responseWriter, 501, msg)
		return
	}

	result := []domain.Chirp{}
	for _, dbChirp := range chirpList {
		result = append(result, domain.Chirp{
			Id:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body, UserId: dbChirp.UserID})
	}

	respondWithJSON(responseWriter, 200, result)
}
