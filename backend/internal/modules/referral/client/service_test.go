package client

import (
	"context"
	"errors"
	"testing"

	"github.com/keep/sunny/ent"
	appErrors "github.com/keep/sunny/pkg/errors"
	"github.com/stretchr/testify/require"
)

type fakeRepository struct {
	user          *ent.User
	invitee       *ent.User
	inviter       *ent.User
	registerInput RegisterInput
	err           error
}

func (repository *fakeRepository) GetUser(context.Context, int) (*ent.User, error) {
	return repository.user, repository.err
}

func (repository *fakeRepository) SetReferralCode(_ context.Context, _ int, code string) (*ent.User, error) {
	repository.user.ReferralCode = &code
	return repository.user, repository.err
}

func (repository *fakeRepository) Register(_ context.Context, input RegisterInput) (*ent.User, *ent.User, error) {
	repository.registerInput = input
	return repository.invitee, repository.inviter, repository.err
}

func (repository *fakeRepository) Dashboard(context.Context, int) (*ent.User, error) {
	return repository.user, repository.err
}

func TestGenerateInvitationReturnsExistingCode(t *testing.T) {
	t.Parallel()
	code := "ABC12345"
	repository := &fakeRepository{user: &ent.User{ID: 1, ReferralCode: &code}}

	actual, err := NewService(repository).GenerateInvitation(context.Background(), 1)

	require.NoError(t, err)
	require.Equal(t, code, actual)
}

func TestRegisterNormalizesInputAndRewardsInviter(t *testing.T) {
	t.Parallel()
	repository := &fakeRepository{
		invitee: &ent.User{ID: 2, Name: "Bob", Email: "bob@example.com"},
		inviter: &ent.User{ID: 1, CreditBalance: RewardCredits},
	}

	result, err := NewService(repository).Register(context.Background(), RegisterInput{
		Code: " abcd2345 ", Name: " Bob ", Email: " BOB@Example.com ",
	})

	require.NoError(t, err)
	require.Equal(t, "ABCD2345", repository.registerInput.Code)
	require.Equal(t, "Bob", repository.registerInput.Name)
	require.Equal(t, "bob@example.com", repository.registerInput.Email)
	require.Equal(t, RewardCredits, result.Reward)
	require.Equal(t, RewardCredits, result.InviterCreditBalance)
}

func TestRegisterMapsDuplicateEmailToConflict(t *testing.T) {
	t.Parallel()
	repository := &fakeRepository{err: ErrEmailExists}

	_, err := NewService(repository).Register(context.Background(), RegisterInput{})

	var apiErr *appErrors.APIError
	require.True(t, errors.As(err, &apiErr))
	require.Equal(t, 409, apiErr.Code)
}

func TestGenerateCode(t *testing.T) {
	t.Parallel()
	code, err := generateCode(8)
	require.NoError(t, err)
	require.Len(t, code, 8)
	for _, character := range code {
		require.Contains(t, codeAlphabet, string(character))
	}
}
