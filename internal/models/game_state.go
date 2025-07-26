package models

import "sync"

// GameState represents the current state of the game
type GameState struct {
	Players      map[string]*Player `json:"players"`
	PlayerCount  int                `json:"playerCount"`  // 총 플레이어 수
	Mu           sync.RWMutex
}

// NewGameState creates a new game state
func NewGameState() *GameState {
	return &GameState{
		Players: make(map[string]*Player),
	}
} 