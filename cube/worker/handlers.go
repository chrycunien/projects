package worker

import (
	"cube/stats"
	"cube/task"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a *Api) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if a.Worker.Stats != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(*a.Worker.Stats)
		return
	}

	w.WriteHeader(200)
	stats := stats.GetStats()
	json.NewEncoder(w).Encode(stats)
}

func (a *Api) StartTaskHanlder(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	te := task.TaskEvent{}
	if err := d.Decode(&te); err != nil {
		msg := fmt.Sprintf("Error unmarshalling body: %v\n", err)
		log.Printf("%s", msg)
		w.WriteHeader(http.StatusBadRequest)
		resp := ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			Message:        msg,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	a.Worker.AddTask(te.Task)
	log.Printf("Added Task %v\n", te.Task.ID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(te.Task)
}

func (a *Api) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(a.Worker.GetTasks())
}

func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		log.Println("No taskID passed in the request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tid, err := uuid.Parse(taskID)
	if err != nil {
		log.Println("Invalid taskID form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, ok := a.Worker.Db[tid]
	if !ok {
		log.Printf("No task with ID %v found\n", tid)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	taskToStop := a.Worker.Db[tid]
	taskCopy := *taskToStop
	taskCopy.State = task.Completed
	a.Worker.AddTask(taskCopy)
	log.Printf("Added task %v to stop container %v\n", taskToStop.ID, taskToStop.ContainerID)
	w.WriteHeader(http.StatusNoContent)
}
