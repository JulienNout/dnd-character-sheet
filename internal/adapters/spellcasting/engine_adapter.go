package spellcastingadapter

import (
	"strings"

	"modules/dndcharactersheet/internal/ports"
	legacy "modules/dndcharactersheet/internal/spellcasting"
)

// EngineAdapter wraps the legacy spellcasting helpers behind the SpellcastingEngine port.
type EngineAdapter struct{}

func NewEngineAdapter() *EngineAdapter { return &EngineAdapter{} }

func (e *EngineAdapter) AssignSpellcasting(class string, level int) (interface{}, error) {
	sc := legacy.AssignSpellcasting(class, level)
	return sc, nil
}

func (e *EngineAdapter) LearnSpell(sc interface{}, class string, spellName string) (interface{}, string, error) {
	spells, err := legacy.LoadSpells("5e-SRD-Spells.csv")
	if err != nil {
		return sc, "", err
	}
	var found *legacy.Spell
	lname := strings.ToLower(spellName)
	lclass := strings.ToLower(class)
	for _, s := range spells {
		if strings.EqualFold(s.Name, lname) && strings.Contains(strings.ToLower(s.Class), lclass) {
			found = &s
			break
		}
	}
	// If not found, still try to learn to surface legacy validation message
	scTyped, _ := sc.(legacy.CharacterSpellcasting)
	if found == nil {
		msg := legacy.LearnSpell(&scTyped, legacy.Spell{Name: spellName})
		return scTyped, msg, nil
	}
	msg := legacy.LearnSpell(&scTyped, *found)
	return scTyped, msg, nil
}

func (e *EngineAdapter) PrepareSpell(sc interface{}, class string, spellName string) (interface{}, string, error) {
	spells, err := legacy.LoadSpells("5e-SRD-Spells.csv")
	if err != nil {
		return sc, "", err
	}
	var found *legacy.Spell
	lname := strings.ToLower(spellName)
	lclass := strings.ToLower(class)
	for _, s := range spells {
		if strings.EqualFold(s.Name, lname) && strings.Contains(strings.ToLower(s.Class), lclass) {
			found = &s
			break
		}
	}
	scTyped, _ := sc.(legacy.CharacterSpellcasting)
	if found == nil {
		msg := legacy.PrepareSpell(&scTyped, legacy.Spell{Name: spellName})
		return scTyped, msg, nil
	}
	msg := legacy.PrepareSpell(&scTyped, *found)
	return scTyped, msg, nil
}

func (e *EngineAdapter) FormatSpellSlots(sc interface{}, class string, level int) string {
	scTyped, _ := sc.(legacy.CharacterSpellcasting)
	return legacy.FormatSpellSlots(&scTyped, class, level)
}

func (e *EngineAdapter) FormatCantrips(sc interface{}) string {
	scTyped, _ := sc.(legacy.CharacterSpellcasting)
	return legacy.FormatCantrips(&scTyped)
}

var _ ports.SpellcastingEngine = (*EngineAdapter)(nil)
