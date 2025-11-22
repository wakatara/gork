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
