package sockets

import (
	"gateway_service/api/middleware"
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
	mux.Handle("/ws/{roomName}", middleware.Chain(http.HandlerFunc(s.socketHandler), middleware.LoggingMiddleware, middleware.Auth))
	mux.HandleFunc("GET /wss", s.debugHandler)
}
