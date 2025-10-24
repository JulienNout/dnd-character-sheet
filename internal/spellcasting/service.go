package spellcasting

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// FormatCantrips returns a formatted string for cantrips (level 0 spells) from SpellSlots or KnownSpells
func FormatCantrips(cs *CharacterSpellcasting) string {
	if cs == nil {
		return ""
	}
	// If cantrips are tracked in SpellSlots (level 0)
	if slots, exists := cs.SpellSlots[0]; exists {
		return fmt.Sprintf("Spell slots:\n  Level 0: %d\n", slots)
	}
	// If cantrips are tracked in KnownSpells
	if len(cs.KnownSpells) > 0 {
		var cantrips []string
		for _, spellName := range cs.KnownSpells {
			if strings.Contains(strings.ToLower(spellName), "cantrip") || strings.Contains(strings.ToLower(spellName), "level 0") {
				cantrips = append(cantrips, spellName)
			}
		}
		if len(cantrips) > 0 {
			return fmt.Sprintf("Cantrips: %s\n", strings.Join(cantrips, ", "))
		}
	}
	return ""
}

// Returns a formatted string for a character's spell slots
func FormatSpellSlots(cs *CharacterSpellcasting, class string, level int) string {
	if cs == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("Spell slots:\n")
	// Print cantrips as Level 0 using GetCantripsKnown if applicable
	cantrips := GetCantripsKnown(class, level)
	if cantrips > 0 {
		sb.WriteString(fmt.Sprintf("  Level 0: %d\n", cantrips))
	}
	// Print actual spell slots (levels 1-9)
	if cs.SpellSlots != nil {
		for lvl := 1; lvl <= 9; lvl++ {
			if slots, exists := cs.SpellSlots[lvl]; exists {
				sb.WriteString(fmt.Sprintf("  Level %d: %d\n", lvl, slots))
			}
		}
	}
	return sb.String()
}

// Default spell slots by level
func GetDefaultSpellSlots(class string, level int) map[int]int {
	class = strings.ToLower(class)
	switch class {
	case "wizard", "cleric", "druid":
		return FullCasterSlots[level]
	case "paladin", "ranger":
		return HalfCasterSlots[level]
	case "warlock":
		return PactCasterSlots[level]
	case "bard", "sorcerer":
		return FullCasterSlots[level]
	default:
		return map[int]int{}
	}
}

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
