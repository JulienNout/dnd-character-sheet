1. [x] Character CRU
   - [x] Create new characters
   - [x] Read their character sheet
   - [x] Update relevant bits: level and proficiency bonus
   - [x] Characters have the following attributes: name, race, class, level, ability scores, background and proficiency bonus
   - [x] Ability scores are determined through Standard Array (15/14/13/12/10/8) and Race
   - [x] Classes only deal with main classes, we'll ignore subclasses
   - [x] Add skill proficiencies on character creation. We'll ignore racial skill proficiencies.
   - [x] Determine skill modifiers based on ability scores and proficiency bonus
   - [x] We'll ignore level-up bonuses and unlocks other than proficiency bonus and spell slots.
   - [x] Anything else listed in the SRD we'll ignore

2. [x] Equipment management
   - [x] Add equipment to/from a character: weapons, armor and shield
   - [x] We'll ignore the concept of "inventory": unequipped items such as a backup dagger, rope, torches, etc.
   - [x] Any other gear listed in the SRD we'll ignore

3. [x] Spellcasting
   - [x] When applicable :)
   - [x] Add known/prepared spells from spell list for your class. It's okay to pick random spells as long as they fit the slot level
   - [x] Add max spell slots per level
   - [x] Anything else related to spells listed in the SRD we'll ignore

4. [x] Integrate external information → Manually tested by Loek once all automatic tests pass
   - [ ] Enrich spells and items from external API: https://www.dnd5eapi.co/. Use Go's built-in concurrency to send multiple requests at the same time.
     > PLEASE READ CAREFULLY: this is an API built and maintained by volunteers. They've rate-limited the API to 50 requests per second. Theoretically, you could gather all info on spells and equipment in about 7.5 seconds. But be nice to these people and don't run up their server costs. Send 5–10 requests per second at most, and only test with small batches (8 or 10 spells) instead of your whole database at a time. If it works for 10 spells, it will work for all 319. You can also download their Docker image to test with their API locally: https://github.com/5e-bits/5e-srd-api?tab=readme-ov-file#how-to-run
   - [x] Spells: school, range
   - [x] Weapons: category, range (normal range, ignore long range), two-handed
   - [x] Armor: armor class, dexterity bonus
   - [x] If you want to do more and build a fully fledged character sheet generator, feel free! We will only look at these properties for this course and exam.

5. [ ] Combat stat calculation
   - [ ] Armor class, initiative, passive perception (character sheet calls it "passive wisdom")
   - [ ] Spell casting: spellcasting ability, spell save DC, spell attack bonus

6. [ ] User interfaces
   - [ ] CLI: create a character, view it and list all characters
   - [ ] HTML frontend: a list of all characters that links to their character sheets (HTML template provided on Learn)
