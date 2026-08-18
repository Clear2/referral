package api_router

import (
	"github.com/gin-gonic/gin"
)

func NewAPIRouter(app *gin.Engine) *gin.RouterGroup {
	return app.Group("/api/v1")
}
