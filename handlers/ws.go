package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"chatapp/internal/chat"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(hub *chat.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {

		username := c.Query("username")
		roomName := c.Query("room")

		if username == "" {
			username = "Anonymous"
		}
		if roomName == "" {
			roomName = "general"
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}

		client := &chat.Client{
			Conn:     conn,
			Username: username,
			Room:     roomName,
		}

		hub.JoinRoom(roomName, client)
		hub.Broadcast(roomName, []byte(username+" joined "+roomName))

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				hub.LeaveRoom(client)
				hub.Broadcast(roomName, []byte(username+" left "+roomName))
				break
			}

			formatted := []byte("[" + roomName + "] " + username + ": " + string(msg))
			hub.Broadcast(roomName, formatted)
		}
	}
}
