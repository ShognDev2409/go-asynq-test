package main

import (
	"log/slog"
	"os"

	"github.com/ShognDev2409/go-asynq-test/echo-worker/status"
	"github.com/ShognDev2409/go-asynq-test/echo-worker/tasks"
	"github.com/hibiken/asynq"
)

// RunWorker sets up the Asynq server with a concurrency limit and handlers.
func RunWorker(redisAddr string, concurrency int) {
	// Create a new Asynq server with configured concurrency.
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{Concurrency: concurrency},
	)

	// Initialize ServeMux and status manager.
	mux := asynq.NewServeMux()
	statusMgr := status.NewManager()

	// Register the buy-ticket task handler.
	mux.Handle(tasks.TaskTypeBuyTicket, tasks.NewHandler(statusMgr))

	// Start processing tasks.
	slog.Info("starting worker", "redis", redisAddr, "concurrency", concurrency)
	if err := srv.Run(mux); err != nil {
		slog.Error("worker stopped", "error", err)
		os.Exit(1)
	}
}
