package game

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/sangjinsu/websocket-multiplayer/internal/models"
)

// Colors for players
var colors = []string{"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7", "#DDA0DD", "#98D8C8", "#F7DC6F"}

// Game represents the game instance
type Game struct {
	State *models.GameState
}

// NewGame creates a new game instance
func NewGame() *Game {
	return &Game{
		State: models.NewGameState(),
	}
}

// GenerateID generates a unique player ID
func (g *Game) GenerateID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// GetRandomColor returns a random color for a player
func (g *Game) GetRandomColor() string {
	return colors[rand.Intn(len(colors))]
}

// GetRandomPosition returns a random starting position that doesn't collide with other players
func (g *Game) GetRandomPosition() (float64, float64) {
	maxAttempts := 100
	playerRadius := 15.0
	minDistance := playerRadius * 2
	
	for attempt := 0; attempt < maxAttempts; attempt++ {
		x := float64(rand.Intn(800))
		y := float64(rand.Intn(600))
		
		// Check bounds
		if x < 15 || x > 785 || y < 15 || y > 585 {
			continue
		}
		
		// Check collision with other players
		collision := false
		for _, player := range g.State.Players {
			dx := x - player.X
			dy := y - player.Y
			distance := dx*dx + dy*dy
			
			if distance < minDistance*minDistance {
				collision = true
				break
			}
		}
		
		if !collision {
			return x, y
		}
	}
	
	// If no valid position found after max attempts, return center
	return 400, 300
}

// AddPlayer adds a player to the game
func (g *Game) AddPlayer(player *models.Player) {
	g.State.Mu.Lock()
	defer g.State.Mu.Unlock()
	
	// Assign player number
	g.State.PlayerCount++
	player.PlayerNum = g.State.PlayerCount
	
	// Use custom name if provided, otherwise use default
	if player.Name == "" {
		player.Name = fmt.Sprintf("Player %d", player.PlayerNum)
	}
	
	player.JoinedAt = time.Now()
	player.LastSeen = time.Now()
	
	g.State.Players[player.ID] = player
}

// RemovePlayer removes a player from the game
func (g *Game) RemovePlayer(playerID string) {
	g.State.Mu.Lock()
	defer g.State.Mu.Unlock()
	
	if _, exists := g.State.Players[playerID]; exists {
		delete(g.State.Players, playerID)
		
		// Reorder remaining players
		g.reorderPlayers()
	}
}

// reorderPlayers reorders player numbers after a player leaves
func (g *Game) reorderPlayers() {
	// Sort players by join time
	type playerInfo struct {
		id       string
		joinedAt time.Time
	}
	
	var players []playerInfo
	for id, player := range g.State.Players {
		players = append(players, playerInfo{id: id, joinedAt: player.JoinedAt})
	}
	
	// Sort by join time
	sort.Slice(players, func(i, j int) bool {
		return players[i].joinedAt.Before(players[j].joinedAt)
	})
	
	// Reassign player numbers
	for i, playerInfo := range players {
		if player, exists := g.State.Players[playerInfo.id]; exists {
			player.PlayerNum = i + 1
			player.Name = fmt.Sprintf("Player %d", player.PlayerNum)
		}
	}
	
	g.State.PlayerCount = len(g.State.Players)
}

// GetPlayer returns a player by ID
func (g *Game) GetPlayer(playerID string) *models.Player {
	g.State.Mu.RLock()
	defer g.State.Mu.RUnlock()
	return g.State.Players[playerID]
}

// 입력 큐: playerID -> 최근 입력
var inputQueue = make(map[string]string)

// ApplyInput: WASD 입력에 따라 속도 변경
func (g *Game) ApplyInput(playerID, key string) {
	g.State.Mu.Lock()
	defer g.State.Mu.Unlock()
	p, ok := g.State.Players[playerID]
	if !ok {
		return
	}
	const speed = 2.5
	switch key {
	case "w":
		p.Vy -= speed
	case "s":
		p.Vy += speed
	case "a":
		p.Vx -= speed
	case "d":
		p.Vx += speed
	}
}

// Tick: 모든 플레이어의 위치/속도/충돌/반동 등 물리 연산 수행
func (g *Game) Tick() {
	g.State.Mu.Lock()
	defer g.State.Mu.Unlock()
	const (
		friction = 0.98
		playerRadius = 15.0
		minDistance = playerRadius * 2
		w = 800.0
		h = 600.0
	)
	// 1. 속도 적용 및 마찰
	for _, p := range g.State.Players {
		p.X += p.Vx
		p.Y += p.Vy
		p.Vx *= friction
		p.Vy *= friction
		// 2. 경계 처리
		if p.X < playerRadius {
			p.X = playerRadius
			p.Vx *= -0.7
		}
		if p.X > w-playerRadius {
			p.X = w - playerRadius
			p.Vx *= -0.7
		}
		if p.Y < playerRadius {
			p.Y = playerRadius
			p.Vy *= -0.7
		}
		if p.Y > h-playerRadius {
			p.Y = h - playerRadius
			p.Vy *= -0.7
		}
	}
	// 3. 플레이어 간 충돌(탄성)
	players := g.State.Players
	for idA, a := range players {
		for idB, b := range players {
			if idA >= idB {
				continue
			}
			dx := b.X - a.X
			dy := b.Y - a.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < minDistance && dist > 0 {
				overlap := minDistance - dist
				// 각자 반씩 밀어내기
				pushX := (dx / dist) * (overlap / 2)
				pushY := (dy / dist) * (overlap / 2)
				a.X -= pushX
				a.Y -= pushY
				b.X += pushX
				b.Y += pushY
				// 탄성 충돌(속도 교환)
				nx, ny := dx/dist, dy/dist
				va := a.Vx*nx + a.Vy*ny
				vb := b.Vx*nx + b.Vy*ny
				a.Vx += (vb - va) * nx
				a.Vy += (vb - va) * ny
				b.Vx += (va - vb) * nx
				b.Vy += (va - vb) * ny
			}
		}
	}
}

// GetAllPlayers returns all players (순수 데이터만)
func (g *Game) GetAllPlayers() map[string]*models.Player {
	g.State.Mu.RLock()
	defer g.State.Mu.RUnlock()
	players := make(map[string]*models.Player, len(g.State.Players))
	for id, p := range g.State.Players {
		players[id] = &models.Player{
			ID:        p.ID,
			PlayerNum: p.PlayerNum,
			Name:      p.Name,
			X:         p.X,
			Y:         p.Y,
			Vx:        p.Vx,
			Vy:        p.Vy,
			Color:     p.Color,
			JoinedAt:  p.JoinedAt,
			LastSeen:  p.LastSeen,
		}
	}
	return players
}

// UpdatePlayerPosition updates a player's position with bounce collision
func (g *Game) UpdatePlayerPosition(playerID string, x, y float64) {
	g.State.Mu.Lock()
	defer g.State.Mu.Unlock()
	
	if player, exists := g.State.Players[playerID]; exists {
		// Check if the new position is within bounds (canvas: 800x600, player radius: 15)
		if x < 15 || x > 785 || y < 15 || y > 585 {
			// Position is outside bounds, don't update
			return
		}
		
		// Check collision with other players and calculate bounce
		playerRadius := 15.0
		minDistance := playerRadius * 2 // Minimum distance between player centers
		bounceX := x
		bounceY := y
		hasCollision := false
		
		for id, otherPlayer := range g.State.Players {
			if id == playerID {
				continue // Skip self
			}
			
			// Calculate distance between players
			dx := x - otherPlayer.X
			dy := y - otherPlayer.Y
			distance := dx*dx + dy*dy // Using squared distance for efficiency
			
			if distance < minDistance*minDistance {
				hasCollision = true
				
				// Calculate bounce direction
				angle := math.Atan2(dy, dx)
				
				// Calculate bounce position
				bounceX = otherPlayer.X + math.Cos(angle)*minDistance
				bounceY = otherPlayer.Y + math.Sin(angle)*minDistance
				
				// Keep bounce position within bounds
				if bounceX < 15 {
					bounceX = 15
				} else if bounceX > 785 {
					bounceX = 785
				}
				if bounceY < 15 {
					bounceY = 15
				} else if bounceY > 585 {
					bounceY = 585
				}
				
				break // Handle first collision only
			}
		}
		
		if hasCollision {
			// Update to bounce position
			player.X = bounceX
			player.Y = bounceY
		} else {
			// Update to original position
			player.X = x
			player.Y = y
		}
		
		player.LastSeen = time.Now()
	}
} 