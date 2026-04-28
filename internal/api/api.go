package api

import (
	"sync/atomic"

	"github.com/Michael-cmd-sys/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
}
