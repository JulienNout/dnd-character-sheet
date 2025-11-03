package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"modules/dndcharactersheet/internal/ports"
)

// APIAdapter calls a DnD-style API and implements the enricher ports.
type APIAdapter struct {
	baseURL string
	client  *http.Client
}

// NewAPIAdapter creates an adapter pointing to baseURL (e.g., https://www.dnd5eapi.co/api/2014)
func NewAPIAdapter(baseURL string) *APIAdapter {
	return &APIAdapter{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func toAPIIndex(name string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "-"))
}

// internal shapes mirror the external API responses we care about.
type weaponResp struct {
	Name     string `json:"name"`
	Category string `json:"weapon_category"`
	Range    struct {
		Normal int `json:"normal"`
	} `json:"range"`
	TwoHanded bool `json:"two_handed"`
}

type armorResp struct {
	Name       string `json:"name"`
	ArmorClass struct {
		Base     int  `json:"base"`
		DexBonus bool `json:"dex_bonus"`
	} `json:"armor_class"`
}

type spellResp struct {
	Name   string `json:"name"`
	Range  string `json:"range"`
	School struct {
		Name string `json:"name"`
	} `json:"school"`
}

func (a *APIAdapter) fetchJSON(url string, out interface{}) error {
	resp, err := a.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api returned status %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// GetWeapon fetches weapon details and maps to ports.WeaponInfo.
func (a *APIAdapter) GetWeapon(name string) (*ports.WeaponInfo, error) {
	idx := toAPIIndex(name)
	url := fmt.Sprintf("%s/equipment/%s", a.baseURL, idx)
	var w weaponResp
	if err := a.fetchJSON(url, &w); err != nil {
		return nil, err
	}
	return &ports.WeaponInfo{
		Name:      w.Name,
		Category:  w.Category,
		Range:     w.Range.Normal,
		TwoHanded: w.TwoHanded,
	}, nil
}

// GetArmor fetches armor details and maps to ports.ArmorInfo.
func (a *APIAdapter) GetArmor(name string) (*ports.ArmorInfo, error) {
	idx := toAPIIndex(name)
	url := fmt.Sprintf("%s/equipment/%s", a.baseURL, idx)
	var ar armorResp
	if err := a.fetchJSON(url, &ar); err != nil {
		return nil, err
	}
	return &ports.ArmorInfo{
		Name:     ar.Name,
		BaseAC:   ar.ArmorClass.Base,
		DexBonus: ar.ArmorClass.DexBonus,
	}, nil
}

// GetSpell fetches spell details and maps to ports.SpellInfo.
func (a *APIAdapter) GetSpell(name string) (*ports.SpellInfo, error) {
	idx := toAPIIndex(name)
	url := fmt.Sprintf("%s/spells/%s", a.baseURL, idx)
	var sp spellResp
	if err := a.fetchJSON(url, &sp); err != nil {
		return nil, err
	}
	return &ports.SpellInfo{
		Name:   sp.Name,
		Range:  sp.Range,
		School: sp.School.Name,
	}, nil
}
