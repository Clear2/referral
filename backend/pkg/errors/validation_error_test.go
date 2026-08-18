package errors

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestValidationMessage(t *testing.T) {
	type request struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=12,max=32"`
	}

	err := validator.New().Struct(request{
		Email:    "not-an-email",
		Password: "short",
	})

	require.Equal(t, "请输入有效的邮箱地址；密码至少需要12位", ValidationMessage(err))
}

func TestWrapAPIErrorDoesNotMutateSharedError(t *testing.T) {
	wrapped := WrapAPIError(ErrInternalServerError, errors.New("database detail"))

	require.Equal(t, "服务器内部错误", wrapped.Message)
	require.Error(t, wrapped.Err)
	require.NoError(t, ErrInternalServerError.Err)
}
