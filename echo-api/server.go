package main

import (
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
	"os"
)

// RunServer sets up Echo routes and starts the HTTP server.
func RunServer(httpAddr, redisAddr string) {
	// Create an Asynq client for enqueuing tasks.
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// Initialize HTTP handler with Asynq client.
	handler := NewHandler(client)

	// Create a new Echo instance and register middleware.
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())

	// Register the buy-ticket endpoint.
	e.POST("/buy-ticket", handler.BuyTicket)

	// Start the HTTP server.
	slog.Info("starting HTTP server", "addr", httpAddr)
	if err := e.Start(httpAddr); err != nil {
		slog.Error("HTTP server failed", "error", err)
		os.Exit(1)
	}
}
