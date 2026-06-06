// Package i18n holds the game's narrative text. German (de) is the primary,
// default language for story, rooms and riddles; the technical surface
// (commands, code, terminal) stays English by design (see decision D001).
package i18n

import "fmt"

// Lang identifies a supported language.
type Lang string

const (
	DE Lang = "de" // German — the narrative default
	EN Lang = "en" // English — fallback / optional
)

// current is the active narrative language. German by default.
var current = DE

// SetLang switches the active narrative language.
func SetLang(l Lang) { current = l }

// CurrentLang reports the active narrative language.
func CurrentLang() Lang { return current }

// Message keys. Using constants keeps lookups typo-proof.
const (
	KeyBannerSubtitle = "banner.subtitle"
	KeyWelcome        = "welcome"
	KeyPrompt         = "prompt"
	KeyHelpHeader     = "help.header"
	KeyHelpHint       = "help.hint"
	KeyUnknownCommand = "unknown.command"
	KeyUnknownInput   = "unknown.input"
	KeyGoodbye        = "goodbye"
	KeyEasterEgg      = "easteregg"

	// Built-in command summaries (shown in /help).
	KeyCmdHelp = "cmd.help"
	KeyCmdQuit = "cmd.quit"
)

// catalog maps a message key to its text per language.
var catalog = map[string]map[Lang]string{
	KeyBannerSubtitle: {
		DE: "Der Lehrmeister des Grimm-Clubs — ein Verlies aus Code und Märchen",
		EN: "TheGrimmClub's tutor — a dungeon of code and fairy tales",
	},
	KeyWelcome: {
		DE: "Willkommen, Lehrling. Vor dir liegt das Verlies der Brüder Grimm,\n" +
			"verwoben mit Nanobots, Androiden und echter künstlicher Intelligenz.\n" +
			"Tippe %s, um die Schriftrolle der Befehle zu sehen, oder %s zum Gehen.",
		EN: "Welcome, apprentice. Before you lies the dungeon of the Brothers Grimm,\n" +
			"woven with nanobots, androids and real artificial intelligence.\n" +
			"Type %s to see the scroll of commands, or %s to leave.",
	},
	KeyPrompt: {
		DE: "grimm> ",
		EN: "grimm> ",
	},
	KeyHelpHeader: {
		DE: "Schriftrolle der Befehle:",
		EN: "Scroll of commands:",
	},
	KeyHelpHint: {
		DE: "Befehle beginnen mit »/«. Alles andere flüsterst du dem Verlies zu.",
		EN: "Commands begin with \"/\". Anything else you whisper to the dungeon.",
	},
	KeyUnknownCommand: {
		DE: "Diesen Zauber (%s) kennt das Verlies nicht. Versuche %s.",
		EN: "The dungeon does not know that spell (%s). Try %s.",
	},
	KeyUnknownInput: {
		DE: "Das Verlies hört dein Flüstern »%s«, doch noch antwortet es nicht.",
		EN: "The dungeon hears your whisper \"%s\", but does not yet answer.",
	},
	KeyGoodbye: {
		DE: "Die Fackeln verlöschen. Bis bald, Lehrling.",
		EN: "The torches go dark. Until next time, apprentice.",
	},
	KeyEasterEgg: {
		DE: "Du sprichst die alten Worte — und plötzlich trägt dich nichts mehr.\n" +
			"Die Schwerkraft hat heute frei.\n",
		EN: "You speak the old words — and suddenly nothing holds you down.\n" +
			"Gravity has the day off.\n",
	},
	KeyCmdHelp: {
		DE: "zeigt die Schriftrolle der Befehle",
		EN: "show the scroll of commands",
	},
	KeyCmdQuit: {
		DE: "verlässt das Verlies",
		EN: "leave the dungeon",
	},
}

// T returns the text for key in the active language, formatted with args.
// Falls back to German, then to the raw key if a translation is missing.
func T(key string, args ...any) string {
	text := lookup(key)
	if len(args) == 0 {
		return text
	}
	return fmt.Sprintf(text, args...)
}

func lookup(key string) string {
	if byLang, ok := catalog[key]; ok {
		if s, ok := byLang[current]; ok {
			return s
		}
		if s, ok := byLang[DE]; ok {
			return s
		}
	}
	return key
}
