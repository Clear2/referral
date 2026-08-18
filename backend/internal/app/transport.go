package app

import (
	"github.com/gin-gonic/gin"
	"github.com/keep/sunny/config"
	"github.com/keep/sunny/internal/router"
	apiRouter "github.com/keep/sunny/internal/router/api"
	"github.com/keep/sunny/pkg/httpserver"
	"go.uber.org/fx"
)

var transportModule = fx.Module("transport",
	fx.Provide(
		fx.Annotate(router.NewRouter, fx.ParamTags("", "", `name:"tracing_filter"`)),
		apiRouter.NewAPIRouter,
		func(cfg config.Config, handler *gin.Engine) httpserver.Server {
			return httpserver.New(handler, httpserver.Port(cfg.HTTP().Port))
		},
	),
)
