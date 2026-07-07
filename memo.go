package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Memo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func LoadMemos(path string) ([]Memo, error) {

	path, err := expandPath(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Memo{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var memos []Memo
	if err := json.NewDecoder(file).Decode(&memos); err != nil {
		return nil, err
	}

	return memos, nil
}

func SaveMemos(path string, memos []Memo) error {

	path, err := expandPath(path)
	if err != nil {
		return err
	}

	dirPath := filepath.Dir(path)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(dirPath, ".memos-*.json")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	encoder := json.NewEncoder(tmpFile)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(memos); err != nil {
		tmpFile.Close()
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}

	return nil
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		path = filepath.Join(home, path[2:])
	}
	return path, nil
}
