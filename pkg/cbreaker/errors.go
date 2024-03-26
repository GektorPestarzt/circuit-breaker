package cbreaker

import "errors"

var ErrNotAvailable = errors.New("circuit breaker: service is not available")
