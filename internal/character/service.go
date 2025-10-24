package characterModel

import (
	backgroundModel "modules/dndcharactersheet/internal/background"
	classModel "modules/dndcharactersheet/internal/class"
	"sort"
	"strings"
)

type CharacterService struct{}

func NewCharacterService() *CharacterService {
	return &CharacterService{}
}
func (cs *CharacterService) GetProficiencyBonus(level int) int {
	switch {
	case level >= 1 && level <= 4:
		return 2
	case level >= 5 && level <= 8:
		return 3
	case level >= 9 && level <= 12:
		return 4
	case level >= 13 && level <= 16:
		return 5
	case level >= 17:
		return 6
	default:
		return 0
	}
}

func (cs *CharacterService) AbilityModifier(score int) int {
	result := (score - 10) / 2
	if (score-10)%2 < 0 {
		result--
	}
	return result
}

func (cs *CharacterService) ApplyRacialBonuses(character *Character) {
	switch strings.ToLower(character.Race) {
	case "dwarf":
		character.Con += 2
	case "hill dwarf":
		character.Con += 2
		character.Wis += 1
	case "elf":
		character.Dex += 2
	case "high elf":
		character.Dex += 2
		character.Int += 1
	case "halfling":
		character.Dex += 2
	case "lightfoot halfling", "lightfoot":
		character.Dex += 2
		character.Cha += 1
	case "human":
		character.Str += 1
		character.Dex += 1
		character.Con += 1
		character.Int += 1
		character.Wis += 1
		character.Cha += 1
	case "dragonborn":
		character.Str += 2
		character.Cha += 1
	case "gnome":
		character.Int += 2
	case "rock gnome":
		character.Int += 2
		character.Con += 1
	case "half-elf":
		character.Cha += 2
		character.Dex += 1
		character.Con += 1
	case "half orc":
		character.Str += 2
		character.Con += 1
	case "tiefling":
		character.Int += 1
		character.Cha += 2
	}
}

func (cs *CharacterService) CombineSkillProficiencies(background backgroundModel.Background, class classModel.Class, userSkills []string) []string {
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
