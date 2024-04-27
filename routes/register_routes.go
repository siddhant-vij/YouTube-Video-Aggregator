package routes

import (
	"net/http"
	"sync"

	"github.com/jasonlvhit/gocron"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/config"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/controllers"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/middlewares"
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
	go executeFetchCronJob()
	go executeDeleteCronJob()

	mux.Handle("/", middlewares.AuthMiddleware(http.HandlerFunc(GenerateReponse), apiConfig))

	mux.Handle("/refresh", middlewares.AuthMiddleware(http.HandlerFunc(GenerateReponse), apiConfig))

	mux.Handle("/follow/{channel_id}", middlewares.AuthMiddleware(http.HandlerFunc(FollowChannel), apiConfig))

	mux.Handle("/unfollow/{channel_id}", middlewares.AuthMiddleware(http.HandlerFunc(UnfollowChannel), apiConfig))

	mux.Handle("/addNewChannel/{channel_id}", middlewares.AuthMiddleware(http.HandlerFunc(AddNewChannel), apiConfig))
}
