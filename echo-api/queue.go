package main

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

// TaskTypeBuyTicket is the type identifier for buy-ticket tasks.
const TaskTypeBuyTicket = "buy:ticket"

// BuyTicketPayload defines the payload for buy-ticket tasks.
type BuyTicketPayload struct {
	Event string `json:"event"`
}

// NewTask creates a new Asynq task for buying a ticket with the given payload.
func NewTask(payload BuyTicketPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskTypeBuyTicket, data), nil
}
