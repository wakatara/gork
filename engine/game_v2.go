package engine

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
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
	rand      *rand.Rand // Random number generator for thief AI
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
		rand:  rand.New(rand.NewSource(time.Now().UnixNano())),
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

	// Initialize game flags
	g.Flags["GRUNLOCK"] = true // Grating starts unlocked (can be opened from either side)
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

	// The Ghosts/Spirits - Block entrance to Land of the Dead (GHOSTS in ZIL 1dungeon.zil lines 109-115)
	ghosts := NewNPC(
		"ghosts",
		"evil spirits",
		"The way through the gate is barred by evil spirits, who jeer at your attempts to pass.",
	)
	ghosts.Location = "entrance-to-hades"
	ghosts.Strength = 0 // Can't be fought
	ghosts.Weapon = ""
	ghosts.Hostile = false // Don't attack, just block passage
	ghosts.Flags.IsAggressive = false
	ghosts.Flags.CanFight = false
	ghosts.Flags.IsAlive = true
	g.NPCs["ghosts"] = ghosts
	g.Rooms["entrance-to-hades"].AddNPC("ghosts")

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
		case "look-in", "look-on":
			// Multi-word verb "look in" or "look on"
			result = g.handleLookIn(cmd.DirectObject)
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
		case "unlock":
			result = g.handleUnlock(cmd.DirectObject, cmd.IndirectObject)
		case "lock":
			result = g.handleLock(cmd.DirectObject, cmd.IndirectObject)
		case "read":
			result = g.handleRead(cmd.DirectObject)
		case "turn":
			// Handle both "turn on lamp" and "turn lamp on"
			objName := cmd.DirectObject
			if objName == "" {
				objName = cmd.IndirectObject
			}

			// Special handling for prayer book (BLACK-BOOK in ZIL lines 2192-2200)
			if objName != "" {
				item := g.findItem(objName)
				if item != nil && item.ID == "book" && cmd.Preposition == "" {
					result = `Beside page 569, there is only one other page with any legible printing on it. Most of it is unreadable, but the subject seems to be the banishment of evil. Apparently, certain noises, lights, and prayers are efficacious in this regard.`
					break
				}
			}

			if cmd.Preposition == "on" {
				result = g.handleTurnOn(objName)
			} else if cmd.Preposition == "off" {
				result = g.handleTurnOff(objName)
			} else {
				result = "Turn it on or off?"
			}
		case "turn-on":
			// Multi-word verb "turn on"
			result = g.handleTurnOn(cmd.DirectObject)
		case "turn-off":
			// Multi-word verb "turn off"
			result = g.handleTurnOff(cmd.DirectObject)
		case "light":
			result = g.handleTurnOn(cmd.DirectObject)
		case "extinguish":
			result = g.handleTurnOff(cmd.DirectObject)
		case "inventory":
			result = g.handleInventory()
		case "help":
			result = g.handleHelp()
		case "put":
			result = g.handlePut(cmd.DirectObject, cmd.Preposition, cmd.IndirectObject)
		case "put-on":
			// Multi-word "put on"
			result = g.handlePut(cmd.DirectObject, "on", cmd.IndirectObject)
		case "give":
			result = g.handleGive(cmd.DirectObject, cmd.IndirectObject)
		case "attack":
			result = g.handleAttack(cmd.DirectObject)
		case "wave":
			result = g.handleWave(cmd.DirectObject)
		case "climb":
			result = g.handleClimb(cmd.DirectObject)
		case "climb-up":
			// Multi-word "climb up"
			result = "You can't climb that."
		case "climb-down":
			// Multi-word "climb down"
			result = "You can't climb down that."
		case "climb-on":
			// Multi-word "climb on"
			result = "You can't climb on that."
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
		case "exorcise":
			result = g.handleExorcise(cmd.DirectObject)
		case "pray":
			result = g.handlePray()
		case "ulysses", "odysseus":
			result = g.handleOdysseus()
		case "wait":
			result = g.handleWait()
		case "eat":
			result = g.handleEat(cmd.DirectObject)
		case "drink":
			result = g.handleDrink(cmd.DirectObject)
		case "drink-from":
			// Multi-word "drink from"
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
		case "break":
			result = g.handleBreak(cmd.DirectObject)
		case "burn":
			result = g.handleBurn(cmd.DirectObject)
		case "search":
			result = g.handleSearch(cmd.DirectObject)
		case "jump":
			result = g.handleJump()
		case "swim":
			result = g.handleSwim()
		case "blow":
			result = g.handleBlow(cmd.DirectObject)
		case "blow-out":
			// Multi-word "blow out"
			result = g.handleTurnOff(cmd.DirectObject)
		case "blow-up":
			// Multi-word "blow up"
			result = "That would be dangerous."
		case "blow-in":
			// Multi-word "blow in"
			result = "That doesn't help."
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
			result = g.handleSave(cmd)
		case "restore":
			result = g.handleRestore(cmd)
		case "score":
			result = g.handleScore()
		case "restart":
			result = "Restart is not yet implemented."
		case "enter":
			result = g.handleEnter(cmd)
		case "exit", "leave":
			result = g.handleMove("out")
		case "throw":
			result = g.handleThrow(cmd)
		case "kill":
			result = g.handleAttack(cmd.DirectObject)
		case "yell", "scream", "shout":
			result = g.handleYell()
		case "board":
			result = g.handleBoard(cmd.DirectObject)
		case "disembark":
			result = g.handleMove("out")
		case "brief":
			result = "Brief mode is now on."
		case "verbose":
			result = "Verbose mode is now on."
		case "superbrief":
			result = "Superbrief mode is now on."
		case "diagnose":
			result = g.handleDiagnose()
		case "version":
			result = "ZORK I: The Great Underground Empire\nGo Edition Version 1.0\nOriginal game Copyright (c) 1981, 1982, 1983 Infocom, Inc."
		case "say", "speak":
			result = g.handleSay(cmd)
		case "find", "where":
			result = "You'll have to find it yourself."
		case "curse", "damn", "shit", "fuck":
			result = "Such language in a high-class establishment like this!"
		case "xyzzy", "plugh":
			result = "A hollow voice says \"Fool.\""
		case "win":
			result = "Preposterous!"
		case "bug":
			result = "No bugs here. This is a feature-complete implementation."
		case "chomp":
			result = "I don't think the dungeon master would approve."
		case "zork":
			result = "At your service!"
		case "echo":
			result = g.handleEcho()
		case "script":
			result = "Scripting is not implemented in this version."
		case "unscript":
			result = "Scripting is not implemented in this version."
		case "cross", "ford":
			result = "You can't cross that."
		case "kick", "taunt":
			result = "Kicking things won't help."
		case "melt", "liquify":
			result = "You have nothing to melt it with."
		case "repent", "sigh":
			result = "It's a bit late for that."
		case "sleep":
			result = "This is no time for sleeping!"
		case "wake":
			result = "The dungeon master does not allow sleeping in the dungeon."
		case "wish":
			result = "You have been granted 1 wish. Too bad you just used it up."
		case "mumble":
			result = "You mumble to yourself. Nothing happens."
		default:
			result = "I don't understand how to \"" + cmd.Verb + "\" something."
		}
	}

	// Process NPC turns after every command (including grues!)
	npcResult := g.processNPCTurns()
	if npcResult != "" {
		result += "\n\n" + npcResult
	}

	// Process lamp fuel depletion
	lampResult := g.processLampFuel()
	if lampResult != "" {
		result += "\n\n" + lampResult
	}

	// Process candles fuel depletion
	candlesResult := g.processCandlesFuel()
	if candlesResult != "" {
		result += "\n\n" + candlesResult
	}

	// Process thief behavior
	thiefResult := g.processThiefTurn()
	if thiefResult != "" {
		result += "\n\n" + thiefResult
	}

	// Process sword glowing
	swordResult := g.processSwordGlow()
	if swordResult != "" {
		result += "\n\n" + swordResult
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

// processLampFuel handles lamp fuel depletion each turn (I-LANTERN in ZIL)
func (g *GameV2) processLampFuel() string {
	lamp := g.Items["lamp"]
	if lamp == nil || !lamp.Flags.IsLit || lamp.Fuel <= 0 {
		return ""
	}

	// Decrement fuel
	lamp.Fuel--

	// Check for warning messages at specific fuel levels
	// ZIL LAMP-TABLE: 100, 70, 15, 0
	switch lamp.Fuel {
	case 230: // After 100 turns (330 - 100)
		return "The lamp appears a bit dimmer."
	case 160: // After 170 turns (330 - 170 = 160)
		return "The lamp is definitely dimmer now."
	case 145: // After 185 turns (330 - 185 = 145)
		return "The lamp is nearly out."
	case 0:
		// Lamp has died
		lamp.Flags.IsLit = false
		result := "The lamp has gone out."

		// Check if player is now in darkness
		room := g.Rooms[g.Location]
		if room != nil && room.Flags.IsDark && !g.hasLight() {
			result += "\n\nOh, no! You have walked into the slavering fangs of a lurking grue!\n\n****  You have died  ****"
			g.GameOver = true
		}

		return result
	}

	return ""
}

// processCandlesFuel depletes candle fuel each turn (I-CANDLES in ZIL lines 2321-2326)
func (g *GameV2) processCandlesFuel() string {
	candles := g.Items["candles"]
	if candles == nil || !candles.Flags.IsLit || candles.Fuel <= 0 {
		return ""
	}

	// Decrement fuel
	candles.Fuel--

	// Check for warning messages at specific fuel levels
	// ZIL CANDLE-TABLE (lines 2406-2414): 20, 10, 5, 0
	switch candles.Fuel {
	case 20:
		return "The candles grow shorter."
	case 10:
		return "The candles are becoming quite short."
	case 5:
		return "The candles won't last long now."
	case 0:
		// Candles have burned out
		candles.Flags.IsLit = false
		result := "You'd better have more light than from the candles."

		// Check if player is now in darkness
		room := g.Rooms[g.Location]
		if room != nil && room.Flags.IsDark && !g.hasLight() {
			result += "\n\nOh, no! You have walked into the slavering fangs of a lurking grue!\n\n****  You have died  ****"
			g.GameOver = true
		}

		return result
	}

	return ""
}

// processThiefTurn handles thief AI: movement, stealing, depositing treasures (I-THIEF in ZIL lines 3890-3931)
func (g *GameV2) processThiefTurn() string {
	thief := g.NPCs["thief"]
	if thief == nil || !thief.Flags.IsAlive {
		return ""
	}

	var result string

	// Check if thief is in player's room
	thiefInRoom := (thief.Location == g.Location)

	// 1. If thief is at treasure-room and player is not there, deposit treasures
	if thief.Location == "treasure-room" && !thiefInRoom {
		// Deposit treasures silently
		g.depositThiefTreasures()
	}

	// 2. If thief is in same room as player, handle interactions
	if thiefInRoom && !g.GameOver {
		// Random chance to steal from player or room
		if g.randomChance(40) { // 40% chance to steal
			stolen := g.thiefStealTreasures()
			if stolen != "" {
				result = stolen
			}
		}
	}

	// 3. Move thief to next room
	g.moveThiefToNextRoom()

	// 4. If thief ended up in player's room after moving, maybe reveal presence
	if thief.Location == g.Location && g.randomChance(30) {
		if len(thief.Inventory) > 0 {
			result += "\nA seedy-looking individual with a large bag just wandered through the room."
		}
	}

	return result
}

// moveThiefToNextRoom moves thief to a random adjacent room
func (g *GameV2) moveThiefToNextRoom() {
	thief := g.NPCs["thief"]
	if thief == nil {
		return
	}

	currentRoom := g.Rooms[thief.Location]
	if currentRoom == nil {
		return
	}

	// Get all possible exits
	var possibleRooms []string
	for _, exit := range currentRoom.Exits {
		// Don't go to sacred rooms
		targetRoom := g.Rooms[exit.To]
		if targetRoom != nil {
			possibleRooms = append(possibleRooms, exit.To)
		}
	}

	// Pick a random room
	if len(possibleRooms) > 0 {
		newRoom := possibleRooms[g.randomInt(len(possibleRooms))]

		// Remove from old room
		oldRoom := g.Rooms[thief.Location]
		if oldRoom != nil {
			oldRoom.RemoveNPC("thief")
		}

		// Move to new room
		thief.Location = newRoom
		newRoomObj := g.Rooms[newRoom]
		if newRoomObj != nil {
			newRoomObj.AddNPC("thief")
		}
	}
}

// thiefStealTreasures attempts to steal treasures from room or player
func (g *GameV2) thiefStealTreasures() string {
	thief := g.NPCs["thief"]
	if thief == nil {
		return ""
	}

	var stolen []string

	// Try to steal from room first
	room := g.Rooms[g.Location]
	if room != nil {
		for _, itemID := range room.Contents {
			item := g.Items[itemID]
			if item != nil && item.Flags.IsTreasure && g.randomChance(75) {
				// Steal it!
				item.Location = "thief-inventory"
				thief.Inventory = append(thief.Inventory, itemID)
				room.RemoveItem(itemID)
				stolen = append(stolen, item.Name)
			}
		}
	}

	// Try to steal from player
	for i := len(g.Player.Inventory) - 1; i >= 0; i-- {
		itemID := g.Player.Inventory[i]
		item := g.Items[itemID]
		if item != nil && item.Flags.IsTreasure && g.randomChance(50) {
			// Steal it!
			item.Location = "thief-inventory"
			thief.Inventory = append(thief.Inventory, itemID)
			g.Player.Inventory = append(g.Player.Inventory[:i], g.Player.Inventory[i+1:]...)
			stolen = append(stolen, item.Name)
		}
	}

	if len(stolen) > 0 {
		return `A seedy-looking individual with a large bag just wandered through the room. On the way through, he quietly abstracted some valuables from your possession, mumbling something about "Doing unto others before..."`
	}

	return ""
}

// depositThiefTreasures moves treasures from thief to treasure-room
func (g *GameV2) depositThiefTreasures() {
	thief := g.NPCs["thief"]
	if thief == nil {
		return
	}

	treasureRoom := g.Rooms["treasure-room"]
	if treasureRoom == nil {
		return
	}

	// Move all treasures from thief to treasure-room
	for _, itemID := range thief.Inventory {
		item := g.Items[itemID]
		if item != nil && item.Flags.IsTreasure {
			item.Location = "treasure-room"
			treasureRoom.AddItem(itemID)
		}
	}

	// Clear thief inventory
	thief.Inventory = []string{}
}

// randomChance returns true with given percentage probability
func (g *GameV2) randomChance(percent int) bool {
	return g.randomInt(100) < percent
}

// randomInt returns a random integer from 0 to n-1
func (g *GameV2) randomInt(n int) int {
	if n <= 0 {
		return 0
	}
	return g.rand.Intn(n)
}

// processSwordGlow updates sword glow level based on nearby enemies (I-SWORD in ZIL lines 3851-3879)
func (g *GameV2) processSwordGlow() string {
	sword := g.Items["sword"]
	if sword == nil {
		return ""
	}

	// Only glow when player is holding it
	isHolding := false
	for _, itemID := range g.Player.Inventory {
		if itemID == "sword" {
			isHolding = true
			break
		}
	}

	if !isHolding {
		return ""
	}

	oldGlow := sword.GlowLevel
	newGlow := 0

	// Check current room for hostile NPCs
	if g.isRoomInfested(g.Location) {
		newGlow = 2 // Very bright
	} else {
		// Check adjacent rooms for hostile NPCs
		currentRoom := g.Rooms[g.Location]
		if currentRoom != nil {
			for _, exit := range currentRoom.Exits {
				if g.isRoomInfested(exit.To) {
					newGlow = 1 // Faint glow
					break
				}
			}
		}
	}

	// Only report changes in glow level
	if newGlow == oldGlow {
		return ""
	}

	sword.GlowLevel = newGlow

	switch newGlow {
	case 2:
		return "Your sword has begun to glow very brightly."
	case 1:
		return "Your sword is glowing with a faint blue glow."
	case 0:
		return "Your sword is no longer glowing."
	}

	return ""
}

// isRoomInfested checks if a room has hostile NPCs (INFESTED? in ZIL lines 3881-3886)
func (g *GameV2) isRoomInfested(roomID string) bool {
	room := g.Rooms[roomID]
	if room == nil {
		return false
	}

	// Check for hostile NPCs in the room
	for _, npcID := range room.NPCs {
		npc := g.NPCs[npcID]
		if npc != nil && npc.Flags.IsAlive && npc.Hostile {
			return true
		}
	}

	return false
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

		// Special case: trap door is hidden until rug is moved (in living-room only, always visible in cellar)
		if item != nil && item.ID == "trap-door" && g.Location == "living-room" && !g.Flags["trap-door-open"] {
			continue // Skip trap door in living-room if rug hasn't been moved
		}

		if item != nil && !item.Flags.IsInvisible {
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
		// Special case: examining mirrors (MIRROR-MIRROR in ZIL lines 994-999)
		if item.ID == "mirror-1" || item.ID == "mirror-2" {
			if g.Flags["mirror-mung"] {
				return "The mirror is broken into many pieces."
			}
			return "There is an ugly person staring back at you."
		}

		// Special case: examining sword (SWORD-FCN in ZIL lines 2432-2442)
		if item.ID == "sword" {
			result := item.Description
			switch item.GlowLevel {
			case 1:
				result += " It is glowing with a faint blue glow."
			case 2:
				result += " It is glowing very brightly."
			}
			return result
		}

		// Special case: examining kitchen window (KITCHEN-WINDOW-F in ZIL lines 247-250)
		if item.ID == "kitchen-window" && !g.Flags["window-opened-once"] {
			return "The window is slightly ajar, but not enough to allow entry."
		}

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

	// Special case: taking mirrors (MIRROR-MIRROR in ZIL lines 1000-1002)
	if item.ID == "mirror-1" || item.ID == "mirror-2" {
		return "The mirror is many times your size. Give up."
	}

	if !item.Flags.IsTakeable {
		return "You can't take the " + item.Name + "."
	}

	// Special case: Taking the rug reveals the trap door
	if item.ID == "rug" && g.Location == "living-room" {
		g.Flags["trap-door-open"] = true
		// Trap door is already in the room (global object), just needs to be revealed
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

	// Special handling for prayer book (BLACK-BOOK in ZIL lines 2187-2189)
	if item.ID == "book" {
		return "The book is already open to page 569."
	}

	// Special handling for kitchen window (KITCHEN-WINDOW-F in ZIL lines 242-246)
	if item.ID == "kitchen-window" {
		if g.Flags["window-open"] {
			return "It's already open."
		}
		g.Flags["window-open"] = true
		g.Flags["window-opened-once"] = true // Track that window has been opened at least once
		item.Flags.IsOpen = true
		return "With great effort, you open the window far enough to allow entry."
	}

	// Special handling for trap door (TRAP-DOOR-FCN in ZIL lines 504-529)
	if item.ID == "trap-door" {
		// Can only open from living-room
		if g.Location == "living-room" {
			if !g.Flags["trap-door-open"] {
				return "The rug must be moved first."
			}
			if item.Flags.IsOpen {
				return "It's already open."
			}
			item.Flags.IsOpen = true
			return "The door reluctantly opens to reveal a rickety staircase descending into darkness."
		} else if g.Location == "cellar" {
			// From cellar, door is locked from above
			if !item.Flags.IsOpen {
				return "The door is locked from above."
			}
			return "It's already open."
		}
		return "You can't see any trap-door here."
	}

	// Special handling for grating (GRATE-FUNCTION in ZIL)
	if item.ID == "grating" || item.ID == "grate" {
		// Check if grate is unlocked
		if !g.Flags["GRUNLOCK"] {
			return "The grating is locked."
		}

		// Check if already open
		if item.Flags.IsOpen {
			return "It's already open."
		}

		// Open the grate
		item.Flags.IsOpen = true
		g.Flags["grate-open"] = true

		// Different messages depending on location
		var result string
		if g.Location == "grating-clearing" {
			result = "The grating opens."
		} else {
			result = "The grating opens to reveal trees above you."
		}

		// If opening from below (grating-room) and leaves haven't been revealed yet
		if g.Location == "grating-room" && !g.Flags["grate-revealed"] {
			result += "\nA pile of leaves falls onto your head and to the ground."
			g.Flags["grate-revealed"] = true

			// Add leaves to the room if they exist
			leaves := g.Items["pile-of-leaves"]
			if leaves != nil {
				leaves.Location = g.Location
				room := g.Rooms[g.Location]
				if room != nil {
					room.AddItem("pile-of-leaves")
				}
			}
		}

		return result
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

	// Special handling for prayer book (BLACK-BOOK in ZIL lines 2190-2191)
	if item.ID == "book" {
		return "As hard as you try, the book cannot be closed."
	}

	// Special handling for kitchen window
	if item.ID == "kitchen-window" {
		if !g.Flags["window-open"] {
			return "It's already closed."
		}
		g.Flags["window-open"] = false
		item.Flags.IsOpen = false
		return "The window closes (more easily than it opened)."
	}

	// Special handling for trap door (TRAP-DOOR-FCN in ZIL lines 504-529)
	if item.ID == "trap-door" {
		if g.Location == "living-room" {
			if !item.Flags.IsOpen {
				return "It's already closed."
			}
			item.Flags.IsOpen = false
			return "The door swings shut and closes."
		} else if g.Location == "cellar" {
			// From cellar, closing also locks it
			if !item.Flags.IsOpen {
				return "It's already closed."
			}
			item.Flags.IsOpen = false
			return "The door closes and locks."
		}
		return "You can't see any trap-door here."
	}

	// Special handling for grating (GRATE-FUNCTION in ZIL)
	if item.ID == "grating" || item.ID == "grate" {
		if !item.Flags.IsOpen {
			return "It's already closed."
		}
		item.Flags.IsOpen = false
		g.Flags["grate-open"] = false
		return "The grating is closed."
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

// handleUnlock unlocks an object with a tool (GRATE-FUNCTION in ZIL)
func (g *GameV2) handleUnlock(objName string, toolName string) string {
	if objName == "" {
		return "What do you want to unlock?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	// Special handling for grate/grating
	if item.ID == "grating" || item.ID == "grate" {
		// Check location - can only unlock from inside (grating-room)
		if g.Location != "grating-room" {
			if g.Location == "grating-clearing" {
				return "You can't reach the lock from here."
			}
			return "You can't unlock that from here."
		}

		// Check for keys
		if toolName == "" {
			return "Unlock it with what?"
		}

		tool := g.findItem(toolName)
		if tool == nil || tool.ID != "keys" {
			return "Can you unlock a grating with a " + toolName + "?"
		}

		// Check if already unlocked
		if g.Flags["GRUNLOCK"] {
			return "It's already unlocked."
		}

		// Success! Unlock the grate
		g.Flags["GRUNLOCK"] = true
		return "The grate is unlocked."
	}

	return "You can't unlock that."
}

// handleLock locks an object (GRATE-FUNCTION in ZIL)
func (g *GameV2) handleLock(objName string, toolName string) string {
	if objName == "" {
		return "What do you want to lock?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	// Special handling for grate/grating
	if item.ID == "grating" || item.ID == "grate" {
		// Check location - can only lock from inside (grating-room)
		if g.Location != "grating-room" {
			if g.Location == "grating-clearing" {
				return "You can't lock it from this side."
			}
			return "You can't lock that from here."
		}

		// Check if already locked
		if !g.Flags["GRUNLOCK"] {
			return "It's already locked."
		}

		// Success! Lock the grate
		g.Flags["GRUNLOCK"] = false
		return "The grate is locked."
	}

	return "You can't lock that."
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

	// Special case: reading book during ceremony (LLD-ROOM M-BEG in ZIL lines 1102-1113)
	if item.ID == "book" && g.Location == "entrance-to-hades" && g.Flags["XC"] && !g.Flags["LLD-FLAG"] {
		// Complete the ceremony!
		g.Flags["LLD-FLAG"] = true

		// Remove ghosts
		ghosts := g.NPCs["ghosts"]
		if ghosts != nil {
			ghosts.Flags.IsAlive = false
			room := g.Rooms["entrance-to-hades"]
			if room != nil {
				room.RemoveNPC("ghosts")
			}
		}

		return `Each word of the prayer reverberates through the hall in a deafening confusion. As the last word fades, a voice, loud and commanding, speaks: "Begone, fiends!" A heart-stopping scream fills the cavern, and the spirits, sensing a greater power, flee through the walls.`
	}

	// If item has Text content, return that; otherwise return Description
	// (V-READ in ZIL gverbs.zil lines 1143-1147)
	if item.Text != "" {
		return item.Text
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

	// Special case: lighting candles during bell ceremony (LLD-ROOM M-END in ZIL lines 1115-1125)
	if item.ID == "candles" && g.Location == "entrance-to-hades" && g.Flags["XB"] && !g.Flags["XC"] && !g.Flags["LLD-FLAG"] {
		g.Flags["XC"] = true
		g.Flags["candles-ceremony-active"] = true
		return `The flames flicker wildly and appear to dance. The earth beneath your feet trembles, and your legs nearly buckle beneath you. The spirits cower at your unearthly power.`
	}

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

			// Special case: trap door is hidden until rug is moved (in living-room only)
			if item != nil && item.ID == "trap-door" && g.Location == "living-room" && !g.Flags["trap-door-open"] {
				continue // Skip trap door if rug hasn't been moved
			}

			if item != nil && !item.Flags.IsInvisible && item.HasAlias(name) {
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

func (g *GameV2) hasItemInInventory(itemID string) bool {
	for _, id := range g.Player.Inventory {
		if id == itemID {
			return true
		}
	}
	return false
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
		// Troll accepts treasures as bribes (TROLL-FCN in ZIL lines 709-751)
		if item.Flags.IsTreasure {
			// Troll is bribed and leaves
			g.Flags["troll-dead"] = true // Opens passages (same flag as when killed)
			room := g.Rooms[npc.Location]

			// Troll drops his axe
			axe := g.Items["axe"]
			if axe != nil {
				axe.Location = npc.Location
				if room != nil {
					room.AddItem("axe")
				}
			}

			// Remove troll from room
			if room != nil {
				room.RemoveNPC("troll")
			}
			npc.Location = "" // Troll leaves the dungeon

			// Treasure is consumed (troll eats it)
			delete(g.Items, item.ID)

			return "The troll, who is not overly proud, graciously accepts the gift and not having the most discriminating tastes, gleefully eats it.\n\nThe troll, satiated, contentedly waddles off into the darkness, his axe clattering to the floor. The passages are now open."
		}

		// Non-treasure items
		// In ZIL, troll accepts anything but only leaves for treasures
		delete(g.Items, item.ID) // Troll eats it anyway
		return "The troll, who is not overly proud, graciously accepts the gift and not having the most discriminating tastes, gleefully eats it.\n\nHowever, the troll is still blocking the passages."

	case "cyclops":
		// Cyclops puzzle (CYCLOPS-FCN in ZIL lines 1543-1574)
		// Two-part puzzle: 1) give lunch (hot peppers) 2) give water to put him to sleep
		if item.ID == "lunch" {
			// Give hot peppers - makes cyclops thirsty but doesn't solve puzzle
			delete(g.Items, item.ID)
			return "The cyclops says \"Mmm Mmm. I love hot peppers! But oh, could I use a drink. Perhaps I could drink the blood of that thing.\" From the gleam in his eye, it could be surmised that you are \"that thing\"."
		}
		if item.ID == "water" || (item.ID == "bottle" && g.Items["water"] != nil && g.Items["water"].Location == "bottle") {
			// Give water - cyclops drinks and falls asleep, sets CYCLOPS-FLAG
			// Only works after giving hot peppers (not checking in this simplified version)
			delete(g.Items, "water")
			// Put empty bottle back in room
			bottle := g.Items["bottle"]
			if bottle != nil {
				bottle.Location = g.Location
				bottle.Flags.IsOpen = true
				room := g.Rooms[g.Location]
				if room != nil {
					room.AddItem("bottle")
				}
			}
			// Cyclops falls asleep
			g.Flags["cyclops-flag"] = true
			npc.Flags.IsAggressive = false
			npc.Flags.CanFight = false
			return "The cyclops takes the bottle, checks that it's open, and drinks the water. A moment later, he lets out a yawn that nearly blows you over, and then falls fast asleep (what did you put in that drink, anyway?)."
		}
		if item.ID == "garlic" {
			return "The cyclops may be hungry, but there is a limit."
		}
		// Default for other items
		return "The cyclops is not so stupid as to eat THAT!"

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
			// Troll drops axe and vanishes (TROLL-FCN F-DEAD in ZIL)
			g.Flags["troll-dead"] = true
			room := g.Rooms[npc.Location]

			// Troll drops his axe
			axe := g.Items["axe"]
			if axe != nil {
				axe.Location = npc.Location
				if room != nil {
					room.AddItem("axe")
				}
			}

			// Remove troll from room
			if room != nil {
				room.RemoveNPC("troll")
			}
			result.WriteString("Almost as soon as the troll breathes his last breath, a cloud of sinister black fog envelops him, and when the fog lifts, the carcass has disappeared.\n\nThe troll's axe clatters to the floor.")

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

	// Special case: waving sceptre (SCEPTRE-FUNCTION in ZIL lines 2592-2619)
	if item.ID == "sceptre" {
		// At Aragain Falls or End of Rainbow
		if g.Location == "aragain-falls" || g.Location == "end-of-rainbow" {
			if !g.Flags["rainbow-flag"] {
				// Solidify the rainbow
				g.Flags["rainbow-flag"] = true

				// Make pot-of-gold visible
				pot := g.Items["pot-of-gold"]
				if pot != nil {
					pot.Flags.IsInvisible = false
				}

				result := "Suddenly, the rainbow appears to become solid and, I venture, walkable (I think the giveaway was the stairs and bannister)."

				// Extra message if at end-of-rainbow and pot is there
				if g.Location == "end-of-rainbow" && pot != nil && pot.Location == "end-of-rainbow" {
					result += "\n\nA shimmering pot of gold appears at the end of the rainbow."
				}

				return result
			} else {
				// Make rainbow insubstantial again
				g.Flags["rainbow-flag"] = false

				// If anyone is on the rainbow, they fall!
				// For now, just make it insubstantial
				return "The rainbow seems to have become somewhat run-of-the-mill."
			}
		}

		// On the rainbow itself - DEADLY!
		if g.Location == "on-rainbow" {
			g.Flags["rainbow-flag"] = false
			g.GameOver = true
			return "The structural integrity of the rainbow is severely compromised, leaving you hanging in midair, supported only by water vapor. Bye.\n\n****  You have died  ****"
		}

		// Anywhere else
		return "A dazzling display of color briefly emanates from the sceptre."
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
		// Trap door is already in the room (global object), just needs to be revealed
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

	// Special case: bell ceremony (LLD-ROOM in ZIL 1actions.zil lines 1083-1101)
	if item.ID == "bell" {
		// Check if we're at entrance-to-hades and ceremony not complete
		if g.Location == "entrance-to-hades" && !g.Flags["LLD-FLAG"] {
			// Ring bell - Step 1 of ceremony
			g.Flags["XB"] = true
			g.Flags["bell-ceremony-active"] = true
			g.Flags["bell-ceremony-turn"] = true

			// Bell becomes hot and drops
			item.Location = g.Location
			room := g.Rooms[g.Location]
			if room != nil {
				room.AddItem("bell")
			}
			// Remove from inventory
			for i, id := range g.Player.Inventory {
				if id == "bell" {
					g.Player.Inventory = append(g.Player.Inventory[:i], g.Player.Inventory[i+1:]...)
					break
				}
			}

			// If player has candles, they drop too
			result := `The bell suddenly becomes red hot and falls to the ground. The wraiths, as if paralyzed, stop their jeering and slowly turn to face you. On their ashen faces, the expression of a long-forgotten terror takes shape.`

			if g.hasItemInInventory("candles") {
				candles := g.Items["candles"]
				if candles != nil {
					candles.Location = g.Location
					candles.Flags.IsLit = false
					if room != nil {
						room.AddItem("candles")
					}
					for i, id := range g.Player.Inventory {
						if id == "candles" {
							g.Player.Inventory = append(g.Player.Inventory[:i], g.Player.Inventory[i+1:]...)
							break
						}
					}
					result += "\nIn your confusion, the candles drop to the ground (and they are out)."
				}
			}

			return result
		}

		// Normal bell ringing elsewhere
		return "Ding, dong. The bell echoes throughout the dungeon."
	}

	return "How does one ring a " + item.Name + "?"
}

// handleExorcise attempts to exorcise spirits (V-EXORCISE in ZIL gverbs.zil lines 643+)
func (g *GameV2) handleExorcise(objName string) string {
	// Check if we're at entrance-to-hades
	if g.Location != "entrance-to-hades" {
		return "There is nothing here to exorcise."
	}

	// Check if spirits are already banished
	if g.Flags["LLD-FLAG"] {
		return "The spirits have already been banished."
	}

	// Check if player has all three items (GHOSTS-F and LLD-ROOM in ZIL)
	hasBook := g.hasItemInInventory("book")
	hasBell := g.hasItemInInventory("bell")
	hasCandles := g.hasItemInInventory("candles")

	if hasBook && hasBell && hasCandles {
		return "You must perform the ceremony."
	}

	return "You aren't equipped for an exorcism."
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

	// Special case: praying while holding coffin (coffin puzzle)
	// Sets coffin-cure flag allowing passage down from south-temple
	for _, itemID := range g.Player.Inventory {
		if itemID == "coffin" {
			g.Flags["coffin-cure"] = true
			return "Your prayer is answered! You feel a lightness, as if a burden has been lifted."
		}
	}

	// Special case: in temple areas
	if strings.Contains(room.ID, "temple") || strings.Contains(room.ID, "altar") {
		return "Your prayer is heard, but not answered."
	}

	return "If you pray enough, your prayers may be answered."
}

// handleOdysseus handles saying "ulysses" or "odysseus" (V-ODYSSEUS in ZIL lines 945-961)
func (g *GameV2) handleOdysseus() string {
	// Only works in cyclops-room with cyclops present and awake
	if g.Location != "cyclops-room" {
		return "Wasn't he a sailor?"
	}

	// Check if cyclops is present and awake
	room := g.Rooms[g.Location]
	if room == nil {
		return "Wasn't he a sailor?"
	}

	cyclopPresent := false
	for _, npcID := range room.NPCs {
		if npcID == "cyclops" {
			cyclopPresent = true
			break
		}
	}

	if !cyclopPresent {
		return "Wasn't he a sailor?"
	}

	// If cyclops is already asleep, can't use the word
	if g.Flags["cyclops-flag"] {
		return "No use talking to him. He's fast asleep."
	}

	// Cyclops flees! This knocks down the east wall and opens passage to strange-passage
	g.Flags["cyclops-flag"] = true  // Cyclops is gone
	g.Flags["magic-flag"] = true     // East passage is now open

	// Remove cyclops from room
	room.RemoveNPC("cyclops")
	cyclops := g.NPCs["cyclops"]
	if cyclops != nil {
		cyclops.Flags.IsAlive = false
		cyclops.Flags.CanFight = false
		cyclops.Location = ""
	}

	return "The cyclops, hearing the name of his father's deadly nemesis, flees the room by knocking down the wall on the east of the room."
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

	// Special case: rubbing mirrors (MIRROR-MIRROR in ZIL lines 971-1012)
	if (item.ID == "mirror-1" || item.ID == "mirror-2") && !g.Flags["mirror-mung"] {
		// Determine the two rooms
		room1 := "mirror-room-1"
		room2 := "mirror-room-2"

		// Determine which room we're in and which is the other
		var fromRoom, toRoom string
		if g.Location == room1 {
			fromRoom = room1
			toRoom = room2
		} else if g.Location == room2 {
			fromRoom = room2
			toRoom = room1
		} else {
			return "You feel nothing unexpected."
		}

		// Swap ALL items between the two rooms (excluding mirrors and NPCs)
		fromRoomObj := g.Rooms[fromRoom]
		toRoomObj := g.Rooms[toRoom]

		if fromRoomObj == nil || toRoomObj == nil {
			return "You feel nothing unexpected."
		}

		// Collect items to move (excluding mirrors themselves)
		var fromItems []string
		for _, itemID := range fromRoomObj.Contents {
			if itemID != "mirror-1" && itemID != "mirror-2" {
				fromItems = append(fromItems, itemID)
			}
		}

		var toItems []string
		for _, itemID := range toRoomObj.Contents {
			if itemID != "mirror-1" && itemID != "mirror-2" {
				toItems = append(toItems, itemID)
			}
		}

		// Move items from fromRoom to toRoom
		for _, itemID := range fromItems {
			if item := g.Items[itemID]; item != nil {
				fromRoomObj.RemoveItem(itemID)
				item.Location = toRoom
				toRoomObj.AddItem(itemID)
			}
		}

		// Move items from toRoom to fromRoom
		for _, itemID := range toItems {
			if item := g.Items[itemID]; item != nil {
				toRoomObj.RemoveItem(itemID)
				item.Location = fromRoom
				fromRoomObj.AddItem(itemID)
			}
		}

		// Teleport player to the other room
		g.Location = toRoom

		return "There is a rumble from deep within the earth and the room shakes.\n\n" + g.handleLook()
	}

	return "You feel nothing unexpected."
}

// handleBreak breaks/smashes something (V-MUNG in ZIL)
func (g *GameV2) handleBreak(objName string) string {
	if objName == "" {
		return "Break what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	// Special case: breaking mirrors (MIRROR-MIRROR in ZIL lines 1003-1012)
	if item.ID == "mirror-1" || item.ID == "mirror-2" {
		if g.Flags["mirror-mung"] {
			return "Haven't you done enough damage already?"
		}

		// Break the mirror
		g.Flags["mirror-mung"] = true
		g.Flags["lucky"] = false

		return "You have broken the mirror. I hope you have a seven years' supply of good luck handy."
	}

	// Default: can't break most things
	return "You can't break that."
}

// handleBurn burns something (BLACK-BOOK burn handling in ZIL lines 2201-2205)
func (g *GameV2) handleBurn(objName string) string {
	if objName == "" {
		return "Burn what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	// Special case: burning the prayer book is DEADLY (BLACK-BOOK in ZIL)
	if item.ID == "book" {
		// Remove the book
		if item.Location == "inventory" {
			for i, id := range g.Player.Inventory {
				if id == item.ID {
					g.Player.Inventory = append(g.Player.Inventory[:i], g.Player.Inventory[i+1:]...)
					break
				}
			}
		} else {
			room := g.Rooms[item.Location]
			if room != nil {
				room.RemoveItem(item.ID)
			}
		}
		item.Location = "REMOVED"

		// Game over!
		g.GameOver = true
		return `A booming voice says "Wrong, cretin!" and you notice that you have turned into a pile of dust. How, I can't imagine.

****  You have died  ****`
	}

	// Default: can't burn most things
	return "You can't burn that."
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

// handleSave saves the current game state
func (g *GameV2) handleSave(cmd *Command) string {
	// Get filename from direct object (if provided)
	filename := cmd.DirectObject

	// Save the game
	if err := g.Save(filename); err != nil {
		return fmt.Sprintf("Failed to save game: %s", err)
	}

	// Get the actual save path to show user
	if filename == "" {
		filename = fmt.Sprintf("gork_save_%s.json", time.Now().Format("20060102_150405"))
	}
	if filepath.Ext(filename) != ".json" {
		filename += ".json"
	}

	savePath, _ := GetSavePath(filename)
	return fmt.Sprintf("Game saved to: %s", savePath)
}

// handleRestore restores a saved game state
func (g *GameV2) handleRestore(cmd *Command) string {
	filename := cmd.DirectObject

	// If no filename provided, list available saves
	if filename == "" {
		saves, err := ListSaves()
		if err != nil {
			return fmt.Sprintf("Failed to list saves: %s", err)
		}

		if len(saves) == 0 {
			return "No saved games found."
		}

		result := "Available saved games:\n"
		for i, save := range saves {
			result += fmt.Sprintf("  %d. %s\n", i+1, save)
		}
		result += "\nUse 'restore <filename>' to load a save."
		return result
	}

	// Restore the game
	if err := g.Restore(filename); err != nil {
		return fmt.Sprintf("Failed to restore game: %s", err)
	}

	// Return the current room description after restore
	return fmt.Sprintf("Game restored.\n\n%s", g.handleLook())
}

// handleEnter handles the ENTER command (V-ENTER in ZIL)
// ENTER alone tries to go IN
// ENTER <object> tries to go through/board the object
func (g *GameV2) handleEnter(cmd *Command) string {
	if cmd.DirectObject == "" {
		// ENTER alone = try to go IN
		return g.handleMove("in")
	}

	// ENTER <object> = try to go through it or board it
	// This is like V-THROUGH in ZIL
	item := g.findItem(cmd.DirectObject)
	if item == nil {
		return "You can't see any " + cmd.DirectObject + " here."
	}

	// For now, just try to move through it as a direction
	// In full ZIL this would handle boats, vehicles, etc.
	return "You can't enter that."
}

// handleThrow handles throwing objects (V-THROW in ZIL)
func (g *GameV2) handleThrow(cmd *Command) string {
	if cmd.DirectObject == "" {
		return "Throw what?"
	}

	// Find the object
	item := g.findItem(cmd.DirectObject)
	if item == nil {
		return "You don't have that."
	}

	if item.Location != "inventory" {
		return "You're not holding the " + item.Name + "."
	}

	// Special case: throwing at something
	if cmd.IndirectObject != "" {
		return fmt.Sprintf("The %s bounces harmlessly off the %s.", item.Name, cmd.IndirectObject)
	}

	// Just drop it
	return g.handleDrop(cmd.DirectObject)
}

// handleYell handles yelling (V-YELL in ZIL)
func (g *GameV2) handleYell() string {
	// Check if we're in the echo room (canyon-view in Zork I)
	if g.Location == "canyon-view" {
		return "Your voice echoes back: \"HELLO!\""
	}

	return "You scream loudly. Nothing happens."
}

// handleBoard handles boarding vehicles (V-BOARD in ZIL)
func (g *GameV2) handleBoard(objName string) string {
	if objName == "" {
		return "Board what?"
	}

	item := g.findItem(objName)
	if item == nil {
		return "You can't see any " + objName + " here."
	}

	// Check if it's the boat
	if objName == "boat" || objName == "raft" {
		if item.Location != g.Location && item.Location != "inventory" {
			return "The boat isn't here."
		}
		return "You are now in the boat."
	}

	return "You can't board that."
}

// handleDiagnose handles the DIAGNOSE command (V-DIAGNOSE in ZIL)
func (g *GameV2) handleDiagnose() string {
	// In original Zork, this would report health status
	// For now, simple implementation
	return "You are in perfect health."
}

// handleSay handles SAY command (V-SAY in ZIL)
func (g *GameV2) handleSay(cmd *Command) string {
	if cmd.DirectObject == "" {
		return "Say what?"
	}

	word := strings.ToLower(cmd.DirectObject)

	// Magic words
	switch word {
	case "xyzzy", "plugh":
		return "A hollow voice says \"Fool.\""
	case "hello":
		return g.handleYell()
	default:
		return "Nothing happens."
	}
}

// handleEcho handles ECHO command (V-ECHO in ZIL)
func (g *GameV2) handleEcho() string {
	return g.handleYell()
}
