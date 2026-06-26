package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateUser inserts a new user and returns its id + created_at.
func CreateUser(ctx context.Context, pool *pgxpool.Pool, email, hashedPassword string) (string, time.Time, error) {
	const q = `
		INSERT INTO users (id, email, password_hash, created_at)
		VALUES (uuid_generate_v4(), $1, $2, NOW())
		RETURNING id, created_at
	`

	var (
		id        string
		createdAt time.Time
	)
	err := pool.QueryRow(ctx, q, email, hashedPassword).Scan(&id, &createdAt)
	if err != nil {
		return "", time.Time{}, err
	}

	//we add user_id in table parent_profiles, children in order not to fetch a db error
	const q2 = `
		INSERT INTO parent_profiles (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING
	`
	_, err = pool.Exec(ctx, q2, id)
	if err != nil {
		return "", time.Time{}, err
	}

	const q3 = `
		INSERT INTO children (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING
	`
	_, err = pool.Exec(ctx, q3, id)
	if err != nil {
		return "", time.Time{}, err
	}

	const q4 = `
		INSERT INTO matching_preferences (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING
	`
	_, err = pool.Exec(ctx, q4, id)
	if err != nil {
		return "", time.Time{}, err
	}

	return id, createdAt, nil
}

// for login
func GetUserByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (string, string, error) {
	const q = `SELECT id, password_hash FROM users WHERE email=$1`

	var (
		id     string
		pwHash string
	)

	err := pool.QueryRow(ctx, q, email).Scan(&id, &pwHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", pgx.ErrNoRows
		}
		return "", "", err
	}

	return id, pwHash, nil
}

func GetUserEmailCreatedAt(ctx context.Context, pool *pgxpool.Pool, userID string) (string, time.Time, error) {
	const q = `SELECT email, created_at FROM users WHERE id=$1`

	var created_at time.Time
	var email string

	err := pool.QueryRow(ctx, q, userID).Scan(&email, &created_at)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", time.Time{} ,  pgx.ErrNoRows
		}
		return "", time.Time{},  err
	}

	return email, created_at, nil
}

