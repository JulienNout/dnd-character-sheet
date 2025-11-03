package race

import "testing"

func TestGetRacialSkillProficiencies(t *testing.T) {
	cases := []struct {
		name   string
		race   string
		expect []string
	}{
		{"dwarf basic", "dwarf", []string{"history"}},
		{"dwarf hyphen", "mountain-dwarf", []string{"history"}},
		{"elf basic", "elf", []string{"perception"}},
		{"elf subtype", "wood elf", []string{"perception"}},
		{"half orc space", "half orc", []string{"intimidation"}},
		{"half-orc hyphen", "half-orc", []string{"intimidation"}},
		{"unknown race", "human", nil},
	}

	for _, tc := range cases {
		// capture range variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := GetRacialSkillProficiencies(tc.race)
			if len(got) != len(tc.expect) {
				t.Fatalf("expected %v, got %v", tc.expect, got)
			}
			for i := range got {
				if got[i] != tc.expect[i] {
					t.Fatalf("expected %v, got %v", tc.expect, got)
				}
			}
		})
	}
}
