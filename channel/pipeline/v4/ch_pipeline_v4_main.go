package main

import "fmt"

func counter(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, v := range nums {
			out <- v
		}
	}()
	return out
}

func squarer(inCh <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for v := range inCh {
			out <- v * v
		}
	}()
	return out
}

func printer(inCh <-chan int) {
	for v := range inCh {
		fmt.Println(v)
	}
}

func main() {
	in := counter(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	outCh := squarer(in)

	printer(outCh)
}
