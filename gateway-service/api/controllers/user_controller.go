package controllers

import (
	"encoding/json"
	"gateway_service/internal"
	"log"
	"net/http"
	"user_service/models"
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

func (Controller *UserController) createUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response := map[string]interface{}{"Message": "Something went wrong", "status": http.StatusInternalServerError}
		json.NewEncoder(w).Encode(response)
		return
	}

	pres, err := Controller.userServiceClient.CreateNewUser(&user)
	if err != nil {
		response := map[string]interface{}{"Message": "Something went wrong", "status": http.StatusInternalServerError}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{"Message": pres.Message, "status": http.StatusOK}
	json.NewEncoder(w).Encode(response)
}

func (u *UserController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /user/{uid}", u.getUserInfo)
	mux.HandleFunc("POST /user", u.createUser)
}
