package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/CHINMAYVIVEK/go-worker-pool/work"
)

func main() {
	// Create a new worker pool with 5 workers and a task channel size of 5
	wp, err := work.NewPool(5, 5)
	if err != nil {
		panic(err)
	}

	// Create a context for canceling the worker pool
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure resources are released when done

	// Start the worker pool
	wp.Start(ctx)

	// Add 20 tasks to the worker pool
	for i := 0; i < 20; i++ {
		task := work.NewTask(func() error {
			const urlString = "https://www.google.com"
			fmt.Println("Fetching", urlString)
			res, err := http.Get(urlString)
			if err != nil {
				return err
			}
			fmt.Printf("Fetched url %s, status code: %d\n", urlString, res.StatusCode)
			return nil
		}, func(err error) {
			fmt.Println("Error fetching:", err)
		})
		wp.AddTaskNonBlocking(task) // Add task without blocking
	}

	// Wait for all tasks to complete
	counter := 0
	for completed := range wp.TaskCompleted() {
		if completed {
			counter++
		}
		if counter == 20 {
			wp.Stop() // Stop the worker pool after all tasks are done
			return
		}
	}
}
