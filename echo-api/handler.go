package main

import (
	"context"
	// "encoding/json"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	client *asynq.Client
}

// NewHandler returns a new Handler with the given Asynq client.
func NewHandler(client *asynq.Client) *Handler {
	return &Handler{client: client}
}

// BuyTicket handles HTTP requests to enqueue a buy-ticket task.
func (h *Handler) BuyTicket(c echo.Context) error {
	// Decode request payload (expect JSON with Event field).
	var payload BuyTicketPayload
	if err := c.Bind(&payload); err != nil {
		slog.Error("invalid request payload", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	// Log receipt of request.
	slog.Info("received buy-ticket request", "payload", payload)

	// Create the task.
	task, err := NewTask(payload)
	if err != nil {
		slog.Error("failed to create task", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	// Enqueue the task.
	info, err := h.client.EnqueueContext(context.Background(), task)
	if err != nil {
		slog.Error("failed to enqueue task", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "enqueue failure"})
	}

	// Log successful enqueue.
	slog.Info("enqueued task", "id", info.ID, "queue", info.Queue)

	// Return task ID to the client.
	return c.JSON(http.StatusAccepted, map[string]string{"id": info.ID})
}
