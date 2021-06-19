package main

import (
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

func main() {
	// load .env files if set in directory
	_ = godotenv.Load()

	// config server and listen
	httpServer(getServerPort(), isDebug())
}
