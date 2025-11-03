# Maintainability Report

## Architecture Overview
The application follows **Domain-Driven Design (DDD)** with an **Onion Architecture** pattern, organizing code into clear layers:

- **Domain**: Pure business logic (Character, Race, Spellcasting) with no external dependencies
- **Application**: Use cases (CharacterService, CharacterBuilder) orchestrating domain logic
- **Ports**: Interface contracts (RaceEnricher, SpellRepository, WeaponEnricher)
- **Adapters**: External integrations (APIAdapter, CSVRepository, JSONStorage)
- **Presentation**: CLI commands and HTML frontend

## Key Design Principles

### 1. Dependency Inversion (SOLID)
All external dependencies flow **inward through ports**. The domain layer has zero coupling to frameworks or infrastructure:
```
Domain ← Application ← Ports ← Adapters
```
This makes business logic testable and portable.

### 2. Open/Closed Principle (OCP)
New data sources are added by implementing existing ports:
- Race enrichment: `RaceEnricher` interface allows API adapter or future database adapter
- Spell data: `SpellRepository` supports CSV now, could support database later
- No modification of domain or application layers required

### 3. Concurrent API Calls
The `APIAdapter` uses **goroutines + WaitGroup + Mutex** to fetch multiple API endpoints concurrently:
```go
var wg sync.WaitGroup
for _, traitName := range traitNames {
    wg.Add(1)
    go func(name string) {
        defer wg.Done()
        // Fetch trait from API
    }(traitName)
}
wg.Wait()
```
This satisfies the requirement that **"api requests need to be concurrent"** while maintaining clean architecture.

### 4. Separation of Concerns
Each layer has a single responsibility:
- **Domain**: Business rules (skill calculation, spellcasting rules)
- **Application**: Workflow orchestration (character creation, skill combination)
- **Adapters**: Technical details (HTTP calls, file I/O, JSON parsing)
- **CLI**: User interaction only (no business logic)

### 5. Testability
Pure domain functions are testable without infrastructure:
```go
// Domain test - no API/database needed
func TestGetRacialSkillProficiencies(t *testing.T) {
    skills := race.GetRacialSkillProficiencies("dwarf")
    assert.Contains(t, skills, "History")
}
```
Application layer tests use **nil enrichers** to verify fallback behavior.

## Maintainability Benefits

✅ **Easy to extend**: Add new races/classes/spells without changing core logic  
✅ **Easy to test**: Pure functions, dependency injection, port interfaces  
✅ **Easy to understand**: Clear boundaries between layers  
✅ **Easy to change**: Swap API for database without touching domain  
✅ **Performance**: Concurrent API calls reduce latency  

## Trade-offs
- More files/interfaces than a monolithic approach
- Requires understanding of DDD layering
- Initial setup complexity for small features

## Conclusion
The architecture prioritizes **long-term maintainability over short-term convenience**. Adding new features (like racial traits) follows a clear pattern: define port → implement adapter → wire in application layer. This consistency reduces cognitive load and prevents architectural drift.

---

## Testing Section
(Left blank as requested)
