package admin

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/keep/sunny/ent"
	appErrors "github.com/keep/sunny/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type serviceRepository interface {
	List(context.Context, ListInput) ([]*ent.User, int, error)
	Get(context.Context, int) (*ent.User, error)
	Create(context.Context, CreateInput, string) (*ent.User, error)
	Update(context.Context, int, UpdateInput) error
	SetStatus(context.Context, int, bool) error
	SetRoles(context.Context, int, []int) error
	ResetPassword(context.Context, int, string) error
}

type Service struct{ repository serviceRepository }

func NewService(repository *Repository) *Service { return newService(repository) }

func newService(repository serviceRepository) *Service { return &Service{repository: repository} }

func adminError(err error) error {
	if err == nil {
		return nil
	}
	if ent.IsNotFound(err) {
		return appErrors.ErrNotFound
	}
	if ent.IsConstraintError(err) {
		return appErrors.NewAPIError(409, "邮箱已注册")
	}
	return appErrors.WrapAPIError(appErrors.ErrInternalServerError, err)
}

func toView(item *ent.User) UserView {
	out := UserView{ID: item.ID, Name: item.Name, Email: item.Email, Enabled: item.Enabled, ReferralCode: item.ReferralCode,
		CreditBalance: item.CreditBalance, SuccessfulReferrals: len(item.Edges.SentReferrals), CreditTransactions: len(item.Edges.CreditTransactions),
		RoleIDs: make([]int, 0, len(item.Edges.Roles)), Roles: make([]RoleView, 0, len(item.Edges.Roles)), CreatedAt: item.CreateTime, UpdatedAt: item.UpdateTime}
	for _, role := range item.Edges.Roles {
		out.RoleIDs = append(out.RoleIDs, role.ID)
		out.Roles = append(out.Roles, RoleView{ID: role.ID, Name: role.Name, Code: role.Code})
	}
	sort.Ints(out.RoleIDs)
	return out
}

func (service *Service) List(ctx context.Context, input ListInput) (*ListView, error) {
	if input.Page == 0 {
		input.Page = 1
	}
	if input.PageSize == 0 {
		input.PageSize = 20
	}
	input.Query = strings.TrimSpace(input.Query)
	items, total, err := service.repository.List(ctx, input)
	if err != nil {
		return nil, adminError(err)
	}
	out := &ListView{Items: make([]UserView, 0, len(items)), Pagination: PaginationView{Page: input.Page, PageSize: input.PageSize, Total: total}}
	if total > 0 {
		out.Pagination.TotalPages = (total + input.PageSize - 1) / input.PageSize
	}
	for _, item := range items {
		out.Items = append(out.Items, toView(item))
	}
	return out, nil
}

func (service *Service) Get(ctx context.Context, id int) (*DetailView, error) {
	item, err := service.repository.Get(ctx, id)
	if err != nil {
		return nil, adminError(err)
	}
	out := &DetailView{User: toView(item), Referrals: make([]ReferralView, 0, len(item.Edges.SentReferrals)), CreditTransactions: make([]CreditView, 0, len(item.Edges.CreditTransactions))}
	for _, record := range item.Edges.SentReferrals {
		out.Referrals = append(out.Referrals, ReferralView{ID: record.ID, InviteeID: record.InviteeID, Name: record.Edges.Invitee.Name, Email: record.Edges.Invitee.Email, Reward: record.Reward, CreatedAt: record.CreateTime})
	}
	for _, record := range item.Edges.CreditTransactions {
		out.CreditTransactions = append(out.CreditTransactions, CreditView{ID: record.ID, ReferralID: record.ReferralID, Amount: record.Amount, Reason: record.Reason, CreatedAt: record.CreateTime})
	}
	return out, nil
}

func (service *Service) Create(ctx context.Context, input CreateInput) (*UserView, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if input.Password != input.ConfirmPassword {
		return nil, appErrors.NewAPIError(http.StatusBadRequest, "两次输入的密码不一致")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, appErrors.WrapAPIError(appErrors.ErrInternalServerError, fmt.Errorf("hash password: %w", err))
	}
	input.Password = ""
	input.ConfirmPassword = ""
	item, err := service.repository.Create(ctx, input, string(hash))
	if err != nil {
		return nil, adminError(err)
	}
	out := toView(item)
	return &out, nil
}

func (service *Service) Update(ctx context.Context, id int, input UpdateInput) error {
	return adminError(service.repository.Update(ctx, id, input))
}
func (service *Service) SetStatus(ctx context.Context, operatorID, id int, enabled bool) error {
	if operatorID == id && !enabled {
		return appErrors.NewAPIError(http.StatusForbidden, "不能禁用当前登录用户")
	}
	return adminError(service.repository.SetStatus(ctx, id, enabled))
}
func (service *Service) SetRoles(ctx context.Context, operatorID, id int, roleIDs []int) error {
	if operatorID == id {
		return appErrors.NewAPIError(http.StatusForbidden, "不能修改当前登录用户的角色")
	}
	return adminError(service.repository.SetRoles(ctx, id, roleIDs))
}
func (service *Service) ResetPassword(ctx context.Context, id int, input ResetPasswordInput) error {
	if input.Password != input.ConfirmPassword {
		return appErrors.NewAPIError(http.StatusBadRequest, "两次输入的密码不一致")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return appErrors.WrapAPIError(appErrors.ErrInternalServerError, fmt.Errorf("hash password: %w", err))
	}
	return adminError(service.repository.ResetPassword(ctx, id, string(hash)))
}
