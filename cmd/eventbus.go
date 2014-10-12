package main

import (
	"encoding/json"
	"github.com/aeden/eventbus"
	"github.com/aeden/eventbus/fileserver"
	"io/ioutil"
	"log"
	"net"
	"os"
)

var (
	fileServerHost     = os.Getenv("HTTP_FILE_SERVER_HOST")
	fileServerPort     = os.Getenv("HTTP_FILE_SERVER_PORT")
	eventBusServerHost = os.Getenv("HTTP_EVENTBUS_SERVER_HOST")
	eventBusServerPort = os.Getenv("HTTP_EVENTBUS_SERVER_PORT")
	eventStore         = eventbus.NewNullEventStore()
)

func loadServicesConfig() *eventbus.ServicesConfig {
	file, e := ioutil.ReadFile("services.json")
	if e != nil {
		log.Printf("Error reading services config: %s", e)
	}

	servicesConfig := eventbus.ServicesConfig{}
	json.Unmarshal(file, &servicesConfig.Services)
	return &servicesConfig
}

func main() {
	servicesConfig := loadServicesConfig()

	fileServerHostAndPort := net.JoinHostPort(fileServerHost, fileServerPort)
	eventBusHostAndPort := net.JoinHostPort(eventBusServerHost, eventBusServerPort)

	eventbus.StartWebsocketHub()

	go fileserver.StartFileServer(fileServerHostAndPort, eventBusHostAndPort)
	eventbus.StartEventBusServer(eventBusHostAndPort, fileServerHostAndPort, servicesConfig, eventStore)

}
