package endpoints

import (
	"fmt"
	"net/http"
)

func (a *ApiConfig) MetricsEndpoint(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.Header().Add(KEY_CONTENT_TYPE, CONTENT_TYPE_HTML)
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write(fmt.Appendf(nil, "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", a.FileServerHits.Load()))
}

func (a *ApiConfig) ResetMetricsEndpoint(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.WriteHeader(http.StatusOK)
	a.FileServerHits.Store(0)
}
