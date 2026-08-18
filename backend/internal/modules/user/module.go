package user

import (
	"github.com/keep/sunny/internal/modules/user/admin"
	"github.com/keep/sunny/internal/modules/user/client"
	"go.uber.org/fx"
)

var Module = fx.Module("user", client.Module, admin.Module)
