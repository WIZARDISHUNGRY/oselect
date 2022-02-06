package oselect

import (
	"testing"
)

func nothing(int) { panic("i shall never be called") }
func nada()       {}

func BenchmarkRecv4Default(b *testing.B) {
	chan0 := make(chan int)
	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int)

	for n := 0; n < b.N; n++ {
		Recv4Default(
			chan0, nothing,
			chan1, nothing,
			chan2, nothing,
			chan3, nothing,
			nada,
		)
	}
}

func BenchmarkSelect4Default_Recv(b *testing.B) {
	chan0 := make(chan int)
	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int)

	for n := 0; n < b.N; n++ {
		Select4Default(
			Recv(chan0, nothing),
			Recv(chan1, nothing),
			Recv(chan2, nothing),
			Recv(chan3, nothing),
			nada,
		)
	}
}

func BenchmarkSelect4Default_Recv_preroll(b *testing.B) {
	chan0 := make(chan int)
	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int)
	r0 := Recv(chan0, nothing)
	r1 := Recv(chan1, nothing)
	r2 := Recv(chan2, nothing)
	r3 := Recv(chan3, nothing)

	for n := 0; n < b.N; n++ {
		Select4Default(
			r0, r1, r2, r3,
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

func BenchmarkRecv4(b *testing.B) {
	chan0 := make(chan int)
	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int, 1)

	for n := 0; n < b.N; n++ {

		chan3 <- 1

		Recv4(
			chan0, nothing,
			chan1, nothing,
			chan2, nothing,
			chan3, func(int) {},
		)
	}
}

func Benchmark_select_4(b *testing.B) {
	chan0 := make(chan int)
	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int, 1)
	f3 := func(int) {}

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
			f3(v3)
		}
	}
}
