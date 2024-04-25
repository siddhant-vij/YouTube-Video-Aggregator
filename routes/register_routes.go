package routes

import (
	"context"
	"net/http"
	"sync"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/config"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/database"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/utils"
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

func GetDataFromDB(w http.ResponseWriter, r *http.Request) {
	apiConfig.Mutex.RLock()
	defer apiConfig.Mutex.RUnlock()

	channels, err := apiConfig.DBQueries.GetAllChannels(context.TODO())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	channelFollows, err := apiConfig.DBQueries.GetAllChannelFollows(context.TODO())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	videos, err := apiConfig.DBQueries.GetAllVideos(context.TODO())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var response struct {
		Channels       []database.Channel
		ChannelFollows []database.ChannelFollow
		Videos         []database.Video
	}
	response.Channels = channels
	response.ChannelFollows = channelFollows
	response.Videos = videos
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/db/test", GetDataFromDB)
}
