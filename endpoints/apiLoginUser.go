package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/faxter/chirpy/domain"
	"github.com/faxter/chirpy/internal/auth"
	"github.com/faxter/chirpy/internal/database"
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

	jwt, err := auth.MakeJWT(dbUser.ID, a.Secret, 1*time.Hour)
	if err != nil {
		logmsg := fmt.Sprintf("could not create JWT: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	refreshToken := auth.MakeRefreshToken()
	a.Queries.CreateRefreshToken(request.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	})

	user := domain.User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        jwt,
		RefreshToken: refreshToken,
	}
	respondWithJSON(responseWriter, 200, user)
}

func (a *ApiConfig) RefreshLoginEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		logmsg := fmt.Sprintf("Error extracting token from authorization header: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 401, logmsg)
		return
	}

	refreshToken, err := a.Queries.GetRefreshToken(request.Context(), token)
	if err != nil {
		logmsg := fmt.Sprintf("Error looking up refresh token in database: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 401, logmsg)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		logmsg := fmt.Sprintln("refresh token is expired!")
		fmt.Println(logmsg)
		respondWithError(responseWriter, 401, logmsg)
		return
	}

	if refreshToken.RevokedAt.Valid {
		logmsg := fmt.Sprintf("refresh token has been revoked at %s!\n", refreshToken.RevokedAt.Time.Format(time.RFC1123))
		fmt.Println(logmsg)
		respondWithError(responseWriter, 401, logmsg)
		return
	}

	newAccessToken, err := auth.MakeJWT(refreshToken.UserID, a.Secret, 1*time.Hour)
	if err != nil {
		logmsg := fmt.Sprintf("Error creating JWT: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(responseWriter, 200, response{Token: newAccessToken})
}

func (a *ApiConfig) RevokeLoginEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		logmsg := fmt.Sprintf("Error extracting token from authorization header: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	err = a.Queries.RevokeRefreshToken(request.Context(), token)
	if err != nil {
		logmsg := fmt.Sprintf("Error revoking the refresh token: %s", err)
		fmt.Println(logmsg)
		respondWithError(responseWriter, 500, logmsg)
		return
	}

	respondWithJSON(responseWriter, 204, "")
}
