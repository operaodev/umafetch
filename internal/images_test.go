package internal

import (
	"image"
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
	uma := Uma{
		ID:    1001,
		Title: "Special Week",
	}

	if err := downloadImage(&uma); err != nil {
		t.Fatalf("downloadImage() error: %v", err)
	}

	if err := extractColors(&uma); err != nil {
		t.Fatalf("extractColors() error: %v", err)
	}

	if uma.MainColor == "" {
		t.Error("MainColor is empty")
	}
	if uma.SubColor == "" {
		t.Error("SubColor is empty")
	}
	if len(uma.MainColor) != 7 || uma.MainColor[0] != '#' {
		t.Errorf("MainColor is not a valid hex color: %q", uma.MainColor)
	}
	if len(uma.SubColor) != 7 || uma.SubColor[0] != '#' {
		t.Errorf("SubColor is not a valid hex color: %q", uma.SubColor)
	}

	t.Logf("Colors: main=%s sub=%s", uma.MainColor, uma.SubColor)
}

func TestGetColors(t *testing.T) {
	uma := Uma{
		ID:    1001,
		Title: "Special Week",
	}

	if err := downloadImage(&uma); err != nil {
		t.Fatalf("downloadImage() error: %v", err)
	}

	file, err := os.Open(uma.Image)
	if err != nil {
		t.Fatalf("open image error: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatalf("decode image error: %v", err)
	}

	colors := getColors(img)
	if len(colors) == 0 {
		t.Fatal("getColors returned no colors")
	}

	t.Logf("All distinct colors (%d):", len(colors))
	for i, c := range colors {
		t.Logf("  [%d] %s", i, c)
	}
}
