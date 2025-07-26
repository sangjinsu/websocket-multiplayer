package models

import (
	"time"

	"github.com/gofiber/websocket/v2"
)

// Player represents a connected player
type Player struct {
	ID          string    `json:"id"`
	PlayerNum   int       `json:"playerNum"`   // 접속 순서 (1, 2, 3...)
	Name        string    `json:"name"`        // 플레이어 이름 (Player 1, Player 2...)
	X           float64   `json:"x"`
	Y           float64   `json:"y"`
	Vx          float64   `json:"vx"`
	Vy          float64   `json:"vy"`
	Color       string    `json:"color"`
	JoinedAt    time.Time `json:"joinedAt"`    // 최초 접속 시간
	LastSeen    time.Time `json:"lastSeen"`    // 마지막 활동 시간
	Conn        *websocket.Conn
}

// PlayerMove represents a player movement
type PlayerMove struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// PlayerJoin represents a player joining event
type PlayerJoin struct {
	ID    string `json:"id"`
	Color string `json:"color"`
}

// PlayerLogin represents a player login request
type PlayerLogin struct {
	Name         string  `json:"name"`
	Color        string  `json:"color"`
	LastPosition *struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"lastPosition"`
} 