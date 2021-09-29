package main

import (
	"fmt"
	"time"
)

type Task struct {
	f func() error
}

func NewTask(f func() error) *Task {
	return &Task{
		f: f,
	}
}

func (t *Task) Execute() {
	t.f()
}

type Pool struct {
	EntryChannel chan *Task
	workerNum    int
	JobChannel   chan *Task
}

func NewPool(cap int) *Pool {
	return &Pool{
		EntryChannel: make(chan *Task),
		workerNum:    cap,
		JobChannel:   make(chan *Task),
	}
}

func (pool *Pool) worker(workerId int) {
	for task := range pool.JobChannel {
		task.Execute()
		fmt.Println("Worker ID", workerId, ",done")
	}
}

func (pool *Pool) Run() {
	for i := 0; i < pool.workerNum; i++ {
		fmt.Println("Start worker:", i)
		go pool.worker(i)
	}

	for task := range pool.EntryChannel {
		pool.JobChannel <- task
	}

	close(pool.JobChannel)
	fmt.Println("Done,close jobs-channel")

	close(pool.EntryChannel)
	fmt.Println("Done,close entry-channel")
}

func main() {
	t := NewTask(func() error {
		fmt.Println("Create task,", time.Now().Format("2006-01-02 15:01:05"))
		return nil
	})

	pool := NewPool(3)

	go func() {
		for {
			pool.EntryChannel <- t
		}
	}()

	pool.Run()
}
