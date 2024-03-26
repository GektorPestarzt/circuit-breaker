package main

import (
	"cbreaker/pkg/cbreaker"
	"errors"
	"sync"
	"time"
)

func main() {
	service := func() error {
		return errors.New("error")
	}
	timeout := 3 * time.Second
	times := 10
	breaker := cbreaker.NewCircuitBreaker(times, timeout)

	var wg sync.WaitGroup
	for i := 0; i < times-1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = breaker.Execute(service)
		}()
	}

	wg.Wait()
}
