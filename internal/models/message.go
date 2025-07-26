package models

// Message represents a websocket message
type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

// PlayerCollision represents a collision between two players
type PlayerCollision struct {
	MyID       string  `json:"myId"`
	PartnerID  string  `json:"partnerId"`
	MyNewX     float64 `json:"myNewX"`
	MyNewY     float64 `json:"myNewY"`
	PartnerX   float64 `json:"partnerX"`
	PartnerY   float64 `json:"partnerY"`
} 