package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"modules/dndcharactersheet/internal/ports"
)

func TestAPIAdapter_GetWeaponArmorSpell(t *testing.T) {
	// Setup a test HTTP server that returns canned JSON for equipment and spells
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/equipment/longsword":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"Longsword","weapon_category":"Martial","range":{"normal":5},"two_handed":false}`))
			return
		case "/equipment/chain-mail":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"Chain Mail","armor_class":{"base":16,"dex_bonus":false}}`))
			return
		case "/spells/acid-arrow":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"Acid Arrow","range":"90 feet","school":{"name":"Evocation"}}`))
			return
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer ts.Close()

	adapter := NewAPIAdapter(ts.URL)

	// Weapon
	w, err := adapter.GetWeapon("longsword")
	if err != nil {
		t.Fatalf("GetWeapon error: %v", err)
	}
	if w.Name != "Longsword" || w.Category == "" || w.Range != 5 {
		t.Fatalf("unexpected weapon: %+v", w)
	}

	// Armor
	a, err := adapter.GetArmor("chain-mail")
	if err != nil {
		t.Fatalf("GetArmor error: %v", err)
	}
	if a.Name != "Chain Mail" || a.BaseAC != 16 {
		t.Fatalf("unexpected armor: %+v", a)
	}

	// Spell
	s, err := adapter.GetSpell("acid-arrow")
	if err != nil {
		t.Fatalf("GetSpell error: %v", err)
	}
	if s.Name != "Acid Arrow" || s.School != "Evocation" || s.Range != "90 feet" {
		t.Fatalf("unexpected spell: %+v", s)
	}

	// Also verify returned types implement the ports DTO shape (sanity)
	var _ *ports.WeaponInfo = w
	var _ *ports.ArmorInfo = a
	var _ *ports.SpellInfo = s
}
