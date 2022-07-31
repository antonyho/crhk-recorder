package resolver

import (
	"testing"
)

func TestCloudfrontCookie_Assigned(t *testing.T) {
	cases := []struct {
		testName string
		cookies  CloudfrontCookie
		wanted   bool
	}{
		{
			"complete",
			CloudfrontCookie{"policy-dummy", "keypair-dummy", "sig-dummy"},
			true,
		},
		{
			"empty policy",
			CloudfrontCookie{"", "keypair-dummy", "sig-dummy"},
			false,
		},
		{
			"empty keypair",
			CloudfrontCookie{"policy-dummy", "", "sig-dummy"},
			false,
		},
		{
			"empty signature",
			CloudfrontCookie{"policy-dummy", "keypair-dummy", ""},
			false,
		},
	}

	for _, c := range cases {
		got := c.cookies.Assigned()
		if got != c.wanted {
			t.Logf("Wanted: %v Got: %v", c.wanted, got)
		}
	}
}
