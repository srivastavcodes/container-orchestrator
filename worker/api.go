package worker

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// TODO: turn write the api in grpc

type ErrResponse struct {
	HTTPStatusCode int
	Message        string
}

type WorkerApi struct {
	Host   string
	Port   string
	Worker *Worker
	Router *chi.Mux
}

func (wa *WorkerApi) initRouter() {
	wa.Router = chi.NewRouter()

	wa.Router.Route("/tasks", func(r chi.Router) {
		r.Post("/", wa.StartTaskHandler)
		r.Get("/", wa.GetTasksHandler)
		r.Route("/{taskId}", func(r chi.Router) {
			r.Delete("/", wa.StopTaskHandler)
		})
	})
}

func (wa *WorkerApi) Start() {
	wa.initRouter()
	addr := wa.Host + ":" + wa.Port
	log.Fatal(http.ListenAndServe(addr, wa.Router))
}
