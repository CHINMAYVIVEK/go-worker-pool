# Go Worker Pool
## Fire and forget worker pool 
This project implements a simple worker pool in Go, allowing for concurrent task execution with a specified number of workers.

## Features

- **Worker Pool**: Create a pool of workers to process tasks concurrently.
- **Task Management**: Add tasks to the pool, which can execute functions and handle errors.
- **Graceful Shutdown**: Workers can be stopped gracefully when all tasks are completed.

## Getting Started

### Prerequisites

- Go 1.23 or later

### Installation

Clone the repository:

```bash
git clone https://github.com/chinmayvivek/go-worker-pool.git
cd go-worker-pool
