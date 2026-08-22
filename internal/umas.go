package internal

import (
	"encoding/json"
	"fmt"
	"math/rand"
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

func umasPath() (string, error) {
	localDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(localDir, AppDir, "umas.json"), nil
}

// FindUma encuentra una UMA aleatoria o específica según la configuración.
func FindUma(config Config) (*Uma, error) {
	name := config.Theme.Name
	outfit := config.Theme.Outfit

	umas, err := FindUmas()
	if err != nil {
		return nil, err
	}
	if len(umas) == 0 {
		return nil, fmt.Errorf("no umas found")
	}

	result := make([]Uma, 0, len(umas))
	for _, uma := range umas {
		matchName := name == nil || uma.Name == *name
		matchOutfit := outfit == nil || uma.Order == *outfit
		if matchName && matchOutfit {
			result = append(result, uma)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no uma found with name=%v outfit=%v", name, outfit)
	}

	return &result[rand.Intn(len(result))], nil
}

// FindUmas encuentra las umas desde el archivo de umas.
func FindUmas() ([]Uma, error) {
	path, err := umasPath()
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

// UmasExist verifica si el archivo de umas existe.
func UmasExist() bool {
	path, err := umasPath()
	if err != nil {
		return false
	}

	info, err := os.Stat(path)
	return err == nil && info.Size() > 0
}

// SaveUmas guarda las umas en el archivo de umas.
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

	path, err := umasPath()
	if err != nil {
		return err
	}

	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, json, 0644)
}
