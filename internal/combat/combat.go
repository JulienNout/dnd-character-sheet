package combat

import (
	"modules/dndcharactersheet/internal/api"
	characterModel "modules/dndcharactersheet/internal/character"
	"strings"
)

// CalculateArmorClass returns the armor class for a character using real-time API enrichment.
func CalculateArmorClass(char *characterModel.Character, service *characterModel.CharacterService) int {
	// Barbarian Unarmored Defense: AC = 10 + Dex mod + Con mod (+ shield bonus if equipped) if no armor
	if strings.ToLower(char.Class) == "barbarian" && char.Armor == "" {
		ac := 10 + service.AbilityModifier(char.Dex) + service.AbilityModifier(char.Con)
		if char.Shield != "" {
			ac += 2 // D&D 5e shield bonus
		}
		return ac
	}
	// Monk Unarmored Defense: AC = 10 + Dex mod + Wis mod if no armor and no shield
	if strings.ToLower(char.Class) == "monk" && char.Armor == "" && char.Shield == "" {
		return 10 + service.AbilityModifier(char.Dex) + service.AbilityModifier(char.Wis)
	}

	baseAC := 10
	dexMod := service.AbilityModifier(char.Dex)

	// Get armor AC from API if equipped, fallback to hardcoded values if API fails
	if char.Armor != "" {
		apiIndex := api.ToAPIIndex(char.Armor)
		armor, err := api.GetArmor(apiIndex)
		if err == nil && armor != nil {
			baseAC = armor.ArmorClass.Base
			if armor.ArmorClass.DexBonus {
				baseAC += dexMod
			}
		} else {
			// Fallback: hardcoded AC values for common armor
			switch apiIndex {
			case "padded":
				baseAC = 11 + dexMod
			case "leather-armor":
				baseAC = 11 + dexMod
			case "studded-leather":
				baseAC = 12 + dexMod
			case "hide-armor":
				baseAC = 12 + min(dexMod, 2)
			case "chain-shirt":
				baseAC = 13 + min(dexMod, 2)
			case "scale-mail":
				baseAC = 14 + min(dexMod, 2)
			case "breastplate":
				baseAC = 14 + min(dexMod, 2)
			case "half-plate":
				baseAC = 15 + min(dexMod, 2)
			case "ring-mail":
				baseAC = 14
			case "chain-mail":
				baseAC = 16
			case "splint":
				baseAC = 17
			case "plate-armor":
				baseAC = 18
			default:
				baseAC += dexMod
			}
		}
	} else {
		baseAC += dexMod
	}

	// Add shield bonus if equipped (assume +2 for D&D 5e shields)
	if char.Shield != "" {
		shield, err := api.GetArmor(char.Shield)
		if err == nil && shield != nil {
			// If shield AC is in API, use it, else default to +2
			if shield.ArmorClass.Base > 2 {
				baseAC += shield.ArmorClass.Base
			} else {
				baseAC += 2
			}
		} else {
			baseAC += 2
		}
	}
	return baseAC
}

// CalculateInitiative returns the initiative bonus for a character.
func CalculateInitiative(char *characterModel.Character, service *characterModel.CharacterService) int {
	return service.AbilityModifier(char.Dex)
}

// CalculatePassivePerception returns the passive perception for a character.
func CalculatePassivePerception(char *characterModel.Character, service *characterModel.CharacterService) int {
	base := 10 + service.AbilityModifier(char.Wis)
	// If proficient in Perception, add proficiency bonus
	for _, skill := range char.SkillProficiencies {
		if skill == "Perception" || skill == "perception" {
			base += char.Proficiency
			break
		}
	}
	return base
}

// min returns the smaller of two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
