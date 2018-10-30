package main

import (
	"fmt"
	"strings"
	"github.com/Shopify/sarama"
	"github.com/golang/glog"
)

func main() {

	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	if err := config.Validate(); err != nil {
		glog.Fatal(err)
	}

	client, err := sarama.NewClient(strings.Split("kafka:9092", ","), config)
	if err != nil {
		glog.Fatal(err)
	}
	defer client.Close()

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		glog.Fatal(err)
	}
	defer producer.Close()

	topic := "test"
	partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder("test"),
	})
	if err != nil {
		glog.Fatal(err)
	}
	fmt.Printf("send message successful to %s with partition %d offset %d", topic, partition, offset)
}
