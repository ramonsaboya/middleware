package main

import (
	"errors"
	"sync"
)

// ThreadSafeQueue Fixed size thread safe queue
type ThreadSafeQueue struct {
	queue chan interface{}
	cond  *sync.Cond
}

// NewThreadSafeQueue Instantiates a ThreadSafeQueue with given capacity
func NewThreadSafeQueue(capacity int) *ThreadSafeQueue {
	return &ThreadSafeQueue{
		make(chan interface{}, capacity),
		sync.NewCond(&sync.Mutex{}),
	}
}

// Push Pushes a new value at the tail of the queue
func (q *ThreadSafeQueue) Push(value interface{}) error {
	q.Lock()
	defer q.Unlock()

	return q.UnsafePush(value)
}

// UnsafePush Pushs a new value at the tail of the queue
// without acquiring locks
func (q *ThreadSafeQueue) UnsafePush(value interface{}) error {
	select {
	case q.queue <- value:
	default:
		return errors.New("queue is full")
	}

	q.Broadcast()

	return nil
}

// Pop Remove and returns value at the front of the queue
func (q *ThreadSafeQueue) Pop() (interface{}, error) {
	q.Lock()
	defer q.Unlock()

	return q.UnsafePop()
}

// UnsafePop Remove and returns value at the front of the queue
// without acquiring locks
func (q *ThreadSafeQueue) UnsafePop() (interface{}, error) {
	select {
	case value, ok := <-q.queue:
		if ok {
			return value, nil
		}
		return nil, errors.New("internal channel is closed")
	default:
		return nil, errors.New("empty queue")
	}
}

// Empty Returns whether the queue is empty
func (q *ThreadSafeQueue) Empty() bool {
	return q.Len() == 0
}

// Len Returns current size of the queue
func (q *ThreadSafeQueue) Len() int {
	return len(q.queue)
}

// Cap Returns maximum capacity of the queue
func (q *ThreadSafeQueue) Cap() int {
	return cap(q.queue)
}

// Lock Acquires queue's lock or wait until available
func (q *ThreadSafeQueue) Lock() {
	q.cond.L.Lock()
}

// Unlock Unlocks queue's lock
func (q *ThreadSafeQueue) Unlock() {
	q.cond.L.Unlock()
}

// Wait Wait for a signal or broadcast from the queue
func (q *ThreadSafeQueue) Wait() {
	q.cond.Wait()
}

// Signal Signals whichever coroutine is waiting the longest
func (q *ThreadSafeQueue) Signal() {
	q.cond.Signal()
}

// Broadcast Signals all waiting coroutines
func (q *ThreadSafeQueue) Broadcast() {
	q.cond.Broadcast()
}
