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

		// House area - ZIL intentionally blocks passage through house
		// (north and south sides don't connect directly - must go around)
		"behind-house:north:north-of-house": true, // ZIL: north/south sides don't connect through house
		"behind-house:south:south-of-house": true,
		"north-of-house:west:west-of-house":  true, // ZIL: west side doesn't connect back east
		"south-of-house:west:west-of-house":  true,
		"west-of-house:north:north-of-house": true,
		"west-of-house:south:south-of-house": true,
		"north-of-house:east:behind-house":   true, // ZIL: behind-house only has conditional west→kitchen
		"south-of-house:east:behind-house":   true,
		"south-of-house:south:forest-3":      true, // ZIL: forest-3 north→clearing (maze-like)
		"west-of-house:west:forest-1":        true, // ZIL: forest-1 east→path (maze-like)

		// Forest area - ZIL intentionally has maze-like non-bidirectional navigation
		"forest-1:north:grating-clearing": true, // grating-clearing south→path (not back to forest-1)
		"forest-1:south:forest-3":         true, // forest-3 north→clearing (not back to forest-1)
		"forest-3:west:forest-1":          true, // forest-1 east→path (not back to forest-3)
		"forest-3:nw:south-of-house":      true, // south-of-house has no se exit back
		"grating-clearing:east:forest-2":  true, // forest-2 west→path (not back to grating-clearing)
		"grating-clearing:west:forest-1":  true, // forest-1 east→path (not back to grating-clearing)
		"clearing:east:canyon-view":       true, // canyon-view nw→clearing (not west back)
		"mountains:north:forest-2":        true, // forest-2 south→clearing (not back to mountains)
		"mountains:south:forest-2":        true, // forest-2 south→clearing (not north back to mountains)

		// Dam/Canyon area - ZIL intentionally has one-way passages
		"dam-base:north:dam-room":            true, // dam-room south→deep-canyon (not back to dam-base, return via down)
		"dam-room:east:dam-base":             true, // dam-base has no west exit (return via up/north)
		"dam-room:south:deep-canyon":         true, // deep-canyon has no north exit
		"deep-canyon:east:dam-room":          true, // dam-room west→reservoir-south (not back)
		"river-1:west:dam-base":              true, // dam-base has no east exit
		"river-1:land:dam-base":              true, // dam-base has no east exit
		"canyon-view:east:cliff-middle":      true, // cliff-middle has no west exit (go down to canyon)
		"canyon-view:nw:clearing":            true, // clearing has no se exit (already listed above as clearing:east:canyon-view)
		"canyon-bottom:north:end-of-rainbow": true, // end-of-rainbow sw→canyon-bottom (different direction, but bidirectional)
		"end-of-rainbow:sw:canyon-bottom":    true, // canyon-bottom north→end-of-rainbow (different direction, but bidirectional)
		"ew-passage:north:chasm-room":        true, // chasm-room sw→ew-passage (different direction, but bidirectional)

		// Mine area - ZIL intentionally has confusing self-loops and one-ways
		"mine-1:east:mine-1":    true, // intentional self-loop per ZIL
		"mine-2:north:mine-2":   true, // intentional self-loop per ZIL
		"mine-3:south:mine-3":   true, // intentional self-loop per ZIL
		"mine-4:west:mine-4":    true, // intentional self-loop per ZIL
		"mine-1:ne:mine-2":      true, // mine-2 has no sw exit back
		"mine-2:south:mine-1":   true, // mine-1 north→gas-room (not back to mine-2)
		"mine-2:se:mine-3":      true, // mine-3 has no nw exit back
		"mine-3:east:mine-2":    true, // mine-2 has no west exit back
		"mine-3:sw:mine-4":      true, // mine-4 has no ne exit back
		"mine-4:north:mine-3":   true, // mine-3 south is self-loop, not back to mine-4
		"gas-room:east:mine-1":  true, // mine-1 has no west exit back
		"gas-room:south:mine-1": true, // mine-1 north goes to gas-room but different direction

		// Mirror rooms - ZIL intentionally has confusing non-bidirectional passages
		"mirror-room-1:east:small-cave":        true, // small-cave west→twisting-passage (triangle maze)
		"mirror-room-1:west:twisting-passage":  true, // twisting-passage east→small-cave (triangle maze)
		"twisting-passage:north:mirror-room-1": true, // mirror-room-1 west→twisting-passage (different direction)
		"twisting-passage:east:small-cave":     true, // small-cave west→twisting-passage (return)
		"small-cave:north:mirror-room-1":       true, // mirror-room-1 east→small-cave (different direction)
		"small-cave:west:twisting-passage":     true, // twisting-passage east→small-cave (return)
		"mirror-room-2:east:tiny-cave":         true, // tiny-cave west→winding-passage (triangle maze)
		"mirror-room-2:west:winding-passage":   true, // winding-passage east→tiny-cave (triangle maze)
		"winding-passage:north:mirror-room-2":  true, // mirror-room-2 west→winding-passage (different direction)
		"winding-passage:east:tiny-cave":       true, // tiny-cave west→winding-passage (return)
		"tiny-cave:north:mirror-room-2":        true, // mirror-room-2 east→tiny-cave (different direction)
		"tiny-cave:west:winding-passage":       true, // winding-passage east→tiny-cave (return)
		"small-cave:down:atlantis-room":        true, // atlantis-room up→small-cave (different direction)
		"small-cave:south:atlantis-room":       true, // atlantis-room has no north exit

		// Misc intentional one-ways
		"river-3:west:white-cliffs-north":  true, // white-cliffs has no east exit (one-way access)
		"river-3:land:white-cliffs-north":  true, // white-cliffs has no east exit (one-way access)
		"river-4:west:white-cliffs-south":  true, // white-cliffs has no east exit (one-way access)
		"strange-passage:west:cyclops-room": true, // cyclops-room east is conditional (magic-flag)
		"strange-passage:in:cyclops-room":   true, // cyclops-room east is conditional (magic-flag)

		// Intentional maze one-ways and wrong reverses (part of puzzle - maze is deliberately confusing)
		"maze-1:north:maze-1":        true, // self-loop, south→maze-2
		"maze-1:west:maze-4":         true, // maze-4 east→dead-end-1 (not back)
		"maze-2:south:maze-1":        true, // maze-1 north→maze-1 self-loop (not back)
		"maze-4:north:maze-1":        true, // maze-1 south→maze-2 (not back)
		"maze-6:west:maze-6":         true, // self-loop, east→maze-7
		"maze-7:east:maze-8":         true, // maze-8 west→maze-8 self-loop (not back)
		"maze-9:east:maze-10":        true, // maze-10 west→maze-13 (not back)
		"maze-9:west:maze-12":        true, // maze-12 east→maze-13 (not back)
		"maze-10:east:maze-9":        true, // maze-9 west→maze-12 (not back)
		"maze-10:west:maze-13":       true, // maze-13 east→maze-9 (not back)
		"maze-12:sw:maze-11":         true, // maze-11 ne→grating-room (not back)
		"maze-12:east:maze-13":       true, // maze-13 west→maze-11 (not back)
		"maze-13:down:maze-12":       true, // maze-12 up→maze-9 (not back)
		"maze-13:east:maze-9":        true, // maze-9 west→maze-12 (not back)
		"dead-end-1:south:maze-4":    true, // maze-4 north→maze-1 (not back)
		// We still detect remaining one-way maze passages by the fact they're in maze rooms

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
