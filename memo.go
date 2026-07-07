package main

import (
	"encoding/json"
	"fmt"
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
		return []Memo{}, nil
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

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("JSONファイルが開けません", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	encoder.SetIndent("", "  ")

	if err = encoder.Encode(memos); err != nil {
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
