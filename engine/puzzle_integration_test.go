package engine

import (
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
