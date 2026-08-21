package internal

import (
	"testing"
	"time"
)

func TestUmasExist(t *testing.T) {
	start := time.Now()
	exists := UmasExist()
	elapsed := time.Since(start)

	t.Logf("UmasExist() = %v (took %v)", exists, elapsed)
}

func TestFindUmas(t *testing.T) {
	if !UmasExist() {
		t.Log("umas.json not found, running SaveUmas()...")
		start := time.Now()
		err := SaveUmas()
		elapsed := time.Since(start)
		if err != nil {
			t.Fatalf("SaveUmas() error: %v (took %v)", err, elapsed)
		}
		t.Logf("SaveUmas() completed (took %v)", elapsed)
	}

	start := time.Now()
	umas, err := FindUmas()
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("FindUmas() error: %v (took %v)", err, elapsed)
	}

	t.Logf("FindUmas() returned %d umas (took %v)", len(umas), elapsed)

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
