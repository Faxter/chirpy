package main

import (
	"database/sql"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/faxter/chirpy/endpoints"
	"github.com/faxter/chirpy/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		os.Exit(1)
	}
	dbQueries := database.New(db)

	cfg := endpoints.ApiConfig{FileServerHits: atomic.Int32{}, Queries: dbQueries}
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
