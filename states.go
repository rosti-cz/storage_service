package main

import (
	"encoding/json"
	"fmt"
)

// reportState sends a message about something has changed like a db is created or password changed
func reportState(dbtype, alias string, state State) error {
	body, err := json.Marshal(&state)
	if err != nil {
		return err
	}

	return nc.Publish(fmt.Sprintf(publishTemplate, dbtype, alias), body)
}
