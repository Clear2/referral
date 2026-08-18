package ent

// UserEntity -.
type UserEntity struct {
	ID            int     `json:"id,omitempty"`
	Name          string  `json:"name,omitempty"`
	Email         string  `json:"email,omitempty"`
	ReferralCode  *string `json:"referralCode,omitempty"`
	CreditBalance int     `json:"creditBalance"`
	CreateTime    string  `json:"createTime,omitempty"`
	UpdateTime    string  `json:"updateTime,omitempty"`
}

// IntoEntity converts ent User to UserEntity.
func (u *User) IntoEntity() *UserEntity {
	return &UserEntity{
		ID:            u.ID,
		Name:          u.Name,
		Email:         u.Email,
		ReferralCode:  u.ReferralCode,
		CreditBalance: u.CreditBalance,
		CreateTime:    u.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:    u.UpdateTime.Format("2006-01-02 15:04:05"),
	}
}
