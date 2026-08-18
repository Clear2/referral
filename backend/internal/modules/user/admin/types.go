package admin

import "time"

type ListInput struct {
	Page        int    `form:"page" binding:"omitempty,min=1"`
	PageSize    int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Query       string `form:"query" binding:"omitempty,max=254"`
	Enabled     *bool  `form:"enabled"`
	AccountType string `form:"account_type" binding:"omitempty,oneof=customer admin"`
}

type CreateInput struct {
	Name            string `json:"name" binding:"required,min=1,max=100"`
	Email           string `json:"email" binding:"required,email,max=254"`
	Password        string `json:"password" binding:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=72"`
	Enabled         *bool  `json:"enabled"`
}

type UpdateInput struct {
	Name  string `json:"name" binding:"required,min=1,max=100"`
	Email string `json:"email" binding:"required,email,max=254"`
}

type StatusInput struct {
	Enabled bool `json:"enabled"`
}
type RolesInput struct {
	RoleIDs []int `json:"role_ids" binding:"dive,min=1"`
}
type ResetPasswordInput struct {
	Password        string `json:"password" binding:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=72"`
}

type RoleView struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type UserView struct {
	ID                  int        `json:"id"`
	Name                string     `json:"name"`
	Email               string     `json:"email"`
	Enabled             bool       `json:"enabled"`
	ReferralCode        *string    `json:"referral_code,omitempty"`
	CreditBalance       int        `json:"credit_balance"`
	SuccessfulReferrals int        `json:"successful_referrals"`
	CreditTransactions  int        `json:"credit_transactions"`
	RoleIDs             []int      `json:"role_ids"`
	Roles               []RoleView `json:"roles"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type ReferralView struct {
	ID        int       `json:"id"`
	InviteeID int       `json:"invitee_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Reward    int       `json:"reward"`
	CreatedAt time.Time `json:"created_at"`
}

type CreditView struct {
	ID         int       `json:"id"`
	ReferralID int       `json:"referral_id"`
	Amount     int       `json:"amount"`
	Reason     string    `json:"reason"`
	CreatedAt  time.Time `json:"created_at"`
}

type DetailView struct {
	User               UserView       `json:"user"`
	Referrals          []ReferralView `json:"referrals"`
	CreditTransactions []CreditView   `json:"credit_transactions"`
}

type PaginationView struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ListView struct {
	Items      []UserView     `json:"items"`
	Pagination PaginationView `json:"pagination"`
}
