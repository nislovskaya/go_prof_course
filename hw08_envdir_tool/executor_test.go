package main

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	dir := "testdata/env"

	env, err := ReadDir(dir)
	require.NoError(t, err, "unexpected error while reading env")

	testCases := []struct {
		name     string
		cmd      []string
		env      Environment
		expected int
	}{
		{
			name:     "Valid arguments",
			cmd:      []string{"ls"},
			env:      env,
			expected: 0,
		},
		{
			name:     "Valid env",
			cmd:      []string{"ls", "-l"},
			env:      Environment{},
			expected: 0,
		},
		{
			name:     "Command is not provided",
			cmd:      []string{},
			env:      Environment{},
			expected: int(syscall.EINVAL),
		},
		{
			name:     "Invalid command",
			cmd:      []string{"invalid_command"},
			env:      env,
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exitCode := RunCmd(tc.cmd, tc.env)
			require.Equal(t, tc.expected, exitCode)
		})
	}
}
