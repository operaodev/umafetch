package internal

import (
	"strings"
	"testing"
)

func TestSeparatorModule(t *testing.T) {
	cfg, err := ConfigLoad()
	if err != nil {
		t.Fatalf("ConfigLoad() error: %v", err)
	}

	result := separatorModule("#FFE7CB", cfg)

	if !strings.Contains(result, `"type":"custom"`) {
		t.Error("missing type:custom")
	}
	if !strings.Contains(result, "#FFE7CB") {
		t.Error("missing color")
	}
	if !strings.Contains(result, cfg.Separator.Decorator) {
		t.Error("missing separator decorator")
	}

	t.Logf("separatorModule: %s", result)
}

func TestBuildTextBlock(t *testing.T) {
	result := buildTextBlock("Slogan", "Run with me! I will become the best!", "#CCCBF9", 20, 30)

	if !strings.Contains(result, `"key":"Slogan"`) {
		t.Error("missing key")
	}
	if !strings.Contains(result, "#CCCBF9") {
		t.Error("missing color")
	}

	t.Logf("buildTextBlock: %s", result)
}

func TestWrapLines(t *testing.T) {
	lines := wrapLines("Short text", 50, 50)
	if len(lines) != 1 {
		t.Errorf("short text should be 1 line, got %d", len(lines))
	}

	longText := "This is a very long text that should be split into multiple lines when it exceeds the limit"
	lines = wrapLines(longText, 20, 30)
	if len(lines) < 2 {
		t.Errorf("long text should be multiple lines, got %d", len(lines))
	}

	t.Logf("wrapLines (%d lines): %v", len(lines), lines)
}

func TestEscJSON(t *testing.T) {
	result := escJSON(`Hello "World"`)
	if !strings.Contains(result, `\"`) {
		t.Errorf("escJSON not escaping quotes: %s", result)
	}
}

func TestRenderUma(t *testing.T) {
	if !UmasExist() {
		if err := SaveUmas(); err != nil {
			t.Fatalf("SaveUmas() error: %v", err)
		}
	}

	path, err := RenderUma()
	if err != nil {
		t.Fatalf("RenderUma() error: %v", err)
	}

	t.Logf("RenderUma() output: %s", path)
}
