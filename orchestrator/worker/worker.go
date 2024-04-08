package worker

import (
	"cube/task"
	"fmt"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]*task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	fmt.Println("I will collect statistics!")
}

func (w *Worker) RunTask() {
	fmt.Println("I will start or stop a task!")
}

func (w *Worker) StartTask() {
	fmt.Println("I will start a task!")
}

func (w *Worker) StopTask(t task.Task) *task.DockerResult {
	config := task.NewConfig(&t)

}
