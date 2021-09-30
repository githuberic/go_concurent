package mutex

import (
	"fmt"
	"sync"
	"testing"
)

var num = 0
func TestVerifyMutex(t *testing.T) {
	mu := sync.Mutex{}
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			mu.Lock()
			defer mu.Unlock()
			num += 1
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("num=", num)
}
