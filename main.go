package main

import (
	"net/http"
	"sync/atomic"

	"github.com/faxter/chirpy/endpoints"
)

func main() {
	cfg := endpoints.ApiConfig{FileServerHits: atomic.Int32{}}
	s := http.NewServeMux()
	s.Handle("/app/", cfg.IncrementsMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	s.HandleFunc("GET /api/healthz", endpoints.ReadinessEndpoint)
	s.HandleFunc("POST /api/validate_chirp", endpoints.ChirpValidatorEndpoint)
	s.HandleFunc("GET /admin/metrics", cfg.MetricsEndpoint)
	s.HandleFunc("POST /admin/reset", cfg.ResetMetricsEndpoint)
	serv := new(http.Server)
	serv.Handler = s
	serv.Addr = ":8080"
	serv.ListenAndServe()
}
