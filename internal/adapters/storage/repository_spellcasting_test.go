package storage

import (
    "os"
    "path/filepath"
    "testing"

    characterpkg "modules/dndcharactersheet/internal/domain/character"
    spkg "modules/dndcharactersheet/internal/domain/spellcasting"
)

func TestRepository_SpellcastingRoundTrip(t *testing.T) {
    dir := t.TempDir()
    file := filepath.Join(dir, "characters.json")

    repo := NewJSONRepository(file)

    // Create a character with spellcasting
    sc := spkg.NewSpellcasting("wizard", 3)
    sc.LearnSpell("magic missile")
    if err := sc.PrepareSpell("magic missile"); err != nil {
        t.Fatalf("prepare should succeed after learning: %v", err)
    }

    ch := &characterpkg.Character{
        Name:        "RepoSpellTest",
        Race:        "human",
        Class:       "wizard",
        Level:       3,
        Str:         10,
        Dex:         12,
        Con:         12,
        Int:         14,
        Wis:         10,
        Cha:         8,
        Spellcasting: sc,
    }

    if err := repo.Save(ch); err != nil {
        t.Fatalf("save failed: %v", err)
    }

    // Ensure file written
    if _, err := os.Stat(file); err != nil {
        t.Fatalf("expected file to exist: %v", err)
    }

    // Load back
    got, err := repo.GetByID("RepoSpellTest")
    if err != nil {
        t.Fatalf("get failed: %v", err)
    }
    if got.Spellcasting == nil {
        t.Fatalf("expected spellcasting to be non-nil after load")
    }

    // Type assertion to domain type
    sc2, ok := got.Spellcasting.(*spkg.Spellcasting)
    if !ok {
        t.Fatalf("expected spellcasting to be *spellcasting.Spellcasting, got %T", got.Spellcasting)
    }

    if len(sc2.KnownSpells) != 1 || sc2.KnownSpells[0] != "magic missile" {
        t.Fatalf("known spells after load mismatch: %#v", sc2.KnownSpells)
    }
    if len(sc2.PreparedSpells) != 1 || sc2.PreparedSpells[0] != "magic missile" {
        t.Fatalf("prepared spells after load mismatch: %#v", sc2.PreparedSpells)
    }

    // Check some slots present
    if len(sc2.SpellSlots) == 0 {
        t.Fatalf("expected spell slots to be present after load")
    }
}
