package controllers

import (
	"encoding/json"
	"gateway_service/api/middleware"
	"gateway_service/internal"
	"io"
	"log"
	"net/http"
	proto "proto/guild"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

type GuildController struct {
	serviceClient *internal.ServiceClientStruct
}

func NewGuildController(serviceClient *internal.ServiceClientStruct) *GuildController {
	return &GuildController{serviceClient: serviceClient}
}

func (g *GuildController) getGuild(w http.ResponseWriter, r *http.Request) {
	guildName := r.PathValue("guildName")

	guild := &proto.Guild{
		Name: guildName,
	}

	pres, err := g.serviceClient.FetchGuild(guild)
	if err != nil {
		log.Println(err)
		response := map[string]interface{}{"Message": status.Convert(err).Message(), "status": http.StatusInternalServerError}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	name, channels := pres.Name, pres.Channels

	response := map[string]interface{}{"guild name": name, "guild channels": channels, "status": http.StatusOK}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (g *GuildController) createGuild(w http.ResponseWriter, r *http.Request) {

	guild := &proto.Guild{}

	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request bytes")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(string(requestBytes))

	if err := protojson.Unmarshal(requestBytes, guild); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pres, err := g.serviceClient.MakeGuild(guild)
	if err != nil {
		log.Println(err)
		response := map[string]interface{}{"Message": status.Convert(err).Message(), "status": http.StatusInternalServerError}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{"Message": pres.Message, "status": http.StatusOK}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (g *GuildController) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /guild/{guildName}", middleware.Chain(http.HandlerFunc(g.getGuild), middleware.LoggingMiddleware, middleware.Auth))
	mux.Handle("POST /guild", middleware.Chain(http.HandlerFunc(g.createGuild), middleware.LoggingMiddleware, middleware.Auth))
}
