package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/saulzepeda/dc-final/proto"
	"github.com/saulzepeda/dc-final/controller"
	"go.nanomsg.org/mangos"
	"go.nanomsg.org/mangos/protocol/req"
	"google.golang.org/grpc"

	// register transports
	_ "go.nanomsg.org/mangos/transport/all"

	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/effect"
    "github.com/anthonynsimon/bild/imgio"
)

var (
	defaultRPCPort = 50051
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

var (
	controllerAddress = ""
	workerName        = ""
	tags              = ""
	jobs_done         = ""
	port              = ""
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("RPC: Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) BlurEffect(ctx context.Context, in *pb.ImageRequest) (*pb.ImageReply, error) {
	new_filename := fmt.Sprintf("%v", in.Img.Index) + "_blur.png"
	path := "../images/" + in.Img.WlId + "/"
	final_path := path + new_filename

	img, err := imgio.Open(in.GetImg().Filepath)

	if err != nil {
		return &pb.ImageReply{Message: "Could not open image " + workerName}, nil 
	}

	blur_img := blur.Gaussian(img, 3.0)

	if err := imgio.Save(final_path, blur_img, imgio.PNGEncoder()); err != nil {
		fmt.Println(err)
		return &pb.ImageReply{Message: "Blur error " + workerName}, nil
	}

	return &pb.ImageReply{Message: fmt.Sprintf("Filename: %v, Workload ID: %v", new_filename, in.Img.WlId)}, nil
}

func (s *server) GrayscaleEffect(ctx context.Context, in *pb.ImageRequest) (*pb.ImageReply, error) {
	new_filename := fmt.Sprintf("%v", in.Img.Index) + "_grayscale.png"
	path := "../images/" + in.Img.WlId + "/"
	final_path := path + new_filename

	img, err := imgio.Open(in.GetImg().Filepath)

	if err != nil {
		return &pb.ImageReply{Message: "Could not open image " + workerName}, nil 
	}

	grayscale_img := effect.Grayscale(img)
	if err := imgio.Save(final_path, grayscale_img, imgio.PNGEncoder()); err != nil {
		fmt.Println(err)
		return &pb.ImageReply{Message: "Grayscale error " + workerName}, nil
	}

	controller.UsageWorkers(workerName)

	return &pb.ImageReply{Message: fmt.Sprintf("Filename: %v, Workload ID: %v", new_filename, in.Img.WlId)}, nil

	
}

func init() {
	flag.StringVar(&controllerAddress, "controller", "tcp://localhost:40899", "Controller address")
	flag.StringVar(&workerName, "worker-name", "hard-worker", "Worker Name")
	flag.StringVar(&tags, "tags", "gpu,superCPU,largeMemory", "Comma-separated worker tags")
}

// joinCluster is meant to join the controller message-passing server
func joinCluster() {
	var sock mangos.Socket
	var err error
	var msg []byte

	if sock, err = req.NewSocket(); err != nil {
		die("can't get new req socket: %s", err.Error())
	}

	log.Printf("Connecting to controller on: %s", controllerAddress)
	if err = sock.Dial(controllerAddress); err != nil {
		die("can't dial on req socket: %s", err.Error())
	}
	info := fmt.Sprintf("%v %v %v %v %v %v", workerName, tags, jobs_done, defaultRPCPort)
		if err = sock.Send([]byte(info)); err != nil {
			die("Cannot send: %s", err.Error())
		}
		log.Printf("Message-Passing: Worker(%s): Received %s\n", workerName, string(msg))
		time.Sleep(time.Second)
}

func getAvailablePort() int {
	port := defaultRPCPort
	for {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
		if err != nil {
			port = port + 1
			continue
		}
		ln.Close()
		break
	}
	return port
}

func main() {
	flag.Parse()

	// Subscribe to Controller
	go joinCluster()

	// Setup Worker RPC Server
	rpcPort := getAvailablePort()
	log.Printf("Starting RPC Service on localhost:%v", rpcPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", rpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
