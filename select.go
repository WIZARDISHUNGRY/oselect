package main

func Select2[T1, T2 any](
	c1 <-chan T1, f1 func(T1),
	c2 <-chan T2, f2 func(T2),
) {
	select {
	case v1 := <-c1:
		f1(v1)
		return
	default:
	}

	select {
	case v1 := <-c1:
		f1(v1)
	case v2 := <-c2:
		f2(v2)
	}
}

func Select2Default[T1, T2 any](
	c1 <-chan T1, f1 func(T1),
	c2 <-chan T2, f2 func(T2),
	df func(),
) {
	select {
	case v1 := <-c1:
		f1(v1)
		return
	default:
	}

	select {
	case v1 := <-c1:
		f1(v1)
	case v2 := <-c2:
		f2(v2)
	default:
		df()
	}
}

func Select3[T1, T2, T3 any](
	c1 <-chan T1, f1 func(T1),
	c2 <-chan T2, f2 func(T2),
	c3 <-chan T3, f3 func(T3),
) {
	select {
	case v1 := <-c1:
		f1(v1)
		return
	default:
	}

	select {
	case v1 := <-c1:
		f1(v1)
		return
	case v2 := <-c2:
		f2(v2)
		return
	default:
	}

	select {
	case v1 := <-c1:
		f1(v1)
	case v2 := <-c2:
		f2(v2)
	case v3 := <-c3:
		f3(v3)
	}

}

func Select3Default[T1, T2, T3 any](
	c1 <-chan T1, f1 func(T1),
	c2 <-chan T2, f2 func(T2),
	c3 <-chan T3, f3 func(T3),
	df func(),
) {
	select {
	case v1 := <-c1:
		f1(v1)
		return
	default:
	}

	select {
	case v1 := <-c1:
		f1(v1)
		return
	case v2 := <-c2:
		f2(v2)
		return
	default:
	}

	select {
	case v1 := <-c1:
		f1(v1)
	case v2 := <-c2:
		f2(v2)
	case v3 := <-c3:
		f3(v3)
	default:
		df()
	}

}
