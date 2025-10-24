package combat

import (
	"fmt"
	characterModel "modules/dndcharactersheet/internal/character"
)

// FormatSpellcastingStats returns a formatted string for spellcasting stats
func FormatSpellcastingStats(char *characterModel.Character, service *characterModel.CharacterService) string {
	stats := CalculateSpellcastingStats(char, service)
	return fmt.Sprintf("Spellcasting ability: %s\nSpell save DC: %d\nSpell attack bonus: %+d\n", stats.Ability, stats.SpellSaveDC, stats.SpellAttackBonus)
}

// SpellcastingStats holds the calculated spellcasting stats for a character
type SpellcastingStats struct {
	Ability          string
	AbilityMod       int
	SpellSaveDC      int
	SpellAttackBonus int
}

// CalculateSpellcastingStats returns spellcasting ability, DC, and attack bonus for a character
func CalculateSpellcastingStats(char *characterModel.Character, service *characterModel.CharacterService) SpellcastingStats {
	var ability string
	var abilityMod int
	switch char.Class {
	case "wizard":
		ability = "intelligence"
		abilityMod = service.AbilityModifier(char.Int)
	case "cleric", "druid", "ranger":
		ability = "wisdom"
		abilityMod = service.AbilityModifier(char.Wis)
	case "bard", "sorcerer", "warlock", "paladin":
		ability = "charisma"
		abilityMod = service.AbilityModifier(char.Cha)
	default:
		ability = "intelligence"
		abilityMod = service.AbilityModifier(char.Int)
	}
	return SpellcastingStats{
		Ability:          ability,
		AbilityMod:       abilityMod,
		SpellSaveDC:      8 + char.Proficiency + abilityMod,
		SpellAttackBonus: char.Proficiency + abilityMod,
	}
}
