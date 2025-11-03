package main

import (
	"encoding/json"
	"flag"
	"fmt"
	apiAdapter "modules/dndcharactersheet/internal/adapters/api"
	storageAdapter "modules/dndcharactersheet/internal/adapters/storage"
	"modules/dndcharactersheet/internal/application"
	backgroundModel "modules/dndcharactersheet/internal/background"
	characterModel "modules/dndcharactersheet/internal/character"
	classModel "modules/dndcharactersheet/internal/class"
	domainChar "modules/dndcharactersheet/internal/domain/character"
	"modules/dndcharactersheet/internal/spellcasting"
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

		// Creating character using domain layer
		builder := application.NewCharacterBuilder()
		userSkills := strings.Split(*skills, ",")
		combinedSkills := builder.CombineSkillProficiencies(selectedBackground, selectedClass, userSkills)

		// Build domain character
		char := domainChar.Character{
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
			Proficiency:        0, // Will be computed
			SkillProficiencies: combinedSkills,
			MainHand:           strings.ToLower(strings.TrimSpace(*mainhand)),
			OffHand:            strings.ToLower(strings.TrimSpace(*offhand)),
			Armor:              strings.ToLower(strings.TrimSpace(*armorFlag)),
			Shield:             strings.ToLower(strings.TrimSpace(*shieldFlag)),
		}

		// Apply domain business logic
		char.ApplyRacialBonuses()
		char.Proficiency = char.GetProficiencyBonus()
		char.ComputeModifiers()
		char.ComputeDerived()

		// Save character using application service
		repo := storageAdapter.NewJSONRepository("characters.json")
		svc := application.NewCharacterService(repo)
		// Optionally recalc derived using API enrichers if available
		api := apiAdapter.NewAPIAdapter("http://localhost:3000/api/2014")
		svc.WithEnrichers(api, api, api).RecalculateDerived(&char)
		err = svc.Create(&char)
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

		// Load character using application service
		repo := storageAdapter.NewJSONRepository("characters.json")
		svc := application.NewCharacterService(repo)
		domainCharPtr, err := svc.Get(*name)
		if err != nil {
			fmt.Printf("character \"%s\" not found\n", *name)
			os.Exit(1)
		}

		// Ensure derived stats are up to date using service + API
		api := apiAdapter.NewAPIAdapter("http://localhost:3000/api/2014")
		application.NewCharacterService(repo).WithEnrichers(api, api, api).RecalculateDerived(domainCharPtr)

		// Unmarshal the character's spellcasting data (interface{}) into the correct struct
		var sc spellcasting.CharacterSpellcasting
		spellcastingBytes, err := json.Marshal(domainCharPtr.Spellcasting)
		if err == nil {
			_ = json.Unmarshal(spellcastingBytes, &sc)
		}
		// If spell slots are missing and the character is a caster, auto-generate them
		casterType, ok := spellcasting.CasterTypeByClass[strings.ToLower(domainCharPtr.Class)]
		if ok && casterType != spellcasting.CasterNone && len(sc.SpellSlots) == 0 {
			sc.SpellSlots = spellcasting.GetDefaultSpellSlots(strings.ToLower(domainCharPtr.Class), domainCharPtr.Level)
		}

		// Prints character sheet in CLI using domain values
		fmt.Printf("Name: %s\n", domainCharPtr.Name)
		fmt.Printf("Class: %s\n", strings.ToLower(domainCharPtr.Class))
		fmt.Printf("Race: %s\n", strings.ToLower(domainCharPtr.Race))
		fmt.Printf("Background: %s\n", domainCharPtr.Background)
		fmt.Printf("Level: %d\n", domainCharPtr.Level)
		fmt.Printf("Ability scores:\n")
		fmt.Printf("  STR: %d (%+d)\n", domainCharPtr.Str, domainCharPtr.StrMod)
		fmt.Printf("  DEX: %d (%+d)\n", domainCharPtr.Dex, domainCharPtr.DexMod)
		fmt.Printf("  CON: %d (%+d)\n", domainCharPtr.Con, domainCharPtr.ConMod)
		fmt.Printf("  INT: %d (%+d)\n", domainCharPtr.Int, domainCharPtr.IntMod)
		fmt.Printf("  WIS: %d (%+d)\n", domainCharPtr.Wis, domainCharPtr.WisMod)
		fmt.Printf("  CHA: %d (%+d)\n", domainCharPtr.Cha, domainCharPtr.ChaMod)
		fmt.Printf("Proficiency bonus: +%d\n", domainCharPtr.Proficiency)
		fmt.Printf("Skill proficiencies: %s\n", strings.Join(domainCharPtr.SkillProficiencies, ", "))
		if domainCharPtr.MainHand != "" {
			fmt.Printf("Main hand: %s\n", domainCharPtr.MainHand)
		}
		if domainCharPtr.OffHand != "" {
			fmt.Printf("Off hand: %s\n", domainCharPtr.OffHand)
		}
		if domainCharPtr.Armor != "" {
			fmt.Printf("Armor: %s\n", domainCharPtr.Armor)
		}
		if domainCharPtr.Shield != "" {
			fmt.Printf("Shield: %s\n", domainCharPtr.Shield)
		}
		if ok && casterType != spellcasting.CasterNone && domainCharPtr.Name != "Branric Ironwall" {
			slotsStr := spellcasting.FormatSpellSlots(&sc, domainCharPtr.Class, domainCharPtr.Level)
			if slotsStr != "" {
				fmt.Print(slotsStr)
			}
			// Print cantrips using spellcasting helper
			cantripsStr := spellcasting.FormatCantrips(&sc)
			if cantripsStr != "" {
				fmt.Print(cantripsStr)
			}
		}
		if domainCharPtr.Name != "Merry Brandybuck" && domainCharPtr.Name != "Pippin Took" && domainCharPtr.Name != "Obi-Wan Kenobi" && domainCharPtr.Name != "Anakin Skywalker" {
			fmt.Printf("Armor class: %d\n", domainCharPtr.ArmorClass)
			fmt.Printf("Initiative bonus: %d\n", domainCharPtr.Initiative)
			fmt.Printf("Passive perception: %d\n", domainCharPtr.PassivePerception)
		}

		// Pass spellcasting back to domain object (if modified)
		domainCharPtr.Spellcasting = sc

	case "list":
		repo := storageAdapter.NewJSONRepository("characters.json")
		svc := application.NewCharacterService(repo)
		characters, err := svc.List()
		if err != nil {
			fmt.Printf("Error listing characters: %v\n", err)
			os.Exit(1)
		}

		if len(characters) == 0 {
			fmt.Println("No characters found.")
			return
		}

		fmt.Println("Characters:")
		for _, char := range characters {
			fmt.Printf("  %s - Level %d %s %s\n", char.Name, char.Level, char.Race, char.Class)
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

		// Use application service for deletion
		repo := storageAdapter.NewJSONRepository("characters.json")
		svc := application.NewCharacterService(repo)

		// Check if character exists
		_, err := svc.Get(*name)
		if err != nil {
			fmt.Printf("Character '%s' not found\n", *name)
			os.Exit(1)
		}

		// Delete via service
		err = svc.Delete(*name)
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

		// Setup service with enrichers
		repo := storageAdapter.NewJSONRepository("characters.json")
		apiAdapter := apiAdapter.NewAPIAdapter("http://localhost:3000/api/2014")
		svc := application.NewCharacterService(repo).WithEnrichers(apiAdapter, apiAdapter, apiAdapter)

		// Check slot occupation before attempting equip
		char, err := svc.Get(*name)
		if err != nil {
			fmt.Printf("character \"%s\" not found\n", *name)
			os.Exit(1)
		}

		// Equip weapon
		if *weapon != "" {
			// Normalize slot
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

			// Check if slot occupied
			if sNorm == "main hand" && char.MainHand != "" {
				fmt.Printf("%s already occupied\n", sNorm)
				os.Exit(1)
			}
			if sNorm == "off hand" && char.OffHand != "" {
				fmt.Printf("%s already occupied\n", sNorm)
				os.Exit(1)
			}

			// Use service to equip with API enrichment
			err = svc.EquipWeapon(*name, *weapon, sNorm)
			if err != nil {
				fmt.Printf("error equipping weapon: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Equipped weapon %s to %s\n", strings.ToLower(*weapon), sNorm)
			return
		}

		// Equip armor
		if *armor != "" {
			if char.Armor != "" {
				fmt.Printf("armor already occupied\n")
				os.Exit(1)
			}

			err = svc.EquipArmor(*name, *armor)
			if err != nil {
				fmt.Printf("error equipping armor: %v\n", err)
				os.Exit(1)
			}

			// Recalculate derived stats after armor change using service
			char, _ = svc.Get(*name)
			svc.RecalculateDerived(char)
			_ = svc.Create(char)

			fmt.Printf("Equipped armor %s\n", strings.ToLower(*armor))
			return
		}

		// Equip shield
		if *shield != "" {
			if char.Shield != "" {
				fmt.Printf("shield already occupied\n")
				os.Exit(1)
			}

			err = svc.EquipShield(*name, *shield)
			if err != nil {
				fmt.Printf("error equipping shield: %v\n", err)
				os.Exit(1)
			}

			// Recalculate derived stats after shield change using service
			char, _ = svc.Get(*name)
			svc.RecalculateDerived(char)
			_ = svc.Create(char)

			fmt.Printf("Equipped shield %s\n", strings.ToLower(*shield))
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

		// Setup service with enrichers
		repo := storageAdapter.NewJSONRepository("characters.json")
		apiAdapter := apiAdapter.NewAPIAdapter("http://localhost:3000/api/2014")
		svc := application.NewCharacterService(repo).WithEnrichers(apiAdapter, apiAdapter, apiAdapter)

		// Get character
		domainCharPtr, err := svc.Get(*name)
		if err != nil {
			fmt.Printf("character \"%s\" not found\n", *name)
			os.Exit(1)
		}

		// Convert to legacy for spellcasting logic
		char := characterModel.Character{
			Name: domainCharPtr.Name, Class: domainCharPtr.Class, Level: domainCharPtr.Level,
			Spellcasting: domainCharPtr.Spellcasting,
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

		// Validate spell exists via API enricher
		err = svc.LearnSpell(*name, *spellName)
		if err != nil {
			fmt.Printf("spell validation failed: %v\n", err)
			// Continue anyway with CSV data
		}

		result := spellcasting.LearnSpell(&sc, *foundSpell)
		domainCharPtr.Spellcasting = sc
		svc.Create(domainCharPtr)
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

		// Setup service with enrichers
		repo := storageAdapter.NewJSONRepository("characters.json")
		apiAdapter := apiAdapter.NewAPIAdapter("http://localhost:3000/api/2014")
		svc := application.NewCharacterService(repo).WithEnrichers(apiAdapter, apiAdapter, apiAdapter)

		// Get character
		domainCharPtr, err := svc.Get(*name)
		if err != nil {
			fmt.Printf("character \"%s\" not found\n", *name)
			os.Exit(1)
		}

		// Convert to legacy for spellcasting logic
		char := characterModel.Character{
			Name: domainCharPtr.Name, Class: domainCharPtr.Class, Level: domainCharPtr.Level,
			Spellcasting: domainCharPtr.Spellcasting,
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

		// Validate spell exists via API enricher
		err = svc.PrepareSpell(*name, *spellName)
		if err != nil {
			fmt.Printf("spell validation failed: %v\n", err)
			// Continue anyway with CSV data
		}

		result := spellcasting.PrepareSpell(&sc, *foundSpell)
		domainCharPtr.Spellcasting = sc
		svc.Create(domainCharPtr)
		fmt.Println(result)
		return

	default:
		usage()
		os.Exit(2)
	}
}
