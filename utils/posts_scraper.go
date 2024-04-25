package utils

import (
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type YouTubeFeedInputParams struct {
	Client          *http.Client
	ChannelBaseURL  string
	PlaylistBaseURL string
}

type Type string

const (
	FTDefault  = Type("")
	FTChannel  = Type("channel")
	FTPlaylist = Type("playlist")
)

type Feed struct {
	Title string `xml:"title"`
	Link  struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Entries []Entry `xml:"entry"`
}

func (f Feed) String() string {
	return fmt.Sprintf("{Title:%q, Link:%s, Entries:%d}", f.Title, f.Link.Href, len(f.Entries))
}

type Entry struct {
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
	} `xml:"http://search.yahoo.com/mrss/ group"`

	Author struct {
		Name string `xml:"name"`
	} `xml:"author"`
}

func (e Entry) String() string {
	tz, _ := time.LoadLocation("Local")

	return fmt.Sprintf("{Title:%q, Link:%s, Published:%s, Author:%s, Description:%s, Image:%s}",
		e.Title, e.Link.Href, e.Published.In(tz).Format(time.RFC3339), e.Author.Name, e.Media.Description, e.Media.Image.URL,
	)
}

func (yt *YouTubeFeedInputParams) url(id string, feedType Type) (string, error) {
	switch feedType {
	case FTChannel, FTDefault:
		return yt.ChannelBaseURL + id, nil
	case FTPlaylist:
		return yt.PlaylistBaseURL + id, nil
	}
	return "", errors.New("unknown feed type")
}

func (yt *YouTubeFeedInputParams) GetFeedPosts(id string, feedType Type) (Feed, error) {
	url, err := yt.url(id, feedType)
	if err != nil {
		return Feed{}, err
	}

	if ValidateFeedURL(url) {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return Feed{}, err
		}

		response, err := yt.Client.Do(request)
		if err != nil {
			return Feed{}, err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return Feed{}, fmt.Errorf("status code error: %d %s", response.StatusCode, response.Status)
		}

		var feed Feed
		if err := xml.NewDecoder(response.Body).Decode(&feed); err != nil {
			return Feed{}, err
		}

		// sort.Slice(feed.Entries, func(i, j int) bool {
		// 	return feed.Entries[i].Published.After(feed.Entries[j].Published)
		// })

		return feed, nil
	}

	return Feed{}, errors.New("not a valid youtube channel/playlist id")
}
