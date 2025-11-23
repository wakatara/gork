package engine

import (
	"fmt"
	"strings"
	"testing"
)

// TestRugTrapDoorPuzzle tests the complete rug/trap door puzzle sequence
func TestRugTrapDoorPuzzle(t *testing.T) {
	g := NewGameV2()

	// Navigate to living room
	g.Process("north")
	g.Process("east")
	g.Process("open window")
	g.Process("in")
	g.Process("west")

	// Verify we're in living room with rug
	result := g.Process("look")
	if !strings.Contains(result, "Living Room") {
		t.Errorf("Expected to be in Living Room, got: %s", result)
	}
	if !strings.Contains(result, "oriental rug") {
		t.Errorf("Expected to see rug in room, got: %s", result)
	}

	// Try to go down without moving rug - should fail
	result = g.Process("down")
	if !strings.Contains(result, "can't go that way") {
		t.Errorf("Expected blocked passage, got: %s", result)
	}

	// Move the rug
	result = g.Process("move rug")
	if !strings.Contains(result, "trap door") {
		t.Errorf("Expected trap door to be revealed, got: %s", result)
	}
	if !g.Flags["trap-door-open"] {
		t.Error("Expected trap-door-open flag to be set")
	}

	// Look again - should see trap door
	result = g.Process("look")
	if !strings.Contains(result, "trap door") {
		t.Errorf("Expected to see trap door in room, got: %s", result)
	}

	// Open the trap door
	result = g.Process("open trap door")
	if !strings.Contains(result, "staircase") {
		t.Errorf("Expected staircase description, got: %s", result)
	}

	// Now go down - should work
	result = g.Process("down")
	if strings.Contains(result, "can't go that way") {
		t.Errorf("Expected to go down, got: %s", result)
	}

	// Should be in cellar (dark without lamp)
	if !strings.Contains(result, "pitch black") && !strings.Contains(result, "Cellar") {
		t.Errorf("Expected to be in cellar, got: %s", result)
	}
}

// TestGratingPuzzle tests unlocking the grating with keys
func TestGratingPuzzle(t *testing.T) {
	g := NewGameV2()

	// Navigate to living room and get keys
	g.Process("north")
	g.Process("east")
	g.Process("open window")
	g.Process("in")
	g.Process("west")

	result := g.Process("take keys")
	if !strings.Contains(result, "Taken") {
		t.Errorf("Expected to take keys, got: %s", result)
	}

	// Navigate to grating-clearing
	g.Process("east")  // to kitchen
	g.Process("out")   // to behind house
	g.Process("north") // to north-of-house
	g.Process("north") // to path
	g.Process("north") // to grating-clearing

	result = g.Process("look")
	if !strings.Contains(result, "Clearing") {
		t.Errorf("Expected to be at clearing, got: %s", result)
	}
	if !strings.Contains(result, "grating") {
		t.Errorf("Expected to see grating, got: %s", result)
	}

	// Try to go down without opening - should fail
	result = g.Process("down")
	if !strings.Contains(result, "closed") {
		t.Errorf("Expected grating to be closed, got: %s", result)
	}

	// Open grating with keys
	result = g.Process("open grating")
	if !strings.Contains(result, "unlocked") {
		t.Errorf("Expected grating to unlock, got: %s", result)
	}
	if !g.Flags["grate-open"] {
		t.Error("Expected grate-open flag to be set")
	}

	// Now go down - should work
	result = g.Process("down")
	if strings.Contains(result, "closed") {
		t.Errorf("Expected to go down, got: %s", result)
	}
}

// TestDamControlsPuzzle tests the dam button mechanics
func TestDamControlsPuzzle(t *testing.T) {
	g := NewGameV2()

	// Set location to maintenance room for testing
	g.Location = "maintenance-room"

	// Initially dam should be closed
	if g.Flags["dam-open"] {
		t.Error("Dam should start closed")
	}
	if g.Flags["low-tide"] {
		t.Error("Reservoir should start full")
	}

	// Push yellow button - opens dam
	result := g.Process("push yellow button")
	if !strings.Contains(result, "rumbling") {
		t.Errorf("Expected dam to open, got: %s", result)
	}
	if !g.Flags["dam-open"] {
		t.Error("Expected dam-open flag to be set")
	}
	if !g.Flags["low-tide"] {
		t.Error("Expected low-tide flag to be set")
	}

	// Push yellow again - should say already open
	result = g.Process("push yellow button")
	if !strings.Contains(result, "already open") {
		t.Errorf("Expected 'already open', got: %s", result)
	}

	// Push blue button - closes dam
	result = g.Process("push blue button")
	if !strings.Contains(result, "rushing") {
		t.Errorf("Expected dam to close, got: %s", result)
	}
	if g.Flags["dam-open"] {
		t.Error("Expected dam-open flag to be cleared")
	}
	if g.Flags["low-tide"] {
		t.Error("Expected low-tide flag to be cleared")
	}

	// Push blue again - should say already closed
	result = g.Process("push blue button")
	if !strings.Contains(result, "already closed") {
		t.Errorf("Expected 'already closed', got: %s", result)
	}
}

// TestMachineBasketPuzzle tests the basket raising/lowering mechanism
func TestMachineBasketPuzzle(t *testing.T) {
	g := NewGameV2()

	// Set location to machine room
	g.Location = "machine-room"

	// Initially basket should be raised (at shaft-room)
	if g.Flags["basket-lowered"] {
		t.Error("Basket should start raised")
	}

	shaftRoom := g.Rooms["shaft-room"]
	if shaftRoom == nil {
		t.Fatal("shaft-room not found")
	}

	// Verify basket is in shaft-room
	hasBasket := false
	for _, itemID := range shaftRoom.Contents {
		if itemID == "raised-basket" {
			hasBasket = true
			break
		}
	}
	if !hasBasket {
		t.Error("Expected raised-basket in shaft-room initially")
	}

	// Push lower button - lowers basket
	result := g.Process("push lower button")
	if !strings.Contains(result, "descends") {
		t.Errorf("Expected basket to descend, got: %s", result)
	}
	if !g.Flags["basket-lowered"] {
		t.Error("Expected basket-lowered flag to be set")
	}

	// Verify basket moved to lower-shaft
	lowerShaft := g.Rooms["lower-shaft"]
	if lowerShaft == nil {
		t.Fatal("lower-shaft not found")
	}

	hasLoweredBasket := false
	for _, itemID := range lowerShaft.Contents {
		if itemID == "lowered-basket" {
			hasLoweredBasket = true
			break
		}
	}
	if !hasLoweredBasket {
		t.Error("Expected lowered-basket in lower-shaft")
	}

	// Push lower again - should say already at bottom
	result = g.Process("push lower button")
	if !strings.Contains(result, "already at the bottom") {
		t.Errorf("Expected 'already at bottom', got: %s", result)
	}

	// Push start button - raises basket
	result = g.Process("push start button")
	if !strings.Contains(result, "ascends") {
		t.Errorf("Expected basket to ascend, got: %s", result)
	}
	if g.Flags["basket-lowered"] {
		t.Error("Expected basket-lowered flag to be cleared")
	}

	// Verify basket back in shaft-room
	hasBasketAgain := false
	for _, itemID := range shaftRoom.Contents {
		if itemID == "raised-basket" {
			hasBasketAgain = true
			break
		}
	}
	if !hasBasketAgain {
		t.Error("Expected raised-basket back in shaft-room")
	}
}

// TestGrueMechanics tests that the grue actually kills the player in darkness
func TestGrueMechanics(t *testing.T) {
	g := NewGameV2()

	// Set location to cellar (dark room) without lamp
	g.Location = "cellar"

	// First look in darkness - should get warning
	result := g.Process("look")
	if !strings.Contains(result, "pitch black") {
		t.Errorf("Expected darkness warning, got: %s", result)
	}
	if !strings.Contains(result, "grue") {
		t.Errorf("Expected grue warning, got: %s", result)
	}

	// Move around in darkness for several turns
	// The grue should eventually attack
	maxTurns := 20
	grueAttacked := false
	for i := 0; i < maxTurns; i++ {
		result = g.Process("look")
		if strings.Contains(result, "slavering fangs") || strings.Contains(result, "You have died") {
			grueAttacked = true
			if !g.GameOver {
				t.Error("Expected GameOver to be true after grue attack")
			}
			break
		}
	}

	if !grueAttacked {
		t.Errorf("Expected grue to attack after %d turns in darkness", maxTurns)
	}
}

// TestLampProtectsFromGrue tests that having a lit lamp prevents grue attacks
func TestLampProtectsFromGrue(t *testing.T) {
	g := NewGameV2()

	// Get the lamp and light it
	lamp := g.Items["lamp"]
	if lamp == nil {
		t.Fatal("Lamp not found")
	}
	lamp.Location = "player-inventory"
	g.Player.Inventory = append(g.Player.Inventory, "lamp")
	lamp.Flags.IsLit = true

	// Set location to cellar (dark room) but with lit lamp
	g.Location = "cellar"

	// Look should work normally with light
	result := g.Process("look")
	if strings.Contains(result, "pitch black") {
		t.Errorf("Should not be dark with lit lamp, got: %s", result)
	}
	if strings.Contains(result, "grue") {
		t.Errorf("Should not warn about grue with lit lamp, got: %s", result)
	}

	// Move around for many turns - grue should never attack
	for i := 0; i < 20; i++ {
		result = g.Process("look")
		if strings.Contains(result, "slavering fangs") || g.GameOver {
			t.Errorf("Grue attacked despite having lit lamp after %d turns", i)
			break
		}
	}

	if g.GameOver {
		t.Error("Game should not be over with lit lamp")
	}
}

// TestTrophyCaseScoring tests that placing treasures in trophy case awards points
func TestTrophyCaseScoring(t *testing.T) {
	g := NewGameV2()

	// Navigate to living room where trophy case is
	g.Process("north")
	g.Process("east")
	g.Process("open window")
	g.Process("in")
	g.Process("west")

	// Verify trophy case is present
	result := g.Process("look")
	if !strings.Contains(result, "trophy case") {
		t.Errorf("Expected to see trophy case, got: %s", result)
	}

	// Get a treasure (use diamond for testing - worth 10 points)
	diamond := g.Items["diamond"]
	if diamond == nil {
		t.Fatal("Diamond not found")
	}
	diamond.Location = "player-inventory"
	g.Player.Inventory = append(g.Player.Inventory, "diamond")

	// Check initial score
	if g.Score != 0 {
		t.Errorf("Initial score should be 0, got: %d", g.Score)
	}

	// Open trophy case
	result = g.Process("open trophy case")
	if !strings.Contains(result, "Opened") {
		t.Errorf("Expected to open trophy case, got: %s", result)
	}

	// Put diamond in trophy case - should award 10 points
	result = g.Process("put diamond in trophy case")
	if !strings.Contains(result, "10 points") {
		t.Errorf("Expected point notification, got: %s", result)
	}

	// Check score increased
	if g.Score != 10 {
		t.Errorf("Expected score to be 10, got: %d", g.Score)
	}

	// Verify scored flag is set
	if !g.Flags["scored-diamond"] {
		t.Error("Expected scored-diamond flag to be set")
	}

	// Try putting same diamond in again (after taking it out) - should NOT award points
	diamond.Location = "player-inventory"
	g.Player.Inventory = append(g.Player.Inventory, "diamond")
	result = g.Process("put diamond in trophy case")
	if strings.Contains(result, "points awarded") {
		t.Errorf("Should not award points twice, got: %s", result)
	}
	if g.Score != 10 {
		t.Errorf("Score should still be 10, got: %d", g.Score)
	}
}

// TestMultipleTreasuresScoring tests scoring with multiple treasures
func TestMultipleTreasuresScoring(t *testing.T) {
	g := NewGameV2()

	// Set up in living room
	g.Location = "living-room"

	// Open trophy case
	trophyCase := g.Items["trophy-case"]
	if trophyCase == nil {
		t.Fatal("Trophy case not found")
	}
	trophyCase.Flags.IsOpen = true

	// Add multiple treasures to inventory
	treasures := []struct {
		id    string
		value int
	}{
		{"diamond", 10},
		{"emerald", 5},
		{"chalice", 10},
	}

	for _, treasure := range treasures {
		item := g.Items[treasure.id]
		if item == nil {
			t.Fatalf("Treasure %s not found", treasure.id)
		}
		item.Location = "player-inventory"
		g.Player.Inventory = append(g.Player.Inventory, treasure.id)
	}

	// Put each treasure in case and verify points
	expectedScore := 0
	for _, treasure := range treasures {
		result := g.Process("put " + treasure.id + " in trophy case")
		expectedScore += treasure.value

		if !strings.Contains(result, fmt.Sprintf("%d points", treasure.value)) {
			t.Errorf("Expected %d points for %s, got: %s", treasure.value, treasure.id, result)
		}

		if g.Score != expectedScore {
			t.Errorf("Expected score %d, got: %d", expectedScore, g.Score)
		}
	}

	// Final score should be 10 + 5 + 10 = 25
	if g.Score != 25 {
		t.Errorf("Final score should be 25, got: %d", g.Score)
	}

	// Check rank
	result := g.Process("score")
	if !strings.Contains(result, "25") {
		t.Errorf("Score command should show 25 points, got: %s", result)
	}
	if !strings.Contains(result, "Amateur Adventurer") {
		t.Errorf("Expected Amateur Adventurer rank for 25 points, got: %s", result)
	}
}
