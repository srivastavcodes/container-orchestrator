package main

import (
	"fmt"
	"orchestrator/manager"
	"orchestrator/node"
	"orchestrator/task"
	"orchestrator/worker"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	t := task.Task{
		ID:     uuid.New(),
		Name:   "Task-1",
		State:  task.Pending,
		Image:  "Image-1",
		Memory: 1024,
		Disk:   1,
	}
	te := task.TaskEvent{
		ID:        uuid.New(),
		State:     task.Pending,
		Timestamp: time.Now(),
		Task:      t,
	}
	fmt.Printf("task: %v\n", t)
	fmt.Printf("task event: %v\n", te)

	w := worker.Worker{
		Name:  "worker-1",
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	fmt.Printf("worker: %v\n", w)
	w.CollectStats()
	w.RunTask()
	w.StartTask()
	w.StopTask()

	m := manager.Manager{
		Pending: *queue.New(),
		TaskDb:  make(map[string][]*task.Task),
		EventDb: make(map[string][]*task.TaskEvent),
		Workers: []string{w.Name},
	}
	fmt.Printf("manager: %v\n", m)
	m.SelectWorker()
	m.UpdateTasks()
	m.SendWork()

	n := node.Node{
		Name:   "Node-1",
		Ip:     "192.168.1.1",
		Cores:  4,
		Memory: 1024,
		Disk:   25,
		Role:   "worker",
	}
	fmt.Printf("node: %v\n", n)

	fmt.Printf("create a test container\n")

	dockerTask, createResult := createContainer()
	if createResult.Error != nil {
		fmt.Printf("%v\n", createResult.Error)
		os.Exit(1)
	}
	time.Sleep(time.Second * 10)

	fmt.Printf("stopping container %s\n", createResult.ContainerID)
	_ = stopContainer(dockerTask, createResult.ContainerID)
}

func createContainer() (*task.Docker, *task.DockerResult) {
	config := task.Config{
		ContainerName: "test-container-1",
		Image:         "postgres:13",
		Env: []string{
			"POSTGRES_USER=cube",
			"POSTGRES_PASSWORD=secret",
		},
	}
	dockerClient, _ := client.NewClientWithOpts(client.FromEnv)

	docker := task.Docker{
		Client: dockerClient,
		Config: config,
	}
	result := docker.Run()

	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil, &result
	}
	fmt.Printf("Container %s is running with config %v\n", result.ContainerID, config)
	return &docker, &result
}

func stopContainer(docker *task.Docker, containerID string) *task.DockerResult {
	result := docker.Stop(containerID)

	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil
	}
	fmt.Printf("Container %s has been stopped and removed\n", result.ContainerID)
	return &result
}
