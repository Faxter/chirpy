package endpoints

import (
	"context"
	"fmt"
	"net/http"
)

func (a *ApiConfig) MetricsEndpoint(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.Header().Add(KEY_CONTENT_TYPE, CONTENT_TYPE_HTML)
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write(fmt.Appendf(nil, "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", a.FileServerHits.Load()))
}

func (a *ApiConfig) ResetEndpoint(responseWriter http.ResponseWriter, _ *http.Request) {
	if a.Platform != "dev" {
		respondWithError(responseWriter, 403, "FORBIDDEN")
		return
	}
	a.FileServerHits.Store(0)
	a.Queries.DropUsers(context.Background())
	responseWriter.WriteHeader(http.StatusOK)
}
