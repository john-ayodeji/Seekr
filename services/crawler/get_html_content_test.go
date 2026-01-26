package crawler

import (
	"strings"
	"testing"
)

func TestGetHTML(t *testing.T) {
	type TestCase struct {
		name     string
		url      string
		contains string
	}

	testCases := []TestCase{
		{
			name:     "Example.com",
			url:      "https://example.com",
			contains: "<title>Example Domain</title>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetHTML(tc.url)
			if err != nil {
				t.Fatal(err)
			}

			if got == "" {
				t.Fatal("expected non-empty HTML")
			}

			if !strings.Contains(got, tc.contains) {
				t.Fatalf(
					"HTML does not contain expected content\nExpected to contain: %q\nGot: %q",
					tc.contains,
					got,
				)
			}
		})
	}
}
