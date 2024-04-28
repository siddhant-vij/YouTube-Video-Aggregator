package controllers

import (
	"context"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/config"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/database"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/services"
)

func FetchNewVideos(config *config.ApiConfig) {
	config.Mutex.RLock()
	channelsToFetch, err := config.DBQueries.GetNumChannelsToFetch(context.TODO(), 10)
	if err != nil {
		log.Fatal(err)
	}
	config.Mutex.RUnlock()

	var channels []services.Channel

	for _, channel := range channelsToFetch {
		channel, err := services.GetChannelVideos(channel.ID)
		if err != nil {
			panic(err)
		}

		channels = append(channels, channel)
	}

	videoParams := createVideoParams(&channelsToFetch, &channels, 5)
	insertVideos(videoParams, config)

	config.Mutex.Lock()
	for _, channel := range channelsToFetch {
		var updateLastFetchedParams = database.UpdateLastFetchedAtParams{
			ID:            channel.ID,
			LastFetchedAt: time.Now(),
		}
		err = config.DBQueries.UpdateLastFetchedAt(context.TODO(), updateLastFetchedParams)
		if err != nil {
			log.Fatal(err)
		}
	}
	config.Mutex.Unlock()
}

func createVideoParams(channelsToFetch *[]database.Channel, channels *[]services.Channel, batchSize int) []database.InsertVideoParams {
	wg := &sync.WaitGroup{}
	ch := make(chan database.InsertVideoParams)
	var params []database.InsertVideoParams

	go func() {
		for param := range ch {
			params = append(params, param)
		}
	}()

	validCounts := make([]int32, len(*channelsToFetch))

	for cidx, channelID := range *channelsToFetch {
		for _, video := range (*channels)[cidx].Videos {
			if atomic.LoadInt32(&validCounts[cidx]) >= int32(batchSize) {
				break
			}
			wg.Add(1)
			go func(video services.Video, channelId string, idx int) {
				defer wg.Done()
				param := createOneVideoParam(&video, channelId)
				if param != nil {
					if atomic.AddInt32(&validCounts[idx], 1) <= int32(batchSize) {
						ch <- *param
					}
				}
			}(video, channelID.ID, cidx)
		}
	}

	wg.Wait()
	close(ch)
	return params
}

func createOneVideoParam(video *services.Video, channelID string) *database.InsertVideoParams {
	if video.Description == "..." || video.Description == "" {
		return nil
	}
	return &database.InsertVideoParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       video.Title,
		Description: video.Description,
		ImageUrl:    video.ImageURL,
		Authors:     video.Authors,
		PublishedAt: video.PublishedAt,
		Url:         video.URL,
		ViewCount:   video.ViewCount,
		StarRating:  video.StarRating,
		StarCount:   video.StarCount,
		ChannelID:   channelID,
	}
}

func insertVideos(params []database.InsertVideoParams, config *config.ApiConfig) {
	wg := &sync.WaitGroup{}

	for _, param := range params {
		wg.Add(1)
		go insertOneVideo(param, config, wg)
	}

	wg.Wait()
}

func insertOneVideo(param database.InsertVideoParams, config *config.ApiConfig, wg *sync.WaitGroup) {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()
	defer wg.Done()

	err := config.DBQueries.InsertVideo(context.TODO(), param)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if isUpdatable(param, config) {
				updateStats(param, config)
			}
		} else {
			log.Fatalf("Failed to insert video: %v", err)
		}
	}
}

func isUpdatable(param database.InsertVideoParams, config *config.ApiConfig) bool {
	vc, sr, sc := getStats(param.Url, config)
	return param.ViewCount != vc ||
		param.StarRating != sr ||
		param.StarCount != sc
}

func getStats(url string, config *config.ApiConfig) (viewCount string, starRating string, starCount string) {
	getResponseRow, err := config.DBQueries.GetStatsForURL(context.TODO(), url)
	if err != nil {
		log.Fatalf("Failed to get video stats: %v", err)
	}
	return getResponseRow.ViewCount, getResponseRow.StarRating, getResponseRow.StarCount
}

func updateStats(param database.InsertVideoParams, config *config.ApiConfig) {
	var updateStatsParams = database.UpdateStatsForURLParams{
		Url:        param.Url,
		ViewCount:  param.ViewCount,
		StarRating: param.StarRating,
		StarCount:  param.StarCount,
	}
	err := config.DBQueries.UpdateStatsForURL(context.TODO(), updateStatsParams)
	if err != nil {
		log.Fatalf("Failed to update video stats: %v", err)
	}
}
