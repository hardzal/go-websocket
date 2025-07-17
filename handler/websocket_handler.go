package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	Upgrader websocket.Upgrader
}

func NewWebSocketHandler() WebSocketHandler {
	return WebSocketHandler{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // allow all incoming request
			},
		},
	}
}

// websocket innit
func (wsh WebSocketHandler) upgradeConn(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := wsh.Upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Upgrade failed: %v", err)
		return nil, err
	}

	return conn, nil
}

// challnge 1
// just read message from client
func (wsh WebSocketHandler) ReceivedMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := wsh.upgradeConn(w, r)
	if err != nil {
		return
	}
	defer conn.Close()
	// Listen for incoming messages
	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Received: %s\\n", message)
		// Echo the message back to the client
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}

// challenge 2
// log message from client
func (wsh WebSocketHandler) LogMessageHandler(w http.ResponseWriter, r *http.Request) {
	// membuat koneksi dengan websocket
	conn, err := wsh.upgradeConn(w, r)

	if err != nil {
		return
	}
	// close in the end
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("error when reading message: ", err)
			break
		}

		log.Printf("Message: %s\n", message)
		log.Printf("Time: %s, IP: %s", time.Now().Format(time.DateTime), r.RemoteAddr)

		if err = conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error when writing message", err)
			break
		}
	}
}

// challenge 3
// push notification handler
func (wsh WebSocketHandler) PushNotificationHandler(w http.ResponseWriter, r *http.Request) {
	c, err := wsh.upgradeConn(w, r)
	if err != nil {
		return
	}

	defer func() {
		log.Println("closing connection")
		c.Close()
	}()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error %s when reading message from client", err)
			return
		}

		if mt == websocket.BinaryMessage {
			err = c.WriteMessage(websocket.TextMessage, []byte("server doesnt support binary messages"))
			if err != nil {
				log.Printf("Error %s when sending message to client.", err)
			}
			return
		}

		log.Printf("Receive message %s", string(message))
		if strings.TrimSpace(string(message)) != "start" {
			err = c.WriteMessage(websocket.TextMessage, []byte("You did not say the magic word!"))

			if err != nil {
				log.Printf("Error %s when sending message to client", err)
				return
			}

			continue
		}

		log.Println("start responding to client....")
		i := 1

		for {
			response := fmt.Sprintf("Notification %d", i)
			err = c.WriteMessage(websocket.TextMessage, []byte(response))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
				return
			}

			i = i + 1
			time.Sleep(2 * time.Second)
		}
	}
}
