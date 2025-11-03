package race

import "strings"

// GetRacialSkillProficiencies returns the list of skill proficiencies granted by race.
// We simplify SRD features (e.g., Stonecunning, Keen Senses, Menacing) to regular proficiencies.
// - Dwarf: Stonecunning -> history
// - Elf: Keen Senses -> perception
// - Half-Orc: Menacing -> intimidation
// Other races return an empty list by default.
func GetRacialSkillProficiencies(race string) []string {
	r := strings.ToLower(strings.TrimSpace(race))
	r = strings.ReplaceAll(r, "-", " ")

	switch r {
	case "dwarf", "hill dwarf", "mountain dwarf":
		return []string{"history"}
	case "elf", "high elf", "wood elf", "dark elf", "drow":
		return []string{"perception"}
	case "half orc", "half-orc":
		return []string{"intimidation"}
	default:
		return nil
	}
}
