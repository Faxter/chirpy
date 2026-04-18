package endpoints

import (
	"sync/atomic"

	"github.com/faxter/chirpy/internal/database"
)

const (
	KEY_CONTENT_TYPE   = "Content-Type"
	CONTENT_TYPE_HTML  = "text/html"
	CONTENT_TYPE_PLAIN = "text/plain; charset=utf-8"
	CONTENT_TYPE_JSON  = "application/json"
	MAX_CHIRP_LENGTH   = 140
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	Queries        *database.Queries
	Platform       string
}

func BadWords() []string {
	return []string{"kerfuffle", "sharbert", "fornax"}
}
