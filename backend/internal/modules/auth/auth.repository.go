package auth

import (
	"context"
	"fmt"

	"github.com/keep/sunny/ent"
	"github.com/keep/sunny/ent/user"
	"github.com/keep/sunny/pkg/postgres"
)

type Repository interface {
	FindOrCreate(ctx context.Context, name, email string) (*ent.User, bool, error)
}

type repository struct{ client *ent.Client }

func NewRepository(pg *postgres.Postgres) Repository { return &repository{client: pg.Client} }

func (repository *repository) FindOrCreate(ctx context.Context, name, email string) (*ent.User, bool, error) {
	entity, err := repository.client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
	if err == nil {
		return entity, false, nil
	}
	if !ent.IsNotFound(err) {
		return nil, false, fmt.Errorf("find Google user by email: %w", err)
	}
	entity, err = repository.client.User.Create().SetName(name).SetEmail(email).Save(ctx)
	if err == nil {
		return entity, true, nil
	}
	if !ent.IsConstraintError(err) {
		return nil, false, fmt.Errorf("create Google user: %w", err)
	}
	entity, err = repository.client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("find concurrently created Google user: %w", err)
	}
	return entity, false, nil
}
