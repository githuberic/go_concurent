package timeout

import (
	"fmt"
	"testing"
	"time"
)

func TestVerify(t *testing.T) {
	// 在这个例子中，假设我们执行了一个外部调用，2秒之后将结果写入c1
	var c = make(chan string, 1)
	go func() {
		time.Sleep(time.Second * 2)
		c <- "result 1"
	}()

	// 这里使用select来实现超时，`res := <-c1`等待通道结果，
	// `<- Time.After`则在等待1秒后返回一个值，因为select首先
	// 执行那些不再阻塞的case，所以这里会执行超时程序，如果
	// `res := <-c1`超过1秒没有执行的话
	select {
	case res := <-c:
		fmt.Println(res)
	case <-time.After(time.Second * 1):
		fmt.Println("timeout 1")
	}

	// 如果我们将超时时间设为3秒，这个时候`res := <-c2`将在
	// 超时case之前执行，从而能够输出写入通道c2的值
	var c2 = make(chan string, 1)
	go func() {
		time.Sleep(time.Second * 2)
		c2 <- "result 2"
	}()

	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(time.Second * 3):
		fmt.Println("timeout 2")
	}
}
