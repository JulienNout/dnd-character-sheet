package ports

// Spell represents a D&D spell with its properties.
type Spell struct {
	Index    string
	Name     string
	Level    int
	School   string
	Classes  []string
	CastTime string
	Range    string
	Duration string
}

// SpellRepository provides access to spell data.
type SpellRepository interface {
	// LoadSpells loads all spells from the data source.
	LoadSpells() ([]Spell, error)

	// FilterByClass filters spells available to a specific class.
	FilterByClass(spells []Spell, class string) []Spell
}
