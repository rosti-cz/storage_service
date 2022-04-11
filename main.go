package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/nats.go"
)

const subscribeTemplate = "admin.storages.%s.%s.events" // storage_type and alias
const publishTemplate = "admin.storages.%s.%s.states"   // storage_type and alias

var config Config
var nc *nats.Conn
var metrics Metrics = Metrics{}

// We have to change name of this function so tests are working without being affected by this.
func _init() {
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	if config.NATSToken != "" {
		nc, err = nats.Connect(config.NATSURL, nats.Token(config.NATSToken))
	} else {
		nc, err = nats.Connect(config.NATSURL)
	}

	metrics.Service = config.MetricsIdent

	if err != nil {
		log.Fatalln(err)
	}
}

// sends metrics to NATS
func sentMetrics(nc *nats.Conn, subject string) {
	metrics.Service = config.MetricsIdent

	receiverMetrics := ReceiverMetric{
		Service: config.MetricsIdent,
		Lines: []string{
			"# HELP storage_service_messages Number of received messages in the current session",
			"# TYPE storage_service_messages counter",
			fmt.Sprintf("storage_service_messages{service=\"%s\"} %d", config.MetricsIdent, metrics.Messages),
		},
	}

	data, err := json.Marshal(receiverMetrics)
	if err != nil {
		log.Println("ERROR: metrics sending:", err)
	}

	err = nc.Publish(subject, data)
	if err != nil {
		log.Println("ERROR: metrics sending:", err)
	}
}

func main() {
	_init()

	defer func() {
		err := nc.Drain()
		if err != nil {
			log.Println("Drain error:", err.Error())
		}
	}()

	// Share metrics with the ecosystem
	go func() {
		for {
			sentMetrics(nc, config.NATSMetricsSubject)
			time.Sleep(15 * time.Second)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for _, database := range strings.Split(config.Databases, ";") {
		databaseParts := strings.Split(database, ":")
		subject := fmt.Sprintf(subscribeTemplate, databaseParts[1], databaseParts[0])

		log.Println("Listening for " + subject)
		_, err := nc.Subscribe(subject, messageHandler)
		if err != nil {
			log.Println("Subscribe error:", err)
		}
	}

	// runtime.Goexit()

	<-sigs
	err := nc.Drain()
	if err != nil {
		log.Println(err)
	}
	os.Exit(0)
}
