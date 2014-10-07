package eventbus

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// HTTP handler for creating events in the event bus.
func eventBusRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Printf("Received request from %s", r.RemoteAddr)
		event := NewEvent()

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&event)
		if err != nil {
			log.Printf("Parser error: %s", err)
			http.Error(w, fmt.Sprintf("Parser error: %s", err), 500)
		} else {
			w.WriteHeader(http.StatusOK)
			Notify(event)
		}

	} else if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
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
func StartEventBusServer(hostAndPort string, corsHostAndPort string) {
	mux := http.NewServeMux()
	mux.Handle("/", NewCorsHandler(corsHostAndPort, http.HandlerFunc(eventBusRequestHandler)))

	log.Printf("Starting EventBus service on %s", hostAndPort)

	server := &http.Server{
		Addr:         hostAndPort,
		Handler:      mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	server.ListenAndServe()
}
