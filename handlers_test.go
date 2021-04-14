package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

// This is integration test and it needs dev VM running
func TestMessageHandlerMySQL(t *testing.T) {
	var err error

	config = Config{
		NATSURL:   "nats://192.168.122.127:4222",
		Databases: "devmysql:mariadb:192.168.122.127:3306:rosti:rosti;devpgsql:pgsql:192.168.122.127:5432:rosti:rosti",
	}

	nc, err = nats.Connect(config.NATSURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer nc.Drain()

	randomName := fmt.Sprintf("test%d", time.Now().Unix())

	msgCreated := nats.Msg{
		Subject: "admin.storages.mariadb.devmysql.events",
		Reply:   "",
		Data:    []byte(`{"event_type": "created", "db_id": 29, "db_name": "` + randomName + `", "username": "` + randomName + `", "password": "test", "extensions": []}`),
		Sub:     &nats.Subscription{Subject: ""},
	}
	msgDeleted := nats.Msg{
		Subject: "admin.storages.mariadb.devmysql.events",
		Reply:   "",
		Data:    []byte(`{"event_type": "deleted", "db_id": 29, "db_name": "` + randomName + `", "username": "` + randomName + `", "password": "test", "extensions": []}`),
		Sub:     &nats.Subscription{Subject: ""},
	}
	msgPasswordChanged := nats.Msg{
		Subject: "admin.storages.mariadb.devmysql.events",
		Reply:   "",
		Data:    []byte(`{"event_type": "password_changed", "db_id": 29, "db_name": "` + randomName + `", "username": "` + randomName + `", "password": "newtest", "extensions": []}`),
		Sub:     &nats.Subscription{Subject: ""},
	}

	assert.Nil(t, _messageHandler(&msgCreated))
	assert.Nil(t, _messageHandler(&msgPasswordChanged))
	assert.Nil(t, _messageHandler(&msgDeleted))
}

func TestMessageHandlerPgSQL(t *testing.T) {
	var err error

	config = Config{
		NATSURL:   "nats://192.168.122.127:4222",
		Databases: "devmysql:mariadb:192.168.122.127:3306:rosti:rosti;devpgsql:pgsql:192.168.122.127:5432:rosti:rosti",
	}

	nc, err = nats.Connect(config.NATSURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer nc.Drain()

	randomName := fmt.Sprintf("test%d", time.Now().Unix())

	msgCreated := nats.Msg{
		Subject: "admin.storages.pgsql.devpgsql.events",
		Reply:   "",
		Data:    []byte(`{"event_type": "created", "db_id": 29, "db_name": "` + randomName + `", "username": "` + randomName + `", "password": "test", "extensions": []}`),
		Sub:     &nats.Subscription{Subject: ""},
	}
	msgDeleted := nats.Msg{
		Subject: "admin.storages.pgsql.devpgsql.events",
		Reply:   "",
		Data:    []byte(`{"event_type": "deleted", "db_id": 29, "db_name": "` + randomName + `", "username": "` + randomName + `", "password": "test", "extensions": []}`),
		Sub:     &nats.Subscription{Subject: ""},
	}
	msgPasswordChanged := nats.Msg{
		Subject: "admin.storages.pgsql.devpgsql.events",
		Reply:   "",
		Data:    []byte(`{"event_type": "password_changed", "db_id": 29, "db_name": "` + randomName + `", "username": "` + randomName + `", "password": "newtest", "extensions": []}`),
		Sub:     &nats.Subscription{Subject: ""},
	}

	assert.Nil(t, _messageHandler(&msgCreated))
	assert.Nil(t, _messageHandler(&msgPasswordChanged))
	assert.Nil(t, _messageHandler(&msgDeleted))
}
