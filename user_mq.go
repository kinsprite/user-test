package main

import (
	"context"
	"log"
	"os"
	"strconv"

	kafka "github.com/segmentio/kafka-go"
)

var kafkaBrokerURL = "kafka:9092"
var w *kafka.Writer

func initMQ() {
	url := os.Getenv("KAFKA_BROKER_URL")

	if url != "" {
		kafkaBrokerURL = url
	}

	// make a writer that produces to topic-A, using the least-bytes distribution
	w = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaBrokerURL},
		Topic:    "user-creating-topic",
		Balancer: &kafka.LeastBytes{},
	})

	if w == nil {
		log.Println("ERROR    connecting to MQS 'user-creating-topic': FAILED")
	} else {
		log.Println("INFO    connecting to MQS 'user-creating-topic': OK")
	}
}

func closeMQ() {
	if w != nil {
		w.Close()
	}
}

func publishUserCreatingMsg(userInfo *UserInfo) {
	value, err := json.Marshal(userInfo)

	if err != nil {
		return
	}

	err = w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(strconv.Itoa(userInfo.ID)),
			Value: value,
		},
	)

	if err != nil {
		log.Println("ERROR    writting message to 'user-creating-topic': ", err)
	} else {
		log.Println("INFO    writting message to 'user-creating-topic': OK")
	}
}
