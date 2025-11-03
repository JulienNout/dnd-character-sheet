package storage

import (
	"encoding/json"
	"fmt"

	stor "modules/dndcharactersheet/internal/adapters/storage/jsonstorage"
	characterpkg "modules/dndcharactersheet/internal/domain/character"
	"modules/dndcharactersheet/internal/domain/spellcasting"
	"modules/dndcharactersheet/internal/ports"
)

// JSONRepository implements ports.CharacterRepository by delegating to the existing
// SingleFileStorage implementation under internal/storage. It performs mapping
// between domain types and the storage model.
type JSONRepository struct {
	backend stor.CharacterStorage
}

// NewJSONRepository creates a repository backed by the provided storage backend.
// filename is passed to the underlying SingleFileStorage. Use empty string to
// use the default file if needed.
func NewJSONRepository(filename string) *JSONRepository {
	return &JSONRepository{backend: stor.NewSingleFileStorage(filename)}
}

func (r *JSONRepository) Save(c *characterpkg.Character) error {
	if c == nil {
		return fmt.Errorf("nil character")
	}
	sm := domainToStorage(c)
	return r.backend.Save(sm)
}

func (r *JSONRepository) GetAll() ([]*characterpkg.Character, error) {
	summaries, err := r.backend.List()
	if err != nil {
		return nil, err
	}

	out := make([]*characterpkg.Character, 0, len(summaries))
	for _, s := range summaries {
		stored, err := r.backend.Load(s.Name)
		if err != nil {
			// skip individual failing loads but log via error aggregation
			return nil, fmt.Errorf("failed loading character %s: %w", s.Name, err)
		}
		d := storageToDomain(&stored)
		out = append(out, d)
	}
	return out, nil
}

func (r *JSONRepository) GetByID(id string) (*characterpkg.Character, error) {
	stored, err := r.backend.Load(id)
	if err != nil {
		return nil, err
	}
	d := storageToDomain(&stored)
	return d, nil
}

func (r *JSONRepository) Delete(id string) error {
	return r.backend.Delete(id)
}

// Ensure JSONRepository satisfies the interface at compile time.
var _ ports.CharacterRepository = (*JSONRepository)(nil)

// domainToStorage converts a domain character to the storage model.
func domainToStorage(d *characterpkg.Character) stor.Character {
	return stor.Character{
		Name:               d.Name,
		Race:               d.Race,
		Class:              d.Class,
		Level:              d.Level,
		Str:                d.Str,
		Dex:                d.Dex,
		Con:                d.Con,
		Int:                d.Int,
		Wis:                d.Wis,
		Cha:                d.Cha,
		Background:         d.Background,
		Proficiency:        d.Proficiency,
		SkillProficiencies: d.SkillProficiencies,
		MainHand:           d.MainHand,
		OffHand:            d.OffHand,
		Armor:              d.Armor,
		Shield:             d.Shield,
		Spellcasting:       marshalSpellcasting(d.Spellcasting),
		StrMod:             d.StrMod,
		DexMod:             d.DexMod,
		ConMod:             d.ConMod,
		IntMod:             d.IntMod,
		WisMod:             d.WisMod,
		ChaMod:             d.ChaMod,
		ArmorClass:         d.ArmorClass,
		Initiative:         d.Initiative,
		PassivePerception:  d.PassivePerception,
		SpellAttackBonus:   d.SpellAttackBonus,
	}
}

// storageToDomain converts a storage model character to the domain type.
func storageToDomain(s *stor.Character) *characterpkg.Character {
	d := &characterpkg.Character{
		Name:               s.Name,
		Race:               s.Race,
		Class:              s.Class,
		Level:              s.Level,
		Str:                s.Str,
		Dex:                s.Dex,
		Con:                s.Con,
		Int:                s.Int,
		Wis:                s.Wis,
		Cha:                s.Cha,
		Background:         s.Background,
		Proficiency:        s.Proficiency,
		SkillProficiencies: s.SkillProficiencies,
		MainHand:           s.MainHand,
		OffHand:            s.OffHand,
		Armor:              s.Armor,
		Shield:             s.Shield,
		Spellcasting:       unmarshalSpellcasting(s.Spellcasting),
		StrMod:             s.StrMod,
		DexMod:             s.DexMod,
		ConMod:             s.ConMod,
		IntMod:             s.IntMod,
		WisMod:             s.WisMod,
		ChaMod:             s.ChaMod,
		ArmorClass:         s.ArmorClass,
		Initiative:         s.Initiative,
		PassivePerception:  s.PassivePerception,
		SpellAttackBonus:   s.SpellAttackBonus,
	}
	return d
}

// marshalSpellcasting converts domain spellcasting to storage format (with JSON-serializable structure).
func marshalSpellcasting(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	// Type assert to domain spellcasting
	sc, ok := data.(*spellcasting.Spellcasting)
	if !ok {
		return data // Return as-is if not our type
	}

	// Create a storage model with JSON tags for serialization
	storageModel := struct {
		CasterType     string      `json:"CasterType"`
		KnownSpells    []string    `json:"KnownSpells"`
		PreparedSpells []string    `json:"PreparedSpells"`
		SpellSlots     map[int]int `json:"SpellSlots"`
	}{
		CasterType:     string(sc.CasterType),
		KnownSpells:    sc.KnownSpells,
		PreparedSpells: sc.PreparedSpells,
		SpellSlots:     sc.SpellSlots,
	}

	return storageModel
}

// unmarshalSpellcasting converts the storage interface{} back to domain spellcasting type.
func unmarshalSpellcasting(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	// When loading from JSON, interface{} will be a map[string]interface{}
	// Use a storage model with JSON tags for deserialization
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var storageModel struct {
		CasterType     string      `json:"CasterType"`
		KnownSpells    []string    `json:"KnownSpells"`
		PreparedSpells []string    `json:"PreparedSpells"`
		SpellSlots     map[int]int `json:"SpellSlots"`
	}

	if err := json.Unmarshal(jsonBytes, &storageModel); err != nil {
		return nil
	}

	// Map to domain model (no JSON tags)
	sc := &spellcasting.Spellcasting{
		CasterType:     spellcasting.CasterType(storageModel.CasterType),
		KnownSpells:    storageModel.KnownSpells,
		PreparedSpells: storageModel.PreparedSpells,
		SpellSlots:     storageModel.SpellSlots,
	}

	return sc
}
