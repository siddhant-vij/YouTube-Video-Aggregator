package config

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/database"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/services"
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

	var channels []services.Channel

	for _, channelID := range channelIDs.Channels {
		channel, err := services.GetChannelVideos(channelID.ChannelID)
		if err != nil {
			panic(err)
		}

		channels = append(channels, channel)
	}

	channelParams := createAllChannelParams(&channelIDs, &channels)
	insertAllChannels(channelParams, config)

	videoParams := createAllVideoParams(&channelIDs, &channels, 5)
	insertAllVideos(videoParams, config)

	log.Println("Database Initialized!")
}

func createAllChannelParams(initDb *InitDB, channels *[]services.Channel) []database.InsertChannelParams {
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
		go createOneChannelParams(&(*channels)[idx], channelID.ChannelID, wg, ch)
	}

	wg.Wait()
	close(ch)
	return params
}

func createOneChannelParams(channel *services.Channel, channelID string, wg *sync.WaitGroup, ch chan<- database.InsertChannelParams) {
	defer wg.Done()
	param := database.InsertChannelParams{
		ID:            channelID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Name:          channel.Name,
		Url:           channel.URL,
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

func createAllVideoParams(initDb *InitDB, channels *[]services.Channel, numVideosPerChannel int) []database.InsertVideoParams {
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
		for _, video := range (*channels)[cidx].Videos {
			if atomic.LoadInt32(&validCounts[cidx]) >= int32(numVideosPerChannel) {
				break
			}
			wg.Add(1)
			go func(video services.Video, channelId string, idx int) {
				defer wg.Done()
				param := createOneVideoParam(&video, channelId)
				if param != nil {
					if atomic.AddInt32(&validCounts[idx], 1) <= int32(numVideosPerChannel) {
						ch <- *param
					}
				}
			}(video, channelID.ChannelID, cidx)
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
