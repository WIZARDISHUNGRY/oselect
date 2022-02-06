package oselect

// Param is an argument for the Select* functions.
type Param[T any] struct {
	RecvChan  <-chan T
	RecvFunc  func(T, bool)
	SendChan  chan<- T
	SendValue T
}

// Send is a helper for generating a Param for sending on a channel
func Send[T any](c chan<- T, v T) Param[T] {
	return Param[T]{
		SendChan:  c,
		SendValue: v,
	}
}

// RecvOK is a helper for registering a callback for receiving on a channel
func RecvOK[T any](c <-chan T, f func(T, bool)) Param[T] {
	return Param[T]{
		RecvChan: c,
		RecvFunc: f,
	}
}

// Recv is a helper for registering a callback for receiving on a channel
func Recv[T any](c <-chan T, f func(T)) Param[T] {
	return RecvOK(c, func(t T, b bool) { f(t) })
}
