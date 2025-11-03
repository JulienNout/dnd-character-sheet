package classModel

// Class represents character class reference data.
// This is a pure domain value object with no infrastructure dependencies.
type Class struct {
	Name               string
	SkillProficiencies []string
	SkillCount         int // How many skills they can choose
}
