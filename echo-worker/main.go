package main

import (
	"flag"
	"log/slog"
	"os"
)

// main parses flags, configures logger, and starts the worker.
func main() {
	// Parse command-line flags for Redis address and concurrency.
	redisAddr := flag.String("redis", "127.0.0.1:6370", "Redis server address")
	concurrency := flag.Int("concurrency", 10, "Maximum concurrent tasks")
	flag.Parse()

	// Initialize structured JSON logger and set as default.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Run the Asynq worker with given settings.
	RunWorker(*redisAddr, *concurrency)
}
