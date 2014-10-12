package eventbus

import (
	"encoding/json"
	"fmt"
	"github.com/aeden/eventbus/middleware"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// configuration options

type option func(*Server)

// The EventBus server.
type Server struct {
	httpServer      *http.Server
	corsHostAndPort string
	eventStore      EventStore
	servicesConfig  *ServicesConfig
}

// Configure a new server that is ready to be started.
func NewServer(opts ...option) *Server {
	server := &Server{
		httpServer: &http.Server{
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
		},
		eventStore:     NewNullEventStore(),
		servicesConfig: &ServicesConfig{},
	}

	for _, opt := range opts {
		opt(server)
	}

	mux := http.NewServeMux()
	mux.Handle("/", middleware.NewCorsHandler(server.corsHostAndPort, newEventBusRequestHandler(server.servicesConfig, server.eventStore)))
	mux.Handle("/ws", newWebSocketHandler(server.corsHostAndPort))
	server.httpServer.Handler = mux

	return server
}

/*
Start the event bus server for handling JSON events over HTTP.

This function starts a handler on the root that is used for POST
requests to construct new events. It also starts a WebSocket
handler on <code>/ws</code> that is used for broadcasting events to
the client or service.
*/
func (server *Server) Start() {
	startWebsocketHub()
	log.Printf("Starting EventBus server %s", server.httpServer.Addr)
	server.httpServer.ListenAndServe()
}

// Configure the host and port of the EventBus server.
func HostAndPort(v string) option {
	return func(server *Server) {
		server.httpServer.Addr = v
	}
}

// Configure the CORS host and port for the EventBus server. This is
// the host and port where JavaScript client calls are coming from.
func CorsHostAndPort(v string) option {
	return func(server *Server) {
		server.corsHostAndPort = v
	}
}

// Configure the services that will attach to the EventBus server.
func Services(in io.Reader) option {
	return func(server *Server) {
		file, e := ioutil.ReadAll(in)
		if e != nil {
			log.Printf("Error reading services config: %s", e)
		}

		servicesConfig := ServicesConfig{}
		json.Unmarshal(file, &servicesConfig.Services)
		server.servicesConfig = &servicesConfig
	}
}

// internal

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
