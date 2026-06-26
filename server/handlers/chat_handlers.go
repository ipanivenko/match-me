package handlers

import (
	"log"
	"matchme-server/database"
	"matchme-server/internal"
	ws "matchme-server/websocket"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUserChats returns all chats for the authenticated user
func GetUserChats(c *gin.Context) {
	userID := c.GetString("userID")
	
	chats, err := database.GetUserChats(c.Request.Context(), internal.DB, userID)
	
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load chats"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"chats": chats})
}

// GetChatMessages returns messages for a specific chat
func GetChatMessages(c *gin.Context) {
	userID := c.GetString("userID")
	chatID := c.Param("chatId")
	before := c.DefaultQuery("before", "")
	limit := 50

	messages, err := database.GetChatMessages(c.Request.Context(), internal.DB, chatID, userID, before, limit)
	if err != nil {
		log.Printf("ERROR getting messages: %v", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied or chat not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// SendMessage sends a new message in a chat
func SendMessage(c *gin.Context) {
	userID := c.GetString("userID")
	chatID := c.Param("chatId")

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		return
	}

	messageID, err := database.SaveChatMessage(c.Request.Context(), internal.DB, chatID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed to send message"})
		return
	}

	createdAt := time.Now()

	// Get the other user's ID
	chat, _ := database.GetChatByID(c.Request.Context(), internal.DB, chatID, userID)
	var recipientID string
	if chat != nil {
		if chat.User1ID == userID {
			recipientID = chat.User2ID
		} else {
			recipientID = chat.User1ID
		}

		if GlobalHub != nil {
			wsMessage := &ws.Message{
				Type:      "new_message",
				ChatID:    chatID,
				MessageID: messageID,
				SenderID:  userID,
				Content:   req.Content,
				CreatedAt: createdAt.Format(time.RFC3339),
			}
			

			GlobalHub.SendToUser(recipientID, wsMessage)
			
			GlobalHub.SendToUser(userID, wsMessage)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": gin.H{
			"id":         messageID,
			"sender_id":  userID,
			"content":    req.Content,
			"created_at": createdAt.Format(time.RFC3339),
		},
	})
}

// MarkMessagesAsRead marks all messages in a chat as read
func MarkMessagesAsRead(c *gin.Context) {
	userID := c.GetString("userID")
	chatID := c.Param("chatId")
	
	// Verify user has access to chat
	_, err := database.GetChatByID(c.Request.Context(), internal.DB, chatID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	
	// Mark messages as read (delete from unread_messages table)
	err = database.MarkChatMessagesAsRead(c.Request.Context(), internal.DB, chatID, userID)
	if err != nil {
		log.Printf("Error marking messages as read: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as read"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// CheckOnlineStatus checks if a user is currently online
func CheckOnlineStatus(c *gin.Context) {
	targetUserID := c.Param("id")
	
	isOnline := false
	if GlobalHub != nil {
		isOnline = GlobalHub.IsOnline(targetUserID)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"user_id": targetUserID,
		"online": isOnline,
	})
}