package controllers

import (
	"encoding/json"
	"gateway_service/internal"
	"log"
	"net/http"
)

type UserController struct {
	userServiceClient *internal.UserServiceClientStruct
}

func NewUserController(userClient *internal.UserServiceClientStruct) *UserController {
	return &UserController{userServiceClient: userClient}
}

func (controller *UserController) getUserInfo(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("uid")

	// Make a grpc call to the user_service when the user_service is implemented
	pres, err := controller.userServiceClient.GetUserDetails(userId)
	if err != nil {
		log.Fatal("Failed to get user from user_service", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"message": pres.Message, "status": pres.StatusCode}
	json.NewEncoder(w).Encode(response)
}

func (u *UserController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/user/{uid}", u.getUserInfo)
}
