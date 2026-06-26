package graphsetup

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"matchme-server/graph"
	"matchme-server/internal"
	"matchme-server/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vektah/gqlparser/v2/ast"
)

func RegisterGraphQL(router gin.IRouter, IsDevMode bool, db *pgxpool.Pool) {

	resolver := &graph.Resolver{
		DB: db,
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	ws := &transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			Subprotocols:    []string{"graphql-transport-ws", "graphql-ws"},
		},

		KeepAlivePingInterval: 25 * time.Second,

		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, *transport.InitPayload, error) {
			var tokenString string
			if auth, ok := initPayload["Authorization"].(string); ok && auth != "" {
				tokenString = strings.TrimPrefix(auth, "Bearer ")
			} else if auth, ok := initPayload["authorization"].(string); ok && auth != "" {
				tokenString = strings.TrimPrefix(auth, "Bearer ")
			} else {
				return nil, nil, fmt.Errorf("missing Authorization in connection payload")
			}

			// Parse and verify JWT
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				return []byte(internal.Cfg.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				return nil, nil, fmt.Errorf("invalid token")
			}

			// Extract user ID from claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, nil, fmt.Errorf("invalid token claims")
			}

			userID, _ := claims["sub"].(string)
			if userID == "" {
				return nil, nil, fmt.Errorf("missing user ID in token")
			}

			newCtx := context.WithValue(ctx, middleware.UserIDKey, userID)
			log.Printf("✅ WebSocket connected for user %s\n", userID)

			return newCtx, &initPayload, nil
		},
	}
	srv.AddTransport(ws)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	playgroundHandler := playground.Handler("GraphQL playground", "/graphql")

	router.POST("/graphql", middleware.GinGqlAuthMiddleware(internal.Cfg.JWTSecret),
	ginAdapter(srv))
	router.GET("/graphql", ginAdapter(srv))
	router.GET("/graphql-ws", ginAdapter(srv))

	if IsDevMode {
		log.Println("Developer mode enabled. GraphQL Playground is available at /playground")
		router.GET("/playground", ginAdapter(playgroundHandler))
	}
}

// ginAdapter converts a standard http.Handler to a gin.HandlerFunc
func ginAdapter(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
