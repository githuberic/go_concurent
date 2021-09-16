package pipeline

import (
	"fmt"
	"testing"
)

func counterV4(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, v := range nums {
			out <- v
		}
	}()
	return out
}

func squarerV4(inCh <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for v := range inCh {
			out <- v * v
		}
	}()
	return out
}

func printerV4(inCh <-chan int) {
	for v := range inCh {
		fmt.Println(v)
	}
}

func TestVerifyV4(t *testing.T) {
	in := counterV4(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	outCh := squarerV4(in)

	printerV4(outCh)
}
