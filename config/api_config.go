package config

import (
	"sync"

	"github.com/google/uuid"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/database"
)

type ApiConfig struct {
	DatabaseURL string

	AuthServerPort     string
	AuthVerifyEndpoint string

	ResourceServerPort string

	ChannelBaseURL string

	DBQueries *database.Queries
	Mutex     *sync.RWMutex

	UserId uuid.UUID
}
