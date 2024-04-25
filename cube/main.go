package main

import (
	"cube/manager"
	"cube/worker"
	"fmt"
	"os"
	"strconv"
)

func main() {
	mhost := os.Getenv("CUBE_MANAGER_HOST")
	mport, _ := strconv.Atoi(os.Getenv("CUBE_MANAGER_PORT"))

	whost := os.Getenv("CUBE_WORKER_HOST")
	wport, _ := strconv.Atoi(os.Getenv("CUBE_WORKER_PORT"))

	fmt.Println("Starting Cube Worker...")

	w1 := worker.NewWorker()
	wapi1 := worker.NewApi(whost, wport, w1)

	w2 := worker.NewWorker()
	wapi2 := worker.NewApi(whost, wport+1, w2)

	w3 := worker.NewWorker()
	wapi3 := worker.NewApi(whost, wport+2, w3)

	go w1.RunTasks()
	go w1.CollectStats()
	go w1.UpdateTasks()
	go wapi1.Start()

	go w2.RunTasks()
	go w2.CollectStats()
	go w2.UpdateTasks()
	go wapi2.Start()

	go w3.RunTasks()
	go w3.CollectStats()
	go w3.UpdateTasks()
	go wapi3.Start()

	fmt.Println("Starting Cube Manager...")
	workers := []string{
		fmt.Sprintf("%s:%d", whost, wport),
		fmt.Sprintf("%s:%d", whost, wport+1),
		fmt.Sprintf("%s:%d", whost, wport+2),
	}
	m := manager.New(workers, "epvm")
	mapi := manager.NewApi(mhost, mport, m)

	go m.ProcessTasks()
	go m.UpdateTasks()
	go m.DoHealthChecks()

	mapi.Start()
}
