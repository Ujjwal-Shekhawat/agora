package sockets

import (
	"gateway_service/internal"
	"net/http"
)

type SocketController struct {
	userServiceClient *internal.ServiceClientStruct
}

func NewSocketController(userClient *internal.ServiceClientStruct) *SocketController {
	return &SocketController{userServiceClient: userClient}
}

func (s *SocketController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ws/{roomName}", s.socketHandler)
	mux.HandleFunc("GET /wss", s.debugHandler)
}
