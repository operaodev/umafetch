package internal

import (
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
