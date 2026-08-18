package rbac

import (
	"context"
	"sort"
	"time"

	"github.com/keep/sunny/ent"
	"github.com/keep/sunny/ent/apiresource"
	"github.com/keep/sunny/ent/menu"
	"github.com/keep/sunny/ent/permission"
	"github.com/keep/sunny/ent/permissionauditlog"
	"github.com/keep/sunny/ent/permissiongroup"
	"github.com/keep/sunny/ent/role"
	"github.com/keep/sunny/ent/user"
	"github.com/keep/sunny/pkg/postgres"
)

type Repository struct{ client *ent.Client }

func NewRepository(pg *postgres.Postgres) *Repository { return &Repository{client: pg.Client} }

func (r *Repository) Snapshot(ctx context.Context) (Snapshot, error) {
	roles, err := r.client.Role.Query().WithPermissions().WithMenus().Order(ent.Asc(role.FieldID)).All(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	permissions, err := r.client.Permission.Query().WithRoles().WithApis().WithMenus().Order(ent.Asc(permission.FieldModule), ent.Asc(permission.FieldCode)).All(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	menus, err := r.client.Menu.Query().Order(ent.Asc(menu.FieldSortOrder), ent.Asc(menu.FieldID)).All(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	groups, err := r.client.PermissionGroup.Query().Order(ent.Asc(permissiongroup.FieldSortOrder), ent.Asc(permissiongroup.FieldID)).All(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	apis, err := r.client.APIResource.Query().Order(ent.Asc(apiresource.FieldPath), ent.Asc(apiresource.FieldMethod)).All(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	audits, err := r.client.PermissionAuditLog.Query().Order(ent.Desc(permissionauditlog.FieldCreateTime)).Limit(100).All(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	out := Snapshot{Roles: make([]RoleView, 0, len(roles)), Permissions: make([]PermissionView, 0, len(permissions)), Menus: make([]MenuView, 0, len(menus)), Groups: make([]PermissionGroupView, 0, len(groups)), APIs: make([]APIView, 0, len(apis)), AuditLogs: make([]AuditView, 0, len(audits))}
	for _, item := range roles {
		view := RoleView{ID: item.ID, Name: item.Name, Code: item.Code, Description: item.Description, Enabled: item.Enabled, IsSystem: item.IsSystem, PermissionIDs: make([]int, 0, len(item.Edges.Permissions)), MenuIDs: make([]int, 0, len(item.Edges.Menus))}
		for _, p := range item.Edges.Permissions {
			view.PermissionIDs = append(view.PermissionIDs, p.ID)
		}
		for _, m := range item.Edges.Menus {
			view.MenuIDs = append(view.MenuIDs, m.ID)
		}
		sort.Ints(view.PermissionIDs)
		sort.Ints(view.MenuIDs)
		out.Roles = append(out.Roles, view)
	}
	for _, item := range permissions {
		view := PermissionView{ID: item.ID, Name: item.Name, Code: item.Code, Module: item.Module, Description: item.Description, Enabled: item.Enabled, RoleCount: len(item.Edges.Roles), GroupID: item.GroupID, APIIDs: make([]int, 0, len(item.Edges.Apis)), MenuIDs: make([]int, 0, len(item.Edges.Menus))}
		for _, api := range item.Edges.Apis {
			view.APIIDs = append(view.APIIDs, api.ID)
		}
		for _, linkedMenu := range item.Edges.Menus {
			view.MenuIDs = append(view.MenuIDs, linkedMenu.ID)
		}
		out.Permissions = append(out.Permissions, view)
	}
	for _, item := range menus {
		out.Menus = append(out.Menus, MenuView{ID: item.ID, Name: item.Name, Path: item.Path, Icon: item.Icon, Component: item.Component, Redirect: item.Redirect, Type: string(item.Type), ParentID: item.ParentID, SortOrder: item.SortOrder, Enabled: item.Enabled})
	}
	for _, item := range groups {
		out.Groups = append(out.Groups, PermissionGroupView{ID: item.ID, Name: item.Name, Module: item.Module, Description: item.Description, ParentID: item.ParentID, SortOrder: item.SortOrder, Enabled: item.Enabled})
	}
	for _, item := range apis {
		out.APIs = append(out.APIs, APIView{ID: item.ID, Name: item.Name, Method: item.Method, Path: item.Path, Description: item.Description, Enabled: item.Enabled})
	}
	for _, item := range audits {
		out.AuditLogs = append(out.AuditLogs, AuditView{ID: item.ID, OperatorID: item.OperatorID, Action: item.Action, TargetType: item.TargetType, TargetID: item.TargetID, RequestID: item.RequestID, IPAddress: item.IPAddress, CreateTime: item.CreateTime.Format(time.RFC3339)})
	}
	return out, nil
}

func (r *Repository) CreateRole(ctx context.Context, in RoleInput) (*ent.Role, error) {
	b := r.client.Role.Create().SetName(in.Name).SetCode(in.Code).SetDescription(in.Description)
	if in.Enabled != nil {
		b.SetEnabled(*in.Enabled)
	}
	return b.Save(ctx)
}
func (r *Repository) UpdateRole(ctx context.Context, id int, in RoleInput) (*ent.Role, error) {
	b := r.client.Role.UpdateOneID(id).SetName(in.Name).SetCode(in.Code).SetDescription(in.Description)
	if in.Enabled != nil {
		b.SetEnabled(*in.Enabled)
	}
	return b.Save(ctx)
}
func (r *Repository) DeleteRole(ctx context.Context, id int) error {
	return r.client.Role.DeleteOneID(id).Exec(ctx)
}

func (r *Repository) CreatePermission(ctx context.Context, in PermissionInput) (*ent.Permission, error) {
	b := r.client.Permission.Create().SetName(in.Name).SetCode(in.Code).SetModule(in.Module).SetDescription(in.Description).SetNillableGroupID(in.GroupID)
	if in.Enabled != nil {
		b.SetEnabled(*in.Enabled)
	}
	return b.Save(ctx)
}
func (r *Repository) UpdatePermission(ctx context.Context, id int, in PermissionInput) (*ent.Permission, error) {
	b := r.client.Permission.UpdateOneID(id).SetName(in.Name).SetCode(in.Code).SetModule(in.Module).SetDescription(in.Description).SetNillableGroupID(in.GroupID)
	if in.GroupID == nil {
		b.ClearGroupID()
	}
	if in.Enabled != nil {
		b.SetEnabled(*in.Enabled)
	}
	return b.Save(ctx)
}
func (r *Repository) DeletePermission(ctx context.Context, id int) error {
	return r.client.Permission.DeleteOneID(id).Exec(ctx)
}

func (r *Repository) CreateMenu(ctx context.Context, in MenuInput) (*ent.Menu, error) {
	b := r.client.Menu.Create().SetName(in.Name).SetPath(in.Path).SetIcon(in.Icon).SetComponent(in.Component).SetRedirect(in.Redirect).SetType(menu.Type(in.Type)).SetSortOrder(in.SortOrder).SetNillableParentID(in.ParentID)
	if in.Enabled != nil {
		b.SetEnabled(*in.Enabled)
	}
	return b.Save(ctx)
}
func (r *Repository) UpdateMenu(ctx context.Context, id int, in MenuInput) (*ent.Menu, error) {
	b := r.client.Menu.UpdateOneID(id).SetName(in.Name).SetPath(in.Path).SetIcon(in.Icon).SetComponent(in.Component).SetRedirect(in.Redirect).SetType(menu.Type(in.Type)).SetSortOrder(in.SortOrder).SetNillableParentID(in.ParentID)
	if in.ParentID == nil {
		b.ClearParentID()
	}
	if in.Enabled != nil {
		b.SetEnabled(*in.Enabled)
	}
	return b.Save(ctx)
}
func (r *Repository) DeleteMenu(ctx context.Context, id int) error {
	return r.client.Menu.DeleteOneID(id).Exec(ctx)
}

func (r *Repository) SetGrants(ctx context.Context, roleID int, in GrantInput) error {
	return r.client.Role.UpdateOneID(roleID).ClearPermissions().ClearMenus().AddPermissionIDs(in.PermissionIDs...).AddMenuIDs(in.MenuIDs...).Exec(ctx)
}
func (r *Repository) SetPermissionResources(ctx context.Context, id int, in ResourceGrantInput) error {
	return r.client.Permission.UpdateOneID(id).ClearApis().ClearMenus().AddAPIIDs(in.APIIDs...).AddMenuIDs(in.MenuIDs...).Exec(ctx)
}
func (r *Repository) CreateGroup(ctx context.Context, in PermissionGroupInput) error {
	b := r.client.PermissionGroup.Create().SetName(in.Name).SetModule(in.Module).SetDescription(in.Description).SetSortOrder(in.SortOrder).SetNillableParentID(in.ParentID)
	if in.Enabled != nil {
		b.SetEnabled(*in.Enabled)
	}
	return b.Exec(ctx)
}
func (r *Repository) UpdateGroup(ctx context.Context, id int, in PermissionGroupInput) error {
	b := r.client.PermissionGroup.UpdateOneID(id).SetName(in.Name).SetModule(in.Module).SetDescription(in.Description).SetSortOrder(in.SortOrder).SetNillableParentID(in.ParentID)
	if in.ParentID == nil {
		b.ClearParentID()
	}
	if in.Enabled != nil {
		b.SetEnabled(*in.Enabled)
	}
	return b.Exec(ctx)
}
func (r *Repository) DeleteGroup(ctx context.Context, id int) error {
	return r.client.PermissionGroup.DeleteOneID(id).Exec(ctx)
}
func (r *Repository) CreateAPI(ctx context.Context, in APIInput) error {
	b := r.client.APIResource.Create().SetName(in.Name).SetMethod(in.Method).SetPath(in.Path).SetDescription(in.Description)
	if in.Enabled != nil {
		b.SetEnabled(*in.Enabled)
	}
	return b.Exec(ctx)
}
func (r *Repository) UpdateAPI(ctx context.Context, id int, in APIInput) error {
	b := r.client.APIResource.UpdateOneID(id).SetName(in.Name).SetMethod(in.Method).SetPath(in.Path).SetDescription(in.Description)
	if in.Enabled != nil {
		b.SetEnabled(*in.Enabled)
	}
	return b.Exec(ctx)
}
func (r *Repository) DeleteAPI(ctx context.Context, id int) error {
	return r.client.APIResource.DeleteOneID(id).Exec(ctx)
}
func (r *Repository) SyncAPI(ctx context.Context, in APIInput) error {
	existing, err := r.client.APIResource.Query().Where(apiresource.MethodEQ(in.Method), apiresource.PathEQ(in.Path)).Only(ctx)
	if ent.IsNotFound(err) {
		if createErr := r.CreateAPI(ctx, in); createErr == nil {
			return nil
		} else if !ent.IsConstraintError(createErr) {
			return createErr
		}
		existing, err = r.client.APIResource.Query().Where(apiresource.MethodEQ(in.Method), apiresource.PathEQ(in.Path)).Only(ctx)
	}
	if err != nil {
		return err
	}
	return r.client.APIResource.UpdateOne(existing).SetName(in.Name).SetEnabled(true).Exec(ctx)
}
func (r *Repository) Audit(ctx context.Context, operatorID int, action, target string, targetID int, requestID, ip string) error {
	return r.client.PermissionAuditLog.Create().SetOperatorID(operatorID).SetAction(action).SetTargetType(target).SetTargetID(targetID).SetRequestID(requestID).SetIPAddress(ip).Exec(ctx)
}
func (r *Repository) IsSystemRole(ctx context.Context, id int) (bool, error) {
	item, err := r.client.Role.Get(ctx, id)
	if err != nil {
		return false, err
	}
	return item.IsSystem, nil
}
func (r *Repository) Access(ctx context.Context, userID int) (AccessView, error) {
	item, err := r.client.User.Query().Where(user.IDEQ(userID), user.EnabledEQ(true)).WithRoles(func(q *ent.RoleQuery) {
		q.Where(role.EnabledEQ(true)).WithPermissions(func(p *ent.PermissionQuery) {
			p.Where(permission.EnabledEQ(true)).WithMenus(func(m *ent.MenuQuery) { m.Where(menu.EnabledEQ(true)) })
		}).WithMenus(func(m *ent.MenuQuery) { m.Where(menu.EnabledEQ(true)) })
	}).Only(ctx)
	if err != nil {
		return AccessView{}, err
	}
	out := AccessView{
		UserID:      userID,
		Roles:       []string{},
		Permissions: []string{},
		MenuIDs:     []int{},
		Menus:       []MenuView{},
	}
	roleSet, permissionSet, menuSet := map[string]bool{}, map[string]bool{}, map[int]bool{}
	for _, assigned := range item.Edges.Roles {
		roleSet[assigned.Code] = true
		for _, p := range assigned.Edges.Permissions {
			permissionSet[p.Code] = true
			for _, m := range p.Edges.Menus {
				menuSet[m.ID] = true
			}
		}
		for _, m := range assigned.Edges.Menus {
			menuSet[m.ID] = true
		}
	}
	for code := range roleSet {
		out.Roles = append(out.Roles, code)
	}
	for code := range permissionSet {
		out.Permissions = append(out.Permissions, code)
	}
	for id := range menuSet {
		out.MenuIDs = append(out.MenuIDs, id)
	}
	sort.Strings(out.Roles)
	sort.Strings(out.Permissions)
	sort.Ints(out.MenuIDs)
	if roleSet["super_admin"] {
		menus, queryErr := r.client.Menu.Query().Where(menu.EnabledEQ(true)).Order(ent.Asc(menu.FieldSortOrder), ent.Asc(menu.FieldID)).All(ctx)
		if queryErr != nil {
			return AccessView{}, queryErr
		}
		out.MenuIDs = make([]int, 0, len(menus))
		out.Menus = make([]MenuView, 0, len(menus))
		for _, m := range menus {
			out.MenuIDs = append(out.MenuIDs, m.ID)
			out.Menus = append(out.Menus, MenuView{ID: m.ID, Name: m.Name, Path: m.Path, Icon: m.Icon, Component: m.Component, Redirect: m.Redirect, Type: string(m.Type), ParentID: m.ParentID, SortOrder: m.SortOrder, Enabled: m.Enabled})
		}
		return out, nil
	}
	if len(out.MenuIDs) > 0 {
		menus, queryErr := r.client.Menu.Query().Where(menu.IDIn(out.MenuIDs...)).Order(ent.Asc(menu.FieldSortOrder), ent.Asc(menu.FieldID)).All(ctx)
		if queryErr != nil {
			return AccessView{}, queryErr
		}
		out.Menus = make([]MenuView, 0, len(menus))
		for _, m := range menus {
			out.Menus = append(out.Menus, MenuView{ID: m.ID, Name: m.Name, Path: m.Path, Icon: m.Icon, Component: m.Component, Redirect: m.Redirect, Type: string(m.Type), ParentID: m.ParentID, SortOrder: m.SortOrder, Enabled: m.Enabled})
		}
	} else {
	}
	return out, nil
}

func (r *Repository) ResourcePermissionCodes(ctx context.Context, method, path string) ([]string, error) {
	items, err := r.client.Permission.Query().Where(
		permission.EnabledEQ(true),
		permission.HasApisWith(apiresource.MethodEQ(method), apiresource.PathEQ(path), apiresource.EnabledEQ(true)),
	).All(ctx)
	if err != nil {
		return nil, err
	}
	codes := make([]string, 0, len(items))
	for _, item := range items {
		codes = append(codes, item.Code)
	}
	return codes, nil
}
