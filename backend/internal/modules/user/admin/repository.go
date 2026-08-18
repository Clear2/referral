package admin

import (
	"context"
	"fmt"
	"strings"

	"github.com/keep/sunny/ent"
	"github.com/keep/sunny/ent/credittransaction"
	"github.com/keep/sunny/ent/permission"
	"github.com/keep/sunny/ent/predicate"
	"github.com/keep/sunny/ent/referral"
	"github.com/keep/sunny/ent/role"
	"github.com/keep/sunny/ent/user"
	"github.com/keep/sunny/pkg/postgres"
)

type Repository struct{ client *ent.Client }

func NewRepository(pg *postgres.Postgres) *Repository { return &Repository{client: pg.Client} }

func (repository *Repository) List(ctx context.Context, input ListInput) ([]*ent.User, int, error) {
	filters := make([]predicate.User, 0, 3)
	if input.Query != "" {
		filters = append(filters, user.Or(user.NameContainsFold(input.Query), user.EmailContainsFold(input.Query)))
	}
	if input.Enabled != nil {
		filters = append(filters, user.EnabledEQ(*input.Enabled))
	}
	administrator := user.HasRolesWith(role.Or(
		role.CodeEQ("super_admin"),
		role.HasPermissionsWith(permission.CodeIn("referral:read", "system:rbac", "system:*")),
	))
	switch input.AccountType {
	case "customer":
		filters = append(filters, user.Not(administrator))
	case "admin":
		filters = append(filters, administrator)
	}
	query := repository.client.User.Query().Where(filters...)
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	items, err := query.WithRoles().WithSentReferrals().WithCreditTransactions().
		Order(ent.Desc(user.FieldCreateTime)).Offset((input.Page - 1) * input.PageSize).Limit(input.PageSize).All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	return items, total, nil
}

func (repository *Repository) Get(ctx context.Context, id int) (*ent.User, error) {
	return repository.client.User.Query().Where(user.IDEQ(id)).WithRoles().
		WithSentReferrals(func(query *ent.ReferralQuery) {
			query.WithInvitee().Order(ent.Desc(referral.FieldCreateTime))
		}).
		WithCreditTransactions(func(query *ent.CreditTransactionQuery) {
			query.Order(ent.Desc(credittransaction.FieldCreateTime))
		}).Only(ctx)
}

func (repository *Repository) Create(ctx context.Context, input CreateInput, passwordHash string) (*ent.User, error) {
	enabled := true
	if input.Enabled != nil {
		enabled = *input.Enabled
	}
	return repository.client.User.Create().SetName(input.Name).SetEmail(input.Email).
		SetPasswordHash(passwordHash).SetEnabled(enabled).Save(ctx)
}

func (repository *Repository) Update(ctx context.Context, id int, input UpdateInput) error {
	return repository.client.User.UpdateOneID(id).SetName(strings.TrimSpace(input.Name)).
		SetEmail(strings.ToLower(strings.TrimSpace(input.Email))).Exec(ctx)
}

func (repository *Repository) SetStatus(ctx context.Context, id int, enabled bool) error {
	return repository.client.User.UpdateOneID(id).SetEnabled(enabled).Exec(ctx)
}

func (repository *Repository) SetRoles(ctx context.Context, id int, roleIDs []int) error {
	return repository.client.User.UpdateOneID(id).ClearRoles().AddRoleIDs(roleIDs...).Exec(ctx)
}

func (repository *Repository) ResetPassword(ctx context.Context, id int, passwordHash string) error {
	return repository.client.User.UpdateOneID(id).SetPasswordHash(passwordHash).Exec(ctx)
}

func (repository *Repository) Audit(ctx context.Context, operatorID int, action string, targetID int, requestID, ip string) error {
	return repository.client.PermissionAuditLog.Create().SetOperatorID(operatorID).SetAction(action).
		SetTargetType("USER").SetTargetID(targetID).SetRequestID(requestID).SetIPAddress(ip).Exec(ctx)
}
