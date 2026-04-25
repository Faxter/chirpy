package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/faxter/chirpy/domain"
	"github.com/faxter/chirpy/internal/auth"
)

func (a *ApiConfig) LoginUserEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	dbUser, err := a.Queries.GetUser(request.Context(), params.Email)
	if err != nil {
		logmsg := fmt.Sprintln("Incorrect email or password")
		fmt.Println(logmsg)
		respondWithError(responseWriter, 401, logmsg)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		logmsg := fmt.Sprintf("Error while checking password: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	if !match {
		logmsg := fmt.Sprintln("Incorrect email or password")
		fmt.Println(logmsg)
		respondWithError(responseWriter, 401, logmsg)
		return
	}

	user := domain.User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	respondWithJSON(responseWriter, 200, user)
}
