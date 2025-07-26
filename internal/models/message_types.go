package models

// MessageType represents the type of websocket message
type MessageType string

const (
	// Welcome message sent to new player
	MessageTypeWelcome MessageType = "welcome"
	
	// Game state sent to player
	MessageTypeGameState MessageType = "game_state"
	
	// Player joined the game
	MessageTypePlayerJoin MessageType = "player_join"
	
	// Player left the game
	MessageTypePlayerLeave MessageType = "player_leave"
	
	// Player moved
	MessageTypePlayerMove MessageType = "player_move"
	
	// Player movement from client
	MessageTypeMove MessageType = "move"
	
	// Player reconnection
	MessageTypeReconnect MessageType = "reconnect"
	
	// Player login
	MessageTypeLogin MessageType = "login"
	
	// Player collision
	MessageTypeCollision MessageType = "collision"

	// Player input
	MessageTypeInput MessageType = "input"
) 