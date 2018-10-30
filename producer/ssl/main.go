package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strings"
	"github.com/Shopify/sarama"
	"github.com/golang/glog"
)

func main() {

	tlsConfig, err := NewTLSConfig(
		"/Users/bborbe/Documents/workspaces/go/src/github.com/bborbe/kafka/secrets/kafka.client.cer.pem",
		"/Users/bborbe/Documents/workspaces/go/src/github.com/bborbe/kafka/secrets/kafka.client.key.pem",
		"/Users/bborbe/Documents/workspaces/go/src/github.com/bborbe/kafka/secrets/kafka.server.cer.pem",
	)
	if err != nil {
		glog.Fatal(err)
	}

	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	config.Net.TLS.Config = tlsConfig
	config.Net.TLS.Enable = true

	if err := config.Validate(); err != nil {
		glog.Fatal(err)
	}

	client, err := sarama.NewClient(strings.Split("kafka:9093", ","), config)
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

func NewTLSConfig(clientCertFile, clientKeyFile, caCertFile string) (*tls.Config, error) {
	tlsConfig := tls.Config{}

	// Load client cert
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return &tlsConfig, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = caCertPool

	tlsConfig.BuildNameToCertificate()
	return &tlsConfig, err
}
