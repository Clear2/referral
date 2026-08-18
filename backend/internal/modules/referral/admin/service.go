package admin

import (
	"context"
	"net/http"
	"strings"

	appErrors "github.com/keep/sunny/pkg/errors"
)

type Service struct{ repository *Repository }

func NewService(repository *Repository) *Service { return &Service{repository: repository} }

func normalize(input ListInput) ListInput {
	if input.Page == 0 {
		input.Page = 1
	}
	if input.PageSize == 0 {
		input.PageSize = 20
	}
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	return input
}
func pageView(page, size, total int) PaginationView {
	pages := 0
	if total > 0 {
		pages = (total + size - 1) / size
	}
	return PaginationView{Page: page, PageSize: size, Total: total, TotalPages: pages}
}
func validate(input ListInput) error {
	if input.CreatedAtFrom != nil && input.CreatedAtTo != nil && input.CreatedAtFrom.After(*input.CreatedAtTo) {
		return appErrors.NewAPIError(http.StatusBadRequest, "开始日期不能晚于结束日期")
	}
	return nil
}
func internal(err error) error {
	if err == nil {
		return nil
	}
	return appErrors.WrapAPIError(appErrors.ErrInternalServerError, err)
}
func (s *Service) Referrals(ctx context.Context, input ListInput) (*ReferralListView, error) {
	input = normalize(input)
	if err := validate(input); err != nil {
		return nil, err
	}
	items, total, err := s.repository.Referrals(ctx, input)
	if err != nil {
		return nil, internal(err)
	}
	out := &ReferralListView{Items: make([]ReferralView, 0, len(items)), Pagination: pageView(input.Page, input.PageSize, total)}
	for _, item := range items {
		out.Items = append(out.Items, ReferralView{ID: item.ID, Inviter: toUserView(item.Edges.Inviter), Invitee: toUserView(item.Edges.Invitee), Reward: item.Reward, CreatedAt: item.CreateTime})
	}
	return out, nil
}
func (s *Service) CreditTransactions(ctx context.Context, input ListInput) (*CreditTransactionListView, error) {
	input = normalize(input)
	if err := validate(input); err != nil {
		return nil, err
	}
	items, total, err := s.repository.CreditTransactions(ctx, input)
	if err != nil {
		return nil, internal(err)
	}
	out := &CreditTransactionListView{Items: make([]CreditTransactionView, 0, len(items)), Pagination: pageView(input.Page, input.PageSize, total)}
	for _, item := range items {
		out.Items = append(out.Items, CreditTransactionView{ID: item.ID, User: toUserView(item.Edges.User), Amount: item.Amount, Reason: item.Reason, ReferralID: item.ReferralID, CreatedAt: item.CreateTime})
	}
	return out, nil
}
func (s *Service) Stats(ctx context.Context) (*StatsView, error) {
	out, err := s.repository.Stats(ctx)
	if err != nil {
		return nil, internal(err)
	}
	return &out, nil
}
