package scheduler

import (
	"strings"
	"strconv"
	"context"
	"log"
	"time"
	"fmt"

	"github.com/saulzepeda/dc-final/controller"
	pb "github.com/saulzepeda/dc-final/proto"
	"google.golang.org/grpc"
)

//const (
//	address     = "localhost:50051"
//	defaultName = "world"
//)

type Job struct {
	Address string
	RPCName string
	Filepath string
	Wl_id string
	Filter_type string
	Actual_worker string
}

func schedule(job Job) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(job.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	fmt.Println("After Connection")
	
	//Update the status worker to Occupied
	controller.StatusWorker(job.Actual_worker)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fmt.Println("After ctx")

	wl_id := job.Wl_id
	path_slices := strings.Split(job.Filepath, "/")
	img_id := strings.Split(path_slices[len(path_slices)-1], "_")
	id_int, _ := strconv.Atoi(img_id[0])
	img := pb.Image{
		WlId: wl_id, 
		WlName: controller.Workloads[wl_id].Name,
		Index: int64(id_int), 
		Filepath: job.Filepath,
		Filter: job.Filter_type,
	}
	if job.Filter_type == "grayscale" {
		fmt.Println("Before GS")
		r, err := c.GrayscaleEffect(ctx, &pb.ImageRequest{Img: &img })
		if err != nil {
			log.Fatalf("Couldn't proccess image: %v", err)
		}
		msg_answer := strings.Split(r.GetMessage(), ",")

		//Update the Workload
		updated_wL := controller.Workload{}
		wl_id = msg_answer[1]
		original_wl := controller.Workloads[wl_id]
		updated_wL = controller.Workload{
			ID: original_wl.ID,
			Filter: original_wl.Filter,
			Name: original_wl.Name,
			Status: "completed",
			Running_jobs: original_wl.Running_jobs + 1,
			Filtered_images: original_wl.Filtered_images,
		}
		
		img_id := strings.Split(strings.Split(msg_answer[0], ".")[0], "_")
		updated_wL.Filtered_images = append(updated_wL.Filtered_images, img_id[0])
		controller.Workloads[wl_id] = updated_wL

	} else if job.Filter_type == "blur" {
		fmt.Println("Before Bl")
		r, err := c.BlurEffect(ctx, &pb.ImageRequest{Img: &img })
		if err != nil {
			log.Fatalf("Couldn't proccess image: %v", err)
		}
		fmt.Println("Image Processed w/bl")
		msg_answer := strings.Split(r.GetMessage(), ",")

		//Update the Workload
		updated_wL := controller.Workload{}
		wl_id = msg_answer[1]
		original_wl := controller.Workloads[wl_id]
		updated_wL = controller.Workload{
			ID: original_wl.ID,
			Filter: original_wl.Filter,
			Name: original_wl.Name,
			Status: "completed",
			Running_jobs: original_wl.Running_jobs + 1,
			Filtered_images: original_wl.Filtered_images,
		}

		img_id := strings.Split(strings.Split(msg_answer[0], ".")[0], "_")
		updated_wL.Filtered_images = append(updated_wL.Filtered_images, img_id[0])
		controller.Workloads[wl_id] = updated_wL
	}

	//Update the status worker to Available
	controller.StatusWorker(job.Actual_worker)
}

func Start(jobs chan Job) error {
	fmt.Println("Hi Scheduler")
	for {
		job := <-jobs
		fmt.Println("Aqui estoy")
		worker := controller.Worker{}
		port := 0
		for _, temp_worker := range controller.Workers {
			if temp_worker.Status == "Available" {
				port = temp_worker.Port
				worker = temp_worker
			}
		}
		controller.UsageWorkers(worker.Name)
		job.Actual_worker = worker.Name
		if port == 0 {
			return nil
		}

		job.Address = "localhost:" + strconv.Itoa(port)
		fmt.Println("Before Schedule")
		schedule(job)
	}
	return nil
}
