package controllers

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/config"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/database"
)

type BookmarkedResponseForUser struct {
	Videos           []database.GetVideosBookmarkedByUserRow
	ChannelsFollowed []database.Channel
	OtherChannels    []database.Channel
}

func AddBookmarkToVideoForUser(config *config.ApiConfig, userId uuid.UUID, videoId uuid.UUID) (bool, error) {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()

	var insertBookmarkParams = database.InsertBookmarkParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userId,
		VideoID:   videoId,
	}
	err := config.DBQueries.InsertBookmark(context.TODO(), insertBookmarkParams)
	if err != nil {
		return false, err
	}
	return true, nil
}

func RemoveBookmarkFromVideoForUser(config *config.ApiConfig, userId uuid.UUID, videoId uuid.UUID) (bool, error) {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()

	var deleteBookmarkParams = database.DeleteBookmarkParams{
		UserID:  userId,
		VideoID: videoId,
	}
	err := config.DBQueries.DeleteBookmark(context.TODO(), deleteBookmarkParams)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetBookmarkedVideosForUser(config *config.ApiConfig, userId uuid.UUID) (BookmarkedResponseForUser, error) {
	config.Mutex.RLock()
	defer config.Mutex.RUnlock()

	videos, err := config.DBQueries.GetVideosBookmarkedByUser(context.TODO(), userId)
	if err != nil {
		return BookmarkedResponseForUser{}, err
	}

	response := BookmarkedResponseForUser{}
	response.Videos = videos
	response.ChannelsFollowed, _ = GetFollowedChannelsForUser(config, userId)
	response.OtherChannels, _ = GetNotYetFollowedChannelsForUser(config, userId)
	return response, nil
}
