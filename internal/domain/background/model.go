package backgroundModel

import (
	"encoding/json"
	"os"
)

type Background struct {
	Name               string   `json:"name"`
	SkillProficiencies []string `json:"skill_proficiencies"`
}

func LoadBackgrounds(filename string) ([]Background, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var backgrounds []Background
	err = json.Unmarshal(data, &backgrounds)
	return backgrounds, err
}
