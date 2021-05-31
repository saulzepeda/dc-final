package main

import (
	/*"fmt"
	
	"math/rand"
	"time"*/
	"log"

	"github.com/saulzepeda/dc-final/api"
	"github.com/saulzepeda/dc-final/controller"
	"github.com/saulzepeda/dc-final/scheduler"
)

func main() {
	log.Println("Welcome to the Distributed and Parallel Image Processing System")

	// Start Controller
	go controller.Start()

	// Start Scheduler
	//jobs := make(chan scheduler.Job)
	go scheduler.Start(api.Jobs)
	// Send sample jobs
	//sampleJob := scheduler.Job{Address: "localhost:50051", RPCName: ""}

	// API
	api.Start()
	
	/*for {
		sampleJob.RPCName = fmt.Sprintf("hello-%v", rand.Intn(10000))
		jobs <- sampleJob
		time.Sleep(time.Second * 5)
	}*/
}
