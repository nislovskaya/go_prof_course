package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer srcFile.Close()

	limit, err = setLimit(srcFile, offset, limit)
	if err != nil {
		return fmt.Errorf("error setting limit: %w", err)
	}

	_, err = srcFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error seeking in source file: %w", err)
	}

	destFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error creating destination file: %w", err)
	}
	defer destFile.Close()

	n, err := copyWithProgress(srcFile, destFile, limit)
	if err != nil {
		return fmt.Errorf("error writing to destination file: %w", err)
	}

	fmt.Printf("\nCopied %d bytes from %s to %s\n", n, fromPath, toPath)
	return nil
}

func setLimit(src *os.File, offset, limit int64) (int64, error) {
	fileInfo, err := src.Stat()
	if err != nil {
		return 0, fmt.Errorf("error getting file info: %w", err)
	}

	if !fileInfo.Mode().IsRegular() {
		return 0, ErrUnsupportedFile
	}

	fileSize := fileInfo.Size()
	if offset > fileSize {
		return 0, ErrOffsetExceedsFileSize
	}
	if limit <= 0 || limit > fileSize-offset {
		limit = fileSize - offset
	}

	return limit, nil
}

func copyWithProgress(src io.Reader, dest io.Writer, limit int64) (int64, error) {
	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(src)

	n, err := io.CopyN(dest, barReader, limit)
	bar.SetTotal(limit)
	bar.Finish()

	if err != nil {
		return n, fmt.Errorf("error writing to destination: %w", err)
	}

	return n, nil
}
