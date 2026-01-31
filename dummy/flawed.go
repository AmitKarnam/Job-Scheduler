package dummy

import (
	"fmt"
	"time"
)

// FlawedWorker is intentionally flawed to provide edge-cases for code review.
// Problems intentionally included:
// - `quit` channel is not initialized in constructor -> nil channel send/receive panic
// - Busy-loop with tiny sleep causing high CPU usage
// - Data race on `count` (no synchronization)
// - Stop() sends to channel without checking nil or closed
type FlawedWorker struct {
	quit  chan bool
	count int
}

// NewFlawedWorker returns a FlawedWorker but forgets to initialize the channel.
func NewFlawedWorker() *FlawedWorker {
	return &FlawedWorker{count: 0}
}

// Start launches the worker goroutine. It spins in a busy loop and modifies
// `count` without synchronization.
func (w *FlawedWorker) Start() {
	go func() {
		for {
			select {
			case <-w.quit:
				fmt.Println("FlawedWorker stopping")
				return
			default:
				// Busy work and unsynchronized increment
				w.count++
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()
}

// Stop attempts to signal the goroutine to stop but may panic if `quit` is nil
// or closed. No idempotency protection.
func (w *FlawedWorker) Stop() {
	w.quit <- true
}

// Count returns the current counter. This is racy when called concurrently.
func (w *FlawedWorker) Count() int {
	return w.count
}

// Example helper that intentionally leaks a goroutine by sleeping forever in a goroutine.
func StartLeakyTimer() {
	go func() {
		// This goroutine never exits; intended as a flaw for review.
		time.Sleep(100 * 365 * 24 * time.Hour)
	}()
}
