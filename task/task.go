package task

import (
	"context"
	"io"
	"log"
	"math"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
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
	Timestamp time.Time
	Task      Task
}

// Config struct to hold Docker container config
type Config struct {
	// Name of the task, also used as the container name
	ContainerName string

	// AttachStdin boolean which determines if stdin should be attached
	AttachStdin bool

	// AttachStdout boolean which determines if stdout should be attached
	AttachStdout bool

	// AttachStderr boolean which determines if stderr should be attached
	AttachStderr bool

	// ExposedPorts list of ports exposed
	ExposedPorts nat.PortSet

	// Cmd to be run inside container (optional)
	Cmd []string

	// Image used to run the container
	Image  string
	Disk   int64 // Disk in GiB
	Memory int64 // Memory in MiB
	Cpu    float64
	Env    []string // Env variables

	// RestartPolicy for the container ["", "always", "unless-stopped", "on-failure"]
	RestartPolicy string
}

type Docker struct {
	Client *client.Client
	Config Config
}

type DockerResult struct {
	Error       error
	Action      string
	ContainerID string
	Result      string
}

func (d *Docker) Run() DockerResult {
	ctx := context.Background()

	reader, err := d.Client.ImagePull(ctx, d.Config.Image, image.PullOptions{})
	if err != nil {
		log.Printf("Error pulling image %s: %v", d.Config.Image, err)
		return DockerResult{Error: err}
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)

	rp := container.RestartPolicy{
		Name: container.RestartPolicyMode(d.Config.RestartPolicy),
	}
	cr := container.Resources{
		Memory:   d.Config.Memory,
		NanoCPUs: int64(d.Config.Cpu * math.Pow(10, 9)),
	}
	cc := container.Config{
		Image:        d.Config.Image,
		Tty:          false,
		Env:          d.Config.Env,
		ExposedPorts: d.Config.ExposedPorts,
	}
	hc := container.HostConfig{
		RestartPolicy:   rp,
		Resources:       cr,
		PublishAllPorts: true,
	}
	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.ContainerName)
	if err != nil {
		log.Printf("Error creating container using image %s: %v", d.Config.Image, err)
		return DockerResult{Error: err}
	}
	err = d.Client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		log.Printf("Error starting container %s: %v", d.Config.Image, err)
		return DockerResult{Error: err}
	}
	out, err := d.Client.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		log.Printf("Error getting container logs %s: %v", resp.ID, err)
		return DockerResult{Error: err}
	}
	defer out.Close()
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return DockerResult{
		ContainerID: resp.ID,
		Action:      "start",
		Result:      "success",
	}
}

func (d *Docker) Stop(containerId string) DockerResult {
	log.Printf("Attempting to stop container %v", containerId)
	ctx := context.Background()

	err := d.Client.ContainerStop(ctx, containerId, container.StopOptions{})
	if err != nil {
		log.Printf("Error stopping container %s: %v\n", containerId, err)
		return DockerResult{Error: err}
	}
	err = d.Client.ContainerRemove(ctx, containerId, container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	})
	if err != nil {
		log.Printf("Error removing container %s: %v\n", containerId, err)
		return DockerResult{Error: err}
	}
	return DockerResult{
		Action: "stop",
		Result: "success",
		Error:  nil,
	}
}
