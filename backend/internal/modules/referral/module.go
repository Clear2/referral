package referral

import (
	"github.com/keep/sunny/internal/modules/referral/admin"
	"github.com/keep/sunny/internal/modules/referral/client"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"referral",
	client.Module,
	admin.Module,
)
