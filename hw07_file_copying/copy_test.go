package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func performCopyAndCheck(t *testing.T, from string, to string, offset int64, limit int64, expectedFile string) {
	t.Helper()

	err := Copy(from, to, offset, limit)
	if err != nil {
		t.Fatalf("Copy() error = %v", err)
	}

	expectedContent, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read expected file: %v", err)
	}

	actualContent, err := os.ReadFile(to)
	if err != nil {
		t.Fatalf("Failed to read actual file after copy: %v", err)
	}

	if !bytes.Equal(expectedContent, actualContent) {
		t.Errorf("Content mismatch; expected %q but got %q", expectedContent, actualContent)
	}

	err = os.Remove(to)
	if err != nil {
		t.Fatalf("Failed to remove file: %v", err)
	}
}

func createTempDir(t *testing.T) string {
	t.Helper()

	tmpDir := "tmp"
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	return tmpDir
}

func TestCopy(t *testing.T) {
	tests := []struct {
		name         string
		from         string
		to           string
		offset       int64
		limit        int64
		expectErr    bool
		expectedFile string
	}{
		{
			"Copy full file",
			"testdata/input.txt",
			"tmp/copy_input_full.txt",
			0,
			0,
			false,
			"testdata/out_offset0_limit0.txt",
		},
		{
			"Copy with limit 10",
			"testdata/input.txt",
			"tmp/copy_input_limit10.txt",
			0,
			10,
			false,
			"testdata/out_offset0_limit10.txt",
		},
		{
			"Copy with limit 1000",
			"testdata/input.txt",
			"tmp/copy_input_limit1000.txt",
			0,
			1000,
			false,
			"testdata/out_offset0_limit1000.txt",
		},
		{
			"Copy with limit 10000",
			"testdata/input.txt",
			"tmp/copy_input_limit10000.txt",
			0,
			10000,
			false,
			"testdata/out_offset0_limit10000.txt",
		},
		{
			"Copy with offset 100 and limit 1000",
			"testdata/input.txt",
			"tmp/copy_input_offset100_limit1000.txt",
			100,
			1000,
			false,
			"testdata/out_offset100_limit1000.txt",
		},
		{
			"Copy with offset 6000 and limit 1000",
			"testdata/input.txt",
			"tmp/copy_input_offset6000_limit1000.txt",
			6000,
			1000,
			false,
			"testdata/out_offset6000_limit1000.txt",
		},
		{
			"Invalid offset",
			"testdata/input.txt",
			"tmp/copy_invalid_offset_input.txt",
			10000,
			0,
			true,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := createTempDir(t)
			tt.to = filepath.Join(tmpDir, filepath.Base(tt.to))

			if tt.expectErr {
				err := Copy(tt.from, tt.to, tt.offset, tt.limit)
				if !errors.Is(err, ErrOffsetExceedsFileSize) {
					t.Errorf("Expected ErrOffsetExceedsFileSize but got: %v", err)
				}
			} else {
				performCopyAndCheck(t, tt.from, tt.to, tt.offset, tt.limit, tt.expectedFile)
			}
		})
	}
}

func TestUnsupportedFile(t *testing.T) {
	tempDir := t.TempDir()
	dirPath := filepath.Join(tempDir, "testdata")

	if mkErr := os.Mkdir(dirPath, 0o755); mkErr != nil {
		t.Fatalf("Failed to create directory for test :%s", mkErr.Error())
	}

	err := Copy(dirPath, filepath.Join(tempDir, "output.txt"), 0, 0)

	if !errors.Is(err, ErrUnsupportedFile) {
		t.Errorf("Expected ErrUnsupportedFile ,got :%s", err.Error())
	}
}
