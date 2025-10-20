package classModel

import (
	"encoding/json"
	"os"
)

type Class struct {
	Name               string   `json:"name"`
	SkillProficiencies []string `json:"skill_proficiencies"`
	SkillCount         int      `json:"skill_count"` // How many skills they can choose
}

func LoadClasses(filename string) ([]Class, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var classes []Class
	err = json.Unmarshal(data, &classes)
	return classes, err
}
