package todo

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const todoDir = ".gtr"
const todoFile = "todos.json"

type Storage struct {
	filePath string
}

func NewStorage() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	gtrDir := filepath.Join(homeDir, todoDir)
	if err := os.MkdirAll(gtrDir, 0755); err != nil {
		return nil, err
	}

	filePath := filepath.Join(gtrDir, todoFile)
	return &Storage{filePath: filePath}, nil
}

func (s *Storage) Load() (*TaskStore, error) {
	data, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &TaskStore{Tasks: []Task{}}, nil
		}
		return nil, err
	}

	var store TaskStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}

	return &store, nil
}

func (s *Storage) Save(store *TaskStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.filePath, data, 0644)
}

func (s *Storage) GetFilePath() string {
	return s.filePath
}
