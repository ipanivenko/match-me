package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func EnsureUserExists(ctx context.Context, pool *pgxpool.Pool, id string) (error, bool) {
	var exists bool
	err := pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`, id).Scan(&exists)
	if err != nil {
		return err, false
	}
	if !exists {
		return nil, false
	}
	return nil, true
}

// Returns: name, photoURL (versioned). If no photo, photoURL = "".
func GetUserNamePhotoURL(ctx context.Context, pool *pgxpool.Pool, userID string) (string, string) {
	var (
		name     string
		publicID *string 
		version  *int32  
	)

	err := pool.QueryRow(ctx, `
		SELECT 
			COALESCE(pp.name, '') AS name,
			up.photo_public_id,
			up.photo_version
		FROM parent_profiles pp
		LEFT JOIN user_photos up ON up.user_id = pp.user_id
		WHERE pp.user_id = $1
		LIMIT 1
	`, userID).Scan(&name, &publicID, &version)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ""
		}
		log.Println("GetUserNamePhotoURL query error:", err)
		return "", ""
	}

	// Build a versioned Cloudinary URL if photo exists
	photoURL := ""
	if publicID != nil && version != nil {
		cloud := os.Getenv("CLOUDINARY_CLOUD_NAME")
		if cloud == "" {
			log.Println("CLOUDINARY_CLOUD_NAME not set; returning empty photo URL")
		} else {
			photoURL = fmt.Sprintf(
				"https://res.cloudinary.com/%s/image/upload/c_fill,w_256,h_256,g_face,f_auto,q_auto/v%d/%s",
				cloud, int(*version), *publicID,
			)
		}
	}

	return name, photoURL
}