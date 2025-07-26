package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/sangjinsu/websocket-multiplayer/internal/game"
	ws "github.com/sangjinsu/websocket-multiplayer/internal/websocket"
)

func main() {
	app := fiber.New()

	// Create game instance
	gameInstance := game.NewGame()

	// Create websocket handler
	wsHandler := ws.NewHandler(gameInstance)

	// Serve static files
	app.Static("/", "./public")

	// WebSocket upgrade handler
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket handler
	app.Get("/ws", websocket.New(wsHandler.HandleWebSocket))

	// Start the server
	log.Println("Server starting on :3000")
	log.Fatal(app.Listen(":3000"))
}