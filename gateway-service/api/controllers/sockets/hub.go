package sockets

import (
	"gateway_service/api/middleware"
	"gateway_service/internal"
	"log"
	"net/http"
	proto "proto/guild"
	"sync"
)

type KafkaCache struct {
	CirculairBuffer map[string]*KafkaBuffer
	mu              sync.Mutex
}

var KafkaCacheManager KafkaCache = KafkaCache{
	CirculairBuffer: make(map[string]*KafkaBuffer),
}

type KafkaBuffer struct {
	Message []internal.KafkaConsumerEvent
	Head    int
}

const (
	size = 10
)

func InitKafkaBuffer(cacheId string) *KafkaBuffer {
	newBuffer := &KafkaBuffer{
		Message: make([]internal.KafkaConsumerEvent, 0, size),
	}
	KafkaCacheManager.mu.Lock()
	_, ok := KafkaCacheManager.CirculairBuffer[cacheId]
	if !ok {
		KafkaCacheManager.CirculairBuffer[cacheId] = newBuffer
	} else {
		log.Fatal("Cant init an exsisting buffer")
	}
	KafkaCacheManager.mu.Unlock()

	return newBuffer
}

func (kafkaBuffer *KafkaBuffer) Add(message internal.KafkaConsumerEvent) {
	if len(kafkaBuffer.Message) < size {
		kafkaBuffer.Message = append(kafkaBuffer.Message, message)
	} else {
		kafkaBuffer.Message[kafkaBuffer.Head] = message
		kafkaBuffer.Head = (kafkaBuffer.Head + 1) % size
	}
}

func (kafkaBuffer *KafkaBuffer) Print() []string {
	x := []string{}
	for _, v := range kafkaBuffer.Message {
		x = append(x, v.Key+": "+v.Message)
	}
	return x
}

type Hub struct {
	rooms map[string]*Room
}

var HubManager = Hub{
	rooms: make(map[string]*Room),
}

func serveSocket(s *SocketController, user string, roomName string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	room, ok := HubManager.rooms[roomName]
	if !ok {
		room = createRoom(roomName)
		HubManager.rooms[roomName] = room
		go room.run(s)
	}

	client := &Client{id: user, room: room, conn: conn, send: make(chan []byte, 1000)}

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

	roomName = roomName + "-" + channelName

	serveSocket(s, user, roomName, w, r)
}

// remove this later
func (s *SocketController) debugHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(HubManager.rooms)
	w.Write([]byte("Done"))
}
