package spellcasting

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

// LoadSpells loads spells from a CSV file
func LoadSpells(filename string) ([]Spell, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var spells []Spell
	for i, rec := range records {
		if i == 0 {
			continue // skip header
		}
		level, _ := strconv.Atoi(rec[1])
		spells = append(spells, Spell{
			Name:  rec[0],
			Level: level,
			Class: rec[2],
		})
	}
	return spells, nil
}

// FilterSpellsByClass returns spells for a given class
func FilterSpellsByClass(spells []Spell, class string) []Spell {
	var filtered []Spell
	for _, s := range spells {
		if strings.EqualFold(s.Class, class) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// CanCastSpells returns true if the caster type is not none
func CanCastSpells(casterType CasterType) bool {
	return casterType == CasterFull || casterType == CasterHalf || casterType == CasterPact || casterType == CasterKnown
}

// LearnSpell attempts to add a spell to the character's known spells
func LearnSpell(cs *CharacterSpellcasting, spell Spell) string {
	switch cs.CasterType {
	case CasterKnown, CasterPact:
		for _, s := range cs.KnownSpells {
			if strings.EqualFold(s, spell.Name) {
				return "Already learned this spell"
			}
		}
		cs.KnownSpells = append(cs.KnownSpells, spell.Name)
		return "Learned spell " + strings.ToLower(spell.Name)
	case CasterFull, CasterHalf:
		return "this class prepares spells and can't learn them"
	default:
		return "this class can't cast spells"
	}
}

// PrepareSpell attempts to add a spell to the character's prepared spells
func PrepareSpell(cs *CharacterSpellcasting, spell Spell) string {
	switch cs.CasterType {
	case CasterFull, CasterHalf:
		if spell.Level == 0 {
			cs.PreparedSpells = append(cs.PreparedSpells, spell.Name)
			return "Prepared spell " + strings.ToLower(spell.Name)
		}
		maxSlot := 0
		for lvl := range cs.SpellSlots {
			if lvl > maxSlot {
				maxSlot = lvl
			}
		}
		if spell.Level > maxSlot {
			return "the spell has higher level than the available spell slots"
		}
		cs.PreparedSpells = append(cs.PreparedSpells, spell.Name)
		return "Prepared spell " + strings.ToLower(spell.Name)
	case CasterKnown, CasterPact:
		return "this class learns spells and can't prepare them"
	default:
		return "this class can't cast spells"
	}
}
