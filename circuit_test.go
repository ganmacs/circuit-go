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

func (fc fakeClock) Now() time.Time {
	return fc.val
}

func NewFakeClock() fakeClock {
	return fakeClock{val: time.Unix(0, 0)}
}

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
	clock := NewFakeClock()
	cb.Clock = &clock

	cb.Run(failFun)
	if cb.state != open {
		t.Errorf("should be open")
	}

	cb.Run(failFun)
	if cb.state != open {
		t.Errorf("should be open")
	}

	// occur timeout
	oldStartIime := cb.startOpenTime
	clock.val = oldStartIime.Add(cb.timeout + 1*time.Second)

	cb.Run(sucessFun)
	if cb.state != closed {
		t.Errorf("should be closed")
	}
}

func TestCircuitFailAgain(t *testing.T) {
	cb := NewCircuitBreaker()
	clock := NewFakeClock()
	cb.Clock = &clock

	cb.Run(failFun)
	if cb.state != open {
		t.Errorf("should be open")
	}

	// occur timeout
	clock.val = cb.startOpenTime.Add(cb.timeout + 1*time.Second)
	oldstartTime := cb.startOpenTime

	cb.Run(failFun)
	if cb.state != open {
		t.Errorf("should be open")
	}

	if oldstartTime == cb.startOpenTime {
		t.Errorf("Set new startOpenTime when fail again")
	}
}
