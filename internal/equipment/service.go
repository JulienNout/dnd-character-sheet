package equipment

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func LoadEquipmentFromCSV(filename string) ([]EquipmentItem, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open equipment csv: %w", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read equipment csv: %w", err)
	}

	if len(records) <= 1 {
		return nil, nil // no data
	}

	header := records[0]
	// map header names to indexes
	idx := make(map[string]int)
	for i, h := range header {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}

	var items []EquipmentItem
	for _, rec := range records[1:] {
		// helper to get column by name if present
		get := func(name string) string {
			if i, ok := idx[strings.ToLower(name)]; ok && i < len(rec) {
				return strings.TrimSpace(rec[i])
			}
			return ""
		}

		item := EquipmentItem{
			Name:     get("name"),
			Category: get("type"),
		}
		// fallback: if "type" column not present, try "category"
		if item.Category == "" {
			item.Category = get("category")
		}

		items = append(items, item)
	}

	return items, nil
}

// FindEquipmentByName finds an equipment item by name (case-insensitive)
func FindEquipmentByName(items []EquipmentItem, name string) *EquipmentItem {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, it := range items {
		if strings.ToLower(it.Name) == name {
			return &it
		}
	}
	return nil
}
