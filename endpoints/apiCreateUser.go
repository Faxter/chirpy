package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/faxter/chirpy/domain"
)

func (a *ApiConfig) CreateUserEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"email"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		logmsg := fmt.Sprintf("Error decoding parameters: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	dbUser, err := a.Queries.CreateUser(request.Context(), params.Body)
	if err != nil {
		logmsg := fmt.Sprintf("Error creating user %s in database: %s", params.Body, err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 501, logmsg)
		return
	}
	user := domain.User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.CreatedAt}
	respondWithJSON(responseWriter, 201, user)
}
