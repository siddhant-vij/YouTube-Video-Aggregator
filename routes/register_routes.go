package routes

import (
	"context"
	"net/http"
	"sync"

	"github.com/siddhant-vij/RSS-Feed-Aggregator/config"
	"github.com/siddhant-vij/RSS-Feed-Aggregator/database"
	"github.com/siddhant-vij/RSS-Feed-Aggregator/utils"
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

	feed, err := apiConfig.DBQueries.GetAllFeeds(context.TODO())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	feedFollows, err := apiConfig.DBQueries.GetAllFeedFollows(context.TODO())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	posts, err := apiConfig.DBQueries.GetAllPosts(context.TODO())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var response struct {
		Feeds       []database.Feed
		FeedFollows []database.FeedFollow
		Posts       []database.Post
	}
	response.Feeds = feed
	response.FeedFollows = feedFollows
	response.Posts = posts
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/db/test", GetDataFromDB)
}
