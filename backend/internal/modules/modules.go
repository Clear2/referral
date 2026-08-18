package modules

import (
	"github.com/keep/sunny/internal/modules/auth"
	"github.com/keep/sunny/internal/modules/rbac"
	"github.com/keep/sunny/internal/modules/referral"
	"github.com/keep/sunny/internal/modules/user"
	"go.uber.org/fx"
)

var APIModule = fx.Module(
	"api",

	user.Module,
	auth.Module,
	referral.Module,
	rbac.Module,
)
