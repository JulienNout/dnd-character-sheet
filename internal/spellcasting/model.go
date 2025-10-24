// GetDefaultSpellSlots returns default spell slots for a class and level
package spellcasting

import "strings"

type CasterType string

const (
	CasterNone  CasterType = "none"
	CasterFull  CasterType = "full"
	CasterHalf  CasterType = "half"
	CasterPact  CasterType = "pact"
	CasterKnown CasterType = "known"
)

type Spell struct {
	Name  string
	Level int
	Class string
}

// CharacterSpellcasting holds spellcasting data for a character
// (to be embedded or referenced in your character model)
type CharacterSpellcasting struct {
	CasterType     CasterType
	KnownSpells    []string
	PreparedSpells []string
	SpellSlots     map[int]int // level -> slots
}

// CasterTypeByClass maps class names to their caster type
var CasterTypeByClass = map[string]CasterType{
	"wizard":   CasterFull,
	"cleric":   CasterFull,
	"druid":    CasterFull,
	"bard":     CasterKnown,
	"sorcerer": CasterKnown,
	"paladin":  CasterHalf,
	"ranger":   CasterHalf,
	"warlock":  CasterPact,
}

// CantripsKnownByClassAndLevel maps class names to a slice where index is level and value is cantrips known
var CantripsKnownByClassAndLevel = map[string][]int{
	"bard":     {0, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
	"cleric":   {0, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
	"druid":    {0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
	"sorcerer": {0, 4, 4, 4, 4, 4, 4, 4, 4, 4, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5},
	"warlock":  {0, 2, 2, 2, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
	"wizard":   {0, 3, 3, 3, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
}

// GetCantripsKnown returns the number of cantrips known for a class and level
func GetCantripsKnown(class string, level int) int {
	class = strings.ToLower(class)
	if arr, ok := CantripsKnownByClassAndLevel[class]; ok {
		if level >= 1 && level < len(arr) {
			return arr[level]
		}
		if level >= len(arr) {
			return arr[len(arr)-1]
		}
	}
	return 0
}

// FullCasterSlots[level][slotLevel] = slots
var FullCasterSlots = map[int]map[int]int{
	1:  {1: 2},
	2:  {1: 3},
	3:  {1: 4, 2: 2},
	4:  {1: 4, 2: 3},
	5:  {1: 4, 2: 3, 3: 2},
	6:  {1: 4, 2: 3, 3: 3},
	7:  {1: 4, 2: 3, 3: 3, 4: 1},
	8:  {1: 4, 2: 3, 3: 3, 4: 2},
	9:  {1: 4, 2: 3, 3: 3, 4: 3, 5: 1},
	10: {1: 4, 2: 3, 3: 3, 4: 3, 5: 2},
	11: {1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1},
	12: {1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1},
	13: {1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1, 7: 1},
	14: {1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1, 7: 1},
	15: {1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1, 7: 1, 8: 1},
	16: {1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1, 7: 1, 8: 1},
	17: {1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1, 7: 1, 8: 1, 9: 1},
	18: {1: 4, 2: 3, 3: 3, 4: 3, 5: 3, 6: 1, 7: 1, 8: 1, 9: 1},
	19: {1: 4, 2: 3, 3: 3, 4: 3, 5: 3, 6: 2, 7: 1, 8: 1, 9: 1},
	20: {1: 4, 2: 3, 3: 3, 4: 3, 5: 3, 6: 2, 7: 2, 8: 1, 9: 1},
}

// HalfCasterSlots[level][slotLevel] = slots (Paladin/Ranger)
var HalfCasterSlots = map[int]map[int]int{
	1:  {},
	2:  {},
	3:  {1: 2},
	4:  {1: 3},
	5:  {1: 4, 2: 2},
	6:  {1: 4, 2: 2},
	7:  {1: 4, 2: 3},
	8:  {1: 4, 2: 3},
	9:  {1: 4, 2: 3, 3: 2},
	10: {1: 4, 2: 3, 3: 2},
	11: {1: 4, 2: 3, 3: 3},
	12: {1: 4, 2: 3, 3: 3},
	13: {1: 4, 2: 3, 3: 3, 4: 1},
	14: {1: 4, 2: 3, 3: 3, 4: 1},
	15: {1: 4, 2: 3, 3: 3, 4: 2},
	16: {1: 4, 2: 3, 3: 3, 4: 2},
	17: {1: 4, 2: 3, 3: 3, 4: 3, 5: 1},
	18: {1: 4, 2: 3, 3: 3, 4: 3, 5: 1},
	19: {1: 4, 2: 3, 3: 3, 4: 3, 5: 2},
	20: {1: 4, 2: 3, 3: 3, 4: 3, 5: 2},
}

// PactCasterSlots[level][slotLevel] = slots (Warlock)
var PactCasterSlots = map[int]map[int]int{
	1:  {1: 1},
	2:  {1: 2},
	3:  {1: 2},
	4:  {2: 2},
	5:  {3: 2},
	6:  {3: 2},
	7:  {4: 2},
	8:  {4: 2},
	9:  {5: 2},
	10: {5: 2},
	11: {5: 3},
	12: {5: 3},
	13: {5: 3},
	14: {5: 3},
	15: {5: 3},
	16: {5: 3},
	17: {5: 4},
	18: {5: 4},
	19: {5: 4},
	20: {5: 4},
}
