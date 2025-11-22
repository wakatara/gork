package engine

import (
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
	case "quit":
		g.GameOver = true
		return "Thanks for playing!"
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

// GetInitialMessage returns the opening text
func (g *GameV2) GetInitialMessage() string {
	return `ZORK I: The Great Underground Empire
Copyright (c) 1981, 1982, 1983 Infocom, Inc. All rights reserved.
ZORK is a registered trademark of Infocom, Inc.
Revision 88 / Serial number 840726

` + g.handleLook()
}
