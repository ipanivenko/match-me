package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetProfileCompletionPercent(ctx context.Context, db *pgxpool.Pool, userID string) (float64, error) {
    var percent float64
    query := `SELECT profile_completion_percent($1);`
    err := db.QueryRow(ctx, query, userID).Scan(&percent)
    if err != nil {
        return 0, fmt.Errorf("get profile completion: %w", err)
    }
    return percent, nil
}
