package cbreaker

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
	"time"
)

var Error = errors.New("error")

func TestCipher(t *testing.T) {
	serviceFailure := func() error {
		return Error
	}
	serviceOk := func() error {
		return nil
	}
	times := 50
	timeout := 2 * time.Second

	Convey("Sync", t, func() {
		breaker := NewCircuitBreaker(times, timeout)
		Convey("One Failure", func() {
			err := breaker.Execute(serviceFailure)
			So(err, ShouldEqual, Error)
			So(breaker.state, ShouldEqual, StateClose)
		})

		Convey("Many Failure Close", func() {
			for i := 0; i < times-1; i++ {
				_ = breaker.Execute(serviceFailure)
			}

			So(breaker.state, ShouldEqual, StateClose)
		})

		Convey("Many Failures Open", func() {
			for i := 0; i < times; i++ {
				_ = breaker.Execute(serviceFailure)
			}

			So(breaker.state, ShouldEqual, StateOpen)
		})

		Convey("Ok In Open", func() {
			for i := 0; i < times; i++ {
				_ = breaker.Execute(serviceFailure)
			}

			time.Sleep(timeout + time.Second)
			err := breaker.Execute(serviceOk)
			So(breaker.state, ShouldEqual, StateClose)
			So(err, ShouldEqual, nil)
		})

		Convey("Failure In Open", func() {
			for i := 0; i < times; i++ {
				_ = breaker.Execute(serviceFailure)
			}

			time.Sleep(timeout + time.Second)
			err := breaker.Execute(serviceFailure)
			So(breaker.state, ShouldEqual, StateOpen)
			So(err, ShouldEqual, Error)
		})
	})

	Convey("Unsync", t, func() {
		breaker := NewCircuitBreaker(times, timeout)
		var wg sync.WaitGroup
		Convey("Many Failure Close", func() {
			for i := 0; i < times-1; i++ {
				wg.Add(1)
				go func() {
					_ = breaker.Execute(serviceFailure)
					wg.Done()
				}()
			}

			wg.Wait()
			So(breaker.state, ShouldEqual, StateClose)
		})

		Convey("Many Failures Open", func() {
			for i := 0; i < times; i++ {
				wg.Add(1)
				go func() {
					_ = breaker.Execute(serviceFailure)
					wg.Done()
				}()
			}

			wg.Wait()
			So(breaker.state, ShouldEqual, StateOpen)
		})

		Convey("Ok In Open", func() {
			for i := 0; i < times; i++ {
				wg.Add(1)
				go func() {
					_ = breaker.Execute(serviceFailure)
					wg.Done()
				}()
			}

			wg.Wait()
			time.Sleep(timeout + time.Second)
			err := breaker.Execute(serviceOk)
			So(breaker.state, ShouldEqual, StateClose)
			So(err, ShouldEqual, nil)
		})

		Convey("Failure In Open", func() {
			for i := 0; i < times; i++ {
				wg.Add(1)
				go func() {
					_ = breaker.Execute(serviceFailure)
					wg.Done()
				}()
			}

			wg.Wait()
			time.Sleep(timeout + time.Second)
			err := breaker.Execute(serviceFailure)
			So(breaker.state, ShouldEqual, StateOpen)
			So(err, ShouldEqual, Error)
		})
	})
}
