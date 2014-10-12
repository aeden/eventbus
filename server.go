package eventbus

import (
	"encoding/json"
	"fmt"
	"github.com/aeden/eventbus/middleware"
	"log"
	"net/http"
	"time"
)

type eventBusRequestHandler struct {
	servicesConfig *ServicesConfig
	eventStore     EventStore
}

func newEventBusRequestHandler(servicesConfig *ServicesConfig, eventStore EventStore) http.Handler {
	return &eventBusRequestHandler{
		servicesConfig: servicesConfig,
		eventStore:     eventStore,
	}
}

func (handler *eventBusRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request from %s", r.RemoteAddr)

	if r.Method == "POST" {
		handler.handlePost(w, r)
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

func (handler *eventBusRequestHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	authContext := handler.prepareAuthContext(w, r)
	log.Printf("Authorization context: %v", authContext)

	event, err := handler.decodeEvent(r)
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

		// The event is persisted here
		err := handler.eventStore.WriteEvent(event)
		if err != nil {
			http.Error(w, "Failed to write event", http.StatusInternalServerError)
			return
		}

		// If the event was successfully persisted, return OK
		w.WriteHeader(http.StatusOK)

		// Route event
		go routeEvent(event)
	}
}

func (handler *eventBusRequestHandler) decodeEvent(r *http.Request) (event *Event, err error) {
	event = NewEvent()
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&event)
	return
}

func (handler *eventBusRequestHandler) prepareAuthContext(w http.ResponseWriter, r *http.Request) (authContext interface{}) {
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

// Start the event bus server for handling JSON events over HTTP.
//
// This function starts a handler on the root that is used for POST requests to construct new events. It also 
// starts a WebSocket handler on `/ws` that is used for broadcasting events to the client or service.
func StartEventBusServer(hostAndPort string, corsHostAndPort string, servicesConfig *ServicesConfig, eventStore EventStore) {
	mux := http.NewServeMux()
	mux.Handle("/", middleware.NewCorsHandler(corsHostAndPort, newEventBusRequestHandler(servicesConfig, eventStore)))
	mux.Handle("/ws", newWebSocketHandler(corsHostAndPort))

	log.Printf("Starting EventBus service on %s", hostAndPort)

	server := &http.Server{
		Addr:         hostAndPort,
		Handler:      mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	server.ListenAndServe()
}
