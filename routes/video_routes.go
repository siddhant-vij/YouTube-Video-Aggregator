package routes

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/controllers"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/utils"
)

func AddBookmark(w http.ResponseWriter, r *http.Request) {
	userId := apiConfig.UserId
	videoId := uuid.MustParse(r.PathValue("video_id"))
	updated, err := controllers.AddBookmarkToVideoForUser(apiConfig, userId, videoId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, updated)
}

func RemoveBookmark(w http.ResponseWriter, r *http.Request) {
	userId := apiConfig.UserId
	videoId := uuid.MustParse(r.PathValue("video_id"))
	updated, err := controllers.RemoveBookmarkFromVideoForUser(apiConfig, userId, videoId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, updated)
}

func GetBookmarkedVideos(w http.ResponseWriter, r *http.Request) {
	userId := apiConfig.UserId
	responseForUser, err := controllers.GetBookmarkedVideosForUser(apiConfig, userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, responseForUser)
}
