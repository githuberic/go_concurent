package loop

import (
	"fmt"
	"testing"
)

func TestVerify(t *testing.T) {
	// 我们遍历queue通道里面的两个数据
	var queue = make(chan string, 2)

	queue <- "one"
	queue <- "two"
	// ok
	close(queue)

	/*
		ok
		go func() {
			close(queue)
		}()
	*/

	/*
		defer func() {
			close(queue)
		}()
	*/

	//defer close(queue)

	// range函数遍历每个从通道接收到的数据，因为queue再发送完两个
	// 数据之后就关闭了通道，所以这里我们range函数在接收到两个数据
	// 之后就结束了。如果上面的queue通道不关闭，那么range函数就不
	// 会结束，从而在接收第三个数据的时候就阻塞了。
	for ele := range queue {
		fmt.Println(ele)
	}
}
