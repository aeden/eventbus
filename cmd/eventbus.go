package main

import (
	"github.com/aeden/eventbus"
	"github.com/aeden/eventbus/fileserver"
	"log"
	"net"
	"os"
)

var (
	fileServerHost     = os.Getenv("HTTP_FILE_SERVER_HOST")
	fileServerPort     = os.Getenv("HTTP_FILE_SERVER_PORT")
	eventBusServerHost = os.Getenv("HTTP_EVENTBUS_SERVER_HOST")
	eventBusServerPort = os.Getenv("HTTP_EVENTBUS_SERVER_PORT")
)

func main() {
	servicesFile := "services.json"
	file, err := os.Open(servicesFile)
	if err != nil {
		log.Printf("Failed to open %s", servicesFile)
		os.Exit(1)
	}

	fileServerHostAndPort := net.JoinHostPort(fileServerHost, fileServerPort)
	eventBusHostAndPort := net.JoinHostPort(eventBusServerHost, eventBusServerPort)

	go fileserver.StartFileServer(fileServerHostAndPort, eventBusHostAndPort)

	server, err := eventbus.NewServer(
		eventbus.HostAndPort(eventBusHostAndPort),
		eventbus.CorsHostAndPort(fileServerHostAndPort),
		eventbus.Services(file),
	)
	if err != nil {
		log.Printf("Failed to start server: %s", err)
		os.Exit(1)
	}
	server.Start()

}
