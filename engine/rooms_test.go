package engine

import (
	"strings"
	"testing"
)

func TestAllRoomsCreated(t *testing.T) {
	g := NewGameV2()

	// Verify we have exactly 110 rooms
	if len(g.Rooms) != 110 {
		t.Errorf("Expected 110 rooms, got %d", len(g.Rooms))
	}
}

func TestKeyRoomsExist(t *testing.T) {
	g := NewGameV2()

	keyRooms := []string{
		"west-of-house",
		"north-of-house",
		"behind-house",
		"kitchen",
		"living-room",
		"cellar",
		"troll-room",
		"maze-1",
		"treasure-room",
		"entrance-to-hades",
		"dam-room",
		"cyclops-room",
		"reservoir",
		"stream-view",
	}

	for _, roomID := range keyRooms {
		if g.Rooms[roomID] == nil {
			t.Errorf("Key room %s does not exist", roomID)
		}
	}
}

func TestRoomConnectivity(t *testing.T) {
	g := NewGameV2()

	tests := []struct {
		from      string
		direction string
		to        string
	}{
		// Around the house (per ZIL - intentionally non-bidirectional)
		{"west-of-house", "north", "north-of-house"},
		{"north-of-house", "west", "west-of-house"}, // north-of-house has no south exit per ZIL
		{"south-of-house", "east", "behind-house"},
		{"behind-house", "north", "north-of-house"},

		// Into the house (requires window-open flag, tested elsewhere)
		{"kitchen", "west", "living-room"},

		// Down to cellar
		{"living-room", "east", "kitchen"},

		// Underground passages
		{"cellar", "north", "troll-room"},
		{"troll-room", "south", "cellar"},
		{"ew-passage", "west", "troll-room"},

		// Forest paths
		{"north-of-house", "north", "path"},
		{"path", "east", "forest-2"},

		// Maze entries
		{"troll-room", "west", "maze-1"},
	}

	for _, tt := range tests {
		room := g.Rooms[tt.from]
		if room == nil {
			t.Errorf("Source room %s does not exist", tt.from)
			continue
		}

		exit := room.Exits[tt.direction]
		if exit == nil {
			t.Errorf("Room %s has no exit in direction %s", tt.from, tt.direction)
			continue
		}

		if exit.To != tt.to {
			t.Errorf("Room %s exit %s leads to %s, expected %s",
				tt.from, tt.direction, exit.To, tt.to)
		}

		// Verify destination room exists
		if g.Rooms[exit.To] == nil {
			t.Errorf("Destination room %s does not exist (from %s via %s)",
				exit.To, tt.from, tt.direction)
		}
	}
}

func TestDarkRooms(t *testing.T) {
	g := NewGameV2()

	darkRooms := []string{
		"cellar",
		"troll-room",
		"ew-passage",
		"maze-1",
		"maze-2",
	}

	for _, roomID := range darkRooms {
		room := g.Rooms[roomID]
		if room == nil {
			t.Errorf("Dark room %s does not exist", roomID)
			continue
		}

		if !room.Flags.IsDark {
			t.Errorf("Room %s should be dark but IsDark flag is false", roomID)
		}
	}
}

func TestOutdoorRooms(t *testing.T) {
	g := NewGameV2()

	outdoorRooms := []string{
		"west-of-house",
		"north-of-house",
		"south-of-house",
		"behind-house",
		"path",
		"clearing",
	}

	for _, roomID := range outdoorRooms {
		room := g.Rooms[roomID]
		if room == nil {
			t.Errorf("Outdoor room %s does not exist", roomID)
			continue
		}

		if !room.Flags.IsOutdoors {
			t.Errorf("Room %s should be outdoors but IsOutdoors flag is false", roomID)
		}
	}
}

func TestConditionalExits(t *testing.T) {
	g := NewGameV2()

	// Troll room east exit should be conditional on troll-dead flag
	trollRoom := g.Rooms["troll-room"]
	if trollRoom == nil {
		t.Fatal("Troll room does not exist")
	}

	eastExit := trollRoom.Exits["east"]
	if eastExit == nil {
		t.Fatal("Troll room should have an east exit")
	}

	if eastExit.Condition != "troll-dead" {
		t.Errorf("Troll room east exit should require troll-dead flag, got %s",
			eastExit.Condition)
	}

	// Living room down exit should be conditional on trap-door-open flag
	livingRoom := g.Rooms["living-room"]
	if livingRoom == nil {
		t.Fatal("Living room does not exist")
	}

	downExit := livingRoom.Exits["down"]
	if downExit == nil {
		t.Fatal("Living room should have a down exit")
	}

	if downExit.Condition != "trap-door-open" {
		t.Errorf("Living room down exit should require trap-door-open flag, got %s",
			downExit.Condition)
	}
}

func TestMazeRooms(t *testing.T) {
	g := NewGameV2()

	// Verify all 15 maze rooms exist
	for i := 1; i <= 15; i++ {
		roomID := ""
		if i < 10 {
			roomID = "maze-" + string(rune('0'+i))
		} else {
			roomID = "maze-1" + string(rune('0'+(i-10)))
		}

		// Use a simpler approach
		switch i {
		case 1:
			roomID = "maze-1"
		case 2:
			roomID = "maze-2"
		case 3:
			roomID = "maze-3"
		case 4:
			roomID = "maze-4"
		case 5:
			roomID = "maze-5"
		case 6:
			roomID = "maze-6"
		case 7:
			roomID = "maze-7"
		case 8:
			roomID = "maze-8"
		case 9:
			roomID = "maze-9"
		case 10:
			roomID = "maze-10"
		case 11:
			roomID = "maze-11"
		case 12:
			roomID = "maze-12"
		case 13:
			roomID = "maze-13"
		case 14:
			roomID = "maze-14"
		case 15:
			roomID = "maze-15"
		}

		if g.Rooms[roomID] == nil {
			t.Errorf("Maze room %s does not exist", roomID)
		}
	}

	// Verify dead end rooms exist
	deadEnds := []string{"dead-end-1", "dead-end-2", "dead-end-3", "dead-end-4"}
	for _, roomID := range deadEnds {
		if g.Rooms[roomID] == nil {
			t.Errorf("Dead end room %s does not exist", roomID)
		}
	}
}

// TestWestEastNavigation tests the specific bug reported: west-of-house -> west -> forest-1 -> east should return
func TestWestEastNavigation(t *testing.T) {
	g := NewGameV2()
	g.Location = "west-of-house"

	// Go west to forest
	result := g.Process("west")
	if !strings.Contains(result, "Forest") {
		t.Errorf("Expected to be in forest after going west, got: %s", result)
	}
	if g.Location != "forest-1" {
		t.Errorf("Expected location to be forest-1, got: %s", g.Location)
	}

	// Go east - per ZIL, forest-1 east leads to path (not back to west-of-house)
	// The forest is intentionally a mild navigation puzzle
	result = g.Process("east")
	if !strings.Contains(result, "path") && !strings.Contains(result, "Path") {
		t.Errorf("Expected to be on path after going east from forest-1, got: %s", result)
	}
	if g.Location != "path" {
		t.Errorf("Expected location to be path, got: %s", g.Location)
	}

	// To return to west-of-house from forest-1, player must discover the correct route
	// (This tests that the ZIL non-bidirectional forest topology is implemented correctly)
}
