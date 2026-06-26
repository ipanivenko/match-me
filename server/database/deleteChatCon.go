package database // Or wherever your DB functions live

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DeleteConnectionAndChat safely deletes both the connection and the associated
// chat in a single, atomic transaction.
func DeleteConnectionAndChat(ctx context.Context, pool *pgxpool.Pool, userID string, targetUserID string) error {
	
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		DELETE FROM connections 
		WHERE (requester_user_id = $1 AND target_user_id = $2)
		   OR (requester_user_id = $2 AND target_user_id = $1)
	`, userID, targetUserID)
	
    if err != nil {
		log.Printf("ERROR deleting connection during tx: %v", err)

		return fmt.Errorf("failed to delete connection: %w", err)
	}

	// 4. Delete the chat (also using 'tx')
	_, err = tx.Exec(ctx, `
		DELETE FROM chats 
		WHERE (user1_id = $1 AND user2_id = $2)
		   OR (user1_id = $2 AND user1_id = $1)
	`, userID, targetUserID)
	
    if err != nil {
	
		log.Printf("ERROR deleting chat during tx: %v", err)
	
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("ERROR committing delete transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}