package sockets

import (
	"gateway_service/internal"
	"net/http"
)

type SocketController struct {
	userServiceClient *internal.UserServiceClientStruct
}

func NewSocketController(userClient *internal.UserServiceClientStruct) *SocketController {
	return &SocketController{userServiceClient: userClient}
}

func (s *SocketController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ws/{roomName}", s.socketHandler)
	mux.HandleFunc("GET /wss", s.debugHandler)
}
