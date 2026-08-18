package auth

import "time"

type GoogleProfile struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	VerifiedEmail bool   `json:"verified_email"`
}

type UserView struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	ReferralCode  *string `json:"referralCode,omitempty"`
	CreditBalance int     `json:"creditBalance"`
}

type LoginView struct {
	AccessToken string    `json:"accessToken"`
	TokenType   string    `json:"tokenType"`
	ExpiresAt   time.Time `json:"expiresAt"`
	IsNewUser   bool      `json:"isNewUser"`
	User        UserView  `json:"user"`
	session     SessionTokens
}

type LoginAPIResponse struct {
	Code    int        `json:"code"`
	Payload *LoginView `json:"payload"`
	Message string     `json:"message"`
}
