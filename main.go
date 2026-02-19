package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sockets/handlers"
	"sockets/internal/chat"
	"os"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	hub := chat.NewHub()

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/ws", handlers.WebSocketHandler(hub))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)

}
