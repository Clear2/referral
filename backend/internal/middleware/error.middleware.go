package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	appErrors "github.com/keep/sunny/pkg/errors"
	"github.com/keep/sunny/pkg/logger"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ErrorHandler captures errors and returns a consistent JSON error response.
func ErrorHandler(logger logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		err := ctx.Errors.Last().Err
		span := trace.SpanFromContext(ctx.Request.Context())

		// API errors carry the public status and message. Check them before
		// unwrapping repository errors so internal database details never leak.
		var apiErr *appErrors.APIError
		if errors.As(err, &apiErr) {
			if apiErr.Code >= 500 {
				span.RecordError(err)
				span.SetStatus(codes.Error, apiErr.Message)
				logger.Errorw(
					"API request failed",
					"error", err.Error(),
					"request_id", requestid.Get(ctx),
					"method", ctx.Request.Method,
					"path", ctx.Request.URL.Path,
				)
			}

			ctx.JSON(
				apiErr.Code,
				gin.H{
					"code":    apiErr.Code,
					"message": apiErr.Message,
				},
			)
			return
		}

		var repErr *appErrors.RepositoryError
		if errors.As(err, &repErr) {
			span.RecordError(err)
			span.SetStatus(codes.Error, repErr.Message)

			logger.Errorw(
				"log from middleware error handler",
				"error", repErr.Message,
				"request_id", requestid.Get(ctx),
				"method", ctx.Request.Method,
				"path", ctx.Request.URL.Path,
				"route", ctx.FullPath(),
				"handler", ctx.HandlerName(),
			)

			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"code":    http.StatusInternalServerError,
					"message": appErrors.ErrInternalServerError.Message,
				},
			)
			return
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Errorw(
			"unhandled API request error",
			"error", err.Error(),
			"request_id", requestid.Get(ctx),
			"method", ctx.Request.Method,
			"path", ctx.Request.URL.Path,
			"route", ctx.FullPath(),
			"handler", ctx.HandlerName(),
		)

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code":    http.StatusInternalServerError,
				"message": appErrors.ErrInternalServerError.Message,
			},
		)
	}
}
