package client

import (
	"context"

	"github.com/keep/sunny/ent"
	"github.com/keep/sunny/ent/user"
	"github.com/keep/sunny/pkg/postgres"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

//go:generate mockgen -source=repository.go -destination=./mocks_user_repo_test.go -package=client_test

var repTracer = otel.Tracer("UserRepository")

// UserRepository -.
type UserRepository interface {
	CreateOne(ctx context.Context, dto UserCreateOneDto) (*ent.User, error)
	GetByID(ctx context.Context, id int) (*ent.User, error)
}

type userRepositoryImpl struct {
	pg *postgres.Postgres
}

func NewUserRepository(pg *postgres.Postgres) UserRepository {
	return &userRepositoryImpl{pg}
}

// CreateOne -.
func (u *userRepositoryImpl) CreateOne(
	ctx context.Context,
	dto UserCreateOneDto,
) (*ent.User, error) {
	return u.pg.Client.User.
		Create().
		SetName(dto.Name).
		SetEmail(dto.Email).
		Save(ctx)
}

// GetByID returns a user by ID.
func (u *userRepositoryImpl) GetByID(
	ctx context.Context,
	id int,
) (*ent.User, error) {
	ctx, span := repTracer.Start(ctx, "UserRepository.GetByID")
	defer span.End()

	span.SetAttributes(
		attribute.Bool("has_user_id", true),
	)

	user, err := u.pg.Client.User.
		Query().
		Where(
			user.IDEQ(id),
		).
		Only(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return user, nil
}
