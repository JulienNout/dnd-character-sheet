package storage

import (
	"os"
	"testing"

	characterpkg "modules/dndcharactersheet/internal/domain/character"
)

func TestJSONRepository_CRUD(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "chars-*.json")
	if err != nil {
		t.Fatal(err)
	}
	filename := tmpFile.Name()
	tmpFile.Close()
	// remove the temp file so the storage implementation treats it as non-existent
	// (SingleFileStorage returns an empty structure when file doesn't exist)
	_ = os.Remove(filename)
	defer os.Remove(filename)

	repo := NewJSONRepository(filename)

	// Create domain character
	c := &characterpkg.Character{
		Name:  "TestChar",
		Race:  "Human",
		Class: "Wizard",
		Level: 3,
		Str:   10,
		Dex:   12,
	}

	if err := repo.Save(c); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// GetByID
	got, err := repo.GetByID("TestChar")
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.Name != c.Name || got.Class != c.Class || got.Level != c.Level {
		t.Fatalf("GetByID returned unexpected character: %+v", got)
	}

	// GetAll
	all, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("GetAll expected 1, got %d", len(all))
	}

	// Delete
	if err := repo.Delete("TestChar"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Ensure deleted
	_, err = repo.GetByID("TestChar")
	if err == nil {
		t.Fatalf("expected error after deletion, got nil")
	}
}
