package ports

// Enricher DTOs and interfaces for external data (weapons, armor, spells).

// WeaponInfo contains basic weapon enrichment data.
type WeaponInfo struct {
	Name      string
	Category  string
	Range     int
	TwoHanded bool
}

// ArmorInfo contains basic armor enrichment data.
type ArmorInfo struct {
	Name     string
	BaseAC   int
	DexBonus bool
}

// SpellInfo contains basic spell enrichment data.
type SpellInfo struct {
	Name   string
	Range  string
	School string
}

// TraitInfo contains racial trait enrichment data.
type TraitInfo struct {
	Index string
	Name  string
	Desc  []string // Description paragraphs
}

// WeaponEnricher fetches weapon details from an external service.
type WeaponEnricher interface {
	GetWeapon(name string) (*WeaponInfo, error)
	// Concurrent batch fetch with polite rate limit (maxPerSecond ~5-10 as per requirements)
	GetWeaponsBatch(names []string, maxPerSecond int) (map[string]*WeaponInfo, error)
}

// ArmorEnricher fetches armor details from an external service.
type ArmorEnricher interface {
	GetArmor(name string) (*ArmorInfo, error)
	GetArmorsBatch(names []string, maxPerSecond int) (map[string]*ArmorInfo, error)
}

// SpellEnricher fetches spell details from an external service.
type SpellEnricher interface {
	GetSpell(name string) (*SpellInfo, error)
	GetSpellsBatch(names []string, maxPerSecond int) (map[string]*SpellInfo, error)
}

// RaceEnricher fetches racial trait details from an external service.
type RaceEnricher interface {
	GetRacialSkillProficiencies(race string) ([]string, error)
	GetRacialTraits(race string) ([]*TraitInfo, error)
}
