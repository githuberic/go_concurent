package e1

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func worker(id int, jobCh <-chan int, retCh chan<- string) {
	cnt := 0

	for job := range jobCh {
		cnt++
		ret := fmt.Sprintf("worker %d processed job: %d, it's the %dth processed by me.", id, job, cnt)
		retCh <- ret
	}
}

func workPool(n int, jobCh <-chan int, retCh chan<- string) {
	for i := 0; i < n; i++ {
		go worker(i, jobCh, retCh)
	}
}

func genJob(n int) <-chan int {
	jobCh := make(chan int, 200)

	go func() {
		for i := 0; i < n; i++ {
			jobCh <- i
		}
		close(jobCh)
	}()
	return jobCh
}

func TestVerify(t *testing.T) {
	jobCh := genJob(10000)
	retCh := make(chan string,10000)

	workPool(5,jobCh,retCh)

	time.Sleep(time.Second * 3)

	close(retCh)

	for ret := range retCh {
		log.Println(ret)
	}
}
// https://lessisbetter.site/2018/12/20/golang-simple-goroutine-pool/
