package admin

import (
	"context"
	"fmt"

	"github.com/keep/sunny/ent"
	"github.com/keep/sunny/ent/credittransaction"
	"github.com/keep/sunny/ent/predicate"
	"github.com/keep/sunny/ent/referral"
	"github.com/keep/sunny/ent/user"
	"github.com/keep/sunny/pkg/postgres"
)

type Repository struct{ client *ent.Client }

func NewRepository(pg *postgres.Postgres) *Repository { return &Repository{client: pg.Client} }
func toUserView(entity *ent.User) UserView {
	return UserView{ID: entity.ID, Name: entity.Name, Email: entity.Email, ReferralCode: entity.ReferralCode, CreditBalance: entity.CreditBalance}
}

func referralFilters(input ListInput) []predicate.Referral {
	filters := make([]predicate.Referral, 0, 4)
	if input.UserID > 0 {
		filters = append(filters, referral.InviterIDEQ(input.UserID))
	}
	if input.Email != "" {
		filters = append(filters, referral.Or(referral.HasInviterWith(user.EmailContainsFold(input.Email)), referral.HasInviteeWith(user.EmailContainsFold(input.Email))))
	}
	if input.CreatedAtFrom != nil {
		filters = append(filters, referral.CreateTimeGTE(*input.CreatedAtFrom))
	}
	if input.CreatedAtTo != nil {
		filters = append(filters, referral.CreateTimeLT(input.CreatedAtTo.AddDate(0, 0, 1)))
	}
	return filters
}
func (r *Repository) Referrals(ctx context.Context, input ListInput) ([]*ent.Referral, int, error) {
	query := r.client.Referral.Query().Where(referralFilters(input)...)
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count referrals: %w", err)
	}
	items, err := query.WithInviter().WithInvitee().Order(ent.Desc(referral.FieldCreateTime)).Offset((input.Page - 1) * input.PageSize).Limit(input.PageSize).All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("list referrals: %w", err)
	}
	return items, total, nil
}
func (r *Repository) CreditTransactions(ctx context.Context, input ListInput) ([]*ent.CreditTransaction, int, error) {
	filters := make([]predicate.CreditTransaction, 0, 4)
	if input.UserID > 0 {
		filters = append(filters, credittransaction.UserIDEQ(input.UserID))
	}
	if input.Email != "" {
		filters = append(filters, credittransaction.HasUserWith(user.EmailContainsFold(input.Email)))
	}
	if input.CreatedAtFrom != nil {
		filters = append(filters, credittransaction.CreateTimeGTE(*input.CreatedAtFrom))
	}
	if input.CreatedAtTo != nil {
		filters = append(filters, credittransaction.CreateTimeLT(input.CreatedAtTo.AddDate(0, 0, 1)))
	}
	query := r.client.CreditTransaction.Query().Where(filters...)
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count credit transactions: %w", err)
	}
	items, err := query.WithUser().Order(ent.Desc(credittransaction.FieldCreateTime)).Offset((input.Page - 1) * input.PageSize).Limit(input.PageSize).All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("list credit transactions: %w", err)
	}
	return items, total, nil
}
func (r *Repository) Stats(ctx context.Context) (StatsView, error) {
	totalUsers, err := r.client.User.Query().Count(ctx)
	if err != nil {
		return StatsView{}, fmt.Errorf("count users: %w", err)
	}
	totalInviters, err := r.client.User.Query().Where(user.HasSentReferrals()).Count(ctx)
	if err != nil {
		return StatsView{}, fmt.Errorf("count inviters: %w", err)
	}
	totalReferrals, err := r.client.Referral.Query().Count(ctx)
	if err != nil {
		return StatsView{}, fmt.Errorf("count referrals: %w", err)
	}
	credits := 0
	if totalReferrals > 0 {
		credits, err = r.client.Referral.Query().Aggregate(ent.Sum(referral.FieldReward)).Int(ctx)
		if err != nil {
			return StatsView{}, fmt.Errorf("sum referral rewards: %w", err)
		}
	}
	return StatsView{TotalUsers: totalUsers, TotalInviters: totalInviters, TotalReferrals: totalReferrals, TotalCreditsAwarded: credits}, nil
}
