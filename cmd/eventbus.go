package main

import (
	"encoding/json"
	"fmt"
	"github.com/aeden/eventbus"
	"log"
	"net/http"
	"time"
)

// HTTP handler for EventBus HTTP requests
func eventBusRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Printf("Received request from %s", r.RemoteAddr)
		event := eventbus.NewEvent()

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&event)
		if err != nil {
			log.Printf("Parser error: %s", err)
			http.Error(w, fmt.Sprintf("Parser error: %s", err), 500)
		} else {
			w.WriteHeader(http.StatusOK)
			eventbus.Notify(event)
		}

	} else if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

// Start a file server for serving HTML, CSS and JS files
func startFileServer(hostAndPort string, corsHostAndPort string) {
	log.Printf("Starting HTTP server on %s", hostAndPort)

	mux := http.NewServeMux()
	mux.Handle("/", corsServer(corsHostAndPort, http.FileServer(http.Dir("static"))))
	mux.Handle("/ws", http.HandlerFunc(eventbus.WebsocketHandler))
	server := &http.Server{
		Addr:    hostAndPort,
		Handler: mux,
	}
	server.ListenAndServe()
}

// Start the event bus server for handling JSON events over HTTP
func startEventBusServer(hostAndPort string, corsHostAndPort string) {
	mux := http.NewServeMux()
	mux.Handle("/", corsServer(corsHostAndPort, http.HandlerFunc(eventBusRequestHandler)))

	log.Printf("Starting EventBus service on %s", hostAndPort)

	server := &http.Server{
		Addr:         hostAndPort,
		Handler:      mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	server.ListenAndServe()
}

// CORS handling
type corsHandler struct {
	corsHostAndPort string
	delegate        http.Handler
}

func corsServer(corsHostAndPort string, handler http.Handler) http.Handler {
	return &corsHandler{
		corsHostAndPort: corsHostAndPort,
		delegate:        handler,
	}
}

func (server *corsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("http://%s", server.corsHostAndPort))
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	server.delegate.ServeHTTP(w, r)
}

// Main function
func main() {
	staticHostAndPort := "localhost:3000"
	eventBusHostAndPort := "localhost:3001"

	eventbus.StartWebsocketHub()

	go startFileServer(staticHostAndPort, eventBusHostAndPort)
	startEventBusServer(eventBusHostAndPort, staticHostAndPort)

}
