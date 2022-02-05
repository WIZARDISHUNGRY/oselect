package oselect

import (
	"testing"
)

func TestSend_DelayedEval(t *testing.T) {
	chan0 := make(chan int, 1)
	chan1 := make(chan int, 1)

	Send2(
		chan0, func() int { return 1 },
		chan1, func() int { t.Fatal("Never called"); return -1 },
	)

	if <-chan0 != 1 {
		t.Fatal("wrong value on channel")
	}
}
