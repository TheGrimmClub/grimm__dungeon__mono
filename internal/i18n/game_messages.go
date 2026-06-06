package i18n

// Engine-facing narrative keys (room descriptions come from YAML content; these
// are the verb-feedback messages). Registered into the shared catalog via init.
const (
	KeyExits          = "game.exits"
	KeyItemsHere      = "game.items_here"
	KeyNoExit         = "game.no_exit"
	KeyWhichDirection = "game.which_direction"
	KeyUnknownDir     = "game.unknown_direction"
	KeyTakeWhat       = "game.take_what"
	KeyTaken          = "game.taken"
	KeyCannotTake     = "game.cannot_take"
	KeyNotHere        = "game.not_here"
	KeyInventoryEmpty = "game.inventory_empty"
	KeyInventoryHead  = "game.inventory_header"
	KeyExamineWhat    = "game.examine_what"
	KeyDontSee        = "game.dont_see"
	KeyUnknownVerb    = "game.unknown_verb"
	KeyWearWhat       = "game.wear_what"
	KeyWorn           = "game.worn"
	KeyCannotWear     = "game.cannot_wear"
	KeyAlreadyWorn    = "game.already_worn"
	KeyWornTag        = "game.worn_tag"
	KeyHeadlampOn     = "game.headlamp_on"
	KeySaved          = "game.saved"
	KeySaveFailed     = "game.save_failed"
	KeySaveDisabled   = "game.save_disabled"
	KeyContinued      = "game.continued"
	KeyVerbHint       = "game.verb_hint"
	KeyEmptyInfo      = "game.empty_info"

	// Help: section headers and one-line descriptions of the game verbs.
	KeyHelpCmdHeader  = "help.cmd_header"
	KeyHelpVerbHeader = "help.verb_header"
	KeyVLook          = "verb.look"
	KeyVGo            = "verb.go"
	KeyVTake          = "verb.take"
	KeyVInspect       = "verb.inspect"
	KeyVWear          = "verb.wear"
	KeyVInventory     = "verb.inventory"

	KeyCmdSave = "cmd.save"
)

func init() {
	add := func(key string, de, en string) {
		catalog[key] = map[Lang]string{DE: de, EN: en}
	}

	add(KeyExits,
		"Ausgänge: %s",
		"Exits: %s")
	add(KeyItemsHere,
		"Hier liegt:",
		"Here lies:")
	add(KeyNoExit,
		"Dorthin führt kein Weg.",
		"There is no path that way.")
	add(KeyWhichDirection,
		"Wohin möchtest du gehen?",
		"Which way do you want to go?")
	add(KeyUnknownDir,
		"»%s« ist keine Richtung, die ich kenne.",
		"\"%s\" is not a direction I know.")
	add(KeyTakeWhat,
		"Was möchtest du nehmen?",
		"What do you want to take?")
	add(KeyTaken,
		"Du nimmst: %s.",
		"You take: %s.")
	add(KeyCannotTake,
		"%s lässt sich nicht nehmen.",
		"%s cannot be taken.")
	add(KeyNotHere,
		"So etwas liegt hier nicht.",
		"There is no such thing here.")
	add(KeyInventoryEmpty,
		"Deine Taschen sind leer.",
		"Your pockets are empty.")
	add(KeyInventoryHead,
		"Du trägst bei dir:",
		"You are carrying:")
	add(KeyExamineWhat,
		"Was möchtest du untersuchen?",
		"What do you want to examine?")
	add(KeyDontSee,
		"So etwas siehst du hier nicht.",
		"You don't see such a thing here.")
	add(KeyUnknownVerb,
		"Das verstehe ich nicht. Versuche: look, go <Richtung>, take <Ding>, inspect <Ding>, wear <Ding>, inventory.",
		"I don't understand. Try: look, go <direction>, take <thing>, inspect <thing>, wear <thing>, inventory.")
	add(KeyWearWhat,
		"Was möchtest du anlegen?",
		"What do you want to wear?")
	add(KeyWorn,
		"Du legst an: %s.",
		"You put on: %s.")
	add(KeyCannotWear,
		"%s lässt sich nicht anlegen.",
		"%s cannot be worn.")
	add(KeyAlreadyWorn,
		"%s trägst du bereits.",
		"You are already wearing %s.")
	add(KeyWornTag,
		"(angelegt)",
		"(worn)")
	add(KeyHeadlampOn,
		"Du setzt den Helm auf und drückst auf die Stirnlampe. Ein Klick — und das "+
			"Dunkel weicht: zum ersten Mal siehst du das Verlies in Farbe.",
		"You put on the helmet and click the headlamp. The dark recoils — for the "+
			"first time you see the dungeon in colour.")
	add(KeySaved,
		"Dein Fortschritt ist in einem Trank versiegelt.",
		"Your progress is sealed in a potion.")
	add(KeySaveFailed,
		"Der Trank zerbrach — dein Fortschritt ließ sich nicht versiegeln.",
		"The potion shattered — your progress could not be saved.")
	add(KeySaveDisabled,
		"Hier gibt es keinen Ort, um einen Trank zu lagern (Speichern deaktiviert).",
		"There is nowhere to store a potion here (saving disabled).")
	add(KeyVerbHint,
		"Sprich mit dem Verlies: look · go <Richtung> · take <Ding> · "+
			"inspect <Ding> · wear <Ding> · inventory. Dinge kannst du auch per "+
			"Nummer wählen (z. B. »take 1«). Befehle wie /help beginnen mit »/«.",
		"Speak to the dungeon: look · go <direction> · take <thing> · "+
			"inspect <thing> · wear <thing> · inventory. You can also pick things by "+
			"number (e.g. \"take 1\"). Commands like /help begin with \"/\".")
	add(KeyEmptyInfo,
		"Du schweigst. Tippe einen Befehl — »look« zum Umsehen, »go north« zum "+
			"Gehen, oder »/help« für die volle Schriftrolle.",
		"You stay silent. Type a command — \"look\" to look around, \"go north\" to "+
			"move, or \"/help\" for the full scroll.")

	add(KeyHelpCmdHeader,
		"Befehle (beginnen mit »/«):",
		"Commands (begin with \"/\"):")
	add(KeyHelpVerbHeader,
		"Spielbefehle (so sprichst du mit dem Verlies):",
		"Game verbs (how you speak to the dungeon):")
	add(KeyVLook, "umsehen im Raum", "look around the room")
	add(KeyVGo, "in eine Richtung gehen (north/south/east/west/up/down)",
		"move in a direction (north/south/east/west/up/down)")
	add(KeyVTake, "einen Gegenstand nehmen (per Name oder Nummer)",
		"take an item (by name or number)")
	add(KeyVInspect, "etwas genauer ansehen (per Name oder Nummer)",
		"inspect something more closely (by name or number)")
	add(KeyVWear, "etwas anlegen (z. B. den Helm)", "wear something (e.g. the helmet)")
	add(KeyVInventory, "zeigt deine Hotbar (was du trägst)",
		"show your hotbar (what you carry)")
	add(KeyContinued,
		"Du nimmst deinen Weg wieder auf, wo du ihn verlassen hast.",
		"You take up your path where you left it.")
	add(KeyCmdSave,
		"versiegelt deinen Fortschritt",
		"seal your progress")
}
