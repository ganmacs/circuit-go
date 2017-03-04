package circuit

import (
	"errors"

	"testing"
)

var (
	sucessFun = func() error { return nil }
	failFun   = func() error { return errors.New("error") }
)

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
