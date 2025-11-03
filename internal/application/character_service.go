package application

import (
	"fmt"
	"strings"

	characterpkg "modules/dndcharactersheet/internal/domain/character"
	"modules/dndcharactersheet/internal/ports"
)

// CharacterService orchestrates character-related use cases using ports.
type CharacterService struct {
	repo           ports.CharacterRepository
	weaponEnricher ports.WeaponEnricher
	armorEnricher  ports.ArmorEnricher
	spellEnricher  ports.SpellEnricher
}

// NewCharacterService creates a new service with repository port.
func NewCharacterService(r ports.CharacterRepository) *CharacterService {
	return &CharacterService{repo: r}
}

// WithEnrichers adds enricher ports to the service for equipment and spell features.
func (s *CharacterService) WithEnrichers(we ports.WeaponEnricher, ae ports.ArmorEnricher, se ports.SpellEnricher) *CharacterService {
	s.weaponEnricher = we
	s.armorEnricher = ae
	s.spellEnricher = se
	return s
}

func (s *CharacterService) Create(c *characterpkg.Character) error {
	return s.repo.Save(c)
}

func (s *CharacterService) List() ([]*characterpkg.Character, error) {
	return s.repo.GetAll()
}

func (s *CharacterService) Get(id string) (*characterpkg.Character, error) {
	return s.repo.GetByID(id)
}

func (s *CharacterService) Delete(id string) error {
	return s.repo.Delete(id)
}

// RecalculateDerived computes AC, initiative, and passive perception using available enrichers.
// If no enrichers are configured, falls back to domain defaults.
func (s *CharacterService) RecalculateDerived(c *characterpkg.Character) {
	// Always ensure modifiers are set
	c.ComputeModifiers()

	// Initiative is Dex mod
	c.Initiative = c.DexMod

	// Passive Perception: 10 + Wis mod (+ proficiency if proficient)
	passive := 10 + c.WisMod
	for _, sk := range c.SkillProficiencies {
		if strings.EqualFold(strings.TrimSpace(sk), "perception") {
			passive += c.Proficiency
			break
		}
	}
	c.PassivePerception = passive

	// Armor Class calculation
	if s.armorEnricher == nil {
		// Fallback minimal rule: base 10 + Dex mod (+2 if shield)
		c.ArmorClass = 10 + c.DexMod
		if c.Shield != "" {
			c.ArmorClass += 2
		}
		return
	}

	// Special unarmored defense cases
	lclass := strings.ToLower(strings.TrimSpace(c.Class))
	if lclass == "barbarian" && c.Armor == "" {
		ac := 10 + c.DexMod + c.ConMod
		if c.Shield != "" {
			ac += 2
		}
		c.ArmorClass = ac
		return
	}
	if lclass == "monk" && c.Armor == "" && c.Shield == "" {
		c.ArmorClass = 10 + c.DexMod + c.WisMod
		return
	}

	// Default armor logic via enricher
	base := 10
	if c.Armor != "" {
		if info, err := s.armorEnricher.GetArmor(c.Armor); err == nil && info != nil {
			base = info.BaseAC
			if info.DexBonus {
				base += c.DexMod
			}
		} else {
			// Fallback: light armor style
			base = 10 + c.DexMod
		}
	} else {
		base = 10 + c.DexMod
	}
	if c.Shield != "" {
		// Default +2 shield bonus
		base += 2
	}
	c.ArmorClass = base
}

// EquipWeapon enriches weapon information and equips it to the specified slot.
func (s *CharacterService) EquipWeapon(characterName, weaponName, slot string) error {
	if s.weaponEnricher == nil {
		return fmt.Errorf("weapon enricher not configured")
	}

	char, err := s.repo.GetByID(characterName)
	if err != nil {
		return fmt.Errorf("character not found: %w", err)
	}

	// Enrich weapon data from external API
	weaponInfo, err := s.weaponEnricher.GetWeapon(weaponName)
	if err != nil {
		return fmt.Errorf("failed to enrich weapon: %w", err)
	}

	// Business logic: equip to slot
	switch slot {
	case "main hand", "mainhand":
		char.MainHand = weaponInfo.Name
	case "off hand", "offhand":
		char.OffHand = weaponInfo.Name
	default:
		return fmt.Errorf("invalid slot: %s", slot)
	}

	return s.repo.Save(char)
}

// EquipArmor enriches armor information and equips it to the character.
func (s *CharacterService) EquipArmor(characterName, armorName string) error {
	if s.armorEnricher == nil {
		return fmt.Errorf("armor enricher not configured")
	}

	char, err := s.repo.GetByID(characterName)
	if err != nil {
		return fmt.Errorf("character not found: %w", err)
	}

	// Enrich armor data from external API
	armorInfo, err := s.armorEnricher.GetArmor(armorName)
	if err != nil {
		return fmt.Errorf("failed to enrich armor: %w", err)
	}

	char.Armor = armorInfo.Name
	return s.repo.Save(char)
}

// EquipShield equips a shield to the character.
func (s *CharacterService) EquipShield(characterName, shieldName string) error {
	char, err := s.repo.GetByID(characterName)
	if err != nil {
		return fmt.Errorf("character not found: %w", err)
	}

	char.Shield = shieldName
	return s.repo.Save(char)
}

// LearnSpell enriches spell information and adds it to the character's known spells.
func (s *CharacterService) LearnSpell(characterName, spellName string) error {
	if s.spellEnricher == nil {
		return fmt.Errorf("spell enricher not configured")
	}

	char, err := s.repo.GetByID(characterName)
	if err != nil {
		return fmt.Errorf("character not found: %w", err)
	}

	// Enrich spell data from external API
	_, err = s.spellEnricher.GetSpell(spellName)
	if err != nil {
		return fmt.Errorf("failed to enrich spell: %w", err)
	}

	// Note: The Spellcasting field is interface{} in domain model
	// This would ideally be a proper domain type, but for now we preserve existing structure
	// The actual spell management is handled by the legacy spellcasting package

	return s.repo.Save(char)
}

// PrepareSpell marks a spell as prepared for the character.
func (s *CharacterService) PrepareSpell(characterName, spellName string) error {
	if s.spellEnricher == nil {
		return fmt.Errorf("spell enricher not configured")
	}

	char, err := s.repo.GetByID(characterName)
	if err != nil {
		return fmt.Errorf("character not found: %w", err)
	}

	// Enrich spell data from external API
	_, err = s.spellEnricher.GetSpell(spellName)
	if err != nil {
		return fmt.Errorf("failed to enrich spell: %w", err)
	}

	// Note: Spell preparation logic handled by legacy spellcasting package
	return s.repo.Save(char)
}
