package eventbus

import (
	"encoding/json"
	"fmt"
	"github.com/aeden/eventbus/middleware"
	nsq "github.com/bitly/go-nsq"
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
	servicesConfig  []ServiceConfig
	nsqProducer     *nsq.Producer
	nsqConsumer     *nsq.Consumer
	nsqTopic        string
	nsqChannel      string
}

// Configure a new server that is ready to be started.
func NewServer(opts ...option) (server *Server, err error) {
	server = &Server{
		httpServer: &http.Server{
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
		},
		servicesConfig: []ServiceConfig{},
		nsqTopic:       "events",
		nsqChannel:     "all",
	}

	for _, opt := range opts {
		opt(server)
	}

	// NSQ producer for sending messages
	producerConfig := nsq.NewConfig()
	err = producerConfig.Validate()
	if err != nil {
		log.Printf("Producer config is not valid: %s", err)
		return
	}
	server.nsqProducer, err = nsq.NewProducer("localhost:4150", producerConfig)
	if err != nil {
		log.Printf("Error connecting to NSQ: %s", err)
		return
	}

	// NSQ consumer for receiving messages
	consumerConfig := nsq.NewConfig()
	err = consumerConfig.Validate()
	if err != nil {
		log.Printf("Consumer config is not valid: %s", err)
		return
	}
	server.nsqConsumer, err = nsq.NewConsumer(server.nsqTopic, server.nsqChannel, consumerConfig)
	if err != nil {
		log.Printf("Error creating consumer: %s", err)
		return
	}
	server.nsqConsumer.AddHandler(nsq.HandlerFunc(nsqHandler))

	// HTTP server for handling WebSocket connections
	mux := http.NewServeMux()
	mux.Handle("/", middleware.NewCorsHandler(server.corsHostAndPort,
		newEventBusRequestHandler(server.servicesConfig, server)))
	mux.Handle("/ws", newWebSocketHandler(server.corsHostAndPort))
	server.httpServer.Handler = mux

	return server, nil
}

/*
Start the event bus server for handling JSON events over HTTP.

This function starts a handler on the root that is used for POST
requests to construct new events. It also starts a WebSocket
handler on /ws that is used for broadcasting events to
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

		servicesConfig := &[]ServiceConfig{}
		json.Unmarshal(file, servicesConfig)
		server.servicesConfig = *servicesConfig
	}
}

// internal

type eventBusRequestHandler struct {
	servicesConfig []ServiceConfig
	server         *Server
}

func newEventBusRequestHandler(servicesConfig []ServiceConfig, server *Server) http.Handler {
	return &eventBusRequestHandler{
		servicesConfig: servicesConfig,
		server:         server,
	}
}

func (handler *eventBusRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request from %s", r.RemoteAddr)

	if r.Method == "POST" {
		handler.handlePost(w, r)
	} else if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
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

		// The event is published to NSQ here
		err := handler.publishEvent(w, event)
		if err != nil {
			http.Error(w, fmt.Sprintf("Internal error: %s", err), http.StatusInternalServerError)
			return
		}

		// If the event was successfully published, return OK
		w.WriteHeader(http.StatusOK)
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
		for _, serviceConfig := range handler.servicesConfig {
			if serviceConfig.Token == authorization {
				authContext = serviceConfig
				return
			}
		}
	}
	return
}

func (handler *eventBusRequestHandler) publishEvent(w http.ResponseWriter, event *Event) (err error) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event JSON: %s", err)
		return
	}

	log.Printf("Publishing event to %s", handler.server.nsqTopic)
	handler.server.nsqProducer.Publish(handler.server.nsqTopic, eventJSON)

	return
}

// routing and sending events

func nsqHandler(message *nsq.Message) (err error) {
	log.Printf("Message: %s", message.Body)
	return
}

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
