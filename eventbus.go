package eventbus

import (
	"encoding/json"
	"log"
)

func Notify(event *Event) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event JSON: %s", err)
		return
	}
	WebsocketHub.Send(eventJSON)
}

func NotifyClient(clientAccessToken string, event *Event) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event JSON: %s", err)
		return
	}
	WebsocketHub.SendToClient(clientAccessToken, eventJSON)
}

func NotifyServices(event *Event) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event JSON: %s", err)
		return
	}
	WebsocketHub.SendToServices(eventJSON)
}
