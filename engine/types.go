package engine

// Room represents a location in the game world
type Room struct {
	ID          string
	Name        string            // Short name ("West of House")
	Description string            // Full description
	Exits       map[string]*Exit  // Direction -> Exit
	Contents    []string          // IDs of items in this room
	NPCs        []string          // IDs of NPCs in this room
	Flags       RoomFlags
	FirstVisit  bool              // True if never visited
	Action      RoomActionHandler // Custom room behavior
}

// RoomFlags holds boolean flags for rooms
type RoomFlags struct {
	IsLit       bool // Room provides light
	IsDark      bool // Room is inherently dark
	IsUnderwater bool
	IsOutdoors  bool
}

// RoomActionHandler handles room-specific events
type RoomActionHandler func(g *GameV2, event RoomEvent) string

type RoomEvent int

const (
	EventEnter RoomEvent = iota
	EventLook
	EventLeave
)

// Exit represents a connection between rooms
type Exit struct {
	To        string // Destination room ID
	Condition string // Optional: game flag that must be true
	Message   string // Message if exit is blocked
}

// Item represents a takeable object in the game
type Item struct {
	ID          string
	Name        string   // Primary name
	Aliases     []string // Alternative names (lamp/lantern)
	Description string   // What you see when examining
	Location    string   // Room ID or "inventory" or container ID
	Flags       ItemFlags
	Weight      int
	Value       int  // For treasures (score)
	Fuel        int  // For light sources (turns remaining, -1 = infinite)
	GlowLevel   int  // For sword: 0=not glowing, 1=faint, 2=bright
	Action      ItemActionHandler
}

// ItemFlags holds boolean flags for items
type ItemFlags struct {
	IsTakeable   bool
	IsContainer  bool
	IsOpen       bool // For containers
	IsTransparent bool // Can see contents when closed
	IsReadable   bool
	IsEdible     bool
	IsDrinkable  bool
	IsWeapon     bool
	IsLightSource bool
	IsLit        bool // For light sources that are on
	IsTreasure   bool
	IsWearable   bool
	IsInvisible  bool // For items that are initially invisible (like pot-of-gold)
}

// ItemActionHandler handles item-specific interactions
type ItemActionHandler func(g *GameV2, action string, item *Item) string

// NPC represents a non-player character
type NPC struct {
	ID          string
	Name        string
	Description string
	Location    string // Room ID
	Flags       NPCFlags
	Strength    int      // For combat
	Weapon      string   // Item ID of weapon
	Inventory   []string // Item IDs
	Hostile     bool
	Action      NPCActionHandler
}

// NPCFlags holds boolean flags for NPCs
type NPCFlags struct {
	IsAggressive bool
	IsFriendly   bool
	IsAlive      bool
	CanTalk      bool
	CanFight     bool
}

// NPCActionHandler handles NPC behavior
type NPCActionHandler func(g *GameV2, npc *NPC, action string) string

// NewRoom creates a new room with default values
func NewRoom(id, name, description string) *Room {
	return &Room{
		ID:          id,
		Name:        name,
		Description: description,
		Exits:       make(map[string]*Exit),
		Contents:    []string{},
		NPCs:        []string{},
		Flags:       RoomFlags{IsLit: true}, // Most rooms are lit
		FirstVisit:  true,
	}
}

// NewItem creates a new item with default values
func NewItem(id, name, description string) *Item {
	return &Item{
		ID:          id,
		Name:        name,
		Aliases:     []string{},
		Description: description,
		Flags:       ItemFlags{},
	}
}

// NewNPC creates a new NPC with default values
func NewNPC(id, name, description string) *NPC {
	return &NPC{
		ID:          id,
		Name:        name,
		Description: description,
		Inventory:   []string{},
		Flags: NPCFlags{
			IsAlive: true,
		},
	}
}

// Room helper methods

func (r *Room) AddExit(direction, destination string) {
	r.Exits[direction] = &Exit{To: destination}
}

func (r *Room) AddConditionalExit(direction, destination, condition, message string) {
	r.Exits[direction] = &Exit{
		To:        destination,
		Condition: condition,
		Message:   message,
	}
}

func (r *Room) AddItem(itemID string) {
	r.Contents = append(r.Contents, itemID)
}

func (r *Room) RemoveItem(itemID string) {
	for i, id := range r.Contents {
		if id == itemID {
			r.Contents = append(r.Contents[:i], r.Contents[i+1:]...)
			return
		}
	}
}

func (r *Room) AddNPC(npcID string) {
	r.NPCs = append(r.NPCs, npcID)
}

func (r *Room) RemoveNPC(npcID string) {
	for i, id := range r.NPCs {
		if id == npcID {
			r.NPCs = append(r.NPCs[:i], r.NPCs[i+1:]...)
			return
		}
	}
}

func (r *Room) HasNPC(npcID string) bool {
	for _, id := range r.NPCs {
		if id == npcID {
			return true
		}
	}
	return false
}

// Item helper methods

func (i *Item) HasAlias(name string) bool {
	if i.Name == name {
		return true
	}
	for _, alias := range i.Aliases {
		if alias == name {
			return true
		}
	}
	return false
}

func (i *Item) IsInRoom(roomID string) bool {
	return i.Location == roomID
}

func (i *Item) IsInInventory() bool {
	return i.Location == "inventory"
}

// NPC helper methods

func (n *NPC) IsInRoom(roomID string) bool {
	return n.Location == roomID
}

func (n *NPC) HasItem(itemID string) bool {
	for _, id := range n.Inventory {
		if id == itemID {
			return true
		}
	}
	return false
}

func (n *NPC) AddItem(itemID string) {
	n.Inventory = append(n.Inventory, itemID)
}

func (n *NPC) RemoveItem(itemID string) {
	for i, id := range n.Inventory {
		if id == itemID {
			n.Inventory = append(n.Inventory[:i], n.Inventory[i+1:]...)
			return
		}
	}
}
