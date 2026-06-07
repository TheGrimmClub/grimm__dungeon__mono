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
	KeyNothingSpecial = "game.nothing_special"
	KeyUnknownVerb    = "game.unknown_verb"
	KeyWearWhat       = "game.wear_what"
	KeyWorn           = "game.worn"
	KeyCannotWear     = "game.cannot_wear"
	KeyAlreadyWorn    = "game.already_worn"
	KeyWornTag        = "game.worn_tag"
	KeyHeadlampOn     = "game.headlamp_on"

	// Puzzles
	KeySolveHint      = "game.solve_hint"
	KeyNoActivePuzzle = "game.no_active_puzzle"
	KeySolveWhat      = "game.solve_what"
	KeyPuzzleSolved   = "game.puzzle_solved"
	KeyPuzzleWrong    = "game.puzzle_wrong"
	KeyPuzzleBroken   = "game.puzzle_broken"
	KeyHintLabel      = "game.hint_label"
	KeyLockedFootnote = "game.locked_footnote"
	KeyVSolve         = "verb.solve"

	// Class selection
	KeyCmdClass     = "cmd.class"
	KeyClassHeader  = "class.header"
	KeyClassChoose  = "class.choose"
	KeyClassChosen  = "class.chosen"
	KeyClassUnknown = "class.unknown"

	// Voice (TTS)
	KeyCmdVoice         = "cmd.voice"
	KeyVoiceOn          = "voice.on"
	KeyVoiceOff         = "voice.off"
	KeyVoiceUnavailable = "voice.unavailable"

	// Alchemist (git)
	KeyCmdAlchemist   = "cmd.alchemist"
	KeyAlchemistNoDir = "alchemist.no_dir"

	// Terminal & editor surfaces
	KeyCmdTerminal    = "cmd.terminal"
	KeyCmdBook        = "cmd.book"
	KeyTerminalEnter  = "terminal.enter"
	KeyTerminalReturn = "terminal.return"
	KeyBookUsage      = "book.usage"
	KeyBookOpen       = "book.open"
	KeyBookClosed     = "book.closed"
	KeyNoEditor       = "book.no_editor"
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
	add(KeyNothingSpecial,
		"Du siehst genauer hin, findest aber nichts Besonderes.",
		"You take a closer look, but find nothing special.")
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
			"Dunkel weicht: zum ersten Mal siehst du das Verlies in Farbe. Im Visier "+
			"erwacht eine Anzeige (Inventar & Karte), und eine Stimme flüstert: "+
			"»Tippe /voice, und ich lese dir vor.«",
		"You put on the helmet and click the headlamp. The dark recoils — for the "+
			"first time you see the dungeon in colour. A display lights up in your "+
			"visor (inventory & map), and a voice whispers: \"Type /voice and I'll "+
			"read to you.\"")

	add(KeySolveHint,
		"(Antworte mit »solve <deine Lösung>«.)",
		"(Answer with \"solve <your solution>\".)")
	add(KeyNoActivePuzzle,
		"Hier gibt es gerade nichts zu lösen.",
		"There is nothing to solve right now.")
	add(KeySolveWhat,
		"Womit möchtest du es versuchen?",
		"What do you want to try?")
	add(KeyPuzzleSolved,
		"Richtig! Mit leisem Knirschen gibt der Weg nach.",
		"Correct! With a soft grind, the way gives.")
	add(KeyPuzzleWrong,
		"Das war nicht die Lösung.",
		"That was not the solution.")
	add(KeyPuzzleBroken,
		"Dieses Rätsel ist zerbrochen — sag dem Spielleiter Bescheid.",
		"This puzzle is broken — tell the game master.")
	add(KeyHintLabel,
		"Tipp:",
		"Hint:")
	add(KeyLockedFootnote,
		"* ein Rätsel versperrt diesen Weg — versuche dort zu »go« und dann »solve«.",
		"* a puzzle seals this way — try to \"go\" there, then \"solve\".")
	add(KeyVSolve, "ein Rätsel lösen (solve <Lösung>)", "solve a puzzle (solve <solution>)")

	add(KeyCmdClass,
		"wähle deinen Pfad (Klasse)",
		"choose your path (class)")
	add(KeyClassHeader,
		"Noch bist du nur »Human«. Welchen Pfad wählst du?",
		"For now you are only \"Human\". Which path do you choose?")
	add(KeyClassChoose,
		"Wähle mit »/class <name>«, z. B. »/class alchemist«.",
		"Choose with \"/class <name>\", e.g. \"/class alchemist\".")
	add(KeyClassChosen,
		"Von nun an bist du %s. Dein Pfad hat begonnen.",
		"From now on you are a %s. Your path has begun.")
	add(KeyClassUnknown,
		"»%s« ist kein Pfad, den ich kenne. Tippe »/class« für die Liste.",
		"\"%s\" is not a path I know. Type \"/class\" for the list.")

	add(KeyCmdVoice,
		"liest den Text laut vor (an/aus)",
		"read the text aloud (on/off)")
	add(KeyVoiceOn,
		"Die Stimme des Helms erwacht — von nun an liest sie dir vor.",
		"The helmet's voice awakens — from now on it reads to you.")
	add(KeyVoiceOff,
		"Die Stimme des Helms verstummt.",
		"The helmet's voice falls silent.")
	add(KeyVoiceUnavailable,
		"Dieser Rechner hat keine Stimme, die vorlesen könnte.",
		"This machine has no voice to read aloud.")

	add(KeyCmdAlchemist,
		"braue mit dem Alchemisten (git): /alchemist <trank>",
		"brew with the alchemist (git): /alchemist <potion>")
	add(KeyAlchemistNoDir,
		"Hier gibt es keinen Kessel — es ist kein Arbeitsverzeichnis eingerichtet.",
		"There is no cauldron here — no working directory is set up.")

	add(KeyCmdTerminal,
		"öffnet eine echte Shell in deinem Kessel",
		"open a real shell in your cauldron")
	add(KeyCmdBook,
		"öffnet das Buch (Editor): /book <datei>",
		"open the book (editor): /book <file>")
	add(KeyTerminalEnter,
		"Du steigst hinab in die Shell deines Kessels …",
		"You descend into your cauldron's shell …")
	add(KeyTerminalReturn,
		"Du tauchst aus der Shell wieder ins Verlies auf.",
		"You surface from the shell back into the dungeon.")
	add(KeyBookUsage,
		"Welches Buch? Tippe »/book <datei>«, z. B. »/book zauber.py«.",
		"Which book? Type \"/book <file>\", e.g. \"/book spell.py\".")
	add(KeyBookOpen,
		"Du schlägst das Buch »%s« auf …",
		"You open the book \"%s\" …")
	add(KeyBookClosed,
		"Du klappst das Buch zu. Vergiss nicht zu brauen (»/alchemist brew«)!",
		"You close the book. Don't forget to brew (\"/alchemist brew\")!")
	add(KeyNoEditor,
		"Kein Editor gefunden. Setze $EDITOR oder installiere micro/nano.",
		"No editor found. Set $EDITOR or install micro/nano.")
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
			"inspect <Ding> · wear <Ding> · solve <Lösung> · inventory. Dinge wählst "+
			"du auch per Nummer (»take 1«). Befehle wie /help beginnen mit »/«.",
		"Speak to the dungeon: look · go <direction> · take <thing> · "+
			"inspect <thing> · wear <thing> · solve <answer> · inventory. You can pick "+
			"things by number (\"take 1\"). Commands like /help begin with \"/\".")
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
