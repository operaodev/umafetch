package internal

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

const baseURL = "https://umapyoi.net/api/v1"

func randomDelay(min, max time.Duration) {
	time.Sleep(min + time.Duration(rand.Int63n(int64(max-min))))
}

// getUmas fetches the umas from the umapyoi API.
func getUmas() ([]Uma, error) {
	randomDelay(100*time.Millisecond, 500*time.Millisecond)
	resp, err := http.Get(baseURL + "/outfit")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	type outfit struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		CharaGameID int    `json:"chara_game_id"`
	}

	var response []*outfit
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	umas := make(map[int][]Uma)
	for _, u := range response {
		umas[u.CharaGameID] = append(umas[u.CharaGameID], Uma{
			ID:       u.CharaGameID,
			OutfitID: &u.ID,
			Title:    strings.Trim(u.Title, "[]"),
			Order:    u.ID % 100,
		})
	}

	characters := make(map[int]*Uma)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	for id := range umas {
		wg.Add(1)
		sem <- struct{}{}
		go func(id int) {
			defer wg.Done()
			defer func() { <-sem }()

			randomDelay(200*time.Millisecond, 800*time.Millisecond)
			character, err := getCharacter(id)
			if err != nil {
				return
			}
			mu.Lock()
			characters[id] = character
			mu.Unlock()
		}(id)
	}
	wg.Wait()

	for id, list := range umas {
		char, ok := characters[id]
		if !ok {
			continue
		}
		for i := range list {
			list[i].Name = char.Name
			list[i].Profile = char.Profile
			list[i].Slogan = char.Slogan
		}
	}

	var result []Uma
	for id, list := range umas {
		if char, ok := characters[id]; ok {
			result = append(result, *char)
		}
		result = append(result, list...)
	}

	return result, nil
}

// getCharacter fetches a specific character by its ID.
func getCharacter(id int) (*Uma, error) {
	resp, err := http.Get(fmt.Sprintf("%s/character/%d", baseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	type character struct {
		Name    string `json:"name_en"`
		NameJP  string `json:"name_jp"`
		Slogan  string `json:"slogan"`
		Profile string `json:"profile"`
	}

	var char character
	if err := json.NewDecoder(resp.Body).Decode(&char); err != nil {
		return nil, err
	}

	var name string
	if char.Name != "" {
		name = char.Name
	} else {
		name = char.NameJP
	}

	return &Uma{
		ID:      id,
		Name:    name,
		Title:   "Tracen Academy",
		Order:   0,
		Profile: char.Profile,
		Slogan:  char.Slogan,
	}, nil
}
