package main

import (
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	testCases := []struct {
		name     string
		files    []string
		expected Environment
	}{
		{
			name: "Normal files",
			files: []string{
				"BAR",
				"HELLO",
			},
			expected: Environment{
				"BAR": {
					Value:      `bar`,
					NeedRemove: false,
				},
				"HELLO": {
					Value:      `"hello"`,
					NeedRemove: false,
				},
			},
		},
		{
			name: "Empty file",
			files: []string{
				"EMPTY",
			},
			expected: Environment{
				"EMPTY": {NeedRemove: true},
			},
		},
		{
			name: "Unset variable",
			files: []string{
				"UNSET",
			},
			expected: Environment{
				"UNSET": {NeedRemove: true},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := filepath.Join("testdata", "env")

			env, err := ReadDir(dir)
			if err != nil {
				t.Fatalf("expected no error but got %v", err)
			}

			for _, file := range tc.files {
				value, exists := env[file]
				if !exists {
					t.Errorf("expected %s not found in env", file)
					continue
				}
				expectedValue := tc.expected[file]
				if value != expectedValue {
					t.Errorf("for %s: expected %+v, got %+v", file, expectedValue, value)
				}
			}
		})
	}
}
