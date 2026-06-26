package services

import (
	"context"
	"log"
	"matchme-server/database"
	"matchme-server/internal"
	"matchme-server/structs"

	"github.com/gin-gonic/gin"
)

func GetConnections(c *gin.Context){
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
	
	connections, err := database.GetConnections(ctx, internal.DB, userID)
	if err != nil {
		log.Println(err)
		c.JSON(500, structs.ErrorResponse{
			Message: CommonErr,
		})
		return
	}

	c.JSON(200, connections)
}