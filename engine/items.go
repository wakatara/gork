package engine

// InitializeItems creates all items from Zork I
// Ported from 1dungeon.zil OBJECT definitions
func InitializeItems(g *GameV2) {
	createTreasures(g)
	createWeapons(g)
	createTools(g)
	createContainers(g)
	createReadableItems(g)
	createLightSources(g)
	createFoodAndDrink(g)
	createFixedObjects(g)
	createSceneryObjects(g)
	createMiscItems(g)
}

// createTreasures creates all treasure items (VALUE > 0)
func createTreasures(g *GameV2) {
	// DIAMOND - 10 points
	diamond := NewItem("diamond", "huge diamond", "There is an enormous diamond here!")
	diamond.Aliases = []string{"diamond", "treasure"}
	diamond.Flags.IsTakeable = true
	diamond.Flags.IsTreasure = true
	diamond.Value = 10
	g.Items["diamond"] = diamond

	// EMERALD - 5 points (in buoy)
	emerald := NewItem("emerald", "large emerald", "There is a large emerald here.")
	emerald.Aliases = []string{"emerald", "treasure"}
	emerald.Location = "buoy"
	emerald.Flags.IsTakeable = true
	emerald.Flags.IsTreasure = true
	emerald.Value = 5
	g.Items["emerald"] = emerald

	// CHALICE - treasure and container
	chalice := NewItem("chalice", "chalice", "There is a silver chalice, intricately engraved, here.")
	chalice.Aliases = []string{"chalice", "cup", "silver", "treasure"}
	chalice.Location = "treasure-room"
	chalice.Flags.IsTakeable = true
	chalice.Flags.IsTreasure = true
	chalice.Flags.IsContainer = true
	chalice.Value = 10
	g.Items["chalice"] = chalice

	// JADE FIGURINE
	jade := NewItem("jade", "jade figurine", "There is a precious jade figurine here!")
	jade.Aliases = []string{"jade", "figurine", "treasure"}
	jade.Flags.IsTakeable = true
	jade.Flags.IsTreasure = true
	jade.Value = 5
	g.Items["jade"] = jade

	// BAG OF COINS
	coins := NewItem("coins", "bag of coins", "There is a bag of coins here.")
	coins.Aliases = []string{"coins", "bag", "bag-of-coins", "treasure"}
	coins.Flags.IsTakeable = true
	coins.Flags.IsTreasure = true
	coins.Value = 5
	g.Items["coins"] = coins

	// PAINTING
	painting := NewItem("painting", "painting", "Fortunately, there is still one chance for you to be a vandal, for on the far wall is a painting.")
	painting.Aliases = []string{"painting", "treasure"}
	painting.Location = "gallery"
	painting.Flags.IsTakeable = true
	painting.Flags.IsTreasure = true
	painting.Value = 4
	g.Items["painting"] = painting

	// BRACELET
	bracelet := NewItem("bracelet", "sapphire bracelet", "There is a sapphire bracelet here.")
	bracelet.Aliases = []string{"bracelet", "sapphire", "treasure"}
	bracelet.Flags.IsTakeable = true
	bracelet.Flags.IsTreasure = true
	bracelet.Value = 5
	g.Items["bracelet"] = bracelet

	// BAUBLE (jeweled egg bauble)
	bauble := NewItem("bauble", "jeweled bauble", "There is a small jeweled bauble here.")
	bauble.Aliases = []string{"bauble", "treasure"}
	bauble.Flags.IsTakeable = true
	bauble.Flags.IsTreasure = true
	bauble.Value = 5
	g.Items["bauble"] = bauble

	// SCARAB
	scarab := NewItem("scarab", "beautiful scarab", "There is a beautiful scarab here.")
	scarab.Aliases = []string{"scarab", "treasure"}
	scarab.Flags.IsTakeable = true
	scarab.Flags.IsTreasure = true
	scarab.Value = 5
	g.Items["scarab"] = scarab

	// POT OF GOLD
	pot := NewItem("pot-of-gold", "pot of gold", "At the end of the rainbow is a pot of gold.")
	pot.Aliases = []string{"pot", "gold", "pot-of-gold", "treasure"}
	pot.Location = "end-of-rainbow"
	pot.Flags.IsTakeable = true
	pot.Flags.IsTreasure = true
	pot.Value = 10
	g.Items["pot-of-gold"] = pot

	// TRIDENT
	trident := NewItem("trident", "crystal trident", "There is a crystal trident here.")
	trident.Aliases = []string{"trident", "crystal", "treasure"}
	trident.Location = "falls"
	trident.Flags.IsTakeable = true
	trident.Flags.IsTreasure = true
	trident.Flags.IsWeapon = true
	trident.Value = 4
	g.Items["trident"] = trident

	// SCEPTRE
	sceptre := NewItem("sceptre", "sceptre", "There is a sceptre, probably that of ancient Egypt itself, here.")
	sceptre.Aliases = []string{"sceptre", "scepter", "treasure"}
	sceptre.Flags.IsTakeable = true
	sceptre.Flags.IsTreasure = true
	sceptre.Flags.IsWeapon = true
	sceptre.Value = 4
	g.Items["sceptre"] = sceptre

	// EGG (containing bauble) - treasure
	egg := NewItem("egg", "jeweled egg", "There is a large nest here, with a jeweled egg resting in it.")
	egg.Aliases = []string{"egg", "treasure"}
	egg.Location = "nest"
	egg.Flags.IsTakeable = true
	egg.Flags.IsContainer = true
	egg.Flags.IsTreasure = true
	egg.Value = 5
	g.Items["egg"] = egg
}

// createWeapons creates all weapon items
func createWeapons(g *GameV2) {
	// SWORD (elvish)
	sword := NewItem("sword", "elvish sword", "There is an elvish sword here.")
	sword.Aliases = []string{"sword", "blade", "elvish"}
	sword.Location = "living-room"
	sword.Flags.IsTakeable = true
	sword.Flags.IsWeapon = true
	g.Items["sword"] = sword
	g.Rooms["living-room"].AddItem("sword")

	// KNIFE
	knife := NewItem("knife", "knife", "There is a knife here.")
	knife.Aliases = []string{"knife"}
	knife.Flags.IsTakeable = true
	knife.Flags.IsWeapon = true
	g.Items["knife"] = knife

	// RUSTY KNIFE
	rustyKnife := NewItem("rusty-knife", "rusty knife", "There is a rusty knife here.")
	rustyKnife.Aliases = []string{"rusty-knife", "knife", "rusty"}
	rustyKnife.Flags.IsTakeable = true
	rustyKnife.Flags.IsWeapon = true
	g.Items["rusty-knife"] = rustyKnife

	// STILETTO
	stiletto := NewItem("stiletto", "stiletto", "There is a wicked-looking stiletto here.")
	stiletto.Aliases = []string{"stiletto", "dagger"}
	stiletto.Flags.IsTakeable = true
	stiletto.Flags.IsWeapon = true
	g.Items["stiletto"] = stiletto

	// AXE (bloody)
	axe := NewItem("axe", "bloody axe", "There is a bloody axe here.")
	axe.Aliases = []string{"axe", "ax"}
	axe.Flags.IsTakeable = true
	axe.Flags.IsWeapon = true
	g.Items["axe"] = axe
}

// createTools creates all tool items
func createTools(g *GameV2) {
	// PUMP (air pump for boat)
	pump := NewItem("pump", "air pump", "There is a hand-held air pump here.")
	pump.Aliases = []string{"pump", "air-pump"}
	pump.Flags.IsTakeable = true
	g.Items["pump"] = pump

	// SCREWDRIVER
	screwdriver := NewItem("screwdriver", "screwdriver", "There is a screwdriver here.")
	screwdriver.Aliases = []string{"screwdriver"}
	screwdriver.Flags.IsTakeable = true
	g.Items["screwdriver"] = screwdriver

	// WRENCH
	wrench := NewItem("wrench", "wrench", "There is a wrench here.")
	wrench.Aliases = []string{"wrench"}
	wrench.Flags.IsTakeable = true
	g.Items["wrench"] = wrench

	// PUTTY
	putty := NewItem("putty", "putty", "There is a tube of putty here.")
	putty.Aliases = []string{"putty", "tube"}
	putty.Flags.IsTakeable = true
	g.Items["putty"] = putty

	// SHOVEL
	shovel := NewItem("shovel", "shovel", "There is a shovel here.")
	shovel.Aliases = []string{"shovel", "spade"}
	shovel.Flags.IsTakeable = true
	g.Items["shovel"] = shovel

	// ROPE
	rope := NewItem("rope", "rope", "There is a rope here.")
	rope.Aliases = []string{"rope"}
	rope.Flags.IsTakeable = true
	g.Items["rope"] = rope
}

// createContainers creates all container items
func createContainers(g *GameV2) {
	// MAILBOX (starting item)
	mailbox := NewItem("mailbox", "small mailbox", "It's a small mailbox.")
	mailbox.Aliases = []string{"mailbox", "box"}
	mailbox.Location = "west-of-house"
	mailbox.Flags.IsContainer = true
	mailbox.Flags.IsOpen = true
	mailbox.Flags.IsTransparent = true
	g.Items["mailbox"] = mailbox
	g.Rooms["west-of-house"].AddItem("mailbox")

	// TROPHY CASE (in living room)
	trophyCase := NewItem("trophy-case", "trophy case", "There is a trophy case here.")
	trophyCase.Aliases = []string{"case", "trophy-case", "trophy"}
	trophyCase.Location = "living-room"
	trophyCase.Flags.IsContainer = true
	trophyCase.Flags.IsOpen = false
	trophyCase.Flags.IsTransparent = true
	g.Items["trophy-case"] = trophyCase
	g.Rooms["living-room"].AddItem("trophy-case")

	// BOTTLE (glass bottle)
	bottle := NewItem("bottle", "glass bottle", "There is a glass bottle here.")
	bottle.Aliases = []string{"bottle", "glass"}
	bottle.Flags.IsTakeable = true
	bottle.Flags.IsContainer = true
	bottle.Flags.IsTransparent = true
	g.Items["bottle"] = bottle

	// COFFIN
	coffin := NewItem("coffin", "coffin", "There is a coffin here.")
	coffin.Aliases = []string{"coffin"}
	coffin.Location = "egypt-room"
	coffin.Flags.IsTakeable = true
	coffin.Flags.IsContainer = true
	coffin.Flags.IsOpen = false
	g.Items["coffin"] = coffin

	// SANDWICH BAG
	bag := NewItem("sandwich-bag", "brown bag", "There is a brown bag here.")
	bag.Aliases = []string{"bag", "brown-bag", "sandwich-bag"}
	bag.Flags.IsTakeable = true
	bag.Flags.IsContainer = true
	g.Items["sandwich-bag"] = bag

	// LARGE BAG (for carrying treasures)
	largeBag := NewItem("large-bag", "large bag", "There is a large leather bag here.")
	largeBag.Aliases = []string{"bag", "large-bag", "leather-bag"}
	largeBag.Flags.IsTakeable = true
	largeBag.Flags.IsContainer = true
	g.Items["large-bag"] = largeBag

	// NEST (bird's nest containing egg)
	nest := NewItem("nest", "bird's nest", "There is a bird's nest here.")
	nest.Aliases = []string{"nest"}
	nest.Location = "up-a-tree"
	nest.Flags.IsTakeable = true
	nest.Flags.IsContainer = true
	nest.Flags.IsOpen = true
	g.Items["nest"] = nest

	// TOOL CHEST
	toolChest := NewItem("tool-chest", "tool chest", "There is a tool chest here.")
	toolChest.Aliases = []string{"chest", "tool-chest"}
	toolChest.Flags.IsContainer = true
	toolChest.Flags.IsOpen = true
	g.Items["tool-chest"] = toolChest

	// TUBE (for putty)
	tube := NewItem("tube", "tube", "There is a small tube here.")
	tube.Aliases = []string{"tube"}
	tube.Flags.IsTakeable = true
	tube.Flags.IsContainer = true
	tube.Flags.IsReadable = true
	g.Items["tube"] = tube

	// TRUNK
	trunk := NewItem("trunk", "trunk", "There is a trunk here.")
	trunk.Aliases = []string{"trunk"}
	trunk.Flags.IsTakeable = true
	trunk.Flags.IsContainer = true
	g.Items["trunk"] = trunk
}

// createReadableItems creates all readable items
func createReadableItems(g *GameV2) {
	// LEAFLET (in mailbox)
	leaflet := NewItem("leaflet", "leaflet", `"WELCOME TO ZORK!

ZORK is a game of adventure, danger, and low cunning. In it you will explore some of the most amazing territory ever seen by mortals. No computer should be without one!"`)
	leaflet.Aliases = []string{"leaflet", "pamphlet", "booklet"}
	leaflet.Location = "mailbox"
	leaflet.Flags.IsTakeable = true
	leaflet.Flags.IsReadable = true
	g.Items["leaflet"] = leaflet

	// BOOK (prayer book)
	book := NewItem("book", "prayer book", "There is a prayer book here.")
	book.Aliases = []string{"book", "prayer-book"}
	book.Flags.IsTakeable = true
	book.Flags.IsReadable = true
	book.Flags.IsContainer = true
	g.Items["book"] = book

	// ADVERTISEMENT
	advertisement := NewItem("advertisement", "advertisement", "There is an advertisement here.")
	advertisement.Aliases = []string{"advertisement", "ad"}
	advertisement.Flags.IsTakeable = true
	advertisement.Flags.IsReadable = true
	g.Items["advertisement"] = advertisement

	// GUIDE (guide book)
	guide := NewItem("guide", "tour guide", `"A Tour of the Great Underground Empire" by Flood Control Dam #3 Public Relations Council`)
	guide.Aliases = []string{"guide", "guidebook", "book"}
	guide.Flags.IsTakeable = true
	guide.Flags.IsReadable = true
	g.Items["guide"] = guide

	// MAP
	mapItem := NewItem("map", "map", "There is a map here.")
	mapItem.Aliases = []string{"map"}
	mapItem.Flags.IsTakeable = true
	mapItem.Flags.IsReadable = true
	g.Items["map"] = mapItem

	// BOAT LABEL
	boatLabel := NewItem("boat-label", "boat label", "There is a label on the boat.")
	boatLabel.Aliases = []string{"label", "boat-label"}
	boatLabel.Flags.IsTakeable = true
	boatLabel.Flags.IsReadable = true
	g.Items["boat-label"] = boatLabel

	// MATCH (matchbook)
	match := NewItem("match", "matchbook", "There is a matchbook here.")
	match.Aliases = []string{"match", "matchbook", "matches"}
	match.Flags.IsTakeable = true
	match.Flags.IsReadable = true
	g.Items["match"] = match

	// OWNERS MANUAL
	manual := NewItem("owners-manual", "owner's manual", "There is an owner's manual here.")
	manual.Aliases = []string{"manual", "owners-manual"}
	manual.Flags.IsTakeable = true
	manual.Flags.IsReadable = true
	g.Items["owners-manual"] = manual
}

// createLightSources creates items that provide light
func createLightSources(g *GameV2) {
	// LAMP (brass lantern)
	lamp := NewItem("lamp", "brass lantern", "The brass lantern is on.")
	lamp.Aliases = []string{"lamp", "lantern", "light"}
	lamp.Location = "living-room"
	lamp.Flags.IsTakeable = true
	lamp.Flags.IsLightSource = true
	lamp.Flags.IsLit = true
	g.Items["lamp"] = lamp
	g.Rooms["living-room"].AddItem("lamp")

	// TORCH
	torch := NewItem("torch", "torch", "There is a burning torch here.")
	torch.Aliases = []string{"torch"}
	torch.Location = "temple"
	torch.Flags.IsTakeable = true
	torch.Flags.IsLightSource = true
	torch.Flags.IsLit = true
	g.Items["torch"] = torch

	// CANDLES
	candles := NewItem("candles", "pair of candles", "There is a pair of candles here.")
	candles.Aliases = []string{"candles", "candle"}
	candles.Location = "entrance-to-hades"
	candles.Flags.IsTakeable = true
	candles.Flags.IsLightSource = true
	candles.Flags.IsLit = false
	g.Items["candles"] = candles

	// BURNED OUT LANTERN
	burnedLamp := NewItem("burned-out-lantern", "burned-out lantern", "There is a burned-out lantern here.")
	burnedLamp.Aliases = []string{"burned-out-lantern", "lantern"}
	burnedLamp.Flags.IsTakeable = true
	g.Items["burned-out-lantern"] = burnedLamp
}

// createFoodAndDrink creates edible and drinkable items
func createFoodAndDrink(g *GameV2) {
	// LUNCH (in sandwich bag)
	lunch := NewItem("lunch", "lunch", "There is a lunch here.")
	lunch.Aliases = []string{"lunch", "sandwich"}
	lunch.Location = "sandwich-bag"
	lunch.Flags.IsTakeable = true
	lunch.Flags.IsEdible = true
	g.Items["lunch"] = lunch

	// GARLIC
	garlic := NewItem("garlic", "clove of garlic", "There is a clove of garlic here.")
	garlic.Aliases = []string{"garlic", "clove"}
	garlic.Flags.IsTakeable = true
	garlic.Flags.IsEdible = true
	g.Items["garlic"] = garlic

	// WATER (in bottle or stream)
	water := NewItem("water", "quantity of water", "There is some water here.")
	water.Aliases = []string{"water", "quantity"}
	water.Flags.IsTakeable = true
	water.Flags.IsDrinkable = true
	g.Items["water"] = water
}

// createFixedObjects creates non-takeable interactive items
func createFixedObjects(g *GameV2) {
	// KITCHEN WINDOW
	window := NewItem("kitchen-window", "small window", "The window is slightly ajar.")
	window.Aliases = []string{"window", "kitchen-window"}
	window.Location = "behind-house"
	window.Flags.IsContainer = false
	window.Flags.IsTakeable = false
	window.Flags.IsOpen = false
	g.Items["kitchen-window"] = window
	g.Rooms["behind-house"].AddItem("kitchen-window")

	// RUG (conceals trap door)
	rug := NewItem("rug", "oriental rug", "There is a large oriental rug here.")
	rug.Aliases = []string{"rug", "oriental-rug"}
	rug.Location = "living-room"
	rug.Flags.IsTakeable = false
	g.Items["rug"] = rug
	g.Rooms["living-room"].AddItem("rug")

	// TRAP DOOR (under rug)
	trapDoor := NewItem("trap-door", "trap door", "There is a trap door here.")
	trapDoor.Aliases = []string{"trap-door", "door", "trapdoor"}
	trapDoor.Location = "living-room"
	trapDoor.Flags.IsTakeable = false
	trapDoor.Flags.IsOpen = false
	g.Items["trap-door"] = trapDoor

	// GRATE (above cellar)
	grate := NewItem("grate", "grating", "There is a grating securely fastened into the ceiling.")
	grate.Aliases = []string{"grate", "grating"}
	grate.Location = "grating-room"
	grate.Flags.IsTakeable = false
	g.Items["grate"] = grate

	// BUTTONS (control panel)
	yellowButton := NewItem("yellow-button", "yellow button", "There is a yellow button here.")
	yellowButton.Aliases = []string{"yellow-button", "yellow", "button"}
	yellowButton.Flags.IsTakeable = false
	g.Items["yellow-button"] = yellowButton

	brownButton := NewItem("brown-button", "brown button", "There is a brown button here.")
	brownButton.Aliases = []string{"brown-button", "brown", "button"}
	brownButton.Flags.IsTakeable = false
	g.Items["brown-button"] = brownButton

	redButton := NewItem("red-button", "red button", "There is a red button here.")
	redButton.Aliases = []string{"red-button", "red", "button"}
	redButton.Flags.IsTakeable = false
	g.Items["red-button"] = redButton

	blueButton := NewItem("blue-button", "blue button", "There is a blue button here.")
	blueButton.Aliases = []string{"blue-button", "blue", "button"}
	blueButton.Flags.IsTakeable = false
	g.Items["blue-button"] = blueButton

	// ALTAR
	altar := NewItem("altar", "altar", "There is a marble altar here.")
	altar.Aliases = []string{"altar"}
	altar.Location = "temple"
	altar.Flags.IsTakeable = false
	g.Items["altar"] = altar

	// BELL (in belfry)
	bell := NewItem("bell", "bell", "There is a large bell here.")
	bell.Aliases = []string{"bell"}
	bell.Location = "belfry"
	bell.Flags.IsTakeable = true
	g.Items["bell"] = bell

	// MACHINE (dam control)
	machine := NewItem("machine", "machine", "There is a massive machine here.")
	machine.Aliases = []string{"machine"}
	machine.Location = "machine-room"
	machine.Flags.IsTakeable = false
	machine.Flags.IsContainer = true
	g.Items["machine"] = machine

	// PEDESTAL
	pedestal := NewItem("pedestal", "pedestal", "There is a pedestal here.")
	pedestal.Aliases = []string{"pedestal"}
	pedestal.Flags.IsTakeable = false
	g.Items["pedestal"] = pedestal
}

// createSceneryObjects creates purely descriptive objects
func createSceneryObjects(g *GameV2) {
	// WHITE HOUSE
	whiteHouse := NewItem("white-house", "white house", "The house is a beautiful colonial house which is painted white. It is clear that the owners must have been extremely wealthy.")
	whiteHouse.Aliases = []string{"house", "white-house", "colonial"}
	whiteHouse.Flags.IsTakeable = false
	g.Items["white-house"] = whiteHouse

	// FOREST (scenery)
	forest := NewItem("forest", "forest", "You are in a forest, with trees in all directions.")
	forest.Aliases = []string{"forest", "trees", "tree"}
	forest.Flags.IsTakeable = false
	g.Items["forest"] = forest

	// MOUNTAIN RANGE
	mountains := NewItem("mountains", "mountain range", "The mountains are impassable.")
	mountains.Aliases = []string{"mountains", "mountain", "range"}
	mountains.Flags.IsTakeable = false
	g.Items["mountains"] = mountains

	// RAINBOW
	rainbow := NewItem("rainbow", "rainbow", "The rainbow seems to have its foot in the vicinity of the building.")
	rainbow.Aliases = []string{"rainbow"}
	rainbow.Location = "canyon-view"
	rainbow.Flags.IsTakeable = false
	g.Items["rainbow"] = rainbow

	// RIVER
	river := NewItem("river", "river", "The Frigid River flows through here.")
	river.Aliases = []string{"river", "frigid-river"}
	river.Flags.IsTakeable = false
	g.Items["river"] = river

	// ENGRAVINGS (on walls)
	engravings := NewItem("engravings", "engravings", "The engravings were incised in the living rock of the cave wall by an unknown hand.")
	engravings.Aliases = []string{"engravings", "inscription"}
	engravings.Flags.IsTakeable = false
	engravings.Flags.IsReadable = true
	g.Items["engravings"] = engravings

	// LEAVES
	leaves := NewItem("leaves", "pile of leaves", "There is a pile of leaves here.")
	leaves.Aliases = []string{"leaves", "pile"}
	leaves.Flags.IsTakeable = false
	g.Items["leaves"] = leaves

	// SAND
	sand := NewItem("sand", "sand", "There is sand here.")
	sand.Aliases = []string{"sand"}
	sand.Flags.IsTakeable = false
	g.Items["sand"] = sand
}

// createMiscItems creates various other items
func createMiscItems(g *GameV2) {
	// INFLATABLE BOAT
	boat := NewItem("boat", "inflatable boat", "There is an inflatable boat here.")
	boat.Aliases = []string{"boat", "inflatable-boat", "raft"}
	boat.Flags.IsTakeable = true
	g.Items["boat"] = boat

	// INFLATED BOAT (boat when inflated)
	inflatedBoat := NewItem("inflated-boat", "inflated boat", "There is an inflated boat here.")
	inflatedBoat.Aliases = []string{"boat", "inflated-boat"}
	inflatedBoat.Flags.IsTakeable = true
	g.Items["inflated-boat"] = inflatedBoat

	// PUNCTURED BOAT
	puncturedBoat := NewItem("punctured-boat", "punctured boat", "There is a punctured boat here.")
	puncturedBoat.Aliases = []string{"boat", "punctured-boat"}
	puncturedBoat.Flags.IsTakeable = true
	g.Items["punctured-boat"] = puncturedBoat

	// SKULL
	skull := NewItem("skull", "skull", "There is a skull here.")
	skull.Aliases = []string{"skull"}
	skull.Flags.IsTakeable = true
	g.Items["skull"] = skull

	// BONES
	bones := NewItem("bones", "pile of bones", "There is a pile of bones here.")
	bones.Aliases = []string{"bones", "pile"}
	bones.Flags.IsTakeable = false
	g.Items["bones"] = bones

	// COAL
	coal := NewItem("coal", "pile of coal", "There is a pile of coal here.")
	coal.Aliases = []string{"coal", "pile"}
	coal.Location = "coal-mine-4"
	coal.Flags.IsTakeable = true
	g.Items["coal"] = coal

	// TIMBER
	timbers := NewItem("timbers", "timber", "There are timber supports here.")
	timbers.Aliases = []string{"timber", "timbers"}
	timbers.Flags.IsTakeable = true
	g.Items["timbers"] = timbers

	// LADDER
	ladder := NewItem("ladder", "wooden ladder", "There is a wooden ladder here.")
	ladder.Aliases = []string{"ladder"}
	ladder.Flags.IsTakeable = true
	g.Items["ladder"] = ladder

	// CANARY (bird)
	canary := NewItem("canary", "canary", "There is a canary here, singing cheerfully.")
	canary.Aliases = []string{"canary", "bird"}
	canary.Flags.IsTakeable = true
	g.Items["canary"] = canary

	// BROKEN CANARY
	brokenCanary := NewItem("broken-canary", "broken canary", "There is a dead canary here.")
	brokenCanary.Aliases = []string{"canary", "bird", "broken-canary"}
	brokenCanary.Flags.IsTakeable = true
	g.Items["broken-canary"] = brokenCanary

	// BROKEN EGG
	brokenEgg := NewItem("broken-egg", "broken egg", "There is a broken jeweled egg here.")
	brokenEgg.Aliases = []string{"egg", "broken-egg"}
	brokenEgg.Flags.IsTakeable = true
	brokenEgg.Flags.IsContainer = true
	brokenEgg.Flags.IsOpen = true
	g.Items["broken-egg"] = brokenEgg

	// BUOY
	buoy := NewItem("buoy", "red buoy", "There is a red buoy here (probably a warning).")
	buoy.Aliases = []string{"buoy", "red-buoy"}
	buoy.Location = "reservoir-south"
	buoy.Flags.IsTakeable = false
	g.Items["buoy"] = buoy

	// RAISED BASKET
	raisedBasket := NewItem("raised-basket", "wicker basket", "There is a large wicker basket here, raised to the ceiling.")
	raisedBasket.Aliases = []string{"basket", "wicker-basket", "raised-basket"}
	raisedBasket.Flags.IsTakeable = false
	raisedBasket.Flags.IsContainer = true
	raisedBasket.Flags.IsOpen = true
	raisedBasket.Flags.IsTransparent = true
	g.Items["raised-basket"] = raisedBasket

	// LOWERED BASKET
	loweredBasket := NewItem("lowered-basket", "wicker basket", "There is a large wicker basket here.")
	loweredBasket.Aliases = []string{"basket", "wicker-basket", "lowered-basket"}
	loweredBasket.Flags.IsTakeable = false
	g.Items["lowered-basket"] = loweredBasket

	// PRAYER (of protection)
	prayer := NewItem("prayer", "prayer", "The prayer seems to be a plea for protection.")
	prayer.Aliases = []string{"prayer"}
	prayer.Flags.IsTakeable = false
	prayer.Flags.IsReadable = true
	g.Items["prayer"] = prayer

	// KEYS
	keys := NewItem("keys", "set of keys", "There is a set of keys here.")
	keys.Aliases = []string{"keys", "key"}
	keys.Flags.IsTakeable = true
	g.Items["keys"] = keys
}
