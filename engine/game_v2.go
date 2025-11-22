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
	// The Troll - First major enemy
	troll := NewNPC(
		"troll",
		"nasty troll",
		"A nasty-looking troll, brandishing a bloody axe, blocks all passages out of the room.",
	)
	troll.Location = "troll-room"
	troll.Strength = 20
	troll.Weapon = "troll-axe" // We'd create this item
	troll.Hostile = true
	troll.Flags.IsAggressive = true
	troll.Flags.CanFight = true
	g.NPCs["troll"] = troll
	g.Rooms["troll-room"].AddNPC("troll")

	// Add more NPCs (thief, cyclops, etc.)...
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

	// Handle movement
	if cmd.Verb == "walk" && cmd.Direction != "" {
		return g.handleMove(cmd.Direction)
	}

	// Handle other verbs
	switch cmd.Verb {
	case "look":
		if cmd.Preposition == "in" || cmd.Preposition == "into" {
			return g.handleLookIn(cmd.IndirectObject)
		}
		return g.handleLook()
	case "examine":
		return g.handleExamine(cmd.DirectObject)
	case "take":
		return g.handleTake(cmd.DirectObject)
	case "drop":
		return g.handleDrop(cmd.DirectObject)
	case "open":
		return g.handleOpen(cmd.DirectObject)
	case "close":
		return g.handleClose(cmd.DirectObject)
	case "read":
		return g.handleRead(cmd.DirectObject)
	case "turn":
		// Handle both "turn on lamp" and "turn lamp on"
		objName := cmd.DirectObject
		if objName == "" {
			objName = cmd.IndirectObject
		}

		if cmd.Preposition == "on" {
			return g.handleTurnOn(objName)
		} else if cmd.Preposition == "off" {
			return g.handleTurnOff(objName)
		}
		return "Turn it on or off?"
	case "inventory":
		return g.handleInventory()
	case "help":
		return g.handleHelp()
	case "put":
		return g.handlePut(cmd.DirectObject, cmd.Preposition, cmd.IndirectObject)
	case "give":
		return g.handleGive(cmd.DirectObject, cmd.IndirectObject)
	case "attack":
		return g.handleAttack(cmd.DirectObject)
	case "wave":
		return g.handleWave(cmd.DirectObject)
	case "climb":
		return g.handleClimb(cmd.DirectObject)
	case "tie":
		return g.handleTie(cmd.DirectObject, cmd.IndirectObject)
	case "untie":
		return g.handleUntie(cmd.DirectObject)
	case "dig":
		return g.handleDig(cmd.DirectObject)
	case "push":
		return g.handlePush(cmd.DirectObject)
	case "pull":
		return g.handlePull(cmd.DirectObject)
	case "ring":
		return g.handleRing(cmd.DirectObject)
	case "pray":
		return g.handlePray()
	case "wait":
		return g.handleWait()
	case "eat":
		return g.handleEat(cmd.DirectObject)
	case "drink":
		return g.handleDrink(cmd.DirectObject)
	case "fill":
		return g.handleFill(cmd.DirectObject, cmd.IndirectObject)
	case "pour":
		return g.handlePour(cmd.DirectObject, cmd.IndirectObject)
	case "listen":
		return g.handleListen()
	case "smell":
		return g.handleSmell(cmd.DirectObject)
	case "touch":
		return g.handleTouch(cmd.DirectObject)
	case "search":
		return g.handleSearch(cmd.DirectObject)
	case "jump":
		return g.handleJump()
	case "swim":
		return g.handleSwim()
	case "blow":
		return g.handleBlow(cmd.DirectObject)
	case "knock":
		return g.handleKnock(cmd.DirectObject)
	case "quit":
		g.GameOver = true
		return "Thanks for playing!"
	case "save":
		return "Save is not yet implemented."
	case "restore":
		return "Restore is not yet implemented."
	case "score":
		return g.handleScore()
	case "restart":
		return "Restart is not yet implemented."
	default:
		return "I don't understand how to \"" + cmd.Verb + "\" something."
	}
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

	// Remove from room
	room := g.Rooms[g.Location]
	if room != nil {
		room.RemoveItem(item.ID)
	}

	// Add to inventory
	item.Location = "inventory"
	g.Player.Inventory = append(g.Player.Inventory, item.ID)

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

	// Special handling for kitchen window
	if item.ID == "kitchen-window" {
		if g.Flags["window-open"] {
			return "It's already open."
		}
		g.Flags["window-open"] = true
		item.Flags.IsOpen = true
		return "With great effort, you open the window far enough to allow entry."
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
		if npc != nil && npc.Name == name {
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

	// Special case: troll likes food
	if npc.ID == "troll" && (item.ID == "lunch" || item.ID == "garlic") {
		// Remove from inventory
		for i, id := range g.Player.Inventory {
			if id == item.ID {
				g.Player.Inventory = append(g.Player.Inventory[:i], g.Player.Inventory[i+1:]...)
				break
			}
		}
		return "The troll grabs the " + item.Name + " and devours it. He looks satisfied and wanders off."
	}

	return "The " + npc.Name + " doesn't want that."
}

// handleAttack attacks an NPC or object (V-ATTACK in ZIL)
func (g *GameV2) handleAttack(objName string) string {
	if objName == "" {
		return "Attack what?"
	}

	// Check for NPC
	npc := g.findNPC(objName)
	if npc != nil {
		if !npc.Flags.CanFight {
			return "You can't attack the " + npc.Name + "."
		}
		// Basic combat - needs full combat system
		return "Attacking the " + npc.Name + " with your bare hands is suicidal."
	}

	// Check for item
	item := g.findItem(objName)
	if item != nil {
		return "Violence isn't the answer to this one."
	}

	return "You can't see any " + objName + " here."
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
		return "You don't have that."
	}

	return "The " + item.Name + " is empty."
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

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	return "No one answers."
}

// handleScore shows the current score (V-SCORE in ZIL)
func (g *GameV2) handleScore() string {
	return fmt.Sprintf("Your score is %d (out of 350), in %d move(s).", g.Score, g.Moves)
}

// GetInitialMessage returns the opening text
func (g *GameV2) GetInitialMessage() string {
	return `ZORK I: The Great Underground Empire
Copyright (c) 1981, 1982, 1983 Infocom, Inc. All rights reserved.
ZORK is a registered trademark of Infocom, Inc.
Revision 88 / Serial number 840726

` + g.handleLook()
}
