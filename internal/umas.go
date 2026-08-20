package internal

import "fmt"

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
	return fmt.Sprintf("https://media.gametora.com/umamusume/characters/profile/%d.png", u.ID)
}

func (u *Uma) ImageOutfitUrl() string {
	return fmt.Sprintf("https://gametora.com/images/umamusume/characters/chara_stand_%d_%d.png", u.ID, u.OutfitID)
}
