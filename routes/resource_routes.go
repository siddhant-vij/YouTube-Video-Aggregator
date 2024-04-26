package routes

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/controllers"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/utils"
)

func GenerateReponse(w http.ResponseWriter, r *http.Request) {
	userId := uuid.MustParse("78b7b1d4-bc92-4488-9531-7d805e689feb")
	// Generated from Auth Service - Will integrate soon

	responseForUser, err := controllers.GenerateReponseForUser(apiConfig, userId, 10, 50)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, responseForUser)
}
