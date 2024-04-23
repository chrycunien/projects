package main

import (
	"cube/manager"
	"cube/task"
	"cube/worker"
	"fmt"
	"os"
	"strconv"
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

// CUBE_HOST=localhost CUBE_PORT=5555 DOCKER_API_VERSION=1.44 go run main.go

func main() {
	mhost := os.Getenv("CUBE_MANAGER_HOST")
	mport, _ := strconv.Atoi(os.Getenv("CUBE_MANAGER_PORT"))

	whost := os.Getenv("CUBE_WORKER_HOST")
	wport, _ := strconv.Atoi(os.Getenv("CUBE_WORKER_PORT"))

	fmt.Println("Starting Cube Worker...")

	w := worker.NewWorker()
	wapi := worker.NewApi(whost, wport, w)

	go w.RunTasks()
	go w.CollectStats()
	go w.UpdateTasks()
	go wapi.Start()

	fmt.Println("Starting Cube Manager...")
	workers := []string{fmt.Sprintf("%s:%d", whost, wport)}
	m := manager.New(workers)
	mapi := manager.NewApi(mhost, mport, m)

	go m.ProcessTasks()
	go m.UpdateTasks()
	go m.DoHealthChecks()

	mapi.Start()
}
