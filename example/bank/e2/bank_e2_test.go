package e2_test

import (
	"go_concurent/example/bank/e2"
	"sync"
	"testing"
)

func TestVerify(t *testing.T) {
	var wg sync.WaitGroup

	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func(amount int) {
			e2.Deposit(amount)
			wg.Done()
		}(i)
	}
	wg.Wait()

	if got, want := e2.Balance(), 1000*(1000+1)/2; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}
