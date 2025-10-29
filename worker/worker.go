package worker

import (
	"fmt"
	"orchestrator/task"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

// Worker is responsible for running docker containers as Task.
type Worker struct {
	// Name is the identifier used to identify an individual worker.
	Name string

	// Queue holds the tasks from the Manager. Tasks are handled in
	// FIFO order
	Queue queue.Queue

	// Db is map of Task uuids to Task to keep track of tasks.
	Db map[uuid.UUID]*task.Task

	// TaskCount is the total number of tasks a worker has at any given
	// time.
	TaskCount int
}

// CollectStats can be used to periodically collects stats about the
// worker.
func (w *Worker) CollectStats() {
	fmt.Println("i will collect stats")
}

// RunTask is responsible to identifying the Task's current state and
// then either starting or stopping the Task depending on the state.
func (w *Worker) RunTask() {
	fmt.Println("i will run tasks")
}

// StartTask starts the Task.
func (w *Worker) StartTask() {
	fmt.Println("i will start tasks")
}

// StopTask starts the Task.
func (w *Worker) StopTask() {
	fmt.Println("i will run tasks")
}
