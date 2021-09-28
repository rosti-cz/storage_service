package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/nats.go"
)

const subscribeTemplate = "admin.storages.%s.%s.events" // storage_type and alias
const publishTemplate = "admin.storages.%s.%s.states"   // storage_type and alias

var config Config
var nc *nats.Conn

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

	if err != nil {
		log.Fatalln(err)
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		err := nc.Drain()
		if err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)
	}()

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

	fmt.Print("Press 'Enter' to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
