# DDD / Onion Architecture

This D&D 5e Character Sheet application follows Domain-Driven Design (DDD) principles organized in an Onion Architecture pattern.

## Architecture Layers

### 1. Domain Layer (`internal/domain/`)
The core business logic with zero external dependencies.

- **`character/character.go`**: Pure domain model
  - `Character` struct with business rules
  - `ComputeModifiers()`: ability score modifiers
  - `ComputeDerived()`: default derived stats (minimal, extensible)
  - `ApplyRacialBonuses()`: racial stat bonuses
  - `GetProficiencyBonus()`: proficiency by level

**Rules:**
- No imports from other layers
- Pure business logic only
- Immutable where possible

### 2. Application Layer (`internal/application/`)
Orchestrates use cases and coordinates domain logic with infrastructure.

- **`character_service.go`**: Main service orchestrating character operations
  - CRUD operations (Create, Get, List, Delete)
  - `EquipWeapon/Armor/Shield`: equipment management via ports
  - `LearnSpell/PrepareSpell`: spellcasting via ports
  - `RecalculateDerived`: recomputes AC/initiative/passive perception using enrichers and business rules

- **`character_builder.go`**: Builds complex character state
  - `CombineSkillProficiencies`: merges background/class/user skill selections

**Rules:**
- Depends on domain and ports (interfaces only)
- Orchestrates workflows
- No direct adapter implementation details

### 3. Ports Layer (`internal/ports/`)
Defines interfaces (contracts) for external systems.

- **`character_repository.go`**: Persistence contract
  - `Save`, `GetByID`, `GetAll`, `Delete`

- **`enrichers.go`**: External data enrichment contracts
  - `WeaponEnricher`, `ArmorEnricher`, `SpellEnricher`
  - `WeaponInfo`, `ArmorInfo`, `SpellInfo` types

- **`spellcasting.go`**: Spellcasting engine contract
  - `AssignSpellcasting`, `LearnSpell`, `PrepareSpell`
  - `FormatSpellSlots`, `FormatCantrips`

**Rules:**
- Interfaces only (no implementations)
- Defines contracts for adapters
- Technology-agnostic

### 4. Adapters Layer (`internal/adapters/`)
Implements ports using specific technologies.

#### Storage Adapter (`adapters/storage/`)
- **`repository.go`**: Implements `CharacterRepository` port
  - Maps between domain `Character` and storage `Character`
  - Delegates to `internal/storage` backend

#### API Adapter (`adapters/api/`)
- **`client_adapter.go`**: Implements enricher ports
  - `APIAdapter` calls D&D 5e REST API
  - `GetWeapon`, `GetArmor`, `GetSpell` methods
  - Configurable base URL

#### Spellcasting Adapter (`adapters/spellcasting/`)
- **`engine_adapter.go`**: Implements `SpellcastingEngine` port
  - Temporary bridge to legacy CSV-based spellcasting
  - Formats spell slots and cantrips for display

**Rules:**
- Implements port interfaces
- Contains technology-specific code
- No business logic

### 5. Infrastructure (`internal/storage/`)
Low-level persistence implementation.

- **`model.go`**: Storage types
  - `Character`: JSON persistence schema
  - `CharacterSummary`: list view
  - `CharacterStorage`: storage interface

- **`single_file_storage.go`**: File-based storage
  - JSON file operations
  - CRUD for characters

**Rules:**
- Owned by adapters (not directly used by application)
- Storage-specific concerns only

### 6. Presentation Layer (`main.go`)
CLI interface - the outer layer.

- Command parsing (create, view, list, delete, equip, learn-spell, prepare-spell)
- Calls application service
- Displays results

**Rules:**
- Only depends on application service and adapters
- No direct domain manipulation
- Wires dependencies (DI)

## Dependency Flow

```
main.go (CLI)
    ↓
application/ (orchestration)
    ↓
ports/ (interfaces) ← domain/ (business logic)
    ↑
adapters/ (implementations)
    ↓
storage/, external APIs
```

**Key principle**: Dependencies point inward. Domain has zero dependencies. Application depends only on domain + ports. Adapters implement ports.

## Design Patterns

1. **Repository Pattern**: `CharacterRepository` abstracts persistence
2. **Adapter Pattern**: Enrichers adapt external APIs to domain needs
3. **Dependency Injection**: main.go wires concrete implementations
4. **Open/Closed Principle**: Add features by implementing new ports/adapters without changing existing code

## Legacy Code Status

**Removed (refactored into DDD layers):**
- ✅ `internal/combat` → `application/character_service.RecalculateDerived`
- ✅ `internal/equipment` → `adapters/api/client_adapter`

**Remaining (to be addressed):**
- `internal/spellcasting`: Bridged via adapter; can be replaced with domain-based engine
- `internal/character`: Legacy storage model (minimal usage remaining)
- `internal/api`: Legacy API client (overlaps with adapters/api)
- `internal/background`, `internal/class`: Data loaders (low priority)

## Testing Strategy

- **Domain**: Pure unit tests (no mocks needed)
- **Application**: Test with mock ports
- **Adapters**: Integration tests with real/mock external systems
- **Storage adapter**: CRUD test with temp file (`repository_test.go`)

## Next Steps

1. Replace spellcasting adapter with domain-based engine
2. Consolidate API clients (merge `internal/api` into `adapters/api`)
3. Remove remaining legacy storage references
4. Add comprehensive tests for application layer
