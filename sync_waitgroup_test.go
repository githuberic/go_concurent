package go_concurrent

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestVerifyV3(t *testing.T) {
	var waitGroup sync.WaitGroup

	for i := 0; i < 10; i++ {
		waitGroup.Add(1) // 添加需要等待goroutine的数量
		go func() {
			fmt.Println("WaitGroup")
			time.Sleep(time.Second)
			waitGroup.Done() // 减少需要等待goroutine的数量 相当于Add(-1)
		}()
	}
	waitGroup.Wait()
	// 执行阻塞，直到所有的需要等待的goroutine数量变成0
	fmt.Println("over")
}

func TestVerifyV6(t *testing.T) {
	ch := make(chan int, 5)

	for i := 0; i < 5; i++ {
		go func(v int) {
			ch <- v
		}(i)
	}

	go func() {
		for v := range ch {
			fmt.Println("Value=", v)
		}
	}()

	go func() {
		fmt.Println("over")
	}()
}
