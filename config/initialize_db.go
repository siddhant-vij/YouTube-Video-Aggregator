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
	"github.com/siddhant-vij/RSS-Feed-Aggregator/database"
	"github.com/siddhant-vij/RSS-Feed-Aggregator/utils"
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
		os.Exit(1)
	}
	reader := json.NewDecoder(jsonFile)

	if err := reader.Decode(&channelIDs); err != nil {
		log.Fatal(err)
		os.Exit(1)
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

	channelFollowParams := createChannelFollowParams(&channelIDs, uuid.MustParse("25614f84-e0e3-4e48-bd5f-81d3ef36d8d3"))
	insertAllChannelFollows(channelFollowParams, config)

	videoParams := createVideoParams(&channelIDs, &feeds, 5)
	// -1 = All Videos for each channel
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
		os.Exit(1)
	}
}

func createChannelFollowParams(initDb *InitDB, user_id uuid.UUID) []database.InsertChannelFollowParams {
	params := make([]database.InsertChannelFollowParams, 0)
	for _, channelID := range initDb.Channels {
		var parameter = database.InsertChannelFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user_id,
			ChannelID: channelID.ChannelID,
		}
		params = append(params, parameter)
	}
	return params
}

func insertAllChannelFollows(params []database.InsertChannelFollowParams, config *ApiConfig) {
	wg := &sync.WaitGroup{}
	for _, insertChannelFollowParam := range params {
		wg.Add(1)
		go insertOneChannelFollow(insertChannelFollowParam, config, wg)
	}
	wg.Wait()
}

func insertOneChannelFollow(insertChannelFollowParam database.InsertChannelFollowParams, config *ApiConfig, wg *sync.WaitGroup) {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()
	defer wg.Done()
	err := config.DBQueries.InsertChannelFollow(context.TODO(), insertChannelFollowParam)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func createVideoParams(initDb *InitDB, feeds *[]utils.Feed, numVideosPerChannel int) []database.InsertVideoParams {
	params := make([]database.InsertVideoParams, 0)
	for cidx, channelID := range initDb.Channels {
		for eidx, entry := range (*feeds)[cidx].Entries {
			if eidx == numVideosPerChannel {
				break
			}
			var parameter = database.InsertVideoParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       entry.Title,
				Description: string(entry.Media.Description),
				ImageUrl:    entry.Media.Image.URL,
				Authors:     entry.Author.Name,
				PublishedAt: entry.Published,
				Url:         entry.Link.Href,
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
		os.Exit(1)
	}
}
