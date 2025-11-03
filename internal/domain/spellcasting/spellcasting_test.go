package spellcasting

import "testing"

func TestLearnSpell_Duplicate(t *testing.T) {
	sc := NewSpellcasting("sorcerer", 3) // known-spell caster

	if added := sc.LearnSpell("fireball"); !added {
		t.Fatalf("expected first learn to add spell")
	}
	if added := sc.LearnSpell("fireball"); added {
		t.Fatalf("expected second learn to be ignored (duplicate)")
	}
	if len(sc.KnownSpells) != 1 || sc.KnownSpells[0] != "fireball" {
		t.Fatalf("known spells mismatch: %#v", sc.KnownSpells)
	}
}

func TestPrepareSpell_RequiresKnown(t *testing.T) {
	sc := NewSpellcasting("wizard", 3) // prepared caster, but domain requires known

	if err := sc.PrepareSpell("magic missile"); err == nil {
		t.Fatalf("expected error when preparing unknown spell")
	}

	sc.LearnSpell("magic missile")

	if err := sc.PrepareSpell("magic missile"); err != nil {
		t.Fatalf("expected prepare to succeed after learning, got %v", err)
	}
	if err := sc.PrepareSpell("magic missile"); err == nil {
		t.Fatalf("expected duplicate prepare to fail")
	}
	if len(sc.PreparedSpells) != 1 || sc.PreparedSpells[0] != "magic missile" {
		t.Fatalf("prepared spells mismatch: %#v", sc.PreparedSpells)
	}
}
