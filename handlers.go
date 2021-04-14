package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/rosti-cz/storage_service/mysql"
	"github.com/rosti-cz/storage_service/pgsql"
)

func report(dbtype, alias string, stateMessage string, message Message, isError bool) {
	err := reportState(dbtype, alias, State{
		DBID:    message.DBID,
		DBName:  message.DBName,
		Error:   isError,
		Message: stateMessage,
	})
	if err != nil {
		log.Println("ERROR: report state:", err.Error())
	}
}

func messageHandler(m *nats.Msg) {
	message := Message{}
	err := json.Unmarshal(m.Data, &message)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("Received a message: %v\n", message)

	var backend Backend

	dbtype := strings.Split(m.Subject, ".")[2]
	alias := strings.Split(m.Subject, ".")[3]

	databaseLine := config.DatabasesMap()[alias+":"+dbtype]

	// MariaDB/MySQL backed setup
	if dbtype == "mysql" || dbtype == "mariadb" {
		port, err := strconv.Atoi(databaseLine.Port)
		if err != nil {
			log.Println("Port issue in config:", err)
		}
		backend = &mysql.MySQLBackend{
			Username: databaseLine.Username,
			Password: databaseLine.Password,
			Hostname: databaseLine.Hostname,
			Port:     port,
		}
	} else if dbtype == "pgsql" { // PostgreSQL backend setup
		port, err := strconv.Atoi(databaseLine.Port)
		if err != nil {
			log.Println("Port issue in config:", err)
		}
		backend = &pgsql.PGSQLBackend{
			Username: databaseLine.Username,
			Password: databaseLine.Password,
			Hostname: databaseLine.Hostname,
			Port:     port,
		}
	} else {
		log.Println("ERROR: database backend not found")
		report(dbtype, alias, "wrong backend", message, true)
		return
	}

	// Message processing
	if message.EventType == "created" {
		err = backend.CreateUser(message.Username, message.Password, message.DBName)
		if err != nil {
			log.Println("ERROR: backend problem:", err.Error())
			report(dbtype, alias, "backend problem", message, true)
			return
		}
		err = backend.CreateDatabase(message.DBName, message.Username, message.Extensions)
		if err != nil {
			log.Println("ERROR: backend problem:", err.Error())
			report(dbtype, alias, "backend problem", message, true)
			return
		}
		report(dbtype, alias, "created", message, false)
	}

	if message.EventType == "password_changed" {
		err = backend.ChangePassword(message.Username, message.Password)
		if err != nil {
			log.Println("ERROR: backend problem:", err.Error())
			report(dbtype, alias, "backend problem", message, true)
			return
		}
		report(dbtype, alias, "password changed", message, false)
	}

	if message.EventType == "deleted" {
		err = backend.DropDatabase(message.DBName)
		if err != nil {
			log.Println("ERROR: backend problem:", err.Error())
			report(dbtype, alias, "backend problem", message, true)
			return
		}
		err = backend.DropUser(message.Username)
		if err != nil {
			log.Println("ERROR: backend problem:", err.Error())
			report(dbtype, alias, "backend problem", message, true)
			return
		}
		report(dbtype, alias, "deleted", message, false)
	}
}
