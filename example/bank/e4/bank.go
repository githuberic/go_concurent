package e4

import "sync"

var (
	lock    sync.RWMutex
	balance int
)

func Deposit(amount int) {
	lock.Lock()
	defer lock.Unlock()
	balance += amount
}

func Balance() int {
	lock.RLock()
	defer lock.RUnlock()
	return balance
}
