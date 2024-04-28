package feed

import (
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

const channelBaseURL = "https://www.youtube.com/feeds/videos.xml?channel_id="

var client *http.Client = &http.Client{}

type FeedChannel struct {
	Title string `xml:"title"`
	Link  struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Entries []FeedVideo `xml:"entry"`
}

type FeedVideo struct {
	Title string `xml:"title"`
	Link  struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Published time.Time `xml:"published"`

	Media struct {
		Description template.HTML `xml:"description"`
		Image       struct {
			URL string `xml:"url,attr"`
		} `xml:"thumbnail"`
		Community struct {
			StarRating struct {
				Count   string `xml:"count,attr"`
				Average string `xml:"average,attr"`
			} `xml:"starRating"`
			Statistics struct {
				Views string `xml:"views,attr"`
			} `xml:"statistics"`
		} `xml:"community"`
	} `xml:"http://search.yahoo.com/mrss/ group"`

	Author struct {
		Name string `xml:"name"`
	} `xml:"author"`
}

func GetFeed(channelId string) (FeedChannel, error) {
	url := channelBaseURL + channelId

	if ValidateFeedURL(url) {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return FeedChannel{}, err
		}

		response, err := client.Do(request)
		if err != nil {
			return FeedChannel{}, err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return FeedChannel{}, fmt.Errorf("status code error: %d %s", response.StatusCode, response.Status)
		}

		var feed FeedChannel
		if err := xml.NewDecoder(response.Body).Decode(&feed); err != nil {
			return FeedChannel{}, err
		}

		return feed, nil
	}

	return FeedChannel{}, errors.New("not a valid youtube channel id")
}
