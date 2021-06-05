package e1

import (
	"fmt"
	"sync"
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

func sequare(inCh <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for v := range inCh {
			out <- v * v
		}
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup
	wg.Add(len(cs))

	collect := func(in <-chan int) {
		defer wg.Done()
		for n := range in {
			out <- n
		}
	}

	for _, v := range cs {
		go collect(v)
	}

	/*
	错误方式：直接等待是bug，死锁，因为merge写了out，main却没有读
	wg.Wait()
	close(out)
	 */

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func TestVerify(t *testing.T) {
	in := producer(1, 2, 3, 4)

	c1 := sequare(in)
	c2 := sequare(in)
	c3 := sequare(in)

	for ret := range merge(c1, c2, c3) {
		fmt.Printf("%3d",ret)
	}
	fmt.Println()
}

// https://lessisbetter.site/2018/11/28/golang-pipeline-fan-model/