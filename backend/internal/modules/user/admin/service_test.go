package admin

import (
	"context"
	"errors"
	"testing"

	"github.com/keep/sunny/ent"
	appErrors "github.com/keep/sunny/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type fakeServiceRepository struct {
	items        []*ent.User
	total        int
	createdInput CreateInput
	passwordHash string
	created      *ent.User
	resetID      int
	err          error
	listInput    ListInput
}

func TestServicePreventsSelfLockout(t *testing.T) {
	t.Parallel()
	service := newService(&fakeServiceRepository{})
	for name, action := range map[string]func() error{
		"disable self":     func() error { return service.SetStatus(context.Background(), 7, 7, false) },
		"change own roles": func() error { return service.SetRoles(context.Background(), 7, 7, []int{1}) },
	} {
		t.Run(name, func(t *testing.T) {
			var apiErr *appErrors.APIError
			if err := action(); !errors.As(err, &apiErr) || apiErr.Code != 403 {
				t.Fatalf("action error = %v, want 403", err)
			}
		})
	}
}

func (repository *fakeServiceRepository) List(_ context.Context, input ListInput) ([]*ent.User, int, error) {
	repository.listInput = input
	return repository.items, repository.total, repository.err
}
func (repository *fakeServiceRepository) Get(context.Context, int) (*ent.User, error) {
	return repository.created, repository.err
}
func (repository *fakeServiceRepository) Create(_ context.Context, input CreateInput, hash string) (*ent.User, error) {
	repository.createdInput, repository.passwordHash = input, hash
	return repository.created, repository.err
}
func (repository *fakeServiceRepository) Update(context.Context, int, UpdateInput) error {
	return repository.err
}
func (repository *fakeServiceRepository) SetStatus(context.Context, int, bool) error {
	return repository.err
}
func (repository *fakeServiceRepository) SetRoles(context.Context, int, []int) error {
	return repository.err
}
func (repository *fakeServiceRepository) ResetPassword(_ context.Context, id int, hash string) error {
	repository.resetID, repository.passwordHash = id, hash
	return repository.err
}

func TestServiceListUsesPaginationDefaults(t *testing.T) {
	t.Parallel()
	repository := &fakeServiceRepository{total: 21, items: []*ent.User{}}

	result, err := newService(repository).List(context.Background(), ListInput{})
	if err != nil {
		t.Fatalf("List(%v) error = %v, want nil", ListInput{}, err)
	}
	if result.Pagination.Page != 1 || result.Pagination.PageSize != 20 || result.Pagination.TotalPages != 2 {
		t.Errorf("List(%v).Pagination = %+v, want page=1 pageSize=20 totalPages=2", ListInput{}, result.Pagination)
	}
}

func TestServiceListPreservesAccountType(t *testing.T) {
	t.Parallel()
	repository := &fakeServiceRepository{items: []*ent.User{}}

	_, err := newService(repository).List(context.Background(), ListInput{AccountType: "admin"})
	if err != nil {
		t.Fatal(err)
	}
	if repository.listInput.AccountType != "admin" {
		t.Fatalf("account type = %q, want admin", repository.listInput.AccountType)
	}
}

func TestServiceCreateNormalizesAndHashesPassword(t *testing.T) {
	t.Parallel()
	repository := &fakeServiceRepository{created: &ent.User{ID: 7, Name: "Alice", Email: "alice@example.com", Enabled: true}}
	input := CreateInput{Name: " Alice ", Email: " ALICE@Example.com ", Password: "secret123", ConfirmPassword: "secret123"}

	result, err := newService(repository).Create(context.Background(), input)
	if err != nil {
		t.Fatalf("Create(%v) error = %v, want nil", input, err)
	}
	if result.ID != 7 {
		t.Errorf("Create(%v).ID = %d, want 7", input, result.ID)
	}
	if repository.createdInput.Name != "Alice" || repository.createdInput.Email != "alice@example.com" {
		t.Errorf("Create(%v) normalized input = %+v, want trimmed name and lowercase email", input, repository.createdInput)
	}
	if repository.createdInput.Password != "" || repository.createdInput.ConfirmPassword != "" {
		t.Errorf("Create(%v) forwarded plaintext passwords to repository", input)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repository.passwordHash), []byte(input.Password)); err != nil {
		t.Errorf("Create(%v) password hash does not match input password: %v", input, err)
	}
}

func TestServiceResetPasswordHashesPassword(t *testing.T) {
	t.Parallel()
	repository := &fakeServiceRepository{}

	err := newService(repository).ResetPassword(context.Background(), 9, ResetPasswordInput{Password: "newsecret123", ConfirmPassword: "newsecret123"})
	if err != nil {
		t.Fatalf("ResetPassword(9) error = %v, want nil", err)
	}
	if repository.resetID != 9 {
		t.Errorf("ResetPassword(9) repository ID = %d, want 9", repository.resetID)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repository.passwordHash), []byte("newsecret123")); err != nil {
		t.Errorf("ResetPassword(9) password hash does not match input password: %v", err)
	}
}

func TestServiceRejectsMismatchedPasswords(t *testing.T) {
	t.Parallel()

	for name, action := range map[string]func() error{
		"create": func() error {
			_, err := newService(&fakeServiceRepository{}).Create(context.Background(), CreateInput{Password: "secret123", ConfirmPassword: "different123"})
			return err
		},
		"reset": func() error {
			return newService(&fakeServiceRepository{}).ResetPassword(context.Background(), 9, ResetPasswordInput{Password: "secret123", ConfirmPassword: "different123"})
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var apiErr *appErrors.APIError
			if err := action(); !errors.As(err, &apiErr) || apiErr.Code != 400 {
				t.Fatalf("action error = %v, want 400", err)
			}
		})
	}
}
