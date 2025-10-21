package spellcasting

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

// LoadSpells loads spells from a CSV file
func LoadSpells(filename string) ([]Spell, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var spells []Spell
	for i, rec := range records {
		if i == 0 {
			continue // skip header
		}
		level, _ := strconv.Atoi(rec[1])
		spells = append(spells, Spell{
			Name:  rec[0],
			Level: level,
			Class: rec[2],
		})
	}
	return spells, nil
}

// FilterSpellsByClass returns spells for a given class
func FilterSpellsByClass(spells []Spell, class string) []Spell {
	var filtered []Spell
	for _, s := range spells {
		if strings.EqualFold(s.Class, class) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
