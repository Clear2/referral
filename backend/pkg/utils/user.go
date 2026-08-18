package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keep/sunny/pkg/errors"
)

// MustGetUser retrieves the current request's userId from the Gin context.
// Returns ErrSessionUnauthorized if the userId does not exist or is not an int.
func MustGetUser(c *gin.Context) (int, error) {
	val, exists := c.Get("userId")
	if !exists {
		return 0, errors.ErrSessionUnauthorized
	}

	id, ok := val.(int)
	if !ok {
		return 0, errors.ErrSessionUnauthorized
	}

	return id, nil
}

// HandlerWithUser is a handler wrapper that injects userId into the handler.
// The handler signature must be func(c *gin.Context, userId int).
// If userId is missing or invalid, the request is aborted with a 401 Unauthorized response.
func HandlerWithUser(handler func(c *gin.Context, userId int)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := MustGetUser(ctx)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		handler(ctx, userId)
	}
}
