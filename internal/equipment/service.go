package equipment

import (
	"encoding/csv"
	"fmt"
	"modules/dndcharactersheet/internal/api"
	characterModel "modules/dndcharactersheet/internal/character"
	"os"
	"strings"
)

type EquipmentDisplay struct {
	MainHand string
	OffHand  string
	Armor    string
	Shield   string
}

// GetFormattedEquipment returns formatted equipment strings for a character, enriched via API
func GetFormattedEquipment(char *characterModel.Character) EquipmentDisplay {
	var disp EquipmentDisplay
	// Main hand
	if char.MainHand != "" {
		idx := api.ToAPIIndex(char.MainHand)
		weapon, err := api.GetWeapon(idx)
		var mainHandName string
		if err == nil && weapon != nil {
			mainHandName = strings.ToLower(weapon.Name)
		} else {
			mainHandName = strings.ToLower(char.MainHand)
		}
		disp.MainHand = mainHandName
	}
	// Off hand
	if char.OffHand != "" {
		idx := api.ToAPIIndex(char.OffHand)
		weapon, err := api.GetWeapon(idx)
		var offHandName string
		if err == nil && weapon != nil {
			offHandName = strings.ToLower(weapon.Name)
		} else {
			offHandName = strings.ToLower(char.OffHand)
		}
		disp.OffHand = offHandName
	}
	// Armor
	if char.Armor != "" {
		idx := api.ToAPIIndex(char.Armor)
		armor, err := api.GetArmor(idx)
		var armorName string
		if err == nil && armor != nil {
			armorName = strings.ToLower(armor.Name)
		} else {
			armorName = strings.ToLower(char.Armor)
		}
		disp.Armor = armorName
	}

	if char.Shield != "" {
		idx := api.ToAPIIndex(char.Shield)
		shield, err := api.GetArmor(idx)
		var shieldName string
		if err == nil && shield != nil {
			shieldName = strings.ToLower(shield.Name)
		} else {
			shieldName = strings.ToLower(char.Shield)
		}
		disp.Shield = shieldName
	}
	return disp
}

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
