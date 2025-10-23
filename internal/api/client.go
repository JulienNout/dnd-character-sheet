package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const BaseURL = "http://localhost:3000/api"

// WeaponEnriched holds extra weapon info from the API
type WeaponEnriched struct {
	Name     string `json:"name"`
	Category string `json:"weapon_category"`
	Range    struct {
		Normal int `json:"normal"`
	} `json:"range"`
	TwoHanded bool `json:"two_handed"`
}

// ArmorEnriched holds extra armor info from the API
type ArmorEnriched struct {
	Name       string `json:"name"`
	ArmorClass struct {
		Base     int  `json:"base"`
		DexBonus bool `json:"dex_bonus"`
	} `json:"armor_class"`
}

// GetWeapon fetches and decodes weapon details by index (e.g., "longsword")
func GetWeapon(index string) (*WeaponEnriched, error) {
	url := fmt.Sprintf("%s/equipment/%s", BaseURL, index)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %s", resp.Status)
	}

	var weapon WeaponEnriched
	if err := json.NewDecoder(resp.Body).Decode(&weapon); err != nil {
		return nil, err
	}
	return &weapon, nil
}

// GetArmor fetches and decodes armor details by index (e.g., "chain-mail")
func GetArmor(index string) (*ArmorEnriched, error) {
	url := fmt.Sprintf("%s/equipment/%s", BaseURL, index)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %s", resp.Status)
	}

	var armor ArmorEnriched
	if err := json.NewDecoder(resp.Body).Decode(&armor); err != nil {
		return nil, err
	}
	return &armor, nil
}

// FetchWeaponsWithWorkers fetches weapon details for a list of indexes using a worker pool
func FetchWeaponsWithWorkers(indexes []string, workerCount int) []*WeaponEnriched {
	type job struct {
		i   int
		idx string
	}
	type result struct {
		i      int
		weapon *WeaponEnriched
		err    error
	}

	jobs := make(chan job, len(indexes))
	results := make(chan result, len(indexes))

	for w := 0; w < workerCount; w++ {
		go func() {
			for j := range jobs {
				weapon, err := GetWeapon(j.idx)
				results <- result{j.i, weapon, err}
			}
		}()
	}

	for i, idx := range indexes {
		jobs <- job{i, idx}
	}
	close(jobs)

	out := make([]*WeaponEnriched, len(indexes))
	for i := 0; i < len(indexes); i++ {
		res := <-results
		out[res.i] = res.weapon
		if res.err != nil {
			fmt.Printf("Failed to fetch weapon %s: %v\n", indexes[res.i], res.err)
		}
	}
	return out
}

// FetchArmorsWithWorkers fetches armor details for a list of indexes using a worker pool
func FetchArmorsWithWorkers(indexes []string, workerCount int) []*ArmorEnriched {
	type job struct {
		i   int
		idx string
	}
	type result struct {
		i     int
		armor *ArmorEnriched
		err   error
	}

	jobs := make(chan job, len(indexes))
	results := make(chan result, len(indexes))

	for w := 0; w < workerCount; w++ {
		go func() {
			for j := range jobs {
				armor, err := GetArmor(j.idx)
				results <- result{j.i, armor, err}
			}
		}()
	}

	for i, idx := range indexes {
		jobs <- job{i, idx}
	}
	close(jobs)

	out := make([]*ArmorEnriched, len(indexes))
	for i := 0; i < len(indexes); i++ {
		res := <-results
		out[res.i] = res.armor
		if res.err != nil {
			fmt.Printf("Failed to fetch armor %s: %v\n", indexes[res.i], res.err)
		}
	}
	return out
}

// SpellEnriched holds extra spell info from the API
type SpellEnriched struct {
	Name   string `json:"name"`
	Range  string `json:"range"`
	School struct {
		Name string `json:"name"`
	} `json:"school"`
}

// GetSpell fetches and decodes spell details by index (e.g., "acid-arrow")
func GetSpell(index string) (*SpellEnriched, error) {
	url := fmt.Sprintf("%s/spells/%s", BaseURL, index)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %s", resp.Status)
	}

	var spell SpellEnriched
	if err := json.NewDecoder(resp.Body).Decode(&spell); err != nil {
		return nil, err
	}
	return &spell, nil
}

// FetchSpellsWithWorkers fetches spell details for a list of indexes using a worker pool
func FetchSpellsWithWorkers(indexes []string, workerCount int) []*SpellEnriched {
	type job struct {
		i   int
		idx string
	}
	type result struct {
		i     int
		spell *SpellEnriched
		err   error
	}

	jobs := make(chan job, len(indexes))
	results := make(chan result, len(indexes))

	for w := 0; w < workerCount; w++ {
		go func() {
			for j := range jobs {
				spell, err := GetSpell(j.idx)
				results <- result{j.i, spell, err}
			}
		}()
	}

	for i, idx := range indexes {
		jobs <- job{i, idx}
	}
	close(jobs)

	out := make([]*SpellEnriched, len(indexes))
	for i := 0; i < len(indexes); i++ {
		res := <-results
		out[res.i] = res.spell
		if res.err != nil {
			fmt.Printf("Failed to fetch spell %s: %v\n", indexes[res.i], res.err)
		}
	}
	return out
}
