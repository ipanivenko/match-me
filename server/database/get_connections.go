package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConnections(ctx context.Context, pool *pgxpool.Pool, userID string) ([]string, error) {
	const query = `
		SELECT 
			CASE 
				WHEN c.requester_user_id = $1 THEN c.target_user_id
				ELSE c.requester_user_id
			END AS other_user_id
		FROM connections c
		WHERE (c.requester_user_id = $1 OR c.target_user_id = $1)
		  AND c.status = 'accepted'
	`
	
	rows, err := pool.Query(ctx, query, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []string

	for rows.Next() {
		var connection string
		if err := rows.Scan(&connection); err != nil {
			continue
		}
		connections = append(connections, connection)
	}

	return connections, nil
}

