package endpoints

import (
	"matchme-server/database"
	"matchme-server/internal"
	"github.com/gin-gonic/gin"
)

func DeleteMePhoto(c *gin.Context) {
    userID := c.GetString("userID")

    if err := database.DeleteUserPhoto(c, internal.DB, userID); err != nil {
        c.JSON(500, gin.H{"error": "Failed to delete photo"})
        return
    }

    c.Status(204)
}
