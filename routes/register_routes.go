package routes

import (
	"net/http"
	"sync"

	"github.com/jasonlvhit/gocron"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/config"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/controllers"
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

func executeFetchCronJob() {
	gocron.Every(1).Minute().Do(controllers.FetchNewVideos, apiConfig)
	<-gocron.Start()
}

func executeDeleteCronJob() {
	gocron.Every(1).Week().Do(controllers.DeleteOldVideos, apiConfig)
	<-gocron.Start()
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", GenerateReponse)
	mux.HandleFunc("/refresh", GenerateReponse)
	mux.HandleFunc("/follow/{channel_id}", FollowChannel)
	mux.HandleFunc("/unfollow/{channel_id}", UnfollowChannel)
	mux.HandleFunc("/addNewChannel/{channel_id}", AddNewChannel)

	go executeFetchCronJob()
	go executeDeleteCronJob()
}
