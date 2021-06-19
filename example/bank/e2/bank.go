package e2

var (
	// 用来保护balance的二进制信号量
	sema    = make(chan struct{}, 1)
	balance int
)

func Deposit(amount int) {
	// acquire token
	sema <- struct{}{}
	balance += amount
	// release token
	<-sema
}

func Balance() int {
	sema <- struct{}{}
	b := balance
	<-sema
	return b
}
