package oselect

import (
	"testing"
)

func nothing(int) {}
func nada()       {}

func BenchmarkSelect4Default(b *testing.B) {
	chan0 := make(chan int)
	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int)

	for n := 0; n < b.N; n++ {
		Select4Default(
			chan0, nothing,
			chan1, nothing,
			chan2, nothing,
			chan3, nothing,
			nada,
		)
	}
}
func Benchmark_select_4_default(b *testing.B) {

	chan0 := make(chan int)
	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int)

	for n := 0; n < b.N; n++ {

		select {
		case v0 := <-chan0:
			nothing(v0)
			return
		default:
		}
		select {
		case v0 := <-chan0:
			nothing(v0)
			return
		case v1 := <-chan1:
			nothing(v1)
			return
		default:
		}
		select {
		case v0 := <-chan0:
			nothing(v0)
			return
		case v1 := <-chan1:
			nothing(v1)
			return
		case v2 := <-chan2:
			nothing(v2)
			return
		default:
		}
		select {
		case v0 := <-chan0:
			nothing(v0)
		case v1 := <-chan1:
			nothing(v1)
		case v2 := <-chan2:
			nothing(v2)
		case v3 := <-chan3:
			nothing(v3)
		default:
			nada()
		}
	}
}

func BenchmarkSelect4(b *testing.B) {
	chan0 := make(chan int)
	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int, 1)

	for n := 0; n < b.N; n++ {

		chan3 <- 1

		Select4(
			chan0, nothing,
			chan1, nothing,
			chan2, nothing,
			chan3, nothing,
		)
	}
}
func Benchmark_select_4(b *testing.B) {
	chan0 := make(chan int)
	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int, 1)

	for n := 0; n < b.N; n++ {

		chan3 <- 1

		select {
		case v0 := <-chan0:
			nothing(v0)
			return
		default:
		}
		select {
		case v0 := <-chan0:
			nothing(v0)
			return
		case v1 := <-chan1:
			nothing(v1)
			return
		default:
		}
		select {
		case v0 := <-chan0:
			nothing(v0)
			return
		case v1 := <-chan1:
			nothing(v1)
			return
		case v2 := <-chan2:
			nothing(v2)
			return
		default:
		}
		select {
		case v0 := <-chan0:
			nothing(v0)
		case v1 := <-chan1:
			nothing(v1)
		case v2 := <-chan2:
			nothing(v2)
		case v3 := <-chan3:
			nothing(v3)
		}
	}
}
