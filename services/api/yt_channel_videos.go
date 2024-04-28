package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const (
	ytSearchBaseURL = "https://www.googleapis.com/youtube/v3/search"
	ytVideoBaseURL  = "https://www.googleapis.com/youtube/v3/videos"
)

var ytClient *http.Client = &http.Client{}

type YTApiChannel struct {
	Items []struct {
		Id struct {
			VideoId string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			ChannelId   string `json:"channelId"`
			Title       string `json:"channelTitle"`
			LiveContent string `json:"liveBroadcastContent"`
		} `json:"snippet"`
	} `json:"items"`
	Name   string       `json:"-"`
	URL    string       `json:"-"`
	Videos []YTApiVideo `json:"-"`
}

type YTApiVideo struct {
	Items []struct {
		Snippet struct {
			Title       string    `json:"title"`
			Description string    `json:"description"`
			PublishedAt time.Time `json:"publishedAt"`
			Thumbnails  struct {
				Standard struct {
					URL string `json:"url"`
				} `json:"standard"`
			} `json:"thumbnails"`
			Authors string `json:"channelTitle"`
		} `json:"snippet"`
		Statistics struct {
			ViewCount string `json:"viewCount"`
			StarCount string `json:"likeCount"`
		} `json:"statistics"`
	} `json:"items"`
	VideoId string `json:"-"`
}

func GetYtApiChannel(apiKey, channelID string) (YTApiChannel, error) {
	request, err := http.NewRequest("GET", ytSearchBaseURL, nil)
	if err != nil {
		return YTApiChannel{}, err
	}

	request.Header.Add("Accept", "application/json")

	query := request.URL.Query()
	query.Add("part", "id")
	query.Add("part", "snippet")
	query.Add("maxResults", "15") // Universe is max. 15 videos (~Feed)
	query.Add("channelId", channelID)
	query.Add("order", "date")
	query.Add("key", apiKey)
	request.URL.RawQuery = query.Encode()

	resp, err := ytClient.Get(request.URL.String())
	if err != nil {
		return YTApiChannel{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return YTApiChannel{}, errors.New("quota exhausted")
	}

	decoder := json.NewDecoder(resp.Body)

	var channel YTApiChannel
	err = decoder.Decode(&channel)
	if err != nil {
		return YTApiChannel{}, err
	}

	var videoIds []string
	for _, item := range channel.Items {
		if item.Snippet.LiveContent == "live" {
			continue
		}
		videoIds = append(videoIds, item.Id.VideoId)
	}

	channel.Videos, err = getYtApiVideos(apiKey, videoIds)
	if err != nil {
		return YTApiChannel{}, err
	}

	channel.Name = channel.Items[0].Snippet.Title
	channel.URL = "https://www.youtube.com/channel/" + channelID
	return channel, nil
}

func getYtApiVideos(apiKey string, videoIds []string) ([]YTApiVideo, error) {
	chVideo := make(chan YTApiVideo)
	chError := make(chan error)

	for _, videoId := range videoIds {
		go getYtApiVideo(apiKey, videoId, chVideo, chError)
	}

	var videos []YTApiVideo
	for range videoIds {
		select {
		case video := <-chVideo:
			videos = append(videos, video)
		case err := <-chError:
			return nil, err
		}
	}
	return videos, nil
}

func getYtApiVideo(apiKey, videoId string, chVideo chan<- YTApiVideo, chError chan<- error) {
	request, err := http.NewRequest("GET", ytVideoBaseURL, nil)
	if err != nil {
		chError <- err
		return
	}

	request.Header.Add("Accept", "application/json")

	query := request.URL.Query()
	query.Add("part", "snippet")
	query.Add("part", "statistics")
	query.Add("id", videoId)
	query.Add("key", apiKey)
	request.URL.RawQuery = query.Encode()

	resp, err := ytClient.Get(request.URL.String())
	if err != nil {
		chError <- err
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		chError <- errors.New("quota exhausted")
		return
	}

	decoder := json.NewDecoder(resp.Body)

	var video YTApiVideo
	err = decoder.Decode(&video)
	if err != nil {
		chError <- err
		return
	}

	video.VideoId = videoId
	chVideo <- video
}
