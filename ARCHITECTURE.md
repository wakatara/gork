# GORK Architecture

## Design Decision: Separate Structs vs. Generic Objects

### The Problem

The original Zork was written in ZIL (a dialect of a Lisp dialect), which used a
generic `OBJECT` type for everything (rooms, items, NPCs). In Lisp, this makes
sense because:

- Lisp has NO static typing
- _Everything_ is a list/atom
- Flexibility through properties

### Why I Chose Dedicated Types in Go

For the purists who might scream about what I've done to the port (though heck,
you could probably implement this whole thing in emacs):

**I chose separate structs (`Room`, `Item`, `NPC`) instead of a generic `Object`
because:**

✅ **Type Safety**

```go
// Before: Everything possible, nothing checked
obj.Exits         // Could exist on a lamp!
obj.Strength      // Makes no sense for a room!

// After: Compiler-enforced correctness
room.Exits        // ✓ Only rooms have exits
npc.Strength      // ✓ Only NPCs have strength
item.Exits        // ✗ Compile error!
```

✅ **Memory Efficiency**

- Generic Object: Every object needs ~20 fields (most unused)
- Separate types: Only relevant fields per type

✅ **Code Clarity**

```go
func enterRoom(room *Room) { ... }      // Clear intent
func takeItem(item *Item) { ... }        // Clear intent
func talkToNPC(npc *NPC) { ... }         // Clear intent
```

✅ **IDE Support**

- Autocomplete knows what fields are available (and I ❤️my inventoryim)
- Can't accidentally access wrong fields
- Better documentation

✅ **Maintainability**

- Easy to add room-specific features
- Easy to add NPC-specific behavior
- Changes don't affect other types

## Current Type Hierarchy

```
GameV2
├── Rooms:     map[string]*Room
├── Items:     map[string]*Item
├── NPCs:      map[string]*NPC
├── Player:    *Player
└── Flags:     map[string]bool (global game state)
```

### Room

```go
type Room struct {
    ID          string
    Name        string
    Description string
    Exits       map[string]*Exit
    Contents    []string // Item IDs
    NPCs        []string // NPC IDs
    Flags       RoomFlags
}
```

**Purpose:** Locations player can visit
**Examples:** west-of-house, kitchen, cellar, troll-room
**Key Features:**

- Directional exits (north, south, etc.)
- Can be lit or dark
- Can contain items and NPCs
- Conditional exits (e.g., trap door requires rug-moved flag)

### Item

```go
type Item struct {
    ID          string
    Name        string
    Aliases     []string
    Description string
    Location    string // Room ID or "inventory"
    Flags       ItemFlags
    Weight      int
    Value       int // For treasures
}
```

**Purpose:** Objects that can be manipulated
**Examples:** lamp, sword, mailbox, leaflet
**Key Features:**

- Can be takeable or fixed
- Can be containers (mailbox contains leaflet)
- Can be light sources (lamp)
- Can be treasures (for scoring)

### NPC

```go
type NPC struct {
    ID          string
    Name        string
    Description string
    Location    string // Room ID
    Strength    int
    Weapon      string  // Item ID
    Inventory   []string
    Hostile     bool
}
```

**Purpose:** Characters you can interact with
**Examples:** troll, thief, cyclops
**Key Features:**

- Can be hostile or friendly
- Have health/strength for combat
- Carry weapons and items
- Can move between rooms (AI)

## Game State Management

### Player Inventory

```go
type Player struct {
    Inventory []string // Item IDs
    MaxWeight int
    Health    int
}
```

### Global Flags

```go
g.Flags["troll-dead"]    = true
g.Flags["rug-moved"]     = true
g.Flags["window-open"]   = false
```

Used for:

- Quest state (troll defeated, treasure found)
- Puzzle state (rug moved, door unlocked)
- Conditional exits (can't go north if troll-dead is false)

## Room Connections

### Simple Exit

```go
room.AddExit("north", "north-of-house")
```

### Conditional Exit

```go
// Can only go down if rug-moved flag is true
livingRoom.AddConditionalExit("down", "cellar", "rug-moved", "You can't go that way.")
```

### Complex Example (Troll Room)

```go
trollRoom.AddExit("east", "east-west-passage")     // Always available
trollRoom.AddExit("west", "maze-1")                 // Always available
trollRoom.AddConditionalExit("north", "north-of-chasm", "troll-dead", "The troll blocks your way.")
```

## Light and Darkness

### Room Lighting

```go
room.Flags.IsLit = true   // Room provides light (outdoors, lit rooms)
room.Flags.IsDark = true  // Room is inherently dark (cellar, caves)
```

### Light Sources

```go
lamp.Flags.IsLightSource = true
lamp.Flags.IsLit = true  // Currently providing light
```

### Light Check

```go
func (g *GameV2) hasLight() bool {
    // 1. Check if current room is lit
    if g.Rooms[g.Location].Flags.IsLit {
        return true
    }
    // 2. Check for lit lamp in inventory
    for _, itemID := range g.Player.Inventory {
        if item.Flags.IsLightSource && item.Flags.IsLit {
            return true
        }
    }
    return false
}
```

### Darkness Behavior

```
> down
It is pitch black. You are likely to be eaten by a grue.
```

## Action Dispatch

When player enters command, the flow is:

1. **Parse** - Convert input to `Command` struct
2. **Execute** - Route to appropriate handler
3. **Update** - Modify game state
4. **Respond** - Return text to display

### Example: Taking an Item

```go
> take lamp

1. Parser: "take lamp" -> Command{Verb: "take", DirectObject: "lamp"}
2. Game: handleTake("lamp")
3. Find item in current room
4. Check if takeable
5. Move from room to inventory
6. Return: "Taken."
```

## File Structure

```
gork/
├── engine/
│   ├── types.go       # Room, Item, NPC, Exit types
│   ├── game_v2.go     # Game state and logic
│   ├── parser.go      # Natural language parsing
│   └── vocabulary.go  # Word database
├── world/
│   ├── rooms.go       # Room definitions (to be created)
│   ├── items.go       # Item definitions (to be created)
│   └── npcs.go        # NPC definitions (to be created)
├── ui/
│   └── terminal.go    # Display and colors
└── cmd/gork/
    └── main.go        # Entry point
```

## Next Steps

1. **Extract world definitions** - Move room/item/NPC creation from `game_v2.go` to `world/` package
2. **Add more rooms** - Port ~50 rooms from original ZIL
3. **Add more items** - Implement all items from original game
4. **Add more NPCs** - Troll, thief, cyclops behaviors
5. **Combat system** - Implement fighting mechanics
6. **Puzzle system** - Dam, mirrors, maze, etc.
7. **Save/restore** - Serialize game state

## Lessons Learned

**Don't blindly port old architectures.**

The original Zork used generic objects because Lisp had no type system. Go has
strong typing. Use it! (and useful modern consstructs in general!). The
refactored architecture is:

- More maintainable
- More performant
- More Go-idiomatic
- Easier to understand (at least to me)
- Safer (compiler catches bugs)

**Modern language features exist for a reason.**

Don't fight your language trying to emulate another one.
