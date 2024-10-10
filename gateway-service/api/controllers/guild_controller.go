package controllers

import (
	"encoding/json"
	"gateway_service/api/middleware"
	"gateway_service/internal"
	"io"
	"log"
	"net/http"
	proto "proto/guild"
	"regexp"
	"time"

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

	name, channels, members := pres.Name, pres.Channels, pres.Members

	response := map[string]interface{}{"guild name": name, "guild channels": channels, "guild_members": members, "status": http.StatusOK}
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

	var requestData map[string]interface{}
	err = json.Unmarshal(requestBytes, &requestData)
	if err != nil {
		log.Println("Error unmarshaling request bytes:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	creator, ok := r.Context().Value(middleware.AuthUserString).(string)
	if !ok {
		log.Println("Error getting username from token string")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	requestData["creator"] = creator

	responseBytes, err := json.Marshal(requestData)
	if err != nil {
		log.Println("Error marshaling response bytes", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := protojson.Unmarshal(responseBytes, guild); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !regexp.MustCompile(`^[A-Za-z0-9]*$`).MatchString(guild.Name) {
		response := map[string]interface{}{"Message": "Guild name should be alnum nonly", "status": http.StatusInternalServerError}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
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

func (g *GuildController) joinGuild(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.AuthUserString).(string)
	if !ok {
		log.Println(r.Context())
		response := map[string]interface{}{"Message": "Something went wrong token", "status": http.StatusInternalServerError}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Println("User want to join guild", user)

	guildName := r.URL.Query().Get("guildName")
	if guildName == "" {
		response := map[string]interface{}{"Message": "Guild name is not valid", "status": http.StatusInternalServerError}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	guild := &proto.GuildMember{
		Name:      user,
		GuildName: guildName,
	}

	pres, err := g.serviceClient.JoinGuild(guild)
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

func (g *GuildController) getMessage(w http.ResponseWriter, r *http.Request) {
	guildName, channelName := r.URL.Query().Get("guild"), r.URL.Query().Get("channel")
	if guildName == "" || channelName == "" {
		log.Println("guild name or channel name missing here")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userName, ok := r.Context().Value(middleware.AuthUserString).(string)
	if !ok {
		log.Println("Username not of type string?")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	validMember := false
	for _, member := range pres.Members {
		if member == userName {
			validMember = true
		}
	}

	if !validMember {
		response := map[string]interface{}{"Message": "Not authorized to access", "status": http.StatusInternalServerError}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	msgRequest := &proto.GuildMessagesRequest{
		Name:    guildName,
		Channel: channelName,
	}

	messages, err := g.serviceClient.GetMessages(msgRequest)
	if err != nil {
		log.Println(err)
		response := map[string]interface{}{"Message": status.Convert(err).Message(), "status": http.StatusInternalServerError}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	fetchedMessages := []struct {
		Message   string
		Timestamp time.Time
	}{}
	for _, message := range messages.Messages {
		name := message.Key
		if name == userName {
			name = "You"
		}
		fetchedMessages = append(fetchedMessages, struct {
			Message   string
			Timestamp time.Time
		}{Message: name + ": " + message.Value, Timestamp: time.Unix(message.Timestamp.Seconds, int64(message.Timestamp.Nanos))})
	}

	response := map[string]interface{}{"Message": fetchedMessages, "status": http.StatusOK}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (g *GuildController) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /guild/{guildName}", middleware.Chain(http.HandlerFunc(g.getGuild), middleware.LoggingMiddleware, middleware.Auth))
	mux.Handle("POST /guild", middleware.Chain(http.HandlerFunc(g.createGuild), middleware.LoggingMiddleware, middleware.Auth))
	mux.Handle("POST /join", middleware.Chain(http.HandlerFunc(g.joinGuild), middleware.LoggingMiddleware, middleware.Auth))
	mux.Handle("GET /messages", middleware.Chain(http.HandlerFunc(g.getMessage), middleware.LoggingMiddleware, middleware.Auth))
}
