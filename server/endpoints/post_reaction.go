package endpoints

import (
	"context"
	"log"
	"matchme-server/database"
	"matchme-server/internal"
	"matchme-server/structs"

	"github.com/gin-gonic/gin"
)

// POST /api/recommendations/:targetUserId/reaction
func PostReaction(c *gin.Context) {
	userID := c.GetString("userID")
	targetID := c.Param("targetUserId")

	var req struct {
		Reaction string `json:"reaction"` // "like" | "dislike"
	}
	
	if err := c.BindJSON(&req); err != nil || (req.Reaction != "like" && req.Reaction != "dislike") {
		c.JSON(400, structs.ErrorResponse{
			Message: "reaction must be 'like' or 'dislike'",
		})
		return
	}

	ctx := context.Background()
	
	// Record the reaction in user_reactions (for tracking history)
	err := database.UpsertReaction(ctx, internal.DB, userID, targetID, database.Reaction(req.Reaction))
	if err != nil {
		log.Printf("ERROR updating reaction: %v", err)
		c.JSON(500, structs.ErrorResponse{
			Message: "db error",
		})
		return
	}

	// If it's a LIKE, create connection request immediately
	if req.Reaction == "like" {
		log.Printf("👍 User %s liked %s - creating connection request", userID, targetID)
		
		connectionID, err := database.CreateConnectionRequest(ctx, internal.DB, userID, targetID)
		if err != nil {
			if err == database.ErrConnectionExists {
				log.Printf("Connection already exists between users")
			} else {
				log.Printf("ERROR creating connection request: %v", err)
			}
		} else {
			log.Printf("✅ Connection request created: %s (from %s to %s)", connectionID, userID, targetID)
		}
	}

	c.JSON(200, gin.H{
		"reaction": req.Reaction,
		"status":   "ok",
	})
}