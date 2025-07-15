package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}

	defer func(){
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

func (wsh webSocketHandler) NewRun(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
    conn, err := wsh.upgrader.Upgrade(w, r, nil)
    if err != nil {
       fmt.Println("Error upgrading:", err)
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

func main() {
	websocketHandler := webSocketHandler{
		upgrader: websocket.Upgrader{},
	}

	http.Handle("/ws", websocketHandler)
	http.HandleFunc("/ws/new", websocketHandler.NewRun)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}