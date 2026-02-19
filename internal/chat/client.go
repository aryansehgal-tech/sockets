package chat

import "github.com/gorilla/websocket"

type Client struct {
	Conn		*websocket.Conn
	Username 	string
	Room		string
}