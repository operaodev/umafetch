package internal

import (
	"fmt"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestImageDir(t *testing.T) {
	dir, err := imageDir()
	if err != nil {
		t.Fatalf("imageDir() error: %v", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("imageDir() path does not exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("imageDir() path is not a directory: %s", dir)
	}
}

func TestDownloadImage(t *testing.T) {
	outfitID := 100103
	uma := Uma{
		ID:       1001,
		OutfitID: &outfitID,
		Title:    "Supreme Commander of the Rising Sun",
	}

	if err := downloadImage(&uma); err != nil {
		t.Fatalf("downloadImage() error: %v", err)
	}

	if uma.Image == "" {
		t.Fatal("uma.Image was not set after download")
	}

	if _, err := os.Stat(uma.Image); err != nil {
		t.Fatalf("downloaded file does not exist: %s", uma.Image)
	}

	t.Logf("Image downloaded to: %s", uma.Image)
}

func TestDownloadImageProfile(t *testing.T) {
	uma := Uma{
		ID:    1001,
		Title: "Special Week",
	}

	if err := downloadImage(&uma); err != nil {
		t.Fatalf("downloadImage() error: %v", err)
	}

	if uma.Image == "" {
		t.Fatal("uma.Image was not set after download")
	}

	if _, err := os.Stat(uma.Image); err != nil {
		t.Fatalf("downloaded file does not exist: %s", uma.Image)
	}

	t.Logf("Image downloaded to: %s", uma.Image)
}

func TestDownloadImageCached(t *testing.T) {
	outfitID := 99999
	uma := Uma{
		ID:       1001,
		OutfitID: &outfitID,
		Title:    "Test Cached",
	}

	dir, err := imageDir()
	if err != nil {
		t.Fatalf("imageDir() error: %v", err)
	}

	cacheURL := uma.ImageUrl()
	cachePath := filepath.Join(dir, filepath.Base(cacheURL))

	os.WriteFile(cachePath, []byte("cached"), 0644)

	if err := downloadImage(&uma); err != nil {
		t.Fatalf("downloadImage() error: %v", err)
	}

	if uma.Image != cachePath {
		t.Errorf("uma.Image = %q, want %q", uma.Image, cachePath)
	}

	content, _ := os.ReadFile(uma.Image)
	if string(content) != "cached" {
		t.Error("file was re-downloaded instead of using cache")
	}
}

func TestDownloadImageInvalid(t *testing.T) {
	uma := Uma{
		ID:    9999,
		Title: "Nonexistent",
	}

	err := downloadImage(&uma)
	if err == nil {
		t.Fatal("expected error for invalid image, got nil")
	}

	if uma.Image != "" {
		t.Errorf("uma.Image should be empty on error, got %q", uma.Image)
	}

	t.Logf("Correctly caught error: %v", err)
}

func TestExtractColors(t *testing.T) {
	umas := []Uma{
		{ID: 1006, OutfitID: new(100602), Title: "Oguri Cap"},
		{ID: 1007, OutfitID: new(100703), Title: "Gold Ship"},
	}

	for _, uma := range umas {
		if err := downloadImage(&uma); err != nil {
			t.Fatalf("downloadImage(%s) error: %v", uma.Title, err)
		}

		if err := extractColors(&uma); err != nil {
			t.Fatalf("extractColors(%s) error: %v", uma.Title, err)
		}

		if uma.MainColor == "" {
			t.Errorf("%s: MainColor is empty", uma.Title)
		}
		if uma.SubColor == "" {
			t.Errorf("%s: SubColor is empty", uma.Title)
		}
		if len(uma.MainColor) != 7 || uma.MainColor[0] != '#' {
			t.Errorf("%s: MainColor is not a valid hex color: %q", uma.Title, uma.MainColor)
		}
		if len(uma.SubColor) != 7 || uma.SubColor[0] != '#' {
			t.Errorf("%s: SubColor is not a valid hex color: %q", uma.Title, uma.SubColor)
		}

		t.Logf("%s: main=%s sub=%s", uma.Title, uma.MainColor, uma.SubColor)
		paintColor(t, "main", uma.MainColor)
		paintColor(t, "sub", uma.SubColor)
	}
}

func ptr(v int) *int { return &v }

func paintColor(t *testing.T, label, hex string) {
	var r, g, b byte
	fmt.Sscanf(hex, "#%02X%02X%02X", &r, &g, &b)
	t.Logf("\033[48;2;%d;%d;%dm  %s %s  \033[0m", r, g, b, label, hex)
}
