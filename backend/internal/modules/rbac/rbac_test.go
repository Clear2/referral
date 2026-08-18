package rbac

import (
	"context"
	"errors"
	"testing"

	"github.com/keep/sunny/ent/enttest"
	appErrors "github.com/keep/sunny/pkg/errors"
	"github.com/keep/sunny/pkg/postgres"
	_ "github.com/mattn/go-sqlite3"
)

func testService(t *testing.T) (*Service, *Repository) {
	t.Helper()
	client := enttest.Open(t, "sqlite3", "file:rbac?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() { _ = client.Close() })
	repository := NewRepository(&postgres.Postgres{Client: client})
	return NewService(repository), repository
}

func TestUnassignedFirstUserDoesNotBootstrapAdmin(t *testing.T) {
	service, repository := testService(t)
	account, err := repository.client.User.Create().SetName("first").SetEmail("first@example.com").Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	access, err := service.Access(context.Background(), account.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(access.Roles) != 0 || len(access.Permissions) != 0 {
		t.Fatalf("unassigned access = %+v, want no privileges", access)
	}
}

func TestSuperAdminAllowsEveryPermission(t *testing.T) {
	service, repository := testService(t)
	adminRole, err := repository.client.Role.Create().SetName("admin").SetCode("super_admin").SetIsSystem(true).Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	account, err := repository.client.User.Create().SetName("admin").SetEmail("admin@example.com").AddRoles(adminRole).Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	allowed, err := service.Allowed(context.Background(), account.ID, "referral:read")
	if err != nil || !allowed {
		t.Fatalf("Allowed() = %v, %v, want true, nil", allowed, err)
	}
}

func TestAllowedAnyAcceptsOneAdministrativePermission(t *testing.T) {
	service, repository := testService(t)
	permission, err := repository.client.Permission.Create().SetName("Referral read").SetCode("referral:read").SetModule("referral").Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assignedRole, err := repository.client.Role.Create().SetName("operator").SetCode("operator").AddPermissions(permission).Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	account, err := repository.client.User.Create().SetName("operator").SetEmail("operator@example.com").AddRoles(assignedRole).Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	allowed, err := service.AllowedAny(context.Background(), account.ID, "referral:read", "system:rbac")
	if err != nil || !allowed {
		t.Fatalf("AllowedAny() = %v, %v, want true, nil", allowed, err)
	}
}

func TestAllowedAnyRejectsUserWithoutAdministrativePermission(t *testing.T) {
	service, repository := testService(t)
	account, err := repository.client.User.Create().SetName("user").SetEmail("user@example.com").Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	allowed, err := service.AllowedAny(context.Background(), account.ID, "referral:read", "system:rbac")
	if err != nil || allowed {
		t.Fatalf("AllowedAny() = %v, %v, want false, nil", allowed, err)
	}
}

func TestSnapshotUsesEmptyRelationshipArrays(t *testing.T) {
	_, repository := testService(t)
	if _, err := repository.client.Role.Create().SetName("empty role").SetCode("empty_role").Save(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := repository.client.Permission.Create().SetName("empty permission").SetCode("empty:read").SetModule("empty").Save(context.Background()); err != nil {
		t.Fatal(err)
	}

	snapshot, err := repository.Snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Roles[0].PermissionIDs == nil || snapshot.Roles[0].MenuIDs == nil {
		t.Fatalf("role relationships must be empty arrays: %+v", snapshot.Roles[0])
	}
	if snapshot.Permissions[0].APIIDs == nil || snapshot.Permissions[0].MenuIDs == nil {
		t.Fatalf("permission relationships must be empty arrays: %+v", snapshot.Permissions[0])
	}
}

func TestSystemRoleCannotBeChanged(t *testing.T) {
	service, repository := testService(t)
	adminRole, err := repository.client.Role.Create().SetName("admin").SetCode("super_admin").SetIsSystem(true).Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	err = service.UpdateRole(context.Background(), adminRole.ID, RoleInput{Name: "renamed", Code: "renamed"})
	var apiErr *appErrors.APIError
	if !errors.As(err, &apiErr) || apiErr.Code != 403 {
		t.Fatalf("UpdateRole() error = %v, want 403", err)
	}
}

func TestAccessIncludesMenusGrantedThroughPermission(t *testing.T) {
	service, repository := testService(t)
	menuItem, err := repository.client.Menu.Create().SetName("Dashboard").SetPath("/").Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	permissionItem, err := repository.client.Permission.Create().SetName("Referral read").SetCode("referral:read").SetModule("referral").AddMenus(menuItem).Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assignedRole, err := repository.client.Role.Create().SetName("operator").SetCode("operator").AddPermissions(permissionItem).Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	account, err := repository.client.User.Create().SetName("operator").SetEmail("menus@example.com").AddRoles(assignedRole).Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	access, err := service.Access(context.Background(), account.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(access.Menus) != 1 || access.Menus[0].Path != "/" {
		t.Fatalf("Access().Menus = %+v, want dashboard menu", access.Menus)
	}
}

func TestAccessIncludesAllEnabledMenusForSuperAdmin(t *testing.T) {
	service, repository := testService(t)
	if _, err := repository.client.Menu.Create().SetName("Dashboard").SetPath("/").SetSortOrder(20).Save(context.Background()); err != nil {
		t.Fatal(err)
	}
	usersMenu, err := repository.client.Menu.Create().SetName("Users").SetPath("/users").SetSortOrder(10).Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repository.client.Menu.Create().SetName("Disabled").SetPath("/disabled").SetEnabled(false).Save(context.Background()); err != nil {
		t.Fatal(err)
	}
	adminRole, err := repository.client.Role.Create().SetName("admin").SetCode("super_admin").SetIsSystem(true).Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	account, err := repository.client.User.Create().SetName("admin").SetEmail("admin-menus@example.com").AddRoles(adminRole).Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	access, err := service.Access(context.Background(), account.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(access.Menus) != 2 || access.Menus[0].ID != usersMenu.ID {
		t.Fatalf("Access().Menus = %+v, want all enabled menus in sort order", access.Menus)
	}
	if len(access.MenuIDs) != 2 || access.MenuIDs[0] != usersMenu.ID {
		t.Fatalf("Access().MenuIDs = %+v, want enabled menu IDs in sort order", access.MenuIDs)
	}
}

func TestSyncAPIIsIdempotentAndRefreshesRouteMetadata(t *testing.T) {
	_, repository := testService(t)
	enabled := true
	first := APIInput{Name: "old handler", Method: "GET", Path: "/api/v1/example", Enabled: &enabled}
	if err := repository.SyncAPI(context.Background(), first); err != nil {
		t.Fatal(err)
	}
	second := APIInput{Name: "new handler", Method: "GET", Path: "/api/v1/example", Enabled: &enabled}
	if err := repository.SyncAPI(context.Background(), second); err != nil {
		t.Fatal(err)
	}

	items, err := repository.client.APIResource.Query().All(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Name != second.Name || !items[0].Enabled {
		t.Fatalf("SyncAPI() resources = %+v, want one refreshed resource", items)
	}
}

func TestAllowedResourceUsesPermissionAPIMap(t *testing.T) {
	service, repository := testService(t)
	api, err := repository.client.APIResource.Create().SetName("Users list").SetMethod("GET").SetPath("/api/v1/admin/users").Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	permissionItem, err := repository.client.Permission.Create().SetName("Users read").SetCode("users:read").SetModule("users").AddApis(api).Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assignedRole, err := repository.client.Role.Create().SetName("operator").SetCode("operator").AddPermissions(permissionItem).Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	account, err := repository.client.User.Create().SetName("operator").SetEmail("api@example.com").AddRoles(assignedRole).Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	allowed, err := service.AllowedResource(context.Background(), account.ID, "GET", "/api/v1/admin/users")
	if err != nil || !allowed {
		t.Fatalf("AllowedResource() = %v, %v, want true, nil", allowed, err)
	}
	allowed, err = service.AllowedResource(context.Background(), account.ID, "DELETE", "/api/v1/admin/users")
	if err != nil || allowed {
		t.Fatalf("AllowedResource() for unmapped operation = %v, %v, want false, nil", allowed, err)
	}
}
