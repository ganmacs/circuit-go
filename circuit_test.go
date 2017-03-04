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

	cb.Run(failFun)
	if cb.state != open {
		t.Errorf("should be close")
	}
}
