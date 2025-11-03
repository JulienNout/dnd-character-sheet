package spellcastingadapter

import (
	"fmt"
	"strings"

	"modules/dndcharactersheet/internal/domain/spellcasting"
	"modules/dndcharactersheet/internal/ports"
)

// EngineAdapter implements the SpellcastingEngine port using domain spellcasting logic.
type EngineAdapter struct {
	spellRepo ports.SpellRepository
}

// NewEngineAdapter creates a new spellcasting engine adapter.
func NewEngineAdapter(spellRepo ports.SpellRepository) *EngineAdapter {
	return &EngineAdapter{
		spellRepo: spellRepo,
	}
}

func (e *EngineAdapter) AssignSpellcasting(class string, level int) (interface{}, error) {
	sc := spellcasting.NewSpellcasting(class, level)
	return sc, nil
}

func (e *EngineAdapter) LearnSpell(sc interface{}, class string, spellName string) (interface{}, string, error) {
	scTyped, ok := sc.(*spellcasting.Spellcasting)
	if !ok {
		return sc, "invalid spellcasting data", nil
	}

	// Check caster type
	if scTyped.CasterType != spellcasting.CasterKnown && scTyped.CasterType != spellcasting.CasterPact {
		return sc, "this class prepares spells and can't learn them", nil
	}

	// Load spells to validate the spell exists and is available to the class
	spells, err := e.spellRepo.LoadSpells()
	if err != nil {
		return sc, "", err
	}

	classSpells := e.spellRepo.FilterByClass(spells, class)

	lname := strings.ToLower(spellName)
	found := false
	for _, s := range classSpells {
		if strings.EqualFold(s.Name, lname) {
			found = true
			break
		}
	}

	if !found {
		return sc, fmt.Sprintf("spell '%s' not found or not available to %s", spellName, class), nil
	}

	// Use domain method to learn the spell
	if !scTyped.LearnSpell(spellName) {
		return sc, "already learned this spell", nil
	}

	return scTyped, fmt.Sprintf("learned spell %s", strings.ToLower(spellName)), nil
}

func (e *EngineAdapter) PrepareSpell(sc interface{}, class string, spellName string) (interface{}, string, error) {
	scTyped, ok := sc.(*spellcasting.Spellcasting)
	if !ok {
		return sc, "invalid spellcasting data", nil
	}

	// Check caster type
	if scTyped.CasterType != spellcasting.CasterFull && scTyped.CasterType != spellcasting.CasterHalf {
		return sc, "this class learns spells and can't prepare them", nil
	}

	// Load spells to get level information
	spells, err := e.spellRepo.LoadSpells()
	if err != nil {
		return sc, "", err
	}

	classSpells := e.spellRepo.FilterByClass(spells, class)

	lname := strings.ToLower(spellName)
	var foundSpell *ports.Spell
	for _, s := range classSpells {
		if strings.EqualFold(s.Name, lname) {
			foundSpell = &s
			break
		}
	}

	if foundSpell == nil {
		return sc, fmt.Sprintf("spell '%s' not found or not available to %s", spellName, class), nil
	}

	// Check spell level against available slots (unless it's a cantrip)
	if foundSpell.Level > 0 {
		maxSlot := 0
		for lvl := range scTyped.SpellSlots {
			if lvl > maxSlot {
				maxSlot = lvl
			}
		}
		if foundSpell.Level > maxSlot {
			return sc, "the spell has higher level than the available spell slots", nil
		}
	}

	// Use domain method to prepare the spell
	if err := scTyped.PrepareSpell(spellName); err != nil {
		return sc, err.Error(), nil
	}

	return scTyped, fmt.Sprintf("prepared spell %s", strings.ToLower(spellName)), nil
}

func (e *EngineAdapter) FormatSpellSlots(sc interface{}, class string, level int) string {
	scTyped, ok := sc.(*spellcasting.Spellcasting)
	if !ok {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("Spell slots:\n")

	// Print cantrips as Level 0 using domain function
	cantrips := spellcasting.GetCantripsKnown(class, level)
	if cantrips > 0 {
		sb.WriteString(fmt.Sprintf("  Level 0: %d\n", cantrips))
	}

	// Print actual spell slots (levels 1-9)
	if scTyped.SpellSlots != nil {
		for lvl := 1; lvl <= 9; lvl++ {
			if slots, exists := scTyped.SpellSlots[lvl]; exists {
				sb.WriteString(fmt.Sprintf("  Level %d: %d\n", lvl, slots))
			}
		}
	}
	return sb.String()
}

func (e *EngineAdapter) FormatCantrips(sc interface{}) string {
	scTyped, ok := sc.(*spellcasting.Spellcasting)
	if !ok {
		return ""
	}

	// If cantrips are tracked in SpellSlots (level 0)
	if slots, exists := scTyped.SpellSlots[0]; exists {
		return fmt.Sprintf("Spell slots:\n  Level 0: %d\n", slots)
	}

	// If cantrips are tracked in KnownSpells
	if len(scTyped.KnownSpells) > 0 {
		var cantrips []string
		for _, spellName := range scTyped.KnownSpells {
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

var _ ports.SpellcastingEngine = (*EngineAdapter)(nil)
