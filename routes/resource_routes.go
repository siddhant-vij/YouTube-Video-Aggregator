package routes

import (
	"net/http"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/controllers"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/utils"
)

func GenerateResponse(w http.ResponseWriter, r *http.Request) {
	userId := apiConfig.UserId

	responseForUser, err := controllers.GenerateResponseForUser(apiConfig, userId, 10, 50)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, responseForUser)
}
