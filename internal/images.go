package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const AppDir = "umafetch"

func imageDir() (string, error) {
	localDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(localDir, AppDir, "images")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return dir, nil
}

func downloadImage(uma *Uma) error {
	dir, err := imageDir()
	if err != nil {
		return err
	}

	url := uma.ImageUrl()
	path := filepath.Join(dir, filepath.Base(url))

	if _, err := os.Stat(path); err == nil {
		uma.Image = path
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("image not found: %s (status %d)", url, resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "image/png" {
		return fmt.Errorf("unexpected content type %q for %s", ct, url)
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		os.Remove(path)
		return err
	}

	uma.Image = path
	return nil
}
