# 对map并行读写
如果map由多协程同时读和写就会出现 fatal error:concurrent map read and map write的错误

### 原因
因为map为引用类型，所以即使函数传值调用，参数副本依然指向映射m, 所以多个goroutine并发写同一个映射m， 
写过多线程程序的同学都知道，对于共享变量，资源，并发读写会产生竞争的， 故共享资源遭到破坏

### 直接读写
会出现 fatal error: concurrent map writes

```go
package v1

import "testing"

var mMap map[int]int

func TestVerify(t *testing.T) {
	mMap = make(map[int]int)

	for i := 0; i < 1000; i++ {
		go func() {
			mMap[i] = i
		}()

		go readMap(i)
	}
}

func readMap(index int) int {
	return mMap[index]
}
```

### 解决方案一 加锁控制(mutex)

```go
package v2

import (
	"sync"
	"testing"
)

type ConcurrentMap struct {
	sync.Mutex
	Map map[int]int
}

func (m *ConcurrentMap) readMap(index int) (int, bool) {
	m.Lock()
	defer m.Unlock()
	value, ok := m.Map[index]
	return value, ok
}

func (m *ConcurrentMap) writeMap(index int, value int) {
	m.Lock()
	defer m.Unlock()
	m.Map[index] = value
}

func TestVerify(t *testing.T) {
	var mMap = &ConcurrentMap{
		Map: make(map[int]int),
	}

	for i := 0; i < 1000; i++ {
		go func() {
			mMap.writeMap(i, i)
		}()

		go mMap.readMap(i)
	}
}
```

### 解决方案一 加锁控制(rwmutex)
```go

package _map

import (
	"log"
	"sync"
	"testing"
)

type RWMutexMap struct {
	sync.RWMutex
	Map map[int]int
}

func (m *RWMutexMap) Get(index int) (int, bool) {
	m.RLock()
	defer m.RUnlock()
	value, ok := m.Map[index]
	return value, ok
}

func (m *RWMutexMap) Set(index int, value int) {
	m.Lock()
	defer m.Unlock()
	m.Map[index] = value
}

func TestVerifyRWMutex(t *testing.T) {
	var mMap = &RWMutexMap{
		Map: make(map[int]int),
	}

	for i := 1; i < 1000; i++ {
		go func() {
			mMap.Set(i, i)
		}()

		go func() {
			value, ok := mMap.Get(i)
			if ok {
				log.Print(value)
			}
		}()
	}
}
```

### 解决方案二 利用channel串行化处理
在go语言里，提倡用信道通讯的方式来替代显式的同步机制。

```go
package _map

import (
	"log"
	"testing"
)

type ChannelMap struct {
	Map map[int]int
	ch  chan func()
}

func NewChannelMap() *ChannelMap {
	m := &ChannelMap{
		Map: make(map[int]int),
		ch:  make(chan func()),
	}

	go func() {
		for {
			(<-m.ch)()
		}
	}()

	return m
}

func (m *ChannelMap) add(index int, value int) {
	m.ch <- func() {
		m.Map[index] = value
	}
}

func (m *ChannelMap) del(index int) {
	m.ch <- func() {
		delete(m.Map, index)
	}
}

func (m *ChannelMap) find(index int) (data int) {
	// 每次查询都要创建一个信道
	m.ch <- func() {
		if res, ok := m.Map[index]; ok {
			data = res
		}
	}
	return
}

func TestVerifyChMap(t *testing.T) {
	mMap := NewChannelMap()

	for i := 0; i < 1000; i++ {
		go func() {
			mMap.add(i, i)
		}()

		go func() {
			value := mMap.find(i)
			if value > 0 {
				log.Printf("Value=%d",value)
			}
		}()
	}
}
```

但是发现有的时候用信道通讯方式实现的似乎也不是很好（暂不考虑效率问题）。

每次调用find都要创建一个通道;

### 解决方案二 利用channel串行化处理(通道预分配)

<pre>利用预分配以及可回收的channel来提高资源利用率。
这个技术在多个goroutine等待一个主动对象返回自己的数据时会比较有用。
例如网游服务器中登录服务器里每个玩家的连接用一个goroutine来处理&#xff1b;
另外一个主动对象代表帐号服务器连接用于验证帐号合法性。
玩家goroutine会把各自的输入的玩家帐号密码发送给这个主动对象&#xff0c;
并阻塞等待主动对象返回验证结果。因为有多个玩家同时发起帐号验证请求&#xff0c;
所以主动对象需要把返回结果进行分发&#xff0c;因此可以在发送请求的时候申请一个信道并等待这个信道。
</pre>

```go
package _map

import (
	"log"
	"testing"
)

type ChannelPoolMap struct {
	Map    map[int]int
	ch     chan func()
	tokens chan chan *int
}

func NewConcurrentMap() *ChannelPoolMap {
	m := &ChannelPoolMap{
		Map:    make(map[int]int),
		ch:     make(chan func()),
		tokens: make(chan chan *int, 128),
	}

	for i := 0; i < cap(m.tokens); i++ {
		m.tokens <- make(chan *int)
	}

	go func() {
		for {
			(<-m.ch)()
		}
	}()

	return m
}

func (m *ChannelPoolMap) add(index int, value int) {
	m.ch <- func() {
		m.Map[index] = value
	}
}

func (m *ChannelPoolMap) del(index int) {
	m.ch <- func() {
		delete(m.Map, index)
	}
}

func (m *ChannelPoolMap) find(index int) *int {
	// 每次查询都要创建一个信道
	c := <-m.tokens

	m.ch <- func() {
		res, ok := m.Map[index]
		if !ok {
			c <- nil
			//data = res
		} else {
			inf := res
			c <- &inf
		}
	}
	inf := <-c

	// 回收信道
	m.tokens <- c

	return inf
}

func TestVerify(t *testing.T) {
	mMap := NewConcurrentMap()

	for i := 0; i < 1000; i++ {
		go func() {
			mMap.add(i, i)
		}()

		go func() {
			value := mMap.find(i)
			if value != nil && *value > 0 {
				log.Print(*value)
			}
		}()
	}
}
```

### 解决方案三 利用sync.Map
```go

```


<blockquote>
mutex/rwmutex的区别，
性能比较差异（benchmark）
</blockquote>

