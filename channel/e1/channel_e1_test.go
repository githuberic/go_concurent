package e1

import (
	"log"
	"testing"
)

func producer(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, v := range nums {
			out <- v
		}
	}()
	return out
}

func consumer(inCh <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for n := range inCh {
			out <- n * n
		}
	}()
	return out
}

func TestVerify(t *testing.T) {
	in := producer(1, 2, 3, 4)
	retCh := consumer(in)

	for n := range retCh {
		log.Print("value=",n)
	}
}

// https://mp.weixin.qq.com/s/YB5XZ5NatniHSYBQ3AHONw