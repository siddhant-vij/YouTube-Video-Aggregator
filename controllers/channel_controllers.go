package controllers

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/config"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/database"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/services"
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

	channel, err := services.GetChannelVideos(channelId)
	if err != nil {
		panic(err)
	}

	var channelParams = database.InsertChannelParams{
		ID:            channelId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Name:          channel.Name,
		Url:           channel.URL,
		LastFetchedAt: time.Now(),
	}

	err = config.DBQueries.InsertChannel(context.TODO(), channelParams)
	if err != nil {
		return err
	}

	videoParams := createAllVideoParams(&(channel.Videos), 5, channelId)
	insertAllVideos(videoParams, config)

	return nil
}

func createAllVideoParams(videos *[]services.Video, numVideos int, channelID string) []database.InsertVideoParams {
	var params []database.InsertVideoParams
	count := 0
	for _, video := range *videos {
		if count >= numVideos {
			break
		}
		if video.Description == "..." || video.Description == "" {
			continue
		}
		param := database.InsertVideoParams{
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
