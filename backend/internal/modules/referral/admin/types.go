package admin

import "time"

type ListInput struct {
	Page          int        `form:"page" binding:"omitempty,min=1"`
	PageSize      int        `form:"page_size" binding:"omitempty,min=1,max=100"`
	UserID        int        `form:"user_id" binding:"omitempty,min=1"`
	Email         string     `form:"email" binding:"omitempty,max=254"`
	CreatedAtFrom *time.Time `form:"created_at_from" time_format:"2006-01-02"`
	CreatedAtTo   *time.Time `form:"created_at_to" time_format:"2006-01-02"`
}

type UserView struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	ReferralCode  *string `json:"referralCode,omitempty"`
	CreditBalance int     `json:"creditBalance"`
}
type PaginationView struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}
type ReferralView struct {
	ID        int       `json:"id"`
	Inviter   UserView  `json:"inviter"`
	Invitee   UserView  `json:"invitee"`
	Reward    int       `json:"reward"`
	CreatedAt time.Time `json:"createdAt"`
}
type ReferralListView struct {
	Items      []ReferralView `json:"items"`
	Pagination PaginationView `json:"pagination"`
}
type CreditTransactionView struct {
	ID         int       `json:"id"`
	User       UserView  `json:"user"`
	Amount     int       `json:"amount"`
	Reason     string    `json:"reason"`
	ReferralID int       `json:"referralId"`
	CreatedAt  time.Time `json:"createdAt"`
}
type CreditTransactionListView struct {
	Items      []CreditTransactionView `json:"items"`
	Pagination PaginationView          `json:"pagination"`
}
type StatsView struct {
	TotalUsers          int `json:"totalUsers"`
	TotalInviters       int `json:"totalInviters"`
	TotalReferrals      int `json:"totalReferrals"`
	TotalCreditsAwarded int `json:"totalCreditsAwarded"`
}

type ReferralListAPIResponse struct {
	Code    int               `json:"code"`
	Payload *ReferralListView `json:"payload"`
	Message string            `json:"message"`
}
type CreditTransactionListAPIResponse struct {
	Code    int                        `json:"code"`
	Payload *CreditTransactionListView `json:"payload"`
	Message string                     `json:"message"`
}
type StatsAPIResponse struct {
	Code    int        `json:"code"`
	Payload *StatsView `json:"payload"`
	Message string     `json:"message"`
}
