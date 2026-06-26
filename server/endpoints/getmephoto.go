package endpoints

import (
	"context"
	"log"
	"matchme-server/database"
	"matchme-server/helpers"
	"matchme-server/internal"
	"matchme-server/structs"
	"github.com/gin-gonic/gin"
)

type UserById struct {
	UserID    string `json:"userId"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatarurl"`
}
//GET users/{id}
func GetNameAndPhoto(c *gin.Context) {
    serveNameAndPhoto(c, c.Param("id"))
}
//GET me
func GetMeNameAndPhoto(c *gin.Context) {
    id := c.GetString("userID") // guaranteed by AuthRequired()
    serveNameAndPhoto(c, id)
}

// GET /me/email - returns only the authenticated user's email
func GetMyEmail(c *gin.Context) {
	userID := c.GetString("userID")
	
	ctx := context.Background()
	
	var email string
	err := internal.DB.QueryRow(ctx, 
		"SELECT email FROM users WHERE id = $1", 
		userID,
	).Scan(&email)
	
	if err != nil {
		log.Println(err)
		c.JSON(500, structs.ErrorResponse{
			Message: "failed to get email",
		})
		return
	}
	
	c.JSON(200, gin.H{
		"email": email,
	})
}

func serveNameAndPhoto(c *gin.Context, id string) {
	var out UserById
	if !helpers.IsValidID(id) {
		c.JSON(400, structs.ErrorResponse{
			Message: "invalid id",
		})
		return
	}

	ctx := context.Background()
	err, Exists := database.EnsureUserExists(ctx, internal.DB, id)
	if err != nil {
		log.Println(err)
		c.JSON(500, structs.ErrorResponse{
			Message: "db error",
		})
		return
	}

	if !Exists {
	c.JSON(404, structs.ErrorResponse{  // Было 400, стало 404
		Message: "user does not exist",
	})
	return
}

	out.Name, out.AvatarUrl = database.GetUserNamePhotoURL(ctx, internal.DB, id)
	out.UserID = id
	
	c.JSON(200, out)

}

