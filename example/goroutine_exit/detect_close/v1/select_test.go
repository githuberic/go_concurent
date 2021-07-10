package v1

import (
	"fmt"
	"testing"
)

func producer(n int) <-chan int {
	retCh := make(chan int)

	go func() {
		defer func() {
			close(retCh)
			retCh = nil
			fmt.Println("producer exit")
		}()

		for i := 0; i < n; i++ {
			fmt.Printf("send %d\n", i)
			retCh <- i
		}
	}()
	return retCh
}

func consumer(intCh <-chan int) <-chan struct{} {
	retCh := make(chan struct{})

	go func() {
		defer func() {
			fmt.Println("worker exit")
			retCh <- struct{}{}
			close(retCh)
		}()

		// Using for-range to exit goroutine
		// range has the ability to detect the close/end of a channel
		for ch := range intCh {
			fmt.Printf("Process %d\n", ch)
		}
	}()

	return retCh
}

func TestVerify(t *testing.T) {
	retCh := producer(3)

	// Wait consumer exit
	endCh := consumer(retCh)
	<-endCh

	fmt.Println("main exit")
}
