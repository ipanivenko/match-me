package database

import (
	"context"
	"fmt"
	"matchme-server/structs"
	"time"
	"log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserChats(ctx context.Context, pool *pgxpool.Pool, userID string) ([]structs.Chat, error) {
	const query = `
		SELECT 
			c.id, 
			c.user1_id, 
			c.user2_id, 
			c.created_at,
			COALESCE(
				(SELECT COUNT(*) 
				 FROM unread_messages um 
				 JOIN messages m ON m.id = um.message_id
				 WHERE um.user_id = $1 AND m.chat_id = c.id),
				0
			) as unread_count
		FROM chats c
		WHERE c.user1_id = $1 OR c.user2_id = $1 
		ORDER BY c.created_at DESC`
	
	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user chats: %w", err)
	}
	defer rows.Close()
	
	var chats []structs.Chat
	for rows.Next() {
		var chat structs.Chat
		var unreadCount int
		err := rows.Scan(&chat.ID, &chat.User1ID, &chat.User2ID, &chat.CreatedAt, &unreadCount)
		if err != nil {
			continue
		}
		chat.UnreadCount = unreadCount
		chats = append(chats, chat)
	}
	
	return chats, nil
}

// GetChatByID returns a specific chat if the user is a participant
func GetChatByID(ctx context.Context, pool *pgxpool.Pool, chatID, userID string) (*structs.Chat, error) {
	
	const query = `
		SELECT id, user1_id, user2_id, created_at 
		FROM chats 
		WHERE id = $1 AND (user1_id = $2 OR user2_id = $2)`

	var chat structs.Chat
	err := pool.QueryRow(ctx, query, chatID, userID).Scan(
		&chat.ID, &chat.User1ID, &chat.User2ID, &chat.CreatedAt,
	)
	if err != nil {
		log.Printf("ERROR in GetChatByID: %v", err)
		return nil, fmt.Errorf("chat not found or access denied: %w", err)
	}
	
	return &chat, nil
}

// GetChatMessages returns messages for a chat with pagination
func GetChatMessages(ctx context.Context, pool *pgxpool.Pool, chatID, userID, before string, limit int) ([]structs.ChatMessage, error) {
	// Verify user has access to this chat
	_, err := GetChatByID(ctx, pool, chatID, userID)
	if err != nil {
		return nil, err
	}

	var query string
	var rows pgx.Rows

	if before != "" {
		// Pagination using UUID
		query = `
			SELECT id, chat_id, sender_id, content, created_at 
			FROM messages 
			WHERE chat_id = $1 AND id::text < $2
			ORDER BY created_at DESC 
			LIMIT $3`
		rows, err = pool.Query(ctx, query, chatID, before, limit)
	} else {
		// First page - no pagination
		query = `
			SELECT id, chat_id, sender_id, content, created_at 
			FROM messages 
			WHERE chat_id = $1
			ORDER BY created_at DESC 
			LIMIT $2`
		rows, err = pool.Query(ctx, query, chatID, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []structs.ChatMessage
	for rows.Next() {
		var msg structs.ChatMessage
		err := rows.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.CreatedAt)
		if err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	// Reverse order so oldest messages are first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
// SaveChatMessage saves a message to the database
func SaveChatMessage(ctx context.Context, pool *pgxpool.Pool, chatID, senderID, content string) (string, error) {
	// Verify sender has access to this chat
	_, err := GetChatByID(ctx, pool, chatID, senderID)
	if err != nil {
		return "", err
	}

	const query = `
		INSERT INTO messages (chat_id, sender_id, content, created_at) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`

	var messageID string
	err = pool.QueryRow(ctx, query, chatID, senderID, content, time.Now().UTC()).Scan(&messageID)
	if err != nil {
		return "", fmt.Errorf("failed to save message: %w", err)
	}

	return messageID, nil
}

// MarkChatMessagesAsRead marks all messages in a chat as read for the user
func MarkChatMessagesAsRead(ctx context.Context, pool *pgxpool.Pool, chatID, userID string) error {
	const query = `
		DELETE FROM unread_messages
		WHERE user_id = $1 
		  AND message_id IN (
			SELECT m.id 
			FROM messages m 
			WHERE m.chat_id = $2
		  )`
	
	_, err := pool.Exec(ctx, query, userID, chatID)
	return err
}