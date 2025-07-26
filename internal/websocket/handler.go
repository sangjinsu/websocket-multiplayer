package ws

import (
	"encoding/json"
	"log"
	"math"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/sangjinsu/websocket-multiplayer/internal/game"
	"github.com/sangjinsu/websocket-multiplayer/internal/models"
)

// Handler represents the websocket handler
type Handler struct {
	game *game.Game
	tickOnce sync.Once
	lastGameState map[string]*models.Player // 이전 게임 상태 저장
}

// NewHandler creates a new websocket handler
func NewHandler(game *game.Game) *Handler {
	return &Handler{
		game: game,
	}
}

// HandleWebSocket handles websocket connections
func (h *Handler) HandleWebSocket(c *websocket.Conn) {
	// Generate unique player ID
	playerID := h.game.GenerateID()
	
	// Create temporary player for connection
	player := &models.Player{
		ID:       playerID,
		Conn:     c,
		LastSeen: time.Now(),
	}

	log.Printf("New connection established: %s", playerID)

	// 최초 1회만: 서버 전체에서 60fps 물리 tick 루프 실행
	h.tickOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(16 * time.Millisecond)
			defer ticker.Stop()
			for {
				<-ticker.C
				h.game.Tick()
				h.broadcastGameState()
			}
		}()
	})

	// Handle incoming messages
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("Player %s disconnected: %v", playerID, err)
			break
		}

		var message models.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		h.handleMessage(player, message)
	}

	// Get player name before removal
	playerName := ""
	if player.Name != "" {
		playerName = player.Name
	}

	// Remove player from game
	h.game.RemovePlayer(playerID)

	// Broadcast player leave
	h.broadcastPlayerLeave(playerID)

	if playerName != "" {
		log.Printf("Player %s (%s) left the game", playerName, playerID)
	} else {
		log.Printf("Player %s left the game", playerID)
	}
}

func (h *Handler) handleMessage(player *models.Player, message models.Message) {
	switch message.Type {
	case models.MessageTypeLogin:
		// Handle player login
		if payload, ok := message.Payload.(map[string]any); ok {
			name, _ := payload["name"].(string)
			color, _ := payload["color"].(string)
			
			log.Printf("Login attempt from player: %s (ID: %s)", name, player.ID)
			
			// Set player properties
			player.Name = name
			if color != "" {
				player.Color = color
			} else {
				player.Color = h.game.GetRandomColor()
			}
			
			// Set position
			if lastPos, ok := payload["lastPosition"].(map[string]any); ok {
				if x, ok := lastPos["x"].(float64); ok {
					player.X = x
				}
				if y, ok := lastPos["y"].(float64); ok {
					player.Y = y
				}
			} else {
				// Use random position if no last position
				player.X, player.Y = h.game.GetRandomPosition()
			}
			
			// Add player to game
			h.game.AddPlayer(player)
			
			// Send welcome message
			welcomeMsg := models.Message{
				Type: models.MessageTypeWelcome,
				Payload: map[string]interface{}{
					"id":        player.ID,
					"playerNum": player.PlayerNum,
					"name":      player.Name,
					"color":     player.Color,
				},
			}
			h.sendMessage(player.Conn, welcomeMsg)
			
			// Broadcast new player to all other players
			h.broadcastPlayerJoin(player)
			
			// Send current game state to new player
			h.sendGameState(player.Conn)
			
			log.Printf("Player %s (%s) joined the game", player.Name, player.ID)
		}
		
	case models.MessageTypeInput:
		if payload, ok := message.Payload.(map[string]any); ok {
			// Handle WASD input
			if key, ok := payload["key"].(string); ok {
				h.game.ApplyInput(player.ID, key)
			}
			
			// Handle touch/click movement input
			if vx, ok := payload["vx"].(float64); ok {
				if vy, ok := payload["vy"].(float64); ok {
					h.game.ApplyVelocityInput(player.ID, vx, vy)
				}
			}
		}
		return
		
	case models.MessageTypeCollision:
		// Handle player collision
		if payload, ok := message.Payload.(map[string]any); ok {
			myID, _ := payload["myId"].(string)
			partnerID, _ := payload["partnerId"].(string)
			myNewX, _ := payload["myNewX"].(float64)
			myNewY, _ := payload["myNewY"].(float64)
			partnerX, _ := payload["partnerX"].(float64)
			partnerY, _ := payload["partnerY"].(float64)
			
			// Update my position
			h.game.UpdatePlayerPosition(myID, myNewX, myNewY)
			
			// Calculate partner's bounce position (opposite direction)
			partner := h.game.GetPlayer(partnerID)
			if partner != nil {
				// Get collision angle from client or calculate it
				var collisionAngle float64
				if angle, ok := payload["collisionAngle"].(float64); ok {
					collisionAngle = angle
				} else {
					// Fallback: calculate angle from positions
					collisionAngle = math.Atan2(myNewY - partnerY, myNewX - partnerX)
				}
				
				// Partner should move in the opposite direction (add π to angle)
				oppositeAngle := collisionAngle + math.Pi
				
				// Calculate partner's new position in opposite direction
				partnerNewX := partnerX + math.Cos(oppositeAngle) * 30 // 30 = minDistance
				partnerNewY := partnerY + math.Sin(oppositeAngle) * 30
				
				// Keep partner position within bounds
				if partnerNewX < 15 {
					partnerNewX = 15
				} else if partnerNewX > 785 {
					partnerNewX = 785
				}
				if partnerNewY < 15 {
					partnerNewY = 15
				} else if partnerNewY > 585 {
					partnerNewY = 585
				}
				
				// Update partner position
				h.game.UpdatePlayerPosition(partnerID, partnerNewX, partnerNewY)
				
				// Broadcast both movements
				h.broadcastPlayerMove(player)
				h.broadcastPlayerMove(partner)
				
				log.Printf("Collision between %s and %s - opposite bounce applied", 
					player.Name, partner.Name)
			}
		}
		
	case models.MessageTypeReconnect:
		// Handle player reconnection with existing ID
		if payload, ok := message.Payload.(map[string]interface{}); ok {
			if id, ok := payload["id"].(string); ok {
				// Check if player ID already exists
				if existingPlayer := h.game.GetPlayer(id); existingPlayer != nil {
					// Update existing player's connection
					existingPlayer.Conn = player.Conn
					existingPlayer.LastSeen = time.Now()
					
					// Remove the new player and use existing one
					h.game.RemovePlayer(player.ID)
					
					// Send welcome message with existing info
					welcomeMsg := models.Message{
						Type: models.MessageTypeWelcome,
						Payload: map[string]interface{}{
							"id":    existingPlayer.ID,
							"color": existingPlayer.Color,
						},
					}
					h.sendMessage(player.Conn, welcomeMsg)
					
					// Send current game state
					h.sendGameState(player.Conn)
					
					log.Printf("Player %s reconnected", id)
					return
				}
			}
		}
	}
}

func (h *Handler) broadcastPlayerJoin(player *models.Player) {
	msg := models.Message{
		Type: models.MessageTypePlayerJoin,
		Payload: map[string]interface{}{
			"id":        player.ID,
			"playerNum": player.PlayerNum,
			"name":      player.Name,
			"x":         player.X,
			"y":         player.Y,
			"color":     player.Color,
		},
	}

	// 원본 State.Players에서 Conn이 있는 플레이어들에게만 전송
	h.game.State.Mu.RLock()
	defer h.game.State.Mu.RUnlock()
	for _, p := range h.game.State.Players {
		if p.ID != player.ID && p.Conn != nil {
			h.sendMessage(p.Conn, msg)
		}
	}
}

func (h *Handler) broadcastPlayerLeave(playerID string) {
	msg := models.Message{
		Type: models.MessageTypePlayerLeave,
		Payload: map[string]string{
			"id": playerID,
		},
	}

	// 원본 State.Players에서 Conn이 있는 플레이어들에게만 전송
	h.game.State.Mu.RLock()
	defer h.game.State.Mu.RUnlock()
	for _, p := range h.game.State.Players {
		if p.Conn != nil {
			h.sendMessage(p.Conn, msg)
		}
	}
}

func (h *Handler) broadcastPlayerMove(player *models.Player) {
	msg := models.Message{
		Type: models.MessageTypePlayerMove,
		Payload: map[string]any{
			"id": player.ID,
			"x":  player.X,
			"y":  player.Y,
		},
	}

	// 원본 State.Players에서 Conn이 있는 플레이어들에게만 전송
	h.game.State.Mu.RLock()
	defer h.game.State.Mu.RUnlock()
	for _, p := range h.game.State.Players {
		if p.ID != player.ID && p.Conn != nil {
			h.sendMessage(p.Conn, msg)
		}
	}
}

func (h *Handler) sendGameState(conn *websocket.Conn) {
	players := h.game.GetAllPlayers()
	
	// Create a copy without websocket connections for JSON serialization
	playersCopy := make(map[string]*models.Player)
	for id, player := range players {
		playersCopy[id] = &models.Player{
			ID:        player.ID,
			PlayerNum: player.PlayerNum,
			Name:      player.Name,
			X:         player.X,
			Y:         player.Y,
			Color:     player.Color,
			JoinedAt:  player.JoinedAt,
			LastSeen:  player.LastSeen,
		}
	}

	msg := models.Message{
		Type:    models.MessageTypeGameState,
		Payload: playersCopy,
	}

	h.sendMessage(conn, msg)
}

func (h *Handler) sendMessage(conn *websocket.Conn, message models.Message) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Error sending message: %v", err)
	}
} 

// 모든 플레이어에게 현재 상태 브로드캐스트 (변경사항이 있을 때만)
func (h *Handler) broadcastGameState() {
	players := h.game.GetAllPlayers()
	
	// 현재 상태를 JSON으로 직렬화하여 변경사항 확인
	currentState, err := json.Marshal(players)
	if err != nil {
		log.Printf("Error marshaling current game state: %v", err)
		return
	}
	
	// 이전 상태와 비교하여 변경사항이 있는지 확인
	if h.lastGameState != nil {
		lastState, _ := json.Marshal(h.lastGameState)
		if string(currentState) == string(lastState) {
			return // 변경사항이 없으면 브로드캐스트하지 않음
		}
	}
	
	// 변경사항이 있으면 브로드캐스트
	msg := models.Message{
		Type:    models.MessageTypeGameState,
		Payload: players,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling game state: %v", err)
		return
	}
	
	// 원본 State.Players에서 Conn이 있는 플레이어들에게만 전송
	h.game.State.Mu.RLock()
	defer h.game.State.Mu.RUnlock()
	for _, p := range h.game.State.Players {
		if p.Conn != nil {
			_ = p.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
	
	// 현재 상태를 이전 상태로 저장
	h.lastGameState = players
} 