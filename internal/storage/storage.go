package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/anselmo/burrow/internal/mission"
	"github.com/anselmo/burrow/internal/pet"
)

type State struct {
	Pet      *pet.Pet          `json:"pet"`
	Mission  *mission.Mission  `json:"mission,omitempty"`
	Items    []mission.Item    `json:"items,omitempty"`
}

type Log struct {
	Entries []pet.LogEntry `json:"entries"`
}

func dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(home, ".burrow")
	return d, os.MkdirAll(d, 0700)
}

func LoadState() (*State, error) {
	d, err := dir()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(d, "pet.json"))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var s State
	return &s, json.Unmarshal(data, &s)
}

func SaveState(s *State) error {
	d, err := dir()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(d, "pet.json"), data, 0600)
}

func LoadLog() (*Log, error) {
	d, err := dir()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(d, "log.json"))
	if os.IsNotExist(err) {
		return &Log{}, nil
	}
	if err != nil {
		return nil, err
	}
	var l Log
	return &l, json.Unmarshal(data, &l)
}

func SaveLog(l *Log) error {
	d, err := dir()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(d, "log.json"), data, 0600)
}
