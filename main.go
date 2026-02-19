package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins (dev only)
	},
}

type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) run() {
	for {
		select {

		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()

		case conn := <-h.unregister:
			h.mu.Lock()
			delete(h.clients, conn)
			conn.Close()
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					client.Close()
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func main() {
	r := gin.Default()
	hub := newHub()
	go hub.run()

	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, htmlPage)
	})

	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}

		hub.register <- conn

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				hub.unregister <- conn
				break
			}
			hub.broadcast <- msg
		}
	})

	r.Run(":8080")
}

const htmlPage = `
<!DOCTYPE html>
<html>
<head>
	<title>Go Chat</title>
</head>
<body>
	<h2>Simple Broadcast Chat</h2>

	<div id="joinSection">
		<input id="username" type="text" placeholder="Enter username" />
		<button onclick="connect()">Join</button>
	</div>

	<p id="status"></p>

	<br>

	<input id="msg" type="text" placeholder="Type message..." disabled />
	<button id="sendBtn" onclick="send()" disabled>Send</button>

	<ul id="messages"></ul>

	<script>
		let ws;

		function connect() {
			const username = document.getElementById("username").value;
			if (!username) {
				alert("Enter username first");
				return;
			}

			ws = new WebSocket("ws://" + location.host + "/ws?username=" + username);

			ws.onopen = function() {
				document.getElementById("status").innerText = "You joined as: " + username;

				document.getElementById("msg").disabled = false;
				document.getElementById("sendBtn").disabled = false;

				document.getElementById("joinSection").style.display = "none";
			};

			ws.onmessage = function(event) {
				const li = document.createElement("li");
				li.innerText = event.data;
				document.getElementById("messages").appendChild(li);
			};
		}

		function send() {
			const input = document.getElementById("msg");
			if (ws && input.value !== "") {
				ws.send(input.value);
				input.value = "";
			}
		}
	</script>
</body>
</html>
`
