package application

import (
	"reflect"
	"testing"

	bg "modules/dndcharactersheet/internal/domain/background"
	cl "modules/dndcharactersheet/internal/domain/class"
)

func TestCombineSkillProficiencies_IncludesRacialAndSorts(t *testing.T) {
	builder := NewCharacterBuilder(nil) // nil enricher - will fallback to domain

	background := bg.Background{
		Name:               "acolyte",
		SkillProficiencies: []string{"insight", "religion"},
	}
	class := cl.Class{
		Name:               "rogue",
		SkillProficiencies: []string{"acrobatics", "athletics", "deception", "insight"},
		SkillCount:         2,
	}
	user := []string{"athletics"}

	// dwarf grants history
	got := builder.CombineSkillProficiencies("dwarf", background, class, user)

	// Expect: first two class skills (acrobatics, athletics), then user (athletics), then background (insight, religion), then racial (history) -> sorted
	expect := []string{"acrobatics", "athletics", "athletics", "history", "insight", "religion"}
	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("expected %v, got %v", expect, got)
	}
}
