package internal

import (
	_ "embed"
	"os"
	"path/filepath"
)

type Template string

const (
	TemplateSmall Template = "small"
	TemplateLarge Template = "large"
)

//go:embed templates/config_large.jsonc
var templateLarge string

//go:embed templates/config_small.jsonc
var templateSmall string

// GenerateDefaultTemplateSmall genera el template small por defecto.
func GenerateDefaultTemplateSmall() error {
	return writeDefaultTemplate("config_small.jsonc", templateSmall)
}

// GenerateDefaultTemplateLarge genera el template large por defecto.
func GenerateDefaultTemplateLarge() error {
	return writeDefaultTemplate("config_large.jsonc", templateLarge)
}

// FindTemplate encuentra el template según la configuración.
func FindTemplate(config Config) (string, error) {
	var name string
	if config.Template == TemplateLarge {
		name = "config_large.jsonc"
	} else {
		name = "config_small.jsonc"
	}

	dir, err := templateDir()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func templateDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(configDir, AppDir, "layouts")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func writeDefaultTemplate(name, content string) error {
	dir, err := templateDir()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644)
}
