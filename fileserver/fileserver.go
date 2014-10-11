package fileserver

import (
	"github.com/aeden/eventbus/middleware"
	"log"
	"net/http"
)

// Start a file server for serving HTML, CSS and JS files
func StartFileServer(hostAndPort string, corsHostAndPort string) {
	log.Printf("Starting HTTP server on %s", hostAndPort)

	mux := http.NewServeMux()
	mux.Handle("/", middleware.NewCorsHandler(corsHostAndPort, http.FileServer(http.Dir("static"))))

	server := &http.Server{
		Addr:    hostAndPort,
		Handler: mux,
	}
	server.ListenAndServe()
}
