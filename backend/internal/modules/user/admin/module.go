package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/keep/sunny/internal/modules/rbac"
	"go.uber.org/fx"
)

var Module = fx.Module("user-admin", fx.Provide(NewRepository, NewService, NewController), fx.Invoke(RegisterRoutes))

func RegisterRoutes(router *gin.RouterGroup, ctl *Controller, access *rbac.Controller) {
	group := router.Group("/admin/users")
	group.Use(access.Require("system:rbac"))
	group.GET("", ctl.List)
	group.POST("", ctl.Create)
	group.GET("/:id", ctl.Get)
	group.PUT("/:id", ctl.Update)
	group.PUT("/:id/status", ctl.SetStatus)
	group.PUT("/:id/roles", ctl.SetRoles)
	group.POST("/:id/reset-password", ctl.ResetPassword)
}
