package main

import "fmt"

func main2() {

	// { // ordering is stable
	// 	ch1 := make(chan string, 1)
	// 	ch1 <- "hello world"

	// 	ch2 := make(chan string, 1)
	// 	ch2 <- "goodbye"

	// 	Select2Default(
	// 		ch1, func(s string) {
	// 			fmt.Println("1", s)
	// 		},
	// 		ch2, func(s string) {
	// 			fmt.Println("2", s)
	// 		},
	// 		func() {
	// 			fmt.Println("no winner")
	// 		},
	// 	)
	// }

	{ // ordering is not stable
		ch1 := make(chan string, 1)
		ch1 <- "hello world"

		ch2 := make(chan string, 1)
		ch2 <- "goodbye"

		var n, s string
		select {
		case s = <-ch1:
			n = "1"
		case s = <-ch2:
			n = "2"
		default:
			n = "d"
		}
		fmt.Println(n, s)
	}

}
