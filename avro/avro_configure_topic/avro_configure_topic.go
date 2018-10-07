package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
)

var (
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
	admin, err := sarama.NewClusterAdmin(strings.Split(*brokers, ","), config)
	if err != nil {
		glog.Exit(err)
	}

	err = admin.AlterConfig(sarama.TopicResource, topic, map[string]*string{
		"retention.ms":    pointer("-1"),
		"retention.bytes": pointer("-1"),
		"cleanup.policy":  pointer("compact"),
	}, false)
	if err != nil {
		glog.Exit(err)
	}
	fmt.Println("configured")

	entries, err := admin.DescribeConfig(sarama.ConfigResource{
		Type:        sarama.TopicResource,
		Name:        topic,
		ConfigNames: []string{"retention.ms", "retention.bytes", "cleanup.policy"},
	})
	if err != nil {
		glog.Exit(err)
	}
	for _, entry := range entries {
		fmt.Printf("%+v\n", entry)
	}
}

func pointer(value string) *string {
	return &value
}
