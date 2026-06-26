package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetRequests(ctx context.Context, pool *pgxpool.Pool, userID string) ([]string, error) {
	const query = `
	SELECT ur.user_id
FROM user_reactions ur
WHERE ur.target_user_id = $1
  AND ur.reaction = 'like'
  AND (ur.is_match IS DISTINCT FROM TRUE)
  AND NOT EXISTS (
    SELECT 1
    FROM user_reactions ur2
    WHERE ur2.user_id = $1
      AND ur2.target_user_id = ur.user_id
      AND ur2.reaction = 'dislike'
  );`

  rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []string
	for rows.Next(){
		var request string
		rows.Scan(&request)
		requests = append(requests, request)
	}
	
	return requests, nil
}
