package task

import (
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

// Task represents a task that the user wants to run on our cluster.
type Task struct {
	// ID is the 128 bit universally unique identifier to distinguish between
	// individual tasks.
	ID          uuid.UUID
	ContainerID string

	// Name is the name of the task for better readability.
	Name string

	// State represents the current state of the task. See State more info.
	State State

	// Image the docker image a task will use.
	Image string
	CPU   float64

	// Memory and Disk help the system identify the number of resources a task
	// needs.
	Memory int64
	Disk   int64 // See Memory for docs.

	// ExposedPorts and PortBindings are used by docker to ensure the machine
	// allocates the proper network ports for the task.
	ExposedPorts nat.PortSet
	PortBindings map[string]string // See ExposedPorts for docs.

	// RestartPolicy tells the system what to do when a task stops or fails
	// unexpectedly.
	RestartPolicy string

	// StartTime and FinishTime lets the user know when the task started and
	// when it stopped.
	StartTime  time.Time
	FinishTime time.Time // see StartTime for docs.
}

// TaskEvent represents a change in the current state of the Task.
type TaskEvent struct {
	ID        uuid.UUID
	State     State
	TimeStamp time.Time
	Task      Task
}
