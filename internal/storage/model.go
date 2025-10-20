package storage

import characterModel "modules/dndcharactersheet/internal/character"

// CharacterSummary contains basic info about a character for listing
type CharacterSummary struct {
	Name  string `json:"name"`
	Race  string `json:"race"`
	Class string `json:"class"`
	Level int    `json:"level"`
}

// CharacterStorage defines the interface for character persistence operations
type CharacterStorage interface {
	// Save stores a character to persistent storage
	Save(character characterModel.Character) error

	// Load retrieves a character by name from persistent storage
	Load(name string) (characterModel.Character, error)

	// List returns a summary of all stored characters
	List() ([]CharacterSummary, error)

	// Delete removes a character from persistent storage
	Delete(name string) error
}
