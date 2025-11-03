package character

import "strings"

// Character is the core domain entity representing a player character.
// Keep this package free of application or infrastructure concerns.
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
	// Data for frontend display
	StrMod            int `json:"str_mod"`
	DexMod            int `json:"dex_mod"`
	ConMod            int `json:"con_mod"`
	IntMod            int `json:"int_mod"`
	WisMod            int `json:"wis_mod"`
	ChaMod            int `json:"cha_mod"`
	ArmorClass        int `json:"armor_class"`
	Initiative        int `json:"initiative"`
	PassivePerception int `json:"passive_perception"`
	SpellAttackBonus  int `json:"spell_attack_bonus,omitempty"`
}

// ComputeModifiers computes the ability modifiers from ability scores.
// Note: This uses integer division like many simple implementations: (score-10)/2.
func (c *Character) ComputeModifiers() {
	c.StrMod = abilityMod(c.Str)
	c.DexMod = abilityMod(c.Dex)
	c.ConMod = abilityMod(c.Con)
	c.IntMod = abilityMod(c.Int)
	c.WisMod = abilityMod(c.Wis)
	c.ChaMod = abilityMod(c.Cha)
}

// ComputeDerived computes derived values used by the frontend.
func (c *Character) ComputeDerived() {
	// Minimal defaults: base AC 10 + dex modifier. More complex armor rules belong in domain/equipment.
	c.ArmorClass = 10 + c.DexMod
	c.Initiative = c.DexMod
	c.PassivePerception = 10 + c.WisMod
}

// GetProficiencyBonus computes proficiency bonus from character level.
func (c *Character) GetProficiencyBonus() int {
	switch {
	case c.Level >= 1 && c.Level <= 4:
		return 2
	case c.Level >= 5 && c.Level <= 8:
		return 3
	case c.Level >= 9 && c.Level <= 12:
		return 4
	case c.Level >= 13 && c.Level <= 16:
		return 5
	case c.Level >= 17:
		return 6
	default:
		return 0
	}
}

// ApplyRacialBonuses applies ability score increases based on race.
func (c *Character) ApplyRacialBonuses() {
	race := strings.ToLower(c.Race)
	switch race {
	case "dwarf":
		c.Con += 2
	case "hill dwarf":
		c.Con += 2
		c.Wis += 1
	case "elf":
		c.Dex += 2
	case "high elf":
		c.Dex += 2
		c.Int += 1
	case "halfling":
		c.Dex += 2
	case "lightfoot halfling", "lightfoot":
		c.Dex += 2
		c.Cha += 1
	case "human":
		c.Str += 1
		c.Dex += 1
		c.Con += 1
		c.Int += 1
		c.Wis += 1
		c.Cha += 1
	case "dragonborn":
		c.Str += 2
		c.Cha += 1
	case "gnome":
		c.Int += 2
	case "rock gnome":
		c.Int += 2
		c.Con += 1
	case "half-elf":
		c.Cha += 2
		c.Dex += 1
		c.Con += 1
	case "half orc", "half-orc":
		c.Str += 2
		c.Con += 1
	case "tiefling":
		c.Int += 1
		c.Cha += 2
	}
}

// abilityMod computes floor((score-10)/2) with correct rounding for negatives.
func abilityMod(score int) int {
	base := (score - 10) / 2
	if (score-10)%2 < 0 {
		base--
	}
	return base
}
