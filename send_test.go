package oselect

import (
	"fmt"
	"testing"
)

// TestCompoundExample shows what the runtime cannot do. Why are you doing this?
func TestCompoundExample(t *testing.T) {
	t.Skip()

	chan0 := make(chan int)
	chan1 := make(chan int, 1)
	chan2 := make(chan int, 1)

	chan2 <- 1

	select {
	case chan0 <- <-chan2:
		t.Fatal("chan0 is unbuffered")
	case <-chan0:
		t.Fatal("nothing in chan0")
	case chan1 <- (<-chan2) + 1:
		fmt.Println("this would be cool")
	default:
		t.Fatal("We don't reach here either! We time out! Weird!")
	}
}

// TestCompoundExample2 shows what the runtime cannot do. Really why would you make your life this miserable?
func TestCompoundExample2(t *testing.T) {
	t.Skip()

	chan0 := make(chan int)
	chan1 := make(chan int, 1)
	chan2 := make(chan int, 1)

	chan2 <- 1

	select {
	case chan0 <- <-chan2:
		t.Fatal("chan0 is unbuffered")
	case chan1 <- <-chan2:
		fmt.Println("this would be cool but why")
	default:
		t.Fatal("We don't reach here either! We time out! Weird!")
	}
}
