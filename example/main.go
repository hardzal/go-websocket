package example

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// Thread-safe list of connected clients
var clients = make(map[*Client]bool)
var clientsLock sync.Mutex

var broadcast = make(chan []byte) // Channel to broadcast messages

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (NOT for production)
	},
}

// WebSocket handler
func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte)}

	// Register client
	clientsLock.Lock()
	clients[client] = true
	clientsLock.Unlock()

	log.Println("New client connected")

	// Start goroutine to send messages to this client
	go handleSending(client)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		log.Printf("Received message: %s", msg)
		broadcast <- msg
	}

	// Cleanup
	clientsLock.Lock()
	delete(clients, client)
	clientsLock.Unlock()
	conn.Close()
	log.Println("Client disconnected")
}

// Send messages from send channel to the client connection
func handleSending(client *Client) {
	for msg := range client.send {
		err := client.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}

func handleBroadcast() {
	for {
		msg := <-broadcast
		clientsLock.Lock()
		for client := range clients {
			select {
			case client.send <- msg:
			default:
				log.Println("Client send channel full, skipping")
			}
		}
		clientsLock.Unlock()
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	go handleBroadcast()

	log.Println("WebSocket server running on http://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
