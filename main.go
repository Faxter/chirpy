package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

const (
	KEY_CONTENT_TYPE   = "Content-Type"
	CONTENT_TYPE_PLAIN = "text/plain; charset=utf-8"
	CONTENT_TYPE_HTML  = "text/html"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (a *apiConfig) incrementsMetrics(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		a.fileServerHits.Add(1)
		handler.ServeHTTP(writer, req)
	})
}

func main() {
	cfg := apiConfig{fileServerHits: atomic.Int32{}}
	s := http.NewServeMux()
	s.Handle("/app/", cfg.incrementsMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	s.HandleFunc("GET /api/healthz", readinessEndpoint)
	s.HandleFunc("GET /admin/metrics", cfg.metricsEndpoint)
	s.HandleFunc("POST /admin/reset", cfg.resetMetricsEndpoint)
	serv := new(http.Server)
	serv.Handler = s
	serv.Addr = ":8080"
	serv.ListenAndServe()
}

func readinessEndpoint(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.Header().Add(KEY_CONTENT_TYPE, CONTENT_TYPE_PLAIN)
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte("OK"))
}

func (a *apiConfig) metricsEndpoint(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.Header().Add(KEY_CONTENT_TYPE, CONTENT_TYPE_HTML)
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write(fmt.Appendf(nil, "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", a.fileServerHits.Load()))
}

func (a *apiConfig) resetMetricsEndpoint(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.WriteHeader(http.StatusOK)
	a.fileServerHits.Store(0)
}
