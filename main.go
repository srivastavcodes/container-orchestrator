package main

import (
	"fmt"
	"orchestrator/task"
	"orchestrator/worker"
	"time"

	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	db := make(map[uuid.UUID]*task.Task)

	wr := worker.Worker{
		Queue: *queue.New(),
		Db:    db,
	}
	currTask := task.Task{
		ID:    uuid.New(),
		Name:  "test-container-1",
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}
	// first time the wr will see the task
	fmt.Println("starting task")
	wr.AddTask(currTask)

	result := wr.RunTask()
	if result.Error != nil {
		panic(fmt.Sprintf("yo whatdafuk? err=%v", result.Error))
	}
	currTask.ContainerID = result.ContainerID

	fmt.Printf("task %s is running on container %s\n", currTask.ID, currTask.ContainerID)
	fmt.Println("sleepy time")
	time.Sleep(30 * time.Second)

	fmt.Printf("stopping task %s\n", currTask.ID)
	currTask.State = task.Completed
	wr.AddTask(currTask)

	result = wr.RunTask()
	if result.Error != nil {
		panic(fmt.Sprintf("yo whatdafuk? err=%v", result.Error))
	}
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
