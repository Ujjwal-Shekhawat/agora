package sockets

import (
	"fmt"
	"gateway_service/internal"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Message []byte
	Sender  *Client
}

type Room struct {
	name       string
	clients    map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func createRoom(name string) *Room {
	return &Room{
		name:       name,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (r *Room) run() {
	channelName := strings.Join(strings.Split(r.name, "-")[1:], "-")
	guildName := strings.Join(strings.Split(r.name, "-")[:1], "-")
	gID := time.Now().UnixMilli()
	kafkaConsumerRef, err := internal.KafkaConsumer(fmt.Sprint(gID))
	if err != nil {
		log.Println("Something waent wron while creating the room")
		return
	}
	exitSignal := make(chan struct{})
	kafkaconsumer := internal.ConsumerTopic(kafkaConsumerRef, "guild-"+guildName, exitSignal)
	defer kafkaConsumerRef.Close()
	defer close(exitSignal)
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
		case message := <-r.broadcast:
			internal.PublishMessage("guild-"+guildName, []byte(message.Sender.id+"-"+channelName), message.Message)
		case message := <-kafkaconsumer:
			for client := range r.clients {
				if client.id+"-"+channelName != message.Key {
					userName := strings.Split(message.Key, "-")[0]
					select {
					case client.send <- []byte(userName + ": " + message.Message):
					default:
						delete(r.clients, client)
						close(client.send)
					}
				}
			}
		}
	}
}
