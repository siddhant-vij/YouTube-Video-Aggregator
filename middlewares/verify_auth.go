package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	
	"github.com/siddhant-vij/YouTube-Video-Aggregator/config"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/utils"
)

func AuthMiddleware(next http.Handler, config *config.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atCookie, err := r.Cookie("access_token")
		if err != nil {
			config.AuthStatus = "false: "
			utils.RespondWithError(w, http.StatusBadRequest, config.GetAuthStatus() + "access token is missing")
			return
		}

		endpoint := config.VerifyEndpoint
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			config.AuthStatus = "false: "
			utils.RespondWithError(w, http.StatusBadRequest, config.GetAuthStatus() + err.Error())
			return
		}

		req.AddCookie(atCookie)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			config.AuthStatus = "false: "
			utils.RespondWithError(w, http.StatusBadRequest, config.GetAuthStatus() + err.Error())
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
			config.AuthStatus = "false: "
			utils.RespondWithError(w, http.StatusBadRequest, config.GetAuthStatus() + "Failed to authenticate user")
			return
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&verifyResponse)
		if err != nil {
			config.AuthStatus = "false: "
			utils.RespondWithError(w, http.StatusBadRequest, config.GetAuthStatus() + "Cannot decode response")
		}

		config.UserId = uuid.MustParse(verifyResponse.UserID)

		next.ServeHTTP(w, r)
	})
}
