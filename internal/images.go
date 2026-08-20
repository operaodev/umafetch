package internal

import (
	"os"
	"os/exec"
	"path/filepath"
)

const AppDir = "umafetch"

func imageDir() (string, error) {
	localDir, err := os.UserCacheDir()
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
	path := filepath.Join(dir, url)

	if _, err := os.Stat(path); err == nil {
		return nil
	}

	cmd := exec.Command("curl", "-s", "-L", "-o", path, url)
	if err := cmd.Run(); err != nil {
		return err
	}

	uma.Image = path

	return nil
}
