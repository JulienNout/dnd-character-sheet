package api

import (
	"fmt"
	"testing"
)

func TestFetchSpellsWithWorkers(t *testing.T) {
	indexes := []string{"acid-arrow", "fireball", "mage-armor"}
	results := FetchSpellsWithWorkers(indexes, 3)
	for i, s := range results {
		if s == nil {
			t.Errorf("Spell %s not enriched", indexes[i])
		} else {
			fmt.Printf("Spell: %+v\n", s)
		}
	}
}

func TestFetchWeaponsWithWorkers(t *testing.T) {
	indexes := []string{"longsword", "shortbow"}
	results := FetchWeaponsWithWorkers(indexes, 2)
	for i, w := range results {
		if w == nil {
			t.Errorf("Weapon %s not enriched", indexes[i])
		} else {
			fmt.Printf("Weapon: %+v\n", w)
		}
	}
}

func TestFetchArmorsWithWorkers(t *testing.T) {
	indexes := []string{"chain-mail", "leather-armor"}
	results := FetchArmorsWithWorkers(indexes, 2)
	for i, a := range results {
		if a == nil {
			t.Errorf("Armor %s not enriched", indexes[i])
		} else {
			fmt.Printf("Armor: %+v\n", a)
		}
	}
}
