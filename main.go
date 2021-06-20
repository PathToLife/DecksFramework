package main

import (
	"decksframework/decks"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

func getServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
}

func isDebug() bool {
	debugStr := strings.ToLower(os.Getenv("DEBUG"))
	return debugStr == "true"
}

var deckStore decks.DeckStore

func main() {
	// load .env files if set in directory
	_ = godotenv.Load()

	// init in memory deck storage, replace with database call in future
	deckStore = decks.NewDeckStore()

	// start server and listen
	startApiServer(getServerPort(), isDebug())
}
