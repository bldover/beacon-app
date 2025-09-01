package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func GetCacheFilePath(filename string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	execDir := filepath.Dir(execPath)
	return filepath.Join(execDir, filename), nil
}

func WriteJSONFile(filePath string, data interface{}) error {
	tempPath := filePath + ".tmp"

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = os.WriteFile(tempPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	err = os.Rename(tempPath, filePath)
	if err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

func ReadJSONFile(filePath string, target interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func IsFileStale(filePath string, maxAge time.Duration) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return true
	}

	return time.Since(info.ModTime()) > maxAge
}
