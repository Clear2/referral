package rbac

import (
	"context"
	"net/http"
	"strings"

	"github.com/keep/sunny/ent"
	appErrors "github.com/keep/sunny/pkg/errors"
)

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

func publicError(err error) error {
	switch {
	case err == nil:
		return nil
	case ent.IsNotFound(err):
		return appErrors.ErrNotFound
	case ent.IsConstraintError(err):
		return appErrors.WrapAPIError(appErrors.ErrConflict, err)
	default:
		return appErrors.WrapAPIError(appErrors.ErrInternalServerError, err)
	}
}

func normalizeCode(value string) string { return strings.ToLower(strings.TrimSpace(value)) }

func (s *Service) Snapshot(ctx context.Context) (Snapshot, error) {
	out, err := s.repo.Snapshot(ctx)
	return out, publicError(err)
}
func (s *Service) CreateRole(ctx context.Context, in RoleInput) error {
	in.Name = strings.TrimSpace(in.Name)
	in.Code = normalizeCode(in.Code)
	_, err := s.repo.CreateRole(ctx, in)
	return publicError(err)
}
func (s *Service) UpdateRole(ctx context.Context, id int, in RoleInput) error {
	system, err := s.repo.IsSystemRole(ctx, id)
	if err != nil {
		return publicError(err)
	}
	if system {
		return appErrors.NewAPIError(http.StatusForbidden, "系统角色不可修改")
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Code = normalizeCode(in.Code)
	_, err = s.repo.UpdateRole(ctx, id, in)
	return publicError(err)
}
func (s *Service) DeleteRole(ctx context.Context, id int) error {
	system, err := s.repo.IsSystemRole(ctx, id)
	if err != nil {
		return publicError(err)
	}
	if system {
		return appErrors.NewAPIError(http.StatusForbidden, "系统角色不可删除")
	}
	return publicError(s.repo.DeleteRole(ctx, id))
}
func (s *Service) CreatePermission(ctx context.Context, in PermissionInput) error {
	in.Name = strings.TrimSpace(in.Name)
	in.Code = normalizeCode(in.Code)
	in.Module = normalizeCode(in.Module)
	_, err := s.repo.CreatePermission(ctx, in)
	return publicError(err)
}
func (s *Service) UpdatePermission(ctx context.Context, id int, in PermissionInput) error {
	in.Name = strings.TrimSpace(in.Name)
	in.Code = normalizeCode(in.Code)
	in.Module = normalizeCode(in.Module)
	_, err := s.repo.UpdatePermission(ctx, id, in)
	return publicError(err)
}
func (s *Service) DeletePermission(ctx context.Context, id int) error {
	return publicError(s.repo.DeletePermission(ctx, id))
}
func (s *Service) CreateMenu(ctx context.Context, in MenuInput) error {
	in.Name = strings.TrimSpace(in.Name)
	_, err := s.repo.CreateMenu(ctx, in)
	return publicError(err)
}
func (s *Service) UpdateMenu(ctx context.Context, id int, in MenuInput) error {
	if in.ParentID != nil && *in.ParentID == id {
		return appErrors.NewAPIError(http.StatusBadRequest, "菜单不能以自身作为上级")
	}
	in.Name = strings.TrimSpace(in.Name)
	_, err := s.repo.UpdateMenu(ctx, id, in)
	return publicError(err)
}
func (s *Service) DeleteMenu(ctx context.Context, id int) error {
	return publicError(s.repo.DeleteMenu(ctx, id))
}
func (s *Service) SetGrants(ctx context.Context, id int, in GrantInput) error {
	system, err := s.repo.IsSystemRole(ctx, id)
	if err != nil {
		return publicError(err)
	}
	if system {
		return appErrors.NewAPIError(http.StatusForbidden, "系统角色授权不可修改")
	}
	return publicError(s.repo.SetGrants(ctx, id, in))
}
func (s *Service) SetPermissionResources(ctx context.Context, id int, in ResourceGrantInput) error {
	return publicError(s.repo.SetPermissionResources(ctx, id, in))
}
func (s *Service) CreateGroup(ctx context.Context, in PermissionGroupInput) error {
	return publicError(s.repo.CreateGroup(ctx, in))
}
func (s *Service) UpdateGroup(ctx context.Context, id int, in PermissionGroupInput) error {
	if in.ParentID != nil && *in.ParentID == id {
		return appErrors.ErrBadRequest
	}
	return publicError(s.repo.UpdateGroup(ctx, id, in))
}
func (s *Service) DeleteGroup(ctx context.Context, id int) error {
	return publicError(s.repo.DeleteGroup(ctx, id))
}
func (s *Service) CreateAPI(ctx context.Context, in APIInput) error {
	in.Method = strings.ToUpper(in.Method)
	return publicError(s.repo.CreateAPI(ctx, in))
}
func (s *Service) UpdateAPI(ctx context.Context, id int, in APIInput) error {
	in.Method = strings.ToUpper(in.Method)
	return publicError(s.repo.UpdateAPI(ctx, id, in))
}
func (s *Service) DeleteAPI(ctx context.Context, id int) error {
	return publicError(s.repo.DeleteAPI(ctx, id))
}
func (s *Service) Access(ctx context.Context, userID int) (AccessView, error) {
	out, err := s.repo.Access(ctx, userID)
	return out, publicError(err)
}
func (s *Service) Allowed(ctx context.Context, userID int, code string) (bool, error) {
	return s.AllowedAny(ctx, userID, code)
}

func (s *Service) AllowedAny(ctx context.Context, userID int, codes ...string) (bool, error) {
	access, err := s.Access(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, item := range access.Permissions {
		if item == "system:*" {
			return true, nil
		}
		for _, code := range codes {
			if item == code {
				return true, nil
			}
		}
	}
	for _, assignedRole := range access.Roles {
		if assignedRole == "super_admin" {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) AllowedResource(ctx context.Context, userID int, method, path string) (bool, error) {
	codes, err := s.repo.ResourcePermissionCodes(ctx, method, path)
	if err != nil || len(codes) == 0 {
		return false, err
	}
	return s.AllowedAny(ctx, userID, codes...)
}
