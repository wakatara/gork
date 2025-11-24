package engine

// Vocabulary stores all known words and their canonical forms
// This is based on the extensive vocabulary in ZIL (gsyntax.zil, 1dungeon.zil)
// ZIL had 684 vocabulary words!
type Vocabulary struct {
	verbs       map[string]string // synonym -> canonical verb
	objects     map[string]string // synonym -> canonical object
	prepositions map[string]bool   // valid prepositions
	directions  map[string]string // direction synonym -> canonical direction
}

// NewVocabulary creates and initializes the vocabulary from ZIL sources
func NewVocabulary() *Vocabulary {
	v := &Vocabulary{
		verbs:       make(map[string]string),
		objects:     make(map[string]string),
		prepositions: make(map[string]bool),
		directions:  make(map[string]string),
	}

	v.initVerbs()
	v.initObjects()
	v.initPrepositions()
	v.initDirections()

	return v
}

// initVerbs initializes verb synonyms from ZIL's syntax definitions (gsyntax.zil)
func (v *Vocabulary) initVerbs() {
	// Movement verbs (V?WALK)
	v.addVerb("walk", "walk", "go", "run", "proceed")

	// Manipulation verbs (V?TAKE, V?DROP, etc.)
	v.addVerb("take", "take", "get", "grab", "carry", "hold")
	v.addVerb("take", "pick up") // multi-word: "pick up" -> "take"
	v.addVerb("drop", "drop", "release", "discard")
	v.addVerb("drop", "put down") // multi-word: "put down" -> "drop"

	// Examination verbs (V?EXAMINE, V?LOOK)
	v.addVerb("examine", "examine", "x", "inspect", "describe", "what")
	v.addVerb("examine", "look at") // multi-word: "look at" -> "examine"
	v.addVerb("look", "look", "l")

	// Container interaction (V?OPEN, V?CLOSE)
	v.addVerb("open", "open", "unlock")
	v.addVerb("close", "close", "shut")

	// Reading (V?READ)
	v.addVerb("read", "read", "peruse")

	// Inventory (V?INVENTORY)
	v.addVerb("inventory", "inventory", "i", "inv")

	// Giving (V?GIVE)
	v.addVerb("give", "give", "offer", "hand")
	v.addVerb("put", "put", "place", "insert", "set")

	// Light control (V?LAMP-ON, V?LAMP-OFF)
	v.addVerb("light", "light", "ignite")
	v.addVerb("extinguish", "extinguish", "douse", "put-out")
	v.addVerb("turn", "turn", "switch", "rotate")

	// Combat (V?ATTACK, V?MUNG)
	v.addVerb("attack", "attack", "fight", "hit", "kill", "strike", "murder")
	v.addVerb("break", "break", "smash", "mung", "destroy")

	// System commands (V?QUIT, V?RESTART, V?SAVE, V?RESTORE)
	v.addVerb("quit", "quit", "q", "exit")
	v.addVerb("restart", "restart")
	v.addVerb("save", "save")
	v.addVerb("restore", "restore", "load")
	v.addVerb("score", "score")
	v.addVerb("version", "version")
	v.addVerb("help", "help", "?")

	// Verbosity commands (V?VERBOSE, V?BRIEF, V?SUPER-BRIEF)
	v.addVerb("verbose", "verbose", "long")
	v.addVerb("brief", "brief", "normal")
	v.addVerb("superbrief", "superbrief", "short", "terse")

	// Communication (V?TELL, V?HELLO)
	v.addVerb("tell", "tell", "say", "speak")
	v.addVerb("hello", "hello", "hi", "greet")

	// Waiting (V?WAIT)
	v.addVerb("wait", "wait", "z")

	// Magic words and special commands
	v.addVerb("xyzzy", "xyzzy")
	v.addVerb("plugh", "plugh")
	v.addVerb("frobozz", "frobozz")
	v.addVerb("ulysses", "ulysses", "odysseus") // V-ODYSSEUS in ZIL

	// More action verbs from gverbs.zil
	v.addVerb("climb", "climb", "scale")
	v.addVerb("drink", "drink", "quaff", "sip")
	v.addVerb("eat", "eat", "consume", "devour")
	v.addVerb("fill", "fill")
	v.addVerb("pour", "pour", "empty")
	v.addVerb("listen", "listen", "hear")
	v.addVerb("smell", "smell", "sniff")
	v.addVerb("touch", "touch", "feel", "rub")
	v.addVerb("wave", "wave", "brandish")
	v.addVerb("pull", "pull", "tug")
	v.addVerb("push", "push", "press")
	v.addVerb("move", "move", "shift")
	v.addVerb("turn-on", "turn-on", "activate", "switch-on")
	v.addVerb("turn-off", "turn-off", "deactivate", "switch-off")
	v.addVerb("enter", "enter", "board", "climb-in")
	v.addVerb("exit", "exit", "leave", "disembark", "climb-out")
	v.addVerb("search", "search")
	v.addVerb("jump", "jump", "leap")
	v.addVerb("knock", "knock", "rap")
	v.addVerb("pray", "pray")
	v.addVerb("dig", "dig", "excavate")
	v.addVerb("swim", "swim")
	v.addVerb("tie", "tie", "attach", "fasten")
	v.addVerb("untie", "untie", "unfasten", "detach")
	v.addVerb("blow", "blow")
	v.addVerb("ring", "ring")
	v.addVerb("count", "count")
	v.addVerb("unlock", "unlock")
	v.addVerb("lock", "lock")
	v.addVerb("inflate", "inflate", "blow-up", "fill")
	v.addVerb("deflate", "deflate")
	v.addVerb("plug", "plug", "repair", "patch")
}

// initObjects initializes object synonyms from ZIL object definitions (1dungeon.zil, gglobals.zil)
func (v *Vocabulary) initObjects() {
	// The opening scene objects (crucial for validation!)
	v.addObject("mailbox", "mailbox", "box", "letter-box")
	v.addObject("leaflet", "leaflet", "pamphlet", "booklet")
	v.addObject("lamp", "lamp", "lantern", "light")
	v.addObject("sword", "sword", "blade", "elvish")
	v.addObject("bottle", "bottle", "flask")
	v.addObject("water", "water", "h2o", "liquid")
	v.addObject("keys", "keys", "key", "set-of-keys")

	// White house area (WEST-OF-HOUSE, EAST-OF-HOUSE)
	v.addObject("white-house", "house", "building", "home", "white house", "white-house")
	v.addObject("kitchen-window", "window", "kitchen window", "kitchen-window")
	v.addObject("door", "door", "front-door", "entrance")
	v.addObject("board", "board", "boards", "planks")

	// Trophy case and living room
	v.addObject("trophy-case", "case", "trophy", "display", "trophy-case")
	v.addObject("table", "table", "desk")
	v.addObject("rug", "rug", "carpet", "mat")
	v.addObject("trap-door", "trap", "trapdoor", "hatch")

	// Treasures (there are 19 treasures worth 350 points total)
	v.addObject("coins", "coins", "bag", "bag-of-coins")
	v.addObject("chalice", "chalice", "goblet", "cup")
	v.addObject("painting", "painting", "picture", "canvas")
	v.addObject("jewels", "jewels", "gems", "jewel")
	v.addObject("diamond", "diamond")
	v.addObject("emerald", "emerald")
	v.addObject("sapphire", "sapphire")
	v.addObject("trident", "trident")
	v.addObject("egg", "egg", "jeweled-egg")
	v.addObject("sceptre", "sceptre", "scepter")
	v.addObject("bracelet", "bracelet")

	// Enemies and NPCs (from 1actions.zil)
	v.addObject("troll", "troll", "monster")
	v.addObject("thief", "thief", "robber", "bandit")
	v.addObject("cyclops", "cyclops", "giant")
	v.addObject("grue", "grue") // The famous grue!

	// Common objects
	v.addObject("me", "me", "myself", "self")
	v.addObject("it", "it", "them")
	v.addObject("all", "all", "everything")

	// Tools and utility items
	v.addObject("knife", "knife", "rusty-knife", "blade")
	v.addObject("axe", "axe", "ax")
	v.addObject("rope", "rope", "line", "cord")
	v.addObject("match", "match", "matches")
	v.addObject("candle", "candle", "candles")
	v.addObject("torch", "torch")
	v.addObject("wrench", "wrench", "spanner")
	v.addObject("screwdriver", "screwdriver")

	// Dungeon features
	v.addObject("grate", "grate", "grating", "bars")
	v.addObject("stairs", "stairs", "steps", "stairway", "staircase")
	v.addObject("ladder", "ladder")
	v.addObject("chimney", "chimney")
	v.addObject("slide", "slide")

	// Environment
	v.addObject("ground", "ground", "floor", "dirt", "earth")
	v.addObject("wall", "wall", "walls")
	v.addObject("ceiling", "ceiling")
	v.addObject("forest", "forest", "trees", "woods")
	v.addObject("river", "river", "stream")

	// Dam/reservoir area
	v.addObject("dam", "dam")
	v.addObject("reservoir", "reservoir", "lake")
	v.addObject("bolt", "bolt")
	v.addObject("bubble", "bubble")
	v.addObject("button", "button", "switch")

	// More items from 1dungeon.zil
	v.addObject("book", "book", "guidebook", "guide")
	v.addObject("bell", "bell")
	v.addObject("canary", "canary", "bird")
	v.addObject("garlic", "garlic", "clove")
	v.addObject("coffin", "coffin", "casket")
	v.addObject("basket", "basket")
	v.addObject("boat", "boat", "raft")
	v.addObject("lunch", "lunch", "sandwich", "food")
	v.addObject("troll", "troll", "nasty", "monster")
	v.addObject("thief", "thief", "robber", "bandit")
	v.addObject("cyclops", "cyclops", "giant")
	v.addObject("bat", "bat", "vampire-bat")
	v.addObject("ghosts", "ghosts", "spirits", "ghost", "spirit")
	v.addObject("rope", "rope")
	v.addObject("bell", "bell")
	v.addObject("tree", "tree", "trees")
	v.addObject("lever", "lever", "handle")
	v.addObject("whistle", "whistle")

	// Newly added items (37 total)
	v.addObject("platinum-bar", "bar", "platinum", "platinum-bar")
	v.addObject("sapphire", "sapphire", "gem")
	v.addObject("ivory-torch", "ivory-torch", "ivory")
	v.addObject("trunk-of-jewels", "trunk-of-jewels", "jewels")
	v.addObject("pearl", "pearl")
	v.addObject("clam", "clam", "shell")
	v.addObject("matches", "matches", "matchbook", "book-of-matches")
	v.addObject("mirror", "mirror", "looking-glass")
	v.addObject("pile-of-leaves", "pile", "pile-of-leaves")
	v.addObject("grating", "grating")
	v.addObject("cyclops-corpse", "cyclops-corpse", "corpse", "body")
	v.addObject("thief-corpse", "thief-corpse")
	v.addObject("reservoir", "reservoir", "lake")
	v.addObject("stream", "stream", "brook")
	v.addObject("glacier", "glacier", "ice")
	v.addObject("slide", "slide")
	v.addObject("brick", "brick")
	v.addObject("statue", "statue", "idol")
	v.addObject("air-pump", "air-pump")
	v.addObject("cyclops-treasure", "cyclops-treasure", "treasure-chest")
	v.addObject("crystal-sphere", "crystal-sphere", "crystal", "sphere", "ball", "crystal-ball")
	v.addObject("rainbow-arc", "rainbow-arc")
	v.addObject("shrunken-heads", "shrunken-heads", "heads", "head")
	v.addObject("flask", "flask", "vial")
	v.addObject("sword-holder", "sword-holder", "holder", "mount")
	v.addObject("crypt", "crypt", "tomb")
	v.addObject("granite-wall", "granite-wall", "wall", "granite")
	v.addObject("wooden-door", "wooden-door")
	v.addObject("iron-door", "iron-door")
	v.addObject("candle", "candle")
	v.addObject("chain", "chain")
	v.addObject("hook", "hook")
	v.addObject("pillar", "pillar", "column")
	v.addObject("altar-cloth", "altar-cloth", "cloth")
	v.addObject("stick", "stick", "walking-stick")
	v.addObject("volcano", "volcano")
	v.addObject("railing", "railing", "rail")

	// Dam control buttons
	v.addObject("yellow-button", "yellow-button", "yellow button", "yellow")
	v.addObject("blue-button", "blue-button", "blue button", "blue")
	v.addObject("brown-button", "brown-button", "brown button", "brown")
	v.addObject("red-button", "red-button", "red button", "red")

	// Machine control buttons
	v.addObject("start-button", "start-button", "start button", "start")
	v.addObject("launch-button", "launch-button", "launch button", "launch")
	v.addObject("lower-button", "lower-button", "lower button", "lower")
}

// initPrepositions initializes valid prepositions from ZIL
// (from gparser.zil - there are 18 prepositions)
func (v *Vocabulary) initPrepositions() {
	preps := []string{
		"in", "into", "inside",
		"on", "onto", "upon",
		"at", "to", "toward", "towards",
		"with", "using",
		"from",
		"through", "across",
		"under", "underneath", "beneath",
		"behind",
		"over",
		"off",
		"for",
		"about",
		"around",
		"down",
		"up",
		"out",
		"away",
	}

	for _, prep := range preps {
		v.prepositions[prep] = true
	}
}

// initDirections initializes direction synonyms (from ZIL: <DIRECTIONS ...>)
func (v *Vocabulary) initDirections() {
	v.addDirection("north", "north", "n")
	v.addDirection("south", "south", "s")
	v.addDirection("east", "east", "e")
	v.addDirection("west", "west", "w")
	v.addDirection("northeast", "northeast", "ne")
	v.addDirection("northwest", "northwest", "nw")
	v.addDirection("southeast", "southeast", "se")
	v.addDirection("southwest", "southwest", "sw")
	v.addDirection("up", "up", "u")
	v.addDirection("down", "down", "d")
	v.addDirection("in", "in")
	v.addDirection("out", "out")
}

// Helper methods to add words with multiple synonyms

func (v *Vocabulary) addVerb(canonical string, synonyms ...string) {
	for _, syn := range synonyms {
		v.verbs[syn] = canonical
	}
}

func (v *Vocabulary) addObject(canonical string, synonyms ...string) {
	for _, syn := range synonyms {
		v.objects[syn] = canonical
	}
}

func (v *Vocabulary) addDirection(canonical string, synonyms ...string) {
	for _, syn := range synonyms {
		v.directions[syn] = canonical
	}
}

// Lookup methods

func (v *Vocabulary) GetVerb(word string) string {
	return v.verbs[word]
}

func (v *Vocabulary) GetObject(word string) string {
	return v.objects[word]
}

func (v *Vocabulary) GetDirection(word string) string {
	return v.directions[word]
}

func (v *Vocabulary) IsPreposition(word string) bool {
	return v.prepositions[word]
}

// IsKnownWord checks if a word exists in any vocabulary category
func (v *Vocabulary) IsKnownWord(word string) bool {
	if v.GetVerb(word) != "" {
		return true
	}
	if v.GetObject(word) != "" {
		return true
	}
	if v.GetDirection(word) != "" {
		return true
	}
	return v.IsPreposition(word)
}
