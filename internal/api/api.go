package api

import (
	"sync/atomic"
	"time"

	"github.com/Michael-cmd-sys/chirpy/internal/database"
	"github.com/google/uuid"
)

type ApiConfig struct {
	DB             *database.Queries
	Platform       string
	FileserverHits atomic.Int32
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
