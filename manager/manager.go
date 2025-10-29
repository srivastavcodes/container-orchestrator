package manager

import (
	"fmt"
	"orchestrator/task"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

// Manager accepts request from users, schedules tasks onto worker
// machines, keeps track of tasks, their states, and the machines
// on which they run.
type Manager struct {
	// Pending is the queue holding the tasks to assign to worker in FIFO
	// order.
	Pending queue.Queue

	// TaskDb is an in-memory database holding Tasks.
	TaskDb map[string][]*task.Task

	// EventDb is an in-memory database holding TaskEvents.
	EventDb map[string][]*task.TaskEvent

	// Workers the Manager keeps track of in the cluster.
	Workers []string

	// WorkerTaskMap is a map of jobs assigned to worker.
	WorkerTaskMap map[string][]uuid.UUID

	// TaskWorkerMap is a map of worker running a given Task.
	TaskWorkerMap map[uuid.UUID]string
}

// SelectWorker schedules tasks onto workers. This method is responsible for
// looking at the requirements specified in a Task and evaluating the
// resources available in the pool of workers to find the best suited worker
// to run the Task.
func (m *Manager) SelectWorker() {
	fmt.Println("I will select an appropriate worker.")
}

// UpdateTasks keeps track of tasks, their states and the machine on which they
// run. This method calls the worker.CollectStats for stats.
func (m *Manager) UpdateTasks() {
	fmt.Println("I will update tasks.")
}

// SendWork sends tasks to workers.
func (m *Manager) SendWork() {
	fmt.Println("I will send work to worker.")
}
