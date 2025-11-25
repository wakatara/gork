package engine

// InitializeRooms creates all 110 rooms from Zork I
// Ported from 1dungeon.zil lines 1239-2660
func InitializeRooms(g *GameV2) {
	createAboveGroundRooms(g)
	createHouseRooms(g)
	createCellarAndVicinity(g)
	createMazeRooms(g)
	createCyclopsArea(g)
	createReservoirArea(g)
	createMirrorRooms(g)
	createRoundRoomArea(g)
	createHadesArea(g)
	createTempleArea(g)
	createDamArea(g)
	createRiverArea(g)
	createCoalMineArea(g)
}

// createAboveGroundRooms creates the outdoor area rooms
func createAboveGroundRooms(g *GameV2) {
	// WEST-OF-HOUSE - The famous starting location
	westOfHouse := NewRoom(
		"west-of-house",
		"West of House",
		"You are standing in an open field west of a white house, with a boarded front door.",
	)
	westOfHouse.Flags.IsOutdoors = true
	westOfHouse.AddExit("north", "north-of-house") // per ZIL
	westOfHouse.AddExit("south", "south-of-house") // per ZIL
	westOfHouse.AddExit("ne", "north-of-house") // per ZIL
	westOfHouse.AddExit("se", "south-of-house") // per ZIL
	westOfHouse.AddExit("west", "forest-1") // per ZIL
	// Note: EAST blocked by boarded door
	westOfHouse.AddConditionalExit("sw", "stone-barrow", "won-flag", "")
	westOfHouse.AddConditionalExit("in", "stone-barrow", "won-flag", "")
	g.Rooms["west-of-house"] = westOfHouse

	// STONE-BARROW - Victory area
	stoneBarrow := NewRoom(
		"stone-barrow",
		"Stone Barrow",
		"You are standing in front of a massive barrow of stone. In the east face is a huge stone door which is open. You cannot see into the dark of the tomb.",
	)
	stoneBarrow.Flags.IsOutdoors = true
	stoneBarrow.AddExit("ne", "west-of-house")
	g.Rooms["stone-barrow"] = stoneBarrow

	// NORTH-OF-HOUSE
	northOfHouse := NewRoom(
		"north-of-house",
		"North of House",
		"You are facing the north side of a white house. There is no door here, and all the windows are boarded up. To the north a narrow path winds through the trees.",
	)
	northOfHouse.Flags.IsOutdoors = true
	northOfHouse.AddExit("north", "path") // per ZIL
	northOfHouse.AddExit("west", "west-of-house") // per ZIL
	northOfHouse.AddExit("sw", "west-of-house") // per ZIL
	northOfHouse.AddExit("east", "behind-house") // per ZIL (EAST-OF-HOUSE)
	northOfHouse.AddExit("se", "behind-house") // per ZIL
	// Note: SOUTH blocked by house ("The windows are all boarded.")
	g.Rooms["north-of-house"] = northOfHouse

	// SOUTH-OF-HOUSE
	southOfHouse := NewRoom(
		"south-of-house",
		"South of House",
		"You are facing the south side of a white house. There is no door here, and all the windows are boarded.",
	)
	southOfHouse.Flags.IsOutdoors = true
	southOfHouse.AddExit("west", "west-of-house") // per ZIL
	southOfHouse.AddExit("nw", "west-of-house") // per ZIL
	southOfHouse.AddExit("east", "behind-house") // per ZIL (EAST-OF-HOUSE)
	southOfHouse.AddExit("ne", "behind-house") // per ZIL
	southOfHouse.AddExit("south", "forest-3") // per ZIL
	// Note: NORTH blocked by house ("The windows are all boarded.")
	g.Rooms["south-of-house"] = southOfHouse

	// BEHIND-HOUSE (EAST-OF-HOUSE in ZIL)
	behindHouse := NewRoom(
		"behind-house",
		"Behind House",
		"You are behind the white house. A path leads into the forest to the east. In one corner of the house there is a small window which is slightly ajar.",
	)
	behindHouse.Flags.IsOutdoors = true
	behindHouse.AddExit("north", "north-of-house") // per ZIL
	behindHouse.AddExit("nw", "north-of-house") // per ZIL
	behindHouse.AddExit("south", "south-of-house") // per ZIL
	behindHouse.AddExit("sw", "south-of-house") // per ZIL
	behindHouse.AddExit("east", "clearing") // per ZIL
	// Kitchen access via window only (per ZIL)
	behindHouse.AddConditionalExit("west", "kitchen", "window-open", "The window is closed.")
	behindHouse.AddConditionalExit("in", "kitchen", "window-open", "The window is closed.")
	g.Rooms["behind-house"] = behindHouse

	// FOREST-1
	forest1 := NewRoom(
		"forest-1",
		"Forest",
		"This is a forest, with trees in all directions. To the east, there appears to be sunlight.",
	)
	forest1.Flags.IsOutdoors = true
	forest1.AddExit("north", "grating-clearing") // per ZIL
	forest1.AddExit("east", "path") // per ZIL
	forest1.AddExit("south", "forest-3") // per ZIL
	g.Rooms["forest-1"] = forest1

	// FOREST-2
	forest2 := NewRoom(
		"forest-2",
		"Forest",
		"This is a dimly lit forest, with large trees all around.",
	)
	forest2.Flags.IsOutdoors = true
	forest2.AddExit("east", "mountains") // per ZIL
	forest2.AddExit("south", "clearing") // per ZIL
	forest2.AddExit("west", "path") // per ZIL
	g.Rooms["forest-2"] = forest2

	// MOUNTAINS
	mountains := NewRoom(
		"mountains",
		"Forest",
		"The forest thins out, revealing impassable mountains.",
	)
	mountains.Flags.IsOutdoors = true
	mountains.AddExit("south", "forest-2") // matches forest-2 east to mountains
	mountains.AddExit("north", "forest-2") // additional path
	mountains.AddExit("west", "forest-2")
	g.Rooms["mountains"] = mountains

	// FOREST-3
	forest3 := NewRoom(
		"forest-3",
		"Forest",
		"This is a dimly lit forest, with large trees all around.",
	)
	forest3.Flags.IsOutdoors = true
	forest3.AddExit("north", "clearing") // per ZIL
	forest3.AddExit("west", "forest-1") // per ZIL
	forest3.AddExit("nw", "south-of-house") // per ZIL
	g.Rooms["forest-3"] = forest3

	// PATH (Forest Path)
	path := NewRoom(
		"path",
		"Forest Path",
		"This is a path winding through a dimly lit forest. The path heads north-south here. One particularly large tree with some low branches stands at the edge of the path.",
	)
	path.Flags.IsOutdoors = true
	path.AddExit("north", "grating-clearing") // per ZIL
	path.AddExit("east", "forest-2") // per ZIL
	path.AddExit("south", "north-of-house") // per ZIL
	path.AddExit("west", "forest-1") // per ZIL
	path.AddExit("up", "up-a-tree")
	g.Rooms["path"] = path

	// UP-A-TREE
	upATree := NewRoom(
		"up-a-tree",
		"Up a Tree",
		"You are about 10 feet above the ground nestled among some large branches. The nearest branch above you is above your reach.",
	)
	upATree.Flags.IsOutdoors = true
	upATree.AddExit("down", "path")
	g.Rooms["up-a-tree"] = upATree

	// GRATING-CLEARING
	gratingClearing := NewRoom(
		"grating-clearing",
		"Clearing",
		"You are in a clearing near a large grating that descends into the ground.",
	)
	gratingClearing.Flags.IsOutdoors = true
	gratingClearing.AddExit("east", "forest-2") // per ZIL
	gratingClearing.AddExit("west", "forest-1") // per ZIL
	gratingClearing.AddExit("south", "path") // per ZIL
	gratingClearing.AddConditionalExit("down", "grating-room", "grate-open", "The grating is closed.")
	g.Rooms["grating-clearing"] = gratingClearing

	// CLEARING
	clearing := NewRoom(
		"clearing",
		"Clearing",
		"You are in a small clearing in a well marked forest path that extends to the east and west.",
	)
	clearing.Flags.IsOutdoors = true
	clearing.AddExit("east", "canyon-view") // per ZIL
	clearing.AddExit("north", "forest-2") // per ZIL
	clearing.AddExit("south", "forest-3") // per ZIL
	clearing.AddExit("west", "behind-house") // per ZIL (ZIL calls it east-of-house)
	g.Rooms["clearing"] = clearing

	// CANYON-VIEW
	canyonView := NewRoom(
		"canyon-view",
		"Canyon View",
		"You are at the top of the Great Canyon on its west wall. From here there is a marvelous view of the canyon and parts of the Frigid River upstream. Across the canyon, the walls of the White Cliffs join the mighty ramparts of the Flathead Mountains to the east. Following the Canyon upstream to the north, Aragain Falls may be seen, complete with rainbow. The mighty Frigid River flows out from a great dark cavern. To the west and south can be seen an immense forest, stretching for miles around. A path leads northwest. It is possible to climb down into the canyon from here.",
	)
	canyonView.Flags.IsOutdoors = true
	canyonView.AddExit("east", "cliff-middle") // per ZIL
	canyonView.AddExit("down", "cliff-middle") // per ZIL
	canyonView.AddExit("nw", "clearing") // per ZIL
	// Note: ZIL has no west exit - intentionally non-bidirectional with clearing
	g.Rooms["canyon-view"] = canyonView
}

// createHouseRooms creates the white house interior
func createHouseRooms(g *GameV2) {
	// KITCHEN
	kitchen := NewRoom(
		"kitchen",
		"Kitchen",
		"You are in the kitchen of the white house. A table seems to have been used recently for the preparation of food. A passage leads to the west and a dark staircase can be seen leading upward. A dark chimney leads down and to the east is a small window which is open.",
	)
	kitchen.AddConditionalExit("east", "behind-house", "window-open", "The window is closed.")
	kitchen.AddConditionalExit("out", "behind-house", "window-open", "The window is closed.")
	kitchen.AddExit("west", "living-room")
	kitchen.AddExit("up", "attic")
	g.Rooms["kitchen"] = kitchen

	// ATTIC
	attic := NewRoom(
		"attic",
		"Attic",
		"This is the attic. The only exit is a stairway leading down.",
	)
	attic.AddExit("down", "kitchen")
	g.Rooms["attic"] = attic

	// LIVING-ROOM
	livingRoom := NewRoom(
		"living-room",
		"Living Room",
		"You are in the living room. There is a doorway to the east, a wooden door with strange gothic lettering to the west (which appears to be nailed shut), a trophy case, and a large oriental rug in the center of the room.",
	)
	livingRoom.AddExit("east", "kitchen")
	livingRoom.AddConditionalExit("west", "strange-passage", "magic-flag", "The door is nailed shut.")
	livingRoom.AddConditionalExit("down", "cellar", "trap-door-open", "You can't go that way.")
	g.Rooms["living-room"] = livingRoom
}

// createCellarAndVicinity creates cellar area rooms
func createCellarAndVicinity(g *GameV2) {
	// CELLAR
	cellar := NewRoom(
		"cellar",
		"Cellar",
		"You are in a dark and damp cellar with a narrow passageway leading north, and a crawlway to the south. On the west is the bottom of a steep metal ramp which is unclimbable.",
	)
	cellar.Flags.IsLit = false
	cellar.Flags.IsDark = true
	cellar.AddExit("north", "troll-room")
	cellar.AddExit("south", "east-of-chasm")
	cellar.AddConditionalExit("up", "living-room", "trap-door-open", "The trap door is closed.")
	g.Rooms["cellar"] = cellar

	// TROLL-ROOM
	trollRoom := NewRoom(
		"troll-room",
		"The Troll Room",
		"This is a small room with passages to the east and south and a forbidding hole leading west. Bloodstains and deep scratches (perhaps made by an axe) mar the walls.",
	)
	trollRoom.Flags.IsLit = false
	trollRoom.Flags.IsDark = true
	trollRoom.AddExit("south", "cellar")
	trollRoom.AddConditionalExit("east", "ew-passage", "troll-dead", "The troll fends you off with a menacing gesture.")
	trollRoom.AddConditionalExit("west", "maze-1", "troll-dead", "The troll fends you off with a menacing gesture.")
	g.Rooms["troll-room"] = trollRoom

	// EAST-OF-CHASM
	eastOfChasm := NewRoom(
		"east-of-chasm",
		"East of Chasm",
		"You are on the east edge of a chasm, the bottom of which cannot be seen. A narrow passage goes north, and the path you are on continues to the east.",
	)
	eastOfChasm.Flags.IsLit = false
	eastOfChasm.Flags.IsDark = true
	eastOfChasm.AddExit("north", "cellar")
	eastOfChasm.AddExit("east", "gallery")
	g.Rooms["east-of-chasm"] = eastOfChasm

	// GALLERY
	gallery := NewRoom(
		"gallery",
		"Gallery",
		"This is an art gallery. Most of the paintings have been stolen by vandals with exceptional taste. The vandals left through either the north or west exits.",
	)
	gallery.AddExit("west", "east-of-chasm")
	gallery.AddExit("north", "studio")
	g.Rooms["gallery"] = gallery

	// STUDIO
	studio := NewRoom(
		"studio",
		"Studio",
		"This appears to have been an artist's studio. The walls and floors are splattered with paints of 69 different colors. Strangely enough, nothing of value is hanging here. At the south end of the room is an open door (also covered with paint). A dark and narrow chimney leads up from a fireplace; although you might be able to get up it, it seems unlikely you could get back down.",
	)
	studio.AddExit("south", "gallery")
	g.Rooms["studio"] = studio
}

// createMazeRooms creates the twisty maze (15 rooms + 4 dead ends)
func createMazeRooms(g *GameV2) {
	// MAZE-1
	maze1 := NewRoom("maze-1", "Maze", "This is part of a maze of twisty little passages, all alike.")
	maze1.Flags.IsLit = false
	maze1.Flags.IsDark = true
	maze1.AddExit("east", "troll-room")
	maze1.AddExit("north", "maze-1") // self-loop per ZIL source
	maze1.AddExit("south", "maze-2") // per ZIL source
	maze1.AddExit("west", "maze-4")
	g.Rooms["maze-1"] = maze1

	// MAZE-2
	maze2 := NewRoom("maze-2", "Maze", "This is part of a maze of twisty little passages, all alike.")
	maze2.Flags.IsLit = false
	maze2.Flags.IsDark = true
	maze2.AddExit("south", "maze-1")
	maze2.AddExit("east", "maze-3")
	g.Rooms["maze-2"] = maze2

	// MAZE-3
	maze3 := NewRoom("maze-3", "Maze", "This is part of a maze of twisty little passages, all alike.")
	maze3.Flags.IsLit = false
	maze3.Flags.IsDark = true
	maze3.AddExit("west", "maze-2")
	maze3.AddExit("north", "maze-4")
	maze3.AddExit("up", "maze-5")
	g.Rooms["maze-3"] = maze3

	// MAZE-4
	maze4 := NewRoom("maze-4", "Maze", "This is part of a maze of twisty little passages, all alike.")
	maze4.Flags.IsLit = false
	maze4.Flags.IsDark = true
	maze4.AddExit("west", "maze-3") // per ZIL source
	maze4.AddExit("north", "maze-1") // per ZIL source
	maze4.AddExit("east", "dead-end-1") // per ZIL source
	g.Rooms["maze-4"] = maze4

	// DEAD-END-1
	deadEnd1 := NewRoom("dead-end-1", "Dead End", "You have come to a dead end in the maze.")
	deadEnd1.Flags.IsLit = false
	deadEnd1.Flags.IsDark = true
	deadEnd1.AddExit("south", "maze-4")
	g.Rooms["dead-end-1"] = deadEnd1

	// MAZE-5
	maze5 := NewRoom("maze-5", "Maze", "This is part of a maze of twisty little passages, all alike. A skeleton, probably the remains of a luckless adventurer, lies here.")
	maze5.Flags.IsLit = false
	maze5.Flags.IsDark = true
	maze5.AddExit("east", "dead-end-2")
	maze5.AddExit("north", "maze-3")
	maze5.AddExit("sw", "maze-6")
	g.Rooms["maze-5"] = maze5

	// DEAD-END-2
	deadEnd2 := NewRoom("dead-end-2", "Dead End", "You have come to a dead end in the maze.")
	deadEnd2.Flags.IsLit = false
	deadEnd2.Flags.IsDark = true
	deadEnd2.AddExit("west", "maze-5")
	g.Rooms["dead-end-2"] = deadEnd2

	// MAZE-6
	maze6 := NewRoom("maze-6", "Maze", "This is part of a maze of twisty little passages, all alike.")
	maze6.Flags.IsLit = false
	maze6.Flags.IsDark = true
	maze6.AddExit("down", "maze-5")
	maze6.AddExit("east", "maze-7") // FIX: matches maze-7 west
	maze6.AddExit("west", "maze-6")
	maze6.AddExit("up", "maze-9")
	g.Rooms["maze-6"] = maze6

	// MAZE-7
	maze7 := NewRoom("maze-7", "Maze", "This is part of a maze of twisty little passages, all alike.")
	maze7.Flags.IsLit = false
	maze7.Flags.IsDark = true
	maze7.AddExit("up", "maze-14")
	maze7.AddExit("west", "maze-6")
	maze7.AddExit("east", "maze-8")
	maze7.AddExit("south", "maze-15")
	g.Rooms["maze-7"] = maze7

	// MAZE-8
	maze8 := NewRoom("maze-8", "Maze", "This is part of a maze of twisty little passages, all alike.")
	maze8.Flags.IsLit = false
	maze8.Flags.IsDark = true
	maze8.AddExit("ne", "maze-7")
	maze8.AddExit("west", "maze-8") // self-loop per ZIL source
	maze8.AddExit("se", "dead-end-3")
	g.Rooms["maze-8"] = maze8

	// DEAD-END-3
	deadEnd3 := NewRoom("dead-end-3", "Dead End", "You have come to a dead end in the maze.")
	deadEnd3.Flags.IsLit = false
	deadEnd3.Flags.IsDark = true
	deadEnd3.AddExit("north", "maze-8")
	g.Rooms["dead-end-3"] = deadEnd3

	// MAZE-9
	maze9 := NewRoom("maze-9", "Maze", "This is part of a maze of twisty little passages, all alike.")
	maze9.Flags.IsLit = false
	maze9.Flags.IsDark = true
	maze9.AddExit("north", "maze-6")
	maze9.AddExit("east", "maze-10")
	maze9.AddExit("south", "maze-13")
	maze9.AddExit("west", "maze-12") // FIX: matches maze-12 east
	maze9.AddExit("nw", "maze-9")
	g.Rooms["maze-9"] = maze9

	// MAZE-10
	maze10 := NewRoom("maze-10", "Maze", "This is part of a maze of twisty little passages, all alike.")
	maze10.Flags.IsLit = false
	maze10.Flags.IsDark = true
	maze10.AddExit("east", "maze-9")
	maze10.AddExit("west", "maze-13") // FIX: matches maze-13 east
	maze10.AddExit("up", "maze-11")
	g.Rooms["maze-10"] = maze10

	// MAZE-11
	maze11 := NewRoom("maze-11", "Maze", "This is part of a maze of twisty little passages, all alike.")
	maze11.Flags.IsLit = false
	maze11.Flags.IsDark = true
	maze11.AddExit("ne", "grating-room") // FIX: matches grating-room sw
	maze11.AddExit("down", "maze-10")
	maze11.AddExit("nw", "maze-13")
	maze11.AddExit("sw", "maze-12")
	g.Rooms["maze-11"] = maze11

	// GRATING-ROOM
	gratingRoom := NewRoom(
		"grating-room",
		"Grating Room",
		"You are in a small room near a grating in the ceiling which admits a dim light. There are passages to the south and the southwest.",
	)
	gratingRoom.Flags.IsLit = false
	gratingRoom.Flags.IsDark = true
	gratingRoom.AddExit("sw", "maze-11")
	gratingRoom.AddConditionalExit("up", "grating-clearing", "grate-open", "The grating is closed.")
	g.Rooms["grating-room"] = gratingRoom

	// MAZE-12
	maze12 := NewRoom("maze-12", "Maze", "This is part of a maze of twisty little passages, all alike.")
	maze12.Flags.IsLit = false
	maze12.Flags.IsDark = true
	maze12.AddExit("sw", "maze-11")
	maze12.AddExit("east", "maze-13") // FIX: matches maze-13 west
	maze12.AddExit("up", "maze-9") // per ZIL source
	maze12.AddExit("north", "dead-end-4")
	g.Rooms["maze-12"] = maze12

	// DEAD-END-4
	deadEnd4 := NewRoom("dead-end-4", "Dead End", "You have come to a dead end in the maze.")
	deadEnd4.Flags.IsLit = false
	deadEnd4.Flags.IsDark = true
	deadEnd4.AddExit("south", "maze-12")
	g.Rooms["dead-end-4"] = deadEnd4

	// MAZE-13
	maze13 := NewRoom("maze-13", "Maze", "This is part of a maze of twisty little passages, all alike.")
	maze13.Flags.IsLit = false
	maze13.Flags.IsDark = true
	maze13.AddExit("east", "maze-9") // FIX: matches maze-9 south
	maze13.AddExit("down", "maze-12")
	maze13.AddExit("south", "maze-10")
	maze13.AddExit("west", "maze-11") // per ZIL source
	g.Rooms["maze-13"] = maze13

	// MAZE-14
	maze14 := NewRoom("maze-14", "Maze", "This is part of a maze of twisty little passages, all alike.")
	maze14.Flags.IsLit = false
	maze14.Flags.IsDark = true
	maze14.AddExit("west", "maze-15")
	maze14.AddExit("nw", "maze-14")
	maze14.AddExit("ne", "maze-7")
	maze14.AddExit("south", "maze-7")
	g.Rooms["maze-14"] = maze14

	// MAZE-15
	maze15 := NewRoom("maze-15", "Maze", "This is part of a maze of twisty little passages, all alike.")
	maze15.Flags.IsLit = false
	maze15.Flags.IsDark = true
	maze15.AddExit("west", "maze-14")
	maze15.AddExit("south", "maze-7")
	maze15.AddExit("se", "cyclops-room")
	g.Rooms["maze-15"] = maze15
}

// createCyclopsArea creates cyclops and treasure room
func createCyclopsArea(g *GameV2) {
	// CYCLOPS-ROOM
	cyclopsRoom := NewRoom(
		"cyclops-room",
		"Cyclops Room",
		"This is a large room hewn out of solid rock. A cyclops, who looks prepared to swing a large club, blocks the stairway leading upward.",
	)
	cyclopsRoom.Flags.IsLit = false
	cyclopsRoom.Flags.IsDark = true
	cyclopsRoom.AddExit("nw", "maze-15")
	cyclopsRoom.AddConditionalExit("east", "strange-passage", "magic-flag", "The east wall is solid rock.")
	cyclopsRoom.AddConditionalExit("up", "treasure-room", "cyclops-dead", "The cyclops doesn't look like he'll let you past.")
	g.Rooms["cyclops-room"] = cyclopsRoom

	// STRANGE-PASSAGE
	strangePassage := NewRoom(
		"strange-passage",
		"Strange Passage",
		"This is a long passage. To the west is one entrance. On the east there is an old wooden door, with a large opening in it (about cyclops sized).",
	)
	strangePassage.AddExit("west", "cyclops-room")
	strangePassage.AddExit("in", "cyclops-room")
	strangePassage.AddExit("east", "living-room")
	g.Rooms["strange-passage"] = strangePassage

	// TREASURE-ROOM
	treasureRoom := NewRoom(
		"treasure-room",
		"Treasure Room",
		"This is a large room, whose east wall is solid granite. A number of discarded bags, which crumble at your touch, are scattered about on the floor. There is an exit down a staircase.",
	)
	treasureRoom.AddExit("down", "cyclops-room")
	g.Rooms["treasure-room"] = treasureRoom
}

// createReservoirArea creates reservoir and dam area
func createReservoirArea(g *GameV2) {
	// RESERVOIR-SOUTH
	reservoirSouth := NewRoom(
		"reservoir-south",
		"Reservoir South",
		"You are in a long room on the south shore of a large lake, far too deep and wide for crossing.",
	)
	reservoirSouth.AddExit("se", "deep-canyon")
	reservoirSouth.AddExit("sw", "chasm-room")
	reservoirSouth.AddExit("east", "dam-room")
	reservoirSouth.AddExit("west", "stream-view")
	reservoirSouth.AddConditionalExit("north", "reservoir", "low-tide", "You would drown.")
	g.Rooms["reservoir-south"] = reservoirSouth

	// RESERVOIR
	reservoir := NewRoom(
		"reservoir",
		"Reservoir",
		"You are on the reservoir. The water is cold and provides little buoyancy.",
	)
	reservoir.AddExit("north", "reservoir-north")
	reservoir.AddConditionalExit("south", "reservoir-south", "low-tide", "You would drown.")
	reservoir.AddExit("up", "in-stream")
	reservoir.AddExit("west", "in-stream")
	g.Rooms["reservoir"] = reservoir

	// RESERVOIR-NORTH
	reservoirNorth := NewRoom(
		"reservoir-north",
		"Reservoir North",
		"You are in a long room on the north shore of a large lake, far too deep and wide for crossing.",
	)
	reservoirNorth.AddExit("north", "atlantis-room")
	reservoirNorth.AddConditionalExit("south", "reservoir", "low-tide", "You would drown.")
	g.Rooms["reservoir-north"] = reservoirNorth

	// STREAM-VIEW
	streamView := NewRoom(
		"stream-view",
		"Stream View",
		"You are standing on a path beside a gently flowing stream. The path follows the stream, which flows from west to east.",
	)
	streamView.AddExit("east", "reservoir-south")
	g.Rooms["stream-view"] = streamView

	// IN-STREAM
	inStream := NewRoom(
		"in-stream",
		"Stream",
		"You are on the gently flowing stream. The upstream route is too narrow to navigate, and the downstream route is invisible due to twisting walls. There is a narrow beach to land on.",
	)
	inStream.AddExit("land", "stream-view")
	inStream.AddExit("down", "reservoir")
	inStream.AddExit("east", "reservoir")
	g.Rooms["in-stream"] = inStream
}

// createMirrorRooms creates mirror room area
func createMirrorRooms(g *GameV2) {
	// MIRROR-ROOM-1
	mirrorRoom1 := NewRoom(
		"mirror-room-1",
		"Mirror Room",
		"You are in a large square room with tall ceilings. On the south wall is an enormous mirror which fills the entire wall. There are exits on the other three sides of the room.",
	)
	mirrorRoom1.AddExit("north", "cold-passage") // per ZIL
	mirrorRoom1.AddExit("west", "twisting-passage") // per ZIL
	mirrorRoom1.AddExit("east", "small-cave") // per ZIL
	g.Rooms["mirror-room-1"] = mirrorRoom1

	// MIRROR-ROOM-2
	mirrorRoom2 := NewRoom(
		"mirror-room-2",
		"Mirror Room",
		"You are in a large square room with tall ceilings. On the north wall is an enormous mirror which fills the entire wall. There are exits on the other three sides of the room.",
	)
	mirrorRoom2.AddExit("west", "winding-passage") // per ZIL
	mirrorRoom2.AddExit("north", "narrow-passage") // per ZIL
	mirrorRoom2.AddExit("east", "tiny-cave") // per ZIL
	g.Rooms["mirror-room-2"] = mirrorRoom2

	// SMALL-CAVE
	smallCave := NewRoom(
		"small-cave",
		"Cave",
		"This is a tiny cave with entrances west and north, and a staircase leading down.",
	)
	smallCave.AddExit("north", "mirror-room-1") // per ZIL
	smallCave.AddExit("west", "twisting-passage") // per ZIL
	smallCave.AddExit("down", "atlantis-room") // per ZIL
	smallCave.AddExit("south", "atlantis-room") // per ZIL
	g.Rooms["small-cave"] = smallCave

	// TINY-CAVE
	tinyCave := NewRoom(
		"tiny-cave",
		"Cave",
		"This is a tiny cave with entrances west and north, and a dark, forbidding staircase leading down.",
	)
	tinyCave.AddExit("north", "mirror-room-2") // per ZIL
	tinyCave.AddExit("west", "winding-passage") // per ZIL
	tinyCave.AddExit("down", "entrance-to-hades") // per ZIL
	g.Rooms["tiny-cave"] = tinyCave

	// COLD-PASSAGE
	coldPassage := NewRoom(
		"cold-passage",
		"Cold Passage",
		"This is a cold and damp corridor where a long east-west passageway turns into a southward path.",
	)
	coldPassage.AddExit("south", "mirror-room-1")
	coldPassage.AddExit("west", "slide-room")
	g.Rooms["cold-passage"] = coldPassage

	// NARROW-PASSAGE
	narrowPassage := NewRoom(
		"narrow-passage",
		"Narrow Passage",
		"This is a long and narrow corridor where a long north-south passageway briefly narrows even further.",
	)
	narrowPassage.AddExit("north", "round-room")
	narrowPassage.AddExit("south", "mirror-room-2")
	g.Rooms["narrow-passage"] = narrowPassage

	// WINDING-PASSAGE
	windingPassage := NewRoom(
		"winding-passage",
		"Winding Passage",
		"This is a winding passage. It seems that there are only exits on the east and north.",
	)
	windingPassage.AddExit("north", "mirror-room-2") // per ZIL
	windingPassage.AddExit("east", "tiny-cave") // per ZIL
	g.Rooms["winding-passage"] = windingPassage

	// TWISTING-PASSAGE
	twistingPassage := NewRoom(
		"twisting-passage",
		"Twisting Passage",
		"This is a winding passage. It seems that there are only exits on the east and north.",
	)
	twistingPassage.AddExit("north", "mirror-room-1") // per ZIL
	twistingPassage.AddExit("east", "small-cave") // per ZIL
	g.Rooms["twisting-passage"] = twistingPassage

	// ATLANTIS-ROOM
	atlantisRoom := NewRoom(
		"atlantis-room",
		"Atlantis Room",
		"This is an ancient room, long under water. There is an exit to the south and a staircase leading up.",
	)
	atlantisRoom.AddExit("up", "small-cave")
	atlantisRoom.AddExit("south", "reservoir-north")
	g.Rooms["atlantis-room"] = atlantisRoom
}

// createRoundRoomArea creates round room and vicinity
func createRoundRoomArea(g *GameV2) {
	// EW-PASSAGE (East-West Passage)
	ewPassage := NewRoom(
		"ew-passage",
		"East-West Passage",
		"This is a narrow east-west passageway. There is a narrow stairway leading down at the north end of the room.",
	)
	ewPassage.Flags.IsLit = false
	ewPassage.Flags.IsDark = true
	ewPassage.AddExit("east", "round-room")
	ewPassage.AddExit("west", "troll-room")
	ewPassage.AddExit("down", "chasm-room")
	ewPassage.AddExit("north", "chasm-room")
	ewPassage.AddExit("ne", "chasm-room") // ADD: reverse of chasm-room sw
	g.Rooms["ew-passage"] = ewPassage

	// ROUND-ROOM
	roundRoom := NewRoom(
		"round-room",
		"Round Room",
		"This is a circular stone room with passages in all directions. Several of them have unfortunately been blocked by cave-ins.",
	)
	roundRoom.Flags.IsLit = false
	roundRoom.Flags.IsDark = true
	roundRoom.AddExit("east", "loud-room")
	roundRoom.AddExit("west", "ew-passage")
	roundRoom.AddExit("north", "ns-passage")
	roundRoom.AddExit("south", "narrow-passage")
	roundRoom.AddExit("se", "engravings-cave")
	g.Rooms["round-room"] = roundRoom

	// DEEP-CANYON
	deepCanyon := NewRoom(
		"deep-canyon",
		"Deep Canyon",
		"You are on the south edge of a deep canyon. Passages lead off to the east, northwest, and southwest.",
	)
	deepCanyon.AddExit("nw", "reservoir-south")
	deepCanyon.AddExit("east", "dam-room") // matches dam-room west
	deepCanyon.AddExit("sw", "ns-passage")
	deepCanyon.AddExit("down", "loud-room")
	g.Rooms["deep-canyon"] = deepCanyon

	// DAMP-CAVE
	dampCave := NewRoom(
		"damp-cave",
		"Damp Cave",
		"This cave has exits to the west and east, and narrows to a crack toward the south. The earth is particularly damp here.",
	)
	dampCave.Flags.IsLit = false
	dampCave.Flags.IsDark = true
	dampCave.AddExit("west", "loud-room")
	dampCave.AddExit("east", "white-cliffs-north")
	g.Rooms["damp-cave"] = dampCave

	// LOUD-ROOM
	loudRoom := NewRoom(
		"loud-room",
		"Loud Room",
		"This is a large room with a ceiling which cannot be detected from the ground. There is a narrow passage from east to west and a stone stairway leading upward. The room is extremely noisy. In fact, it is difficult to hear yourself think.",
	)
	loudRoom.Flags.IsLit = false
	loudRoom.Flags.IsDark = true
	loudRoom.AddExit("east", "damp-cave")
	loudRoom.AddExit("west", "round-room")
	loudRoom.AddExit("up", "deep-canyon")
	g.Rooms["loud-room"] = loudRoom

	// NS-PASSAGE (North-South Passage)
	nsPassage := NewRoom(
		"ns-passage",
		"North-South Passage",
		"This is a high north-south passage, which forks to the northeast.",
	)
	nsPassage.Flags.IsLit = false
	nsPassage.Flags.IsDark = true
	nsPassage.AddExit("north", "chasm-room")
	nsPassage.AddExit("ne", "deep-canyon")
	nsPassage.AddExit("south", "round-room")
	g.Rooms["ns-passage"] = nsPassage

	// CHASM-ROOM
	chasmRoom := NewRoom(
		"chasm-room",
		"Chasm",
		"A chasm runs southwest to northeast and the path follows it. You are on the south side of the chasm, where a crack opens into a passage.",
	)
	chasmRoom.Flags.IsLit = false
	chasmRoom.Flags.IsDark = true
	chasmRoom.AddExit("ne", "reservoir-south")
	chasmRoom.AddExit("sw", "ew-passage")
	chasmRoom.AddExit("up", "ew-passage")
	chasmRoom.AddExit("south", "ns-passage") // per ZIL
	g.Rooms["chasm-room"] = chasmRoom
}

// createHadesArea creates entrance to hades and land of the dead
func createHadesArea(g *GameV2) {
	// ENTRANCE-TO-HADES
	entranceToHades := NewRoom(
		"entrance-to-hades",
		"Entrance to Hades",
		"You are outside a large gateway, on which is inscribed \"Abandon every hope, all ye who enter here.\" The gate is open; through it you can see a desolation, with a pile of mangled bodies in one corner. Thousands of voices, lamenting some hideous fate, can be heard.",
	)
	entranceToHades.AddExit("up", "tiny-cave")
	entranceToHades.AddConditionalExit("in", "land-of-living-dead", "LLD-FLAG", "Some invisible force prevents you from passing through the gate.")
	entranceToHades.AddConditionalExit("south", "land-of-living-dead", "LLD-FLAG", "Some invisible force prevents you from passing through the gate.")
	g.Rooms["entrance-to-hades"] = entranceToHades

	// LAND-OF-LIVING-DEAD
	landOfLivingDead := NewRoom(
		"land-of-living-dead",
		"Land of the Dead",
		"You have entered the Land of the Living Dead. Thousands of lost souls can be heard weeping and moaning. In the corner are stacked the remains of dozens of previous adventurers less fortunate than yourself. A passage exits to the north.",
	)
	landOfLivingDead.AddExit("out", "entrance-to-hades")
	landOfLivingDead.AddExit("north", "entrance-to-hades")
	g.Rooms["land-of-living-dead"] = landOfLivingDead
}

// createTempleArea creates engravings, dome, egypt, temple
func createTempleArea(g *GameV2) {
	// ENGRAVINGS-CAVE
	engravingsCave := NewRoom(
		"engravings-cave",
		"Engravings Cave",
		"You have entered a low cave with passages leading northwest and east.",
	)
	engravingsCave.Flags.IsLit = false
	engravingsCave.Flags.IsDark = true
	engravingsCave.AddExit("nw", "round-room")
	engravingsCave.AddExit("east", "dome-room")
	g.Rooms["engravings-cave"] = engravingsCave

	// EGYPT-ROOM
	egyptRoom := NewRoom(
		"egypt-room",
		"Egyptian Room",
		"This is a room which looks like an Egyptian tomb. There is an ascending staircase to the west.",
	)
	egyptRoom.AddExit("west", "north-temple")
	egyptRoom.AddExit("up", "north-temple")
	g.Rooms["egypt-room"] = egyptRoom

	// DOME-ROOM
	domeRoom := NewRoom(
		"dome-room",
		"Dome Room",
		"You are at the periphery of a large dome, which forms the ceiling of another room below. Protecting you from a precipitous drop is a wooden railing which circles the dome.",
	)
	domeRoom.AddExit("west", "engravings-cave")
	domeRoom.AddConditionalExit("down", "torch-room", "dome-flag", "You cannot go down without fracturing many bones.")
	g.Rooms["dome-room"] = domeRoom

	// TORCH-ROOM
	torchRoom := NewRoom(
		"torch-room",
		"Torch Room",
		"This is a large room with a prominent doorway leading to a down staircase. Above you is a large dome. Up around the edge of the dome (20 feet up) is a wooden railing. In the center of the room there is a white marble pedestal.",
	)
	torchRoom.AddExit("south", "north-temple")
	torchRoom.AddExit("down", "north-temple")
	torchRoom.AddExit("in", "north-temple") // ADD: reverse of north-temple out
	g.Rooms["torch-room"] = torchRoom

	// NORTH-TEMPLE
	northTemple := NewRoom(
		"north-temple",
		"Temple",
		"This is the north end of a large temple. On the east wall is an ancient inscription, probably a prayer in a long-forgotten language. Below the prayer is a staircase leading down. The west wall is solid granite. The exit to the north end of the room is through huge marble pillars.",
	)
	northTemple.AddExit("down", "egypt-room")
	northTemple.AddExit("east", "egypt-room")
	northTemple.AddExit("north", "torch-room")
	northTemple.AddExit("out", "torch-room")
	northTemple.AddExit("up", "torch-room")
	northTemple.AddExit("south", "south-temple")
	g.Rooms["north-temple"] = northTemple

	// SOUTH-TEMPLE
	southTemple := NewRoom(
		"south-temple",
		"Altar",
		"This is the south end of a large temple. In front of you is what appears to be an altar. In one corner is a small hole in the floor which leads into darkness. You probably could not get back up it.",
	)
	southTemple.AddExit("north", "north-temple")
	southTemple.AddConditionalExit("down", "tiny-cave", "coffin-cure", "You haven't a prayer of getting the coffin down there.")
	g.Rooms["south-temple"] = southTemple
}

// createDamArea creates dam rooms
func createDamArea(g *GameV2) {
	// DAM-ROOM
	damRoom := NewRoom(
		"dam-room",
		"Dam",
		"You are standing on the top of the Flood Control Dam #3, which was quite a tourist attraction in times far distant. There are paths to the north, south, and west, and a scramble down.",
	)
	damRoom.AddExit("south", "deep-canyon") // per ZIL
	damRoom.AddExit("down", "dam-base") // per ZIL
	damRoom.AddExit("east", "dam-base") // per ZIL
	damRoom.AddExit("north", "dam-lobby") // per ZIL
	damRoom.AddExit("west", "reservoir-south") // per ZIL
	g.Rooms["dam-room"] = damRoom

	// DAM-LOBBY
	damLobby := NewRoom(
		"dam-lobby",
		"Dam Lobby",
		"This room appears to have been the waiting room for groups touring the dam. There are open doorways here to the north and east marked \"Private\", and there is a path leading south over the top of the dam.",
	)
	damLobby.AddExit("south", "dam-room")
	damLobby.AddExit("north", "maintenance-room")
	damLobby.AddExit("east", "maintenance-room")
	g.Rooms["dam-lobby"] = damLobby

	// MAINTENANCE-ROOM
	maintenanceRoom := NewRoom(
		"maintenance-room",
		"Maintenance Room",
		"This is what appears to have been the maintenance room for Flood Control Dam #3. Apparently, this room has been ransacked recently, for most of the valuable equipment is gone. On the wall in front of you is a group of buttons colored blue, yellow, brown, and red. There are doorways to the west and south.",
	)
	maintenanceRoom.AddExit("south", "dam-lobby")
	maintenanceRoom.AddExit("west", "dam-lobby")
	g.Rooms["maintenance-room"] = maintenanceRoom

	// DAM-BASE
	damBase := NewRoom(
		"dam-base",
		"Dam Base",
		"You are at the base of Flood Control Dam #3, which looms above you and to the north. The river Frigid is flowing by here. Along the river are the White Cliffs which seem to form giant walls stretching from north to south along the shores of the river as it winds its way downstream.",
	)
	damBase.Flags.IsOutdoors = true
	damBase.AddExit("north", "dam-room") // per ZIL
	damBase.AddExit("up", "dam-room")    // per ZIL
	// Note: ZIL has no east or west exits - intentional one-ways:
	// - dam-room east→dam-base (one-way, return via up/north)
	// - river-1 west→dam-base (one-way, no return)
	g.Rooms["dam-base"] = damBase
}

// createRiverArea creates frigid river rooms
func createRiverArea(g *GameV2) {
	// RIVER-1
	river1 := NewRoom(
		"river-1",
		"Frigid River",
		"You are on the Frigid River in the vicinity of the Dam. The river flows quietly here. There is a landing on the west shore.",
	)
	river1.Flags.IsOutdoors = true
	river1.AddExit("west", "dam-base")
	river1.AddExit("land", "dam-base")
	river1.AddExit("down", "river-2")
	g.Rooms["river-1"] = river1

	// RIVER-2
	river2 := NewRoom(
		"river-2",
		"Frigid River",
		"The river turns a corner here making it impossible to see the Dam. The White Cliffs loom on the east bank and large rocks prevent landing on the west.",
	)
	river2.Flags.IsOutdoors = true
	river2.AddExit("down", "river-3")
	g.Rooms["river-2"] = river2

	// RIVER-3
	river3 := NewRoom(
		"river-3",
		"Frigid River",
		"The river descends here into a valley. There is a narrow beach on the west shore below the cliffs. In the distance a faint rumbling can be heard.",
	)
	river3.Flags.IsOutdoors = true
	river3.AddExit("down", "river-4")
	river3.AddExit("land", "white-cliffs-north")
	river3.AddExit("west", "white-cliffs-north")
	g.Rooms["river-3"] = river3

	// WHITE-CLIFFS-NORTH
	whiteCliffsNorth := NewRoom(
		"white-cliffs-north",
		"White Cliffs Beach",
		"You are on a narrow strip of beach which runs along the base of the White Cliffs. There is a narrow path heading south along the Cliffs and a tight passage leading west into the cliffs themselves.",
	)
	whiteCliffsNorth.Flags.IsOutdoors = true
	whiteCliffsNorth.AddConditionalExit("south", "white-cliffs-south", "deflate", "The path is too narrow.")
	whiteCliffsNorth.AddConditionalExit("west", "damp-cave", "deflate", "The path is too narrow.")
	g.Rooms["white-cliffs-north"] = whiteCliffsNorth

	// WHITE-CLIFFS-SOUTH
	whiteCliffsSouth := NewRoom(
		"white-cliffs-south",
		"White Cliffs Beach",
		"You are on a rocky, narrow strip of beach beside the Cliffs. A narrow path leads north along the shore.",
	)
	whiteCliffsSouth.Flags.IsOutdoors = true
	whiteCliffsSouth.AddConditionalExit("north", "white-cliffs-north", "deflate", "The path is too narrow.")
	g.Rooms["white-cliffs-south"] = whiteCliffsSouth

	// RIVER-4
	river4 := NewRoom(
		"river-4",
		"Frigid River",
		"The river is running faster here and the sound ahead appears to be that of rushing water. On the east shore is a sandy beach. A small area of beach can also be seen below the cliffs on the west shore.",
	)
	river4.Flags.IsOutdoors = true
	river4.AddExit("down", "river-5")
	river4.AddExit("west", "white-cliffs-south")
	river4.AddExit("east", "sandy-beach")
	g.Rooms["river-4"] = river4

	// RIVER-5
	river5 := NewRoom(
		"river-5",
		"Frigid River",
		"The sound of rushing water is nearly unbearable here. On the east shore is a large landing area.",
	)
	river5.Flags.IsOutdoors = true
	river5.AddExit("east", "shore")
	river5.AddExit("land", "shore")
	g.Rooms["river-5"] = river5

	// SHORE
	shore := NewRoom(
		"shore",
		"Shore",
		"You are on the east shore of the river. The water here seems somewhat treacherous. A path travels from north to south here, the south end quickly turning around a sharp corner.",
	)
	shore.Flags.IsOutdoors = true
	shore.AddExit("north", "sandy-beach")
	shore.AddExit("south", "aragain-falls")
	shore.AddExit("west", "river-5") // ADD: reverse of river-5 east/land
	g.Rooms["shore"] = shore

	// SANDY-BEACH
	sandyBeach := NewRoom(
		"sandy-beach",
		"Sandy Beach",
		"You are on a large sandy beach on the east shore of the river, which is flowing quickly by. A path runs beside the river to the south here, and a passage is partially buried in sand to the northeast.",
	)
	sandyBeach.Flags.IsOutdoors = true
	sandyBeach.AddExit("ne", "sandy-cave")
	sandyBeach.AddExit("south", "shore")
	sandyBeach.AddExit("west", "river-4") // ADD: reverse of river-4 east
	g.Rooms["sandy-beach"] = sandyBeach

	// SANDY-CAVE
	sandyCave := NewRoom(
		"sandy-cave",
		"Sandy Cave",
		"This is a sand-filled cave whose exit is to the southwest.",
	)
	sandyCave.AddExit("sw", "sandy-beach")
	g.Rooms["sandy-cave"] = sandyCave

	// ARAGAIN-FALLS
	aragainFalls := NewRoom(
		"aragain-falls",
		"Aragain Falls",
		"You are at the top of Aragain Falls, an enormous waterfall with a drop of about 450 feet. The only path here is on the north end.",
	)
	aragainFalls.Flags.IsOutdoors = true
	aragainFalls.AddExit("north", "shore")
	aragainFalls.AddConditionalExit("west", "on-rainbow", "rainbow-flag", "")
	aragainFalls.AddConditionalExit("up", "on-rainbow", "rainbow-flag", "")
	g.Rooms["aragain-falls"] = aragainFalls

	// ON-RAINBOW
	onRainbow := NewRoom(
		"on-rainbow",
		"On the Rainbow",
		"You are on top of a rainbow (I bet you never thought you would walk on a rainbow), with a magnificent view of the Falls. The rainbow travels east-west here.",
	)
	onRainbow.Flags.IsOutdoors = true
	onRainbow.AddExit("west", "end-of-rainbow")
	onRainbow.AddExit("east", "aragain-falls")
	g.Rooms["on-rainbow"] = onRainbow

	// END-OF-RAINBOW
	endOfRainbow := NewRoom(
		"end-of-rainbow",
		"End of Rainbow",
		"You are on a small, rocky beach on the continuation of the Frigid River past the Falls. The beach is narrow due to the presence of the White Cliffs. The river canyon opens here and sunlight shines in from above. A rainbow crosses over the falls to the east and a narrow path continues to the southwest.",
	)
	endOfRainbow.Flags.IsOutdoors = true
	endOfRainbow.AddConditionalExit("up", "on-rainbow", "rainbow-flag", "")
	endOfRainbow.AddConditionalExit("ne", "on-rainbow", "rainbow-flag", "")
	endOfRainbow.AddConditionalExit("east", "on-rainbow", "rainbow-flag", "")
	endOfRainbow.AddExit("sw", "canyon-bottom")
	g.Rooms["end-of-rainbow"] = endOfRainbow

	// CANYON-BOTTOM
	canyonBottom := NewRoom(
		"canyon-bottom",
		"Canyon Bottom",
		"You are beneath the walls of the river canyon which may be climbable here. The lesser part of the runoff of Aragain Falls flows by below. To the north is a narrow path.",
	)
	canyonBottom.Flags.IsOutdoors = true
	canyonBottom.AddExit("up", "cliff-middle")
	canyonBottom.AddExit("north", "end-of-rainbow")
	g.Rooms["canyon-bottom"] = canyonBottom

	// CLIFF-MIDDLE
	cliffMiddle := NewRoom(
		"cliff-middle",
		"Rocky Ledge",
		"You are on a ledge about halfway up the wall of the river canyon. You can see from here that the main flow from Aragain Falls twists along a passage which it is impossible for you to enter. Below you is the canyon bottom. Above you is more cliff, which appears climbable.",
	)
	cliffMiddle.Flags.IsOutdoors = true
	cliffMiddle.AddExit("up", "canyon-view")
	cliffMiddle.AddExit("down", "canyon-bottom")
	g.Rooms["cliff-middle"] = cliffMiddle
}

// createCoalMineArea creates coal mine rooms
func createCoalMineArea(g *GameV2) {
	// MINE-ENTRANCE
	mineEntrance := NewRoom(
		"mine-entrance",
		"Mine Entrance",
		"You are standing at the entrance of what might have been a coal mine. The shaft enters the west wall, and there is another exit on the south end of the room.",
	)
	mineEntrance.AddExit("south", "slide-room")
	mineEntrance.AddExit("in", "squeeky-room")
	mineEntrance.AddExit("west", "squeeky-room")
	g.Rooms["mine-entrance"] = mineEntrance

	// SQUEEKY-ROOM
	squeekyRoom := NewRoom(
		"squeeky-room",
		"Squeaky Room",
		"You are in a small room. Strange squeaky sounds may be heard coming from the passage at the north end. You may also escape to the east.",
	)
	squeekyRoom.AddExit("north", "bat-room")
	squeekyRoom.AddExit("east", "mine-entrance")
	squeekyRoom.AddExit("out", "mine-entrance") // ADD: reverse of mine-entrance in
	g.Rooms["squeeky-room"] = squeekyRoom

	// BAT-ROOM
	batRoom := NewRoom(
		"bat-room",
		"Bat Room",
		"You are in a small room which has doors only to the east and south.",
	)
	batRoom.Flags.IsOutdoors = false
	batRoom.AddExit("south", "squeeky-room")
	batRoom.AddExit("east", "shaft-room")
	g.Rooms["bat-room"] = batRoom

	// SHAFT-ROOM
	shaftRoom := NewRoom(
		"shaft-room",
		"Shaft Room",
		"This is a large room, in the middle of which is a small shaft descending through the floor into darkness below. To the west and the north are exits from this room. Constructed over the top of the shaft is a metal framework to which a heavy iron chain is attached.",
	)
	shaftRoom.AddExit("west", "bat-room")
	shaftRoom.AddExit("north", "smelly-room")
	g.Rooms["shaft-room"] = shaftRoom

	// SMELLY-ROOM
	smellyRoom := NewRoom(
		"smelly-room",
		"Smelly Room",
		"This is a small nondescript room. However, from the direction of a small descending staircase a foul odor can be detected. To the south is a narrow tunnel.",
	)
	smellyRoom.AddExit("down", "gas-room")
	smellyRoom.AddExit("south", "shaft-room")
	g.Rooms["smelly-room"] = smellyRoom

	// GAS-ROOM
	gasRoom := NewRoom(
		"gas-room",
		"Gas Room",
		"This is a small room which smells strongly of coal gas. There is a short climb up some stairs and a narrow tunnel leading east.",
	)
	gasRoom.Flags.IsOutdoors = false
	gasRoom.AddExit("up", "smelly-room")
	gasRoom.AddExit("east", "mine-1")
	gasRoom.AddExit("south", "mine-1") // ADD: reverse of mine-1 north
	g.Rooms["gas-room"] = gasRoom

	// LADDER-TOP
	ladderTop := NewRoom(
		"ladder-top",
		"Ladder Top",
		"This is a very small room. In the corner is a rickety wooden ladder, leading downward. It might be safe to descend. There is also a staircase leading upward.",
	)
	ladderTop.AddExit("down", "ladder-bottom")
	ladderTop.AddExit("up", "mine-4")
	g.Rooms["ladder-top"] = ladderTop

	// LADDER-BOTTOM
	ladderBottom := NewRoom(
		"ladder-bottom",
		"Ladder Bottom",
		"This is a rather wide room. On one side is the bottom of a narrow wooden ladder. To the west and the south are passages leaving the room.",
	)
	ladderBottom.AddExit("south", "dead-end-5")
	ladderBottom.AddExit("west", "timber-room")
	ladderBottom.AddExit("up", "ladder-top")
	g.Rooms["ladder-bottom"] = ladderBottom

	// DEAD-END-5
	deadEnd5 := NewRoom(
		"dead-end-5",
		"Dead End",
		"You have come to a dead end in the mine.",
	)
	deadEnd5.AddExit("north", "ladder-bottom")
	g.Rooms["dead-end-5"] = deadEnd5

	// TIMBER-ROOM
	timberRoom := NewRoom(
		"timber-room",
		"Timber Room",
		"This is a long and narrow passage, which is cluttered with broken timbers. A wide passage comes from the east and turns at the west end of the room into a very narrow passageway. From the west comes a strong draft.",
	)
	timberRoom.Flags.IsOutdoors = false
	timberRoom.AddExit("east", "ladder-bottom")
	timberRoom.AddConditionalExit("west", "lower-shaft", "empty-handed", "You cannot fit through this passage with that load.")
	g.Rooms["timber-room"] = timberRoom

	// LOWER-SHAFT
	lowerShaft := NewRoom(
		"lower-shaft",
		"Drafty Room",
		"This is a small drafty room in which is the bottom of a long shaft. To the south is a passageway and to the east a very narrow passage. In the shaft can be seen a heavy iron chain.",
	)
	lowerShaft.Flags.IsOutdoors = false
	lowerShaft.AddExit("south", "machine-room")
	lowerShaft.AddConditionalExit("out", "timber-room", "empty-handed", "You cannot fit through this passage with that load.")
	lowerShaft.AddConditionalExit("east", "timber-room", "empty-handed", "You cannot fit through this passage with that load.")
	g.Rooms["lower-shaft"] = lowerShaft

	// MACHINE-ROOM
	machineRoom := NewRoom(
		"machine-room",
		"Machine Room",
		"This is a large room full of assorted pieces of machinery. The room smells of burned resistors. Along one wall of the room are three buttons marked \"Start\", \"Launch\", and \"Lower\".",
	)
	machineRoom.AddExit("north", "lower-shaft")
	g.Rooms["machine-room"] = machineRoom

	// MINE-1
	mine1 := NewRoom(
		"mine-1",
		"Coal Mine",
		"This is a nondescript part of a coal mine.",
	)
	mine1.AddExit("north", "gas-room") // per ZIL
	mine1.AddExit("east", "mine-1")    // per ZIL (intentional self-loop)
	mine1.AddExit("ne", "mine-2")      // per ZIL
	// Note: ZIL has no west exit, gas-room east is one-way
	g.Rooms["mine-1"] = mine1

	// MINE-2
	mine2 := NewRoom(
		"mine-2",
		"Coal Mine",
		"This is a nondescript part of a coal mine.",
	)
	mine2.AddExit("north", "mine-2") // per ZIL (intentional self-loop)
	mine2.AddExit("south", "mine-1") // per ZIL
	mine2.AddExit("se", "mine-3")    // per ZIL
	// Note: ZIL has no sw or west exits - mine is intentionally confusing
	g.Rooms["mine-2"] = mine2

	// MINE-3
	mine3 := NewRoom(
		"mine-3",
		"Coal Mine",
		"This is a nondescript part of a coal mine.",
	)
	mine3.AddExit("south", "mine-3") // per ZIL (self-loop)
	mine3.AddExit("sw", "mine-4") // per ZIL
	mine3.AddExit("east", "mine-2") // per ZIL
	g.Rooms["mine-3"] = mine3

	// MINE-4
	mine4 := NewRoom(
		"mine-4",
		"Coal Mine",
		"This is a nondescript part of a coal mine.",
	)
	mine4.AddExit("north", "mine-3") // per ZIL
	mine4.AddExit("west", "mine-4") // per ZIL (self-loop)
	mine4.AddExit("down", "ladder-top") // per ZIL
	g.Rooms["mine-4"] = mine4

	// SLIDE-ROOM
	slideRoom := NewRoom(
		"slide-room",
		"Slide Room",
		"This is a small chamber, which appears to have been part of a coal mine. On the south wall of the chamber the letters \"Granite Wall\" are etched in the rock. To the east is a long passage, and there is a steep metal slide twisting downward. To the north is a small opening.",
	)
	slideRoom.AddExit("east", "cold-passage")
	slideRoom.AddExit("north", "mine-entrance")
	slideRoom.AddExit("down", "cellar")
	g.Rooms["slide-room"] = slideRoom
}
