package database

import (
	"context"
	"errors"
	"matchme-server/structs"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connection errors for error's handling
var (
	ErrConnectionExists   = errors.New("connection already exists")
	ErrConnectionNotFound = errors.New("connection not found")
)

// GetUserConnections retrieves list of connected users (accepted connections)
func GetUserConnections(ctx context.Context, pool *pgxpool.Pool, userID string) ([]string, error) {
	const query = `
		SELECT CASE 
			WHEN requester_user_id = $1 THEN target_user_id::text
			ELSE requester_user_id::text
		END as connected_user_id
		FROM connections 
		WHERE (requester_user_id = $1 OR target_user_id = $1) 
			AND status = 'accepted'
		ORDER BY updated_at DESC`

	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []string
	for rows.Next() {
		var connectedUserID string
		if err := rows.Scan(&connectedUserID); err != nil {
			continue // Skip rows with scan errors
		}
		connections = append(connections, connectedUserID)
	}

	return connections, nil
}

// CreateConnectionRequest creates a new connection request between users
func CreateConnectionRequest(ctx context.Context, pool *pgxpool.Pool, requesterID, targetID string) (string, error) {
	// Check if connection already exists (in either direction)
	var exists bool
	err := pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM connections 
			WHERE (requester_user_id = $1 AND target_user_id = $2)
				OR (requester_user_id = $2 AND target_user_id = $1)
		)`, requesterID, targetID).Scan(&exists)
	
	if err != nil {
		return "", err
	}
	
	if exists {
		return "", ErrConnectionExists
	}

	// Create new connection request with pending status
	const query = `
		INSERT INTO connections (requester_user_id, target_user_id, status)
		VALUES ($1, $2, 'pending')
		RETURNING id::text`

	var connectionID string
	err = pool.QueryRow(ctx, query, requesterID, targetID).Scan(&connectionID)
	return connectionID, err
}

// GetIncomingConnectionRequests retrieves pending connection requests for a user
func GetIncomingConnectionRequests(ctx context.Context, pool *pgxpool.Pool, userID string) ([]structs.Connection, error) {
	const query = `
		SELECT id::text, requester_user_id::text, target_user_id::text, 
			   status, created_at, updated_at
		FROM connections 
		WHERE target_user_id = $1 AND status = 'pending'
		ORDER BY created_at DESC`

	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []structs.Connection
	for rows.Next() {
		var req structs.Connection
		err := rows.Scan(&req.ID, &req.RequesterUserID, &req.TargetUserID, 
						 &req.Status, &req.CreatedAt, &req.UpdatedAt)
		if err != nil {
			continue // Skip rows with scan errors
		}
		requests = append(requests, req)
	}

	return requests, nil
}

// UpdateConnectionStatus updates the status of a connection request
// Only the target user can accept/reject pending requests
func UpdateConnectionStatus(ctx context.Context, pool *pgxpool.Pool, connectionID, userID, status string) error {
	const query = `
		UPDATE connections 
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND target_user_id = $3 AND status = 'pending'`

	result, err := pool.Exec(ctx, query, status, connectionID, userID)
	if err != nil {
		return err
	}
	
	// Check if any rows were affected (connection exists and belongs to user)
	if result.RowsAffected() == 0 {
		return ErrConnectionNotFound
	}

	var requesterID, targetID string
		var chatID string
	if status == "accepted" {
    // Get the two user IDs from the specific connection
    err := pool.QueryRow(ctx, 
        "SELECT requester_user_id, target_user_id FROM connections WHERE id = $1", 
        connectionID,
    ).Scan(&requesterID, &targetID)
    if err != nil {
        return err
    }

    // Create chat room with these two users
      err = pool.QueryRow(ctx, `
        INSERT INTO chats (id, user1_id, user2_id, created_at) 
        VALUES (gen_random_uuid(), $1, $2, NOW())
        RETURNING id`,
        requesterID, targetID,
    ).Scan(&chatID)
    if err != nil {
        return err
    }
	}
	
	return nil
}