package controllers

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/config"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/database"
)

type FullResponseForUser struct {
	Videos           []database.GetUserVideosRow
	ChannelsFollowed []database.Channel
	OtherChannels    []database.Channel
}

func GetAllChannelFollowsForUser(config *config.ApiConfig, userId uuid.UUID) ([]database.ChannelFollow, error) {
	config.Mutex.RLock()
	defer config.Mutex.RUnlock()

	channelFollows, err := config.DBQueries.GetChannelFollowsForUser(context.TODO(), userId)
	if err != nil {
		return nil, err
	}

	return channelFollows, nil
}

func InsertNumFollowsForUser(config *config.ApiConfig, userId uuid.UUID, numChannels int32) error {
	channels, err := getNumChannels(config, numChannels)
	if err != nil {
		return err
	}

	config.Mutex.Lock()
	defer config.Mutex.Unlock()

	for _, channel := range channels {
		var channelFollowParams = database.InsertChannelFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    userId,
			ChannelID: channel.ID,
		}
		err := config.DBQueries.InsertChannelFollow(context.TODO(), channelFollowParams)
		if err != nil {
			return err
		}
	}

	return nil
}

func getNumChannels(config *config.ApiConfig, numChannels int32) ([]database.Channel, error) {
	config.Mutex.RLock()
	defer config.Mutex.RUnlock()

	channels, err := config.DBQueries.GetNumChannelsByCreatedAt(context.TODO(), numChannels)
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func GetResponseVideosForUser(config *config.ApiConfig, userId uuid.UUID, numVideos int32) ([]database.GetUserVideosRow, error) {
	config.Mutex.RLock()
	defer config.Mutex.RUnlock()

	var userVideoParams = database.GetUserVideosParams{
		UserID: userId,
		Limit:  numVideos,
	}

	videos, err := config.DBQueries.GetUserVideos(context.TODO(), userVideoParams)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func GetFollowedChannelsForUser(config *config.ApiConfig, userId uuid.UUID) ([]database.Channel, error) {
	config.Mutex.RLock()
	defer config.Mutex.RUnlock()

	channels, err := config.DBQueries.GetUserFollowedChannels(context.TODO(), userId)
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func GetNotYetFollowedChannelsForUser(config *config.ApiConfig, userId uuid.UUID) ([]database.Channel, error) {
	config.Mutex.RLock()
	defer config.Mutex.RUnlock()

	channels, err := config.DBQueries.GetOtherChannelsForUser(context.TODO(), userId)
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func GenerateResponseForUser(config *config.ApiConfig, userId uuid.UUID, numFollowChannels, numVideos int32) (FullResponseForUser, error) {
	channelFollows, err := GetAllChannelFollowsForUser(config, userId)
	if err != nil {
		return FullResponseForUser{}, err
	}
	if len(channelFollows) == 0 {
		InsertNumFollowsForUser(config, userId, numFollowChannels)
	}
	videos, err := GetResponseVideosForUser(config, userId, numVideos)
	if err != nil {
		return FullResponseForUser{}, err
	}

	response := FullResponseForUser{}
	response.Videos = videos
	response.ChannelsFollowed, _ = GetFollowedChannelsForUser(config, userId)
	response.OtherChannels, _ = GetNotYetFollowedChannelsForUser(config, userId)
	return response, nil
}
