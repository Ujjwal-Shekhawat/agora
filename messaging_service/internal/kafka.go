package internal

import (
	"log"
	"message_persistance/config"
	"message_persistance/db"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type kafkaConsumer struct {
	consumer *kafka.Consumer
}

var kconsumer *kafkaConsumer = &kafkaConsumer{
	consumer: nil,
}

func InitKafka() error {
	cfg := config.LoadConfig()
	groupId := "message_persistance"

	kafkaConsumerConfig := &kafka.ConfigMap{
		"bootstrap.servers":                  cfg.KafkaBrokers,
		"group.id":                           groupId,
		"auto.offset.reset":                  "earliest",
		"topic.metadata.refresh.interval.ms": 60000,
	}

	consumer, err := kafka.NewConsumer(kafkaConsumerConfig)
	if err != nil {
		log.Println(err)
		return err
	}

	kconsumer.consumer = consumer

	return nil
}

func Start() {
	if kconsumer.consumer == nil {
		log.Fatal("consumer is set nil")
	}

	err := kconsumer.consumer.Subscribe("^guild-.*", nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		event := kconsumer.consumer.Poll(100)
		switch e := event.(type) {
		case *kafka.Message:
			log.Println("Persisting", string(e.Key), string(e.Value))
			guildName := e.TopicPartition.Topic
			key := string(e.Key)
			username := strings.Split(key, "-")[0]
			channelName := strings.Split(key, "-")[1]
			timeStamp := e.Timestamp
			err := db.ExecQuery("INSERT INTO messages (guild_name, user_name, channel_name, user_message, timestamp) VALUES (?, ?, ?, ?, ?)", guildName, username, channelName, string(e.Value), timeStamp)
			if err != nil {
				log.Println(err)
			}
		case kafka.Error:
			log.Println("Encounter kafka error", e)
		}
	}
}
