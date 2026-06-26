package services

import (
	"errors"
	"log"
	"matchme-server/database"
	"matchme-server/helpers"
	"matchme-server/internal"
	"matchme-server/structs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// POST /login
func Login(c *gin.Context) {
	var input structs.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, structs.ErrorResponse{
			Message: CommonErr})
		return
	}

	id, pwHash, err := database.GetUserByEmail(c.Request.Context(), internal.DB, input.Email)

	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(401, structs.ErrorResponse{
			Message: "user does not exist",
		})
		return
	}
	
	if err != nil {
		c.JSON(500, structs.ErrorResponse{
			Message: CommonErr,
		})
		return
	}

	if !helpers.IsCorrectPassword(pwHash, input.Password) {
		c.JSON(401, structs.ErrorResponse{
			Message: "invalid credentials",
		})
		return
	}

	access, err := helpers.MakeAccessToken(id)
	if err != nil {
		log.Println(err)
		c.JSON(500, structs.ErrorResponse{
			Message: CommonErr,
		})
		return
	}
	c.JSON(200, gin.H{
		"user_id":      id,
		"access_token": access,
	})
}
