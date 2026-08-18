package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/keep/sunny/ent"
	"github.com/keep/sunny/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var srvTracer = otel.Tracer("UserService")

type UserService interface {
	CreateOne(ctx context.Context, dto UserCreateOneDto) (*ent.UserEntity, error)
	GetByID(ctx context.Context, id int) (*ent.UserEntity, error)
}

type userServiceImpl struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) UserService {
	return &userServiceImpl{userRepository: userRepository}
}

// CreateUser -.
func (us *userServiceImpl) CreateOne(ctx context.Context, dto UserCreateOneDto) (*ent.UserEntity, error) {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.Email = strings.ToLower(strings.TrimSpace(dto.Email))
	user, err := us.userRepository.CreateOne(ctx, dto)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, errors.NewAPIError(http.StatusConflict, "邮箱已注册")
		}
		return nil, errors.WrapAPIError(
			errors.ErrInternalServerError,
			errors.NewRepositoryError(
				err.Error(),
				err,
			),
		)
	}

	return user.IntoEntity(), nil
}

// GetByID returns a user by ID.
func (us *userServiceImpl) GetByID(ctx context.Context, id int) (*ent.UserEntity, error) {
	ctx, span := srvTracer.Start(ctx, "UserService.GetByID")
	defer span.End()

	span.SetAttributes(
		attribute.Bool("has_user_id", true),
	)

	user, err := us.userRepository.GetByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			span.SetAttributes(
				attribute.Bool("user.not_found", true),
			)
			return nil, errors.NewAPIError(
				http.StatusNotFound,
				fmt.Sprintf("user %d not found", id),
			)
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapAPIError(
			errors.ErrInternalServerError,
			errors.NewRepositoryError(
				err.Error(),
				err,
			),
		)
	}

	return user.IntoEntity(), nil
}
