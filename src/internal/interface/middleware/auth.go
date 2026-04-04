package middleware

import (
	"net/http"
	"strings"
	"wetalk-academy/config"
	"wetalk-academy/internal/interface/dto/response"
	"wetalk-academy/package/logger"
	"wetalk-academy/package/util"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(conf *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Errorf("[Err] Missing Authorization header")
			c.JSON(http.StatusUnauthorized, response.APIResponse{
				Success: false,
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Errorf("[Err] Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, response.APIResponse{
				Success: false,
				Message: "Invalid authorization header format. Expected 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Verify JWT token
		claims, err := util.VerifyJWT(tokenString, conf.Auth.JWTSecret)
		if err != nil {
			logger.Errorf("[Err] Invalid JWT token: %v", err)
			c.JSON(http.StatusUnauthorized, response.APIResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("userID", claims.UserID)

		c.Next()
	}
}

func OptionalAuthMiddleware(conf *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		// Verify JWT token
		claims, err := util.VerifyJWT(tokenString, conf.Auth.JWTSecret)
		if err != nil {
			logger.Warnf("[Warn] Invalid JWT token in optional auth: %v", err)
			c.Next()
			return
		}

		c.Set("userID", claims.UserID)

		c.Next()
	}
}
