package controllers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/config"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/database"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/utils"
)

func FetchNewVideos(config *config.ApiConfig) {
	config.Mutex.RLock()
	channelsToFetch, err := config.DBQueries.GetNumChannelsToFetch(context.TODO(), 10)
	if err != nil {
		log.Fatal(err)
	}
	config.Mutex.RUnlock()

	var feeds []utils.Feed

	for _, channel := range channelsToFetch {
		var ytFeedInputParams = utils.YouTubeFeedInputParams{
			Client:         &http.Client{},
			ChannelBaseURL: config.ChannelBaseURL,
		}

		feed, err := ytFeedInputParams.GetFeed(channel.ID, utils.FTChannel)
		if err != nil {
			panic(err)
		}

		feeds = append(feeds, feed)
	}

	videoParams := createVideoParams(&channelsToFetch, &feeds, 5)
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

func createVideoParams(channels *[]database.Channel, feeds *[]utils.Feed, batchSize int) []database.InsertVideoParams {
	wg := &sync.WaitGroup{}
	ch := make(chan database.InsertVideoParams)
	var params []database.InsertVideoParams

	go func() {
		for param := range ch {
			params = append(params, param)
		}
	}()

	validCounts := make([]int32, len(*channels))

	for cidx, channelID := range *channels {
		for _, entry := range (*feeds)[cidx].Entries {
			if atomic.LoadInt32(&validCounts[cidx]) >= int32(batchSize) {
				break
			}
			wg.Add(1)
			go func(entry utils.Entry, channelId string, idx int) {
				defer wg.Done()
				param := createOneVideoParam(&entry, channelId)
				if param != nil {
					if atomic.AddInt32(&validCounts[idx], 1) <= int32(batchSize) {
						ch <- *param
					}
				}
			}(entry, channelID.ID, cidx)
		}
	}

	wg.Wait()
	close(ch)
	return params
}

func createOneVideoParam(entry *utils.Entry, channelID string) *database.InsertVideoParams {
	description := string(entry.Media.Description)
	if description == "" {
		return nil
	}
	return &database.InsertVideoParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       entry.Title,
		Description: utils.ShortenText(description),
		ImageUrl:    entry.Media.Image.URL,
		Authors:     entry.Author.Name,
		PublishedAt: entry.Published,
		Url:         entry.Link.Href,
		ViewCount:   utils.ShortenViewCount(entry.Media.Community.Statistics.Views),
		StarRating:  entry.Media.Community.StarRating.Average,
		StarCount:   utils.ShortenStarCount(entry.Media.Community.StarRating.Count),
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
