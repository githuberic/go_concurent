# 原子操作的基础知识

CPU 提供了基础的原子操作，不过，不同架构的系统的原子操作是不一样的。

### 单核-原子操作

对于单处理器单核系统来说，如果一个操作是由一个 CPU 指令来实现的，那么它就是原子操作，比如它的 XCHG 和 INC 等指令。

如果操作是基于多条指令来实现的，那么，执行的过程中可能会被中断，并执行上下文切换，这样的话，原子性的保证就被打破了，因为这个时候，操作可能只执行了一半。

### 多核-原子操作

在多处理器多核系统中，原子操作的实现就比较复杂了。

由于 cache 的存在，单个核上的单个指令进行原子操作的时候，你要确保其它处理器或者核不访问此原子操作的地址，或者是确保其它处理器或者核总是访问原子操作之后的最新的值。

x86 架构中提供了指令前缀 LOCK，LOCK 保证了指令（比如 LOCK CMPXCHG op1、op2）不会受其它处理器或 CPU 核的影响，有些指令（比如 XCHG）本身就提供 Lock 的机制。

不同的 CPU 架构提供的原子操作指令的方式也是不同的，比如对于多核的 MIPS 和 ARM，提供了 LL/SC（Load Link/Store Conditional）指令，可以帮助实现原子操作（ARMLL/SC 指令 LDREX 和 STREX）。

### golang-原子操作

因为不同的 CPU 架构甚至不同的版本提供的原子操作的指令是不同的，所以，要用一种编程语言实现支持不同架构的原子操作是相当有难度的。

不过，还好这些都不需要你操心，因为 Go 提供了一个通用的原子操作的 API，将更底层的不同的架构下的实现封装成 atomic 包，提供了修改类型的<b>原子操作（atomic read-modify-write，RMW）</b>

和加载存储类型的原子操作（Load 和 Store）的 API。 有的代码也会因为架构的不同而不同。有时看起来貌似一个操作是原子操作，但实际上，对于不同的架构来说，情况是不一样的。比如下面的代码的第 4 行，是将一个 64 位的值赋值给变量 i：

```go

const x int64 = 1 + 1<<33

func main() {
    var i = x
    _ = i
}
```

如果你使用 GOARCH=386 的架构去编译这段代码，那么，第 5 行其实是被拆成了两个指令，分别操作低 32 位和高 32 位（使用 GOARCH=386 go tool compile -N -l test.go；

GOARCH=386 go tool objdump -gnu test.o 反编译试试）：

![GOARCH=386](./img/atomoic_01.png)

如果 GOARCH=amd64 的架构去编译这段代码，那么，第 5 行其中的赋值操作其实是一条指令：

![GOARCH=amd64](./img/atomoic_02.png)

### atomic 提供的方法

目前的 Go 的泛型的特性还没有发布，Go 的标准库中的很多实现会显得非常啰嗦，多个类型会实现很多类似的方法，尤其是 atomic 包，最为明显。

相信泛型支持之后，atomic 的 API 会清爽很多。

<b>atomic 为了支持 int32、int64、uint32、uint64、uintptr、Pointer（Add 方法不支持）类型， 分别提供了 AddXXX、CompareAndSwapXXX、SwapXXX、LoadXXX、StoreXXX 等方法。</b>

不过，你也不要担心，你只要记住了一种数据类型的方法的意义，其它数据类型的方法也是一样的。

关于 atomic，还有一个地方你一定要记住，<b>atomic 操作的对象是一个地址，你需要把可寻址的变量的地址作为参数传递给方法，而不是把变量的值传递给方法。</b>

#### Add

首先，我们来看 Add 方法的签名：

![Add方法](./img/atomoic_03.png)

其实，Add 方法就是给第一个参数地址中的值增加一个 delta 值。

对于有符号的整数来说，delta 可以是一个负数，相当于减去一个值。对于无符号的整数和 uinptr 类型来说，怎么实现减去一个值呢？毕竟，atomic 并没有提供单独的减法操作。

#### CAS （CompareAndSwap）

以 int32 为例，我们学习一下 CAS 提供的功能。在 CAS 的方法签名中，需要提供要操作的地址、原数据值、新值，如下所示：

```go
func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
```

这个方法会比较当前 addr 地址里的值是不是 old，如果不等于 old，就返回 false；如果等于 old，就把此地址的值替换成 new 值，返回 true。这就相当于“判断相等才替换”。

如果使用伪代码来表示这个原子操作，代码如下：

```go
if *addr == old {
  *addr = new
  return true
}
return false
```

它支持的类型和方法如图所示：

![CompareAndSwap](./img/atomoic_04.png)

#### Swap

如果不需要比较旧值，只是比较粗暴地替换的话，就可以使用 Swap 方法，它替换后还可以返回旧值，伪代码如下：

```go
old = *addr
*addr = new
return old
```

#### Load

Load 方法会取出 addr 地址中的值，即使在多处理器、多核、有 CPU cache 的情况下，这个操作也能保证 Load 是一个原子操作。

它支持的数据类型和方法如图所示：

![Load方法](./img/atomoic_05.png)

#### Store

Store 方法会把一个值存入到指定的 addr 地址中，即使在多处理器、多核、有 CPU cache 的情况下，这个操作也能保证 Store 是一个原子操作。

别的 goroutine 通过 Load 读取出来，不会看到存取了一半的值。

它支持的数据类型和方法如图所示：

![Store](./img/atomoic_06.png)

#### Value 类型

刚刚说的都是一些比较常见的类型，其实，atomic 还提供了一个特殊的类型：Value。

它可以原子地存取对象类型，但也只能存取，不能 CAS 和 Swap，常常用在配置变更等场景中。

![Value类型](./img/atomoic_07.png)

接下来，我以一个配置变更的例子，来演示 Value 类型的使用。这里定义了一个 Value 类型的变量 config， 用来存储配置信息。

```go

type Config struct {
    NodeName string
    Addr     string
    Count    int32
}

func loadNewConfig() Config {
    return Config{
        NodeName: "北京",
        Addr:     "10.77.95.27",
        Count:    rand.Int31(),
    }
}
func main() {
    var config atomic.Value
    config.Store(loadNewConfig())
    var cond = sync.NewCond(&sync.Mutex{})

    // 设置新的config
    go func() {
        for {
            time.Sleep(time.Duration(5+rand.Int63n(5)) * time.Second)
            config.Store(loadNewConfig())
            cond.Broadcast() // 通知等待着配置已变更
        }
    }()

    go func() {
        for {
            cond.L.Lock()
            cond.Wait()                 // 等待变更信号
            c := config.Load().(Config) // 读取新的配置
            fmt.Printf("new config: %+v\n", c)
            cond.L.Unlock()
        }
    }()

    select {}
}
```

### 第三方库的扩展

不过有一点让人觉得不爽的是，或者是让熟悉面向对象编程的程序员不爽的是，函数调用有一点点麻烦。

所以，有些人就对这些函数做了进一步的包装，跟 atomic 中的 Value 类型类似，这些类型也提供了面向对象的使用方式，

比如关注度比较高的uber-go/atomic，它定义和封装了几种与常见类型相对应的原子操作类型，这些类型提供了原子操作的方法。

这些类型包括 Bool、Duration、Error、Float64、Int32、Int64、String、Uint32、Uint64 等。

比如 Bool 类型，提供了 CAS、Store、Swap、Toggle 等原子方法，还提供 String、MarshalJSON、UnmarshalJSON 等辅助方法，确实是一个精心设计的 atomic 扩展库。关于这些方法，你一看名字就能猜出来它们的功能，我就不多说了。

### 使用 atomic 实现 Lock-Free queue

atomic 常常用来实现 Lock-Free 的数据结构，这次我会给你展示一个 Lock-Free queue 的实现。

Lock-Free queue 最出名的就是 Maged M. Michael 和 Michael L. Scott 1996 年发表的论文中的算法，算法比较简单，容易实现

，伪代码的每一行都提供了注释，我就不在这里贴出伪代码了，因为我们使用 Go 实现这个数据结构的代码几乎和伪代码一样：

```go

package queue
import (
  "sync/atomic"
  "unsafe"
)
// lock-free的queue
type LKQueue struct {
  head unsafe.Pointer
  tail unsafe.Pointer
}
// 通过链表实现，这个数据结构代表链表中的节点
type node struct {
  value interface{}
  next  unsafe.Pointer
}
func NewLKQueue() *LKQueue {
  n := unsafe.Pointer(&node{})
  return &LKQueue{head: n, tail: n}
}
// 入队
func (q *LKQueue) Enqueue(v interface{}) {
  n := &node{value: v}
  for {
    tail := load(&q.tail)
    next := load(&tail.next)
    if tail == load(&q.tail) { // 尾还是尾
      if next == nil { // 还没有新数据入队
        if cas(&tail.next, next, n) { //增加到队尾
          cas(&q.tail, tail, n) //入队成功，移动尾巴指针
          return
        }
      } else { // 已有新数据加到队列后面，需要移动尾指针
        cas(&q.tail, tail, next)
      }
    }
  }
}
// 出队，没有元素则返回nil
func (q *LKQueue) Dequeue() interface{} {
  for {
    head := load(&q.head)
    tail := load(&q.tail)
    next := load(&head.next)
    if head == load(&q.head) { // head还是那个head
      if head == tail { // head和tail一样
        if next == nil { // 说明是空队列
          return nil
        }
        // 只是尾指针还没有调整，尝试调整它指向下一个
        cas(&q.tail, tail, next)
      } else {
        // 读取出队的数据
        v := next.value
                // 既然要出队了，头指针移动到下一个
        if cas(&q.head, head, next) {
          return v // Dequeue is done.  return
        }
      }
    }
  }
}

// 将unsafe.Pointer原子加载转换成node
func load(p *unsafe.Pointer) (n *node) {
  return (*node)(atomic.LoadPointer(p))
}

// 封装CAS,避免直接将*node转换成unsafe.Pointer
func cas(p *unsafe.Pointer, old, new *node) (ok bool) {
  return atomic.CompareAndSwapPointer(
    p, unsafe.Pointer(old), unsafe.Pointer(new))
}
```

这个 lock-free 的实现使用了一个辅助头指针（head），头指针不包含有意义的数据，只是一个辅助的节点，这样的话，出队入队中的节点会更简单。

入队的时候，通过 CAS 操作将一个元素添加到队尾，并且移动尾指针。

出队的时候移除一个节点，并通过 CAS 操作移动 head 指针，同时在必要的时候移动尾指针。

#### 最后，我还想和你讨论一个额外的问题：对一个地址的赋值是原子操作吗？

这是一个很有趣的问题，如果是原子操作，还要 atomic 包干什么？官方的文档中并没有特意的介绍，不过，在一些 issue 或者论坛中，

每当有人谈到这个问题时，总是会被建议用 atomic 包。

Dave Cheney就谈到过这个问题，讲得非常好。

我来给你总结一下他讲的知识点，这样你就比较容易理解使用 atomic 和直接内存操作的区别了。

在现在的系统中，write 的地址基本上都是对齐的（aligned）。 

比如，32 位的操作系统、CPU 以及编译器，write 的地址总是 4 的倍数，64 位的系统总是 8 的倍数（还记得 WaitGroup 针对 64 位系统和 32 位系统对 state1 的字段不同的处理吗）。

对齐地址的写，不会导致其他人看到只写了一半的数据，因为它通过一个指令就可以实现对地址的操作。

如果地址不是对齐的话，那么，处理器就需要分成两个指令去处理，如果执行了一个指令，其它人就会看到更新了一半的错误的数据，这被称做撕裂写（torn write） 。

所以，你可以认为赋值操作是一个原子操作，这个“原子操作”可以认为是保证数据的完整性。

###### 多处理器

但是，对于现代的多处理多核的系统来说，由于 cache、指令重排，可见性等问题，我们对原子操作的意义有了更多的追求。

在多核系统中，一个核对地址的值的更改，在更新到主内存中之前，是在多级缓存中存放的。

这时，多个核看到的数据可能是不一样的，其它的核可能还没有看到更新的数据，还在使用旧的数据。

多处理器多核心系统为了处理这类问题，使用了一种叫做内存屏障（memory fence 或 memory barrier）的方式。

一个写内存屏障会告诉处理器，必须要等到它管道中的未完成的操作（特别是写操作）都被刷新到内存中，再进行操作。

此操作还会让相关的处理器的 CPU 缓存失效，以便让它们从主存中拉取最新的值。

atomic 包提供的方法会提供内存屏障的功能，所以，atomic 不仅仅可以保证赋值的数据完整性，还能保证数据的可见性，

一旦一个核更新了该地址的值，其它处理器总是能读取到它的最新值。

但是，需要注意的是，因为需要处理器之间保证数据的一致性，atomic 的操作也是会降低性能的。

![总结](./img/atomoic_08.jpg)
