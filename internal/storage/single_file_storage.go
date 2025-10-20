package storage

import (
	"encoding/json"
	"fmt"
	characterModel "modules/dndcharactersheet/internal/character"
	"os"
)

// SingleFileStorage stores all characters in one JSON file
type SingleFileStorage struct {
	filename string
}

// CharactersFile represents the structure of the single JSON file
type CharactersFile struct {
	Characters []characterModel.Character `json:"characters"`
}

// NewSingleFileStorage creates a new single file storage instance
func NewSingleFileStorage(filename string) *SingleFileStorage {
	return &SingleFileStorage{
		filename: filename,
	}
}

// loadCharactersFile loads all characters from the JSON file
func (sfs *SingleFileStorage) loadCharactersFile() (*CharactersFile, error) {
	charactersFile := &CharactersFile{
		Characters: []characterModel.Character{},
	}

	// Check if file exists
	if _, err := os.Stat(sfs.filename); os.IsNotExist(err) {
		return charactersFile, nil // Return empty structure if file doesn't exist
	}

	// Read the file
	data, err := os.ReadFile(sfs.filename)
	if err != nil {
		return nil, fmt.Errorf("error reading characters file: %v", err)
	}

	// Parse JSON
	err = json.Unmarshal(data, charactersFile)
	if err != nil {
		return nil, fmt.Errorf("error parsing characters file: %v", err)
	}

	return charactersFile, nil
}

// saveCharactersFile saves all characters to the JSON file
func (sfs *SingleFileStorage) saveCharactersFile(charactersFile *CharactersFile) error {
	data, err := json.MarshalIndent(charactersFile, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling characters: %v", err)
	}

	err = os.WriteFile(sfs.filename, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing characters file: %v", err)
	}

	return nil
}

// Save stores a character in the single JSON file
func (sfs *SingleFileStorage) Save(character characterModel.Character) error {
	charactersFile, err := sfs.loadCharactersFile()
	if err != nil {
		return err
	}

	// Check if character already exists (update if it does)
	found := false
	for i, existingChar := range charactersFile.Characters {
		if existingChar.Name == character.Name {
			charactersFile.Characters[i] = character
			found = true
			break
		}
	}

	// If not found, add as new character
	if !found {
		charactersFile.Characters = append(charactersFile.Characters, character)
	}

	return sfs.saveCharactersFile(charactersFile)
}

// Load retrieves a character by name from the single JSON file
func (sfs *SingleFileStorage) Load(name string) (characterModel.Character, error) {
	charactersFile, err := sfs.loadCharactersFile()
	if err != nil {
		return characterModel.Character{}, err
	}

	// Find the character by name
	for _, character := range charactersFile.Characters {
		if character.Name == name {
			return character, nil
		}
	}

	return characterModel.Character{}, fmt.Errorf("character '%s' not found", name)
}

// List returns a summary of all stored characters
func (sfs *SingleFileStorage) List() ([]CharacterSummary, error) {
	charactersFile, err := sfs.loadCharactersFile()
	if err != nil {
		return nil, err
	}

	var summaries []CharacterSummary
	for _, character := range charactersFile.Characters {
		summary := CharacterSummary{
			Name:  character.Name,
			Race:  character.Race,
			Class: character.Class,
			Level: character.Level,
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// Delete removes a character from the single JSON file
func (sfs *SingleFileStorage) Delete(name string) error {
	charactersFile, err := sfs.loadCharactersFile()
	if err != nil {
		return err
	}

	// Find and remove the character
	found := false
	for i, character := range charactersFile.Characters {
		if character.Name == name {
			// Remove character by slicing
			charactersFile.Characters = append(charactersFile.Characters[:i], charactersFile.Characters[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("character '%s' not found", name)
	}

	return sfs.saveCharactersFile(charactersFile)
}
