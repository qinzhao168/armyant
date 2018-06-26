package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/yiqinguo/armyant/pkg/models"
	pb "github.com/yiqinguo/armyant/pkg/server"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	//address = "10.200.48.13:50051"
	address     = "10.6.24.110:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewApiserverClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	query(ctx, c)
	//newJob(ctx, c)
	//report(ctx, c)
}

func report(ctx context.Context, c pb.ApiserverClient) {
	str := `{"spec":{"numberOfConnections":100,"testType":"number-of-requests","numberOfRequests":20000,"method":"GET","url":"http://10.6.24.113/","timeoutSeconds":2,"client":"fasthttp"},"result":{"bytesRead":84519020,"bytesWritten":1196048,"timeTakenSeconds":67.17559,"req2xx":20000,"latency":{"mean":193284.89,"stddev":2836210.5,"max":67174264,"percentiles":{"50":305,"75":380,"90":707,"95":90373,"99":1997471}},"rps":{"mean":297.81464,"stddev":1304.9563,"max":12249.383,"percentiles":{"50":0,"75":0,"90":0,"95":2220.106,"99":7550.1147}}}}`

	var stats *models.Stats
	err := json.Unmarshal([]byte(str), &stats)
	resp, err := c.ReportResult(ctx, stats)
	fmt.Println("------------:", resp, err)
}

func query(ctx context.Context, c pb.ApiserverClient) {
	jr := &models.GetJobRequest{}
	jr.JobId = "106277000"
	//jr.JobId = "06277000"
	resp, err := c.GetJobResult(ctx, jr)
	fmt.Println("--------------: ", err)
	body, err := json.Marshal(resp)
	fmt.Println("--------------: ", err, string(body))
}

func newJob(ctx context.Context, c pb.ApiserverClient) {
	id := time.Now().Nanosecond()
	job := &models.Job{
		Id:          fmt.Sprint(id),
		HistoryId:   fmt.Sprint(id),
		Connections: 500,
		Method:      "GET",
		Requests:    100000,
		Url:         "http://10.6.24.113/",
		Name:        "test",
		Master:      address,
	}
	_ = job
	//r, err := c.ReportResult(ctx, &models.Stats{})
	r, err := c.CreateJob(ctx, job)
	//r, err := c.CreateJob(ctx, &models.Job{Name: name})
	//r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
	body, err := json.Marshal(r)
	log.Println(string(body), err)
}
