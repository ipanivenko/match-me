package structs

import "time"

// Chat represents a chat room between two users
type Chat struct {
	ID        string    `json:"id"`
	User1ID   string    `json:"user1_id"`
	User2ID   string    `json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
	UnreadCount int       `json:"unread_count"`

}

// ChatMessage represents a message in a chat
type ChatMessage struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Connection represents a connection request between users
type Connection struct {
	ID              string    `json:"id"`
	RequesterUserID string    `json:"requester_user_id"`
	TargetUserID    string    `json:"target_user_id"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}