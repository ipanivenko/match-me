package services

import (
	"log"
	"matchme-server/database"
	"matchme-server/helpers"
	"matchme-server/internal"
	"matchme-server/structs"

	"github.com/gin-gonic/gin"
)

var CommonErr = "Something went wrong, please, try again later."

func Register(c *gin.Context) {
	var input structs.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err)
		c.JSON(400, structs.ErrorResponse{
			Message: CommonErr,
		})
		return
	}

	if !helpers.IsValidEmail(input.Email) {
		c.JSON(400, structs.ErrorResponse{
			Message: "Email is invalid, please provide a valid email.",
		})
		return
	}

	if !helpers.IsUniqEmail(input.Email) {
		c.JSON(400, structs.ErrorResponse{
			Message: "Email already exists.",
		})
		return
	}

	if !helpers.IsValidPassword(input.Password) {
		log.Println("invalid password")
		c.JSON(400, structs.ErrorResponse{
			Message: "Password should be at least 6 symbols, one letter.",
		})
		return
	}

	hashedPassword, err := helpers.HashPassword(input.Password)
	if err != nil {
		log.Println(err)
		c.JSON(500, structs.ErrorResponse{			
			Message: CommonErr,
		})
		return
	}

	id, createdAt, err := database.CreateUser(c.Request.Context(), internal.DB, input.Email, hashedPassword)
	if err != nil {
		log.Println(err)
		c.JSON(500, structs.ErrorResponse{
			Message:  CommonErr,
		})
		return
	}


	c.JSON(201, gin.H{"id": id, "email": input.Email, "created_at": createdAt})
}

