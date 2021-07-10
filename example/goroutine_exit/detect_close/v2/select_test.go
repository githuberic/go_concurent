package v2

import (
	"fmt"
	"testing"
	"time"
)

func producer(n int) <-chan int {
	retCh := make(chan int)

	go func() {
		defer func() {
			close(retCh)
			retCh = nil
			fmt.Println("Producer exit")
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

	t := time.NewTicker(time.Millisecond * 500)
	processedCnt := 0

	go func() {
		defer func() {
			fmt.Println("Consumer-worker exit")
			retCh <- struct{}{}
			close(retCh)
		}()

		for {
			select {
			case x, ok := <-intCh:
				if !ok {
					return
				}
				fmt.Printf("Consumer-Process %d\n", x)
			case <-t.C:
				fmt.Printf("Working, processedCnt = %d\n", processedCnt)
			}
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
