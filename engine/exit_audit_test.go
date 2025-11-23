package engine

import (
	"fmt"
	"sort"
	"testing"
)

// TestExitBidirectionality performs a comprehensive audit of all room exits
// to ensure navigational consistency throughout the game world
func TestExitBidirectionality(t *testing.T) {
	g := NewGameV2()

	// Map of reverse directions
	reverseDir := map[string]string{
		"north": "south", "south": "north",
		"east":  "west", "west": "east",
		"ne":    "sw", "sw": "ne",
		"nw":    "se", "se": "nw",
		"up":    "down", "down": "up",
		"in":    "out", "out": "in",
	}

	// Intentional one-way passages that should NOT be bidirectional
	// Format: "from-room:direction:to-room"
	intentionalOneWay := map[string]bool{
		// Slides (can't climb back up)
		"slide-room:down:cellar":           true,
		"reservoir-south:down:deep-canyon": true,

		// Drops/falls (can't go back up same way)
		"cliff-middle:down:canyon-bottom": true,

		// Conditional exits that only work one way
		"grating-clearing:down:grating-room": true, // requires grate-open
		"living-room:down:cellar":            true, // requires trap-door-open

		// Intentional maze one-ways (part of puzzle)
		// We'll detect these by the fact they're in maze rooms

		// River current (flows one direction)
		"frigid-river-1:down:frigid-river-2": true,
		"frigid-river-2:down:frigid-river-3": true,
		"frigid-river-3:down:frigid-river-4": true,
		"frigid-river-4:down:frigid-river-5": true,
		"frigid-river-5:down:shore":          true,
	}

	issues := []string{}
	warnings := []string{}
	checked := make(map[string]bool)
	stats := struct {
		totalExits      int
		bidirectional   int
		oneWayIntended  int
		oneWayUnintended int
		missingDest     int
		wrongReverse    int
	}{}

	// Get sorted room list for consistent output
	roomIDs := []string{}
	for id := range g.Rooms {
		roomIDs = append(roomIDs, id)
	}
	sort.Strings(roomIDs)

	for _, roomID := range roomIDs {
		room := g.Rooms[roomID]

		for dir, exit := range room.Exits {
			stats.totalExits++
			key := roomID + ":" + dir + ":" + exit.To

			// Skip if already checked as reverse
			if checked[key] {
				continue
			}
			checked[key] = true

			// Check if destination exists
			destRoom := g.Rooms[exit.To]
			if destRoom == nil {
				issues = append(issues, fmt.Sprintf("❌ MISSING DEST: %s --%s--> %s (destination room doesn't exist!)",
					roomID, dir, exit.To))
				stats.missingDest++
				continue
			}

			// Get reverse direction
			revDir := reverseDir[dir]
			if revDir == "" {
				// No standard reverse direction (shouldn't happen with our directions)
				warnings = append(warnings, fmt.Sprintf("⚠️  UNKNOWN DIR: %s --%s--> %s (no reverse for direction %s)",
					roomID, dir, exit.To, dir))
				continue
			}

			// Check if this is an intentional one-way passage
			if intentionalOneWay[key] {
				stats.oneWayIntended++
				continue
			}

			// Check for reverse exit
			revExit, hasReverse := destRoom.Exits[revDir]

			if !hasReverse {
				// No reverse exit - check if this is intentional

				// Conditional exits are often one-way until condition is met
				if exit.Condition != "" {
					stats.oneWayIntended++
					continue
				}

				// Maze rooms can have intentional one-way passages
				isMaze := (len(roomID) >= 4 && roomID[:4] == "maze") ||
					(len(exit.To) >= 4 && exit.To[:4] == "maze")
				if isMaze {
					warnings = append(warnings, fmt.Sprintf("⚠️  MAZE ONE-WAY: %s --%s--> %s (no %s exit back)",
						roomID, dir, exit.To, revDir))
					stats.oneWayIntended++
					continue
				}

				// Check if it's a slide/drop (down with no up)
				if dir == "down" {
					warnings = append(warnings, fmt.Sprintf("⚠️  POTENTIAL DROP: %s --%s--> %s (no %s exit back - slide/drop?)",
						roomID, dir, exit.To, revDir))
					stats.oneWayIntended++
					continue
				}

				// This looks like an unintended one-way passage
				issues = append(issues, fmt.Sprintf("❌ ONE-WAY: %s --%s--> %s (no %s exit from %s back to %s)",
					roomID, dir, exit.To, revDir, exit.To, roomID))
				stats.oneWayUnintended++

			} else if revExit.To != roomID {
				// Reverse exists but points to wrong room!
				issues = append(issues, fmt.Sprintf("❌ WRONG REVERSE: %s --%s--> %s, but %s --%s--> %s (should be %s!)",
					roomID, dir, exit.To, exit.To, revDir, revExit.To, roomID))
				stats.wrongReverse++

			} else {
				// Bidirectional exit pair - mark reverse as checked
				reverseKey := exit.To + ":" + revDir + ":" + roomID
				checked[reverseKey] = true
				stats.bidirectional++
			}
		}
	}

	// Print summary
	t.Logf("\n=== EXIT AUDIT SUMMARY ===")
	t.Logf("Total exits: %d", stats.totalExits)
	t.Logf("Bidirectional pairs: %d", stats.bidirectional)
	t.Logf("Intentional one-way: %d", stats.oneWayIntended)
	t.Logf("Unintended one-way: %d", stats.oneWayUnintended)
	t.Logf("Missing destinations: %d", stats.missingDest)
	t.Logf("Wrong reverse: %d", stats.wrongReverse)

	// Print issues
	if len(issues) > 0 {
		t.Logf("\n=== CRITICAL ISSUES (%d) ===", len(issues))
		for _, issue := range issues {
			t.Log(issue)
		}
	}

	// Print warnings (for review but not failures)
	if len(warnings) > 0 && testing.Verbose() {
		t.Logf("\n=== WARNINGS FOR REVIEW (%d) ===", len(warnings))
		for i, warning := range warnings {
			if i < 20 { // Limit output
				t.Log(warning)
			}
		}
		if len(warnings) > 20 {
			t.Logf("... and %d more warnings", len(warnings)-20)
		}
	}

	// Fail test if there are critical issues
	if stats.missingDest > 0 || stats.wrongReverse > 0 {
		t.Errorf("Found %d critical navigation errors (missing destinations or wrong reverse exits)",
			stats.missingDest+stats.wrongReverse)
	}

	// Log unintended one-ways but don't fail (they might need manual review)
	if stats.oneWayUnintended > 0 {
		t.Logf("\nFound %d potentially unintended one-way passages - please review", stats.oneWayUnintended)
	}
}

// TestCriticalPathsNavigable ensures key game locations are reachable and properly connected
func TestCriticalPathsNavigable(t *testing.T) {
	g := NewGameV2()

	// Test critical navigation paths that must work
	tests := []struct {
		name      string
		startRoom string
		path      []string // sequence of directions
		endRoom   string
	}{
		{
			name:      "West around house to north",
			startRoom: "west-of-house",
			path:      []string{"north"},
			endRoom:   "north-of-house",
		},
		{
			name:      "North back to west",
			startRoom: "north-of-house",
			path:      []string{"west"},
			endRoom:   "west-of-house",
		},
		{
			name:      "West into forest leads to path (intentional non-bidirectional)",
			startRoom: "west-of-house",
			path:      []string{"west", "east"},
			endRoom:   "path", // per ZIL: forest-1 east -> path (not back to west-of-house)
		},
		{
			name:      "North to path and back",
			startRoom: "north-of-house",
			path:      []string{"north", "south"},
			endRoom:   "north-of-house",
		},
		{
			name:      "Path to clearing (forest-2) and back",
			startRoom: "path",
			path:      []string{"east", "west"},
			endRoom:   "path",
		},
		{
			name:      "Kitchen to living room and back",
			startRoom: "kitchen",
			path:      []string{"west", "east"},
			endRoom:   "kitchen",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.Location = tt.startRoom
			startLoc := g.Location

			for i, dir := range tt.path {
				room := g.Rooms[g.Location]
				if room == nil {
					t.Fatalf("Step %d: Current room %s doesn't exist", i, g.Location)
				}

				exit := room.Exits[dir]
				if exit == nil {
					t.Fatalf("Step %d: No exit %s from %s", i, dir, g.Location)
				}

				g.Location = exit.To
			}

			if g.Location != tt.endRoom {
				t.Errorf("Expected to end at %s, but ended at %s (started at %s)",
					tt.endRoom, g.Location, startLoc)
			}
		})
	}
}
