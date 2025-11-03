# DDD / Onion Architecture

This D&D 5e Character Sheet application follows Domain-Driven Design (DDD) principles organized in an Onion Architecture pattern.

## Architecture Layers

### 1. Domain Layer (`internal/domain/`)
The core business logic with zero external dependencies.

#### Character Entity (`character/`)
- **`character.go`**: Pure domain model
  - `Character` struct with business rules
  - `ComputeModifiers()`: ability score modifiers
  - `ComputeDerived()`: default derived stats (minimal, extensible)
  - `ApplyRacialBonuses()`: racial stat bonuses
  - `GetProficiencyBonus()`: proficiency by level

#### Spellcasting Entity (`spellcasting/`)
- **`spellcasting.go`**: Spellcasting domain model
  - `Spellcasting` struct with CasterType, KnownSpells, PreparedSpells, SpellSlots
  - `NewSpellcasting(class, level)`: creates spellcasting for a character
  - `LearnSpell(spellName)`: adds to known spells (prevents duplicates)
  - `PrepareSpell(spellName)`: adds to prepared spells (requires known, prevents duplicates)
  - Domain errors: `ErrSpellNotKnown`, `ErrSpellAlreadyPrepared`

- **`rules.go`**: Spellcasting business rules
  - `GetCasterType(class)`: determines caster progression (full/half/pact/known/none)
  - `GetSpellSlots(casterType, level)`: returns spell slots by level
  - `GetCantripsKnown(class, level)`: cantrips available per class/level
  - Spell slot tables for full/half/pact casters

#### Reference Data (`background/`, `class/`)
- **`background/model.go`**: Background definitions with skill proficiencies
- **`class/model.go`**: Class definitions with skill counts and proficiencies

**Rules:**
- No imports from other layers (no ports, adapters, application)
- Pure business logic only
- No JSON tags or infrastructure concerns
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

- **`spell_repository.go`**: Spell data access contract
  - `LoadSpells()`: loads all spell data
  - `FilterByClass(spells, class)`: filters spells for a class
  - `Spell` type with Index, Name, Level, Classes

**Rules:**
- Interfaces only (no implementations)
- Defines contracts for adapters
- Technology-agnostic

### 4. Adapters Layer (`internal/adapters/`)
Implements ports using specific technologies.

#### Storage Adapter (`adapters/storage/`)
- **`repository.go`**: Implements `CharacterRepository` port
  - Maps between domain `Character` and storage model
  - `unmarshalSpellcasting()`: reconstructs typed spellcasting from JSON
  - Delegates to `jsonstorage/` backend

- **`jsonstorage/`**: JSON file storage implementation
  - `model.go`: storage types (`Character`, `CharacterSummary`, `CharacterStorage` interface)
  - `single_file_storage.go`: JSON file CRUD operations
  - No domain knowledge, pure infrastructure

- **`repository_spellcasting_test.go`**: Tests spellcasting persistence round-trip

#### API Adapter (`adapters/api/`)
- **`client_adapter.go`**: Implements enricher ports
  - `APIAdapter` calls D&D 5e REST API (configurable base URL)
  - `GetWeapon`, `GetArmor`, `GetSpell` methods
  - Implements `WeaponEnricher`, `ArmorEnricher`, `SpellEnricher`

#### Spellcasting Adapter (`adapters/spellcasting/`)
- **`engine_adapter.go`**: Implements `SpellcastingEngine` port
  - Uses domain `spellcasting.NewSpellcasting()` for initialization
  - Uses domain `LearnSpell()` and `PrepareSpell()` methods
  - Validates spells via `SpellRepository` (CSV)
  - Formats spell slots and cantrips for display using domain rules

- **`spell_repository.go`**: Implements `SpellRepository` port
  - `CSVSpellRepository`: loads spells from CSV file
  - Filters spells by class
  - Pure infrastructure (CSV parsing)

**Rules:**
- Implements port interfaces
- Contains technology-specific code (HTTP, JSON, CSV)
- No business logic (delegates to domain)

### 5. Presentation Layer (`main.go`)
CLI interface - the outer layer.

- Command parsing (create, view, list, delete, equip, learn-spell, prepare-spell)
- Wires dependencies: creates adapters (storage, API, spellcasting), injects into application service
- Calls application service for all operations
- Displays results (including type-asserting domain spellcasting for known/prepared spell display)

**Rules:**
- Only depends on application service and adapters
- No direct domain manipulation
- Wires dependencies (Dependency Injection)

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

## Folder Structure

```
internal/
├── domain/                       ← Core business logic (no external deps)
│   ├── character/
│   │   └── character.go         ← Character entity with business rules
│   ├── spellcasting/
│   │   ├── spellcasting.go      ← Spellcasting entity (learn/prepare logic)
│   │   ├── rules.go             ← Caster types, slot tables, cantrips
│   │   └── spellcasting_test.go ← Domain unit tests
│   ├── background/
│   │   └── model.go             ← Background reference data
│   └── class/
│       └── model.go             ← Class reference data
│
├── application/                  ← Use case orchestration
│   ├── character_service.go     ← CRUD, equipment, spells, derived stats
│   └── character_builder.go     ← Skill combining logic
│
├── ports/                        ← Interface contracts
│   ├── character_repository.go  ← Persistence port
│   ├── enrichers.go             ← Equipment/spell data ports
│   ├── spellcasting.go          ← Spellcasting engine port
│   └── spell_repository.go      ← Spell data access port
│
└── adapters/                     ← Infrastructure implementations
    ├── storage/
    │   ├── repository.go        ← CharacterRepository implementation
    │   ├── repository_test.go   ← CRUD tests
    │   ├── repository_spellcasting_test.go  ← Round-trip tests
    │   └── jsonstorage/
    │       ├── model.go         ← Storage schema
    │       └── single_file_storage.go  ← JSON file backend
    ├── api/
    │   ├── client_adapter.go    ← D&D 5e API enrichers
    │   └── client_adapter_test.go
    └── spellcasting/
        ├── engine_adapter.go    ← Spellcasting engine (domain + CSV)
        └── spell_repository.go  ← CSV spell loader
```

## Legacy Code Status

**Removed (refactored into DDD layers):**
- ✅ `internal/combat` → `application/character_service.RecalculateDerived`
- ✅ `internal/equipment` → `adapters/api/client_adapter`
- ✅ `internal/character` → `domain/character` + `application/character_builder`
- ✅ `internal/api` → `adapters/api/client_adapter`
- ✅ `internal/spellcasting` → `domain/spellcasting` (rules/entity) + `adapters/spellcasting` (CSV/engine)
- ✅ `internal/storage` → `adapters/storage/jsonstorage`
- ✅ `internal/background` → `domain/background`
- ✅ `internal/class` → `domain/class`

**Result:** Clean DDD/Onion structure with all code in appropriate layers.

## Testing Strategy

- **Domain**: Pure unit tests (no mocks needed)
  - `spellcasting_test.go`: tests learn/prepare logic and edge cases
- **Application**: Test with mock ports
- **Adapters**: Integration tests with real/mock external systems
  - `repository_test.go`: CRUD with temp file
  - `repository_spellcasting_test.go`: spellcasting persistence round-trip
  - `client_adapter_test.go`: API calls (requires local server)

## Next Steps

1. Replace spellcasting adapter with domain-based engine
2. Consolidate API clients (merge `internal/api` into `adapters/api`)
3. Remove remaining legacy storage references
4. Add comprehensive tests for application layer
