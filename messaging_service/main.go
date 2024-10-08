package main

import (
	"log"
	"message_persistance/db"
	"message_persistance/internal"
)

func main() {
	if err := db.InitSession(); err != nil {
		log.Fatal("Failed to init cassandra session")
	}

	err := internal.InitKafka()
	if err != nil {
		log.Fatal("Failed to init kafka consumer")
	}

	log.Println("Started persistence service")
	internal.Start()
}
