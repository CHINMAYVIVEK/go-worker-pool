package work

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Executer interface defines methods for executing a task and handling errors
type Executer interface {
	Execute() error
	OnError(error)
}

// Pool represents a worker pool
type Pool struct {
	numWorkers    int           // Number of workers in the pool
	tasks         chan Executer // Channel for tasks
	start         sync.Once     // Ensures the pool starts only once
	stop          sync.Once     // Ensures the pool stops only once
	taskCompleted chan bool     // Channel for task completion notifications
	quit          chan struct{} // Channel for quitting
}

// NewPool initializes a new worker pool
func NewPool(numWorkers, taskChanSize int) (*Pool, error) {
	if numWorkers <= 0 {
		return nil, errors.New("number of workers must be greater than 0")
	}
	if taskChanSize < 0 {
		return nil, errors.New("task channel size must be a positive number")
	}
	return &Pool{
		numWorkers:    numWorkers,
		tasks:         make(chan Executer, taskChanSize),
		taskCompleted: make(chan bool),
		quit:          make(chan struct{}),
	}, nil
}

// Start begins the worker pool
func (p *Pool) Start(ctx context.Context) {
	p.start.Do(func() {
		p.startWorker(ctx)
	})
}

// Stop shuts down the worker pool
func (p *Pool) Stop() {
	p.stop.Do(func() {
		close(p.quit) // Notify all workers to stop
	})
}

// AddTask adds a task to the pool (blocking)
func (p *Pool) AddTask(t Executer) {
	select {
	case p.tasks <- t:
	case <-p.quit:
	}
}

// AddTaskNonBlocking adds a task to the pool without blocking
func (p *Pool) AddTaskNonBlocking(t Executer) {
	go func() {
		select {
		case p.tasks <- t:
		case <-p.quit:
		}
	}()
}

// TaskCompleted returns a channel to notify when tasks are completed
func (p *Pool) TaskCompleted() <-chan bool {
	return p.taskCompleted
}

// startWorker initializes worker goroutines
func (p *Pool) startWorker(ctx context.Context) {
	for i := 0; i < p.numWorkers; i++ {
		go func(workerNum int) {
			fmt.Println("Worker", workerNum, "started")
			for {
				select {
				case <-ctx.Done(): // Stop if context is done
					return
				case <-p.quit: // Stop if quit signal received
					return
				case task, ok := <-p.tasks: // Retrieve task from the queue
					if !ok {
						return // Exit if tasks channel is closed
					}
					if err := task.Execute(); err != nil {
						task.OnError(err) // Handle task error
					}
					p.taskCompleted <- true // Notify task completion
					fmt.Println("Worker", workerNum, "finished task")
				}
			}
		}(i)
	}
}
