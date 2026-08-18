package auth

import (
	"context"
	"testing"

	"github.com/keep/sunny/config"
	"github.com/keep/sunny/ent/enttest"
	"github.com/keep/sunny/ent/role"
	"github.com/keep/sunny/ent/user"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func TestBootstrapAdminUpdatesCredentialsAndAssignsConfiguredUsers(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:"+t.Name()+"?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() { _ = client.Close() })
	ctx := context.Background()
	adminRole, err := client.Role.Create().SetName("Super Admin").SetCode("super_admin").SetIsSystem(true).Save(ctx)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}
	existingHash, err := bcrypt.GenerateFromPassword([]byte("old-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash old password: %v", err)
	}
	bootstrapUser, err := client.User.Create().SetName("Old Name").SetEmail("admin@example.com").SetPasswordHash(string(existingHash)).Save(ctx)
	if err != nil {
		t.Fatalf("create bootstrap user: %v", err)
	}
	configuredUser, err := client.User.Create().SetName("Operator").SetEmail("operator@example.com").Save(ctx)
	if err != nil {
		t.Fatalf("create configured user: %v", err)
	}

	manager := &Manager{client: client}
	err = manager.bootstrapAdmin(ctx, config.Admin{
		UserIDs: []int{configuredUser.ID},
		Emails:  []string{" ADMIN@EXAMPLE.COM ", "operator@example.com"},
		Bootstrap: config.AdminBootstrap{
			Enabled: true, Name: "Super Admin", Email: "admin@example.com", Password: "new-password",
		},
	})
	if err != nil {
		t.Fatalf("bootstrapAdmin(): %v", err)
	}

	updated, err := client.User.Query().Where(user.IDEQ(bootstrapUser.ID)).WithRoles().Only(ctx)
	if err != nil {
		t.Fatalf("query bootstrap user: %v", err)
	}
	if updated.Name != "Super Admin" || bcrypt.CompareHashAndPassword([]byte(updated.PasswordHash), []byte("new-password")) != nil {
		t.Fatalf("bootstrap credentials were not updated")
	}
	for _, userID := range []int{bootstrapUser.ID, configuredUser.ID} {
		hasRole, queryErr := client.User.Query().Where(user.IDEQ(userID), user.HasRolesWith(role.IDEQ(adminRole.ID))).Exist(ctx)
		if queryErr != nil || !hasRole {
			t.Fatalf("user %d super_admin = %v, err=%v, want true", userID, hasRole, queryErr)
		}
	}
}

func TestBootstrapAdminPreservesNameWhenConfiguredNameIsTaken(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:"+t.Name()+"?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() { _ = client.Close() })
	ctx := context.Background()
	if _, err := client.Role.Create().SetName("Super Admin Role").SetCode("super_admin").SetIsSystem(true).Save(ctx); err != nil {
		t.Fatalf("create role: %v", err)
	}
	if _, err := client.User.Create().SetName("Super Admin").SetEmail("other@example.com").Save(ctx); err != nil {
		t.Fatalf("create conflicting user: %v", err)
	}
	bootstrapUser, err := client.User.Create().SetName("Existing Admin").SetEmail("admin@example.com").Save(ctx)
	if err != nil {
		t.Fatalf("create bootstrap user: %v", err)
	}

	manager := &Manager{client: client}
	err = manager.bootstrapAdmin(ctx, config.Admin{Bootstrap: config.AdminBootstrap{
		Enabled: true, Name: "Super Admin", Email: "admin@example.com", Password: "new-password",
	}})
	if err != nil {
		t.Fatalf("bootstrapAdmin(): %v", err)
	}

	updated, err := client.User.Get(ctx, bootstrapUser.ID)
	if err != nil {
		t.Fatalf("get bootstrap user: %v", err)
	}
	if updated.Name != "Existing Admin" {
		t.Fatalf("bootstrap name = %q, want preserved name", updated.Name)
	}
	if bcrypt.CompareHashAndPassword([]byte(updated.PasswordHash), []byte("new-password")) != nil {
		t.Fatal("bootstrap password was not updated")
	}
}
