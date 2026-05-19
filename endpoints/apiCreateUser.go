package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/faxter/chirpy/domain"
	"github.com/faxter/chirpy/internal/auth"
	"github.com/faxter/chirpy/internal/database"
	"github.com/google/uuid"
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
		UpdatedAt: dbUser.CreatedAt,
		IsPremium: dbUser.IsChirpyRed}
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

	updatedUser := domain.User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		IsPremium: dbUser.IsChirpyRed}
	respondWithJSON(responseWriter, 200, updatedUser)
}

func (a *ApiConfig) UpgradeUserEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	apiKey, err := auth.GetAPIKey(request.Header)
	if err != nil {
		logmsg := fmt.Sprintf("Error extracting API key from header: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 401, logmsg)
		return
	}

	if apiKey != a.PolkaKey {
		logmsg := fmt.Sprintln("incorrect API key!")
		fmt.Println(logmsg)
		respondWithError(responseWriter, 401, logmsg)
		return
	}

	type data struct {
		UserId string `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
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

	if params.Event != "user.upgraded" {
		respondWithError(responseWriter, 204, "")
		return
	}

	userId, err := uuid.Parse(params.Data.UserId)
	if err != nil {
		logmsg := fmt.Sprintf("Error parsing user id: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 404, logmsg)
		return
	}

	err = a.Queries.UpgradeUserToRed(request.Context(), userId)
	if err != nil {
		logmsg := fmt.Sprintf("Error upgrading user in database: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 404, logmsg)
	}

	respondWithJSON(responseWriter, 204, "")
}
