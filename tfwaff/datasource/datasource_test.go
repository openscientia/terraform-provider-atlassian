package datasource

import (
	"testing"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		TestName string
		Input    string
		Expect   string
	}{
		{
			TestName: "pascal case",
			Input:    "FooBarBaz",
			Expect:   "nil",
		},
		{
			TestName: "camel case",
			Input:    "fooBarBaz",
			Expect:   "'name' must be in pascal case, e.g., FooBarBaz",
		},
		{
			TestName: "prose",
			Input:    "foo bar baz",
			Expect:   "'name' must be in pascal case, e.g., FooBarBaz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			result := Create("abc", tt.Input, true, true)
			if result != nil {
				if result.Error() != tt.Expect {
					t.Errorf("got %s, expected %s", result.Error(), tt.Expect)
				}
			}
		})
	}
}
