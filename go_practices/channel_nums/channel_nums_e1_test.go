package channel_nums

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestVerify(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 3)

	for i := 0; i < 10; i++ {
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			log.Println(i)
			time.Sleep(time.Second * 1)
			<-ch
		}(i)
	}
	wg.Wait()
}
