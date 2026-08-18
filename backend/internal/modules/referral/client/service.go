package client

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/keep/sunny/ent"
	appErrors "github.com/keep/sunny/pkg/errors"
)

const codeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

type Service interface {
	GenerateInvitation(context.Context, int) (string, error)
	Register(context.Context, RegisterInput) (*RegistrationView, error)
	Dashboard(context.Context, int) (*DashboardView, error)
}
type service struct{ repository Repository }

func NewService(repository Repository) Service { return &service{repository: repository} }

func (s *service) GenerateInvitation(ctx context.Context, userID int) (string, error) {
	entity, err := s.repository.GetUser(ctx, userID)
	if err != nil {
		return "", publicError(err)
	}
	if entity.ReferralCode != nil {
		return *entity.ReferralCode, nil
	}
	for range 5 {
		code, codeErr := generateCode(8)
		if codeErr != nil {
			return "", appErrors.WrapAPIError(appErrors.ErrInternalServerError, codeErr)
		}
		entity, err = s.repository.SetReferralCode(ctx, userID, code)
		if err == nil {
			return *entity.ReferralCode, nil
		}
		if !ent.IsConstraintError(err) {
			return "", publicError(err)
		}
	}
	return "", appErrors.NewAPIError(http.StatusServiceUnavailable, "邀请码生成失败，请重试")
}
func (s *service) Register(ctx context.Context, input RegisterInput) (*RegistrationView, error) {
	input.Code = strings.ToUpper(strings.TrimSpace(input.Code))
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	invitee, inviter, err := s.repository.Register(ctx, input)
	if err != nil {
		return nil, publicError(err)
	}
	return &RegistrationView{Invitee: toUserView(invitee), Reward: RewardCredits, InviterCreditBalance: inviter.CreditBalance}, nil
}
func (s *service) Dashboard(ctx context.Context, userID int) (*DashboardView, error) {
	entity, err := s.repository.Dashboard(ctx, userID)
	if err != nil {
		return nil, publicError(err)
	}
	out := &DashboardView{User: toUserView(entity), SuccessfulReferrals: len(entity.Edges.SentReferrals), Referrals: make([]ReferralView, 0, len(entity.Edges.SentReferrals)), CreditTransactions: make([]CreditTransactionView, 0, len(entity.Edges.CreditTransactions))}
	for _, item := range entity.Edges.SentReferrals {
		out.Referrals = append(out.Referrals, ReferralView{ID: item.ID, Invitee: toUserView(item.Edges.Invitee), Reward: item.Reward, CreatedAt: item.CreateTime})
	}
	for _, item := range entity.Edges.CreditTransactions {
		if item.Reason == "referral_reward" && item.Amount > 0 {
			out.TotalCreditsEarned += item.Amount
		}
		out.CreditTransactions = append(out.CreditTransactions, CreditTransactionView{ID: item.ID, Amount: item.Amount, Reason: item.Reason, ReferralID: item.ReferralID, CreatedAt: item.CreateTime})
	}
	return out, nil
}
func generateCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	for i := range bytes {
		bytes[i] = codeAlphabet[int(bytes[i])%len(codeAlphabet)]
	}
	return string(bytes), nil
}
func toUserView(entity *ent.User) UserView {
	return UserView{ID: entity.ID, Name: entity.Name, Email: entity.Email, ReferralCode: entity.ReferralCode, CreditBalance: entity.CreditBalance}
}
func publicError(err error) error {
	switch {
	case errors.Is(err, ErrInvitationNotFound):
		return appErrors.NewAPIError(http.StatusNotFound, "邀请码不存在")
	case errors.Is(err, ErrUserNotFound):
		return appErrors.NewAPIError(http.StatusNotFound, "用户不存在")
	case errors.Is(err, ErrEmailExists):
		return appErrors.NewAPIError(http.StatusConflict, "邮箱已注册")
	default:
		return appErrors.WrapAPIError(appErrors.ErrInternalServerError, err)
	}
}
