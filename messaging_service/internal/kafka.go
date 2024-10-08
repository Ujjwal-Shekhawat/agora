package internal

import (
	"log"
	"message_persistance/config"

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
		"debug":                              "consumer,topic,metadata",
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
		case kafka.Error:
			log.Println("Encounter kafka error", e)
		}
	}
}
