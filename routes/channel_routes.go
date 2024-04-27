package routes

import (
	"net/http"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/controllers"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/utils"
)

func FollowChannel(w http.ResponseWriter, r *http.Request) {
	userId := apiConfig.UserId

	channelId := r.PathValue("channel_id")
	err := controllers.AddChannelFollowForUser(apiConfig, userId, channelId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseForUser, err := controllers.GenerateReponseForUser(apiConfig, userId, 10, 50)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, responseForUser)
}

func UnfollowChannel(w http.ResponseWriter, r *http.Request) {
	userId := apiConfig.UserId

	channelId := r.PathValue("channel_id")
	err := controllers.RemoveChannelFollowForUser(apiConfig, userId, channelId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseForUser, err := controllers.GenerateReponseForUser(apiConfig, userId, 10, 50)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, responseForUser)
}

func AddNewChannel(w http.ResponseWriter, r *http.Request) {
	userId := apiConfig.UserId

	channelId := r.PathValue("channel_id")
	err := controllers.AddChannelAndVideos(apiConfig, channelId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = controllers.AddChannelFollowForUser(apiConfig, userId, channelId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseForUser, err := controllers.GenerateReponseForUser(apiConfig, userId, 10, 50)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, responseForUser)
}
