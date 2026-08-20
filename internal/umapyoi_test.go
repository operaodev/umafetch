package internal

import (
	"testing"
)

func TestGetUmas(t *testing.T) {
	umas, err := GetUmas()
	if err != nil {
		t.Fatalf("GetUmas returned error: %v", err)
	}

	if len(umas) == 0 {
		t.Fatal("GetUmas returned empty slice")
	}

	for i, u := range umas {
		if u.ID == 0 {
			t.Errorf("[%d] ID is 0", i)
		}
		if u.Name == "" {
			t.Errorf("[%d] Name is empty (ID=%d)", i, u.ID)
		}
		if u.Title == "" {
			t.Errorf("[%d] Title is empty (ID=%d, Name=%s)", i, u.ID, u.Name)
		}
	}

	t.Logf("Total umas: %d", len(umas))
}
