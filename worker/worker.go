package worker

import (
	"errors"
	"fmt"
	"log"
	"orchestrator/task"
	"time"

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
func (w *Worker) RunTask() task.DockerResult {
	currTask := w.Queue.Dequeue()
	if currTask == nil {
		log.Println("No task in the queue")
		return task.DockerResult{Error: nil}
	}
	queuedTask := currTask.(task.Task)

	persistedTask := w.Db[queuedTask.ID]
	if persistedTask == nil {
		persistedTask = &queuedTask
		w.Db[queuedTask.ID] = &queuedTask
	}
	var result task.DockerResult

	if task.ValidStateTransition(persistedTask.State, queuedTask.State) {
		switch queuedTask.State {
		case task.Scheduled:
			result = w.StartTask(queuedTask)
		case task.Completed:
			result = w.StopTask(queuedTask)
		default:
			result.Error = errors.New("docker? I don't even know her")
		}
	} else {
		err := fmt.Errorf(
			"invalid state transition from %v to %v",
			persistedTask.State, queuedTask.State,
		)
		result.Error = err
	}
	return result
}

// AddTask enqueues a Task in the worker's queue.
func (w *Worker) AddTask(currTask task.Task) {
	w.Queue.Enqueue(currTask)
}

// StartTask starts the Task.
func (w *Worker) StartTask(currTask task.Task) task.DockerResult {
	currTask.StartTime = time.Now().UTC()
	var (
		config = task.NewConfig(&currTask)
		docker = task.NewDocker(config)
		result = docker.Run()
	)
	if result.Error != nil {
		log.Printf("Error running task %v: %v", currTask.ID, result.Error)
		currTask.State = task.Failed
		w.Db[currTask.ID] = &currTask
		return result
	}
	currTask.ContainerID = result.ContainerID
	currTask.State = task.Running
	w.Db[currTask.ID] = &currTask
	log.Printf("Container %v running for task %v\n", currTask.ContainerID, currTask.ID)
	return result
}

// StopTask starts the Task.
func (w *Worker) StopTask(currTask task.Task) task.DockerResult {
	config := task.NewConfig(&currTask)
	docker := task.NewDocker(config)

	result := docker.Stop(currTask.ContainerID)
	if result.Error != nil {
		log.Printf("Error stopping container %v: %v\n", currTask.ContainerID, result.Error)
	}
	currTask.FinishTime = time.Now()
	currTask.State = task.Completed
	w.Db[currTask.ID] = &currTask
	log.Printf("Stopped and removed container %v for task %v\n", currTask.ContainerID, currTask.ID)
	return result
}
