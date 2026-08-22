package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateDefaultTemplateLarge(t *testing.T) {
	if err := GenerateDefaultTemplateLarge(); err != nil {
		t.Fatalf("GenerateDefaultTemplateLarge() error: %v", err)
	}

	dir, _ := templateDir()
	path := filepath.Join(dir, "config_large.jsonc")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("template file was not created: %v", err)
	}

	data, _ := os.ReadFile(path)
	if len(data) == 0 {
		t.Error("template file is empty")
	}
}

func TestGenerateDefaultTemplateSmall(t *testing.T) {
	if err := GenerateDefaultTemplateSmall(); err != nil {
		t.Fatalf("GenerateDefaultTemplateSmall() error: %v", err)
	}

	dir, _ := templateDir()
	path := filepath.Join(dir, "config_small.jsonc")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("template file was not created: %v", err)
	}

	data, _ := os.ReadFile(path)
	if len(data) == 0 {
		t.Error("template file is empty")
	}
}

func TestFindTemplate(t *testing.T) {
	config := Config{
		Template: TemplateLarge,
	}
	content, err := FindTemplate(config)
	if err != nil {
		t.Fatalf("FindTemplate() error: %v", err)
	}

	if len(content) == 0 {
		t.Error("FindTemplate() returned empty content")
	}
}
