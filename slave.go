package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/yiqinguo/armyant/pkg"
	"github.com/yiqinguo/armyant/pkg/models"
	"github.com/yiqinguo/armyant/pkg/server"

	"golang.org/x/net/context"
)

func main() {
	job := pkg.ParseArgs()

	gclient, err := server.NewGrpcClient(job.Master)
	if err != nil {
		log.Fatalf("connection master error: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go gclient.ReportStatus(ctx, &models.Status{})

	client := pkg.NewBombardierClient(job)
	stats, err := client.Run()
	if err != nil {
		log.Printf("exec benchmark error: %v", err)
		return
	}
	body, _ := json.Marshal(stats)
	log.Println(string(body))
	stats.JobId = job.Id
	stats.HistoryId = job.HistoryId
	stats.InstancdId = os.Getenv("POD_NAME")

	err = gclient.ReportResult(context.Background(), &stats)
	if err != nil {
		log.Println("report benchmark result error: %v", err)
	}

}
