package main

import (
	// "log"
	// "net/http"

	// "github.com/siddhant-vij/RSS-Feed-Aggregator/middlewares"
	// "github.com/siddhant-vij/RSS-Feed-Aggregator/routes"

	"fmt"
	"net/http"

	"github.com/siddhant-vij/RSS-Feed-Aggregator/utils"
)

// func main() {
// 	mux := http.NewServeMux()
// 	corsMux := middlewares.CorsMiddleware(mux)
// 	routes.RegisterRoutes(mux)

// 	serverAddr := "localhost:" + routes.ResourceServerPort()
// 	log.Fatal(http.ListenAndServe(serverAddr, corsMux))
// }

func main() {
	// channelId := "UCUyeluBRhGPCW4rPe_UvBZQ"
	playlistId := "PLQnljOFTspQUNnO4p00ua_C5mKTfldiYT"
	var ytFeedInputParams = utils.YouTubeFeedInputParams{
		Client:          &http.Client{},
		ChannelBaseURL:  "https://www.youtube.com/feeds/videos.xml?channel_id=",
		PlaylistBaseURL: "https://www.youtube.com/feeds/videos.xml?playlist_id=",
	}

	feedPosts, err := ytFeedInputParams.GetFeedPosts(playlistId, utils.FTPlaylist)
	if err != nil {
		panic(err)
	}

	fmt.Println(feedPosts.Title)
	fmt.Println(feedPosts.Link.Href)

	for _, entry := range feedPosts.Entries {
		fmt.Println("-----------------------")
		fmt.Println(entry.Title)
		fmt.Println(entry.Media.Description)
		fmt.Println(entry.Media.Image.URL)
		fmt.Println(entry.Author.Name)
		fmt.Println(entry.Published)
		fmt.Println(entry.Link.Href)
	}
}
