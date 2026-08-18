package auth

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/keep/sunny/config"
	"github.com/keep/sunny/internal/modules/rbac"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"auth",
	fx.Provide(
		NewManager,
		NewGoogleProvider,
		NewRepository,
		func(manager *Manager) TokenIssuer { return manager },
		NewService,
		NewController,
	),
	fx.Invoke(func(router *gin.RouterGroup, manager *Manager, googleController *Controller, access *rbac.Controller, lifecycle fx.Lifecycle, cfg config.Config) {
		router.POST("/auth/login-with-account", manager.login)
		router.POST("/auth/registration-code", manager.registrationCode)
		router.POST("/auth/register", manager.register)
		router.POST("/auth/refresh", manager.refresh)
		router.POST("/auth/logout", manager.logout)
		router.GET("/admin/session", manager.Require(), access.RequireAny("referral:read", "system:rbac"), manager.session)
		router.GET("/auth/google/login", googleController.GoogleLogin)
		router.GET("/auth/google/callback", googleController.GoogleCallback)
		lifecycle.Append(fx.Hook{OnStart: func(ctx context.Context) error { return manager.bootstrap(ctx, cfg) }})
	}),
)
