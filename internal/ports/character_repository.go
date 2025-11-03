package ports

import characterpkg "modules/dndcharactersheet/internal/domain/character"

// CharacterRepository defines persistence operations for domain Characters.
// Adapters in internal/adapters/* should implement this interface.
type CharacterRepository interface {
	Save(c *characterpkg.Character) error
	GetAll() ([]*characterpkg.Character, error)
	GetByID(id string) (*characterpkg.Character, error)
	Delete(id string) error
}
