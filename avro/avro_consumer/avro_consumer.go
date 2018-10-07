package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/sample_kafka/avro/avro"
	"github.com/golang/glog"
)

var (
	addr    = flag.String("addr", ":8002", "The address to bind to")
	brokers = flag.String("brokers", "kafka:9092", "The Kafka brokers to connect to, as a comma separated list")
	verbose = flag.Bool("verbose", false, "Turn on Sarama logging")
)

const topic = "sample_avro"

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Set("logtostderr", "true")
	flag.Parse()

	glog.V(0).Infof("addr %s", *addr)
	glog.V(0).Infof("brokers %s", *brokers)
	glog.V(0).Infof("verbose %v", *verbose)

	if *verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	if *brokers == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0

	httpServer := &http.Server{
		Addr: *addr,
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			client, err := sarama.NewClient(strings.Split(*brokers, ","), config)
			if err != nil {
				glog.Warning(err)
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			defer client.Close()

			consumer, err := sarama.NewConsumerFromClient(client)
			if err != nil {
				glog.Warning(err)
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			defer consumer.Close()

			partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
			if err != nil {
				glog.Warning(err)
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			defer partitionConsumer.Close()

			ctx, cancelFunc := context.WithTimeout(req.Context(), time.Second)
			defer cancelFunc()

			for {
				select {
				case err := <-partitionConsumer.Errors():
					glog.Warning(err)
					http.Error(resp, err.Error(), http.StatusInternalServerError)
					return
				case msg := <-partitionConsumer.Messages():
					b := bytes.NewBuffer(msg.Value)
					demoSchema, err := avro.DeserializeDemoSchema(b)
					if err != nil {
						glog.Warning(err)
						http.Error(resp, err.Error(), http.StatusInternalServerError)
						return
					}
					fmt.Fprintf(resp, "found schema: %+v\n", demoSchema)
				case <-ctx.Done():
					return
				}
			}
		}),
	}

	glog.V(0).Infof("Listening for requests on %s", *addr)
	log.Fatal(httpServer.ListenAndServe())
}
