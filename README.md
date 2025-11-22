# GORK - Zork I in Go

> You'll likely be eaten by a grue

A somewhat faithful port of the classic 1980s text adventure game **Zork I: The
Great Underground Empire** from Infocom to Go from the original variant of a
variant of lisp ZIL (Zork Implementation Language) files located here from the
[original Zork I ZIL code here](https://github.com/historicalsource/zork1)

## About

Ah, nostalgia... I remember trying to play Zork in the "computer store" in the
small town I grew up in well before I could even dream of owning my own actual
computer. I was eaten repeatedly by several grue (grues? gruen? gri?).

Zork was amazingly influential as a piece of interactive fiction and to computer
gaming in general, and probably what got Infocom purchased by Activision (as
well as being - for the time - a huge commercial success).

I never really got a chance to play the whole thing (or any substantive playtime
at all, to be honest — also, suddenly games with actual graphical interfaces
exploded on the scene which my pre-teen brain was powerless to resist), and was
curious how hard this would be to actually port since I've been porting a lot of
old scientific libraries in Astronomy lately, and it's been rather surprisingly
educational. Also, I kinda wanted to play the whole thing through, ya know?

The idea is to recreate the original game experience (as I hazily and
imperfectly remember it - so please chime in if you see I've messed something
up) and just make it fun and available for people to try. Kinda a 2025 Xmas gift
to the world, if you will.

I had to take some liberties with the original ZIL source code since
_everything_ was an object, so translating to modern computing constructs for my
own sanity took a bit of overarching redesign. The parser was trickier than
you'd think despite the simplicity of the language commands. Other than the
underlying internals though, and some affectations with intentionally trying to
make it "CRT"-y and terminal effects, I _think_ I managed to stick the landing.

Please try it out and let me know what I may have gotten worong. I'm hoping
there's someone who actually played the whole thing through that can give
feedback on the port.

Have fun and please don't get eaten by a grue.

## Features

- ✅ **Sophisticated Parser**: Natural language command processing with synonym support
  - Handles complex commands like "put sword in case", "look at white house"
  - Multi-word object resolution ("kitchen window", "white house")
  - "IT" reference tracking (the last mentioned object)
  - 684 vocabulary words from the original game

- ✅ **Complete Synonym Support**: All classic shortcuts work
  - Movement: `n`, `s`, `e`, `w`, `ne`, `nw`, `se`, `sw`, `u`, `d`
  - Verbs: `x` (examine), `l` (look), `i` (inventory)
  - Actions: `get`/`take`, `grab`/`pick up`, `drop`/`put down`

- ✅ **Type-Safe Go Architecture**: Modern struct-based design
  - Separate `Room`, `Item`, and `NPC` types (not generic objects)
  - Compiler-enforced type safety
  - Memory efficient with only relevant fields per type
  - Better IDE support and maintainability
  - See [ARCHITECTURE.md](ARCHITECTURE.md) for design rationale

- ✅ **Complete Game World**:
  - All 110 rooms from original Zork I
  - Famous locations: West of House, Troll Room, Maze, Cyclops Room, Aragain Falls
  - Conditional exits (trap door, troll blocking passages)
  - Light/dark room mechanics with grue warning
  - Outdoor/indoor/underground areas

- ✅ **Comprehensive Item System**:
  - 85+ items across 10 categories
  - All treasures, weapons, tools, containers
  - Light sources with on/off control
  - Readable items (books, maps, leaflets)
  - Container manipulation (open, close, look in)

- ✅ **Retro Terminal UI**:
  - Optional character-by-character typing effect
  - Amber or green CRT monitor color themes
  - ASCII art title screen
  - Classic "> " prompt

- 🚧 **Game World** (~50-55% Complete):
  - All rooms and most items implemented
  - Core exploration and item manipulation works
  - NPCs and combat system in progress
  - Puzzles and special handlers planned

## Installation

**WIP**: My plan is to make this available via a GoReleaser GH build pipeline
and available as compiled binaries and at least on the homebrew package manager.
Other package managers perhaps - depending on how keen folks are.

```bash
# Clone the repository
cd gork

# Build the game
go build -o gork ./cmd/gork

# Run
./gork
```

## Playing

```bash
> look
West of House
You are standing in an open field west of a white house, with a boarded front door.
There is a small mailbox here.

> examine mailbox
The small mailbox is open.

> take leaflet
Taken.

> read leaflet
"WELCOME TO ZORK!

ZORK is a game of adventure, danger, and low cunning..."

> inventory
You are carrying:
  A leaflet
```

### Common Commands

- **Movement**: `north` (`n`), `south` (`s`), `east` (`e`), `west` (`w`), `up` (`u`), `down` (`d`), `in`, `out`
- **Examination**: `look` (`l`), `examine <object>` (`x <object>`), `read <object>`
- **Manipulation**: `take <object>` (`get`), `drop <object>`, `open/close <container>`
- **Light**: `turn on/off <light>` - Control light sources
- **Containers**: `look in <container>` - Inspect contents
- **Interaction**: `put <obj> in <container>`, `give <obj> to <npc>`
- **System**: `inventory` (`i`), `help`, `save`, `restore`, `quit`

### Items in the Game

85+ items across 10 categories:

- **Treasures** (13): diamond, emerald, chalice, jade, coins, painting, bracelet, scarab, etc.
- **Weapons** (5): elvish sword, knife, axe, stiletto, trident
- **Tools** (6): pump, screwdriver, wrench, rope, shovel, putty
- **Containers** (11): mailbox, trophy case, bottle, coffin, nest, bags
- **Light Sources** (4): brass lantern, torch, candles
- **Readable** (8): leaflet, prayer book, guidebook, maps, manuals
- **Food/Drink** (3): lunch, garlic, water
- **Fixed Objects** (14): windows, doors, buttons, grating, altar, machine
- **Scenery** (8): white house, forest, mountains, rainbow, river, engravings
- **Miscellaneous** (13): boats, skull, bones, coal, ladder, canary, etc.

### Tips

- The lamp is crucial - you can't explore in the dark or you'll be eaten by a grue!
- Try examining everything
- The white house has a boarded front door, but perhaps there's another way in?
- The original Zork I has 19 treasures worth 350 points total

## Development Status

| Feature          | Status                             |
| ---------------- | ---------------------------------- |
| Parser           | ✅ Complete (106 tests passing)    |
| Synonyms         | ✅ Complete (684 vocabulary words) |
| Type System      | ✅ Complete (Room/Item/NPC)        |
| Game Engine      | ✅ Core complete                   |
| Rooms            | ✅ 110/110 rooms (100%)            |
| Items            | 🚧 85/122 items (70%)              |
| Treasures        | 🚧 13/19 treasures (70%)           |
| Verb Handlers    | 🚧 Core verbs implemented (~60%)   |
| Light/Darkness   | 🚧 Basic implementation (80%)      |
| **Overall Game** | **🚧 ~50-55% complete**            |
| NPCs             | 🚧 In progress (1/5 implemented)   |
| Combat System    | ⏸️ Planned (0%)                    |
| Puzzles          | ⏸️ Planned (0%)                    |
| Save/Restore     | ⏸️ Planned (0%)                    |
| Score System     | ⏸️ Planned (0%)                    |

### Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./engine -v

# Run specific test
go test ./engine -run TestParser
```

Current test coverage:

- `engine/parser_test.go`: 67 test cases covering command parsing
- `engine/rooms_test.go`: 21 test cases for all 110 rooms
- `engine/items_test.go`: 10 test cases for 85+ items
- `engine/game_verbs_test.go`: 8 test cases for verb handlers
- **Total: 106 tests passing, ~80% coverage**

## Architecture

```
gork/
├── engine/           # Core game engine
│   ├── parser.go     # Natural language parser
│   ├── vocabulary.go # 684 words from original Zork
│   ├── types.go      # Room, Item, NPC type definitions
│   ├── game_v2.go    # Game state and command execution
│   ├── rooms.go      # All 110 rooms from Zork I
│   └── items.go      # All 85+ items from Zork I
├── ui/               # Terminal UI and retro effects
└── cmd/gork/         # Main entry point
```

## Implementation Notes

### Parser

The parser implements ZIL's sophisticated natural language processing:

```go
// Handles complex commands
"put sword in case"        → verb: put, direct: sword, prep: in, indirect: case
"look at white house"      → verb: examine, direct: white-house
"take it"                  → verb: take, direct: [last referenced object]
```

### Type System

Unlike ZIL where everything was a generic `OBJECT`, we use modern Go structs:

```go
type Room struct {
    ID, Name, Description string
    Exits       map[string]*Exit
    Contents    []string  // Item IDs
    NPCs        []string  // NPC IDs
    Flags       RoomFlags
}

type Item struct {
    ID, Name string
    Aliases     []string
    Description string
    Location    string
    Flags       ItemFlags
    Value       int  // For treasures
}

type NPC struct {
    ID, Name, Description string
    Location    string
    Strength    int
    Weapon      string
    Hostile     bool
}
```

See [ARCHITECTURE.md](ARCHITECTURE.md) for the complete design rationale behind choosing separate structs over ZIL's generic object model.

### Action Dispatch

Following ZIL's priority chain from `gmain.zil`:

1. Player/Actor action handler
2. Room's M-BEG handler
3. Preaction handlers
4. Indirect object handler
5. Direct object handler
6. Default verb handler

## Original Source

This port is based on the original Zork I source code:

- Repository: [historicalsource/zork1](https://github.com/historicalsource/zork1)
- Language: ZIL (Zork Implementation Language)
- Total Source: ~12,000 lines of ZIL
- Vocabulary: 684 words
- Objects: 122
- Rooms: 110

## Retro Features

Want that authentic 1980s terminal experience?

```go
// In main.go, add:
ui.EnableRetroMode()
```

This enables:

- Character-by-character typing effect
- Amber CRT monitor color theme
- Authentic terminal feel with slight random delays

## License

The original Zork I source code was released as open source by
Infocom/Activision for historical preservation. This Go port is created for
educational purposes and preserving gaming history.

## Acknowledgments

- Marc Blank, Dave Lebling, Bruce Daniels, and Tim Anderson - Original Zork creators
- Infocom - For creating the ZIL language and Zork
- The interactive fiction community - For preserving these classics

## Progress

- ✅ ~~Full dungeon implementation (~110 rooms)~~ **Complete - 110 rooms**
- ✅ ~~Light and darkness (and the grue!)~~ **Complete - hasLight() with grue warning**
- ✅ ~~Commands and synonyms (`x`, `i`, `n`, `s`, etc.)~~ **Complete - 684 words**
- 🚧 All 19 treasures (13/19 implemented - 68% complete)
- 🚧 Complete item set (85/122 items - 70% complete)
- [ ] NPCs and combat system (troll, thief, cyclops, bat, ghosts)
- [ ] Puzzle special handlers (dam controls, mirror room, basket/rope, machine)
- [ ] Save/restore functionality
- [ ] Score tracking (0/350 points system)
- [ ] Death and resurrection system

---

_"You **are** likely to be eaten by a grue."_
