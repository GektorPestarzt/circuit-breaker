package cbreaker

import (
	"sync"
	"time"
)

type State int

const (
	StateClose State = iota
	StateOpen
	StateHalf
)

type CircuitBreaker struct {
	threshold int
	timeout   time.Duration

	mutex          sync.Mutex
	lastFailure    time.Time
	failureCounter int
	state          State
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold: threshold,
		timeout:   timeout,
		state:     StateClose,
	}
}

func (cb *CircuitBreaker) Execute(operation func() error) error {
	err := cb.beforeRequest()
	if err != nil {
		return err
	}

	err = operation()

	cb.afterRequest(err == nil)
	return err
}

func (cb *CircuitBreaker) beforeRequest() error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.state == StateOpen && time.Since(cb.lastFailure) > cb.timeout {
		cb.state = StateHalf
	}

	if cb.state == StateOpen {
		return ErrNotAvailable
	}
	return nil
}

func (cb *CircuitBreaker) afterRequest(isOk bool) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case StateClose:
		if !isOk {
			cb.failureCounter++
			if cb.failureCounter >= cb.threshold {
				cb.state = StateOpen
				cb.lastFailure = time.Now()
			}
		}
	case StateHalf:
		if !isOk {
			cb.state = StateOpen
			cb.lastFailure = time.Now()
			return
		}
		cb.state = StateClose
	default:
		return
	}
}
