package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"orchestrator/task"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (wa *WorkerApi) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024))

	dec.DisallowUnknownFields()
	var event task.TaskEvent

	if err := dec.Decode(&event); err != nil {
		errMsg := fmt.Sprintf("Error unmarshalling body: %v\n", err)
		log.Printf(errMsg)

		w.WriteHeader(http.StatusBadRequest)
		errRes := ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			Message:        errMsg,
		}
		json.NewEncoder(w).Encode(errRes)
		return
	}
	wa.Worker.AddTask(event.Task)

	log.Printf("Added Task %v\n", event.Task)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(event.Task)
}

func (wa *WorkerApi) GetTasksHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wa.Worker.GetTasks())
}

func (wa *WorkerApi) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "taskId")
	if taskId == "" {
		log.Printf("No \"taskId\" provided in request\n")
		w.WriteHeader(http.StatusBadRequest)
	}
	tId, _ := uuid.Parse(taskId)

	taskToStop, ok := wa.Worker.Db[tId]
	if !ok {
		log.Printf("No task with ID %v found\n", tId)
		w.WriteHeader(http.StatusBadRequest)
	}
	// copying because if we change the pointer to completed, the verification
	// for the state transition will fail; can't do completed -> completed
	taskCopy := *taskToStop

	taskCopy.State = task.Completed
	wa.Worker.AddTask(taskCopy)

	log.Printf("Added task %v to stop container %v\n", taskToStop.ID, taskToStop.ContainerID)
	w.WriteHeader(http.StatusNoContent)
}
