package main

import (
	"log"
	"net/http"

	"github.com/hardzal/go-websocket/routes"
)

func main() {

	mux := http.NewServeMux()
	routes.RegisterWebSocketRoutes(mux)

	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
