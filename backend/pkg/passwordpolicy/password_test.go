package passwordpolicy

import "testing"

func TestValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{name: "strong", password: "Referral@2026", valid: true},
		{name: "short", password: "Aa1!", valid: false},
		{name: "missing uppercase", password: "referral@2026", valid: false},
		{name: "missing lowercase", password: "REFERRAL@2026", valid: false},
		{name: "missing digit", password: "Referral@Test", valid: false},
		{name: "missing special", password: "Referral2026", valid: false},
		{name: "space is not special", password: "Referral 2026", valid: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := Validate(test.password)
			if (err == nil) != test.valid {
				t.Fatalf("Validate() error = %v, valid = %v", err, test.valid)
			}
		})
	}
}
