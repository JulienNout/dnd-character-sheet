package characterModel

type Character struct {
	Name               string   `json:"name"`
	Race               string   `json:"race"`
	Class              string   `json:"class"`
	Level              int      `json:"level"`
	Str                int      `json:"str"`
	Dex                int      `json:"dex"`
	Con                int      `json:"con"`
	Int                int      `json:"int"`
	Wis                int      `json:"wis"`
	Cha                int      `json:"cha"`
	Background         string   `json:"background"`
	Proficiency        int      `json:"proficiency"`
	SkillProficiencies []string `json:"skill_proficiencies"`
}

// Returns the proficiency bonus based on character level
func GetProficiencyBonus(level int) int {
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
