package testutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// BuildFilePath constructs an absolute file path by searching for the project root
// directory (identified by go.mod) and joining it with the provided fileName.
//
// Returns an error if the current working directory cannot be determined or if
// go.mod is not found in any parent directory.
func BuildFilePath(fileName string) (string, error) {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current working directory: %w", err) //nolint
	}

	for {
		rootPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(rootPath); err == nil {
			fullPath := filepath.Join(dir, fileName)
			return fullPath, nil
		}
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	return "", errors.New("go.mod not found") //nolint
}
