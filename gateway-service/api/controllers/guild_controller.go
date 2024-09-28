package controllers

import (
	"encoding/json"
	"gateway_service/api/middleware"
	"gateway_service/internal"
	"net/http"
)

type GuildController struct {
	userServiceClient *internal.UserServiceClientStruct
}

func NewGuildController(userServiceClient *internal.UserServiceClientStruct) *GuildController {
	return &GuildController{userServiceClient: userServiceClient}
}

func (g *GuildController) createGuild(w http.ResponseWriter, r *http.Request) {
	guildName := r.PathValue("guildName")
	if guildName == "" {
		response := map[string]interface{}{"Message": "Invalid guid name", "status": http.StatusNotAcceptable}
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{"Message": "Guild Created", "status": http.StatusOK}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (g *GuildController) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /guilds/{guildName}", middleware.Chain(http.HandlerFunc(g.createGuild), middleware.LoggingMiddleware, middleware.Auth))
}
