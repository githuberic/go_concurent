package e3_test

import (
	"go_concurent/example/bank/e3"
	"sync"
	"testing"
)

func TestVerify(t *testing.T)  {
	var wg sync.WaitGroup

	for i :=1; i <= 1000; i++ {
		wg.Add(1)
		go func(amount int) {
			e3.Deposit(amount)
			wg.Done()
		}(i)
	}
	wg.Wait()

	if got, want := e3.Balance(), (1000+1)*1000/2; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}
