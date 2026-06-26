package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SaveUserPhoto(ctx context.Context, pool *pgxpool.Pool, userID, publicID string, version int) error {
    _, err := pool.Exec(ctx, `
        INSERT INTO user_photos (user_id, photo_public_id, photo_version)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id)
        DO UPDATE SET
            photo_public_id = EXCLUDED.photo_public_id,
            photo_version   = EXCLUDED.photo_version
    `, userID, publicID, version)
    return err
}

func DeleteUserPhoto(ctx context.Context, pool *pgxpool.Pool, userID string) error {
    _, err := pool.Exec(ctx, `DELETE FROM user_photos WHERE user_id = $1`, userID)
    return err
}

