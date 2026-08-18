package client_test

import (
	"context"
	"testing"

	ent "github.com/keep/sunny/ent"
	user "github.com/keep/sunny/internal/modules/user/client"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type test struct {
	name string
	mock func()
	res  interface{}
	err  error
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := NewMockUserRepository(ctrl)
	userService := user.NewUserService(repository)

	ctx := context.Background()
	inputDto := user.UserCreateOneDto{
		Name:  "Alice",
		Email: "alice@example.com",
	}
	mockUser := &ent.User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
	}

	tests := []test{
		{
			name: "create one empty result",
			mock: func() {
				repository.EXPECT().CreateOne(ctx, inputDto).Return(mockUser, nil)
			},
			res: mockUser,
			err: nil,
		},
	}

	for _, tc := range tests { //nolint:paralleltest // data races here
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			res, err := userService.CreateOne(ctx, inputDto)

			require.Equal(t, res.ID, mockUser.ID)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}
