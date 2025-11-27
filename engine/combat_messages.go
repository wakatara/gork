package engine

// Combat flavor text messages ported from ZIL source
// From 1actions.zil:3611-3788

// Hero attack messages (HERO-MELEE table, 1actions.zil:3611-3648)
var heroMelee = map[int][]string{
	CombatMissed: {
		"Your {weapon} misses the {npc} by an inch.",
		"A good slash, but it misses the {npc} by a mile.",
		"You charge, but the {npc} jumps nimbly aside.",
		"Clang! Crash! The {npc} parries.",
		"A quick stroke, but the {npc} is on guard.",
		"A good stroke, but it's too slow; the {npc} dodges.",
	},
	CombatUnconscious: {
		"Your {weapon} crashes down, knocking the {npc} into dreamland.",
		"The {npc} is battered into unconsciousness.",
		"A furious exchange, and the {npc} is knocked out!",
		"The haft of your {weapon} knocks out the {npc}.",
		"The {npc} is knocked out!",
	},
	CombatKilled: {
		"It's curtains for the {npc} as your {weapon} removes his head.",
		"The fatal blow strikes the {npc} square in the heart: He dies.",
		"The {npc} takes a fatal blow and slumps to the floor dead.",
	},
	CombatLightWound: {
		"The {npc} is struck on the arm; blood begins to trickle down.",
		"Your {weapon} pinks the {npc} on the wrist, but it's not serious.",
		"Your stroke lands, but it was only the flat of the blade.",
		"The blow lands, making a shallow gash in the {npc}'s arm!",
	},
	CombatSeriousWound: {
		"The {npc} receives a deep gash in his side.",
		"A savage blow on the thigh! The {npc} is stunned but can still fight!",
		"Slash! Your blow lands! That one hit an artery, it could be serious!",
		"Slash! Your stroke connects! This could be serious!",
	},
	CombatStagger: {
		"The {npc} is staggered, and drops to his knees.",
		"The {npc} is momentarily disoriented and can't fight back.",
		"The force of your blow knocks the {npc} back, stunned.",
		"The {npc} is confused and can't fight back.",
		"The quickness of your thrust knocks the {npc} back, stunned.",
	},
	CombatLoseWeapon: {
		"The {npc}'s weapon is knocked to the floor, leaving him unarmed.",
		"The {npc} is disarmed by a subtle feint past his guard.",
	},
}

// Troll counter-attack messages (TROLL-MELEE table, 1actions.zil:3689-3729)
var trollMelee = map[int][]string{
	CombatMissed: {
		"The troll swings his axe, but it misses.",
		"The troll's axe barely misses your ear.",
		"The axe sweeps past as you jump aside.",
		"The axe crashes against the rock, throwing sparks!",
	},
	CombatUnconscious: {
		"The flat of the troll's axe hits you delicately on the head, knocking you out.",
	},
	CombatKilled: {
		"The troll neatly removes your head.",
		"The troll's axe stroke cleaves you from the nave to the chops.",
		"The troll's axe removes your head.",
	},
	CombatLightWound: {
		"The axe gets you right in the side. Ouch!",
		"The flat of the troll's axe skins across your forearm.",
		"The troll's swing almost knocks you over as you barely parry in time.",
		"The troll swings his axe, and it nicks your arm as you dodge.",
	},
	CombatSeriousWound: {
		"The troll charges, and his axe slashes you on your {weapon} arm.",
		"An axe stroke makes a deep wound in your leg.",
		"The troll's axe swings down, gashing your shoulder.",
	},
	CombatStagger: {
		"The troll hits you with a glancing blow, and you are momentarily stunned.",
		"The troll swings; the blade turns on your armor but crashes broadside into your head.",
		"You stagger back under a hail of axe strokes.",
		"The troll's mighty blow drops you to your knees.",
	},
	CombatLoseWeapon: {
		"The axe hits your {weapon} and knocks it spinning.",
		"The troll swings, you parry, but the force of his blow knocks your {weapon} away.",
		"The axe knocks your {weapon} out of your hand. It falls to the floor.",
	},
}

// Thief counter-attack messages (THIEF-MELEE table, 1actions.zil:3735-3788)
var thiefMelee = map[int][]string{
	CombatMissed: {
		"The thief stabs nonchalantly with his stiletto and misses.",
		"You dodge as the thief comes in low.",
		"You parry a lightning thrust, and the thief salutes you with a grim nod.",
		"The thief tries to sneak past your guard, but you twist away.",
	},
	CombatUnconscious: {
		"Shifting in the midst of a thrust, the thief knocks you unconscious with the haft of his stiletto.",
		"The thief knocks you out.",
	},
	CombatKilled: {
		"Finishing you off, the thief inserts his blade into your heart.",
		"The thief comes in from the side, feints, and inserts the blade into your ribs.",
		"The thief bows formally, raises his stiletto, and with a wry grin, ends the battle and your life.",
	},
	CombatLightWound: {
		"A quick thrust pinks your left arm, and blood starts to trickle down.",
		"The thief draws blood, raking his stiletto across your arm.",
		"The stiletto flashes faster than you can follow, and blood wells from your leg.",
		"The thief slowly approaches, strikes like a snake, and leaves you wounded.",
	},
	CombatSeriousWound: {
		"The thief strikes like a snake! The resulting wound is serious.",
		"The thief stabs a deep cut in your upper arm.",
		"The stiletto touches your forehead, and the blood obscures your vision.",
		"The thief strikes at your wrist, and suddenly your grip is slippery with blood.",
	},
	CombatStagger: {
		"The butt of his stiletto cracks you on the skull, and you stagger back.",
		"The thief rams the haft of his blade into your stomach, leaving you out of breath.",
		"The thief attacks, and you fall back desperately.",
	},
	CombatLoseWeapon: {
		"A long, theatrical slash. You catch it on your {weapon}, but the thief twists his knife, and the {weapon} goes flying.",
		"The thief neatly flips your {weapon} out of your hands, and it drops to the floor.",
		"You parry a low thrust, and your {weapon} slips out of your hand.",
	},
}

// Cyclops counter-attack messages (CYCLOPS-MELEE table, 1actions.zil:3654-3683)
var cyclopsMelee = map[int][]string{
	CombatMissed: {
		"The Cyclops misses, but the backwash almost knocks you over.",
		"The Cyclops rushes you, but runs into the wall.",
	},
	CombatUnconscious: {
		"The Cyclops sends you crashing to the floor, unconscious.",
	},
	CombatKilled: {
		"The Cyclops breaks your neck with a massive smash.",
	},
	CombatLightWound: {
		"A quick punch, but it was only a glancing blow.",
		"A glancing blow from the Cyclops' fist.",
	},
	CombatSeriousWound: {
		"The monster smashes his huge fist into your chest, breaking several ribs.",
		"The Cyclops almost knocks the wind out of you with a quick punch.",
	},
	CombatStagger: {
		"The Cyclops lands a punch that knocks the wind out of you.",
		"Heedless of your weapons, the Cyclops tosses you against the rock wall of the room.",
	},
	CombatLoseWeapon: {
		"The Cyclops grabs your {weapon}, tastes it, and throws it to the ground in disgust.",
		"The monster grabs you on the wrist, squeezes, and you drop your {weapon} in pain.",
	},
}
