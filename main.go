package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"chatapp/handlers"
	"chatapp/internal/chat"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	hub := chat.NewHub()

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/ws", handlers.WebSocketHandler(hub))

	r.Run("0.0.0.0:8080")
}
