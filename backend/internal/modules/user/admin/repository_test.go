package admin

import (
	"context"
	"testing"

	"github.com/keep/sunny/ent/enttest"
	"github.com/keep/sunny/pkg/postgres"
	_ "github.com/mattn/go-sqlite3"
)

func TestListSeparatesCustomersAndAdministrators(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:"+t.Name()+"?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() { _ = client.Close() })
	ctx := context.Background()
	adminPermission, err := client.Permission.Create().SetName("Referral read").SetCode("referral:read").SetModule("referral").Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	adminRole, err := client.Role.Create().SetName("operator").SetCode("operator").AddPermissions(adminPermission).Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = client.User.Create().SetName("Customer").SetEmail("customer@example.com").Save(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err = client.User.Create().SetName("Admin").SetEmail("admin@example.com").AddRoles(adminRole).Save(ctx); err != nil {
		t.Fatal(err)
	}
	repository := NewRepository(&postgres.Postgres{Client: client})

	customers, customerTotal, err := repository.List(ctx, ListInput{Page: 1, PageSize: 20, AccountType: "customer"})
	if err != nil || customerTotal != 1 || len(customers) != 1 || customers[0].Email != "customer@example.com" {
		t.Fatalf("customers = %+v, total = %d, err = %v", customers, customerTotal, err)
	}
	administrators, adminTotal, err := repository.List(ctx, ListInput{Page: 1, PageSize: 20, AccountType: "admin"})
	if err != nil || adminTotal != 1 || len(administrators) != 1 || administrators[0].Email != "admin@example.com" {
		t.Fatalf("administrators = %+v, total = %d, err = %v", administrators, adminTotal, err)
	}
}
