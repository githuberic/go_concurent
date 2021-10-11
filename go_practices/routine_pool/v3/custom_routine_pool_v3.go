package main

import (
	"fmt"
	"time"
)

// 定义任务Task类型,每一个任务Task都可以抽象成一个函数
type Task struct {
	f func() error
}

//通过NewTask来创建一个Task
func NewTask(f func() error) *Task {
	return &Task{
		f: f,
	}
}

// 执行Task任务的方法
func (t *Task) Execute() {
	t.f() //调用任务所绑定的函数
}

/* 有关协程池的定义及操作 */
type Pool struct {
	//对外接收Task的入口
	EntryChannel chan *Task
	//协程池最大worker数量,限定Goroutine的个数
	workerNum int
	//协程池内部的任务就绪队列
	JobChannel chan *Task
}

//创建一个协程池
func NewPool(cap int) *Pool {
	return &Pool{
		EntryChannel: make(chan *Task),
		workerNum:    cap,
		JobChannel:   make(chan *Task),
	}
}

//协程池创建一个worker并且开始工作
func (pool *Pool) worker(workerId int) {
	//worker不断的从JobsChannel内部任务队列中拿任务
	for task := range pool.JobChannel {
		//如果拿到任务,则执行task任务
		task.Execute()
		fmt.Println("Worker ID", workerId, ",done")
	}
}

// 让协程池Pool开始工作
func (pool *Pool) Run() {
	//1,首先根据协程池的worker数量限定,开启固定数量的Worker,
	for i := 0; i < pool.workerNum; i++ {
		fmt.Println("Start worker:", i)
		// 每一个Worker用一个Goroutine承载
		go pool.worker(i)
	}

	//2, 从EntryChannel协程池入口取外界传递过来的任务
	// 并且将任务送进JobsChannel中
	for task := range pool.EntryChannel {
		pool.JobChannel <- task
	}

	// 3, 执行完毕需要关闭JobsChannel
	close(pool.JobChannel)

	//4, 执行完毕需要关闭EntryChannel'
	close(pool.EntryChannel)
}

func main() {
	t := NewTask(func() error {
		fmt.Println("Create task,", time.Now().Format("2006-01-02 15:01:05"))
		return nil
	})

	pool := NewPool(3)

	//开一个协程 不断的向 Pool 输送打印一条时间的task任务
	go func() {
		for {
			pool.EntryChannel <- t
		}
	}()

	// 启动协程池p
	pool.Run()
}
