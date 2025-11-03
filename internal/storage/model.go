package storage

// CharacterSummary contains basic info about a character for listing
type CharacterSummary struct {
	Name  string `json:"name"`
	Race  string `json:"race"`
	Class string `json:"class"`
	Level int    `json:"level"`
}

// Character is the persisted representation written to/read from JSON.
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
	Spellcasting       interface{} `json:"spellcasting"`
	StrMod             int         `json:"str_mod"`
	DexMod             int         `json:"dex_mod"`
	ConMod             int         `json:"con_mod"`
	IntMod             int         `json:"int_mod"`
	WisMod             int         `json:"wis_mod"`
	ChaMod             int         `json:"cha_mod"`
	ArmorClass         int         `json:"armor_class"`
	Initiative         int         `json:"initiative"`
	PassivePerception  int         `json:"passive_perception"`
	SpellAttackBonus   int         `json:"spell_attack_bonus,omitempty"`
}

// CharacterStorage defines the interface for character persistence operations
type CharacterStorage interface {
	// Save stores a character to persistent storage
	Save(character Character) error

	// Load retrieves a character by name from persistent storage
	Load(name string) (Character, error)

	// List returns a summary of all stored characters
	List() ([]CharacterSummary, error)

	// Delete removes a character from persistent storage
	Delete(name string) error
}
