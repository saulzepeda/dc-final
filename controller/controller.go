package controller

import (
	"fmt"
	"os"
	"time"
	"strconv"
	"strings"

	"go.nanomsg.org/mangos"
	"go.nanomsg.org/mangos/protocol/rep"

	// register transports
	_ "go.nanomsg.org/mangos/transport/all"
)

var sock mangos.Socket

var controllerAddress = "tcp://localhost:40899"

var Workloads = make(map[string]Workload)
type Workload struct{
	ID string
	Filter string
	Name string
	Status string
	Running_jobs int
	Filtered_images []string
}

var Workers = make(map[string]Worker)
type Worker struct{
	Name     string `json:"name"`
	Tags     string `json:"tags"`
	Status   string `json:"status"`
	Usage    int    `json:"usage"`
	URL      string `json:"url"`
	Port     int    `json:"port"`
	Jobs_done int    `json:"jobs_done"`
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}

func Start() {
	var answer []byte
	var sock mangos.Socket
	var err error
	fmt.Println("Hi Controller")
	if sock, err = rep.NewSocket(); err != nil {
		die("can't get new rep socket: %s", err)
	}
	if err = sock.Listen(controllerAddress); err != nil {
		die("can't listen on rep socket: %s", err.Error())
	}

	for{
		if answer, err = sock.Recv();err !=nil{
			break
		}
		new_worker := Worker{}
		info := strings.Split(string(answer)," ")
		new_worker.Name = info[0]
		new_worker.Status = "Available"
		new_worker.Tags = info[1]
		port_int, _ := strconv.Atoi(info[3])
		new_worker.Port = port_int
		jobs_done_int, _ := strconv.Atoi(info[2])
		new_worker.Jobs_done = jobs_done_int
		new_worker.URL = "localhost:" + info[3]
		_, ok := Workers[new_worker.Name]
		if !ok {
			Workers[new_worker.Name] = new_worker
		}

		fmt.Println(Workers[new_worker.Name].Name, " serves in localhost:", Workers[new_worker.Name].Port, "\n")

	}
}

func StatusWorker(name string) { //Update the status worker
	worker, ok := Workers[name]
	if ok {
		if worker.Status == "Available" {
			worker.Status = "Occupied"
		} else {
			worker.Status = "Available"
		}
	}
} 

func UsageWorkers(worker_name string){ //Update the usage of each worker
	worker, ok := Workers[worker_name]
	if ok {
		/*Workers[worker_name].Usage++
		Workers[worker_name].Jobs_done++*/
		worker.Usage++
		worker.Jobs_done++
		Workers[worker_name] = worker
	}
}

