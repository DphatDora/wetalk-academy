package middleware

import (
	"errors"
	"net/http"
	"strings"
	"wetalk-academy/config"
	"wetalk-academy/internal/interface/dto/response"
	"wetalk-academy/package/logger"
	"wetalk-academy/package/util"

	"github.com/gin-gonic/gin"
)

// resolveToken parse token from header and inject into context
func resolveToken(c *gin.Context, conf *config.Config) error {
	newCtx := logger.ContextWithClientIP(c.Request.Context(), c.ClientIP())

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.Request = c.Request.WithContext(newCtx)
		return nil
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.Request = c.Request.WithContext(newCtx)
		return errors.New("invalid Authorization header format, expected 'Bearer <token>'")
	}

	tokenString := parts[1]

	claims, err := util.VerifyJWT(tokenString, conf.Auth.JWTSecret)
	if err != nil {
		c.Request = c.Request.WithContext(newCtx)
		return err
	}

	c.Set("userID", claims.UserID)

	newCtx = logger.ContextWithUserID(newCtx, claims.UserID)
	newCtx = logger.ContextWithToken(newCtx, tokenString)
	c.Request = c.Request.WithContext(newCtx)

	return nil
}

func AuthMiddleware(conf *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := resolveToken(c, conf); err != nil {
			logger.Errorf("[Err] AuthMiddleware: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.APIResponse{
				Success: false,
				Message: "Unauthorized: " + err.Error(),
			})
			return
		}
		c.Next()
	}
}

func OptionalAuthMiddleware(conf *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := resolveToken(c, conf); err != nil {
			// Warning if has token but invalid
			if c.GetHeader("Authorization") != "" {
				logger.Warnf("[Warn] OptionalAuth invalid token: %v", err)
			}
		}
		c.Next()
	}
}
