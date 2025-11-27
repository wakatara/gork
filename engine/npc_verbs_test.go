package engine

import (
	"strings"
	"testing"
)

func TestNPCVerbInteractions(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*GameV2) // Setup function to position NPCs/items
		command  string
		contains []string // Strings that should appear in output
		notContains []string // Strings that should NOT appear
	}{
		// EXAMINE NPC tests
		{
			name: "examine troll",
			setup: func(g *GameV2) {
				g.Location = "troll-room"
			},
			command:  "examine troll",
			contains: []string{"nasty-looking troll", "bloody axe"},
		},
		{
			name: "examine thief",
			setup: func(g *GameV2) {
				// Move thief to current room
				g.NPCs["thief"].Location = "west-of-house"
				g.Location = "west-of-house"
				g.Rooms["west-of-house"].AddNPC("thief")
			},
			command:  "examine thief",
			contains: []string{"suspicious-looking individual"},
		},

		// TAKE NPC tests
		{
			name: "take living npc",
			setup: func(g *GameV2) {
				g.Location = "troll-room"
			},
			command:  "take troll",
			contains: []string{"wouldn't hear of it"},
		},

		// TALK/TELL NPC tests
		{
			name: "talk to troll",
			setup: func(g *GameV2) {
				g.Location = "troll-room"
			},
			command:  "talk to troll",
			contains: []string{"grunts", "axe menacingly"},
		},
		{
			name: "tell thief hello",
			setup: func(g *GameV2) {
				g.NPCs["thief"].Location = "west-of-house"
				g.Location = "west-of-house"
				g.Rooms["west-of-house"].AddNPC("thief")
			},
			command:  "tell thief hello",
			contains: []string{"nothing to say"},
		},
		{
			name: "ask cyclops about lunch",
			setup: func(g *GameV2) {
				g.Location = "cyclops-room"
			},
			command:  "ask cyclops",
			contains: []string{"Me eat you"},
		},
		{
			name: "hello to bat",
			setup: func(g *GameV2) {
				g.Location = "bat-room"
			},
			command:  "hello bat",
			contains: []string{"screeches"},
		},

		// SMELL NPC tests
		{
			name: "smell troll",
			setup: func(g *GameV2) {
				g.Location = "troll-room"
			},
			command:  "smell troll",
			contains: []string{"rotten meat", "wet dog"},
		},
		{
			name: "smell thief",
			setup: func(g *GameV2) {
				g.NPCs["thief"].Location = "west-of-house"
				g.Location = "west-of-house"
				g.Rooms["west-of-house"].AddNPC("thief")
			},
			command:  "smell thief",
			contains: []string{"sweat", "greed"},
		},
		{
			name: "smell cyclops",
			setup: func(g *GameV2) {
				g.Location = "cyclops-room"
			},
			command:  "smell cyclops",
			contains: []string{"garlic", "peppers"},
		},
		{
			name: "smell bat",
			setup: func(g *GameV2) {
				g.Location = "bat-room"
			},
			command:  "smell bat",
			contains: []string{"bat guano"},
		},
		{
			name: "smell ghosts",
			setup: func(g *GameV2) {
				g.Location = "entrance-to-hades"
			},
			command:  "smell ghosts",
			contains: []string{"no smell", "incorporeal"},
		},

		// TOUCH NPC tests
		{
			name: "touch aggressive troll",
			setup: func(g *GameV2) {
				g.Location = "troll-room"
				// Troll should already be in troll-room from initialization
			},
			command:  "touch troll",
			contains: []string{"extremely dangerous"},
		},
		{
			name: "touch dead npc",
			setup: func(g *GameV2) {
				g.Location = "troll-room"
				// Troll should already be in troll-room from initialization
				g.NPCs["troll"].Flags.IsAlive = false
			},
			command:  "touch troll",
			contains: []string{"quite dead"},
		},

		// SEARCH NPC tests
		{
			name: "search living npc",
			setup: func(g *GameV2) {
				g.Location = "troll-room"
			},
			command:  "search troll",
			contains: []string{"wouldn't appreciate"},
		},
		{
			name: "search dead npc with inventory",
			setup: func(g *GameV2) {
				g.Location = "troll-room"
				g.NPCs["troll"].Flags.IsAlive = false
				g.NPCs["troll"].Inventory = []string{"axe"}
				// Make sure axe exists and is in NPC's inventory
				if axe := g.Items["axe"]; axe != nil {
					axe.Location = "troll"
				}
			},
			command:  "search troll",
			contains: []string{"Searching", "find"},
		},
		{
			name: "search dead npc with empty inventory",
			setup: func(g *GameV2) {
				g.Location = "troll-room"
				g.NPCs["troll"].Flags.IsAlive = false
				g.NPCs["troll"].Inventory = []string{}
			},
			command:  "search troll",
			contains: []string{"nothing of interest"},
		},

		// THROW AT NPC tests
		{
			name: "throw item at npc",
			setup: func(g *GameV2) {
				g.Location = "troll-room"
				// Give player the knife
				g.Player.Inventory = append(g.Player.Inventory, "knife")
				if knife := g.Items["knife"]; knife != nil {
					knife.Location = "inventory"
				}
			},
			command:  "throw knife at troll",
			contains: []string{"bounces harmlessly"},
		},
		{
			name: "throw at non-existent target",
			setup: func(g *GameV2) {
				g.Location = "west-of-house"
				g.Player.Inventory = append(g.Player.Inventory, "knife")
				if knife := g.Items["knife"]; knife != nil {
					knife.Location = "inventory"
				}
			},
			command:  "throw knife at tree",
			contains: []string{"can't see"}, // Tree isn't at west-of-house
		},

		// Test that inappropriate commands fail gracefully
		{
			name: "try to open npc",
			setup: func(g *GameV2) {
				g.Location = "troll-room"
			},
			command:     "open troll",
			notContains: []string{"The troll is now open"}, // Should fail
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGameV2("test")

			// Run setup
			if tt.setup != nil {
				tt.setup(g)
			}

			// Execute command
			result := g.Process(tt.command)

			// Check for expected strings
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected output to contain %q, got:\n%s", expected, result)
				}
			}

			// Check for unexpected strings
			for _, unexpected := range tt.notContains {
				if strings.Contains(result, unexpected) {
					t.Errorf("Expected output NOT to contain %q, got:\n%s", unexpected, result)
				}
			}
		})
	}
}

func TestNPCTalkResponses(t *testing.T) {
	g := NewGameV2("test")

	tests := []struct {
		npcID    string
		location string
		command  string
		contains string
	}{
		{"troll", "troll-room", "talk to troll", "grunts"},
		{"thief", "maze-1", "talk to thief", "nothing to say"},
		{"cyclops", "cyclops-room", "ask cyclops", "eat you"},
		{"bat", "bat-room", "talk to bat", "screeches"},
		{"ghosts", "entrance-to-hades", "talk to ghosts", "moan"},
	}

	for _, tt := range tests {
		t.Run("talk_to_"+tt.npcID, func(t *testing.T) {
			g.Location = tt.location
			result := g.Process(tt.command)

			if !strings.Contains(result, tt.contains) {
				t.Errorf("Expected %q to contain %q, got: %s", tt.command, tt.contains, result)
			}
		})
	}
}

func TestSearchDeadThief(t *testing.T) {
	g := NewGameV2("test")
	g.Location = "maze-1"

	// Setup: Kill thief and give him some loot
	thief := g.NPCs["thief"]
	thief.Flags.IsAlive = false
	thief.Inventory = []string{"sword", "lamp"}

	// Make items exist and be in thief's inventory
	if sword := g.Items["sword"]; sword != nil {
		sword.Location = "thief"
	}
	if lamp := g.Items["lamp"]; lamp != nil {
		lamp.Location = "thief"
	}

	// Search the dead thief
	result := g.Process("search thief")

	// Should find items
	if !strings.Contains(result, "Searching") {
		t.Errorf("Expected to find items on dead thief, got: %s", result)
	}

	// Items should now be in the room
	sword := g.Items["sword"]
	lamp := g.Items["lamp"]

	if sword.Location != "maze-1" {
		t.Errorf("Expected sword to be in maze-1, got: %s", sword.Location)
	}
	if lamp.Location != "maze-1" {
		t.Errorf("Expected lamp to be in maze-1, got: %s", lamp.Location)
	}

	// Thief's inventory should be empty
	if len(thief.Inventory) != 0 {
		t.Errorf("Expected thief inventory to be empty after search, got %d items", len(thief.Inventory))
	}
}
