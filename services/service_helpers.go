package services

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/siddhant-vij/YouTube-Video-Aggregator/services/api"
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

var apiKeys []string

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey1 := os.Getenv("YT_API_KEY_1")
	apiKey2 := os.Getenv("YT_API_KEY_2")
	apiKey3 := os.Getenv("YT_API_KEY_3")
	apiKey4 := os.Getenv("YT_API_KEY_4")
	apiKey5 := os.Getenv("YT_API_KEY_5")
	apiKey6 := os.Getenv("YT_API_KEY_6")

	apiKeys = []string{apiKey1, apiKey2, apiKey3, apiKey4, apiKey5, apiKey6}
}

func isYouTubeFeedServerWorking() bool {
	client := &http.Client{}
	url := "https://www.youtube.com/feeds/videos.xml?channel_id=UCBR8-60-B28hp2BmDPdntcQ"
	// YouTube's Official Channel - it's that simple?
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func GetChannelVideos(channelID string) (Channel, error) {
	if isYouTubeFeedServerWorking() {
		feedChannel, err := feed.GetFeed(channelID)
		if err != nil {
			return Channel{}, err
		}
		return convertFeedToChannel(&feedChannel), nil
	} else {
		for _, apiKey := range apiKeys {
			ytChannel, err := api.GetYtApiChannel(apiKey, channelID)
			if err != nil {
				if strings.Contains(err.Error(), "quota exhausted") {
					continue
				}
				return Channel{}, err
			}
			return convertYtToChannel(&ytChannel), nil
		}
	}
	return Channel{}, errors.New("no fetching error: rss feed server not working and all api keys exhausted in a single day")
}

func convertFeedToChannel(feedChannel *feed.FeedChannel) Channel {
	return Channel{
		Name:   feedChannel.Title,
		URL:    feedChannel.Link.Href,
		Videos: convertFeedEntriesToVideos(&feedChannel.Entries),
	}
}

func convertYtToChannel(ytChannel *api.YTApiChannel) Channel {
	return Channel{
		Name:   ytChannel.Name,
		URL:    ytChannel.URL,
		Videos: convertYtEntriesToVideos(&ytChannel.Videos),
	}
}

func convertFeedEntriesToVideos(feedVideos *[]feed.FeedVideo) []Video {
	var videos []Video
	for _, feedVideo := range *feedVideos {
		videos = append(videos, convertFeedEntryToVideo(&feedVideo))
	}
	return videos
}

func convertYtEntriesToVideos(ytVideos *[]api.YTApiVideo) []Video {
	var videos []Video
	for _, ytVideo := range *ytVideos {
		videos = append(videos, convertYtEntryToVideo(&ytVideo))
	}
	return videos
}

func convertFeedEntryToVideo(feedVideo *feed.FeedVideo) Video {
	return Video{
		Title:       feedVideo.Title,
		Description: utils.ShortenText(string(feedVideo.Media.Description)),
		ImageURL:    feedVideo.Media.Image.URL,
		Authors:     feedVideo.Author.Name,
		PublishedAt: feedVideo.Published,
		URL:         feedVideo.Link.Href,
		ViewCount:   utils.ShortenViewCount(feedVideo.Media.Community.Statistics.Views),
		StarRating:  feedVideo.Media.Community.StarRating.Average,
		StarCount:   utils.ShortenStarCount(feedVideo.Media.Community.StarRating.Count),
	}
}

func convertYtEntryToVideo(ytVideo *api.YTApiVideo) Video {
	video := ytVideo.Items[0].Snippet
	stats := ytVideo.Items[0].Statistics
	return Video{
		Title:       video.Title,
		Description: utils.ShortenText(video.Description),
		ImageURL:    video.Thumbnails.Standard.URL,
		Authors:     video.Authors,
		PublishedAt: video.PublishedAt,
		URL:         "https://www.youtube.com/watch?v=" + ytVideo.VideoId,
		ViewCount:   utils.ShortenViewCount(stats.ViewCount),
		StarRating:  "5.00",
		StarCount:   utils.ShortenStarCount(stats.StarCount),
	}
}
