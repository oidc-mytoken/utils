package stringutils

import (
	"testing"
)

func TestWrap(t *testing.T) {
	tests := []struct {
		name   string
		string string
		wrap   string
		result string
	}{
		{
			name:   "string empty",
			string: "",
			wrap:   "$",
			result: "$$",
		},
		{
			name:   "wrap empty",
			string: "string",
			wrap:   "",
			result: "string",
		},
		{
			name:   "normal",
			string: "string",
			wrap:   "$",
			result: "$string$",
		},
		{
			name:   "quoting",
			string: "string",
			wrap:   "\"",
			result: "\"string\"",
		},
		{
			name:   "longer wrap",
			string: "string",
			wrap:   "wrap",
			result: "wrapstringwrap",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				wrapped := Wrap(test.string, test.wrap)
				if wrapped != test.result {
					t.Errorf("Wrap() got = '%s', want '%s'", wrapped, test.result)
				}
			},
		)
	}
}

func TestUnwrap(t *testing.T) {
	tests := []struct {
		name   string
		string string
		wrap   string
		result string
	}{
		{
			name:   "string empty",
			string: "",
			wrap:   "$",
			result: "",
		},
		{
			name:   "wrap empty",
			string: "string",
			wrap:   "",
			result: "string",
		},
		{
			name:   "normal",
			string: "$string$",
			wrap:   "$",
			result: "string",
		},
		{
			name:   "quoted",
			string: "\"string\"",
			wrap:   "\"",
			result: "string",
		},
		{
			name:   "longer wrap",
			string: "wrapstringwrap",
			wrap:   "wrap",
			result: "string",
		},
		{
			name:   "double wrap",
			string: "''string''",
			wrap:   "'",
			result: "'string'",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				unwrapped := Unwrap(test.string, test.wrap)
				if unwrapped != test.result {
					t.Errorf("Wrap() got = '%s', want '%s'", unwrapped, test.result)
				}
			},
		)
	}
}
