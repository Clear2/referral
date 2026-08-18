package auth

import "testing"

func TestSafeReturnPath(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct{ input, want string }{
		"root":           {input: "/", want: "/"},
		"path":           {input: "/dashboard?tab=credits", want: "/dashboard?tab=credits"},
		"empty":          {input: "", want: "/"},
		"absolute":       {input: "https://evil.example", want: "/"},
		"protocol local": {input: "//evil.example", want: "/"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := safeReturnPath(test.input); got != test.want {
				t.Fatalf("safeReturnPath(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}

func TestControllerSafeReturnURL(t *testing.T) {
	t.Parallel()
	controller := &Controller{allowedOrigins: map[string]struct{}{"http://localhost:5173": {}}}

	for name, test := range map[string]struct{ input, want string }{
		"allowed origin": {input: "http://localhost:5173/dashboard?tab=credits", want: "http://localhost:5173/dashboard?tab=credits"},
		"relative path":  {input: "/dashboard", want: "/dashboard"},
		"foreign origin": {input: "https://evil.example/dashboard", want: "/"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := controller.safeReturnURL(test.input); got != test.want {
				t.Fatalf("safeReturnURL(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}
