package internal

import (
	"testing"
)

func TestUmasExist(t *testing.T) {
	exists := UmasExist()
	t.Logf("UmasExist() = %v", exists)
}

func TestFindUmas(t *testing.T) {
	if !UmasExist() {
		t.Log("umas.json not found, running SaveUmas()...")
		if err := SaveUmas(); err != nil {
			t.Fatalf("SaveUmas() error: %v", err)
		}
	}

	umas, err := FindUmas()
	if err != nil {
		t.Fatalf("FindUmas() error: %v", err)
	}

	t.Logf("FindUmas() returned %d umas", len(umas))

	if len(umas) == 0 {
		t.Fatal("FindUmas() returned empty slice")
	}

	for i, u := range umas {
		if u.ID == 0 {
			t.Errorf("[%d] ID is 0", i)
		}
		if u.Name == "" {
			t.Errorf("[%d] Name is empty", i)
		}
	}
}

func TestFindUmaRandom(t *testing.T) {
	if !UmasExist() {
		if err := SaveUmas(); err != nil {
			t.Fatalf("SaveUmas() error: %v", err)
		}
	}

	cfg := Config{
		Separator: Separator{Width: 52, Decorator: "\u2500"},
		Template:  TemplateLarge,
	}
	cfg.ConfigSave()

	uma, err := FindUma()
	if err != nil {
		t.Fatalf("FindUma() random error: %v", err)
	}

	if uma.ID == 0 {
		t.Error("random uma has ID 0")
	}
	if uma.Name == "" {
		t.Error("random uma has empty Name")
	}

	t.Logf("Random: %s - %s", uma.Name, uma.Title)
}

func TestFindUmaSpecialWeekOutfit0(t *testing.T) {
	if !UmasExist() {
		if err := SaveUmas(); err != nil {
			t.Fatalf("SaveUmas() error: %v", err)
		}
	}

	name := "Special Week"
	outfit := 0
	cfg := Config{
		Theme: Theme{
			Name:   &name,
			Outfit: &outfit,
		},
		Separator: Separator{Width: 52, Decorator: "\u2500"},
		Template:  TemplateLarge,
	}
	cfg.ConfigSave()

	uma, err := FindUma()
	if err != nil {
		t.Fatalf("FindUma() Special Week outfit 0 error: %v", err)
	}

	if uma.Name != "Special Week" {
		t.Errorf("Name = %q, want %q", uma.Name, "Special Week")
	}
	if uma.Order != 0 {
		t.Errorf("Order = %d, want 0", uma.Order)
	}

	t.Logf("Special Week outfit 0: %s - %s (order=%d)", uma.Name, uma.Title, uma.Order)
}

func TestFindUmaSpecialWeekAnyOutfit(t *testing.T) {
	if !UmasExist() {
		if err := SaveUmas(); err != nil {
			t.Fatalf("SaveUmas() error: %v", err)
		}
	}

	name := "Special Week"
	cfg := Config{
		Theme: Theme{
			Name: &name,
		},
		Separator: Separator{Width: 52, Decorator: "\u2500"},
		Template:  TemplateLarge,
	}
	cfg.ConfigSave()

	uma, err := FindUma()
	if err != nil {
		t.Fatalf("FindUma() Special Week any outfit error: %v", err)
	}

	if uma.Name != "Special Week" {
		t.Errorf("Name = %q, want %q", uma.Name, "Special Week")
	}

	t.Logf("Special Week any: %s - %s (order=%d)", uma.Name, uma.Title, uma.Order)
}
