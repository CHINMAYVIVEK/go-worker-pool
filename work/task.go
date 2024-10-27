package work

// Task represents a unit of work to be executed
type Task struct {
	execute      func() error // Function to execute the task
	errorHandler func(error)  // Function to handle errors
}

// NewTask creates a new task with provided execute and error handling functions
func NewTask(execute func() error, errorHandler func(error)) *Task {
	return &Task{
		execute:      execute,
		errorHandler: errorHandler,
	}
}

// Execute runs the task
func (t *Task) Execute() error {
	return t.execute()
}

// OnError handles errors that occur during task execution
func (t *Task) OnError(err error) {
	t.errorHandler(err)
}
