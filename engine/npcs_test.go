package engine

import (
	"strings"
	"testing"
)

func TestAllNPCsExist(t *testing.T) {
	g := NewGameV2()

	expectedNPCs := []string{"troll", "thief", "cyclops", "bat", "ghosts"}

	for _, npcID := range expectedNPCs {
		npc := g.NPCs[npcID]
		if npc == nil {
			t.Errorf("NPC %s does not exist", npcID)
			continue
		}

		if !npc.Flags.IsAlive {
			t.Errorf("NPC %s should start alive", npcID)
		}
	}

	t.Logf("All 5 NPCs created successfully")
}

func TestTrollCombat(t *testing.T) {
	g := NewGameV2()
	g.Location = "troll-room"

	// Give player a sword
	sword := g.Items["sword"]
	sword.Location = "inventory"
	g.Player.Inventory = append(g.Player.Inventory, "sword")

	// Attack troll
	result := g.Process("attack troll")

	if !strings.Contains(result, "attack") {
		t.Errorf("Expected combat message, got: %s", result)
	}
}

func TestTrollFood(t *testing.T) {
	g := NewGameV2()
	g.Location = "troll-room"

	// Give player lunch (non-treasure)
	lunch := g.Items["lunch"]
	lunch.Location = "inventory"
	g.Player.Inventory = append(g.Player.Inventory, "lunch")

	// Give lunch to troll - he eats it but doesn't leave (only treasures satisfy him)
	result := g.Process("give lunch to troll")

	if !strings.Contains(result, "gleefully eats it") {
		t.Errorf("Expected troll to eat lunch, got: %s", result)
	}
	if !strings.Contains(result, "still blocking") {
		t.Errorf("Expected troll to still block passages, got: %s", result)
	}

	// Verify troll is still there (not satisfied by non-treasures)
	troll := g.NPCs["troll"]
	if troll == nil || troll.Location == "" {
		t.Error("Troll should still be present (non-treasure doesn't satisfy him)")
	}
}

func TestCyclopsCombat(t *testing.T) {
	g := NewGameV2()
	g.Location = "cyclops-room"

	// Verify cyclops exists
	cyclops := g.NPCs["cyclops"]
	if cyclops == nil {
		t.Fatal("Cyclops does not exist")
	}

	if !cyclops.Flags.IsAggressive {
		t.Error("Cyclops should be aggressive")
	}
}

func TestThiefStealing(t *testing.T) {
	g := NewGameV2()
	g.Location = "maze-1" // Where thief starts

	// Give player a treasure
	diamond := g.Items["diamond"]
	diamond.Location = "inventory"
	g.Player.Inventory = append(g.Player.Inventory, "diamond")

	// Give diamond to thief
	result := g.Process("give diamond to thief")

	if !strings.Contains(result, "snatches") || !strings.Contains(result, "runs off") {
		t.Errorf("Expected thief to steal, got: %s", result)
	}

	// Verify thief has the diamond
	thief := g.NPCs["thief"]
	found := false
	for _, itemID := range thief.Inventory {
		if itemID == "diamond" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Thief should have stolen the diamond")
	}
}

func TestBatExists(t *testing.T) {
	g := NewGameV2()

	bat := g.NPCs["bat"]
	if bat == nil {
		t.Fatal("Bat does not exist")
	}

	if bat.Flags.CanFight {
		t.Error("Bat should not be fightable")
	}

	if !bat.Flags.IsAggressive {
		t.Error("Bat should be aggressive (will grab you)")
	}
}

func TestGhostsExist(t *testing.T) {
	g := NewGameV2()

	ghosts := g.NPCs["ghosts"]
	if ghosts == nil {
		t.Fatal("Ghosts do not exist")
	}

	if ghosts.Flags.CanFight {
		t.Error("Ghosts should not be fightable (need exorcism)")
	}
}

func TestNPCLocations(t *testing.T) {
	g := NewGameV2()

	tests := []struct {
		npcID    string
		location string
	}{
		{"troll", "troll-room"},
		{"thief", "maze-1"},
		{"cyclops", "cyclops-room"},
		{"bat", "bat-room"},
		{"ghosts", "entrance-to-hades"},
	}

	for _, tt := range tests {
		npc := g.NPCs[tt.npcID]
		if npc == nil {
			t.Errorf("NPC %s does not exist", tt.npcID)
			continue
		}

		if npc.Location != tt.location {
			t.Errorf("NPC %s should be in %s, got %s", tt.npcID, tt.location, npc.Location)
		}

		// Verify NPC is in the room's NPC list
		room := g.Rooms[tt.location]
		if room == nil {
			t.Errorf("Room %s does not exist", tt.location)
			continue
		}

		found := false
		for _, id := range room.NPCs {
			if id == tt.npcID {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("NPC %s not in room %s's NPC list", tt.npcID, tt.location)
		}
	}
}
