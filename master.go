package main

import (
	"log"

	"github.com/yiqinguo/armyant/pkg"
	"github.com/yiqinguo/armyant/pkg/server"
)

func main() {

	args := pkg.ParseMasterArgs()

	bserver := server.NewBenchmarkServer(args)
	err := bserver.Run(args)
	log.Println(err)
	log.Printf("start grpc server: %v", err)
}
