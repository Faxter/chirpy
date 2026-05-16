package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/faxter/chirpy/domain"
	"github.com/faxter/chirpy/internal/auth"
	"github.com/faxter/chirpy/internal/database"
)

func (a *ApiConfig) CreateUserEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
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

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		logmsg := fmt.Sprintf("Error hashing password: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	dbUser, err := a.Queries.CreateUser(request.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass})
	if err != nil {
		logmsg := fmt.Sprintf("Error creating user %s in database: %s", params.Email, err)
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

func (a *ApiConfig) UpdateUserEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	accessToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		logmsg := fmt.Sprintf("Error extracting access token from header: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 401, logmsg)
		return
	}

	userId, err := auth.ValidateJWT(accessToken, a.Secret)
	if err != nil {
		logmsg := fmt.Sprintf("Error validating user: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 401, logmsg)
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		logmsg := fmt.Sprintf("Error decoding parameters: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		logmsg := fmt.Sprintf("Error hashing password: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	dbUser, err := a.Queries.UpdateUser(request.Context(), database.UpdateUserParams{
		ID:             userId,
		Email:          params.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		logmsg := fmt.Sprintf("Error updating database: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	updatedUser := domain.User{ID: dbUser.ID, CreatedAt: dbUser.CreatedAt, UpdatedAt: dbUser.UpdatedAt, Email: dbUser.Email}
	respondWithJSON(responseWriter, 200, updatedUser)
}
