package circuit

import (
	"errors"
	"time"
)

const (
	open     int = iota
	closed   int = iota
	halfopen int = iota
)
const (
	DefaultTimeout = 60 * time.Second
	DefaultRate    = 0.5
)

type Configs struct {
	rate float64
}

type Clock interface {
	Now() time.Time
}

type CircuitBreaker struct {
	state         int
	step          uint64
	bucket        Bucket
	timeout       time.Duration
	startOpenTime time.Time
	rate          float64
	Clock         Clock
}

type RealClock struct{}

func (t RealClock) Now() time.Time {
	return time.Now()
}

var clock = RealClock{}

type operation func() error

func NewCircuitBreaker() *CircuitBreaker {
	return &CircuitBreaker{
		state:         closed,
		step:          0,
		bucket:        *NewBucket(),
		timeout:       DefaultTimeout,
		startOpenTime: time.Now(),
		rate:          DefaultRate,
		Clock:         clock,
	}
}

func (cb *CircuitBreaker) isTimeout() bool {
	return cb.timeout < cb.Clock.Now().Sub(cb.startOpenTime)
}

func (cb *CircuitBreaker) HalfOpen() {
	cb.state = halfopen
}

func (cb *CircuitBreaker) Sucess() {
	switch cb.state {
	case halfopen:
		cb.bucket.Sucess()
		cb.state = closed
	case closed:
		cb.bucket.Sucess()
	}
}

func (cb *CircuitBreaker) exceedThreshold() bool {
	return cb.bucket.ConsecutiveRate() > cb.rate
}

func (cb *CircuitBreaker) Fail() {
	switch cb.state {
	case halfopen:
		cb.state = open
		cb.bucket.Fail()
		cb.startOpenTime = cb.Clock.Now()
	case closed:
		cb.bucket.Fail()
		if cb.exceedThreshold() {
			cb.state = open
		}
	}
}

func (cb *CircuitBreaker) Run(op operation) error {
	if cb.state == open && cb.isTimeout() {
		cb.HalfOpen()
	}

	if cb.state == open {
		return errors.New("CircuitBreaker is open")
	}

	err := op()

	if err == nil {
		cb.Sucess()
	} else {
		cb.Fail()
	}

	return err
}
