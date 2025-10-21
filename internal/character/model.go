package characterModel

type Character struct {
	Name               string      `json:"name"`
	Race               string      `json:"race"`
	Class              string      `json:"class"`
	Level              int         `json:"level"`
	Str                int         `json:"str"`
	Dex                int         `json:"dex"`
	Con                int         `json:"con"`
	Int                int         `json:"int"`
	Wis                int         `json:"wis"`
	Cha                int         `json:"cha"`
	Background         string      `json:"background"`
	Proficiency        int         `json:"proficiency"`
	SkillProficiencies []string    `json:"skill_proficiencies"`
	MainHand           string      `json:"main_hand,omitempty"`
	OffHand            string      `json:"off_hand,omitempty"`
	Armor              string      `json:"armor,omitempty"`
	Shield             string      `json:"shield,omitempty"`
	Spellcasting       interface{} `json:"spellcasting"` // Spellcasting data handled in service logic
}
