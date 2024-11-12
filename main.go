package main

import (
	"elrek-system_GO/api"
	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	godotenv.Load("../.env")

	slog.Info("Starting API...")
	api.Api()
}
