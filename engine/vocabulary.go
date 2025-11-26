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
// Complete implementation with all ZIL verb synonyms
func (v *Vocabulary) initVerbs() {
	// System commands (V-VERBOSE, V-BRIEF, V-SUPER-BRIEF, etc.)
	v.addVerb("verbose", "verbose", "long")
	v.addVerb("brief", "brief", "normal")
	v.addVerb("superbrief", "superbrief", "super", "short", "terse")
	v.addVerb("diagnose", "diagnose")
	v.addVerb("inventory", "inventory", "i", "inv")
	v.addVerb("quit", "quit", "q")
	v.addVerb("restart", "restart")
	v.addVerb("restore", "restore", "load")
	v.addVerb("save", "save")
	v.addVerb("score", "score")
	v.addVerb("script", "script")
	v.addVerb("unscript", "unscript")
	v.addVerb("version", "version")
	v.addVerb("help", "help", "?")

	// Movement (V-WALK, V-ENTER, V-EXIT, V-CLIMB, etc.)
	v.addVerb("walk", "walk", "go", "run", "proceed", "step")
	v.addVerb("enter", "enter")
	v.addVerb("exit", "exit", "leave")
	v.addVerb("climb", "climb", "sit")
	v.addVerb("climb-up", "climb up") // multi-word
	v.addVerb("climb-down", "climb down") // multi-word
	v.addVerb("climb-on", "climb on") // multi-word
	v.addVerb("board", "board")
	v.addVerb("disembark", "disembark")
	v.addVerb("jump", "jump", "leap", "dive")
	v.addVerb("cross", "cross", "ford")
	v.addVerb("follow", "follow", "pursue", "chase", "come")
	v.addVerb("back", "back")
	v.addVerb("skip", "skip", "hop")
	v.addVerb("swim", "swim", "bathe", "wade")
	v.addVerb("stand", "stand")
	v.addVerb("stay", "stay")

	// Manipulation (V-TAKE, V-DROP, V-PUT, etc.)
	v.addVerb("take", "take", "get", "hold", "carry", "remove", "grab", "catch")
	v.addVerb("take", "pick up") // multi-word
	v.addVerb("drop", "drop", "release", "discard")
	v.addVerb("drop", "put down") // multi-word
	v.addVerb("put", "put", "stuff", "insert", "place", "hide")
	v.addVerb("put-on", "put on") // multi-word
	v.addVerb("give", "give", "donate", "offer", "feed")
	v.addVerb("throw", "throw", "hurl", "chuck", "toss")
	v.addVerb("move", "move", "shift")
	v.addVerb("raise", "raise", "lift")
	v.addVerb("lower", "lower")

	// Examination (V-EXAMINE, V-LOOK, V-READ, V-SEARCH, etc.)
	v.addVerb("examine", "examine", "x", "inspect", "describe", "what", "whats")
	v.addVerb("examine", "look at") // multi-word
	v.addVerb("look", "look", "l", "stare", "gaze")
	v.addVerb("look-in", "look in") // multi-word
	v.addVerb("look-on", "look on") // multi-word
	v.addVerb("read", "read", "skim")
	v.addVerb("search", "search")
	v.addVerb("find", "find", "where", "seek", "see")
	v.addVerb("count", "count")
	v.addVerb("listen", "listen", "hear")
	v.addVerb("smell", "smell", "sniff")
	v.addVerb("touch", "touch", "feel", "rub", "pat", "pet")

	// Container interaction (V-OPEN, V-CLOSE, V-LOCK, V-UNLOCK)
	v.addVerb("open", "open")
	v.addVerb("close", "close", "shut")
	v.addVerb("lock", "lock")
	v.addVerb("unlock", "unlock")

	// Light and fire (V-LAMP-ON, V-LAMP-OFF, V-BURN, etc.)
	v.addVerb("light", "light", "ignite")
	v.addVerb("turn", "turn", "set", "flip")
	v.addVerb("turn-on", "turn on", "switch on", "activate") // multi-word
	v.addVerb("turn-off", "turn off", "switch off", "deactivate") // multi-word
	v.addVerb("extinguish", "extinguish", "douse")
	v.addVerb("burn", "burn", "incinerate", "ignite")
	v.addVerb("melt", "melt", "liquify")

	// Combat and destruction (V-ATTACK, V-KILL, V-CUT, V-BREAK, etc.)
	v.addVerb("attack", "attack", "fight", "hurt", "injure", "hit")
	v.addVerb("kill", "kill", "murder", "slay", "dispatch")
	v.addVerb("break", "break", "smash", "mung", "destroy", "damage", "block")
	v.addVerb("cut", "cut", "slice", "pierce")
	v.addVerb("kick", "kick", "taunt")
	v.addVerb("stab", "stab")
	v.addVerb("strike", "strike")
	v.addVerb("swing", "swing", "thrust")
	v.addVerb("blast", "blast")

	// Consumption (V-EAT, V-DRINK)
	v.addVerb("drink", "drink", "imbibe", "swallow", "quaff", "sip")
	v.addVerb("drink-from", "drink from") // multi-word
	v.addVerb("eat", "eat", "consume", "taste", "bite", "devour")

	// Fluids (V-POUR, V-FILL, V-INFLATE, V-DEFLATE)
	v.addVerb("pour", "pour", "spill", "empty")
	v.addVerb("fill", "fill")
	v.addVerb("inflate", "inflate")
	v.addVerb("deflate", "deflate")
	v.addVerb("pump", "pump")
	v.addVerb("spray", "spray")
	v.addVerb("squeeze", "squeeze")

	// Physical manipulation (V-PUSH, V-PULL, V-RUB, V-TIE, etc.)
	v.addVerb("push", "push", "press")
	v.addVerb("pull", "pull", "tug", "yank")
	v.addVerb("tie", "tie", "fasten", "secure", "attach")
	v.addVerb("untie", "untie", "free", "release", "unfasten", "unattach", "unhook")
	v.addVerb("wave", "wave", "brandish")
	v.addVerb("shake", "shake")
	v.addVerb("knock", "knock", "rap")
	v.addVerb("ring", "ring", "peal")
	v.addVerb("brush", "brush", "clean")
	v.addVerb("lubricate", "lubricate", "oil", "grease")
	v.addVerb("plug", "plug", "glue", "patch", "repair", "fix")
	v.addVerb("puncture", "puncture")
	v.addVerb("dig", "dig", "excavate")
	v.addVerb("poke", "poke")
	v.addVerb("roll", "roll")
	v.addVerb("spin", "spin")
	v.addVerb("slide", "slide")
	v.addVerb("lean", "lean")
	v.addVerb("blow", "blow")
	v.addVerb("blow-out", "blow out") // multi-word
	v.addVerb("blow-up", "blow up") // multi-word
	v.addVerb("blow-in", "blow in") // multi-word

	// Wearing and equipment
	v.addVerb("wear", "wear")

	// Communication (V-TELL, V-HELLO, V-YELL, V-ANSWER, etc.)
	v.addVerb("tell", "tell", "ask")
	v.addVerb("say", "say", "speak")
	v.addVerb("hello", "hello", "hi")
	v.addVerb("yell", "yell", "scream", "shout")
	v.addVerb("answer", "answer", "reply")
	v.addVerb("command", "command")
	v.addVerb("curse", "curse", "shit", "fuck", "damn")

	// Magic and special (V-EXORCISE, V-PRAY, V-WISH, etc.)
	v.addVerb("exorcise", "exorcise", "banish", "cast", "drive", "begone")
	v.addVerb("pray", "pray")
	v.addVerb("wish", "wish")
	v.addVerb("incant", "incant", "chant")
	v.addVerb("enchant", "enchant")
	v.addVerb("disenchant", "disenchant")
	v.addVerb("mumble", "mumble", "sigh")
	v.addVerb("repent", "repent")

	// Time and state
	v.addVerb("wait", "wait", "z")
	v.addVerb("wake", "wake", "awake", "surprise", "startle")
	v.addVerb("sleep", "sleep")

	// Easter eggs and magic words
	v.addVerb("xyzzy", "xyzzy", "plugh") // plugh synonymous per ZIL
	v.addVerb("frobozz", "frobozz")
	v.addVerb("ulysses", "ulysses", "odysseus")
	v.addVerb("zork", "zork")
	v.addVerb("win", "win", "winnage")
	v.addVerb("echo", "echo")
	v.addVerb("chomp", "chomp", "lose", "barf")
	v.addVerb("bug", "bug")
	v.addVerb("treasure", "treasure", "temple")

	// Miscellaneous actions
	v.addVerb("apply", "apply")
	v.addVerb("make", "make")
	v.addVerb("hatch", "hatch")
	v.addVerb("launch", "launch")
	v.addVerb("activate", "activate")
	v.addVerb("send", "send")
	v.addVerb("talk", "talk")
	v.addVerb("play", "play")
	v.addVerb("wind", "wind")
	v.addVerb("pick", "pick")
	v.addVerb("rape", "rape", "molest")
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
	v.addObject("pot-of-gold", "pot", "gold", "pot-of-gold")

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
	v.addObject("book", "book", "guidebook", "guide", "prayer-book", "black-book", "prayer", "black")
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
