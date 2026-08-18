package client

import (
	"context"
	"testing"

	"github.com/keep/sunny/ent/enttest"
	"github.com/keep/sunny/pkg/postgres"
	_ "github.com/mattn/go-sqlite3"
)

func newTestRepository(t *testing.T) *repository {
	t.Helper()
	client := enttest.Open(t, "sqlite3", "file:"+t.Name()+"?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() { _ = client.Close() })
	return &repository{pg: &postgres.Postgres{Client: client}}
}

func TestSetReferralCodeDoesNotOverwriteExistingCode(t *testing.T) {
	repository := newTestRepository(t)
	account, err := repository.pg.Client.User.Create().SetName("Alice").SetEmail("alice@example.com").Save(context.Background())
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	first, err := repository.SetReferralCode(context.Background(), account.ID, "AAAAAAAA")
	if err != nil {
		t.Fatalf("SetReferralCode(first): %v", err)
	}
	second, err := repository.SetReferralCode(context.Background(), account.ID, "BBBBBBBB")
	if err != nil {
		t.Fatalf("SetReferralCode(second): %v", err)
	}
	if *first.ReferralCode != "AAAAAAAA" || *second.ReferralCode != "AAAAAAAA" {
		t.Fatalf("referral codes = (%q, %q), want both AAAAAAAA", *first.ReferralCode, *second.ReferralCode)
	}
}

func TestRegisterPersistsRewardAtomically(t *testing.T) {
	repository := newTestRepository(t)
	inviter, err := repository.pg.Client.User.Create().SetName("Alice").SetEmail("alice@example.com").SetReferralCode("AAAAAAAA").Save(context.Background())
	if err != nil {
		t.Fatalf("create inviter: %v", err)
	}
	invitee, rewarded, err := repository.Register(context.Background(), RegisterInput{Code: "AAAAAAAA", Name: "Bob", Email: "bob@example.com"})
	if err != nil {
		t.Fatalf("Register(): %v", err)
	}
	if invitee.ID == 0 || rewarded.ID != inviter.ID || rewarded.CreditBalance != RewardCredits {
		t.Fatalf("Register() = invitee=%+v inviter=%+v", invitee, rewarded)
	}
	if count, queryErr := repository.pg.Client.Referral.Query().Count(context.Background()); queryErr != nil || count != 1 {
		t.Fatalf("referral count = %d, err=%v, want 1", count, queryErr)
	}
	if count, queryErr := repository.pg.Client.CreditTransaction.Query().Count(context.Background()); queryErr != nil || count != 1 {
		t.Fatalf("credit transaction count = %d, err=%v, want 1", count, queryErr)
	}
}

func TestDashboardReturnsReferralsSentByUser(t *testing.T) {
	repository := newTestRepository(t)
	inviter, err := repository.pg.Client.User.Create().
		SetName("Sunny").
		SetEmail("sunny@example.com").
		SetReferralCode("AAAAAAAA").
		Save(context.Background())
	if err != nil {
		t.Fatalf("create inviter: %v", err)
	}
	invitee, _, err := repository.Register(context.Background(), RegisterInput{
		Code:  "AAAAAAAA",
		Name:  "AAA",
		Email: "aaa@example.com",
	})
	if err != nil {
		t.Fatalf("Register(): %v", err)
	}

	dashboard, err := repository.Dashboard(context.Background(), inviter.ID)
	if err != nil {
		t.Fatalf("Dashboard(): %v", err)
	}
	if len(dashboard.Edges.SentReferrals) != 1 {
		t.Fatalf("sent referrals = %d, want 1", len(dashboard.Edges.SentReferrals))
	}
	got := dashboard.Edges.SentReferrals[0].Edges.Invitee
	if got == nil || got.ID != invitee.ID {
		t.Fatalf("dashboard invitee = %+v, want user %d", got, invitee.ID)
	}
}
