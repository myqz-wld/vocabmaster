package appconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const Filename = "config.json"

type LLMSettings struct {
	Adapter  string `json:"adapter"`
	Model    string `json:"model,omitempty"`
	Thinking string `json:"thinking,omitempty"`
}

type Settings struct {
	LLM LLMSettings `json:"llm"`
}

func Path(dataDir string) string {
	return filepath.Join(dataDir, Filename)
}

func Load(dataDir string) (Settings, error) {
	path := Path(dataDir)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Settings{}, nil
		}
		return Settings{}, fmt.Errorf("read config: %w", err)
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return Settings{}, fmt.Errorf("parse config: %w", err)
	}
	return settings, nil
}

func Save(dataDir string, settings Settings) error {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	data = append(data, '\n')

	tmp, err := os.CreateTemp(dataDir, ".config-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary config: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if err := tmp.Chmod(0600); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("set config permissions: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write config: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("sync config: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close config: %w", err)
	}
	if err := os.Rename(tmpPath, Path(dataDir)); err != nil {
		return fmt.Errorf("replace config: %w", err)
	}
	return nil
}
