package engine

import (
	"strings"
	"testing"
)

func TestPutCommand(t *testing.T) {
	g := NewGameV2()

	tests := []struct {
		name     string
		setup    func()
		command  string
		contains string
	}{
		{
			name: "put leaflet in mailbox",
			setup: func() {
				// Take leaflet first
				g.Process("take leaflet")
			},
			command:  "put leaflet in mailbox",
			contains: "Done",
		},
		{
			name: "put item not in inventory",
			setup: func() {
				// Don't take anything
			},
			command:  "put sword in case",
			contains: "don't have",
		},
		{
			name: "put into closed container",
			setup: func() {
				// Create a closed container
				coffin := g.Items["coffin"]
				coffin.Location = "west-of-house"
				g.Rooms["west-of-house"].AddItem("coffin")
				coffin.Flags.IsOpen = false

				g.Process("take leaflet")
			},
			command:  "put leaflet in coffin",
			contains: "closed",
		},
		{
			name: "put into non-container",
			setup: func() {
				// Try to put something into the lamp (not a container)
				lamp := g.Items["lamp"]
				lamp.Location = "west-of-house"
				g.Rooms["west-of-house"].AddItem("lamp")

				g.Process("take leaflet")
			},
			command:  "put leaflet in lamp",
			contains: "can't put",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset game
			g = NewGameV2()
			g.Location = "west-of-house"

			// Run setup
			if tt.setup != nil {
				tt.setup()
			}

			// Run command
			result := g.Process(tt.command)

			if !strings.Contains(result, tt.contains) {
				t.Errorf("Expected %q to contain %q, got: %q", tt.command, tt.contains, result)
			}
		})
	}
}

func TestGiveCommand(t *testing.T) {
	g := NewGameV2()
	g.Location = "troll-room"

	// Give food to troll
	g.Player.Inventory = append(g.Player.Inventory, "lunch")
	g.Items["lunch"].Location = "inventory"

	result := g.Process("give lunch to troll")

	if !strings.Contains(result, "troll") {
		t.Errorf("Expected troll response, got: %s", result)
	}
}

func TestAttackCommand(t *testing.T) {
	g := NewGameV2()
	g.Location = "troll-room"

	result := g.Process("attack troll")

	if !strings.Contains(result, "suicidal") || !strings.Contains(result, "bare hands") {
		t.Errorf("Expected bare hands warning, got: %s", result)
	}
}

func TestWaveCommand(t *testing.T) {
	g := NewGameV2()

	// Add sceptre to inventory
	g.Player.Inventory = append(g.Player.Inventory, "sceptre")
	g.Items["sceptre"].Location = "inventory"

	result := g.Process("wave sceptre")

	if !strings.Contains(result, "glow") {
		t.Errorf("Expected sceptre to glow, got: %s", result)
	}
}

func TestClimbCommand(t *testing.T) {
	g := NewGameV2()

	result := g.Process("climb tree")

	if result == "" {
		t.Error("Expected response for climb command")
	}
}

func TestTieUntieCommands(t *testing.T) {
	g := NewGameV2()

	// Add rope to inventory
	g.Player.Inventory = append(g.Player.Inventory, "rope")
	g.Items["rope"].Location = "inventory"

	result := g.Process("tie rope to railing")
	if result == "" {
		t.Error("Expected response for tie command")
	}

	result = g.Process("untie rope")
	if !strings.Contains(result, "not tied") {
		t.Errorf("Expected 'not tied' message, got: %s", result)
	}
}

func TestDigCommand(t *testing.T) {
	g := NewGameV2()

	result := g.Process("dig")

	if !strings.Contains(result, "hard") || !strings.Contains(result, "ground") {
		t.Errorf("Expected hard ground message, got: %s", result)
	}
}

func TestPushPullCommands(t *testing.T) {
	g := NewGameV2()
	g.Location = "machine-room"

	// Push button
	result := g.Process("push red button")
	if !strings.Contains(result, "Click") {
		t.Errorf("Expected button click, got: %s", result)
	}

	// Pull something
	result = g.Process("pull lever")
	if result == "" {
		t.Error("Expected response for pull command")
	}
}

func TestRingCommand(t *testing.T) {
	g := NewGameV2()

	// Add bell to room
	g.Items["bell"].Location = g.Location

	result := g.Process("ring bell")

	if !strings.Contains(result, "Ding") {
		t.Errorf("Expected bell to ring, got: %s", result)
	}
}

func TestPrayCommand(t *testing.T) {
	g := NewGameV2()

	result := g.Process("pray")

	if !strings.Contains(result, "pray") {
		t.Errorf("Expected prayer response, got: %s", result)
	}
}

func TestWaitCommand(t *testing.T) {
	g := NewGameV2()

	result := g.Process("wait")

	if !strings.Contains(result, "Time passes") {
		t.Errorf("Expected time to pass, got: %s", result)
	}

	// Test 'z' synonym
	result = g.Process("z")

	if !strings.Contains(result, "Time passes") {
		t.Errorf("Expected time to pass for 'z', got: %s", result)
	}
}

func TestEatCommand(t *testing.T) {
	g := NewGameV2()

	tests := []struct {
		name     string
		item     string
		contains string
	}{
		{"eat lunch", "lunch", "hit the spot"},
		{"eat garlic", "garlic", "vampires"},
		{"eat lamp", "lamp", "would agree"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g = NewGameV2()

			// Add item to inventory
			g.Player.Inventory = append(g.Player.Inventory, tt.item)
			g.Items[tt.item].Location = "inventory"

			result := g.Process("eat " + tt.item)

			if !strings.Contains(result, tt.contains) {
				t.Errorf("Expected %q, got: %s", tt.contains, result)
			}
		})
	}
}

func TestDrinkCommand(t *testing.T) {
	g := NewGameV2()

	// Add water to room
	g.Items["water"].Location = g.Location

	result := g.Process("drink water")

	if !strings.Contains(result, "thirsty") {
		t.Errorf("Expected thirsty response, got: %s", result)
	}
}

func TestFillPourCommands(t *testing.T) {
	g := NewGameV2()

	// Add bottle to inventory
	g.Player.Inventory = append(g.Player.Inventory, "bottle")
	g.Items["bottle"].Location = "inventory"

	result := g.Process("fill bottle")

	if !strings.Contains(result, "nothing") {
		t.Errorf("Expected 'nothing to fill' message, got: %s", result)
	}

	result = g.Process("pour bottle")

	if !strings.Contains(result, "empty") {
		t.Errorf("Expected empty message, got: %s", result)
	}
}

func TestListenCommand(t *testing.T) {
	g := NewGameV2()

	result := g.Process("listen")

	if !strings.Contains(result, "nothing unusual") {
		t.Errorf("Expected nothing unusual, got: %s", result)
	}
}

func TestSmellTouchCommands(t *testing.T) {
	g := NewGameV2()
	g.Location = "west-of-house"

	result := g.Process("smell mailbox")

	if result == "" {
		t.Error("Expected response for smell command")
	}

	result = g.Process("touch mailbox")

	if result == "" {
		t.Error("Expected response for touch command")
	}
}

func TestSearchCommand(t *testing.T) {
	g := NewGameV2()
	g.Location = "west-of-house"

	// Search without object (should be like look)
	result := g.Process("search")

	if !strings.Contains(result, "West of House") {
		t.Errorf("Expected room description, got: %s", result)
	}

	// Search a container (should be like look in)
	result = g.Process("search mailbox")

	if !strings.Contains(result, "leaflet") || !strings.Contains(result, "contains") {
		t.Errorf("Expected container contents, got: %s", result)
	}
}

func TestJumpSwimCommands(t *testing.T) {
	g := NewGameV2()

	result := g.Process("jump")

	if !strings.Contains(result, "fruitlessly") {
		t.Errorf("Expected fruitless jump, got: %s", result)
	}

	result = g.Process("swim")

	if !strings.Contains(result, "no water") {
		t.Errorf("Expected no water message, got: %s", result)
	}
}

func TestBlowKnockCommands(t *testing.T) {
	g := NewGameV2()

	result := g.Process("blow whistle")

	if result == "" {
		t.Error("Expected response for blow command")
	}

	result = g.Process("knock on door")

	if !strings.Contains(result, "No one answers") {
		t.Errorf("Expected no answer, got: %s", result)
	}
}

func TestScoreCommand(t *testing.T) {
	g := NewGameV2()

	result := g.Process("score")

	if !strings.Contains(result, "350") {
		t.Errorf("Expected score out of 350, got: %s", result)
	}

	if !strings.Contains(result, "move") {
		t.Errorf("Expected move count, got: %s", result)
	}
}
