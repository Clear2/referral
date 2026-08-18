package observability

import (
	"github.com/keep/sunny/internal/observability/metrics"
	"github.com/keep/sunny/internal/observability/sentry"
	"github.com/keep/sunny/internal/observability/tracing"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"observability",

	// tracing
	tracing.Module,
	// metrics
	metrics.Module,
	// sentry
	sentry.Module,
)
