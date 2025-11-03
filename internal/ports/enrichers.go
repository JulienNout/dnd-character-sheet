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

// WeaponEnricher fetches weapon details from an external service.
type WeaponEnricher interface {
	GetWeapon(name string) (*WeaponInfo, error)
}

// ArmorEnricher fetches armor details from an external service.
type ArmorEnricher interface {
	GetArmor(name string) (*ArmorInfo, error)
}

// SpellEnricher fetches spell details from an external service.
type SpellEnricher interface {
	GetSpell(name string) (*SpellInfo, error)
}
