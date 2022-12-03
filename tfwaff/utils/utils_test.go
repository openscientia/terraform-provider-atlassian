package utils

import (
	"strconv"
	"testing"
)

func TestIsPascalCase(t *testing.T) {
	tests := []struct {
		TestName string
		Input    string
		Expect   string
	}{
		{
			TestName: "correct",
			Input:    "FooBarBaz",
			Expect:   "true",
		},
		{
			TestName: "incorrect",
			Input:    "fooBarBaz",
			Expect:   "false",
		},
		{
			TestName: "incorrect",
			Input:    "foo_bar",
			Expect:   "false",
		},
		{
			TestName: "incorrect",
			Input:    "foo-bar-baz",
			Expect:   "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			result := strconv.FormatBool(IsPascalCase(tt.Input))
			if result != tt.Expect {
				t.Errorf("got %s, expected %s", result, tt.Expect)
			}
		})
	}
}

func TestGetSnakeCase(t *testing.T) {
	tests := []struct {
		TestName string
		Input    string
		Expect   string
	}{
		{
			TestName: "prose",
			Input:    "foo bar baz",
			Expect:   "foo bar baz",
		},
		{
			TestName: "pascal case",
			Input:    "FooBarBaz",
			Expect:   "foo_bar_baz",
		},
		{
			TestName: "camel case",
			Input:    "fooBarBaz",
			Expect:   "fooBarBaz",
		},
		{
			TestName: "kebab case",
			Input:    "foo-bar-baz",
			Expect:   "foo-bar-baz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			result := GetSnakeCase(tt.Input)
			if result != tt.Expect {
				t.Errorf("got %s, expected %s", result, tt.Expect)
			}
		})
	}
}

func TestGetTitleCase(t *testing.T) {
	tests := []struct {
		TestName string
		Input    string
		Expect   string
	}{
		{
			TestName: "prose",
			Input:    "foo bar baz",
			Expect:   "Foo Bar Baz",
		},
		{
			TestName: "camel case",
			Input:    "fooBarBaz",
			Expect:   "Foobarbaz",
		},
		{
			TestName: "title case all caps",
			Input:    "FOO BAR BAZ",
			Expect:   "Foo Bar Baz",
		},
		{
			TestName: "pascal case",
			Input:    "FooBarBaz",
			Expect:   "Foobarbaz",
		},
		{
			TestName: "kebab case",
			Input:    "foo-bar-baz",
			Expect:   "Foo Bar Baz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			result := GetTitleCase(tt.Input)
			if result != tt.Expect {
				t.Errorf("got %s, expected %s", result, tt.Expect)
			}
		})
	}
}
