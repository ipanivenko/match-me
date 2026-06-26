package endpoints

import (
	"context"
	"log"
	"matchme-server/database"
	"matchme-server/internal"
	"matchme-server/structs"

	"github.com/gin-gonic/gin"
)

// POST /connections/:connectionId/action
// Accepts or rejects a connection request
func PostConnectionAction(c *gin.Context) {
	userID := c.GetString("userID")
	connectionID := c.Param("connectionId")

	var req struct {
		Action string `json:"action"` // "accept" or "reject"
	}

	if err := c.BindJSON(&req); err != nil || (req.Action != "accept" && req.Action != "reject") {
		c.JSON(400, structs.ErrorResponse{
			Message: "action must be 'accept' or 'reject'",
		})
		return
	}

	ctx := context.Background()

	// Determine status based on action
	status := req.Action + "ed" // "accepted" or "rejected"

	// Update connection status
	err := database.UpdateConnectionStatus(ctx, internal.DB, connectionID, userID, status)
	if err != nil {
		if err == database.ErrConnectionNotFound {
			c.JSON(404, structs.ErrorResponse{
				Message: "connection request not found or already processed",
			})
			return
		}
		log.Printf("Error updating connection status: %v", err)
		c.JSON(500, structs.ErrorResponse{
			Message: "db error",
		})
		return
	}

	log.Printf("✅ Connection %s %s by user %s", connectionID, status, userID)

	c.JSON(200, gin.H{
		"action": req.Action,
		"status": status,
	})
}