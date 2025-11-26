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
	// DIAMOND - 10 points - LDESC from ZIL
	diamond := NewItem("diamond", "huge diamond", "The diamond is perfectly cut and huge.")
	diamond.Aliases = []string{"diamond", "treasure", "huge", "enormous"}
	diamond.RoomDescription = "There is an enormous diamond (perfectly cut) here."
	diamond.Flags.IsTakeable = true
	diamond.Flags.IsTreasure = true
	diamond.Value = 10
	g.Items["diamond"] = diamond

	// EMERALD - 10 points (in buoy) - ZIL TVALUE 10
	emerald := NewItem("emerald", "large emerald", "The emerald is large and beautifully cut.")
	emerald.Aliases = []string{"emerald", "treasure", "large"}
	emerald.Location = "buoy"
	emerald.RoomDescription = "There is a large emerald here."
	emerald.Flags.IsTakeable = true
	emerald.Flags.IsTreasure = true
	emerald.Value = 10 // ZIL: 10
	g.Items["emerald"] = emerald

	// CHALICE - treasure and container - ZIL TVALUE 5, LDESC from ZIL
	chalice := NewItem("chalice", "silver chalice", "The chalice is made of silver and intricately engraved.")
	chalice.Aliases = []string{"chalice", "cup", "silver", "treasure", "engravings"}
	chalice.Location = "treasure-room"
	chalice.RoomDescription = "There is a silver chalice, intricately engraved, here."
	chalice.Flags.IsTakeable = true
	chalice.Flags.IsTreasure = true
	chalice.Flags.IsContainer = true
	chalice.Value = 5 // ZIL: 5
	g.Items["chalice"] = chalice

	// JADE FIGURINE - LDESC from ZIL
	jade := NewItem("jade", "jade figurine", "The figurine is carved from exquisite jade.")
	jade.Aliases = []string{"jade", "figurine", "treasure", "exquisite"}
	jade.RoomDescription = "There is an exquisite jade figurine here."
	jade.Flags.IsTakeable = true
	jade.Flags.IsTreasure = true
	jade.Value = 5
	g.Items["jade"] = jade

	// BAG OF COINS - LDESC from ZIL (also known as BAG-OF-COINS in ZIL)
	coins := NewItem("coins", "bag of coins", "The leather bag is old and bulging with coins.")
	coins.Aliases = []string{"coins", "bag", "bag-of-coins", "treasure", "leather"}
	coins.RoomDescription = "An old leather bag, bulging with coins, is here."
	coins.Flags.IsTakeable = true
	coins.Flags.IsTreasure = true
	coins.Value = 5
	g.Items["coins"] = coins
	g.Items["bag-of-coins"] = coins // ZIL uses BAG-OF-COINS

	// PAINTING - ZIL TVALUE 6, FDESC/LDESC from ZIL
	painting := NewItem("painting", "painting", "It's a painting of unparalleled beauty by a neglected genius.")
	painting.Aliases = []string{"painting", "treasure", "art", "canvas"}
	painting.Location = "gallery"
	painting.RoomDescription = "A painting by a neglected genius is here."
	painting.Flags.IsTakeable = true
	painting.Flags.IsTreasure = true
	painting.Flags.IsBurnable = true // BURNBIT in ZIL
	painting.Value = 6 // ZIL: 6
	g.Items["painting"] = painting

	// BRACELET - sapphire-encrusted from ZIL
	bracelet := NewItem("bracelet", "sapphire-encrusted bracelet", "The bracelet is encrusted with sapphires.")
	bracelet.Aliases = []string{"bracelet", "sapphire", "treasure", "jewel"}
	bracelet.RoomDescription = "There is a sapphire-encrusted bracelet here."
	bracelet.Flags.IsTakeable = true
	bracelet.Flags.IsTreasure = true
	bracelet.Value = 5
	g.Items["bracelet"] = bracelet

	// BAUBLE (brass bauble) - ZIL TVALUE 1
	bauble := NewItem("bauble", "brass bauble", "It's a beautiful brass bauble.")
	bauble.Aliases = []string{"bauble", "treasure", "brass"}
	bauble.RoomDescription = "There is a beautiful brass bauble here."
	bauble.Flags.IsTakeable = true
	bauble.Flags.IsTreasure = true
	bauble.Value = 1 // ZIL: 1
	g.Items["bauble"] = bauble

	// SCARAB
	scarab := NewItem("scarab", "beautiful scarab", "The scarab is beautifully carved.")
	scarab.Aliases = []string{"scarab", "treasure"}
	scarab.RoomDescription = "There is a beautiful scarab here."
	scarab.Flags.IsTakeable = true
	scarab.Flags.IsTreasure = true
	scarab.Value = 5
	g.Items["scarab"] = scarab

	// POT OF GOLD - FDESC from ZIL, invisible until rainbow solid
	pot := NewItem("pot-of-gold", "pot of gold", "It's a pot full of gold coins.")
	pot.Aliases = []string{"pot", "gold", "pot-of-gold", "treasure"}
	pot.Location = "end-of-rainbow"
	pot.RoomDescription = "At the end of the rainbow is a pot of gold."
	pot.Flags.IsTakeable = true
	pot.Flags.IsTreasure = true
	pot.Flags.IsInvisible = true // Invisible until rainbow is solidified
	pot.Value = 10
	g.Items["pot-of-gold"] = pot
	g.Rooms["end-of-rainbow"].AddItem("pot-of-gold")

	// TRIDENT - ZIL TVALUE 11, FDESC from ZIL
	trident := NewItem("trident", "crystal trident", "It's Poseidon's own crystal trident, a weapon of great power.")
	trident.Aliases = []string{"trident", "crystal", "treasure", "fork", "poseidon"}
	trident.Location = "falls"
	trident.RoomDescription = "On the shore lies Poseidon's own crystal trident."
	trident.Flags.IsTakeable = true
	trident.Flags.IsTreasure = true
	trident.Flags.IsWeapon = true
	trident.Value = 11 // ZIL: 11
	g.Items["trident"] = trident

	// SCEPTRE - ZIL TVALUE 6
	sceptre := NewItem("sceptre", "sceptre", "The sceptre is encrusted with jewels and appears to be from ancient Egypt.")
	sceptre.Aliases = []string{"sceptre", "scepter", "treasure"}
	sceptre.RoomDescription = "There is a sceptre, probably that of ancient Egypt itself, here."
	sceptre.Flags.IsTakeable = true
	sceptre.Flags.IsTreasure = true
	sceptre.Flags.IsWeapon = true
	sceptre.Value = 6 // ZIL: 6
	g.Items["sceptre"] = sceptre

	// EGG (containing bauble) - treasure, FDESC from ZIL
	egg := NewItem("egg", "jewel-encrusted egg", "The egg is covered with jewels and has a golden clasp. It appears extremely fragile.")
	egg.Aliases = []string{"egg", "treasure", "jeweled", "encrusted", "birds"}
	egg.Location = "nest"
	egg.RoomDescription = "There is a jewel-encrusted egg here."
	egg.Flags.IsTakeable = true
	egg.Flags.IsContainer = true
	egg.Flags.IsTreasure = true
	egg.Value = 5
	g.Items["egg"] = egg

	// PLATINUM BAR - ZIL TVALUE 5, LDESC from ZIL (also known as BAR in ZIL)
	platinumBar := NewItem("platinum-bar", "platinum bar", "It's a large bar of solid platinum.")
	platinumBar.Aliases = []string{"bar", "platinum", "platinum-bar", "treasure", "large"}
	platinumBar.RoomDescription = "On the ground is a large platinum bar."
	platinumBar.Flags.IsTakeable = true
	platinumBar.Flags.IsTreasure = true
	platinumBar.Value = 5 // ZIL: 5
	g.Items["platinum-bar"] = platinumBar
	g.Items["bar"] = platinumBar // ZIL uses BAR

	// SAPPHIRE - 5 points
	sapphire := NewItem("sapphire", "large sapphire", "It's a large, brilliantly cut sapphire.")
	sapphire.Aliases = []string{"sapphire", "gem", "treasure", "large"}
	sapphire.RoomDescription = "There is a large sapphire here."
	sapphire.Flags.IsTakeable = true
	sapphire.Flags.IsTreasure = true
	sapphire.Value = 5
	g.Items["sapphire"] = sapphire

	// IVORY TORCH - 6 points, FDESC from ZIL
	ivoryTorch := NewItem("ivory-torch", "ivory torch", "The torch is made of ivory and burns with an eternal flame.")
	ivoryTorch.Aliases = []string{"ivory-torch", "ivory", "torch", "treasure", "flaming"}
	ivoryTorch.RoomDescription = "Sitting on the pedestal is a flaming torch, made of ivory."
	ivoryTorch.Flags.IsTakeable = true
	ivoryTorch.Flags.IsTreasure = true
	ivoryTorch.Flags.IsLightSource = true
	ivoryTorch.Flags.IsLit = true
	ivoryTorch.Fuel = -1 // Eternal flame, never burns out
	ivoryTorch.Value = 6
	g.Items["ivory-torch"] = ivoryTorch

	// TRUNK OF JEWELS - ZIL TVALUE 5, FDESC/LDESC from ZIL
	trunkOfJewels := NewItem("trunk-of-jewels", "trunk of jewels", "The old trunk is bulging with assorted jewels.")
	trunkOfJewels.Aliases = []string{"trunk", "jewels", "trunk-of-jewels", "treasure", "old"}
	trunkOfJewels.RoomDescription = "There is an old trunk here, bulging with assorted jewels."
	trunkOfJewels.Flags.IsTakeable = true
	trunkOfJewels.Flags.IsTreasure = true
	trunkOfJewels.Flags.IsContainer = true
	trunkOfJewels.Flags.IsInvisible = true // INVISIBLE in ZIL initially
	trunkOfJewels.Value = 5 // ZIL: 5
	g.Items["trunk-of-jewels"] = trunkOfJewels

	// PEARL - 1 point
	pearl := NewItem("pearl", "large pearl", "It's an enormous, lustrous pearl.")
	pearl.Aliases = []string{"pearl", "treasure", "large", "enormous"}
	pearl.RoomDescription = "There is an enormous pearl resting in an open clam here."
	pearl.Flags.IsTakeable = true
	pearl.Flags.IsTreasure = true
	pearl.Value = 1
	g.Items["pearl"] = pearl

	// Mark the oriental rug as a treasure (it already exists in createFixedObjects)
	// We'll update it there to have the treasure flag
}

// createWeapons creates all weapon items
func createWeapons(g *GameV2) {
	// SWORD (elvish) - FDESC from ZIL line 923, SWORD-FCN shows glow when examining
	sword := NewItem("sword", "elvish sword", "The sword is well crafted.")
	sword.Aliases = []string{"sword", "blade", "elvish", "orcrist", "glamdring"}
	sword.RoomDescription = "Above the trophy case hangs an elvish sword of great antiquity."
	sword.Location = "living-room"
	sword.Flags.IsTakeable = true
	sword.Flags.IsWeapon = true
	g.Items["sword"] = sword
	g.Rooms["living-room"].AddItem("sword")

	// KNIFE - FDESC from ZIL, nasty-looking knife on attic table (container item)
	knife := NewItem("knife", "nasty knife", "It's a nasty-looking knife.")
	knife.Aliases = []string{"knife", "knives", "blade", "nasty", "unrusty"}
	knife.Location = "attic-table" // IN attic-table container
	knife.RoomDescription = "On a table is a nasty-looking knife."
	knife.Flags.IsTakeable = true
	knife.Flags.IsWeapon = true
	g.Items["knife"] = knife

	// RUSTY KNIFE - FDESC from ZIL, beside skeleton in maze
	rustyKnife := NewItem("rusty-knife", "rusty knife", "It's an old rusty knife.")
	rustyKnife.Aliases = []string{"rusty-knife", "knife", "rusty", "knives"}
	rustyKnife.RoomDescription = "Beside the skeleton is a rusty knife."
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

	// ROPE - FDESC from ZIL, large coil in attic corner
	rope := NewItem("rope", "rope", "It's a large coil of strong hemp rope.")
	rope.Aliases = []string{"rope", "hemp", "coil", "large"}
	rope.Location = "attic"
	rope.RoomDescription = "A large coil of rope is lying in the corner."
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

	// TROPHY CASE (in living room) - NDESCBIT means it's described in room description
	trophyCase := NewItem("trophy-case", "trophy case", "The trophy case is securely fastened to the wall.")
	trophyCase.Aliases = []string{"case", "trophy-case", "trophy"}
	trophyCase.Location = "living-room"
	trophyCase.Flags.IsContainer = true
	trophyCase.Flags.IsOpen = false
	trophyCase.Flags.IsTransparent = true
	trophyCase.Flags.NoRoomListing = true // NDESCBIT - don't list separately, it's mentioned in room desc
	g.Items["trophy-case"] = trophyCase
	g.Rooms["living-room"].AddItem("trophy-case")

	// BOTTLE (glass bottle) - FDESC from ZIL, on kitchen table
	bottle := NewItem("bottle", "glass bottle", "It's a clear glass bottle that can hold liquids.")
	bottle.Aliases = []string{"bottle", "glass", "clear", "container"}
	bottle.Location = "kitchen-table"
	bottle.RoomDescription = "A bottle is sitting on the table."
	bottle.Flags.IsTakeable = true
	bottle.Flags.IsContainer = true
	bottle.Flags.IsTransparent = true
	g.Items["bottle"] = bottle

	// COFFIN - ZIL TVALUE 15 (treasure!), LDESC from ZIL
	coffin := NewItem("coffin", "gold coffin", "The solid-gold coffin is used for the burial of Ramses II and is intricately decorated.")
	coffin.Aliases = []string{"coffin", "treasure", "casket", "solid", "gold"}
	coffin.Location = "egypt-room"
	coffin.RoomDescription = "The solid-gold coffin used for the burial of Ramses II is here."
	coffin.Flags.IsTakeable = true
	coffin.Flags.IsContainer = true
	coffin.Flags.IsOpen = false
	coffin.Flags.IsTreasure = true
	coffin.Value = 15 // ZIL: 15
	g.Items["coffin"] = coffin
	g.Rooms["egypt-room"].AddItem("coffin")

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

	// NEST (bird's nest containing egg) - FDESC from ZIL
	nest := NewItem("nest", "bird's nest", "It's a small bird's nest.")
	nest.Aliases = []string{"nest", "birds"}
	nest.Location = "up-a-tree"
	nest.RoomDescription = "Beside you on the branch is a small bird's nest."
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
	// LEAFLET (in mailbox) - LDESC from ZIL
	leaflet := NewItem("leaflet", "leaflet", `"WELCOME TO ZORK!

ZORK is a game of adventure, danger, and low cunning. In it you will explore some of the most amazing territory ever seen by mortals. No computer should be without one!"`)
	leaflet.Aliases = []string{"leaflet", "pamphlet", "booklet", "advertisement", "mail", "small"}
	leaflet.Location = "mailbox"
	leaflet.RoomDescription = "A small leaflet is on the ground."
	leaflet.Flags.IsTakeable = true
	leaflet.Flags.IsReadable = true
	g.Items["leaflet"] = leaflet

	// BOOK (prayer book) - ZIL lines 212-231
	book := NewItem("book", "black book", "On the altar is a large black book, open to page 569.")
	book.Aliases = []string{"book", "prayer-book", "prayer", "black", "black-book"}
	book.Flags.IsTakeable = true
	book.Flags.IsReadable = true
	book.Flags.IsContainer = false // Not a container
	book.Flags.IsOpen = true        // Always open to page 569
	book.Text = `Commandment #12592

Oh ye who go about saying unto each: "Hello sailor":
Dost thou know the magnitude of thy sin before the gods?
Yea, verily, thou shalt be ground between two stones.
Shall the angry gods cast thy body into the whirlpool?
Surely, thy eye shall be put out with a sharp stick!
Even unto the ends of the earth shalt thou wander and
Unto the land of the dead shalt thou be sent at last.
Surely thou shalt repent of thy cunning.`
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
	// LAMP (brass lantern) - FDESC/LDESC from ZIL, starts with 330 turns of fuel
	lamp := NewItem("lamp", "brass lantern", "The brass lantern is battery-powered. It is currently on.")
	lamp.Aliases = []string{"lamp", "lantern", "light"}
	lamp.RoomDescription = "A battery-powered brass lantern is on the trophy case."
	lamp.Location = "living-room"
	lamp.Flags.IsTakeable = true
	lamp.Flags.IsLightSource = true
	lamp.Flags.IsLit = true
	lamp.Fuel = 330 // Total turns before lamp dies
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

	// CANDLES - FDESC from ZIL, burn time from I-CANDLES in ZIL line 2641 (40 turns)
	candles := NewItem("candles", "pair of candles", "They are burning candles.")
	candles.Aliases = []string{"candles", "candle", "pair", "burning"}
	candles.Location = "entrance-to-hades"
	candles.RoomDescription = "On the two ends of the altar are burning candles."
	candles.Flags.IsTakeable = true
	candles.Flags.IsLightSource = true
	candles.Flags.IsLit = false
	candles.Fuel = 40 // Burns for 40 turns when lit
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

	// RUG (conceals trap door) - NDESCBIT, mentioned in room description
	rug := NewItem("rug", "oriental rug", "The rug is extremely heavy and cannot be carried.")
	rug.Aliases = []string{"rug", "oriental-rug", "carpet", "treasure"}
	rug.Location = "living-room"
	rug.Flags.IsTakeable = true // Can be taken once trap door is opened
	rug.Flags.IsTreasure = true
	rug.Flags.NoRoomListing = true // NDESCBIT - mentioned in room description
	rug.Value = 15
	g.Items["rug"] = rug
	g.Rooms["living-room"].AddItem("rug")

	// TRAP DOOR (under rug) - global object accessible from both living-room and cellar
	trapDoor := NewItem("trap-door", "trap door", "There is a trap door here.")
	trapDoor.Aliases = []string{"trap-door", "door", "trapdoor", "cover"}
	trapDoor.Location = "GLOBAL" // Special marker for global objects
	trapDoor.Flags.IsTakeable = false
	trapDoor.Flags.IsOpen = false
	g.Items["trap-door"] = trapDoor
	// Add to both rooms so findItem can locate it
	g.Rooms["living-room"].AddItem("trap-door")
	g.Rooms["cellar"].AddItem("trap-door")

	// GRATE (above cellar) - global object accessible from both grating-clearing and grating-room
	grate := NewItem("grate", "grating", "There is a grating securely fastened into the ceiling.")
	grate.Aliases = []string{"grate", "grating"}
	grate.Location = "GLOBAL" // Special marker for global objects
	grate.Flags.IsTakeable = false
	g.Items["grate"] = grate
	// Add to both rooms so findItem can locate it
	g.Rooms["grating-clearing"].AddItem("grate")
	g.Rooms["grating-room"].AddItem("grate")

	// BUTTONS (control panel)
	yellowButton := NewItem("yellow-button", "yellow button", "There is a yellow button here.")
	yellowButton.Aliases = []string{"yellow-button", "yellow button", "yellow"}
	yellowButton.Location = "maintenance-room"
	yellowButton.Flags.IsTakeable = false
	g.Items["yellow-button"] = yellowButton
	g.Rooms["maintenance-room"].AddItem("yellow-button")

	brownButton := NewItem("brown-button", "brown button", "There is a brown button here.")
	brownButton.Aliases = []string{"brown-button", "brown button", "brown"}
	brownButton.Location = "maintenance-room"
	brownButton.Flags.IsTakeable = false
	g.Items["brown-button"] = brownButton
	g.Rooms["maintenance-room"].AddItem("brown-button")

	redButton := NewItem("red-button", "red button", "There is a red button here.")
	redButton.Aliases = []string{"red-button", "red button", "red"}
	redButton.Location = "maintenance-room"
	redButton.Flags.IsTakeable = false
	g.Items["red-button"] = redButton
	g.Rooms["maintenance-room"].AddItem("red-button")

	blueButton := NewItem("blue-button", "blue button", "There is a blue button here.")
	blueButton.Aliases = []string{"blue-button", "blue button", "blue"}
	blueButton.Location = "maintenance-room"
	blueButton.Flags.IsTakeable = false
	g.Items["blue-button"] = blueButton
	g.Rooms["maintenance-room"].AddItem("blue-button")

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
	g.Rooms["machine-room"].AddItem("machine")

	// Machine control buttons
	startButton := NewItem("start-button", "start button", "")
	startButton.Aliases = []string{"start-button", "start button", "start"}
	startButton.Location = "machine-room"
	startButton.Flags.IsTakeable = false
	g.Items["start-button"] = startButton
	g.Rooms["machine-room"].AddItem("start-button")

	launchButton := NewItem("launch-button", "launch button", "")
	launchButton.Aliases = []string{"launch-button", "launch button", "launch"}
	launchButton.Location = "machine-room"
	launchButton.Flags.IsTakeable = false
	g.Items["launch-button"] = launchButton
	g.Rooms["machine-room"].AddItem("launch-button")

	lowerButton := NewItem("lower-button", "lower button", "")
	lowerButton.Aliases = []string{"lower-button", "lower button", "lower"}
	lowerButton.Location = "machine-room"
	lowerButton.Flags.IsTakeable = false
	g.Items["lower-button"] = lowerButton
	g.Rooms["machine-room"].AddItem("lower-button")

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

	// MOUNTAIN RANGE (also MOUNTAIN-RANGE in ZIL)
	mountains := NewItem("mountains", "mountain range", "The mountains are impassable.")
	mountains.Aliases = []string{"mountains", "mountain", "range", "mountain-range"}
	mountains.Flags.IsTakeable = false
	g.Items["mountains"] = mountains
	g.Items["mountain-range"] = mountains // ZIL uses MOUNTAIN-RANGE

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

	// CHIMNEY - LOCAL-GLOBALS object (accessible from kitchen and living room)
	chimney := NewItem("chimney", "chimney", "The chimney is dark and narrow. It leads downward into darkness.")
	chimney.Aliases = []string{"chimney", "dark", "narrow"}
	chimney.Location = "GLOBAL" // Global object
	chimney.Flags.IsTakeable = false
	chimney.Flags.NoRoomListing = true // NDESCBIT - mentioned in room descriptions
	g.Items["chimney"] = chimney
	// Add to both kitchen and living-room
	g.Rooms["kitchen"].AddItem("chimney")
	g.Rooms["living-room"].AddItem("chimney")

	// KITCHEN TABLE - NDESCBIT, mentioned in room description
	kitchenTable := NewItem("kitchen-table", "kitchen table", "It's an ordinary kitchen table.")
	kitchenTable.Aliases = []string{"table", "kitchen-table", "kitchen"}
	kitchenTable.Location = "kitchen"
	kitchenTable.Flags.IsTakeable = false
	kitchenTable.Flags.IsContainer = true
	kitchenTable.Flags.IsOpen = true
	kitchenTable.Flags.NoRoomListing = true // NDESCBIT - mentioned in room description
	g.Items["kitchen-table"] = kitchenTable
	g.Rooms["kitchen"].AddItem("kitchen-table")

	// ATTIC TABLE - NDESCBIT
	atticTable := NewItem("attic-table", "table", "It's a small wooden table.")
	atticTable.Aliases = []string{"table", "attic-table"}
	atticTable.Location = "attic"
	atticTable.Flags.IsTakeable = false
	atticTable.Flags.IsContainer = true
	atticTable.Flags.IsOpen = true
	atticTable.Flags.NoRoomListing = true // NDESCBIT
	g.Items["attic-table"] = atticTable
	g.Rooms["attic"].AddItem("attic-table")

	// BOARD - LOCAL-GLOBALS scenery object
	board := NewItem("board", "board", "The boards appear to be nailed securely across the windows.")
	board.Aliases = []string{"boards", "board"}
	board.Location = "GLOBAL"
	board.Flags.IsTakeable = false
	board.Flags.NoRoomListing = true // NDESCBIT
	g.Items["board"] = board

	// TREE - LOCAL-GLOBALS scenery object
	tree := NewItem("tree", "tree", "The trees are large and storm-tossed.")
	tree.Aliases = []string{"tree", "branch", "trees", "large", "storm"}
	tree.Location = "GLOBAL"
	tree.Flags.IsTakeable = false
	tree.Flags.NoRoomListing = true // NDESCBIT
	g.Items["tree"] = tree

	// CRACK - LOCAL-GLOBALS scenery object
	crack := NewItem("crack", "crack", "It's a narrow crack in the wall.")
	crack.Aliases = []string{"crack", "narrow"}
	crack.Location = "GLOBAL"
	crack.Flags.IsTakeable = false
	crack.Flags.NoRoomListing = true // NDESCBIT
	g.Items["crack"] = crack

	// CLIMBABLE-CLIFF - LOCAL-GLOBALS scenery object
	climbableCliff := NewItem("climbable-cliff", "cliff", "The rocky cliff face is steep but appears climbable.")
	climbableCliff.Aliases = []string{"wall", "cliff", "walls", "ledge", "rocky", "sheer"}
	climbableCliff.Location = "GLOBAL"
	climbableCliff.Flags.IsTakeable = false
	climbableCliff.Flags.NoRoomListing = true // NDESCBIT
	g.Items["climbable-cliff"] = climbableCliff

	// WHITE-CLIFF - LOCAL-GLOBALS scenery object
	whiteCliff := NewItem("white-cliff", "white cliffs", "The White Cliffs of Quendor are massive ramparts of white stone.")
	whiteCliff.Aliases = []string{"cliff", "cliffs", "white"}
	whiteCliff.Location = "GLOBAL"
	whiteCliff.Flags.IsTakeable = false
	whiteCliff.Flags.NoRoomListing = true // NDESCBIT
	g.Items["white-cliff"] = whiteCliff

	// WALL - GLOBAL-OBJECTS scenery
	wall := NewItem("wall", "surrounding wall", "The walls surround you on all sides.")
	wall.Aliases = []string{"wall", "walls", "surrounding"}
	wall.Location = "GLOBAL"
	wall.Flags.IsTakeable = false
	wall.Flags.NoRoomListing = true
	g.Items["wall"] = wall

	// SONGBIRD - LOCAL-GLOBALS scenery object
	songbird := NewItem("songbird", "songbird", "It's a beautiful songbird, singing merrily.")
	songbird.Aliases = []string{"bird", "songbird", "song"}
	songbird.Location = "GLOBAL"
	songbird.Flags.IsTakeable = false
	songbird.Flags.NoRoomListing = true // NDESCBIT
	g.Items["songbird"] = songbird

	// BOARDED-WINDOW - LOCAL-GLOBALS scenery object
	boardedWindow := NewItem("boarded-window", "boarded window", "The window is boarded up securely.")
	boardedWindow.Aliases = []string{"window", "boarded"}
	boardedWindow.Location = "GLOBAL"
	boardedWindow.Flags.IsTakeable = false
	boardedWindow.Flags.NoRoomListing = true // NDESCBIT
	g.Items["boarded-window"] = boardedWindow

	// FRONT-DOOR - At west-of-house
	frontDoor := NewItem("front-door", "door", "The front door is boarded and can't be opened.")
	frontDoor.Aliases = []string{"door", "front", "boarded", "front-door"}
	frontDoor.Location = "west-of-house"
	frontDoor.Flags.IsTakeable = false
	frontDoor.Flags.NoRoomListing = true // NDESCBIT
	g.Items["front-door"] = frontDoor
	g.Rooms["west-of-house"].AddItem("front-door")

	// BARROW-DOOR - At stone-barrow
	barrowDoor := NewItem("barrow-door", "stone door", "It's a huge stone door.")
	barrowDoor.Aliases = []string{"door", "huge", "stone", "barrow-door"}
	barrowDoor.Flags.IsTakeable = false
	barrowDoor.Flags.IsOpen = true
	barrowDoor.Flags.NoRoomListing = true // NDESCBIT
	g.Items["barrow-door"] = barrowDoor

	// BARROW - At stone-barrow
	barrow := NewItem("barrow", "stone barrow", "The barrow is a massive stone structure.")
	barrow.Aliases = []string{"barrow", "tomb", "massive", "stone"}
	barrow.Flags.IsTakeable = false
	barrow.Flags.NoRoomListing = true // NDESCBIT
	g.Items["barrow"] = barrow

	// Additional scenery objects from ZIL

	// TEETH - GLOBAL-OBJECTS
	teeth := NewItem("teeth", "set of teeth", "They're a set of sharp, menacing teeth.")
	teeth.Aliases = []string{"teeth", "overboard"}
	teeth.Location = "GLOBAL"
	teeth.Flags.IsTakeable = false
	teeth.Flags.NoRoomListing = true
	g.Items["teeth"] = teeth

	// GLOBAL-WATER - LOCAL-GLOBALS water (different from bottle water)
	globalWater := NewItem("global-water", "water", "The water is flowing and appears drinkable.")
	globalWater.Aliases = []string{"water", "stream"}
	globalWater.Location = "GLOBAL"
	globalWater.Flags.IsTakeable = false
	globalWater.Flags.NoRoomListing = true
	g.Items["global-water"] = globalWater

	// BODIES - LOCAL-GLOBALS pile of bodies
	bodies := NewItem("bodies", "pile of bodies", "It's a gruesome pile of bodies.")
	bodies.Aliases = []string{"bodies", "pile", "corpses"}
	bodies.Location = "GLOBAL"
	bodies.Flags.IsTakeable = false
	bodies.Flags.NoRoomListing = true
	g.Items["bodies"] = bodies

	// DAM - DAM-ROOM scenery
	dam := NewItem("dam", "dam", "The dam is a massive structure controlling the flow of water.")
	dam.Aliases = []string{"dam", "structure"}
	dam.Flags.IsTakeable = false
	dam.Flags.NoRoomListing = true
	g.Items["dam"] = dam

	// LEAK - MAINTENANCE-ROOM scenery
	leak := NewItem("leak", "leak", "There's a small leak dripping water.")
	leak.Aliases = []string{"leak", "drip"}
	leak.Flags.IsTakeable = false
	leak.Flags.NoRoomListing = true
	g.Items["leak"] = leak

	// CONTROL-PANEL - DAM-ROOM scenery
	controlPanel := NewItem("control-panel", "control panel", "The control panel has various buttons and switches.")
	controlPanel.Aliases = []string{"control-panel", "panel", "controls"}
	controlPanel.Flags.IsTakeable = false
	controlPanel.Flags.NoRoomListing = true
	g.Items["control-panel"] = controlPanel

	// MACHINE-SWITCH - MACHINE-ROOM
	machineSwitch := NewItem("machine-switch", "switch", "It's a large switch on the machine.")
	machineSwitch.Aliases = []string{"switch", "machine-switch"}
	machineSwitch.Flags.IsTakeable = false
	machineSwitch.Flags.NoRoomListing = true
	g.Items["machine-switch"] = machineSwitch

	// BOLT - DAM-ROOM
	bolt := NewItem("bolt", "bolt", "It's a large metal bolt.")
	bolt.Aliases = []string{"bolt"}
	bolt.Flags.IsTakeable = true
	g.Items["bolt"] = bolt

	// BUBBLE - DAM-ROOM
	bubble := NewItem("bubble", "green bubble", "It's a strange green bubble.")
	bubble.Aliases = []string{"bubble", "green"}
	bubble.Flags.IsTakeable = false
	g.Items["bubble"] = bubble

	// BROKEN-LAMP
	brokenLamp := NewItem("broken-lamp", "broken lantern", "The lantern is broken and useless.")
	brokenLamp.Aliases = []string{"broken-lamp", "broken", "lantern"}
	brokenLamp.Flags.IsTakeable = true
	g.Items["broken-lamp"] = brokenLamp

	// HOT-BELL
	hotBell := NewItem("hot-bell", "red hot brass bell", "The bell is glowing red hot. Don't touch it!")
	hotBell.Aliases = []string{"hot-bell", "bell", "red", "hot", "brass"}
	hotBell.Flags.IsTakeable = false
	g.Items["hot-bell"] = hotBell

	// GUNK
	gunk := NewItem("gunk", "small piece of vitreous slag", "It's a small, glassy piece of slag.")
	gunk.Aliases = []string{"gunk", "slag", "vitreous"}
	gunk.Flags.IsTakeable = true
	g.Items["gunk"] = gunk
}

// createMiscItems creates various other items
func createMiscItems(g *GameV2) {
	// INFLATABLE BOAT (also INFLATABLE-BOAT in ZIL)
	boat := NewItem("boat", "inflatable boat", "There is an inflatable boat here.")
	boat.Aliases = []string{"boat", "inflatable-boat", "raft", "pile", "plastic"}
	boat.Flags.IsTakeable = true
	g.Items["boat"] = boat
	g.Items["inflatable-boat"] = boat // ZIL uses INFLATABLE-BOAT

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

	// SKULL - ZIL TVALUE 10 (treasure!) - crystal skull from LAND-OF-LIVING-DEAD
	skull := NewItem("skull", "crystal skull", "The crystal skull is beautifully carved and grinning rather nastily.")
	skull.Aliases = []string{"skull", "head", "treasure", "crystal"}
	skull.RoomDescription = "Lying in one corner of the room is a beautifully carved crystal skull. It appears to be grinning at you rather nastily."
	skull.Flags.IsTakeable = true
	skull.Flags.IsTreasure = true
	skull.Value = 10 // ZIL: 10
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

	// RAISED BASKET (starts at shaft-room, raised to ceiling) - LDESC from ZIL
	raisedBasket := NewItem("raised-basket", "basket", "It's a wicker basket suspended from a chain.")
	raisedBasket.Aliases = []string{"basket", "wicker-basket", "raised-basket", "cage", "dumbwaiter"}
	raisedBasket.Location = "shaft-room"
	raisedBasket.RoomDescription = "At the end of the chain is a basket."
	raisedBasket.Flags.IsTakeable = false
	raisedBasket.Flags.IsContainer = true
	raisedBasket.Flags.IsOpen = true
	raisedBasket.Flags.IsTransparent = true
	g.Items["raised-basket"] = raisedBasket
	g.Rooms["shaft-room"].AddItem("raised-basket")

	// LOWERED BASKET (appears when basket is lowered) - LDESC from ZIL
	loweredBasket := NewItem("lowered-basket", "basket", "It's a wicker basket suspended from a chain.")
	loweredBasket.Aliases = []string{"basket", "wicker-basket", "lowered-basket", "cage", "dumbwaiter"}
	loweredBasket.RoomDescription = "From the chain is suspended a basket."
	loweredBasket.Flags.IsTakeable = false
	loweredBasket.Flags.IsContainer = true
	loweredBasket.Flags.IsOpen = true
	loweredBasket.Flags.IsTransparent = true
	g.Items["lowered-basket"] = loweredBasket

	// PRAYER (of protection)
	prayer := NewItem("prayer", "prayer", "The prayer seems to be a plea for protection.")
	prayer.Aliases = []string{"prayer"}
	prayer.Flags.IsTakeable = false
	prayer.Flags.IsReadable = true
	g.Items["prayer"] = prayer

	// KEYS
	keys := NewItem("keys", "set of keys", "It's just a normal set of keys.")
	keys.Aliases = []string{"keys", "key"}
	keys.RoomDescription = "There is a set of keys here."
	keys.Flags.IsTakeable = true
	keys.Location = "living-room"
	g.Items["keys"] = keys
	g.Rooms["living-room"].AddItem("keys")

	// CLAM (contains pearl)
	clam := NewItem("clam", "giant clam", "There is a giant clam here.")
	clam.Aliases = []string{"clam", "shell"}
	clam.Flags.IsContainer = true
	clam.Flags.IsOpen = true
	g.Items["clam"] = clam

	// MATCHES (book of matches)
	matches := NewItem("matches", "book of matches", "There is a book of matches here.")
	matches.Aliases = []string{"matches", "book-of-matches", "matchbook"}
	matches.Flags.IsTakeable = true
	matches.Flags.IsReadable = true
	g.Items["matches"] = matches

	// MIRROR (in mirror-room-1)
	mirror1 := NewItem("mirror-1", "mirror", "An enormous mirror fills the south wall.")
	mirror1.Aliases = []string{"mirror", "looking-glass"}
	mirror1.Location = "mirror-room-1"
	mirror1.Flags.IsTakeable = false
	g.Items["mirror-1"] = mirror1
	g.Rooms["mirror-room-1"].AddItem("mirror-1")

	// MIRROR (in mirror-room-2)
	mirror2 := NewItem("mirror-2", "mirror", "An enormous mirror fills the north wall.")
	mirror2.Aliases = []string{"mirror", "looking-glass"}
	mirror2.Location = "mirror-room-2"
	mirror2.Flags.IsTakeable = false
	g.Items["mirror-2"] = mirror2
	g.Rooms["mirror-room-2"].AddItem("mirror-2")

	// MIRROR (generic - for other uses)
	mirror := NewItem("mirror", "mirror", "There is a large mirror here.")
	mirror.Aliases = []string{"mirror", "looking-glass"}
	mirror.Flags.IsTakeable = false
	g.Items["mirror"] = mirror

	// PILE (pile of leaves)
	pile := NewItem("pile-of-leaves", "pile of leaves", "There is a pile of leaves here.")
	pile.Aliases = []string{"pile", "leaves", "pile-of-leaves"}
	pile.Flags.IsTakeable = false
	g.Items["pile-of-leaves"] = pile

	// Note: Grate is a global object defined earlier (line ~501) and accessible from both grating-clearing and grating-room

	// CYCLOPS (as an object, not NPC - for examine purposes)
	cyclopsCorp := NewItem("cyclops-corpse", "cyclops corpse", "The body of a dead cyclops is here.")
	cyclopsCorp.Aliases = []string{"corpse", "body", "cyclops"}
	cyclopsCorp.Flags.IsTakeable = false
	g.Items["cyclops-corpse"] = cyclopsCorp

	// THIEF (as corpse/object)
	thiefCorpse := NewItem("thief-corpse", "thief corpse", "The body of a dead thief is here.")
	thiefCorpse.Aliases = []string{"corpse", "body", "thief"}
	thiefCorpse.Flags.IsTakeable = false
	g.Items["thief-corpse"] = thiefCorpse

	// RESERVOIR (large body of water)
	reservoir := NewItem("reservoir", "reservoir", "The reservoir is a large body of water.")
	reservoir.Aliases = []string{"reservoir", "lake"}
	reservoir.Flags.IsTakeable = false
	g.Items["reservoir"] = reservoir

	// STREAM
	stream := NewItem("stream", "stream", "A stream of water flows here.")
	stream.Aliases = []string{"stream", "brook"}
	stream.Flags.IsTakeable = false
	g.Items["stream"] = stream

	// GLACIER
	glacier := NewItem("glacier", "glacier", "A massive glacier fills the cavern.")
	glacier.Aliases = []string{"glacier", "ice"}
	glacier.Flags.IsTakeable = false
	g.Items["glacier"] = glacier

	// SLIDE (ice slide)
	slide := NewItem("slide", "ice slide", "There is a long ice slide here.")
	slide.Aliases = []string{"slide"}
	slide.Flags.IsTakeable = false
	g.Items["slide"] = slide

	// BRICK
	brick := NewItem("brick", "brick", "There is a brick here.")
	brick.Aliases = []string{"brick"}
	brick.Flags.IsTakeable = true
	g.Items["brick"] = brick

	// STATUE (ivory and jade)
	statue := NewItem("statue", "ivory and jade statue", "There is an exquisite statue here.")
	statue.Aliases = []string{"statue", "idol"}
	statue.Flags.IsTakeable = false
	g.Items["statue"] = statue

	// AIR PUMP (for boat)
	airPump := NewItem("air-pump", "air pump", "There is a hand-held air pump here.")
	airPump.Aliases = []string{"air-pump", "pump"}
	airPump.Flags.IsTakeable = true
	g.Items["air-pump"] = airPump

	// CYCLOPS CORPSE TREASURE (what cyclops drops)
	cyclopsTreasure := NewItem("cyclops-treasure", "treasure chest", "There is a small treasure chest here.")
	cyclopsTreasure.Aliases = []string{"chest", "treasure-chest"}
	cyclopsTreasure.Flags.IsTakeable = true
	cyclopsTreasure.Flags.IsContainer = true
	g.Items["cyclops-treasure"] = cyclopsTreasure

	// CRYSTAL (crystal sphere/ball)
	crystal := NewItem("crystal-sphere", "crystal sphere", "There is a crystal sphere here.")
	crystal.Aliases = []string{"crystal", "sphere", "ball", "crystal-ball"}
	crystal.Flags.IsTakeable = true
	g.Items["crystal-sphere"] = crystal

	// RAINBOW (when active)
	rainbowObj := NewItem("rainbow-arc", "rainbow", "A brilliant rainbow arches overhead.")
	rainbowObj.Aliases = []string{"rainbow", "arc"}
	rainbowObj.Flags.IsTakeable = false
	g.Items["rainbow-arc"] = rainbowObj

	// HEADS (shrunken heads)
	heads := NewItem("shrunken-heads", "shrunken heads", "There are four shrunken heads here.")
	heads.Aliases = []string{"heads", "head", "shrunken-heads"}
	heads.Flags.IsTakeable = true
	g.Items["shrunken-heads"] = heads

	// FLASK (crystal flask)
	flask := NewItem("flask", "crystal flask", "There is a crystal flask here.")
	flask.Aliases = []string{"flask", "vial"}
	flask.Flags.IsTakeable = true
	flask.Flags.IsContainer = true
	g.Items["flask"] = flask

	// SWORD-HOLDER (for mounting sword)
	swordHolder := NewItem("sword-holder", "sword holder", "There is a sword holder mounted on the wall.")
	swordHolder.Aliases = []string{"holder", "mount", "sword-holder"}
	swordHolder.Flags.IsTakeable = false
	g.Items["sword-holder"] = swordHolder

	// CRYPT (stone crypt)
	crypt := NewItem("crypt", "stone crypt", "There is a stone crypt here.")
	crypt.Aliases = []string{"crypt", "tomb"}
	crypt.Flags.IsTakeable = false
	crypt.Flags.IsContainer = true
	g.Items["crypt"] = crypt

	// GRANITE WALL
	graniteWall := NewItem("granite-wall", "granite wall", "A massive granite wall blocks your way.")
	graniteWall.Aliases = []string{"wall", "granite", "granite-wall"}
	graniteWall.Flags.IsTakeable = false
	g.Items["granite-wall"] = graniteWall

	// WOODEN DOOR with gothic lettering (in living room) - NDESCBIT, mentioned in room description
	woodenDoor := NewItem("door", "wooden door with strange gothic lettering", "The engravings translate to \"This space intentionally left blank.\"")
	woodenDoor.Aliases = []string{"wooden-door", "door", "lettering", "writing", "front-door", "entrance", "gothic-door"}
	woodenDoor.Location = "living-room"
	woodenDoor.Flags.IsTakeable = false
	woodenDoor.Flags.IsReadable = true
	woodenDoor.Flags.NoRoomListing = true // NDESCBIT - mentioned in room description
	woodenDoor.Text = "The engravings translate to \"This space intentionally left blank.\""
	g.Items["door"] = woodenDoor
	g.Items["wooden-door"] = woodenDoor // ZIL uses WOODEN-DOOR
	g.Rooms["living-room"].AddItem("door")

	// IRON DOOR
	ironDoor := NewItem("iron-door", "iron door", "There is an iron door here.")
	ironDoor.Aliases = []string{"iron-door"}
	ironDoor.Flags.IsTakeable = false
	g.Items["iron-door"] = ironDoor

	// CANDLE (single candle)
	candle := NewItem("candle", "candle", "There is a candle here.")
	candle.Aliases = []string{"candle"}
	candle.Flags.IsTakeable = true
	candle.Flags.IsLightSource = true
	candle.Flags.IsLit = false
	g.Items["candle"] = candle

	// CHAIN (rusty chain)
	chain := NewItem("chain", "rusty chain", "There is a rusty chain here.")
	chain.Aliases = []string{"chain"}
	chain.Flags.IsTakeable = true
	g.Items["chain"] = chain

	// HOOK (brass hook)
	hook := NewItem("hook", "brass hook", "There is a brass hook here.")
	hook.Aliases = []string{"hook"}
	hook.Flags.IsTakeable = true
	g.Items["hook"] = hook

	// PILLAR (marble pillar)
	pillar := NewItem("pillar", "marble pillar", "A massive marble pillar dominates the room.")
	pillar.Aliases = []string{"pillar", "column"}
	pillar.Flags.IsTakeable = false
	g.Items["pillar"] = pillar

	// ALTAR-CLOTH (on altar)
	altarCloth := NewItem("altar-cloth", "altar cloth", "There is an altar cloth here.")
	altarCloth.Aliases = []string{"cloth", "altar-cloth"}
	altarCloth.Flags.IsTakeable = true
	g.Items["altar-cloth"] = altarCloth

	// STICK (walking stick)
	stick := NewItem("stick", "walking stick", "There is a walking stick here.")
	stick.Aliases = []string{"stick", "walking-stick"}
	stick.Flags.IsTakeable = true
	g.Items["stick"] = stick

	// VOLCANO (scenery - for volcano room)
	volcano := NewItem("volcano", "volcano", "A massive volcano looms above.")
	volcano.Aliases = []string{"volcano"}
	volcano.Flags.IsTakeable = false
	g.Items["volcano"] = volcano

	// RAILING (for ledge rooms)
	railing := NewItem("railing", "railing", "There is a wooden railing here.")
	railing.Aliases = []string{"railing", "rail"}
	railing.Flags.IsTakeable = false
	g.Items["railing"] = railing
}
