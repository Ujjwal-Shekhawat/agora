package main

import (
	"log"
	"user_service/config"
	"user_service/db"
	"user_service/internal"
)

func main() {
	cfg := config.LoadConfig()

	if err := db.InitSession(); err != nil {
		log.Fatal(err)
	}

	log.Println("Service started on port: ", cfg.ServerPort)
	if err := internal.StartGrpcServer(cfg); err != nil {
		log.Fatal("Unable to start the grpc server")
	}

}
