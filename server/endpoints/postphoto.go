package endpoints

import (
	"log"
	"matchme-server/database"
	"matchme-server/internal"
	"matchme-server/structs"

	"github.com/gin-gonic/gin"
)

func PostMePhoto(c *gin.Context) {
    userID := c.GetString("userID")

    var body struct {
        PublicID string `json:"public_id" binding:"required"`
        Version  int    `json:"version" binding:"required"`
    }

    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(400, structs.ErrorResponse {
            Message: "Invalid payload",
        })
        return
    }

    expected := "users/" + userID + "/avatar"
    if body.PublicID != expected {
        c.JSON(400, structs.ErrorResponse{
            Message: "Invalid public_id",
        })
        return
    }

    if body.Version <= 0 {
        c.JSON(400, structs.ErrorResponse {
            Message: "Invalid version",
        })
        return
    }

    if err := database.SaveUserPhoto(c, internal.DB, userID, body.PublicID, body.Version); err != nil {
        log.Println(err)
        c.JSON(500, structs.ErrorResponse {
            Message: "Failed to save photo",
        })
        return
    }

    c.Status(204)
}
