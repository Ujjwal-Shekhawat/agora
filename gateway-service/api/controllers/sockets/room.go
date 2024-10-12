package sockets

import (
	"fmt"
	"gateway_service/internal"
	"log"
	"net/http"
	proto "proto/guild"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/gogo/protobuf/types"
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
		broadcast:  make(chan *Message, 1000),
		register:   make(chan *Client, 1000),
		unregister: make(chan *Client, 1000),
	}
}

func (r *Room) run(s *SocketController) {
	channelName := strings.Join(strings.Split(r.name, "-")[1:], "-")
	guildName := strings.Join(strings.Split(r.name, "-")[:1], "-")
	gID := time.Now().UnixMilli()
	buffer := InitKafkaBuffer(fmt.Sprint(gID))
	kafkaConsumerRef, err := internal.KafkaConsumer(fmt.Sprint(gID))
	if err != nil {
		log.Println("Something waent wron while creating the room")
		return
	}
	exitSignal := make(chan struct{})
	kafkaconsumer := internal.ConsumerTopic(kafkaConsumerRef, "guild-"+guildName, exitSignal)
	defer kafkaConsumerRef.Close()
	defer close(exitSignal)
	defer delete(HubManager.rooms, r.name)
	defer delete(KafkaCacheManager.CirculairBuffer, fmt.Sprint(gID))
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			messages, err := s.userServiceClient.GetMessages(&proto.GuildMessagesRequest{
				Name:    guildName,
				Channel: channelName,
			})
			if err != nil {
				delete(r.clients, client)
				close(client.send)
				client.conn.Close()
				continue
			}

			dbMessages := []internal.KafkaConsumerEvent{}
			for _, message := range messages.Messages {
				user := message.Key
				if user == client.id {
					user = "You"
				}
				t, err := types.TimestampFromProto(message.Timestamp)
				if err != nil {
					delete(r.clients, client)
					close(client.send)
					client.conn.Close()
					continue
				}

				dbMessages = append(dbMessages, internal.KafkaConsumerEvent{
					Key:       user,
					Message:   message.Value,
					TimeStamp: t,
				})
			}
			bufferCopy := make([]internal.KafkaConsumerEvent, 0, len(buffer.Message))
			copy(bufferCopy, buffer.Message)

			sort.Slice(bufferCopy, func(i, j int) bool {
				return bufferCopy[i].TimeStamp.Before(bufferCopy[j].TimeStamp)
			})

			log.Println("1")
			dbMessageOt := len(dbMessages) - 1
			if dbMessageOt < 0 {
				dbMessageOt = 0
			}

			cutOffIndex := 0
			if len(dbMessages) > 0 {
				latestDatabaseMessage := dbMessages[dbMessageOt:][0].TimeStamp

				log.Println("2")
				for i, v := range bufferCopy {
					if v.TimeStamp.After(latestDatabaseMessage) {
						cutOffIndex = i
						break
					}
				}
				slices.Reverse(dbMessages)
			}

			for _, v := range dbMessages {
				client.send <- []byte(v.Key + ": " + v.Message)
			}
			for _, v := range bufferCopy[cutOffIndex:] {
				client.send <- []byte(v.Key + ": " + v.Message)
			}
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
				client.conn.Close()
			}
		case message := <-r.broadcast:
			internal.PublishMessage("guild-"+guildName, []byte(message.Sender.id+"-"+channelName), message.Message)
		case message := <-kafkaconsumer:
			buffer.Add(message)
			log.Println(buffer.Print())
			for client := range r.clients {
				if client.id+"-"+channelName != message.Key {
					userName := strings.Split(message.Key, "-")[0]
					select {
					case client.send <- []byte(userName + ": " + message.Message):
					default:
						log.Printf("Client %s is unable to receive message, disconnecting.", client.id)

						delete(r.clients, client)
						close(client.send)
					}
				}
			}
		default:
			if len(r.clients) == 0 {
				log.Println("Closing room", r.name, "no clients connected")
				return
			}
		}
	}
}
