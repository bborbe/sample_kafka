package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
	flag "github.com/bborbe/flagenv"
	"github.com/Shopify/sarama"
	"github.com/bborbe/sample_kafka/avro/avro"
	"github.com/golang/glog"
)

var (
	addr    = flag.String("addr", ":8001", "The address to bind to")
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
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	httpServer := &http.Server{
		Addr: *addr,
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

			producer, err := sarama.NewSyncProducer(strings.Split(*brokers, ","), config)
			if err != nil {
				glog.Warning(err)
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			defer producer.Close()

			demoStruct := avro.DemoSchema{
				IntField:    1,
				DoubleField: 2.3,
				StringField: time.Now().Format(time.RFC3339),
				BoolField:   true,
				BytesField:  []byte{1, 2, 3, 4},
			}
			b := &bytes.Buffer{}
			demoStruct.Serialize(b)
			partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.ByteEncoder(b.Bytes()),
			})
			if err != nil {
				glog.Warning(err)
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(resp, "send message successful to %s with partition %d offset %d", topic, partition, offset)
		}),
	}

	glog.V(0).Infof("Listening for requests on %s", *addr)
	log.Fatal(httpServer.ListenAndServe())
}
