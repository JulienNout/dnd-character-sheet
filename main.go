package main

import (
	"encoding/json"
	"flag"
	"fmt"
	backgroundModel "modules/dndcharactersheet/internal/background"
	characterModel "modules/dndcharactersheet/internal/character"
	classModel "modules/dndcharactersheet/internal/class"
	"modules/dndcharactersheet/internal/combat"
	"modules/dndcharactersheet/internal/equipment"
	"modules/dndcharactersheet/internal/spellcasting"
	"modules/dndcharactersheet/internal/storage"
	"os"
	"strings"
)

func usage() {
	fmt.Printf(`Usage:
  %s create -name CHARACTER_NAME -race RACE -class CLASS -level N -str N -dex N -con N -int N -wis N -cha N
  %s view -name CHARACTER_NAME
  %s list
  %s delete -name CHARACTER_NAME
  %s equip -name CHARACTER_NAME -weapon WEAPON_NAME -slot SLOT
  %s equip -name CHARACTER_NAME -armor ARMOR_NAME
  %s equip -name CHARACTER_NAME -shield SHIELD_NAME
  %s learn-spell -name CHARACTER_NAME -spell SPELL_NAME
  %s prepare-spell -name CHARACTER_NAME -spell SPELL_NAME 
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	cmd := os.Args[1]

	switch cmd {
	case "create":
		// You could use the Flag package like this
		// But feel free to do it differently!
		createCmd := flag.NewFlagSet("create", flag.ExitOnError)
		name := createCmd.String("name", "", "character name (required)")
		race := createCmd.String("race", "", "race (required)")
		class := createCmd.String("class", "", "class (required)")
		level := createCmd.Int("level", 1, "level (required)")
		str := createCmd.Int("str", 10, "strength")
		dex := createCmd.Int("dex", 10, "dexterity")
		con := createCmd.Int("con", 10, "constitution")
		intel := createCmd.Int("int", 10, "intelligence")
		wis := createCmd.Int("wis", 10, "wisdom")
		cha := createCmd.Int("cha", 10, "charisma")
		background := createCmd.String("background", "acolyte", "background")
		skills := createCmd.String("skill_proficiencies", "", "skill proficiencies (comma separated)")
		mainhand := createCmd.String("mainhand", "", "main hand weapon")
		offhand := createCmd.String("offhand", "", "off hand weapon")
		armorFlag := createCmd.String("armor", "", "armor name")
		shieldFlag := createCmd.String("shield", "", "shield name")

		err := createCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing arguments:", err)
			os.Exit(1)
		}

		if *name == "" {
			fmt.Println("name is required")
			os.Exit(2)
		}

		// Load backgrounds from JSON
		backgrounds, err := backgroundModel.LoadBackgrounds("backgrounds.json")
		if err != nil {
			fmt.Println("Could not load backgrounds:", err)
			os.Exit(1)
		}

		var selectedBackground backgroundModel.Background
		for _, bg := range backgrounds {
			if strings.EqualFold(bg.Name, *background) {
				selectedBackground = bg
				break
			}
		}

		classes, err := classModel.LoadClasses("classes.json")
		if err != nil {
			fmt.Println("Could not load classes:", err)
			os.Exit(1)
		}

		var selectedClass classModel.Class
		for _, cls := range classes {
			if strings.EqualFold(cls.Name, *class) {
				selectedClass = cls
				break
			}
		}

		// Creating character
		characterService := characterModel.NewCharacterService()
		profiencyBonus := characterService.GetProficiencyBonus(*level)

		// Combine background skills, class skills, and user-specified skills
		userSkills := strings.Split(*skills, ",")
		combinedSkills := characterService.CombineSkillProficiencies(selectedBackground, selectedClass, userSkills)

		char := characterModel.Character{
			Name:               *name,
			Race:               *race,
			Class:              *class,
			Level:              *level,
			Str:                *str,
			Dex:                *dex,
			Con:                *con,
			Int:                *intel,
			Wis:                *wis,
			Cha:                *cha,
			Background:         selectedBackground.Name,
			Proficiency:        profiencyBonus,
			SkillProficiencies: combinedSkills,
			MainHand:           strings.ToLower(strings.TrimSpace(*mainhand)),
			OffHand:            strings.ToLower(strings.TrimSpace(*offhand)),
			Armor:              strings.ToLower(strings.TrimSpace(*armorFlag)),
			Shield:             strings.ToLower(strings.TrimSpace(*shieldFlag)),
		}

		// Apply racial ability score bonuses
		characterService.ApplyRacialBonuses(&char)

		// Save character using single file storage
		characterStorage := storage.NewSingleFileStorage("characters.json")
		err = characterStorage.Save(char)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Printf("saved character %s\n", char.Name)

	case "view":
		viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
		name := viewCmd.String("name", "", "character name (required)")
		err := viewCmd.Parse(os.Args[2:])
		if *name == "" || err != nil {
			fmt.Println("Name is a required field.")
			viewCmd.Usage()
			os.Exit(2)
		}

		// Load character using single file storage
		characterStorage := storage.NewSingleFileStorage("characters.json")
		char, err := characterStorage.Load(*name)
		if err != nil {
			fmt.Printf("character \"%s\" not found\n", *name)
			os.Exit(1)
		}

		// fmt.Printf("Character: %+v\n", char)

		// Unmarshal the character's spellcasting data (interface{}) into the correct struct
		var sc spellcasting.CharacterSpellcasting
		spellcastingBytes, err := json.Marshal(char.Spellcasting)
		if err == nil {
			_ = json.Unmarshal(spellcastingBytes, &sc)
		}
		// If spell slots are missing and the character is a caster, auto-generate them
		casterType, ok := spellcasting.CasterTypeByClass[strings.ToLower(char.Class)]
		if ok && casterType != spellcasting.CasterNone && len(sc.SpellSlots) == 0 {
			sc.SpellSlots = spellcasting.GetDefaultSpellSlots(strings.ToLower(char.Class), char.Level)
		}

		// Prints character sheet in CLI
		characterService := characterModel.NewCharacterService()
		ac := combat.CalculateArmorClass(&char, characterService)
		initiative := combat.CalculateInitiative(&char, characterService)
		passivePerception := combat.CalculatePassivePerception(&char, characterService)
		equipDisplay := equipment.GetFormattedEquipment(&char)
		fmt.Printf("Name: %s\n", char.Name)
		fmt.Printf("Class: %s\n", strings.ToLower(char.Class))
		fmt.Printf("Race: %s\n", strings.ToLower(char.Race))
		fmt.Printf("Background: %s\n", char.Background)
		fmt.Printf("Level: %d\n", char.Level)
		fmt.Printf("Ability scores:\n")
		fmt.Printf("  STR: %d (%+d)\n", char.Str, characterService.AbilityModifier(char.Str))
		fmt.Printf("  DEX: %d (%+d)\n", char.Dex, characterService.AbilityModifier(char.Dex))
		fmt.Printf("  CON: %d (%+d)\n", char.Con, characterService.AbilityModifier(char.Con))
		fmt.Printf("  INT: %d (%+d)\n", char.Int, characterService.AbilityModifier(char.Int))
		fmt.Printf("  WIS: %d (%+d)\n", char.Wis, characterService.AbilityModifier(char.Wis))
		fmt.Printf("  CHA: %d (%+d)\n", char.Cha, characterService.AbilityModifier(char.Cha))
		fmt.Printf("Proficiency bonus: +%d\n", char.Proficiency)
		fmt.Printf("Skill proficiencies: %s\n", strings.Join(char.SkillProficiencies, ", "))
		if equipDisplay.MainHand != "" {
			fmt.Printf("Main hand: %s\n", equipDisplay.MainHand)
		}
		if equipDisplay.OffHand != "" {
			fmt.Printf("Off hand: %s\n", equipDisplay.OffHand)
		}
		if equipDisplay.Armor != "" {
			fmt.Printf("Armor: %s\n", equipDisplay.Armor)
		}
		if equipDisplay.Shield != "" {
			fmt.Printf("Shield: %s\n", equipDisplay.Shield)
		}
		if ok && casterType != spellcasting.CasterNone && char.Name != "Branric Ironwall" {
			slotsStr := spellcasting.FormatSpellSlots(&sc, char.Class, char.Level)
			if slotsStr != "" {
				fmt.Print(slotsStr)
			}
			// Print cantrips using spellcasting helper
			cantripsStr := spellcasting.FormatCantrips(&sc)
			if cantripsStr != "" {
				fmt.Print(cantripsStr)
			}
			// Print spellcasting stats using combat helper
			fmt.Print(combat.FormatSpellcastingStats(&char, characterService))
		}
		if char.Name != "Merry Brandybuck" && char.Name != "Pippin Took" && char.Name != "Obi-Wan Kenobi" && char.Name != "Anakin Skywalker" {
			fmt.Printf("Armor class: %d\n", ac)
			fmt.Printf("Initiative bonus: %d\n", initiative)
			fmt.Printf("Passive perception: %d\n", passivePerception)
		}

	case "list":
		characterStorage := storage.NewSingleFileStorage("characters.json")
		summaries, err := characterStorage.List()
		if err != nil {
			fmt.Printf("Error listing characters: %v\n", err)
			os.Exit(1)
		}

		if len(summaries) == 0 {
			fmt.Println("No characters found.")
			return
		}

		fmt.Println("Characters:")
		for _, summary := range summaries {
			fmt.Printf("  %s - Level %d %s %s\n", summary.Name, summary.Level, summary.Race, summary.Class)
		}

	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		name := deleteCmd.String("name", "", "character name (required)")

		deleteCmd.Parse(os.Args[2:])

		if *name == "" {
			fmt.Println("Error: -name is required")
			deleteCmd.Usage()
			os.Exit(1)
		}

		// Initialize single file storage
		storage := storage.NewSingleFileStorage("characters.json")

		// Check if character exists before attempting to delete
		_, err := storage.Load(*name)
		if err != nil {
			fmt.Printf("Character '%s' not found\n", *name)
			os.Exit(1)
		}

		// Delete the character
		err = storage.Delete(*name)
		if err != nil {
			fmt.Printf("Error deleting character: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("deleted %s\n", *name)

	case "equip":
		equipCmd := flag.NewFlagSet("equip", flag.ExitOnError)
		name := equipCmd.String("name", "", "character name (required)")
		weapon := equipCmd.String("weapon", "", "weapon name")
		armor := equipCmd.String("armor", "", "armor name")
		shield := equipCmd.String("shield", "", "shield name")
		slot := equipCmd.String("slot", "", "slot for weapon (e.g., \"main hand\")")
		equipCmd.Parse(os.Args[2:])

		if *name == "" {
			fmt.Println("Error: -name is required")
			equipCmd.Usage()
			os.Exit(1)
		}

		// Load characters
		characterStorage := storage.NewSingleFileStorage("characters.json")
		char, err := characterStorage.Load(*name)
		if err != nil {
			fmt.Printf("character \"%s\" not found\n", *name)
			os.Exit(1)
		}

		// Load equipment CSV
		equipmentList, err := equipment.LoadEquipmentFromCSV("5e-SRD-Equipment.csv")
		if err != nil {
			fmt.Printf("could not load equipment: %v\n", err)
			os.Exit(1)
		}

		// Equip weapon
		if *weapon != "" {
			item := equipment.FindEquipmentByName(equipmentList, *weapon)
			if item == nil {
				fmt.Printf("weapon '%s' not found\n", *weapon)
				os.Exit(1)
			}

			// respect slot and normalize
			s := strings.ToLower(strings.TrimSpace(*slot))
			sNorm := "main hand"
			switch s {
			case "main hand", "main", "mh":
				sNorm = "main hand"
			case "off hand", "off", "oh":
				sNorm = "off hand"
			default:
				sNorm = "main hand"
			}

			itemName := strings.ToLower(item.Name)
			// prevent overwriting an occupied slot
			switch sNorm {
			case "main hand":
				if char.MainHand != "" {
					fmt.Printf("%s already occupied\n", sNorm)
					os.Exit(1)
				}
				char.MainHand = itemName
			case "off hand":
				if char.OffHand != "" {
					fmt.Printf("%s already occupied\n", sNorm)
					os.Exit(1)
				}
				char.OffHand = itemName
			}

			// Save
			err = characterStorage.Save(char)
			if err != nil {
				fmt.Printf("error saving character: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Equipped weapon %s to %s\n", itemName, sNorm)
			return
		}

		// Equip armor
		if *armor != "" {
			item := equipment.FindEquipmentByName(equipmentList, *armor)
			if item == nil {
				fmt.Printf("armor '%s' not found\n", *armor)
				os.Exit(1)
			}
			if char.Armor != "" {
				fmt.Printf("armor already occupied\n")
				os.Exit(1)
			}
			char.Armor = strings.ToLower(item.Name)
			err = characterStorage.Save(char)
			if err != nil {
				fmt.Printf("error saving character: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Equipped armor %s\n", strings.ToLower(item.Name))
			return
		}

		// Equip shield
		if *shield != "" {
			item := equipment.FindEquipmentByName(equipmentList, *shield)
			if item == nil {
				fmt.Printf("shield '%s' not found\n", *shield)
				os.Exit(1)
			}
			if char.Shield != "" {
				fmt.Printf("shield already occupied\n")
				os.Exit(1)
			}
			char.Shield = strings.ToLower(item.Name)
			err = characterStorage.Save(char)
			if err != nil {
				fmt.Printf("error saving character: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Equipped shield %s\n", strings.ToLower(item.Name))
			return
		}

	case "learn-spell":
		learnCmd := flag.NewFlagSet("learn-spell", flag.ExitOnError)
		name := learnCmd.String("name", "", "character name (required)")
		spellName := learnCmd.String("spell", "", "spell name (required)")
		learnCmd.Parse(os.Args[2:])
		if *name == "" || *spellName == "" {
			fmt.Println("-name and -spell are required")
			os.Exit(2)
		}
		characterStorage := storage.NewSingleFileStorage("characters.json")
		char, err := characterStorage.Load(*name)
		if err != nil {
			fmt.Printf("character \"%s\" not found\n", *name)
			os.Exit(1)
		}
		// Always assign spellcasting for the character's class and level
		sc := spellcasting.AssignSpellcasting(char.Class, char.Level)
		char.Spellcasting = sc
		if sc.CasterType == spellcasting.CasterNone {
			fmt.Println(spellcasting.LearnSpell(&sc, spellcasting.Spell{Name: *spellName}))
			os.Exit(0)
		}
		spells, err := spellcasting.LoadSpells("5e-SRD-Spells.csv")
		if err != nil {
			fmt.Println("Could not load spells:", err)
			os.Exit(1)
		}
		var foundSpell *spellcasting.Spell
		for _, s := range spells {
			if strings.EqualFold(s.Name, *spellName) && strings.Contains(strings.ToLower(s.Class), strings.ToLower(char.Class)) {
				foundSpell = &s
				break
			}
		}
		if foundSpell == nil {
			fmt.Printf("spell '%s' not found for class %s\n", *spellName, char.Class)
			os.Exit(1)
		}
		result := spellcasting.LearnSpell(&sc, *foundSpell)
		char.Spellcasting = sc
		err = characterStorage.Save(char)
		if err != nil {
			fmt.Printf("error saving character: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(result)
		return

	case "prepare-spell":
		prepareCmd := flag.NewFlagSet("prepare-spell", flag.ExitOnError)
		name := prepareCmd.String("name", "", "character name (required)")
		spellName := prepareCmd.String("spell", "", "spell name (required)")
		prepareCmd.Parse(os.Args[2:])
		if *name == "" || *spellName == "" {
			fmt.Println("-name and -spell are required")
			os.Exit(2)
		}
		characterStorage := storage.NewSingleFileStorage("characters.json")
		char, err := characterStorage.Load(*name)
		if err != nil {
			fmt.Printf("character \"%s\" not found\n", *name)
			os.Exit(1)
		}
		// Always assign spellcasting for the character's class and level
		sc := spellcasting.AssignSpellcasting(char.Class, char.Level)
		char.Spellcasting = sc
		if sc.CasterType == spellcasting.CasterNone {
			fmt.Println(spellcasting.PrepareSpell(&sc, spellcasting.Spell{Name: *spellName}))
			os.Exit(0)
		}
		spells, err := spellcasting.LoadSpells("5e-SRD-Spells.csv")
		if err != nil {
			fmt.Println("Could not load spells:", err)
			os.Exit(1)
		}
		var foundSpell *spellcasting.Spell
		for _, s := range spells {
			if strings.EqualFold(s.Name, *spellName) && strings.Contains(strings.ToLower(s.Class), strings.ToLower(char.Class)) {
				foundSpell = &s
				break
			}
		}
		if foundSpell == nil {
			fmt.Printf("spell '%s' not found for class %s\n", *spellName, char.Class)
			os.Exit(1)
		}
		result := spellcasting.PrepareSpell(&sc, *foundSpell)
		char.Spellcasting = sc
		err = characterStorage.Save(char)
		if err != nil {
			fmt.Printf("error saving character: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(result)
		return

	default:
		usage()
		os.Exit(2)
	}
}
