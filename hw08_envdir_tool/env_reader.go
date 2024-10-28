package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	env := make(Environment, len(dirEntries))

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.Contains(name, "=") {
			return nil, fmt.Errorf("invalid variable name: %s", name)
		}

		filePath := filepath.Join(dir, name)
		value, err := readFileValue(filePath)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
		}

		if value == "" {
			env[name] = EnvValue{NeedRemove: true}
		} else {
			env[name] = EnvValue{Value: value, NeedRemove: false}
		}
	}

	return env, nil
}

func readFileValue(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		value := strings.TrimRight(scanner.Text(), " \t")
		return strings.ReplaceAll(value, "\x00", "\n"), nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}
