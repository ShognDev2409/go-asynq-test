package tasks

import (
   "context"
   "encoding/json"
   "github.com/hibiken/asynq"
   "github.com/ShognDev2409/go-asynq-test/echo-worker/status"
   "log/slog"
   "time"
)

// TaskTypeBuyTicket is the type identifier for buy-ticket tasks.
const TaskTypeBuyTicket = "buy:ticket"

// Payload holds the data for a buy-ticket task.
type Payload struct {
	Event string `json:"event"`
}

// NewHandler returns an asynq.HandlerFunc that processes buy-ticket tasks
// and updates the status manager.
func NewHandler(mgr *status.Manager) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		// Get unique task ID from context.
		id, _ := asynq.GetTaskID(ctx)

		// Mark task as pending.
		mgr.Set(id, status.Pending)

		// Mark task as processing and log.
		mgr.Set(id, status.Processing)
		slog.Info("task started", "id", id)

		// Decode payload and simulate work.
		var p Payload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			slog.Error("failed to decode payload", "error", err)
			mgr.Set(id, status.Failed)
			return err
		}
		time.Sleep(100 * time.Millisecond)

		// Mark task as done and log.
		mgr.Set(id, status.Done)
		slog.Info("task completed", "id", id)
		return nil
	}
}

// NewTask creates a new Asynq task for a buy-ticket job with the given payload.
func NewTask(p Payload) (*asynq.Task, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskTypeBuyTicket, data), nil
}
