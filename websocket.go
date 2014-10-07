package eventbus

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
)

type wsConnection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

type wsHub struct {
	// Registered connections.
	connections map[*wsConnection]*wsConnectionState

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *wsConnection

	// Unregister requests from connections.
	unregister chan *wsConnection
}

type wsCommand struct {
	Action string `json:"action"`
}

type wsIdentifyCommandResponse struct {
	Action string `json:"action"`
	Token  string `json:"token"`
}

type wsConnectionState struct {
	Token string
}

var WebsocketHub = wsHub{
	broadcast:   make(chan []byte),
	register:    make(chan *wsConnection),
	unregister:  make(chan *wsConnection),
	connections: make(map[*wsConnection]*wsConnectionState),
}

func (h *wsHub) Send(message []byte) {
	h.broadcast <- message
}

func (h *wsHub) run() {
	log.Printf("Run websocket hub")
	for {
		select {
		case c := <-h.register:
			log.Printf("Register connection %s", c.ws.RemoteAddr())
			h.connections[c] = &wsConnectionState{}
		case c := <-h.unregister:
			log.Printf("Unregister connection %s", c.ws.RemoteAddr())
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			log.Printf("Broadcast message: %s", m)
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					delete(h.connections, c)
					close(c.send)
				}
			}
		}
	}
}

func (c *wsConnection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from web socket: %s", err)
			break
		}
		log.Printf("Message received: %s", message)

		// Right now this only handles a single command: identify.

		command := &wsCommand{}
		json.Unmarshal(message, command)
		log.Printf("Command received: %s", command.Action)

		connectionState := WebsocketHub.connections[c]
		connectionState.Token = randSeq(16)
		log.Printf("Connection state: %s", connectionState)

		response := &wsIdentifyCommandResponse{
			Action: command.Action,
			Token:  connectionState.Token,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshaling response JSON: %s", err)
			continue
		}
		c.ws.WriteMessage(websocket.TextMessage, responseJSON)
	}
	c.ws.Close()
}

func (c *wsConnection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %s", err)
		return
	}
	log.Printf("Establishing websocket connection")
	c := &wsConnection{send: make(chan []byte, 256), ws: ws}
	WebsocketHub.register <- c
	defer func() { WebsocketHub.unregister <- c }()
	go c.writer()
	c.reader()
}

func StartWebsocketHub() {
	go WebsocketHub.run()
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
