package main

import (
	"gateway_service/api/controllers"
	"gateway_service/config"
	"gateway_service/internal"
	"gateway_service/routes"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	userServiceClient, err := internal.GetUserServiceClient(cfg.UserServiceAddr)
	if err != nil {
		log.Fatal("Failed to create user service client: ", err)
	}

	userController := controllers.NewUserController(userServiceClient)

	router := http.NewServeMux()

	routes.RegisterControllers(router, userController)

	log.Print("Server started on address: ", cfg.ServerPort)

	if err := http.ListenAndServe(cfg.ServerPort, router); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
