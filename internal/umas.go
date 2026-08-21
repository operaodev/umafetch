package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Uma struct {
	ID        int    `json:"id"`
	OutfitID  *int   `json:"outfit_id,omitempty"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Order     int    `json:"order"`
	Profile   string `json:"profile"`
	Slogan    string `json:"slogan"`
	Image     string `json:"image"`
	MainColor string `json:"main_color"`
	SubColor  string `json:"sub_color"`
}

func (u *Uma) ImageUrl() string {
	if u.OutfitID == nil {
		return fmt.Sprintf("https://media.gametora.com/umamusume/characters/profile/%d.png", u.ID)
	}
	return fmt.Sprintf("https://gametora.com/images/umamusume/characters/chara_stand_%d_%d.png", u.ID, *u.OutfitID)
}

func umasDir() (string, error) {
	localDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(localDir, AppDir, "umas.json"), nil
}

func FindUmas() ([]Uma, error) {
	path, err := umasDir()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var umas []Uma
	if err := json.Unmarshal(data, &umas); err != nil {
		return nil, err
	}

	return umas, nil
}

func UmasExist() bool {
	path, err := umasDir()
	if err != nil {
		return false
	}

	info, err := os.Stat(path)
	return err == nil && info.Size() > 0
}

func SaveUmas() error {
	umas, err := getUmas()
	if err != nil {
		return err
	}

	data := make([]Uma, 0)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	for i := range umas {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()

			randomDelay(200*time.Millisecond, 800*time.Millisecond)
			if err := downloadImage(&umas[i]); err != nil {
				return
			}
			if err := extractColors(&umas[i]); err != nil {
				return
			}
			mu.Lock()
			data = append(data, umas[i])
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	path, err := umasDir()
	if err != nil {
		return err
	}

	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, json, 0644)
}