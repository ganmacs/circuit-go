package circuit

import (
	"time"
)

const (
	open     int = iota
	closed   int = iota
	halfopen int = iota
)
const (
	DefaultTimeout = 60 * time.Second
)

type Configs struct {
	rate float64
}

type CircuitBreaker struct {
	state         int
	step          uint64
	bucket        Bucket
	timeout       time.Duration
	startOpenTime time.Time
	rate          float64
}

type operation func() error

func NewCircuitBreaker(config Configs) *CircuitBreaker {
	return &CircuitBreaker{
		state:         closed,
		step:          0,
		bucket:        *NewBucket(),
		timeout:       DefaultTimeout,
		startOpenTime: time.Now(),
	}
}

func (cb *CircuitBreaker) isTimeout() bool {
	return cb.timeout < time.Now().Sub(cb.startOpenTime)
}

func (cb *CircuitBreaker) HalfOpen() {
	cb.state = halfopen
}

func (cb *CircuitBreaker) Sucess() {
	if cb.state == halfopen {
		cb.state = closed
		cb.bucket.Reset()
		cb.bucket.Sucess()
	} else if cb.state == closed {
		cb.bucket.Sucess()
	}
}

func (cb *CircuitBreaker) exceedThreshold() bool {
	return cb.bucket.ConsecutiveRate() > cb.rate
}

func (cb *CircuitBreaker) Fail() {
	if cb.state == halfopen {
		cb.state = open
		cb.bucket.Fail()
		cb.startOpenTime = time.Now()
	} else if cb.state == closed {
		cb.bucket.Fail()
		if cb.exceedThreshold() {
			cb.state = closed
		}
	}
}

func (cb *CircuitBreaker) Run(op operation) error {
	if cb.state == open && cb.isTimeout() {
		cb.HalfOpen()
	}

	if cb.state == open {
		return nil // cb is open todo error type
	}

	err := op()

	if err == nil {
		cb.Sucess()
	} else {
		cb.Fail()
	}

	return err
}
