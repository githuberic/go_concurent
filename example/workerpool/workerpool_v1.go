package main

import (
	"fmt"
	"time"
)

func worker(index int, jobCh <-chan int, retCh chan<- string) {
	for job := range jobCh {
		ret := fmt.Sprintf("worker %d processed job: %d", index, job)
		retCh <- ret
	}
}

func workerPool(count int, jobCh <-chan int, retCh chan<- string) {
	for i := 0; i < count; i++ {
		go worker(i, jobCh, retCh)
	}
}

func geneJob(count int) <-chan int {
	jobCh := make(chan int, 200)

	go func() {
		for i := 0; i < count; i++ {
			jobCh <- i
		}
		close(jobCh)
	}()

	return jobCh
}

func main() {
	jobCh := geneJob(20)
	retCh := make(chan string, 100)
	workerPool(5, jobCh, retCh)

	time.Sleep(1 * time.Second)
	close(retCh)
	for ret := range retCh {
		fmt.Println(ret)
	}
}
