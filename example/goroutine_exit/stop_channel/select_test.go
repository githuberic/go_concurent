package stop_channel

import (
	"fmt"
	"testing"
	"time"
)

func worker(inCh <-chan struct{}) {
	go func() {
		defer fmt.Println("worker exit")

		t := time.NewTicker(time.Millisecond * 500)

		// Using stop channel explicit exit
		for {
			select {
			case <-inCh:
				fmt.Println("Recv stop signal")
				return
			case <-t.C:
				fmt.Println("Working .")
			}
		}
	}()
	return
}

func TestVerify(t *testing.T) {
	stopCh := make(chan struct{})
	worker(stopCh)

	time.Sleep(time.Second * 2)
	close(stopCh)

	time.Sleep(time.Second)
	fmt.Println("main exit")
}
