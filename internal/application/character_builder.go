package application

import (
	"sort"
	"strings"

	backgroundModel "modules/dndcharactersheet/internal/background"
	classModel "modules/dndcharactersheet/internal/class"
)

// CharacterBuilder provides helper functions for building characters.
// These functions coordinate multiple inputs (background, class, user choices)
// and don't belong in the domain entity itself.
type CharacterBuilder struct{}

// NewCharacterBuilder creates a new character builder.
func NewCharacterBuilder() *CharacterBuilder {
	return &CharacterBuilder{}
}

// CombineSkillProficiencies combines skill proficiencies from background, class, and user selections.
// It adds class skills first (up to class skill count), then user-selected skills, then background skills.
// The result is sorted alphabetically.
func (cb *CharacterBuilder) CombineSkillProficiencies(
	background backgroundModel.Background,
	class classModel.Class,
	userSkills []string,
) []string {
	var combined []string

	// Add class skills first (up to the class skill count)
	classSkillsAdded := 0
	for _, skill := range class.SkillProficiencies {
		skill = strings.ToLower(strings.TrimSpace(skill))
		if skill != "" && classSkillsAdded < class.SkillCount {
			combined = append(combined, skill)
			classSkillsAdded++
		}
	}

	// Add user-selected skills
	for _, skill := range userSkills {
		skill = strings.ToLower(strings.TrimSpace(skill))
		if skill != "" {
			combined = append(combined, skill)
		}
	}

	// Add background skills (allow duplicates)
	for _, skill := range background.SkillProficiencies {
		skill = strings.ToLower(strings.TrimSpace(skill))
		if skill != "" {
			combined = append(combined, skill)
		}
	}

	// Sort alphabetically
	sort.Strings(combined)

	return combined
}
