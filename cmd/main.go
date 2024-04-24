package main

import (
	"log"
	"net/http"

	"github.com/siddhant-vij/RSS-Feed-Aggregator/middlewares"
	"github.com/siddhant-vij/RSS-Feed-Aggregator/routes"
)

func main() {
	mux := http.NewServeMux()
	corsMux := middlewares.CorsMiddleware(mux)
	routes.RegisterRoutes(mux)

	serverAddr := "localhost:" + routes.ResourceServerPort()
	log.Fatal(http.ListenAndServe(serverAddr, corsMux))
}
