package routes

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"

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

func getFiveChannels() ([]database.Channel, error) {
	apiConfig.Mutex.RLock()
	defer apiConfig.Mutex.RUnlock()

	channels, err := apiConfig.DBQueries.GetFiveChannels(context.TODO())
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func insertFiveFollowsForUser(user_id uuid.UUID) error {
	channels, err := getFiveChannels()
	if err != nil {
		return err
	}
	apiConfig.Mutex.Lock()
	defer apiConfig.Mutex.Unlock()

	for _, channel := range channels {
		var channelFollowParams = database.InsertChannelFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user_id,
			ChannelID: channel.ID,
		}
		err := apiConfig.DBQueries.InsertChannelFollow(context.TODO(), channelFollowParams)
		if err != nil {
			return err
		}
	}

	return nil
}

func getFollowedChannelsForUser(user_id uuid.UUID) ([]database.Channel, error) {
	apiConfig.Mutex.RLock()
	defer apiConfig.Mutex.RUnlock()

	channels, err := apiConfig.DBQueries.GetUserFollowedChannels(context.TODO(), user_id)
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func getNotYetFollowedChannelsForUser(user_id uuid.UUID) ([]database.Channel, error) {
	apiConfig.Mutex.RLock()
	defer apiConfig.Mutex.RUnlock()

	channels, err := apiConfig.DBQueries.GetOtherChannelsForUser(context.TODO(), user_id)
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func getTwentyVideosForUser(user_id uuid.UUID) ([]database.Video, error) {
	apiConfig.Mutex.RLock()
	defer apiConfig.Mutex.RUnlock()

	var userVideoParams = database.GetUserVideosParams{
		UserID: user_id,
		Limit:  20,
	}

	videos, err := apiConfig.DBQueries.GetUserVideos(context.TODO(), userVideoParams)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func getAllChannelFollows() ([]database.ChannelFollow, error) {
	apiConfig.Mutex.RLock()
	defer apiConfig.Mutex.RUnlock()

	channelFollows, err := apiConfig.DBQueries.GetAllChannelFollows(context.TODO())
	if err != nil {
		return nil, err
	}

	return channelFollows, nil
}

func GetReponseForUser(w http.ResponseWriter, r *http.Request) {
	channelFollows, err := getAllChannelFollows()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	user_id := uuid.MustParse("bc8d331a-470b-4dea-a05b-e62168bd0cba")
	// Generated from Auth Service - Will integrate soon
	if len(channelFollows) == 0 {
		insertFiveFollowsForUser(user_id)
	}
	videos, err := getTwentyVideosForUser(user_id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var response struct {
		Videos           []database.Video
		ChannelsFollowed []database.Channel
		OtherChannels    []database.Channel
	}
	response.Videos = videos
	response.ChannelsFollowed, _ = getFollowedChannelsForUser(user_id)
	response.OtherChannels, _ = getNotYetFollowedChannelsForUser(user_id)
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/response", GetReponseForUser)
}
