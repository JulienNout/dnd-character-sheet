package spellcastingadapter

import (
	"encoding/csv"
	"modules/dndcharactersheet/internal/ports"
	"os"
	"strconv"
	"strings"
)

// CSVSpellRepository implements the SpellRepository port using CSV files.
type CSVSpellRepository struct {
	filename string
}

// NewCSVSpellRepository creates a new CSV-based spell repository.
func NewCSVSpellRepository(filename string) *CSVSpellRepository {
	return &CSVSpellRepository{
		filename: filename,
	}
}

// LoadSpells loads spells from the CSV file.
func (r *CSVSpellRepository) LoadSpells() ([]ports.Spell, error) {
	file, err := os.Open(r.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var spells []ports.Spell
	for i, rec := range records {
		if i == 0 {
			continue // skip header
		}

		// CSV format: Name, Level, Class/Classes, ...
		level, _ := strconv.Atoi(rec[1])
		classes := parseClasses(rec[2])

		spells = append(spells, ports.Spell{
			Index:   toIndex(rec[0]),
			Name:    rec[0],
			Level:   level,
			Classes: classes,
		})
	}
	return spells, nil
}

// FilterByClass filters spells available to a specific class.
func (r *CSVSpellRepository) FilterByClass(spells []ports.Spell, class string) []ports.Spell {
	var filtered []ports.Spell
	class = strings.ToLower(class)

	for _, s := range spells {
		for _, c := range s.Classes {
			if strings.EqualFold(c, class) {
				filtered = append(filtered, s)
				break
			}
		}
	}
	return filtered
}

// parseClasses splits a comma-separated list of classes.
func parseClasses(classStr string) []string {
	if classStr == "" {
		return []string{}
	}

	parts := strings.Split(classStr, ",")
	classes := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			classes = append(classes, trimmed)
		}
	}
	return classes
}

// toIndex converts a spell name to an index format.
func toIndex(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
