package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/faxter/chirpy/endpoints"
	"github.com/faxter/chirpy/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	cfg := loadConfig()
	s := http.NewServeMux()
	s.Handle("/app/", cfg.IncrementsMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	s.HandleFunc("GET /api/healthz", endpoints.ReadinessEndpoint)
	s.HandleFunc("GET /admin/metrics", cfg.MetricsEndpoint)
	s.HandleFunc("POST /admin/reset", cfg.ResetEndpoint)
	s.HandleFunc("POST /api/users", cfg.CreateUserEndpoint)
	s.HandleFunc("PUT /api/users", cfg.UpdateUserEndpoint)
	s.HandleFunc("POST /api/login", cfg.LoginUserEndpoint)
	s.HandleFunc("POST /api/refresh", cfg.RefreshLoginEndpoint)
	s.HandleFunc("POST /api/revoke", cfg.RevokeLoginEndpoint)
	s.HandleFunc("POST /api/chirps", cfg.CreateChirpEndpoint)
	s.HandleFunc("GET /api/chirps", cfg.GetChirpsEndpoint)
	s.HandleFunc("GET /api/chirps/{chirpID}", cfg.GetSingleChirpEndpoint)
	s.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.DeleteSingleChirpEndpoint)
	s.HandleFunc("POST /api/polka/webhooks", cfg.UpgradeUserEndpoint)
	serv := new(http.Server)
	serv.Handler = s
	serv.Addr = ":8080"
	serv.ListenAndServe()
}

func loadConfig() endpoints.ApiConfig {
	godotenv.Load(".env")
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println("could not open connection to database:", dbUrl)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")

	return endpoints.ApiConfig{FileServerHits: atomic.Int32{}, Queries: dbQueries, Platform: platform, Secret: secret}
}
