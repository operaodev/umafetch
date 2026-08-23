package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Template  Template  `json:"template"`
	Separator Separator `json:"separator"`
	Theme     Theme     `json:"theme"`
}

type Theme struct {
	Name   *string `json:"name"`   // nil = random, Special Week, Silence Suzuka
	Outfit *int    `json:"outfit"` // nil = random, 0 = Tracen Academy, 1 = Outfit1, 2 = Outfit2
}

type Separator struct {
	Width     int    `json:"width"`
	Decorator string `json:"decorator"`
}

func (s Separator) Build() string {
	return strings.Repeat(s.Decorator, s.Width)
}

func configPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, AppDir, "config.json"), nil
}

// GenerateDefaultConfig generates a default configuration and saves it to the config file.
func GenerateDefaultConfig() Config {
	config := Config{
		Separator: Separator{
			Width:     52,
			Decorator: "\u2500",
		},
		Template: TemplateLarge,
	}
	config.ConfigSave()
	return config
}

// ConfigLoad loads the configuration from the config file.
func ConfigLoad() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return GenerateDefaultConfig(), nil
	}
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return GenerateDefaultConfig(), err
	}

	return config, nil
}

// ConfigSave saves the configuration to the config file.
func (c Config) ConfigSave() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
