package storage

import (
	"sync"
)

// Queue is a simple thread-safe FIFO queue used to
// enqueue telemetry data for processing.
type Queue struct {
	mu    sync.Mutex
	items []any
}

func NewQueue() *Queue {
	return &Queue{items: make([]any, 0)}
}

// Enqueue a telemetry item to the queue
func (q *Queue) Enqueue(item any) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, item)

	return nil
}

// Dequeue a telemetry item from the queue
func (q *Queue) Dequeue() any {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return nil
	}

	item := q.items[0]
	q.items = q.items[1:]

	return item
}
