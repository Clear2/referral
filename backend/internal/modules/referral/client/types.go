package client

import "time"

const RewardCredits = 100

type RegisterInput struct {
	Code  string `json:"code" binding:"required,len=8,alphanum"`
	Name  string `json:"name" binding:"required,min=1,max=100"`
	Email string `json:"email" binding:"required,email,max=254"`
}
type UserView struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	ReferralCode  *string `json:"referralCode,omitempty"`
	CreditBalance int     `json:"creditBalance"`
}
type InvitationView struct {
	Code string `json:"code"`
	URL  string `json:"url"`
}
type RegistrationView struct {
	Invitee              UserView `json:"invitee"`
	Reward               int      `json:"reward"`
	InviterCreditBalance int      `json:"inviterCreditBalance"`
}
type ReferralView struct {
	ID        int       `json:"id"`
	Invitee   UserView  `json:"invitee"`
	Reward    int       `json:"reward"`
	CreatedAt time.Time `json:"createdAt"`
}
type CreditTransactionView struct {
	ID         int       `json:"id"`
	Amount     int       `json:"amount"`
	Reason     string    `json:"reason"`
	ReferralID int       `json:"referralId"`
	CreatedAt  time.Time `json:"createdAt"`
}
type DashboardView struct {
	User                UserView                `json:"user"`
	SuccessfulReferrals int                     `json:"successfulReferrals"`
	TotalCreditsEarned  int                     `json:"totalCreditsEarned"`
	Referrals           []ReferralView          `json:"referrals"`
	CreditTransactions  []CreditTransactionView `json:"creditTransactions"`
}
type InvitationAPIResponse struct {
	Code    int             `json:"code"`
	Payload *InvitationView `json:"payload"`
	Message string          `json:"message"`
}
type RegistrationAPIResponse struct {
	Code    int               `json:"code"`
	Payload *RegistrationView `json:"payload"`
	Message string            `json:"message"`
}
type DashboardAPIResponse struct {
	Code    int            `json:"code"`
	Payload *DashboardView `json:"payload"`
	Message string         `json:"message"`
}
