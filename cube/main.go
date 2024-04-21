package main

import (
	"cube/task"
	"cube/worker"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test-container-1",
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_USER=cube",
			"POSTGRES_PASSWORD=secret",
		},
	}

	d := task.NewDocker(&c)

	result := d.Run()
	if result.Error != nil {
		fmt.Printf("Error: %v\n", result.Error)
		return nil, &result
	}

	fmt.Printf("Container is %s is running with config: %v\n", result.ContainerID, c)
	return d, &result
}

func stopContainer(d *task.Docker, id string) *task.DockerResult {
	result := d.Stop(id)
	if result.Error != nil {
		fmt.Printf("Error: %v\n", result.Error)
		return &result
	}

	fmt.Printf("Container %s is stopped and removed\n", id)
	return &result
}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				log.Printf("Error running task: %v\n", result.Error)
			}
		} else {
			log.Println("No task to process currently!")
		}

		log.Println("Sleeping for 10 seconds.")
		time.Sleep(10 * time.Second)
	}
}

// CUBE_HOST=localhost CUBE_PORT=5555 DOCKER_API_VERSION=1.44 go run main.go

func main() {
	host := os.Getenv("CUBE_HOST")
	port, _ := strconv.Atoi(os.Getenv("CUBE_PORT"))

	fmt.Println("Starting Cube Worker...")

	w := worker.NewWorker()
	api := worker.NewApi(host, port, w)

	go runTasks(w)
	go w.CollectStats()
	api.Start()
}
