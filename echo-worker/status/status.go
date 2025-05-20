package status

import "sync"

// Status represents the current state of a task.
type Status string

const (
   // Pending indicates task is queued but not yet processed.
   Pending Status = "pending"
   // Processing indicates task is currently being processed.
   Processing Status = "processing"
   // Done indicates task has completed successfully.
   Done Status = "done"
   // Failed indicates task processing failed.
   Failed Status = "failed"
)

// Manager tracks task statuses in memory.
type Manager struct {
   mu sync.RWMutex
   m  map[string]Status
}

// NewManager creates a new Manager instance.
func NewManager() *Manager {
   return &Manager{m: make(map[string]Status)}
}

// Set updates the status for a given task ID.
func (mgr *Manager) Set(id string, s Status) {
   mgr.mu.Lock()
   defer mgr.mu.Unlock()
   mgr.m[id] = s
}

// Get retrieves the status for a given task ID.
func (mgr *Manager) Get(id string) Status {
   mgr.mu.RLock()
   defer mgr.mu.RUnlock()
   return mgr.m[id]
}