package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

func sendMain(i int, wg *sync.WaitGroup) func() {
	var cnt int

	return func() {
		for {
			time.Sleep(time.Second * 2)
			fmt.Println("send email to ", i)
			cnt++

			if cnt > 5 && i == 1 {
				fmt.Println("Exit go-routine id:", i)
				break
			}
		}
		wg.Done()
	}
}

func main() {
	wg := sync.WaitGroup{}

	pool, _ := ants.NewPool(2)

	defer pool.Release()

	for i := 1; i <= 5; i++ {
		pool.Submit(sendMain(i, &wg))
		wg.Add(1)
	}
	wg.Wait()
}
