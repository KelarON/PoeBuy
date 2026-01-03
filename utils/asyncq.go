package utils

import (
	"sync"
)

// AsyncQueue is a thread-safe queue that can hold a fixed number of items.
type AsyncQueue[T any] struct {
	queueChan  chan *T
	queueMutex *sync.Mutex
	queueSize  int
}

// NewAsyncQueue creates a new AsyncQueue with the specified queue size.
func NewAsyncQueue[T any](queueSize int) *AsyncQueue[T] {
	return &AsyncQueue[T]{
		queueChan:  make(chan *T, queueSize),
		queueMutex: &sync.Mutex{},
		queueSize:  queueSize,
	}
}

// Pop removes and returns the oldest item from the queue. If the queue is empty, it blocks until an item is available.
func (q *AsyncQueue[T]) Pop() *T {
	item := <-q.queueChan
	return item
}

// Push adds an item to the queue. If the queue is full, it removes the oldest item before adding the new one.
func (q *AsyncQueue[T]) Push(item *T) {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()
	if len(q.queueChan) == q.queueSize {
		// Remove the oldest item if the queue is full
		_ = <-q.queueChan
	}
	q.queueChan <- item
}
