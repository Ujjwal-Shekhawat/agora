package controllers

import (
	"encoding/json"
	"gateway_service/api/middleware"
	"gateway_service/internal"
	"io"
	"log"
	"net/http"
	proto "proto/user"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

type UserController struct {
	userServiceClient *internal.UserServiceClientStruct
}

func NewUserController(userClient *internal.UserServiceClientStruct) *UserController {
	return &UserController{userServiceClient: userClient}
}

func (controller *UserController) getUserInfo(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("name")

	// Make a grpc call to the user_service when the user_service is implemented
	pres, err := controller.userServiceClient.GetUserDetails(userId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		response := map[string]interface{}{"message": status.Convert(err).Message(), "status": 404}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"message": pres.Message, "status": http.StatusOK}
	json.NewEncoder(w).Encode(response)
}

func (controller *UserController) createUser(w http.ResponseWriter, r *http.Request) {

	user := &proto.User{}

	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request bytes")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := protojson.Unmarshal(requestBytes, user); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pres, err := controller.userServiceClient.CreateNewUser(user)
	if err != nil {
		log.Println(err)
		response := map[string]interface{}{"Message": status.Convert(err).Message(), "status": http.StatusInternalServerError}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	token, err := middleware.GenTok(user.Name)
	if err != nil {
		log.Println("Error generating token for user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{"Message": pres.Message, "token": token, "status": http.StatusOK}
	json.NewEncoder(w).Encode(response)
}

func (controller *UserController) login(w http.ResponseWriter, r *http.Request) {
	loginReq := &proto.LoginReq{}

	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request bytes")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := protojson.Unmarshal(requestBytes, loginReq); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pres, err := controller.userServiceClient.LoginUser(loginReq)
	if err != nil {
		response := map[string]interface{}{"Message": status.Convert(err).Message(), "status": http.StatusInternalServerError}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	token, err := middleware.GenTok(loginReq.Name)
	if err != nil {
		log.Println("Error generating token for user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{"Message": pres.Message, "token": token, "status": http.StatusOK}
	json.NewEncoder(w).Encode(response)
}

func (u *UserController) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /user/{name}", middleware.Chain(http.HandlerFunc(u.getUserInfo), middleware.LoggingMiddleware, middleware.Auth))
	mux.HandleFunc("POST /register/user", u.createUser)
	mux.HandleFunc("POST /login", u.login)
}
