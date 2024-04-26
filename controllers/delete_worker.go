package controllers

import (
	"context"
	"log"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/config"
)

func DeleteOldVideos(config *config.ApiConfig) {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()

	err := config.DBQueries.DeleteOldVideos(context.TODO())
	if err != nil {
		log.Fatalf("Failed to delete old videos: %v", err)
	}
}
