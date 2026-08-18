package auth

import (
	"time"

	"github.com/keep/sunny/ent"
)

type TokenIssuer interface {
	Issue(user *ent.User) (SessionTokens, error)
}

type SessionTokens struct {
	AccessToken      string
	AccessExpiresAt  time.Time
	RefreshToken     string
	RefreshExpiresAt time.Time
}
