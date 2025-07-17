package routes

import (
	"net/http"

	"github.com/hardzal/go-websocket/handler"
)

func RegisterWebSocketRoutes(mux *http.ServeMux) {
	wsHandler := handler.NewWebSocketHandler()

	mux.HandleFunc("/ws", wsHandler.ReceivedMessageHandler)
	mux.HandleFunc("/ws/log", wsHandler.LogMessageHandler)
	mux.HandleFunc("/ws/push", wsHandler.PushNotificationHandler)
}
