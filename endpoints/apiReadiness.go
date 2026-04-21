package endpoints

import "net/http"

func ReadinessEndpoint(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.Header().Add(KEY_CONTENT_TYPE, CONTENT_TYPE_PLAIN)
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte("OK"))
}
