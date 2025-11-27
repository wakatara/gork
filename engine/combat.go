package engine

import (
	"fmt"
	"strings"
)

// Combat outcome constants (from ZIL 1actions.zil:3245-3254)
const (
	CombatMissed = iota + 1
	CombatUnconscious
	CombatKilled
	CombatLightWound
	CombatSeriousWound
	CombatStagger
	CombatLoseWeapon
	CombatHesitate    // For unconscious victims (not used in basic implementation)
	CombatSittingDuck // For unconscious victims (not used in basic implementation)
)

// Combat constants (from ZIL 1actions.zil:3323-3325)
const (
	StrengthMin = 2   // Minimum player strength
	StrengthMax = 7   // Maximum player strength
	ScoreMax    = 350 // Maximum possible score in Zork
)

// Combat result table - 9 entries indexed 0-8 (ZIL uses RANDOM 9 = 1-9)
type CombatResultTable [9]int

// Defense level 1 (strength 1) - weak defense
// ZIL DEF1 has 13 entries, we compress to 9 while maintaining probabilities
// Original: 4 MISSED, 2 STAGGER, 2 UNCONSCIOUS, 5 KILLED (out of 13)
var def1 = CombatResultTable{
	CombatMissed, CombatMissed, CombatMissed,           // 3/9 ≈ 33% (was 31%)
	CombatStagger, CombatStagger,                       // 2/9 ≈ 22% (was 15%)
	CombatUnconscious,                                  // 1/9 ≈ 11% (was 15%)
	CombatKilled, CombatKilled, CombatKilled,           // 3/9 ≈ 33% (was 38%)
}

// Defense level 2A (strength 2, ATT <= 2)
// Original: 5 MISSED, 2 STAGGER, 2 LIGHT-WOUND, 1 UNCONSCIOUS (out of 10)
var def2A = CombatResultTable{
	CombatMissed, CombatMissed, CombatMissed, CombatMissed, // 4/9 ≈ 44% (was 50%)
	CombatStagger, CombatStagger,                            // 2/9 ≈ 22% (was 20%)
	CombatLightWound, CombatLightWound,                      // 2/9 ≈ 22% (was 20%)
	CombatUnconscious,                                       // 1/9 ≈ 11% (was 10%)
}

// Defense level 2B (strength 2, ATT > 2)
// Original: 3 MISSED, 2 STAGGER, 3 LIGHT-WOUND, 1 UNCONSCIOUS, 3 KILLED (out of 12, use 11)
var def2B = CombatResultTable{
	CombatMissed, CombatMissed,                 // 2/9 ≈ 22% (was 27%)
	CombatStagger, CombatStagger,               // 2/9 ≈ 22% (was 18%)
	CombatLightWound, CombatLightWound, CombatLightWound, // 3/9 ≈ 33% (was 27%)
	CombatUnconscious,                          // 1/9 ≈ 11% (was 9%)
	CombatKilled,                               // 1/9 ≈ 11% (was 27%, reduced for balance)
}

// Defense level 3A (strength 3+, ATT-DEF <= -2)
// Original: 5 MISSED, 2 STAGGER, 2 LIGHT-WOUND, 2 SERIOUS-WOUND (out of 11)
var def3A = CombatResultTable{
	CombatMissed, CombatMissed, CombatMissed, CombatMissed, // 4/9 ≈ 44% (was 45%)
	CombatStagger, CombatStagger,                            // 2/9 ≈ 22% (was 18%)
	CombatLightWound, CombatLightWound,                      // 2/9 ≈ 22% (was 18%)
	CombatSeriousWound,                                      // 1/9 ≈ 11% (was 18%)
}

// Defense level 3B (strength 3+, -1 <= ATT-DEF <= 0)
// Original: 3 MISSED, 2 STAGGER, 3 LIGHT-WOUND, 3 SERIOUS-WOUND (out of 11)
var def3B = CombatResultTable{
	CombatMissed, CombatMissed,                              // 2/9 ≈ 22% (was 27%)
	CombatStagger, CombatStagger,                            // 2/9 ≈ 22% (was 18%)
	CombatLightWound, CombatLightWound, CombatLightWound,    // 3/9 ≈ 33% (was 27%)
	CombatSeriousWound, CombatSeriousWound,                  // 2/9 ≈ 22% (was 27%)
}

// Defense level 3C (strength 3+, ATT-DEF >= 1)
// Original: 1 MISSED, 2 STAGGER, 4 LIGHT-WOUND, 3 SERIOUS-WOUND (out of 10, use 9)
var def3C = CombatResultTable{
	CombatMissed,                                            // 1/9 ≈ 11%
	CombatStagger, CombatStagger,                            // 2/9 ≈ 22%
	CombatLightWound, CombatLightWound, CombatLightWound,    // 3/9 ≈ 33% (was 44%)
	CombatSeriousWound, CombatSeriousWound, CombatSeriousWound, // 3/9 ≈ 33%
}

// VillainData contains NPC-specific combat information
// From ZIL VILLAINS table (1actions.zil:3801-3804)
type VillainData struct {
	BestWeapon    string              // Weapon that weakens this NPC
	BestAdvantage int                 // Strength reduction when player has best weapon
	MeleeMessages map[int][]string    // Messages indexed by outcome type
}

var villainData = map[string]VillainData{
	"troll": {
		BestWeapon:    "sword",
		BestAdvantage: 1,
		MeleeMessages: trollMelee,
	},
	"thief": {
		BestWeapon:    "knife",
		BestAdvantage: 1,
		MeleeMessages: thiefMelee,
	},
	"cyclops": {
		BestWeapon:    "", // No weakness
		BestAdvantage: 0,
		MeleeMessages: cyclopsMelee,
	},
}

// fightStrength calculates player's combat strength
// From ZIL FIGHT-STRENGTH routine (1actions.zil:3375-3381)
func (g *GameV2) fightStrength() int {
	// Base strength from score
	// ZIL: STRENGTH-MIN + (SCORE / (SCORE-MAX / (STRENGTH-MAX - STRENGTH-MIN)))
	baseStrength := StrengthMin + (g.Score / (ScoreMax / (StrengthMax - StrengthMin)))

	// Add P?STRENGTH modifier (normally 0, reduced by wounds)
	finalStrength := baseStrength + g.Player.StrengthModifier

	if finalStrength < 1 {
		return 1
	}
	return finalStrength
}

// villainStrength calculates NPC's combat strength
// From ZIL VILLAIN-STRENGTH routine (1actions.zil:3383-3397)
func (g *GameV2) villainStrength(npc *NPC) int {
	strength := npc.Strength

	// If NPC is distracted (thief-engrossed), cap strength at 2
	if npc.ID == "thief" && g.Flags["thief-engrossed"] {
		if strength > 2 {
			strength = 2
		}
		g.Flags["thief-engrossed"] = false
	}

	// If player has NPC's best weapon, reduce NPC strength
	if vData, ok := villainData[npc.ID]; ok && vData.BestWeapon != "" {
		for _, itemID := range g.Player.Inventory {
			if itemID == vData.BestWeapon {
				strength -= vData.BestAdvantage
				if strength < 1 {
					strength = 1
				}
				break
			}
		}
	}

	return strength
}

// selectResultTable chooses the appropriate combat result table
// Based on ZIL logic in HERO-BLOW (1actions.zil:3428-3439)
func selectResultTable(att, def int) CombatResultTable {
	if def == 1 {
		return def1
	}

	if def == 2 {
		if att <= 2 {
			return def2A
		}
		return def2B
	}

	// def >= 3
	diff := att - def
	if diff <= -2 {
		return def3A
	} else if diff >= 1 {
		return def3C
	}
	return def3B
}

// heroBlow performs player attack on NPC
// From ZIL HERO-BLOW routine (1actions.zil:3476-3560)
func (g *GameV2) heroBlow(npc *NPC, weapon *Item) (int, string) {
	att := g.fightStrength()
	def := g.villainStrength(npc)

	// Select result table
	table := selectResultTable(att, def)

	// Roll outcome (RANDOM 9 = 1-9, index 0-8)
	outcomeIndex := g.rand.Intn(9)
	outcome := table[outcomeIndex]

	// 25% chance to convert STAGGER to LOSE-WEAPON (ZIL line 3444-3447)
	if outcome == CombatStagger && g.rand.Intn(100) < 25 {
		// Check if NPC has a weapon to lose
		npcWeapon := g.findNPCWeapon(npc)
		if npcWeapon != nil {
			outcome = CombatLoseWeapon
		}
	}

	// Select random message for this outcome
	messages := heroMelee[outcome]
	if len(messages) == 0 {
		return outcome, fmt.Sprintf("You attack the %s!", npc.Name)
	}

	message := messages[g.rand.Intn(len(messages))]

	// Replace placeholders
	message = strings.ReplaceAll(message, "{weapon}", weapon.Name)
	message = strings.ReplaceAll(message, "{npc}", npc.Name)

	return outcome, message
}

// applyHeroOutcome applies combat outcome to NPC
// From ZIL HERO-BLOW outcome handling (1actions.zil:3452-3474)
func (g *GameV2) applyHeroOutcome(npc *NPC, outcome int) {
	switch outcome {
	case CombatKilled:
		npc.Strength = 0
		npc.Flags.IsAlive = false

	case CombatUnconscious:
		// In full ZIL, NPCs can wake up from unconsciousness
		// For simplicity, treat as death
		npc.Strength = 0
		npc.Flags.IsAlive = false

	case CombatLightWound:
		npc.Strength -= 1
		if npc.Strength <= 0 {
			npc.Strength = 0
			npc.Flags.IsAlive = false
		}

	case CombatSeriousWound:
		npc.Strength -= 2
		if npc.Strength <= 0 {
			npc.Strength = 0
			npc.Flags.IsAlive = false
		}

	case CombatStagger:
		// NPC skips next turn
		g.Flags[npc.ID+"-staggered"] = true

	case CombatLoseWeapon:
		// NPC drops weapon
		weapon := g.findNPCWeapon(npc)
		if weapon != nil {
			weapon.Location = npc.Location
			if room := g.Rooms[npc.Location]; room != nil {
				room.AddItem(weapon.ID)
			}
			npc.Inventory = removeFromSlice(npc.Inventory, weapon.ID)
		}
	}
}

// villainBlow performs NPC counter-attack on player
// From ZIL VILLAIN-BLOW routine (1actions.zil:3413-3474)
func (g *GameV2) villainBlow(npc *NPC) (int, string) {
	// If NPC is staggered, skip turn (ZIL lines 3418-3422)
	if g.Flags[npc.ID+"-staggered"] {
		g.Flags[npc.ID+"-staggered"] = false
		return CombatMissed, fmt.Sprintf("The %s slowly regains his feet.", npc.Name)
	}

	att := g.villainStrength(npc)
	def := g.fightStrength()

	// Select result table
	table := selectResultTable(att, def)

	// Roll outcome
	outcomeIndex := g.rand.Intn(9)
	outcome := table[outcomeIndex]

	// 25% chance to convert STAGGER to LOSE-WEAPON
	if outcome == CombatStagger && g.rand.Intn(100) < 25 {
		playerWeapon := g.findPlayerWeapon()
		if playerWeapon != nil {
			outcome = CombatLoseWeapon
		}
	}

	// Get NPC-specific messages
	vData, ok := villainData[npc.ID]
	if !ok {
		return outcome, fmt.Sprintf("The %s attacks!", npc.Name)
	}

	messages := vData.MeleeMessages[outcome]
	if len(messages) == 0 {
		return outcome, fmt.Sprintf("The %s attacks!", npc.Name)
	}

	message := messages[g.rand.Intn(len(messages))]

	// Replace placeholders (note: {weapon} refers to PLAYER's weapon)
	playerWeapon := g.findPlayerWeapon()
	if playerWeapon != nil {
		message = strings.ReplaceAll(message, "{weapon}", playerWeapon.Name)
	}

	return outcome, message
}

// applyVillainOutcome applies combat outcome to player
// From ZIL VILLAIN-BLOW outcome handling (1actions.zil:3452-3474)
func (g *GameV2) applyVillainOutcome(outcome int) {
	switch outcome {
	case CombatKilled:
		g.Player.Health = 0

	case CombatUnconscious:
		g.Player.Health = 0 // Simplified: treat as death

	case CombatLightWound:
		g.Player.StrengthModifier -= 1

	case CombatSeriousWound:
		g.Player.StrengthModifier -= 2

	case CombatStagger:
		// Player is staggered (in ZIL, misses next turn)
		// We'll just show it in the message, no mechanical effect needed
		// since this is single combat, not turn-based with multiple NPCs

	case CombatLoseWeapon:
		// Player drops weapon
		playerWeapon := g.findPlayerWeapon()
		if playerWeapon != nil {
			playerWeapon.Location = g.Location
			if room := g.Rooms[g.Location]; room != nil {
				room.AddItem(playerWeapon.ID)
			}
			g.Player.Inventory = removeFromSlice(g.Player.Inventory, playerWeapon.ID)
		}
	}
}

// findPlayerWeapon finds the first weapon in player's inventory
func (g *GameV2) findPlayerWeapon() *Item {
	for _, itemID := range g.Player.Inventory {
		item := g.Items[itemID]
		if item != nil && item.Flags.IsWeapon {
			return item
		}
	}
	return nil
}

// findNPCWeapon finds the first weapon in NPC's inventory
func (g *GameV2) findNPCWeapon(npc *NPC) *Item {
	for _, itemID := range npc.Inventory {
		item := g.Items[itemID]
		if item != nil && item.Flags.IsWeapon {
			return item
		}
	}
	return nil
}

// Helper function to remove item from slice
func removeFromSlice(slice []string, item string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}
