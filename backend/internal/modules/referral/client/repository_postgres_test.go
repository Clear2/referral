//go:build integration

package client

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/keep/sunny/ent"
	"github.com/keep/sunny/ent/user"
	"github.com/keep/sunny/pkg/postgres"
	_ "github.com/lib/pq"
)

func newPostgresTestRepository(t *testing.T) *repository {
	t.Helper()
	dsn := os.Getenv("REFERRAL_TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("REFERRAL_TEST_DATABASE_URL is required for integration tests")
	}
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open PostgreSQL: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })
	if err = client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("create PostgreSQL schema: %v", err)
	}
	return &repository{pg: &postgres.Postgres{Client: client}}
}

func TestPostgresReferralRegistrationIsTransactional(t *testing.T) {
	ctx := context.Background()
	repository := newPostgresTestRepository(t)
	inviter, err := repository.pg.Client.User.Create().SetName("Postgres Inviter").SetEmail("pg-inviter@example.com").SetReferralCode("PGTEST22").Save(ctx)
	if err != nil {
		t.Fatalf("create inviter: %v", err)
	}

	invitee, rewarded, err := repository.Register(ctx, RegisterInput{Code: "PGTEST22", Name: "Invitee", Email: "pg-invitee@example.com"})
	if err != nil {
		t.Fatalf("register invitee: %v", err)
	}
	if invitee.ID == 0 || rewarded.CreditBalance != 100 {
		t.Fatalf("registration result = invitee %d, balance %d", invitee.ID, rewarded.CreditBalance)
	}
	dashboard, err := repository.Dashboard(ctx, inviter.ID)
	if err != nil {
		t.Fatalf("load dashboard: %v", err)
	}
	if len(dashboard.Edges.SentReferrals) != 1 || len(dashboard.Edges.CreditTransactions) != 1 {
		t.Fatalf("persisted referrals=%d transactions=%d, want 1 each", len(dashboard.Edges.SentReferrals), len(dashboard.Edges.CreditTransactions))
	}

	_, _, err = repository.Register(ctx, RegisterInput{Code: "MISSING2", Name: "Nobody", Email: "not-created@example.com"})
	if !errors.Is(err, ErrInvitationNotFound) {
		t.Fatalf("invalid invitation error = %v", err)
	}
	if exists, queryErr := repository.pg.Client.User.Query().Where(user.EmailEQ("not-created@example.com")).Exist(ctx); queryErr != nil || exists {
		t.Fatalf("invalid invitation created user: exists=%v err=%v", exists, queryErr)
	}
}

func TestPostgresConcurrentRewardsDoNotLoseCredits(t *testing.T) {
	ctx := context.Background()
	repository := newPostgresTestRepository(t)
	inviter, err := repository.pg.Client.User.Create().SetName("Concurrent Inviter").SetEmail("concurrent-inviter@example.com").SetReferralCode("CONCUR22").Save(ctx)
	if err != nil {
		t.Fatalf("create inviter: %v", err)
	}

	const registrations = 8
	errorsFound := make(chan error, registrations)
	var group sync.WaitGroup
	for index := range registrations {
		group.Add(1)
		go func() {
			defer group.Done()
			_, _, registerErr := repository.Register(ctx, RegisterInput{Code: "CONCUR22", Name: fmt.Sprintf("Invitee %d", index), Email: fmt.Sprintf("concurrent-%d@example.com", index)})
			errorsFound <- registerErr
		}()
	}
	group.Wait()
	close(errorsFound)
	for registerErr := range errorsFound {
		if registerErr != nil {
			t.Fatalf("concurrent registration: %v", registerErr)
		}
	}

	dashboard, err := repository.Dashboard(ctx, inviter.ID)
	if err != nil {
		t.Fatalf("load dashboard: %v", err)
	}
	if dashboard.CreditBalance != registrations*100 || len(dashboard.Edges.SentReferrals) != registrations || len(dashboard.Edges.CreditTransactions) != registrations {
		t.Fatalf("balance=%d referrals=%d transactions=%d", dashboard.CreditBalance, len(dashboard.Edges.SentReferrals), len(dashboard.Edges.CreditTransactions))
	}
}
