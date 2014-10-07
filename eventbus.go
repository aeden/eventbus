package eventbus

import (
	"encoding/json"
	"log"
)

type Event struct {
	Name    string      `json:"name"`
	Data    interface{} `json:"data"`
	Context interface{} `json:"context"`
}

func NewEvent() *Event {
	return &Event{}
}

func Notify(event *Event) {
	// we have the event now, do something useful
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event JSON: %s", err)
		return
	}
	WebsocketHub.Send(eventJSON)
}
