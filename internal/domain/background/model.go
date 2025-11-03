package backgroundModel

// Background represents character background reference data.
// This is a pure domain value object with no infrastructure dependencies.
type Background struct {
	Name               string
	SkillProficiencies []string
}
