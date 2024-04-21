package worker

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ErrResponse struct {
	HTTPStatusCode int
	Message        string
}

func NewApi(host string, port int, worker *Worker) *Api {
	return &Api{
		Host:   host,
		Port:   port,
		Worker: worker,
	}
}

type Api struct {
	Host   string
	Port   int
	Worker *Worker
	Router *chi.Mux
}

func (a *Api) initRounter() {
	a.Router = chi.NewRouter()
	a.Router.Route("/tasks", func(r chi.Router) {
		r.Post("/", a.StartTaskHanlder)
		r.Get("/", a.GetTaskHandler)
		r.Route("/{taskID}", func(r chi.Router) {
			r.Delete("/", a.StopTaskHandler)
		})
	})
	a.Router.Route("/stats", func(r chi.Router) {
		r.Get("/", a.GetStatsHandler)
	})
}

func (a *Api) Start() {
	a.initRounter()
	addr := fmt.Sprintf("%s:%d", a.Host, a.Port)
	http.ListenAndServe(addr, a.Router)
}
