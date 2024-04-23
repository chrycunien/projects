package manager

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ErrResponse struct {
	HTTPStatusCode int
	Message        string
}

func NewApi(host string, port int, manager *Manager) *Api {
	return &Api{
		Host:    host,
		Port:    port,
		Manager: manager,
	}
}

type Api struct {
	Host    string
	Port    int
	Manager *Manager
	Router  *chi.Mux
}

func (a *Api) initRouter() {
	a.Router = chi.NewRouter()
	a.Router.Route("/tasks", func(r chi.Router) {
		r.Post("/", a.StartTaskHandler)
		r.Get("/", a.GetTaskHandler)
		r.Route("/{taskID}", func(r chi.Router) {
			r.Delete("/", a.StopTaskHandler)
		})
	})
}

func (a *Api) Start() {
	a.initRouter()
	addr := fmt.Sprintf("%s:%d", a.Host, a.Port)
	http.ListenAndServe(addr, a.Router)
}
