package config

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
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

	channelParams := createAllChannelParams(&channelIDs, &feeds)
	insertAllChannels(channelParams, config)

	videoParams := createAllVideoParams(&channelIDs, &feeds, 5)
	insertAllVideos(videoParams, config)

	log.Println("Database Initialized!")
}

func createAllChannelParams(initDb *InitDB, feeds *[]utils.Feed) []database.InsertChannelParams {
	wg := &sync.WaitGroup{}
	ch := make(chan database.InsertChannelParams)
	var params []database.InsertChannelParams

	go func() {
		for param := range ch {
			params = append(params, param)
		}
	}()

	for idx, channelID := range initDb.Channels {
		wg.Add(1)
		go createOneChannelParams(&(*feeds)[idx], channelID.ChannelID, wg, ch)
	}

	wg.Wait()
	close(ch)
	return params
}

func createOneChannelParams(feed *utils.Feed, channelID string, wg *sync.WaitGroup, ch chan<- database.InsertChannelParams) {
	defer wg.Done()
	param := database.InsertChannelParams{
		ID:            channelID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Name:          feed.Title,
		Url:           feed.Link.Href,
		LastFetchedAt: time.Now(),
	}
	ch <- param
}

func insertAllChannels(params []database.InsertChannelParams, config *ApiConfig) {
	wg := &sync.WaitGroup{}

	for _, param := range params {
		wg.Add(1)
		go insertOneChannel(param, config, wg)
	}

	wg.Wait()
}

func insertOneChannel(param database.InsertChannelParams, config *ApiConfig, wg *sync.WaitGroup) {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()
	defer wg.Done()
	err := config.DBQueries.InsertChannel(context.TODO(), param)
	if err != nil {
		log.Fatalf("Failed to insert channel: %v", err)
	}
}

func createAllVideoParams(initDb *InitDB, feeds *[]utils.Feed, numVideosPerChannel int) []database.InsertVideoParams {
	wg := &sync.WaitGroup{}
	ch := make(chan database.InsertVideoParams)
	var params []database.InsertVideoParams

	go func() {
		for param := range ch {
			params = append(params, param)
		}
	}()

	validCounts := make([]int32, len(initDb.Channels))

	for cidx, channelID := range initDb.Channels {
		for _, entry := range (*feeds)[cidx].Entries {
			if atomic.LoadInt32(&validCounts[cidx]) >= int32(numVideosPerChannel) {
				break
			}
			wg.Add(1)
			go func(entry utils.Entry, channelId string, idx int) {
				defer wg.Done()
				param := createOneVideoParam(&entry, channelId)
				if param != nil {
					if atomic.AddInt32(&validCounts[idx], 1) <= int32(numVideosPerChannel) {
						ch <- *param
					}
				}
			}(entry, channelID.ChannelID, cidx)
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

func insertAllVideos(params []database.InsertVideoParams, config *ApiConfig) {
	wg := &sync.WaitGroup{}

	for _, param := range params {
		wg.Add(1)
		go func(p database.InsertVideoParams) {
			defer wg.Done()
			insertOneVideo(p, config)
		}(param)
	}

	wg.Wait()
}

func insertOneVideo(param database.InsertVideoParams, config *ApiConfig) {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()
	err := config.DBQueries.InsertVideo(context.TODO(), param)
	if err != nil {
		log.Fatalf("Failed to insert video: %v", err)
	}
}
