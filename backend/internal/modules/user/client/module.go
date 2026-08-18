package client

import (
	"github.com/gin-gonic/gin"
	"github.com/keep/sunny/internal/modules/auth"
	"go.uber.org/fx"
)

var Module = fx.Module("user-client", fx.Provide(NewUserRepository, NewUserService, NewUserController), fx.Invoke(RegisterRoutes))

func RegisterRoutes(router *gin.RouterGroup, controller *UserController, authManager *auth.Manager) {
	group := router.Group("/users")
	group.GET("/me", authManager.Require(), controller.GetMe)
}
