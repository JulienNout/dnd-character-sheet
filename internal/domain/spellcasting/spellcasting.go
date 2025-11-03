package spellcasting

// CasterType represents the spellcasting progression type for a class.
type CasterType string

const (
	CasterNone  CasterType = "none"
	CasterFull  CasterType = "full"
	CasterHalf  CasterType = "half"
	CasterPact  CasterType = "pact"
	CasterKnown CasterType = "known"
)

// Spellcasting represents a character's spellcasting capabilities.
type Spellcasting struct {
	CasterType     CasterType
	KnownSpells    []string
	PreparedSpells []string
	SpellSlots     map[int]int // Spell level -> number of slots
}

// CanCast returns true if this caster type can cast spells.
func (ct CasterType) CanCast() bool {
	return ct != CasterNone
}

// NewSpellcasting creates a new Spellcasting instance for a character.
func NewSpellcasting(class string, level int) *Spellcasting {
	casterType := GetCasterType(class)
	slots := GetSpellSlots(casterType, level)

	return &Spellcasting{
		CasterType:     casterType,
		KnownSpells:    []string{},
		PreparedSpells: []string{},
		SpellSlots:     slots,
	}
}

// LearnSpell adds a spell to the character's known spells if valid.
// Returns true if the spell was added, false if it was already known.
func (s *Spellcasting) LearnSpell(spellName string) bool {
	// Check if spell is already known
	for _, known := range s.KnownSpells {
		if known == spellName {
			return false
		}
	}

	s.KnownSpells = append(s.KnownSpells, spellName)
	return true
}

// PrepareSpell adds a spell to the character's prepared spells if valid.
// Returns an error if the spell is not known or already prepared.
func (s *Spellcasting) PrepareSpell(spellName string) error {
	// Check if spell is known
	known := false
	for _, k := range s.KnownSpells {
		if k == spellName {
			known = true
			break
		}
	}
	if !known {
		return ErrSpellNotKnown
	}

	// Check if spell is already prepared
	for _, p := range s.PreparedSpells {
		if p == spellName {
			return ErrSpellAlreadyPrepared
		}
	}

	s.PreparedSpells = append(s.PreparedSpells, spellName)
	return nil
}

// Domain errors for spellcasting
var (
	ErrSpellNotKnown        = &SpellcastingError{"spell is not known"}
	ErrSpellAlreadyPrepared = &SpellcastingError{"spell is already prepared"}
)

// SpellcastingError represents a domain error in spellcasting operations.
type SpellcastingError struct {
	message string
}

func (e *SpellcastingError) Error() string {
	return e.message
}
