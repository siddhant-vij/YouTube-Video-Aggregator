package middlewares

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/config"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/utils"
)

func AuthMiddleware(next http.Handler, config *config.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atCookie, err := r.Cookie("access_token")
		if err != nil {
			log.Println(err)
			utils.RespondWithError(w, http.StatusBadRequest, "Access token is missing")
			return
		}

		endpoint := config.VerifyEndpoint
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			log.Println(err)
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		req.AddCookie(atCookie)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer resp.Body.Close()

		var verifyResponse struct {
			Token     string
			TokenUuid string
			UserID    string
			ExpiresIn int64
		}

		if resp.StatusCode != http.StatusOK {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to authenticate user")
			return
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&verifyResponse)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Cannot decode response")
			return
		}

		config.UserId = uuid.MustParse(verifyResponse.UserID)

		next.ServeHTTP(w, r)
	})
}
