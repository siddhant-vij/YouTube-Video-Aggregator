package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/siddhant-vij/RSS-Feed-Aggregator/database"
)

func InitializeDB(config *ApiConfig) {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()

	var insertFeedParams = database.InsertFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      "Feed name",
		Url:       "Feed url",
	}
	err := config.DBQueries.InsertFeed(context.TODO(), insertFeedParams)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var insertFeedFollowParams = database.InsertFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    uuid.MustParse("4f78ccdb-9b91-4a26-be73-d0acc5872abb"),
		FeedID:    insertFeedParams.ID,
	}
	err = config.DBQueries.InsertFeedFollow(context.TODO(), insertFeedFollowParams)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var insertPostParams = database.InsertPostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       "Post title",
		ImageTitle:  "Image alt title",
		ImageUrl:    "Image url",
		Authors:     "Author1, Author2",
		PublishedAt: time.Now(),
		Description: "Post description",
		Url:         "Post url",
		FeedID:      insertFeedParams.ID,
	}
	err = config.DBQueries.InsertPost(context.TODO(), insertPostParams)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
