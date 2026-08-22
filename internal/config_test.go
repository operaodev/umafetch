package internal

import (
	"os"
	"testing"
)

func TestGenerateDefaultConfig(t *testing.T) {
	cfg := GenerateDefaultConfig()

	if cfg.Separator.Width != 52 {
		t.Errorf("Separator.Width = %d, want 52", cfg.Separator.Width)
	}
	if cfg.Separator.Decorator != "\u2500" {
		t.Errorf("Separator.Decorator = %q, want %q", cfg.Separator.Decorator, "\u2500")
	}
	if cfg.Template != TemplateLarge {
		t.Errorf("Template = %q, want %q", cfg.Template, TemplateLarge)
	}

	path, _ := configPath()
	if _, err := os.Stat(path); err != nil {
		t.Errorf("config file was not created: %v", err)
	}
}

func TestConfigLoad(t *testing.T) {
	cfg, err := ConfigLoad()
	if err != nil {
		t.Fatalf("ConfigLoad() error: %v", err)
	}

	if cfg.Separator.Width == 0 {
		t.Error("Separator.Width is 0")
	}
	if cfg.Separator.Decorator == "" {
		t.Error("Separator.Decorator is empty")
	}

	t.Logf("ConfigLoad() template=%s separator=%dx%q", cfg.Template, cfg.Separator.Width, cfg.Separator.Decorator)
}

func TestConfigSave(t *testing.T) {
	cfg := Config{
		Separator: Separator{
			Width:     40,
			Decorator: "=",
		},
		Template: TemplateSmall,
	}

	if err := cfg.ConfigSave(); err != nil {
		t.Fatalf("ConfigSave() error: %v", err)
	}

	loaded, err := ConfigLoad()
	if err != nil {
		t.Fatalf("ConfigLoad() after save error: %v", err)
	}

	if loaded.Separator.Width != 40 {
		t.Errorf("loaded Separator.Width = %d, want 40", loaded.Separator.Width)
	}
	if loaded.Template != TemplateSmall {
		t.Errorf("loaded Template = %q, want %q", loaded.Template, TemplateSmall)
	}
}
