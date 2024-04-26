package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/config"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/database"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/utils"
)

func AddChannelFollowForUser(config *config.ApiConfig, userId uuid.UUID, channelId string) error {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()

	var channelFollowParams = database.InsertChannelFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userId,
		ChannelID: channelId,
	}
	err := config.DBQueries.InsertChannelFollow(context.TODO(), channelFollowParams)
	if err != nil {
		return err
	}

	return nil
}

func RemoveChannelFollowForUser(config *config.ApiConfig, userId uuid.UUID, channelId string) error {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()

	err := config.DBQueries.DeleteChannelFollow(context.TODO(), database.DeleteChannelFollowParams{
		UserID:    userId,
		ChannelID: channelId,
	})
	if err != nil {
		return err
	}

	return nil
}

func AddChannelAndVideos(config *config.ApiConfig, channelId string) error {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()

	var ytFeedInputParams = utils.YouTubeFeedInputParams{
		Client:         &http.Client{},
		ChannelBaseURL: config.ChannelBaseURL,
	}

	feed, err := ytFeedInputParams.GetFeed(channelId, utils.FTChannel)
	if err != nil {
		panic(err)
	}

	var channelParams = database.InsertChannelParams{
		ID:            channelId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Name:          feed.Title,
		Url:           feed.Link.Href,
		LastFetchedAt: time.Now(),
	}

	err = config.DBQueries.InsertChannel(context.TODO(), channelParams)
	if err != nil {
		return err
	}

	videoParams := createAllVideoParams(&(feed.Entries), 5, channelId)
	insertAllVideos(videoParams, config)

	return nil
}

func createAllVideoParams(entries *[]utils.Entry, numVideos int, channelID string) []database.InsertVideoParams {
	var params []database.InsertVideoParams
	count := 0
	for _, entry := range *entries {
		if count >= numVideos {
			break
		}
		description := string(entry.Media.Description)
		if description == "" {
			continue
		}
		param := database.InsertVideoParams{
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
		params = append(params, param)
		count++
	}
	return params
}

func insertAllVideos(params []database.InsertVideoParams, config *config.ApiConfig) {
	for _, param := range params {
		err := config.DBQueries.InsertVideo(context.TODO(), param)
		if err != nil {
			log.Fatalf("Failed to insert video: %v", err)
		}
	}
}
