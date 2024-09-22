package sockets

import (
	"log"
	"net/http"
)

type Hub struct {
	rooms map[string]*Room
}

var hub = Hub{
	rooms: make(map[string]*Room),
}

func serveSocket(roomName string, w http.ResponseWriter, r *http.Request) {
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

	client := &Client{room: room, conn: conn, send: make(chan []byte)}

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

	serveSocket(roomName, w, r)
}

// remove this later
func (s *SocketController) debugHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(hub.rooms)
	w.Write([]byte("Done"))
}
