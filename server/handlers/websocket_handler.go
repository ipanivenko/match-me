package handlers

import (
	"log"
	"matchme-server/internal"
	ws "matchme-server/websocket"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

func HandleWebSocket(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from query parameter
		tokenString := c.Query("token")
		
		// If not in query, try Authorization header
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			return
		}

		// Verify JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(internal.Cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			log.Printf("Invalid token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Extract user ID from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}

		// Upgrade to WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		client := ws.NewClient(hub, conn, userID)
		hub.Register <- client

		go client.WritePump()
		go client.ReadPump()
	}
}

// Global hub instance
var GlobalHub *ws.Hub

func init() {
	GlobalHub = ws.NewHub()
	go GlobalHub.Run()
}