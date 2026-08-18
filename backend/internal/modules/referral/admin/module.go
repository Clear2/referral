package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/keep/sunny/internal/modules/rbac"
	"go.uber.org/fx"
)

var Module = fx.Module("admin-referral", fx.Provide(NewRepository, NewService, NewController), fx.Invoke(RegisterRoutes))

func RegisterRoutes(router *gin.RouterGroup, controller *Controller, access *rbac.Controller) {
	group := router.Group("/admin")
	group.Use(access.Require("referral:read"))
	group.GET("/referrals", controller.Referrals)
	group.GET("/credit-transactions", controller.CreditTransactions)
	group.GET("/referral-stats", controller.Stats)
}
