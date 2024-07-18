package main

import (
	"elrek-system_GO/api"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

func main() {
	godotenv.Load(".env")
	godotenv.Load("../.env")
	dbName := os.Getenv("DB_NAME")

	var logfile string

	if dbName == "elrek-system_prod" {
		logfile = "logs/prod.log"
	} else {
		logfile = "logs/dev.log"
	}

	// w := os.Stdout

	file, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed to open log file: %v", err)
	}
	defer file.Close()

	logger := slog.New(
		tint.NewHandler(file, &tint.Options{
			NoColor: !isatty.IsTerminal(file.Fd()),
		}),
	)
	slog.SetDefault(logger)

	slog.Info("Starting API...")
	api.Api()
}
