package main

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ShognDev2409/go-asynq-test/echo-worker/status"
	"github.com/ShognDev2409/go-asynq-test/echo-worker/tasks"
	"github.com/hibiken/asynq"
)

// TestConcurrencyLimit verifies that no more than the configured number of tasks run concurrently.
func TestConcurrencyLimit(t *testing.T) {
	// Status manager to track task status.
	mgr := status.NewManager()

	var current int32
	var maxConcurrent int32

	// Handler tracks concurrency and updates statuses.
	handler := func(ctx context.Context, t *asynq.Task) error {
		id, _ := asynq.GetTaskID(ctx)

		mgr.Set(id, status.Pending)

		now := atomic.AddInt32(&current, 1)
		for {
			prevMax := atomic.LoadInt32(&maxConcurrent)
			if now > prevMax && atomic.CompareAndSwapInt32(&maxConcurrent, prevMax, now) {
				break
			}
			if now <= prevMax {
				break
			}
		}

		mgr.Set(id, status.Processing)
		time.Sleep(50 * time.Millisecond)
		atomic.AddInt32(&current, -1)

		mgr.Set(id, status.Done)
		return nil
	}

	// Start Asynq server with a concurrency limit of 10.
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "127.0.0.1:6370"},
		asynq.Config{Concurrency: 10},
	)
	mux := asynq.NewServeMux()
	mux.Handle(tasks.TaskTypeBuyTicket, asynq.HandlerFunc(handler))

	go func() {
		if err := srv.Run(mux); err != nil {
			t.Errorf("server error: %v", err)
		}
	}()
	// Give the server time to start.
	time.Sleep(100 * time.Millisecond)

	// Enqueue 20 tasks.
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "127.0.0.1:6370"})
	var ids []string
	for i := 0; i < 20; i++ {
		payload := tasks.Payload{Event: "concert"}
		task, err := tasks.NewTask(payload)
		if err != nil {
			t.Fatalf("task creation failed: %v", err)
		}
		info, err := client.EnqueueContext(context.Background(), task)
		if err != nil {
			t.Fatalf("enqueue failed: %v", err)
		}
		ids = append(ids, info.ID)
	}

	// Wait for all tasks to complete (or timeout).
	deadline := time.Now().Add(5 * time.Second)
	for _, id := range ids {
		for {
			if mgr.Get(id) == status.Done {
				break
			}
			if time.Now().After(deadline) {
				t.Fatalf("timeout waiting for task %s to complete", id)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}

	// Verify the concurrency limit was not exceeded.
	if maxConcurrent > 10 {
		t.Errorf("max concurrent tasks = %d; want <= 10", maxConcurrent)
	}

	// Verify all tasks completed successfully.
	for _, id := range ids {
		if mgr.Get(id) != status.Done {
			t.Errorf("task %s status = %s; want done", id, mgr.Get(id))
		}
	}
}
