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
	KeySaved          = "game.saved"
	KeySaveFailed     = "game.save_failed"
	KeySaveDisabled   = "game.save_disabled"
	KeyContinued      = "game.continued"
	KeyVerbHint       = "game.verb_hint"
	KeyCmdSave        = "cmd.save"
)

func init() {
	add := func(key string, de, en string) {
		catalog[key] = map[Lang]string{DE: de, EN: en}
	}

	add(KeyExits,
		"Ausgänge: %s",
		"Exits: %s")
	add(KeyItemsHere,
		"Hier liegt: %s",
		"Here lies: %s")
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
		"Das verstehe ich nicht. Versuche: schau, gehe <Richtung>, nimm <Ding>, untersuche <Ding>, inventar.",
		"I don't understand. Try: look, go <direction>, take <thing>, examine <thing>, inventory.")
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
		"Sprich mit dem Verlies: schau · gehe <Richtung> · nimm <Ding> · "+
			"untersuche <Ding> · inventar. Befehle wie /help beginnen mit »/«.",
		"Speak to the dungeon: look · go <direction> · take <thing> · "+
			"examine <thing> · inventory. Commands like /help begin with \"/\".")
	add(KeyContinued,
		"Du nimmst deinen Weg wieder auf, wo du ihn verlassen hast.",
		"You take up your path where you left it.")
	add(KeyCmdSave,
		"versiegelt deinen Fortschritt",
		"seal your progress")
}
