package sockets

import (
	"gateway_service/internal"
	"log"
	"net/http"
	"strings"

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
	kafkaConsumerRef, err := internal.KafkaConsumer(guildName)
	if err != nil {
		log.Println("Something waent wron while creating the room")
		return
	}
	exitSignal := make(chan struct{})
	kafkaconsumer := internal.ConsumerTopic(kafkaConsumerRef, guildName, exitSignal)
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
			for client := range r.clients {
				if client != message.Sender {
					select {
					case client.send <- message.Message:
					default:
						delete(r.clients, client)
						close(client.send)
					}
				}
			}
		case message := <-kafkaconsumer:
			for client := range r.clients {
				if client.id+"-"+channelName != message.Key {
					select {
					case client.send <- []byte(message.Message):
					default:
						delete(r.clients, client)
						close(client.send)
					}
				}
			}
		}
	}
}
