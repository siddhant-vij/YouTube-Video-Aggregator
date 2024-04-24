package routes

import (
	"net/http"

	"github.com/siddhant-vij/RSS-Feed-Aggregator/config"
	"github.com/siddhant-vij/RSS-Feed-Aggregator/utils"
)

var apiConfig *config.ApiConfig = &config.ApiConfig{}

func init() {
	config.LoadEnv(apiConfig)
}

func ResourceServerPort() string {
	return apiConfig.ResourceServerPort
}

func healthChecker(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, "OK, Resource!")
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthChecker", healthChecker)
}
