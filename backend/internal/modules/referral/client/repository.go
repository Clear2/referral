package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/keep/sunny/ent"
	"github.com/keep/sunny/ent/credittransaction"
	"github.com/keep/sunny/ent/referral"
	"github.com/keep/sunny/ent/user"
	"github.com/keep/sunny/pkg/postgres"
)

var (
	ErrInvitationNotFound = errors.New("invitation not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailExists        = errors.New("email already registered")
)

type Repository interface {
	GetUser(ctx context.Context, userID int) (*ent.User, error)
	SetReferralCode(ctx context.Context, userID int, code string) (*ent.User, error)
	Register(ctx context.Context, input RegisterInput) (*ent.User, *ent.User, error)
	Dashboard(ctx context.Context, userID int) (*ent.User, error)
}

type repository struct{ pg *postgres.Postgres }

func NewRepository(pg *postgres.Postgres) Repository { return &repository{pg: pg} }

func (r *repository) GetUser(ctx context.Context, userID int) (*ent.User, error) {
	entity, err := r.pg.Client.User.Get(ctx, userID)
	if ent.IsNotFound(err) {
		return nil, ErrUserNotFound
	}
	return entity, err
}

func (r *repository) SetReferralCode(ctx context.Context, userID int, code string) (*ent.User, error) {
	updated, err := r.pg.Client.User.Update().
		Where(user.IDEQ(userID), user.ReferralCodeIsNil()).
		SetReferralCode(code).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	entity, err := r.pg.Client.User.Get(ctx, userID)
	if ent.IsNotFound(err) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	if updated == 0 && entity.ReferralCode == nil {
		return nil, fmt.Errorf("referral code conditional update affected no rows")
	}
	return entity, err
}

func (r *repository) Register(ctx context.Context, input RegisterInput) (_ *ent.User, _ *ent.User, resultErr error) {
	tx, err := r.pg.Client.Tx(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if resultErr != nil {
			_ = tx.Rollback()
		}
	}()

	inviter, err := tx.User.Query().Where(user.ReferralCodeEQ(input.Code)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, nil, ErrInvitationNotFound
	}
	if err != nil {
		return nil, nil, fmt.Errorf("find inviter: %w", err)
	}

	invitee, err := tx.User.Create().SetName(input.Name).SetEmail(input.Email).Save(ctx)
	if ent.IsConstraintError(err) {
		return nil, nil, ErrEmailExists
	}
	if err != nil {
		return nil, nil, fmt.Errorf("create invitee: %w", err)
	}

	referralEntity, err := tx.Referral.Create().SetInviterID(inviter.ID).SetInviteeID(invitee.ID).SetReward(RewardCredits).Save(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("create referral: %w", err)
	}

	inviter, err = tx.User.UpdateOneID(inviter.ID).AddCreditBalance(RewardCredits).Save(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("reward inviter: %w", err)
	}

	_, err = tx.CreditTransaction.Create().
		SetUserID(inviter.ID).
		SetReferralID(referralEntity.ID).
		SetAmount(RewardCredits).
		SetReason("referral_reward").
		Save(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("record credit transaction: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf("commit referral: %w", err)
	}
	return invitee, inviter, nil
}

func (r *repository) Dashboard(ctx context.Context, userID int) (*ent.User, error) {
	entity, err := r.pg.Client.User.Query().
		Where(user.IDEQ(userID)).
		WithSentReferrals(func(query *ent.ReferralQuery) {
			query.WithInvitee().Order(ent.Desc(referral.FieldCreateTime))
		}).
		WithCreditTransactions(func(query *ent.CreditTransactionQuery) {
			query.Order(ent.Desc(credittransaction.FieldCreateTime))
		}).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, ErrUserNotFound
	}
	return entity, err
}
