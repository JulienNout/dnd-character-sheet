package spellcasting

import (
	"strings"
)

// AssignSpellcasting initializes the CharacterSpellcasting struct for a character based on class and level
func AssignSpellcasting(class string, level int) CharacterSpellcasting {
	casterType, ok := CasterTypeByClass[strings.ToLower(class)]
	if !ok {
		casterType = CasterNone
	}
	var slots map[int]int
	switch casterType {
	case CasterFull:
		slots = FullCasterSlots[level]
	case CasterHalf:
		slots = HalfCasterSlots[level]
	case CasterPact:
		slots = PactCasterSlots[level]
	default:
		slots = map[int]int{}
	}
	return CharacterSpellcasting{
		CasterType:     casterType,
		KnownSpells:    []string{},
		PreparedSpells: []string{},
		SpellSlots:     slots,
	}
}
