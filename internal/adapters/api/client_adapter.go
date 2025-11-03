package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
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

// GetRacialSkillProficiencies fetches racial traits from the API and maps them to skill proficiencies.
// It queries known SRD traits concurrently and returns matching skills for the given race.
func (a *APIAdapter) GetRacialSkillProficiencies(race string) ([]string, error) {
	rNorm := normalizeRace(race)

	// Known SRD trait -> skill mappings
	checks := []struct {
		path  string
		skill string
	}{
		{"traits/stonecunning", "history"},   // dwarves
		{"traits/keen-senses", "perception"}, // elves
		{"traits/menacing", "intimidation"},  // half-orcs
	}

	type result struct {
		skill string
		match bool
	}

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []string
	)

	wg.Add(len(checks))
	for _, c := range checks {
		c := c // capture range var
		go func() {
			defer wg.Done()
			ok, err := a.raceHasTrait(rNorm, c.path)
			if err != nil || !ok {
				return
			}
			mu.Lock()
			results = append(results, c.skill)
			mu.Unlock()
		}()
	}
	wg.Wait()

	return results, nil
}

type traitResp struct {
	Index string `json:"index"`
	Races []struct {
		Index string `json:"index"`
		Name  string `json:"name"`
	} `json:"races"`
}

func (a *APIAdapter) raceHasTrait(race string, traitPath string) (bool, error) {
	url := fmt.Sprintf("%s/%s", a.baseURL, strings.TrimLeft(traitPath, "/"))
	var tr traitResp
	if err := a.fetchJSON(url, &tr); err != nil {
		return false, err
	}
	for _, rr := range tr.Races {
		if normalizeRace(rr.Index) == race || normalizeRace(rr.Name) == race {
			return true, nil
		}
	}
	return false, nil
}

func normalizeRace(v string) string {
	s := strings.ToLower(strings.TrimSpace(v))
	s = strings.ReplaceAll(s, "-", " ")
	return s
}

// GetRacialTraits fetches trait descriptions from the API concurrently for a given race.
func (a *APIAdapter) GetRacialTraits(race string) ([]*ports.TraitInfo, error) {
	rNorm := normalizeRace(race)

	// Known SRD traits that map to races
	traitPaths := []string{
		"traits/stonecunning",            // dwarf
		"traits/darkvision",              // many races
		"traits/dwarven-resilience",      // dwarf
		"traits/dwarven-combat-training", // dwarf
		"traits/keen-senses",             // elf
		"traits/fey-ancestry",            // elf
		"traits/trance",                  // elf
		"traits/menacing",                // half-orc
		"traits/relentless-endurance",    // half-orc
		"traits/savage-attacks",          // half-orc
		"traits/brave",                   // halfling
		"traits/halfling-nimbleness",     // halfling
		"traits/lucky",                   // halfling
	}

	type traitWithDesc struct {
		Index string   `json:"index"`
		Name  string   `json:"name"`
		Desc  []string `json:"desc"`
		Races []struct {
			Index string `json:"index"`
			Name  string `json:"name"`
		} `json:"races"`
	}

	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		traits []*ports.TraitInfo
	)

	wg.Add(len(traitPaths))
	for _, path := range traitPaths {
		path := path
		go func() {
			defer wg.Done()
			url := fmt.Sprintf("%s/%s", a.baseURL, strings.TrimLeft(path, "/"))
			var tr traitWithDesc
			if err := a.fetchJSON(url, &tr); err != nil {
				return
			}
			// Check if this trait belongs to the character's race
			for _, rr := range tr.Races {
				if normalizeRace(rr.Index) == rNorm || normalizeRace(rr.Name) == rNorm {
					mu.Lock()
					traits = append(traits, &ports.TraitInfo{
						Index: tr.Index,
						Name:  tr.Name,
						Desc:  tr.Desc,
					})
					mu.Unlock()
					break
				}
			}
		}()
	}
	wg.Wait()

	return traits, nil
}

// GetWeaponsBatch fetches multiple weapons concurrently with a simple rate limiter.
func (a *APIAdapter) GetWeaponsBatch(names []string, maxPerSecond int) (map[string]*ports.WeaponInfo, error) {
	if maxPerSecond <= 0 {
		maxPerSecond = 5
	}
	res := make(map[string]*ports.WeaponInfo, len(names))
	var mu sync.Mutex
	var wg sync.WaitGroup
	interval := time.Second / time.Duration(maxPerSecond)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for _, n := range names {
		name := n
		<-ticker.C
		wg.Add(1)
		go func() {
			defer wg.Done()
			if info, err := a.GetWeapon(name); err == nil {
				mu.Lock()
				res[name] = info
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return res, nil
}

// GetArmorsBatch fetches multiple armors concurrently with a simple rate limiter.
func (a *APIAdapter) GetArmorsBatch(names []string, maxPerSecond int) (map[string]*ports.ArmorInfo, error) {
	if maxPerSecond <= 0 {
		maxPerSecond = 5
	}
	res := make(map[string]*ports.ArmorInfo, len(names))
	var mu sync.Mutex
	var wg sync.WaitGroup
	interval := time.Second / time.Duration(maxPerSecond)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for _, n := range names {
		name := n
		<-ticker.C
		wg.Add(1)
		go func() {
			defer wg.Done()
			if info, err := a.GetArmor(name); err == nil {
				mu.Lock()
				res[name] = info
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return res, nil
}

// GetSpellsBatch fetches multiple spells concurrently with a simple rate limiter.
func (a *APIAdapter) GetSpellsBatch(names []string, maxPerSecond int) (map[string]*ports.SpellInfo, error) {
	if maxPerSecond <= 0 {
		maxPerSecond = 5
	}
	res := make(map[string]*ports.SpellInfo, len(names))
	var mu sync.Mutex
	var wg sync.WaitGroup
	interval := time.Second / time.Duration(maxPerSecond)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for _, n := range names {
		name := n
		<-ticker.C
		wg.Add(1)
		go func() {
			defer wg.Done()
			if info, err := a.GetSpell(name); err == nil {
				mu.Lock()
				res[name] = info
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return res, nil
}
