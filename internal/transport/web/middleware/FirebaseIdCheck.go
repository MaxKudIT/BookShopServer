package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

var authClient *auth.Client

func InitFirebaseAuthWithProject(projectID string) error {
	ctx := context.Background()

	opt := option.WithCredentialsFile("../../internal/service-account.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return err
	}

	authClient, err = app.Auth(ctx)
	return err
}

func VerifyTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if authClient == nil {
			c.AbortWithStatusJSON(500, gin.H{
				"error":   "Firebase auth client not initialized",
				"details": "Call InitFirebaseAuthWithProject() first",
			})
			return
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})

			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization format. Use: Bearer <token>"})
			return
		}

		idToken := parts[1]

		token, err := authClient.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("firebase_id", token.UID)
		c.Next()
	}
}
