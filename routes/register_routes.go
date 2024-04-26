package routes

import (
	"net/http"
	"sync"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/config"
)

var apiConfig *config.ApiConfig = &config.ApiConfig{}

func init() {
	apiConfig.Mutex = &sync.RWMutex{}

	config.LoadEnv(apiConfig)
	config.ConnectDB(apiConfig)
	config.InitializeDB(apiConfig)
}

func ResourceServerPort() string {
	return apiConfig.ResourceServerPort
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", GenerateReponse)
	mux.HandleFunc("/refresh", GenerateReponse)
	mux.HandleFunc("/follow/{channel_id}", FollowChannel)
	mux.HandleFunc("/unfollow/{channel_id}", UnfollowChannel)
	mux.HandleFunc("/addNewChannel/{channel_id}", AddNewChannel)
}
