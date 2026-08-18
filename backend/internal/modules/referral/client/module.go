package client

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/keep/sunny/internal/modules/auth"
	"github.com/keep/sunny/pkg/errors"
	"github.com/keep/sunny/pkg/utils"
	"go.uber.org/fx"
)

var Module = fx.Module("referral-client", fx.Provide(NewRepository, NewService, NewController), fx.Invoke(RegisterRoutes))

func RegisterRoutes(router *gin.RouterGroup, controller *Controller, authManager *auth.Manager) {
	owner := func(c *gin.Context) {
		current, err := utils.MustGetUser(c)
		requested, parseErr := strconv.Atoi(c.Param("id"))
		if err != nil || parseErr != nil || current != requested {
			_ = c.AbortWithError(http.StatusForbidden, errors.ErrForbidden)
			return
		}
		c.Next()
	}
	router.POST("/users/:id/referral-code", authManager.Require(), owner, controller.GenerateInvitation)
	router.POST("/referrals/register", controller.Register)
	router.GET("/users/:id/referral-dashboard", authManager.Require(), owner, controller.Dashboard)
}
