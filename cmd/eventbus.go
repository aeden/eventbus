package main

import (
	"github.com/aeden/eventbus"
	"net"
	"os"
)

var (
	fileServerHost = os.Getenv("HTTP_FILE_SERVER_HOST")
	fileServerPort = os.Getenv("HTTP_FILE_SERVER_PORT")
        eventBusServerHost = os.Getenv("HTTP_EVENTBUS_SERVER_HOST")
        eventBusServerPort = os.Getenv("HTTP_EVENTBUS_SERVER_PORT")
)

func main() {
	fileServerHostAndPort := net.JoinHostPort(fileServerHost, fileServerPort)
	eventBusHostAndPort := net.JoinHostPort(eventBusServerHost, eventBusServerPort)

	eventbus.StartWebsocketHub()

	go eventbus.StartFileServer(fileServerHostAndPort, eventBusHostAndPort)
	eventbus.StartEventBusServer(eventBusHostAndPort, fileServerHostAndPort)

}
