package ports

// SpellcastingEngine abstracts spellcasting operations so the application
// layer can orchestrate without depending on legacy implementation details.
// Note: The concrete spellcasting data type is kept as interface{} to avoid
// coupling domain to a specific representation.
type SpellcastingEngine interface {
	AssignSpellcasting(class string, level int) (interface{}, error)
	LearnSpell(sc interface{}, class string, spellName string) (updated interface{}, message string, err error)
	PrepareSpell(sc interface{}, class string, spellName string) (updated interface{}, message string, err error)
	FormatSpellSlots(sc interface{}, class string, level int) string
	FormatCantrips(sc interface{}) string
}
