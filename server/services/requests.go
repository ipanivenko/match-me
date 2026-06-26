package services

import (
	"context"
	"log"
	"matchme-server/database"
	"matchme-server/internal"
	"matchme-server/structs"

	"github.com/gin-gonic/gin"
)

func GetRequests(c *gin.Context) {
	userID := c.GetString("userID")
	ctx := context.Background()
	
	percent, err := database.GetProfileCompletionPercent(ctx, internal.DB, userID)
	if err != nil {
		log.Println(err)
		c.JSON(500, structs.ErrorResponse{
			Message: CommonErr,
		})
		return
	}
	
	if percent < 100 {
		c.JSON(400, structs.ErrorResponse{
			Message: "profile did not complete",
		})
		return
	}
	
	// Get connection requests (with connection IDs)
	connectionRequests, err := database.GetIncomingConnectionRequests(ctx, internal.DB, userID)
	if err != nil {
		log.Println(err)
		c.JSON(500, structs.ErrorResponse{
			Message: CommonErr,
		})
		return
	}
	
	// Extract requester user IDs for response
	var userIDs []string
	connectionMap := make(map[string]string) // user_id -> connection_id
	
	for _, conn := range connectionRequests {
		userIDs = append(userIDs, conn.RequesterUserID)
		connectionMap[conn.RequesterUserID] = conn.ID
	}
	
	log.Printf("Found %d connection requests for user %s", len(userIDs), userID)
	
	c.JSON(200, gin.H{
		"user_ids": userIDs,
		"connection_map": connectionMap, // Frontend can use this to get connection_id
	})
}