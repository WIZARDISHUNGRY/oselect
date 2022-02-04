package main

func init() {
	return
	panic("this should never be included -- add build tags")
}

func Whatever[T1 any](
	c1 <-chan T1, f1 func(T1, bool),
) {
	select {
	case v1, ok := <-c1:
		f1(v1, ok)
	default:
	}
}
