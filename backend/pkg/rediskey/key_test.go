package rediskey

import "testing"

func TestKey(t *testing.T) {
	tests := []struct {
		name        string
		prefix      string
		environment string
		parts       []string
		want        string
	}{
		{name: "namespaced", prefix: "referral", environment: "test", parts: []string{"session", "42"}, want: "test:referral:session:42"},
		{name: "default prefix", environment: "development", parts: []string{"lock"}, want: "development:prefix:lock"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := New(test.prefix, test.environment).Key(test.parts...); got != test.want {
				t.Fatalf("Key() = %q, want %q", got, test.want)
			}
		})
	}
}
