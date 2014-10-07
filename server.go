package eventbus

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type EventBusRequestHandler struct {
	servicesConfig *ServicesConfig
	eventStore     EventStore
}

func NewEventBusRequestHandler(servicesConfig *ServicesConfig, eventStore EventStore) http.Handler {
	return &EventBusRequestHandler{
		servicesConfig: servicesConfig,
		eventStore:     eventStore,
	}
}

func (handler *EventBusRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		authContext := handler.authorizeEventPostClient(w, r)
		log.Printf("Received request from %s", r.RemoteAddr)
		log.Printf("Authorization context: %v", authContext)
		event := NewEvent()
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&event)
		if err != nil {
			log.Printf("Parser error: %s", err)
			http.Error(w, fmt.Sprintf("Parser error: %s", err), 500)
		} else {
			// If client access token isn't present, and auth context is nil, then 401
			clientAccessToken := event.Context["identifier"]
			if clientAccessToken == "" && authContext == nil {
				http.Error(w, "Authorization required", http.StatusUnauthorized)
				return
			}

			// The event should be persisted here
			err := handler.eventStore.WriteEvent(event)
			if err != nil {
				http.Error(w, "Failed to write event", http.StatusInternalServerError)
			}

			// If the event was successfully persisted, return OK
			w.WriteHeader(http.StatusOK)

			// If client access token is present, then send to client
			if clientAccessToken != "" {
				NotifyClient(clientAccessToken, event)
			}

			// Broadcast the event to services
			NotifyServices(event)
		}

	} else if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	} else if r.Method == "GET" {
		json, err := json.Marshal(handler.eventStore.(*InMemoryEventStore).Events)
		if err != nil {
			log.Printf("Error marshaling events JSON: %s", err)
			return
		}
		w.Write(json)

	} else {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func (handler *EventBusRequestHandler) authorizeEventPostClient(w http.ResponseWriter, r *http.Request) (authContext interface{}) {
	authorization := r.Header.Get("Authorization")
	if authorization != "" {
		for _, serviceConfig := range handler.servicesConfig.Services {
			if serviceConfig["token"] == authorization {
				authContext = serviceConfig
				return
			}
		}
	}
	return
}

// Start a file server for serving HTML, CSS and JS files
func StartFileServer(hostAndPort string, corsHostAndPort string) {
	log.Printf("Starting HTTP server on %s", hostAndPort)

	mux := http.NewServeMux()
	mux.Handle("/", NewCorsHandler(corsHostAndPort, http.FileServer(http.Dir("static"))))
	mux.Handle("/ws", http.HandlerFunc(WebsocketHandler))
	server := &http.Server{
		Addr:    hostAndPort,
		Handler: mux,
	}
	server.ListenAndServe()
}

// Start the event bus server for handling JSON events over HTTP
func StartEventBusServer(hostAndPort string, corsHostAndPort string, servicesConfig *ServicesConfig, eventStore EventStore) {
	mux := http.NewServeMux()
	mux.Handle("/", NewCorsHandler(corsHostAndPort, NewEventBusRequestHandler(servicesConfig, eventStore)))

	log.Printf("Starting EventBus service on %s", hostAndPort)

	server := &http.Server{
		Addr:         hostAndPort,
		Handler:      mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	server.ListenAndServe()
}
