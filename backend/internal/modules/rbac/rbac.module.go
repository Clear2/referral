package rbac

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Module("rbac",
	fx.Provide(NewRepository, NewService, NewController),
	fx.Invoke(RegisterRoutes),
)

func RegisterRoutes(router *gin.RouterGroup, ctl *Controller) {
	router.GET("/access/me", ctl.MyAccess)
	group := router.Group("/admin/rbac")
	group.Use(ctl.Require("system:rbac"))
	group.GET("", ctl.Snapshot)
	group.POST("/roles", ctl.CreateRole)
	group.PUT("/roles/:id", ctl.UpdateRole)
	group.DELETE("/roles/:id", ctl.DeleteRole)
	group.PUT("/roles/:id/grants", ctl.SetGrants)
	group.POST("/permissions", ctl.CreatePermission)
	group.PUT("/permissions/:id", ctl.UpdatePermission)
	group.PUT("/permissions/:id/resources", ctl.SetPermissionResources)
	group.DELETE("/permissions/:id", ctl.DeletePermission)
	group.POST("/menus", ctl.CreateMenu)
	group.PUT("/menus/:id", ctl.UpdateMenu)
	group.DELETE("/menus/:id", ctl.DeleteMenu)
	group.POST("/groups", ctl.CreateGroup)
	group.PUT("/groups/:id", ctl.UpdateGroup)
	group.DELETE("/groups/:id", ctl.DeleteGroup)
	group.POST("/apis", ctl.CreateAPI)
	group.POST("/apis/sync", ctl.SyncAPIs)
	group.PUT("/apis/:id", ctl.UpdateAPI)
	group.DELETE("/apis/:id", ctl.DeleteAPI)
}
