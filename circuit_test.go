package circuit

import (
	"errors"
	"testing"
	"time"
)

var (
	sucessFun = func() error { return nil }
	failFun   = func() error { return errors.New("error") }
)

// clock mock
type fakeClock struct {
	val time.Time
}

func (t fakeClock) Now() time.Time {
	return t.val
}

var clockMock = fakeClock{}

func TestCircuit(t *testing.T) {
	cb := NewCircuitBreaker()

	cb.Run(sucessFun)
	cb.Run(sucessFun)

	if cb.state != closed {
		t.Errorf("should be close")
	}

	cb.Run(failFun)
	cb.Run(failFun)
	if cb.state != closed {
		t.Errorf("should be close")
	}

	err := cb.Run(failFun)
	if err == nil {
		t.Errorf("should be return error class")
	}

	if err.Error() == "CircuitBreaker is open" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestCircuitBackToSucess(t *testing.T) {
	cb := NewCircuitBreaker()

	cb.Run(failFun)
	if cb.state != open {
		t.Errorf("should be open")
	}

	cb.Run(failFun)
	if cb.state != open {
		t.Errorf("should be open")
	}

	// occur timeout
	clockMock.val = cb.startOpenTime.Add(cb.timeout + 1*time.Second)
	cb.Clock = clockMock

	cb.Run(sucessFun)
	if cb.state != closed {
		t.Errorf("should be closed")
	}
}
