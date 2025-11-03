package referencedata

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	backgroundModel "modules/dndcharactersheet/internal/domain/background"
	classModel "modules/dndcharactersheet/internal/domain/class"
	"modules/dndcharactersheet/internal/ports"
)

// JSONBackgroundRepository loads background data from a JSON file.
type JSONBackgroundRepository struct {
	filename string
}

// NewJSONBackgroundRepository creates a new background repository.
func NewJSONBackgroundRepository(filename string) *JSONBackgroundRepository {
	return &JSONBackgroundRepository{filename: filename}
}

// LoadBackgrounds loads all backgrounds from the JSON file.
func (r *JSONBackgroundRepository) LoadBackgrounds() ([]backgroundModel.Background, error) {
	data, err := os.ReadFile(r.filename)
	if err != nil {
		return nil, err
	}

	// Use a storage model with JSON tags for deserialization
	var storageBackgrounds []struct {
		Name               string   `json:"name"`
		SkillProficiencies []string `json:"skill_proficiencies"`
	}

	if err := json.Unmarshal(data, &storageBackgrounds); err != nil {
		return nil, err
	}

	// Convert to domain models
	backgrounds := make([]backgroundModel.Background, len(storageBackgrounds))
	for i, sb := range storageBackgrounds {
		backgrounds[i] = backgroundModel.Background{
			Name:               sb.Name,
			SkillProficiencies: sb.SkillProficiencies,
		}
	}

	return backgrounds, nil
}

// FindByName searches for a background by name (case-insensitive).
func (r *JSONBackgroundRepository) FindByName(name string) (*backgroundModel.Background, error) {
	backgrounds, err := r.LoadBackgrounds()
	if err != nil {
		return nil, err
	}

	for _, bg := range backgrounds {
		if strings.EqualFold(bg.Name, name) {
			return &bg, nil
		}
	}

	return nil, fmt.Errorf("background not found: %s", name)
}

// Ensure JSONBackgroundRepository satisfies the interface.
var _ ports.BackgroundRepository = (*JSONBackgroundRepository)(nil)

// JSONClassRepository loads class data from a JSON file.
type JSONClassRepository struct {
	filename string
}

// NewJSONClassRepository creates a new class repository.
func NewJSONClassRepository(filename string) *JSONClassRepository {
	return &JSONClassRepository{filename: filename}
}

// LoadClasses loads all classes from the JSON file.
func (r *JSONClassRepository) LoadClasses() ([]classModel.Class, error) {
	data, err := os.ReadFile(r.filename)
	if err != nil {
		return nil, err
	}

	// Use a storage model with JSON tags for deserialization
	var storageClasses []struct {
		Name               string   `json:"name"`
		SkillProficiencies []string `json:"skill_proficiencies"`
		SkillCount         int      `json:"skill_count"`
	}

	if err := json.Unmarshal(data, &storageClasses); err != nil {
		return nil, err
	}

	// Convert to domain models
	classes := make([]classModel.Class, len(storageClasses))
	for i, sc := range storageClasses {
		classes[i] = classModel.Class{
			Name:               sc.Name,
			SkillProficiencies: sc.SkillProficiencies,
			SkillCount:         sc.SkillCount,
		}
	}

	return classes, nil
}

// FindByName searches for a class by name (case-insensitive).
func (r *JSONClassRepository) FindByName(name string) (*classModel.Class, error) {
	classes, err := r.LoadClasses()
	if err != nil {
		return nil, err
	}

	for _, cls := range classes {
		if strings.EqualFold(cls.Name, name) {
			return &cls, nil
		}
	}

	return nil, fmt.Errorf("class not found: %s", name)
}

// Ensure JSONClassRepository satisfies the interface.
var _ ports.ClassRepository = (*JSONClassRepository)(nil)
