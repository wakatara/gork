# GORK - Zork I in Go

> You'll likely be eaten by a grue

An somewhat faithful port of the classic 1980s text adventure game **Zork I: The
Great Underground Empire** from Infocom to Go from the original variant of a
variant of lisp ZIL (Zork Implementation Language) [Zork I ZIL code
here](https://github.com/historicalsource/zork1). If you care about the how and
changes I made from the ZIL Code, you can find details in [[./ARCHITECTURE.md]]

## About

Ah, nostalgia... I remember trying to play Zork in the "computer store" in the
small town I grew up in well before I could even dream of owning my own actual
computer. I was eaten repeatedly by grue (grues? gruen? gri?).

Zork was incredibly influential as a piece of interactive fiction and to
computer gaming, and instrumental in Infocom being by Activision. It was also a
runaway commercial success for video games at that time.

I never really got a chance to play the whole thing. I didn't even have my own
computer (a Vic 20!) until years later. And when the nice people at
[nixCraft](https://mastodon.social/@nixCraft/115585413690855037) posted on
Mastodon about the original source code being open sourced. Well, I honestly
wodnered how hard it would be to port (and have been doing a lot of game dev
lately outside of scientific work, anyway.). And kinda wanted to play the whol
thing through, ya know?

The idea is to recreate the original game experience (as I jazily remember it -
yell if you see I've messed something up) and just make it fun and available for
people to try and play through (now, I am wondering about the old Hitchhiker's
Guide to the Galaxy text games and what happened to them. 🤔).

I had to take some serious liberties with the original ZIL source code since
_everything_ was an non-type-safe object, so translating to modern computing
constructs for my own sanity took a bit of overarching redesign though the
gameplay _should be_ completely faithful to the original. The ZIL fiels have no
tests, and while I've created as many as I could, I am flying a bit blind here
(so please playtest the heck outta this thing for me). The parser was trickier
than you'd think despite the simplicity of the language commands. Other than the
underlying internals though, and some affectations with intentionally trying to
make it "CRT"-y with terminal effects, I _think_ I managed to stick the landing.

I also ignored the original save and restore custom text format in favour of
something more modern (json) and Go idiomatic. Also, as a quality of life
addition, there is a `clear` and `cls` command I added to keep everything on one
screen if you choose. But, other than those, yeah... it should be a rather
faitful rendition of Zork I.

Please try it out and let me know what I may have gotten worong. I'm sure there
must be bugs - even with using an AI to list and trace all room exissts and
entrances, I ran across from serious DAG problems early on, though think I
nailed all of those. I'm really hoping there's someone who actually played the
whole thing through that can give feedback on the port with some proper
playtesting.

_Anyways, kinda an early 2025 Xmas gift back to the world, if you will._

**Have fun and please don't feed the grues.**

## Features

- ✅ **Somewhat Sophisticated Parser**: Natural language command processing
  - Handles complex commands like "put letter in mailbox", "look at white house"
  - Multi-word object resolution ("kitchen window", "white house")
  - "it" reference tracking (the last mentioned object)
  - ~685 vocabulary words from the original game
  - understands synonyms (and multiple terms for the same object or verb action)

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
  - Famous locations and puzzles: West of House, Trolls, Maze, Cyclops, Canyon
  - Conditional exits (aka puzzles)
  - Light/dark room mechanics with grue warning
  - Outdoor/indoor/underground areas
  - 122 items across 10 categories
  - All 19 treasures, weapons, tools, and containers
  - Readable items
  - And, of course, container maniuplation (open, close, look in, light etc)

- ✅ **Retro Terminal UI**:
  - Optional character-by-character typing effect
  - Amber or green CRT monitor color themes
  - ASCII art title screen
  - Classic "> " prompt

## Installation

### Homebrew (macOS/Linux)

```sh
brew tap wakatara/tap
brew install gork
```

### Working Go Environment

If you have Go installed, the simplest method:

```sh
go install github.com/wakatara/gork/cmd/gork@latest
```

This will install `gork` to your `$GOPATH/bin` directory. Make sure that directory is in your PATH.

### Pre-built Binaries

Download binaries for your platform from the [releases page](https://github.com/wakatara/gork/releases).

Available formats:

- `.tar.gz` / `.zip` archives (all platforms)
- `.deb` packages (Debian/Ubuntu)
- `.rpm` packages (Red Hat/Fedora/CentOS)

### Build from Source

```sh
git clone https://github.com/wakatara/gork
cd gork
go build -o gork ./cmd/gork

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

- **Movement**: `north` (`n`), `south` (`s`), `east` (`e`), `west` (`w`), `up`
  (`u`), `down` (`d`), `in`, `out`
- **Examination**: `look` (`l`), `examine <object>` (`x <object>`), `read
<object>`
- **Manipulation**: `take <object>` (`get`), `drop <object>`, `open/close
<container>`
- **Light**: `turn on/off <light>` - Control light sources
- **Containers**: `look in <container>` - Inspect contents
- **Interaction**: `put <obj> in <container>`, `give <obj> to <npc>`
- **System**: `inventory` (`i`), `help`, `save`, `restore`, `quit`

### Items in the Game

122 items across 10 categories (ALL items from original Zork I):

- **Treasures** (19): diamond, emerald, chalice, jade, coins, painting,
  bracelet, scarab, sceptre, egg, pot-of-gold, trident, bauble, platinum bar,
  sapphire, ivory torch, trunk of jewels, pearl, oriental rug
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

### Testing

**Important:** Always run tests from the project root directory with `./...` to
test all packages:

```bash
# Run all tests (from project root)
go test ./...

# Run all tests with verbose output
go test ./... -v

# Run only engine tests with verbose output
go test ./engine -v

# Run specific test
go test ./engine -run TestParser

# Run with coverage report
go test ./engine -cover
```

**Note:** Running `go test -v` without `./...` from the root will fail because the root directory has no Go files. Always use `go test ./...` to test all packages recursively.

Current test coverage:

- `engine/parser_test.go`: Parser and command processing
- `engine/rooms_test.go`: All 110 rooms and navigation
- `engine/items_test.go`: All 122 items
- `engine/game_verbs_test.go`: Verb handlers
- `engine/new_verbs_test.go`: New verb handlers
- `engine/puzzle_integration_test.go`: Major puzzle sequences
- `engine/save_test.go`: Save/restore functionality
- **Total: 200 tests passing**

All tests use `t.Logf()` for debug output, so verbose mode (`-v`) will show
detailed test information with proper PASS/FAIL indicators.

## Architecture

```
gork/
├── engine/           # Core game engine
│   ├── parser.go     # Natural language parser
│   ├── vocabulary.go # 684 words from original Zork
│   ├── types.go      # Room, Item, NPC type definitions
│   ├── game_v2.go    # Game state and command execution
│   ├── rooms.go      # All 110 rooms from Zork I
│   └── items.go      # All 122 items from Zork I
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

See [ARCHITECTURE.md](ARCHITECTURE.md) for the complete design rationale behind
choosing separate structs over ZIL's generic object model.

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

MIT License - see [LICENSE](LICENSE) file for details.

This is a modern Go port of the classic 1980s text adventure game "Zork I: The
Great Underground Empire" originally created by Infocom, Inc. The game's source
code in ZIL (Zork Implementation Language) was released as open source by
Microsoft Corporation in 2025 under the MIT License for historical preservation
purposes.

This implementation is also released under the MIT License, with proper
attribution to both the original creators and the ZIL source release.

## Acknowledgments

- Marc Blank, Dave Lebling, Bruce Daniels, and Tim Anderson - Original Zork creators (1977-1979)
- Infocom, Inc. - For creating the ZIL language and publishing Zork I (1980)
- Microsoft Corporation - For releasing the ZIL source code under MIT License (2025)
- The interactive fiction community - For preserving these classics

---

_"You **are** likely to be eaten by a grue."_
