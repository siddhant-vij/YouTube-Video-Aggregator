package config

import (
	"sync"

	"github.com/siddhant-vij/RSS-Feed-Aggregator/database"
)

type ApiConfig struct {
	DatabaseURL string

	AuthServerPort     string
	AuthVerifyEndpoint string

	ResourceServerPort string

	ChannelBaseURL string

	DBQueries *database.Queries
	Mutex     *sync.RWMutex
}
