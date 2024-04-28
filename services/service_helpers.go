package services

import (
	"time"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/services/feed"
	"github.com/siddhant-vij/YouTube-Video-Aggregator/utils"
)

type Channel struct {
	Name   string
	URL    string
	Videos []Video
}

type Video struct {
	Title       string
	Description string
	ImageURL    string
	Authors     string
	PublishedAt time.Time
	URL         string
	ViewCount   string
	StarRating  string
	StarCount   string
}

func isYouTubeFeedServerWorking() bool {
	return true
}

func ValidateChannelURL(url string) bool {
	if isYouTubeFeedServerWorking() {
		return feed.ValidateFeedURL(url)
	} else {
		return false
	}
}

func convertFeedToChannel(feed *feed.FeedChannel) Channel {
	return Channel{
		Name:   feed.Title,
		URL:    feed.Link.Href,
		Videos: convertFeedEntriesToVideos(&feed.Entries),
	}
}

func convertFeedEntriesToVideos(entries *[]feed.FeedVideo) []Video {
	var videos []Video
	for _, entry := range *entries {
		videos = append(videos, convertFeedEntryToVideo(&entry))
	}
	return videos
}

func convertFeedEntryToVideo(entry *feed.FeedVideo) Video {
	return Video{
		Title:       entry.Title,
		Description: utils.ShortenText(string(entry.Media.Description)),
		ImageURL:    entry.Media.Image.URL,
		Authors:     entry.Author.Name,
		PublishedAt: entry.Published,
		URL:         entry.Link.Href,
		ViewCount:   utils.ShortenViewCount(entry.Media.Community.Statistics.Views),
		StarRating:  entry.Media.Community.StarRating.Average,
		StarCount:   utils.ShortenStarCount(entry.Media.Community.StarRating.Count),
	}
}

func GetChannelVideos(channelID string) (Channel, error) {
	if isYouTubeFeedServerWorking() {
		feed, err := feed.GetFeed(channelID)
		if err != nil {
			return Channel{}, err
		}
		return convertFeedToChannel(&feed), nil
	} else {
		return Channel{}, nil
	}
}
