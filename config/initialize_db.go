package config

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/database"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/utils"
)

type InitDB struct {
	Channels []struct {
		ChannelID string `json:"channelId"`
	} `json:"channels"`
}

func InitializeDB(config *ApiConfig) {
	var channelIDs = InitDB{}
	jsonFile, err := os.Open("../init_db.json")
	if err != nil {
		log.Fatal(err)
	}
	reader := json.NewDecoder(jsonFile)

	if err := reader.Decode(&channelIDs); err != nil {
		log.Fatal(err)
	}

	var feeds []utils.Feed

	for _, channelID := range channelIDs.Channels {
		var ytFeedInputParams = utils.YouTubeFeedInputParams{
			Client:         &http.Client{},
			ChannelBaseURL: config.ChannelBaseURL,
		}

		feed, err := ytFeedInputParams.GetFeed(channelID.ChannelID, utils.FTChannel)
		if err != nil {
			panic(err)
		}

		feeds = append(feeds, feed)
	}

	channelParams := createChannelParams(&channelIDs, &feeds)
	insertAllChannels(channelParams, config)

	videoParams := createVideoParams(&channelIDs, &feeds, 5)
	// -1 for all videos of each channel
	insertAllVideos(videoParams, config)
}

func createChannelParams(initDb *InitDB, feeds *[]utils.Feed) []database.InsertChannelParams {
	params := make([]database.InsertChannelParams, 0)
	for idx, channelID := range initDb.Channels {
		var parameter = database.InsertChannelParams{
			ID:        channelID.ChannelID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      (*feeds)[idx].Title,
			Url:       (*feeds)[idx].Link.Href,
		}
		params = append(params, parameter)
	}
	return params
}

func insertAllChannels(params []database.InsertChannelParams, config *ApiConfig) {
	wg := &sync.WaitGroup{}
	for _, insertChannelParam := range params {
		wg.Add(1)
		go insertOneChannel(insertChannelParam, config, wg)
	}
	wg.Wait()
}

func insertOneChannel(insertChannelParam database.InsertChannelParams, config *ApiConfig, wg *sync.WaitGroup) {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()
	defer wg.Done()
	err := config.DBQueries.InsertChannel(context.TODO(), insertChannelParam)
	if err != nil {
		log.Fatal(err)
	}
}

func createVideoParams(initDb *InitDB, feeds *[]utils.Feed, numVideosPerChannel int) []database.InsertVideoParams {
	params := make([]database.InsertVideoParams, 0)
	for cidx, channelID := range initDb.Channels {
		for eidx, entry := range (*feeds)[cidx].Entries {
			if eidx == numVideosPerChannel {
				break
			}
			description := string(entry.Media.Description)
			if description == "" {
				continue
			}
			var parameter = database.InsertVideoParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       entry.Title,
				Description: utils.SummarizeText(description),
				ImageUrl:    entry.Media.Image.URL,
				Authors:     entry.Author.Name,
				PublishedAt: entry.Published,
				Url:         entry.Link.Href,
				ViewCount:   utils.ShortenViewCount(entry.Media.Community.Statistics.Views),
				StarRating:  entry.Media.Community.StarRating.Average,
				StarCount:   utils.ShortenStarCount(entry.Media.Community.StarRating.Count),
				ChannelID:   channelID.ChannelID,
			}
			params = append(params, parameter)
			eidx += 1
		}
	}
	return params
}

func insertAllVideos(params []database.InsertVideoParams, config *ApiConfig) {
	wg := &sync.WaitGroup{}
	for _, insertVideoParam := range params {
		wg.Add(1)
		go insertOneVideo(insertVideoParam, config, wg)
	}
	wg.Wait()
}

func insertOneVideo(insertVideoParam database.InsertVideoParams, config *ApiConfig, wg *sync.WaitGroup) {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()
	defer wg.Done()
	err := config.DBQueries.InsertVideo(context.TODO(), insertVideoParam)
	if err != nil {
		log.Fatal(err)
	}
}
