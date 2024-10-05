package sockets

import (
	"gateway_service/api/middleware"
	"log"
	"net/http"
	proto "proto/guild"
)

type Hub struct {
	rooms map[string]*Room
}

var hub = Hub{
	rooms: make(map[string]*Room),
}

func serveSocket(user string, roomName string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	room, ok := hub.rooms[roomName]
	if !ok {
		room = createRoom()
		hub.rooms[roomName] = room
		go room.run()
	}

	client := &Client{id: user, room: room, conn: conn, send: make(chan []byte)}

	room.register <- client

	go client.readPump()
	go client.writePump()
}

func (s *SocketController) socketHandler(w http.ResponseWriter, r *http.Request) {
	roomName := r.PathValue("roomName")
	if roomName == "" {
		http.Error(w, "Room name is required", http.StatusBadRequest)
		return
	}

	channelName := r.URL.Query().Get("channel")
	if channelName == "" {
		http.Error(w, "Channel not correct", http.StatusInternalServerError)
		return
	}

	user, ok := r.Context().Value(middleware.AuthUserString).(string)
	if !ok {
		log.Println("Something wrong with the user token")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	guild := &proto.Guild{
		Name: roomName,
	}

	pres, err := s.userServiceClient.FetchGuild(guild)
	if err != nil {
		log.Println("Error fetching guild", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	validMember := false
	for _, member := range pres.Members {
		if user == member {
			validMember = true
		}
	}

	if !validMember {
		log.Println("User not in guild", err)
		http.Error(w, "Not registered member", http.StatusInternalServerError)
		return
	}

	vaildChannel := false
	for _, channel := range pres.Channels {
		if channelName == channel {
			vaildChannel = true
		}
	}

	if !vaildChannel {
		log.Println("User not in guild", err)
		http.Error(w, "Channel name dne", http.StatusInternalServerError)
		return
	}

	serveSocket(user, roomName, w, r)
}

// remove this later
func (s *SocketController) debugHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(hub.rooms)
	w.Write([]byte("Done"))
}
