package eventbus

import (
	"encoding/json"
	"log"
)

func routeEvent(event *Event) {
	clientAccessToken := event.Context["identifier"]

	// If client access token is present, then send to client
	if clientAccessToken != "" {
		notifyClient(clientAccessToken, event)
	}

	// Broadcast the event to services
	notifyServices(event)
}

func notify(event *Event) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event JSON: %s", err)
		return
	}
	websocketHub.send(eventJSON)
}

func notifyClient(clientAccessToken string, event *Event) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event JSON: %s", err)
		return
	}
	websocketHub.sendToClient(clientAccessToken, eventJSON)
}

func notifyServices(event *Event) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event JSON: %s", err)
		return
	}
	websocketHub.sendToServices(eventJSON)
}
