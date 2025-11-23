package engine

import (
	"fmt"
	"strings"
)

// GameV2 represents the refactored game state with proper types
type GameV2 struct {
	Rooms     map[string]*Room
	Items     map[string]*Item
	NPCs      map[string]*NPC
	Parser    *Parser
	Player    *Player
	Location  string // Current room ID
	Score     int
	Moves     int
	Flags     map[string]bool // Global game flags (WINDOW-OPEN, TROLL-DEAD, etc.)
	GameOver  bool
	Won       bool
}

// Player represents the player character
type Player struct {
	Inventory []string // Item IDs
	MaxWeight int
	Health    int
}

// NewGameV2 creates a new game with proper type separation
func NewGameV2() *GameV2 {
	g := &GameV2{
		Rooms:  make(map[string]*Room),
		Items:  make(map[string]*Item),
		NPCs:   make(map[string]*NPC),
		Parser: NewParser(),
		Player: &Player{
			Inventory: []string{},
			MaxWeight: 100,
			Health:    100,
		},
		Flags: make(map[string]bool),
	}

	// Initialize world
	g.initializeWorld()

	return g
}

// initializeWorld sets up the initial game state
func (g *GameV2) initializeWorld() {
	// Create all 110 rooms from original Zork I
	InitializeRooms(g)
	// Create all items from original Zork I
	InitializeItems(g)
	// Create NPCs
	g.createNPCs()

	// Set starting location
	g.Location = "west-of-house"
}

func (g *GameV2) createNPCs() {
	// The Troll - Blocks passages, guards treasure
	troll := NewNPC(
		"troll",
		"nasty troll",
		"A nasty-looking troll, brandishing a bloody axe, blocks all passages out of the room.",
	)
	troll.Location = "troll-room"
	troll.Strength = 20
	troll.Weapon = "axe"
	troll.Hostile = true
	troll.Flags.IsAggressive = true
	troll.Flags.CanFight = true
	g.NPCs["troll"] = troll
	g.Rooms["troll-room"].AddNPC("troll")

	// The Thief - Steals treasures, moves around dungeon
	thief := NewNPC(
		"thief",
		"shady thief",
		"A suspicious-looking individual with a bag of stolen goods eyes you warily.",
	)
	thief.Location = "maze-1" // Starts in maze
	thief.Strength = 15
	thief.Weapon = "stiletto"
	thief.Hostile = false // Only hostile if attacked or has loot
	thief.Flags.IsAggressive = false
	thief.Flags.CanFight = true
	thief.Inventory = []string{} // Will steal treasures
	g.NPCs["thief"] = thief
	g.Rooms["maze-1"].AddNPC("thief")

	// The Cyclops - Guards treasure room
	cyclops := NewNPC(
		"cyclops",
		"cyclops",
		"A cyclops, who looks prepared to eat you, blocks the way. He seems to have been eating, for on the ground is a lot of refuse.",
	)
	cyclops.Location = "cyclops-room"
	cyclops.Strength = 30
	cyclops.Weapon = "" // Uses fists
	cyclops.Hostile = true
	cyclops.Flags.IsAggressive = true
	cyclops.Flags.CanFight = true
	g.NPCs["cyclops"] = cyclops
	g.Rooms["cyclops-room"].AddNPC("cyclops")

	// The Bat - Carries player to random dark rooms
	bat := NewNPC(
		"bat",
		"vampire bat",
		"A vampire bat is circling overhead, its beady eyes fixed on you.",
	)
	bat.Location = "bat-room"
	bat.Strength = 5
	bat.Weapon = ""
	bat.Hostile = false // Not directly hostile, just annoying
	bat.Flags.IsAggressive = true // Will grab you
	bat.Flags.CanFight = false    // Can't be fought effectively
	g.NPCs["bat"] = bat
	g.Rooms["bat-room"].AddNPC("bat")

	// The Ghosts - Temple spirits, need exorcism
	ghosts := NewNPC(
		"ghosts",
		"evil spirits",
		"The room is filled with evil spirits. They are making a significant racket.",
	)
	ghosts.Location = "entrance-to-hades"
	ghosts.Strength = 0 // Can't be fought
	ghosts.Weapon = ""
	ghosts.Hostile = false
	ghosts.Flags.IsAggressive = false
	ghosts.Flags.CanFight = false // Need exorcism, not combat
	g.NPCs["ghosts"] = ghosts
	g.Rooms["entrance-to-hades"].AddNPC("ghosts")
}

// Process handles a command - same interface as before
func (g *GameV2) Process(input string) string {
	if g.GameOver {
		return "The game is over. Type RESTART to play again."
	}

	if strings.TrimSpace(input) == "" {
		return ""
	}

	cmd, err := g.Parser.Parse(input)
	if err != nil {
		return err.Error()
	}

	return g.executeCommand(cmd)
}

func (g *GameV2) executeCommand(cmd *Command) string {
	g.Moves++

	var result string

	// Handle movement
	if cmd.Verb == "walk" && cmd.Direction != "" {
		result = g.handleMove(cmd.Direction)
	} else {
		// Handle other verbs
		switch cmd.Verb {
		case "look":
			if cmd.Preposition == "in" || cmd.Preposition == "into" {
				result = g.handleLookIn(cmd.IndirectObject)
			} else {
				result = g.handleLook()
			}
		case "examine":
			result = g.handleExamine(cmd.DirectObject)
		case "take":
			result = g.handleTake(cmd.DirectObject)
		case "drop":
			result = g.handleDrop(cmd.DirectObject)
		case "open":
			result = g.handleOpen(cmd.DirectObject)
		case "close":
			result = g.handleClose(cmd.DirectObject)
		case "read":
			result = g.handleRead(cmd.DirectObject)
		case "turn":
			// Handle both "turn on lamp" and "turn lamp on"
			objName := cmd.DirectObject
			if objName == "" {
				objName = cmd.IndirectObject
			}

			if cmd.Preposition == "on" {
				result = g.handleTurnOn(objName)
			} else if cmd.Preposition == "off" {
				result = g.handleTurnOff(objName)
			} else {
				result = "Turn it on or off?"
			}
		case "inventory":
			result = g.handleInventory()
		case "help":
			result = g.handleHelp()
		case "put":
			result = g.handlePut(cmd.DirectObject, cmd.Preposition, cmd.IndirectObject)
		case "give":
			result = g.handleGive(cmd.DirectObject, cmd.IndirectObject)
		case "attack":
			result = g.handleAttack(cmd.DirectObject)
		case "wave":
			result = g.handleWave(cmd.DirectObject)
		case "climb":
			result = g.handleClimb(cmd.DirectObject)
		case "tie":
			result = g.handleTie(cmd.DirectObject, cmd.IndirectObject)
		case "untie":
			result = g.handleUntie(cmd.DirectObject)
		case "dig":
			result = g.handleDig(cmd.DirectObject)
		case "push":
			result = g.handlePush(cmd.DirectObject)
		case "pull":
			result = g.handlePull(cmd.DirectObject)
		case "move":
			result = g.handleMoveObject(cmd.DirectObject)
		case "ring":
			result = g.handleRing(cmd.DirectObject)
		case "pray":
			result = g.handlePray()
		case "wait":
			result = g.handleWait()
		case "eat":
			result = g.handleEat(cmd.DirectObject)
		case "drink":
			result = g.handleDrink(cmd.DirectObject)
		case "fill":
			result = g.handleFill(cmd.DirectObject, cmd.IndirectObject)
		case "inflate":
			result = g.handleInflate(cmd.DirectObject, cmd.IndirectObject)
		case "deflate":
			result = g.handleDeflate(cmd.DirectObject)
		case "plug":
			result = g.handlePlug(cmd.DirectObject, cmd.IndirectObject)
		case "pour":
			result = g.handlePour(cmd.DirectObject, cmd.IndirectObject)
		case "listen":
			result = g.handleListen()
		case "smell":
			result = g.handleSmell(cmd.DirectObject)
		case "touch":
			result = g.handleTouch(cmd.DirectObject)
		case "search":
			result = g.handleSearch(cmd.DirectObject)
		case "jump":
			result = g.handleJump()
		case "swim":
			result = g.handleSwim()
		case "blow":
			result = g.handleBlow(cmd.DirectObject)
		case "knock":
			// Handle both "knock door" and "knock on door"
			objName := cmd.DirectObject
			if objName == "" && cmd.Preposition == "on" {
				objName = cmd.IndirectObject
			}
			result = g.handleKnock(objName)
		case "quit":
			g.GameOver = true
			return "Thanks for playing!"
		case "save":
			result = "Save is not yet implemented."
		case "restore":
			result = "Restore is not yet implemented."
		case "score":
			result = g.handleScore()
		case "restart":
			result = "Restart is not yet implemented."
		default:
			result = "I don't understand how to \"" + cmd.Verb + "\" something."
		}
	}

	// Process NPC turns after every command (including grues!)
	npcResult := g.processNPCTurns()
	if npcResult != "" {
		result += "\n\n" + npcResult
	}

	return result
}

// processNPCTurns handles NPC behaviors each turn
func (g *GameV2) processNPCTurns() string {
	var result strings.Builder

	// Process grue attacks (must be first - can end the game!)
	grueResult := g.processGrueBehavior()
	if grueResult != "" {
		result.WriteString(grueResult)
		// If grue killed player, don't process other NPCs
		if g.GameOver {
			return strings.TrimSpace(result.String())
		}
	}

	// Process thief roaming and stealing
	thiefResult := g.processThiefBehavior()
	if thiefResult != "" {
		if result.Len() > 0 {
			result.WriteString("\n")
		}
		result.WriteString(thiefResult)
	}

	// Process bat behavior
	batResult := g.processBatBehavior()
	if batResult != "" {
		if result.Len() > 0 {
			result.WriteString("\n")
		}
		result.WriteString(batResult)
	}

	return strings.TrimSpace(result.String())
}

// processThiefBehavior handles thief roaming and treasure stealing
func (g *GameV2) processThiefBehavior() string {
	thief := g.NPCs["thief"]
	if thief == nil || !thief.Flags.IsAlive {
		return ""
	}

	// Thief only acts every few turns (to not be too annoying)
	if g.Moves%3 != 0 {
		return ""
	}

	thiefRoom := g.Rooms[thief.Location]
	if thiefRoom == nil {
		return ""
	}

	// If thief is in same room as player, try to steal treasures
	if thief.Location == g.Location {
		// Look for treasures in room or player inventory
		for _, itemID := range g.Player.Inventory {
			item := g.Items[itemID]
			if item != nil && item.Flags.IsTreasure {
				// Thief steals the treasure!
				// Remove from player inventory
				for i, id := range g.Player.Inventory {
					if id == itemID {
						g.Player.Inventory = append(g.Player.Inventory[:i], g.Player.Inventory[i+1:]...)
						break
					}
				}
				// Add to thief's inventory
				thief.Inventory = append(thief.Inventory, itemID)
				item.Location = "thief-inventory"
				return "The thief steals your " + item.Name + " and runs away laughing!"
			}
		}
	}

	// Thief moves to a random connected room
	exits := thiefRoom.Exits
	if len(exits) > 0 {
		// Pick a random exit from the map
		exitIndex := 0
		targetIndex := g.Moves % len(exits)
		var newLocation string
		for _, exit := range exits {
			if exitIndex == targetIndex {
				newLocation = exit.To
				break
			}
			exitIndex++
		}

		if newLocation != "" {
			// Move thief
			thiefRoom.RemoveNPC("thief")
			thief.Location = newLocation
			newRoom := g.Rooms[newLocation]
			if newRoom != nil {
				newRoom.AddNPC("thief")
			}

			// If player can hear the thief, mention it
			if thief.Location == g.Location {
				return "The thief appears from the shadows!"
			}
		}
	}

	return ""
}

// processBatBehavior handles bat grabbing and moving player
func (g *GameV2) processBatBehavior() string {
	bat := g.NPCs["bat"]
	if bat == nil || !bat.Flags.IsAlive {
		return ""
	}

	// Bat only acts occasionally
	if g.Moves%5 != 0 {
		return ""
	}

	// If bat is in same room as player, it might grab them
	if bat.Location == g.Location {
		// 50% chance to grab player (based on turn number)
		if g.Moves%10 == 0 {
			// Move player to a random adjacent room
			currentRoom := g.Rooms[g.Location]
			if currentRoom != nil && len(currentRoom.Exits) > 0 {
				// Pick a random exit from the map
				exitIndex := 0
				targetIndex := g.Moves % len(currentRoom.Exits)
				var newLocation string
				for _, exit := range currentRoom.Exits {
					if exitIndex == targetIndex {
						newLocation = exit.To
						break
					}
					exitIndex++
				}

				if newLocation != "" {
					g.Location = newLocation
					return "A large vampire bat swoops down, grabs you, and carries you off!\n\n" + g.handleLook()
				}
			}
		}
	}

	return ""
}

// processGrueBehavior handles grue attacks in the dark
func (g *GameV2) processGrueBehavior() string {
	// Check if player is in darkness without light
	room := g.Rooms[g.Location]
	if room == nil {
		return ""
	}

	// If room is dark and player has no light source
	if room.Flags.IsDark && !g.hasLight() {
		// Track consecutive turns in darkness
		if !g.Flags["in-darkness"] {
			// First turn in darkness - just a warning
			g.Flags["in-darkness"] = true
			g.Flags["darkness-turns"] = true // Using flag as counter start
			return "" // Warning already shown by handleLook/handleMove
		}

		// Player has been in darkness for multiple turns
		// In original Zork, grue attacks after 2-3 turns in darkness
		// We'll use turn count modulo to create randomness
		turnMod := g.Moves % 4
		if turnMod == 0 {
			// Grue attacks!
			g.GameOver = true
			return "\nOh, no! You have walked into the slavering fangs of a lurking grue!\n\n****  You have died  ****"
		}

		// Additional warnings
		warningMod := g.Moves % 3
		switch warningMod {
		case 0:
			return "You hear a horrible slavering sound in the darkness nearby..."
		case 1:
			return "The grue is getting closer! You can feel its hot, fetid breath..."
		default:
			return "Something is moving in the darkness. You should find light quickly!"
		}
	} else {
		// Player has light or is in lit room - clear darkness tracking
		if g.Flags["in-darkness"] {
			g.Flags["in-darkness"] = false
			delete(g.Flags, "darkness-turns")
		}
	}

	return ""
}

func (g *GameV2) handleMove(direction string) string {
	currentRoom := g.Rooms[g.Location]
	if currentRoom == nil {
		return "You are nowhere!"
	}

	exit := currentRoom.Exits[direction]
	if exit == nil {
		return "You can't go that way."
	}

	// Check condition if present
	if exit.Condition != "" && !g.Flags[exit.Condition] {
		if exit.Message != "" {
			return exit.Message
		}
		return "You can't go that way."
	}

	// Check if destination room exists
	destRoom := g.Rooms[exit.To]
	if destRoom == nil {
		return "You can't go that way."
	}

	// Check if room is dark and player has no light
	if destRoom.Flags.IsDark && !g.hasLight() {
		return "It is pitch black. You are likely to be eaten by a grue."
	}

	// Move player
	g.Location = exit.To
	destRoom.FirstVisit = false

	// Auto-look at new room
	return g.handleLook()
}

func (g *GameV2) handleLook() string {
	room := g.Rooms[g.Location]
	if room == nil {
		return "You are nowhere!"
	}

	// Check for darkness
	if room.Flags.IsDark && !g.hasLight() {
		return "It is pitch black. You are likely to be eaten by a grue."
	}

	var result strings.Builder
	result.WriteString(room.Name + "\n")
	result.WriteString(room.Description + "\n")

	// List items in room
	for _, itemID := range room.Contents {
		item := g.Items[itemID]
		if item != nil {
			result.WriteString("There is a " + item.Name + " here.\n")
		}
	}

	// List NPCs in room
	for _, npcID := range room.NPCs {
		npc := g.NPCs[npcID]
		if npc != nil && npc.Flags.IsAlive {
			result.WriteString(npc.Description + "\n")
		}
	}

	return strings.TrimSpace(result.String())
}

func (g *GameV2) handleExamine(objName string) string {
	if objName == "" {
		return "What do you want to examine?"
	}

	// Try to find item
	item := g.findItem(objName)
	if item != nil {
		result := item.Description

		// If it's a container, show if it's open/closed
		if item.Flags.IsContainer {
			if item.Flags.IsOpen {
				result += " It is open."
			} else {
				result += " It is closed."
			}

			// Show what's inside if open or transparent
			if item.Flags.IsOpen || item.Flags.IsTransparent {
				var contents []string
				for _, otherItem := range g.Items {
					if otherItem.Location == item.ID {
						contents = append(contents, otherItem.Name)
					}
				}

				if len(contents) > 0 {
					result += "\n\nThe " + item.Name + " contains:\n"
					for _, itemName := range contents {
						result += "  A " + itemName + "\n"
					}
				}
			}
		}

		return strings.TrimSpace(result)
	}

	// Try to find NPC
	npc := g.findNPC(objName)
	if npc != nil {
		return npc.Description
	}

	return "You can't see any " + objName + " here."
}

func (g *GameV2) handleTake(objName string) string {
	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	if !item.Flags.IsTakeable {
		return "You can't take the " + item.Name + "."
	}

	// Special case: Taking the rug reveals the trap door
	if item.ID == "rug" && g.Location == "living-room" {
		g.Flags["trap-door-open"] = true
		trapDoor := g.Items["trap-door"]
		if trapDoor != nil {
			trapDoor.Location = "living-room"
			room := g.Rooms["living-room"]
			if room != nil {
				room.AddItem("trap-door")
			}
		}
	}

	// Remove from room
	room := g.Rooms[g.Location]
	if room != nil {
		room.RemoveItem(item.ID)
	}

	// Add to inventory
	item.Location = "inventory"
	g.Player.Inventory = append(g.Player.Inventory, item.ID)

	if item.ID == "rug" && g.Location == "living-room" {
		return "Taken.\nWith the rug moved aside, you can see a closed trap door beneath it!"
	}

	return "Taken."
}

func (g *GameV2) handleDrop(objName string) string {
	item := g.findItemInInventory(objName)
	if item == nil {
		return "You don't have that."
	}

	// Remove from inventory
	for i, id := range g.Player.Inventory {
		if id == item.ID {
			g.Player.Inventory = append(g.Player.Inventory[:i], g.Player.Inventory[i+1:]...)
			break
		}
	}

	// Add to current room
	item.Location = g.Location
	room := g.Rooms[g.Location]
	if room != nil {
		room.AddItem(item.ID)
	}

	return "Dropped."
}

func (g *GameV2) handleInventory() string {
	if len(g.Player.Inventory) == 0 {
		return "You are empty-handed."
	}

	var result strings.Builder
	result.WriteString("You are carrying:\n")
	for _, itemID := range g.Player.Inventory {
		item := g.Items[itemID]
		if item != nil {
			result.WriteString("  A " + item.Name + "\n")
		}
	}
	return strings.TrimSpace(result.String())
}

func (g *GameV2) handleOpen(objName string) string {
	if objName == "" {
		return "What do you want to open?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	// Special handling for nailed door in living room
	if item.ID == "door" {
		return "The door is solidly nailed shut and cannot be opened."
	}

	// Special handling for kitchen window
	if item.ID == "kitchen-window" {
		if g.Flags["window-open"] {
			return "It's already open."
		}
		g.Flags["window-open"] = true
		item.Flags.IsOpen = true
		return "With great effort, you open the window far enough to allow entry."
	}

	// Special handling for trap door
	if item.ID == "trap-door" {
		if !g.Flags["trap-door-open"] {
			return "The rug must be moved first."
		}
		if item.Flags.IsOpen {
			return "It's already open."
		}
		item.Flags.IsOpen = true
		return "The door reluctantly opens to reveal a rickety staircase descending into darkness."
	}

	// Special handling for grating
	if item.ID == "grating" || item.ID == "grate" {
		// Check if player has keys
		hasKeys := false
		for _, itemID := range g.Player.Inventory {
			if itemID == "keys" {
				hasKeys = true
				break
			}
		}
		if !hasKeys {
			return "The grating is locked."
		}
		if g.Flags["grate-open"] {
			return "It's already open."
		}
		g.Flags["grate-open"] = true
		item.Flags.IsOpen = true
		return "The grating is unlocked and opens easily."
	}

	if !item.Flags.IsContainer {
		return "You can't open that."
	}

	if item.Flags.IsOpen {
		return "It's already open."
	}

	item.Flags.IsOpen = true
	return "Opened."
}

func (g *GameV2) handleClose(objName string) string {
	if objName == "" {
		return "What do you want to close?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	// Special handling for kitchen window
	if item.ID == "kitchen-window" {
		if !g.Flags["window-open"] {
			return "It's already closed."
		}
		g.Flags["window-open"] = false
		item.Flags.IsOpen = false
		return "You close the window."
	}

	if !item.Flags.IsContainer {
		return "You can't close that."
	}

	if !item.Flags.IsOpen {
		return "It's already closed."
	}

	item.Flags.IsOpen = false
	return "Closed."
}

func (g *GameV2) handleRead(objName string) string {
	if objName == "" {
		return "What do you want to read?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	if !item.Flags.IsReadable {
		return "How does one read a " + item.Name + "?"
	}

	return item.Description
}

func (g *GameV2) handleLookIn(objName string) string {
	if objName == "" {
		return "Look in what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	if !item.Flags.IsContainer {
		return "You can't look inside that."
	}

	if !item.Flags.IsOpen && !item.Flags.IsTransparent {
		return "You can't see inside the closed " + item.Name + "."
	}

	// Find items inside this container
	var contents []string
	for _, otherItem := range g.Items {
		if otherItem.Location == item.ID {
			contents = append(contents, otherItem.Name)
		}
	}

	if len(contents) == 0 {
		return "The " + item.Name + " is empty."
	}

	var result strings.Builder
	result.WriteString("The " + item.Name + " contains:\n")
	for _, itemName := range contents {
		result.WriteString("  A " + itemName + "\n")
	}
	return strings.TrimSpace(result.String())
}

func (g *GameV2) handleTurnOn(objName string) string {
	if objName == "" {
		return "What do you want to turn on?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	if !item.Flags.IsLightSource {
		return "You can't turn that on."
	}

	if item.Flags.IsLit {
		return "It's already on."
	}

	item.Flags.IsLit = true
	return "The " + item.Name + " is now on."
}

func (g *GameV2) handleTurnOff(objName string) string {
	if objName == "" {
		return "What do you want to turn off?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	if !item.Flags.IsLightSource {
		return "You can't turn that off."
	}

	if !item.Flags.IsLit {
		return "It's already off."
	}

	item.Flags.IsLit = false
	return "The " + item.Name + " is now off."
}

func (g *GameV2) handleHelp() string {
	room := g.Rooms[g.Location]
	if room == nil {
		return "You are nowhere!"
	}

	var result strings.Builder
	result.WriteString("Available commands:\n\n")
	result.WriteString("Movement: NORTH, SOUTH, EAST, WEST, UP, DOWN, IN, OUT, etc.\n")
	result.WriteString("Actions: TAKE, DROP, OPEN, CLOSE, READ, EXAMINE, LOOK, INVENTORY\n")
	result.WriteString("Light: TURN ON, TURN OFF\n")
	result.WriteString("Other: HELP, QUIT\n\n")

	// Show available exits
	result.WriteString("Obvious exits from here:\n")
	if len(room.Exits) == 0 {
		result.WriteString("  None!\n")
	} else {
		for dir := range room.Exits {
			result.WriteString("  " + strings.ToUpper(dir) + "\n")
		}
	}

	return result.String()
}

// Helper methods

func (g *GameV2) findItem(name string) *Item {
	// Check current room
	room := g.Rooms[g.Location]
	if room != nil {
		for _, itemID := range room.Contents {
			item := g.Items[itemID]
			if item != nil && item.HasAlias(name) {
				return item
			}

			// Check inside containers in the room
			if item != nil && item.Flags.IsContainer && (item.Flags.IsOpen || item.Flags.IsTransparent) {
				for _, otherItem := range g.Items {
					if otherItem.Location == item.ID && otherItem.HasAlias(name) {
						return otherItem
					}
				}
			}
		}
	}

	// Check inventory
	inventoryItem := g.findItemInInventory(name)
	if inventoryItem != nil {
		return inventoryItem
	}

	// Check inside containers in inventory
	for _, itemID := range g.Player.Inventory {
		item := g.Items[itemID]
		if item != nil && item.Flags.IsContainer && (item.Flags.IsOpen || item.Flags.IsTransparent) {
			for _, otherItem := range g.Items {
				if otherItem.Location == item.ID && otherItem.HasAlias(name) {
					return otherItem
				}
			}
		}
	}

	return nil
}

func (g *GameV2) findItemInInventory(name string) *Item {
	for _, itemID := range g.Player.Inventory {
		item := g.Items[itemID]
		if item != nil && item.HasAlias(name) {
			return item
		}
	}
	return nil
}

func (g *GameV2) findNPC(name string) *NPC {
	room := g.Rooms[g.Location]
	if room == nil {
		return nil
	}

	for _, npcID := range room.NPCs {
		npc := g.NPCs[npcID]
		if npc != nil && (npc.ID == name || npc.Name == name || strings.Contains(npc.Name, name)) {
			return npc
		}
	}
	return nil
}

func (g *GameV2) hasLight() bool {
	// Check if current room is lit
	room := g.Rooms[g.Location]
	if room != nil && room.Flags.IsLit {
		return true
	}

	// Check for light source in inventory
	for _, itemID := range g.Player.Inventory {
		item := g.Items[itemID]
		if item != nil && item.Flags.IsLightSource && item.Flags.IsLit {
			return true
		}
	}

	return false
}

// handlePut places an item in/on a container (V-PUT in ZIL)
func (g *GameV2) handlePut(objName string, prep string, containerName string) string {
	if objName == "" {
		return "What do you want to put?"
	}
	if containerName == "" {
		return "Where do you want to put it?"
	}

	// Find the item in inventory
	item := g.findItemInInventory(objName)
	if item == nil {
		return "You don't have that."
	}

	// Find the container
	container := g.findItem(containerName)
	if container == nil {
		return "You can't see any " + containerName + " here."
	}

	if !container.Flags.IsContainer {
		return "You can't put things in the " + container.Name + "."
	}

	if !container.Flags.IsOpen && !container.Flags.IsTransparent {
		return "The " + container.Name + " is closed."
	}

	// Remove from inventory
	for i, id := range g.Player.Inventory {
		if id == item.ID {
			g.Player.Inventory = append(g.Player.Inventory[:i], g.Player.Inventory[i+1:]...)
			break
		}
	}

	// Add to container
	item.Location = container.ID

	// Special case: Putting treasure in trophy case awards points
	if container.ID == "trophy-case" && item.Flags.IsTreasure {
		// Check if this treasure has already been scored
		scoreFlag := "scored-" + item.ID
		if !g.Flags[scoreFlag] {
			g.Flags[scoreFlag] = true
			g.Score += item.Value
			return fmt.Sprintf("Done. (%d points awarded)", item.Value)
		}
	}

	return "Done."
}

// handleGive gives an item to an NPC (V-GIVE in ZIL)
func (g *GameV2) handleGive(objName string, npcName string) string {
	if objName == "" {
		return "What do you want to give?"
	}
	if npcName == "" {
		return "Give it to whom?"
	}

	item := g.findItemInInventory(objName)
	if item == nil {
		return "You don't have that."
	}

	npc := g.findNPC(npcName)
	if npc == nil {
		return "There is no " + npcName + " here."
	}

	// Remove from inventory first
	for i, id := range g.Player.Inventory {
		if id == item.ID {
			g.Player.Inventory = append(g.Player.Inventory[:i], g.Player.Inventory[i+1:]...)
			break
		}
	}

	// Special cases per NPC
	switch npc.ID {
	case "troll":
		if item.ID == "lunch" || item.ID == "garlic" {
			// Troll accepts food and leaves
			g.Flags["troll-satisfied"] = true
			room := g.Rooms[npc.Location]
			if room != nil {
				room.RemoveNPC("troll")
			}
			delete(g.Items, item.ID) // Food is consumed
			return "The troll grabs the " + item.Name + " and devours it. He looks satisfied and wanders off, leaving the path clear."
		}

	case "cyclops":
		if item.ID == "lunch" {
			// Cyclops eats lunch and becomes less hostile
			g.Flags["cyclops-fed"] = true
			npc.Flags.IsAggressive = false
			delete(g.Items, item.ID)
			return "The cyclops eagerly grabs the lunch and stuffs it into his mouth. He seems less inclined to eat you now."
		}

	case "thief":
		// Thief steals valuable items
		if item.Flags.IsTreasure {
			npc.Inventory = append(npc.Inventory, item.ID)
			item.Location = "thief-inventory"
			return "The thief snatches the " + item.Name + " and runs off with a wicked grin!"
		}
	}

	// Give item to NPC (add to their inventory)
	npc.Inventory = append(npc.Inventory, item.ID)
	item.Location = npc.ID

	return "The " + npc.Name + " accepts the " + item.Name + " reluctantly."
}

// handleAttack attacks an NPC or object (V-ATTACK in ZIL)
func (g *GameV2) handleAttack(objName string) string {
	if objName == "" {
		return "Attack what?"
	}

	// Check for NPC
	npc := g.findNPC(objName)
	if npc == nil {
		// Check for item
		item := g.findItem(objName)
		if item != nil {
			return "Violence isn't the answer to this one."
		}
		return "You can't see any " + objName + " here."
	}

	if !npc.Flags.CanFight {
		return "You can't attack the " + npc.Name + "."
	}

	if !npc.Flags.IsAlive {
		return "The " + npc.Name + " is already dead."
	}

	// Check if player has a weapon
	var playerWeapon *Item
	var playerDamage int = 5 // Base damage with bare hands

	for _, itemID := range g.Player.Inventory {
		item := g.Items[itemID]
		if item != nil && item.Flags.IsWeapon {
			playerWeapon = item
			// Different weapons have different damage
			switch itemID {
			case "sword":
				playerDamage = 20
			case "axe":
				playerDamage = 15
			case "knife", "stiletto", "rusty-knife":
				playerDamage = 10
			case "trident":
				playerDamage = 12
			default:
				playerDamage = 8
			}
			break
		}
	}

	if playerWeapon == nil {
		return "Attacking the " + npc.Name + " with your bare hands is suicidal."
	}

	// Combat! Player attacks first
	var result strings.Builder
	result.WriteString("You attack the " + npc.Name + " with your " + playerWeapon.Name + "!\n")

	// Player hits NPC
	npc.Strength -= playerDamage

	if npc.Strength <= 0 {
		// NPC is defeated!
		npc.Flags.IsAlive = false
		npc.Flags.IsAggressive = false

		result.WriteString("The " + npc.Name + " is defeated!\n")

		// Special handling per NPC
		switch npc.ID {
		case "troll":
			// Troll drops treasure and vanishes
			g.Flags["troll-dead"] = true
			// Remove troll from room
			room := g.Rooms[npc.Location]
			if room != nil {
				room.RemoveNPC("troll")
			}
			result.WriteString("Almost as soon as the troll breathes his last breath, a cloud of sinister black fog envelops him, and when the fog lifts, the carcass has disappeared.")

		case "cyclops":
			// Cyclops drops treasure
			g.Flags["cyclops-dead"] = true
			// Add treasure to room
			if treasure := g.Items["cyclops-treasure"]; treasure != nil {
				treasure.Location = npc.Location
				if room := g.Rooms[npc.Location]; room != nil {
					room.AddItem("cyclops-treasure")
				}
			}
			// Replace with corpse
			room := g.Rooms[npc.Location]
			if room != nil {
				room.RemoveNPC("cyclops")
			}
			result.WriteString("The cyclops falls with a thunderous crash. His treasures are now yours!")

		case "thief":
			// Thief drops stolen items
			g.Flags["thief-dead"] = true
			for _, itemID := range npc.Inventory {
				if item := g.Items[itemID]; item != nil {
					item.Location = npc.Location
					if room := g.Rooms[npc.Location]; room != nil {
						room.AddItem(itemID)
					}
				}
			}
			npc.Inventory = []string{}
			// Remove thief from room
			room := g.Rooms[npc.Location]
			if room != nil {
				room.RemoveNPC("thief")
			}
			result.WriteString("The thief falls, and his stolen loot spills across the floor.")
		}

		return strings.TrimSpace(result.String())
	}

	// NPC is still alive and fights back!
	result.WriteString("The " + npc.Name + " is wounded but still fighting!\n")

	// NPC counter-attacks
	npcDamage := npc.Strength / 5 // Simple damage calculation
	if npcDamage < 3 {
		npcDamage = 3
	}

	g.Player.Health -= npcDamage
	result.WriteString(fmt.Sprintf("The %s strikes back, dealing %d damage!\n", npc.Name, npcDamage))

	if g.Player.Health <= 0 {
		g.GameOver = true
		result.WriteString("\n****  You have died  ****\n")
	} else {
		result.WriteString(fmt.Sprintf("Your health: %d\n", g.Player.Health))
	}

	return strings.TrimSpace(result.String())
}

// handleWave waves an item (V-WAVE in ZIL)
func (g *GameV2) handleWave(objName string) string {
	if objName == "" {
		return "Wave what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You don't have that."
	}

	// Special case: waving sceptre
	if item.ID == "sceptre" {
		return "The sceptre glows briefly."
	}

	return "You wave the " + item.Name + " around. Nothing happens."
}

// handleClimb climbs something (V-CLIMB in ZIL)
func (g *GameV2) handleClimb(objName string) string {
	if objName == "" {
		return "Climb what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	// Special cases would go here (ladder, tree, etc.)
	if item.ID == "ladder" {
		return "The ladder is lying on the ground. You can't climb it."
	}

	return "You can't climb that."
}

// handleTie ties something to something else (V-TIE in ZIL)
func (g *GameV2) handleTie(objName string, targetName string) string {
	if objName == "" {
		return "Tie what?"
	}
	if targetName == "" {
		return "Tie it to what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You don't have that."
	}

	// Special case: rope
	if item.ID == "rope" {
		return "You tie the rope, but nothing interesting happens."
	}

	return "You can't tie that."
}

// handleUntie unties something (V-UNTIE in ZIL)
func (g *GameV2) handleUntie(objName string) string {
	if objName == "" {
		return "Untie what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	return "It's not tied."
}

// handleDig digs in the ground (V-DIG in ZIL)
func (g *GameV2) handleDig(objName string) string {
	room := g.Rooms[g.Location]
	if room == nil {
		return "You are nowhere!"
	}

	// Special cases would check for sandy areas, shovel, etc.
	return "The ground is too hard for digging here."
}

// handlePush pushes something (V-PUSH in ZIL)
func (g *GameV2) handlePush(objName string) string {
	if objName == "" {
		return "Push what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	// Special cases would go here (buttons, statues, etc.)
	if strings.Contains(item.ID, "button") {
		// Dam control buttons
		if g.Location == "maintenance-room" {
			switch item.ID {
			case "yellow-button":
				// Yellow button opens the dam gates - drains reservoir
				if g.Flags["dam-open"] {
					return "Click. The gates are already open."
				}
				g.Flags["dam-open"] = true
				g.Flags["low-tide"] = true
				return "Click. You hear a rumbling sound in the distance as the dam gates open and water drains from the reservoir."

			case "blue-button":
				// Blue button closes the dam gates - fills reservoir
				if !g.Flags["dam-open"] {
					return "Click. The gates are already closed."
				}
				g.Flags["dam-open"] = false
				g.Flags["low-tide"] = false
				return "Click. You hear a rushing sound as the dam gates close and water begins to fill the reservoir."

			case "brown-button":
				return "Click. Nothing seems to happen."

			case "red-button":
				return "Click. Nothing seems to happen."
			}
		}

		// Machine control buttons
		if g.Location == "machine-room" {
			switch item.ID {
			case "lower-button":
				// Lower the basket from shaft-room to lower-shaft
				if g.Flags["basket-lowered"] {
					return "Click. The basket is already at the bottom."
				}

				// Transfer basket and its contents
				shaftRoom := g.Rooms["shaft-room"]
				lowerShaft := g.Rooms["lower-shaft"]

				// Remove raised basket from shaft-room
				shaftRoom.RemoveItem("raised-basket")

				// Move all items from raised-basket to lowered-basket
				for _, otherItem := range g.Items {
					if otherItem.Location == "raised-basket" {
						otherItem.Location = "lowered-basket"
					}
				}

				// Add lowered basket to lower-shaft
				loweredBasket := g.Items["lowered-basket"]
				loweredBasket.Location = "lower-shaft"
				lowerShaft.AddItem("lowered-basket")

				g.Flags["basket-lowered"] = true
				return "Click. You hear a whirring sound as the basket descends."

			case "start-button":
				// Raise the basket from lower-shaft to shaft-room
				if !g.Flags["basket-lowered"] {
					return "Click. The basket is already at the top."
				}

				// Transfer basket and its contents
				shaftRoom := g.Rooms["shaft-room"]
				lowerShaft := g.Rooms["lower-shaft"]

				// Remove lowered basket from lower-shaft
				lowerShaft.RemoveItem("lowered-basket")

				// Move all items from lowered-basket to raised-basket
				for _, otherItem := range g.Items {
					if otherItem.Location == "lowered-basket" {
						otherItem.Location = "raised-basket"
					}
				}

				// Add raised basket back to shaft-room
				raisedBasket := g.Items["raised-basket"]
				raisedBasket.Location = "shaft-room"
				shaftRoom.AddItem("raised-basket")

				g.Flags["basket-lowered"] = false
				return "Click. You hear a whirring sound as the basket ascends."

			case "launch-button":
				return "Click. The machine makes a grinding noise but nothing happens."
			}
		}

		return "Click."
	}

	return "Pushing the " + item.Name + " doesn't seem to help."
}

// handlePull pulls something (V-PULL in ZIL)
func (g *GameV2) handlePull(objName string) string {
	if objName == "" {
		return "Pull what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	return "Pulling the " + item.Name + " doesn't seem to help."
}

// handleMoveObject moves an object (V-MOVE in ZIL)
func (g *GameV2) handleMoveObject(objName string) string {
	if objName == "" {
		return "Move what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	// Special case: Moving the rug reveals the trap door
	if item.ID == "rug" && g.Location == "living-room" {
		if g.Flags["trap-door-open"] {
			return "The rug is already moved, and the trap door is visible."
		}
		g.Flags["trap-door-open"] = true
		trapDoor := g.Items["trap-door"]
		if trapDoor != nil {
			trapDoor.Location = "living-room"
			room := g.Rooms["living-room"]
			if room != nil {
				room.AddItem("trap-door")
			}
		}
		return "With the rug moved aside, you can see a closed trap door beneath it!"
	}

	return "Moving the " + item.Name + " doesn't seem to help."
}

// handleRing rings something (V-RING in ZIL)
func (g *GameV2) handleRing(objName string) string {
	if objName == "" {
		return "Ring what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	// Special case: bell
	if item.ID == "bell" {
		return "Ding, dong. The bell echoes throughout the dungeon."
	}

	return "How does one ring a " + item.Name + "?"
}

// handlePray prays (V-PRAY in ZIL)
func (g *GameV2) handlePray() string {
	room := g.Rooms[g.Location]
	if room == nil {
		return "You are nowhere!"
	}

	// Special case: exorcising ghosts at entrance to Hades
	if g.Location == "entrance-to-hades" {
		// Check if ghosts are present
		ghostsPresent := false
		for _, npcID := range room.NPCs {
			if npcID == "ghosts" {
				ghostsPresent = true
				break
			}
		}

		if ghostsPresent {
			// Banish the ghosts
			g.Flags["ghosts-banished"] = true
			room.RemoveNPC("ghosts")
			// Mark ghosts as not alive
			if ghosts := g.NPCs["ghosts"]; ghosts != nil {
				ghosts.Flags.IsAlive = false
			}
			return "The prayer is answered! The evil spirits are dispelled, and a path opens to the south!"
		} else if g.Flags["ghosts-banished"] {
			return "Your prayers have already been answered here."
		}
	}

	// Special case: in temple areas
	if strings.Contains(room.ID, "temple") || strings.Contains(room.ID, "altar") {
		return "Your prayer is heard, but not answered."
	}

	return "If you pray enough, your prayers may be answered."
}

// handleWait waits a turn (V-WAIT in ZIL)
func (g *GameV2) handleWait() string {
	// Time passes... (would trigger turn-based events in full implementation)
	return "Time passes..."
}

// handleEat eats something (V-EAT in ZIL)
func (g *GameV2) handleEat(objName string) string {
	if objName == "" {
		return "Eat what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You don't have that."
	}

	// Special cases for edible items
	if item.ID == "lunch" {
		// Remove from inventory
		for i, id := range g.Player.Inventory {
			if id == item.ID {
				g.Player.Inventory = append(g.Player.Inventory[:i], g.Player.Inventory[i+1:]...)
				break
			}
		}
		delete(g.Items, item.ID)
		return "Thank you very much. It really hit the spot."
	}

	if item.ID == "garlic" {
		// Remove from inventory
		for i, id := range g.Player.Inventory {
			if id == item.ID {
				g.Player.Inventory = append(g.Player.Inventory[:i], g.Player.Inventory[i+1:]...)
				break
			}
		}
		delete(g.Items, item.ID)
		return "What the heck! You won't be bothered by vampires, anyway."
	}

	return "I don't think the " + item.Name + " would agree with you."
}

// handleDrink drinks something (V-DRINK in ZIL)
func (g *GameV2) handleDrink(objName string) string {
	if objName == "" {
		return "Drink what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	// Special case: water
	if item.ID == "water" {
		return "Thank you very much. I was rather thirsty."
	}

	return "I don't think the " + item.Name + " is potable."
}

// handleFill fills a container (V-FILL in ZIL)
func (g *GameV2) handleFill(objName string, sourceName string) string {
	if objName == "" {
		return "Fill what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You don't have that."
	}

	if !item.Flags.IsContainer {
		return "You can't fill that."
	}

	return "There is nothing to fill it with."
}

// handlePour pours from a container (V-POUR in ZIL)
func (g *GameV2) handlePour(objName string, targetName string) string {
	if objName == "" {
		return "Pour what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You don't have that?"
	}

	return "The " + item.Name + " is empty."
}

// handleInflate inflates an object (IBOAT-FUNCTION in ZIL)
func (g *GameV2) handleInflate(objName string, toolName string) string {
	if objName == "" {
		return "Inflate what?"
	}

	// Find the boat
	boat := g.findItem(objName)
	if boat == nil {
		return "You can't see any " + objName + " here."
	}

	// Only works on inflatable boat (deflated)
	if boat.ID != "boat" && boat.ID != "inflatable-boat" {
		return "How can you inflate that?"
	}

	// Check if boat is already inflated
	inflatedBoat := g.Items["inflated-boat"]
	if inflatedBoat != nil && inflatedBoat.Location == g.Location {
		return "Inflating it further would probably burst it."
	}

	// Boat must be on the ground
	if boat.Location != g.Location {
		return "The boat must be on the ground to be inflated."
	}

	// Determine what we're inflating with
	tool := toolName
	if tool == "" {
		tool = "lungs" // Default if no tool specified
	}

	// Check if using pump
	pump := g.findItem(tool)
	if pump != nil && (pump.ID == "pump" || pump.ID == "air-pump") {
		// Success! Inflate the boat
		room := g.Rooms[g.Location]

		// Remove deflated boat
		room.RemoveItem(boat.ID)
		boat.Location = ""

		// Add inflated boat
		inflatedBoat.Location = g.Location
		room.AddItem("inflated-boat")

		// Reset deflate flag (allows passage through narrow areas)
		g.Flags["deflate"] = true

		result := "The boat inflates and appears seaworthy."

		// Check if label hasn't been seen yet
		label := g.Items["boat-label"]
		if label != nil && label.Location == "inflated-boat" {
			result += "\nA tan label is lying inside the boat."
		}

		return result
	} else if tool == "lungs" {
		return "You don't have enough lung power to inflate it."
	} else {
		return "With a " + tool + "? Surely you jest!"
	}
}

// handleDeflate deflates an object (RBOAT-FUNCTION in ZIL)
func (g *GameV2) handleDeflate(objName string) string {
	if objName == "" {
		return "Deflate what?"
	}

	boat := g.findItem(objName)
	if boat == nil {
		return "You can't see any " + objName + " here."
	}

	// Only works on inflated boat
	if boat.ID != "inflated-boat" {
		return "Come on, now!"
	}

	// Can't deflate if player is in the boat
	if g.Location == "inflated-boat" {
		return "You can't deflate the boat while you're in it."
	}

	// Boat must be on ground (not in inventory)
	if boat.Location != g.Location {
		return "The boat must be on the ground to be deflated."
	}

	room := g.Rooms[g.Location]

	// Remove inflated boat
	room.RemoveItem("inflated-boat")
	boat.Location = ""

	// Add deflated boat back
	deflatedBoat := g.Items["boat"]
	if deflatedBoat == nil {
		deflatedBoat = g.Items["inflatable-boat"]
	}
	if deflatedBoat != nil {
		deflatedBoat.Location = g.Location
		room.AddItem(deflatedBoat.ID)
	}

	// Clear deflate flag (blocks passage through narrow areas)
	g.Flags["deflate"] = false

	return "The boat deflates."
}

// handlePlug repairs the punctured boat (DBOAT-FUNCTION in ZIL)
func (g *GameV2) handlePlug(objName string, materialName string) string {
	if objName == "" {
		return "Plug what?"
	}

	boat := g.findItem(objName)
	if boat == nil {
		return "You can't see any " + objName + " here."
	}

	// Only works on punctured boat
	if boat.ID != "punctured-boat" {
		return "That doesn't need plugging."
	}

	// Check for putty
	if materialName == "" {
		return "Plug it with what?"
	}

	material := g.findItem(materialName)
	if material == nil || material.ID != "putty" {
		return "That won't work."
	}

	// Success! Repair the boat
	room := g.Rooms[g.Location]

	// Remove punctured boat
	room.RemoveItem("punctured-boat")
	boat.Location = ""

	// Add deflated boat
	deflatedBoat := g.Items["boat"]
	if deflatedBoat == nil {
		deflatedBoat = g.Items["inflatable-boat"]
	}
	if deflatedBoat != nil {
		deflatedBoat.Location = g.Location
		room.AddItem(deflatedBoat.ID)
	}

	return "Well done. The boat is repaired."
}

// handleListen listens (V-LISTEN in ZIL)
func (g *GameV2) handleListen() string {
	room := g.Rooms[g.Location]
	if room == nil {
		return "You are nowhere!"
	}

	// Special cases for specific rooms
	return "You hear nothing unusual."
}

// handleSmell smells something (V-SMELL in ZIL)
func (g *GameV2) handleSmell(objName string) string {
	if objName == "" {
		return "Smell what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	return "It doesn't smell unusual."
}

// handleTouch touches something (V-TOUCH in ZIL)
func (g *GameV2) handleTouch(objName string) string {
	if objName == "" {
		return "Touch what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	return "You feel nothing unexpected."
}

// handleSearch searches something (V-SEARCH in ZIL)
func (g *GameV2) handleSearch(objName string) string {
	if objName == "" {
		// Search the room
		return g.handleLook()
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	if item.Flags.IsContainer {
		return g.handleLookIn(item.ID)
	}

	return "You find nothing special."
}

// handleJump jumps (V-JUMP in ZIL)
func (g *GameV2) handleJump() string {
	return "You jump on the spot, fruitlessly."
}

// handleSwim swims (V-SWIM in ZIL)
func (g *GameV2) handleSwim() string {
	room := g.Rooms[g.Location]
	if room == nil {
		return "You are nowhere!"
	}

	// Check for water
	if strings.Contains(room.ID, "river") || strings.Contains(room.ID, "reservoir") {
		return "You would drown."
	}

	return "There is no water here."
}

// handleBlow blows something (V-BLOW in ZIL)
func (g *GameV2) handleBlow(objName string) string {
	if objName == "" {
		return "Blow what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You don't have that."
	}

	return "You can't blow that."
}

// handleKnock knocks on something (V-KNOCK in ZIL)
func (g *GameV2) handleKnock(objName string) string {
	if objName == "" {
		return "Knock on what?"
	}

	// For "knock on X", objName might be empty and we need to check for direct object
	// The parser puts "door" in IndirectObject for "knock on door"
	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	return "No one answers."
}

// handleScore shows the current score (V-SCORE in ZIL)
func (g *GameV2) handleScore() string {
	// Determine rank based on score (from original Zork)
	rank := ""
	switch {
	case g.Score >= 350:
		rank = "Master Adventurer"
	case g.Score >= 330:
		rank = "Wizard"
	case g.Score >= 300:
		rank = "Master"
	case g.Score >= 200:
		rank = "Adventurer"
	case g.Score >= 100:
		rank = "Junior Adventurer"
	case g.Score >= 50:
		rank = "Novice Adventurer"
	case g.Score >= 25:
		rank = "Amateur Adventurer"
	default:
		rank = "Beginner"
	}

	return fmt.Sprintf("Your score is %d (out of 350), in %d move(s).\nThis gives you the rank of %s.",
		g.Score, g.Moves, rank)
}

// GetInitialMessage returns the opening text
func (g *GameV2) GetInitialMessage() string {
	return `ZORK I: The Great Underground Empire
Copyright (c) 1981, 1982, 1983 Infocom, Inc. All rights reserved.
ZORK is a registered trademark of Infocom, Inc.
Revision 88 / Serial number 840726

` + g.handleLook()
}
