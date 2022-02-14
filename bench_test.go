package oselect

import (
	"reflect"
	"testing"
)

//go:noinline
func nothing(int) { panic("i shall never be called") }

//go:noinline
func nothingOK(int, bool) { panic("i shall never be called, ok?") }

//go:noinline
func nada() {}

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
			continue
		default:
		}
		select {
		case v0 := <-chan0:
			nothing(v0)
			continue
		case v1 := <-chan1:
			nothing(v1)
			continue
		default:
		}
		select {
		case v0 := <-chan0:
			nothing(v0)
			continue
		case v1 := <-chan1:
			nothing(v1)
			continue
		case v2 := <-chan2:
			nothing(v2)
			continue
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

func Benchmark_reflectDotSelect_4_default_preroll(b *testing.B) {
	chan0 := make(chan int)
	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int)
	def := reflect.SelectCase{Dir: reflect.SelectDefault}
	cases0 := []reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(chan0),
		},
		def,
	}
	cases1 := append([]reflect.SelectCase{}, cases0[0:1]...)
	cases1 = append(cases1, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(chan1),
	},
		def)
	cases2 := append([]reflect.SelectCase{}, cases1[0:2]...)
	cases2 = append(cases2, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(chan2),
	},
		def)
	cases3 := append([]reflect.SelectCase{}, cases2[0:3]...)
	cases3 = append(cases3, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(chan3),
	},
		def)

	for n := 0; n < b.N; n++ {
		chosen, v, _ := reflect.Select(cases0)
		if chosen <= 0 {
			nothing(int(v.Int()))
			continue
		}
		chosen, v, _ = reflect.Select(cases1)
		if chosen <= 1 {
			nothing(int(v.Int()))
			continue
		}
		chosen, v, _ = reflect.Select(cases2)
		if chosen <= 2 {
			nothing(int(v.Int()))
			continue
		}
		chosen, v, _ = reflect.Select(cases3)
		if chosen < 3 {
			nothing(int(v.Int()))
			continue
		}
		nada()
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

func BenchmarkSelect4_Recv(b *testing.B) {
	chan0 := make(chan int)
	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int, 1)

	for n := 0; n < b.N; n++ {

		chan3 <- 1
		Select4(
			Recv(chan0, nothing),
			Recv(chan1, nothing),
			Recv(chan2, nothing),
			Recv(chan3, func(int) {}),
		)
	}
}

func BenchmarkSelect4_RecvOK(b *testing.B) {
	chan0 := make(chan int)
	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int, 1)

	for n := 0; n < b.N; n++ {

		chan3 <- 1
		Select4(
			RecvOK(chan0, nothingOK),
			RecvOK(chan1, nothingOK),
			RecvOK(chan2, nothingOK),
			RecvOK(chan3, func(int, bool) {}),
		)
	}
}

func BenchmarkSelect4_Recv_preroll(b *testing.B) {
	chan0 := make(chan int)
	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int, 1)

	r0 := Recv(chan0, nothing)
	r1 := Recv(chan1, nothing)
	r2 := Recv(chan2, nothing)
	r3 := Recv(chan3, func(int) {})

	for n := 0; n < b.N; n++ {

		chan3 <- 1
		Select4(r0, r1, r2, r3)
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
			continue
		default:
		}
		select {
		case v0 := <-chan0:
			nothing(v0)
			continue
		case v1 := <-chan1:
			nothing(v1)
			continue
		default:
		}
		select {
		case v0 := <-chan0:
			nothing(v0)
			continue
		case v1 := <-chan1:
			nothing(v1)
			continue
		case v2 := <-chan2:
			nothing(v2)
			continue
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
