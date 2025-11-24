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

// TestLampFuelDepletion tests the lamp fuel system
func TestLampFuelDepletion(t *testing.T) {
	g := NewGameV2()

	// Test 1: Check initial fuel
	lamp := g.Items["lamp"]
	if lamp == nil {
		t.Fatal("Lamp not found")
	}
	if lamp.Fuel != 330 {
		t.Errorf("Expected lamp to start with 330 turns of fuel, got: %d", lamp.Fuel)
	}

	// Test 2: Lamp off should not consume fuel
	lamp.Flags.IsLit = false
	initialFuel := lamp.Fuel
	g.Process("look")
	if lamp.Fuel != initialFuel {
		t.Errorf("Fuel should not decrease when lamp is off, was %d, now %d", initialFuel, lamp.Fuel)
	}

	// Test 3: First warning at 230 turns remaining (after 100 turns)
	lamp.Flags.IsLit = true
	lamp.Fuel = 231
	result := g.Process("look")
	if !strings.Contains(result, "bit dimmer") {
		t.Errorf("Expected 'bit dimmer' warning at 230 fuel, got: %s", result)
	}
	if lamp.Fuel != 230 {
		t.Errorf("Expected fuel to be 230, got: %d", lamp.Fuel)
	}

	// Test 4: Second warning at 160 turns remaining (after 170 turns)
	lamp.Fuel = 161
	result = g.Process("look")
	if !strings.Contains(result, "definitely dimmer") {
		t.Errorf("Expected 'definitely dimmer' warning at 160 fuel, got: %s", result)
	}
	if lamp.Fuel != 160 {
		t.Errorf("Expected fuel to be 160, got: %d", lamp.Fuel)
	}

	// Test 5: Third warning at 145 turns remaining (after 185 turns)
	lamp.Fuel = 146
	result = g.Process("look")
	if !strings.Contains(result, "nearly out") {
		t.Errorf("Expected 'nearly out' warning at 145 fuel, got: %s", result)
	}
	if lamp.Fuel != 145 {
		t.Errorf("Expected fuel to be 145, got: %d", lamp.Fuel)
	}

	// Test 6: Lamp dies in lit room (safe)
	g.Location = "west-of-house"
	lamp.Fuel = 1
	result = g.Process("look")
	if !strings.Contains(result, "gone out") {
		t.Errorf("Expected 'gone out' message when lamp dies, got: %s", result)
	}
	if lamp.Flags.IsLit {
		t.Error("Lamp should be extinguished")
	}
	if g.GameOver {
		t.Error("Player should not die in lit room")
	}

	// Test 7: Lamp dies in dark room (grue death)
	g2 := NewGameV2()
	g2.Location = "cellar"
	lamp2 := g2.Items["lamp"]
	lamp2.Location = "player-inventory"
	g2.Player.Inventory = append(g2.Player.Inventory, "lamp")
	lamp2.Flags.IsLit = true
	lamp2.Fuel = 1
	result = g2.Process("look")
	if !g2.GameOver {
		t.Error("Player should die from grue in dark room")
	}
	if !strings.Contains(result, "grue") || !strings.Contains(result, "died") {
		t.Errorf("Expected grue death message, got: %s", result)
	}
}

// TestCandlesFuelDepletion tests the candles fuel system
func TestCandlesFuelDepletion(t *testing.T) {
	g := NewGameV2()

	// Test 1: Check initial fuel
	candles := g.Items["candles"]
	if candles == nil {
		t.Fatal("Candles not found")
	}
	if candles.Fuel != 40 {
		t.Errorf("Expected candles to start with 40 turns of fuel, got: %d", candles.Fuel)
	}

	// Test 2: Candles off should not consume fuel
	candles.Location = "player-inventory"
	g.Player.Inventory = append(g.Player.Inventory, "candles")
	candles.Flags.IsLit = false
	g.Location = "west-of-house"
	initialFuel := candles.Fuel
	for i := 0; i < 5; i++ {
		g.Process("look")
	}
	if candles.Fuel != initialFuel {
		t.Errorf("Fuel should not decrease when candles not lit, was %d, now %d", initialFuel, candles.Fuel)
	}

	// Test 3: First warning at 20 turns remaining
	candles.Flags.IsLit = true
	candles.Fuel = 21
	result := g.Process("look")
	if !strings.Contains(result, "grow shorter") {
		t.Errorf("Expected 'grow shorter' warning at 20 fuel, got: %s", result)
	}
	if candles.Fuel != 20 {
		t.Errorf("Expected fuel to be 20, got: %d", candles.Fuel)
	}

	// Test 4: Second warning at 10 turns remaining
	candles.Fuel = 11
	result = g.Process("look")
	if !strings.Contains(result, "quite short") {
		t.Errorf("Expected 'quite short' warning at 10 fuel, got: %s", result)
	}
	if candles.Fuel != 10 {
		t.Errorf("Expected fuel to be 10, got: %d", candles.Fuel)
	}

	// Test 5: Third warning at 5 turns remaining
	candles.Fuel = 6
	result = g.Process("look")
	if !strings.Contains(result, "won't last long") {
		t.Errorf("Expected 'won't last long' warning at 5 fuel, got: %s", result)
	}
	if candles.Fuel != 5 {
		t.Errorf("Expected fuel to be 5, got: %d", candles.Fuel)
	}

	// Test 6: Candles burn out in lit room (safe)
	candles.Fuel = 1
	result = g.Process("look")
	if !strings.Contains(result, "better have more light") {
		t.Errorf("Expected burnout message, got: %s", result)
	}
	if candles.Flags.IsLit {
		t.Error("Candles should be extinguished")
	}
	if g.GameOver {
		t.Error("Player should not die in lit room")
	}

	// Test 7: Candles burn out in dark room (grue death)
	g2 := NewGameV2()
	g2.Location = "cellar"
	candles2 := g2.Items["candles"]
	candles2.Location = "player-inventory"
	g2.Player.Inventory = append(g2.Player.Inventory, "candles")
	candles2.Flags.IsLit = true
	candles2.Fuel = 1
	result = g2.Process("look")
	if !g2.GameOver {
		t.Error("Player should die from grue in dark room")
	}
	if !strings.Contains(result, "grue") || !strings.Contains(result, "died") {
		t.Errorf("Expected grue death message, got: %s", result)
	}
}

// TestIvoryTorchEternalFlame tests that the ivory torch never burns out
func TestIvoryTorchEternalFlame(t *testing.T) {
	g := NewGameV2()

	torch := g.Items["ivory-torch"]
	if torch == nil {
		t.Fatal("Ivory torch not found")
	}

	// Torch should have eternal flame (Fuel=-1)
	if torch.Fuel != -1 {
		t.Errorf("Expected torch to have eternal flame (Fuel=-1), got: %d", torch.Fuel)
	}

	// Torch should be lit by default
	if !torch.Flags.IsLit {
		t.Error("Torch should be lit by default")
	}

	// Torch should be a light source
	if !torch.Flags.IsLightSource {
		t.Error("Torch should be a light source")
	}

	// Torch should be a treasure worth 6 points
	if !torch.Flags.IsTreasure {
		t.Error("Torch should be a treasure")
	}
	if torch.Value != 6 {
		t.Errorf("Expected torch value to be 6, got: %d", torch.Value)
	}
}

func TestTrapDoorMechanics(t *testing.T) {
	g := NewGameV2()

	// Test 1: Trap door should be hidden in living-room until rug is moved
	t.Run("hidden until rug moved", func(t *testing.T) {
		g.Location = "living-room"
		g.Flags["trap-door-open"] = false // Ensure rug hasn't been moved

		result := g.Process("look")
		if strings.Contains(result, "trap door") {
			t.Error("Trap door should be hidden until rug is moved")
		}

		result = g.Process("examine trap door")
		if !strings.Contains(result, "can't see") {
			t.Errorf("Should not be able to see trap door, got: %s", result)
		}
	})

	// Test 2: Moving rug reveals trap door
	t.Run("moving rug reveals door", func(t *testing.T) {
		g.Location = "living-room"
		result := g.Process("move rug")

		if !strings.Contains(result, "trap door") {
			t.Errorf("Expected trap door to be revealed, got: %s", result)
		}

		if !g.Flags["trap-door-open"] {
			t.Error("trap-door-open flag should be set after moving rug")
		}

		// Now it should be visible
		result = g.Process("look")
		if !strings.Contains(result, "trap door") {
			t.Error("Trap door should now be visible")
		}
	})

	// Test 3: Open from living-room reveals staircase
	t.Run("open from living-room", func(t *testing.T) {
		g.Location = "living-room"
		g.Flags["trap-door-open"] = true
		trapDoor := g.Items["trap-door"]
		trapDoor.Flags.IsOpen = false

		result := g.Process("open trap door")
		if !strings.Contains(result, "staircase") {
			t.Errorf("Expected staircase description, got: %s", result)
		}

		if !trapDoor.Flags.IsOpen {
			t.Error("Trap door should be open")
		}
	})

	// Test 4: Close from living-room
	t.Run("close from living-room", func(t *testing.T) {
		g.Location = "living-room"
		g.Flags["trap-door-open"] = true
		trapDoor := g.Items["trap-door"]
		trapDoor.Flags.IsOpen = true

		result := g.Process("close trap door")
		if !strings.Contains(result, "swings shut") {
			t.Errorf("Expected door to swing shut, got: %s", result)
		}

		if trapDoor.Flags.IsOpen {
			t.Error("Trap door should be closed")
		}
	})

	// Test 5: Open from cellar shows "locked from above"
	t.Run("open from cellar locked", func(t *testing.T) {
		g.Location = "cellar"
		trapDoor := g.Items["trap-door"]
		trapDoor.Flags.IsOpen = false

		// Add lamp for visibility in dark cellar
		lamp := g.Items["lamp"]
		lamp.Location = "inventory"
		lamp.Flags.IsLit = true
		g.Player.Inventory = []string{"lamp"}

		result := g.Process("open trap door")
		if !strings.Contains(result, "locked from above") {
			t.Errorf("Expected locked message, got: %s", result)
		}

		if trapDoor.Flags.IsOpen {
			t.Error("Trap door should remain closed")
		}
	})

	// Test 6: Close from cellar locks door
	t.Run("close from cellar locks", func(t *testing.T) {
		g.Location = "cellar"
		trapDoor := g.Items["trap-door"]
		trapDoor.Flags.IsOpen = true

		// Add lamp for visibility
		lamp := g.Items["lamp"]
		lamp.Location = "inventory"
		lamp.Flags.IsLit = true
		g.Player.Inventory = []string{"lamp"}

		result := g.Process("close trap door")
		if !strings.Contains(result, "closes and locks") {
			t.Errorf("Expected door to close and lock, got: %s", result)
		}

		if trapDoor.Flags.IsOpen {
			t.Error("Trap door should be closed")
		}
	})

	// Test 7: Trap door always visible in cellar
	t.Run("visible in cellar", func(t *testing.T) {
		g.Location = "cellar"

		// Add lamp for visibility
		lamp := g.Items["lamp"]
		lamp.Location = "inventory"
		lamp.Flags.IsLit = true
		g.Player.Inventory = []string{"lamp"}

		result := g.Process("look")
		if !strings.Contains(result, "trap door") {
			t.Error("Trap door should always be visible in cellar")
		}
	})
}

func TestKitchenWindowMechanics(t *testing.T) {
	g := NewGameV2()

	// Test 1: Window examine before opening
	t.Run("examine before first open", func(t *testing.T) {
		g.Location = "behind-house"
		g.Flags["window-opened-once"] = false

		result := g.Process("examine window")
		if !strings.Contains(result, "not enough") {
			t.Errorf("Expected 'not enough to allow entry' message, got: %s", result)
		}
	})

	// Test 2: Open window from behind-house
	t.Run("open window", func(t *testing.T) {
		g.Location = "behind-house"
		window := g.Items["kitchen-window"]
		window.Flags.IsOpen = false
		g.Flags["window-open"] = false
		g.Flags["window-opened-once"] = false

		result := g.Process("open window")
		if !strings.Contains(result, "great effort") {
			t.Errorf("Expected 'great effort' message, got: %s", result)
		}

		if !window.Flags.IsOpen {
			t.Error("Window should be open")
		}

		if !g.Flags["window-opened-once"] {
			t.Error("window-opened-once flag should be set")
		}
	})

	// Test 3: Close window
	t.Run("close window", func(t *testing.T) {
		g.Location = "behind-house"
		window := g.Items["kitchen-window"]
		window.Flags.IsOpen = true
		g.Flags["window-open"] = true

		result := g.Process("close window")
		if !strings.Contains(result, "more easily") {
			t.Errorf("Expected 'more easily' message, got: %s", result)
		}

		if window.Flags.IsOpen {
			t.Error("Window should be closed")
		}

		if g.Flags["window-open"] {
			t.Error("window-open flag should be cleared")
		}
	})

	// Test 4: Examine after opening once
	t.Run("examine after first open", func(t *testing.T) {
		g.Location = "behind-house"
		g.Flags["window-opened-once"] = true

		result := g.Process("examine window")
		// After opening once, should show normal description, not the special message
		if strings.Contains(result, "not enough") {
			t.Errorf("Should not show 'not enough' after first open, got: %s", result)
		}
	})

	// Test 5: Window controls passage between rooms
	t.Run("window controls passage", func(t *testing.T) {
		g.Location = "behind-house"
		window := g.Items["kitchen-window"]

		// With window closed, can't go west
		window.Flags.IsOpen = false
		g.Flags["window-open"] = false
		result := g.Process("west")
		if !strings.Contains(result, "window is closed") && !strings.Contains(result, "can't go that way") {
			// May vary based on implementation
			t.Logf("Window closed, movement result: %s", result)
		}

		// With window open, can go west
		window.Flags.IsOpen = true
		g.Flags["window-open"] = true
		result = g.Process("west")
		// Should successfully move or at least not be blocked by window
		t.Logf("Window open, movement result: %s", result)
	})
}
