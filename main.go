package main

import (
	"decksframework/db"
	"decksframework/server"
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

	// init in memory deck storage, replace with database call in future
	db.InitDb()

	// start server and listen
	server.StartApiServer(getServerPort(), isDebug())
}
