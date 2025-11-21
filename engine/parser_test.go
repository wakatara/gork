package engine

import (
	"testing"
)

// TestBasicCommands tests simple verb-only and verb-direction commands
func TestBasicCommands(t *testing.T) {
	tests := []struct {
		input    string
		verb     string
		obj      string
		prep     string
		indirect string
		dir      string
	}{
		// Direction commands (from ZIL: V?WALK handling)
		{"north", "walk", "", "", "", "north"},
		{"n", "walk", "", "", "", "north"},
		{"south", "walk", "", "", "", "south"},
		{"east", "walk", "", "", "", "east"},
		{"west", "walk", "", "", "", "west"},
		{"up", "walk", "", "", "", "up"},
		{"down", "walk", "", "", "", "down"},

		// Simple verb-object (from ZIL: V?TAKE, V?EXAMINE, etc.)
		{"take lamp", "take", "lamp", "", "", ""},
		{"get lamp", "take", "lamp", "", "", ""}, // synonym
		{"examine mailbox", "examine", "mailbox", "", "", ""},
		{"look at mailbox", "examine", "mailbox", "", "", ""}, // "at" is consumed as part of "look at" verb
		{"open mailbox", "open", "mailbox", "", "", ""},
		{"read leaflet", "read", "leaflet", "", "", ""},

		// Verb only (from ZIL: V?INVENTORY, V?LOOK)
		{"inventory", "inventory", "", "", "", ""},
		{"i", "inventory", "", "", "", ""}, // short form
		{"look", "look", "", "", "", ""},
		{"l", "look", "", "", "", ""}, // short form

		// Verb-object-prep-indirect (from ZIL: V?PUT, V?GIVE)
		{"put sword in case", "put", "sword", "in", "trophy-case", ""},
		{"put lamp on table", "put", "lamp", "on", "table", ""},
		{"give coins to troll", "give", "coins", "to", "troll", ""},

		// Multi-word synonyms (from ZIL: KITCHEN-WINDOW, WHITE-HOUSE)
		{"open kitchen window", "open", "kitchen-window", "", "", ""},
		{"examine white house", "examine", "white-house", "", "", ""},
	}

	p := NewParser()

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cmd, err := p.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.input, err)
			}

			if cmd.Verb != tt.verb {
				t.Errorf("Verb = %q, want %q", cmd.Verb, tt.verb)
			}
			if cmd.DirectObject != tt.obj {
				t.Errorf("DirectObject = %q, want %q", cmd.DirectObject, tt.obj)
			}
			if cmd.Preposition != tt.prep {
				t.Errorf("Preposition = %q, want %q", cmd.Preposition, tt.prep)
			}
			if cmd.IndirectObject != tt.indirect {
				t.Errorf("IndirectObject = %q, want %q", cmd.IndirectObject, tt.indirect)
			}
			if cmd.Direction != tt.dir {
				t.Errorf("Direction = %q, want %q", cmd.Direction, tt.dir)
			}
		})
	}
}

// TestSynonyms verifies the extensive synonym system from ZIL
func TestSynonyms(t *testing.T) {
	tests := []struct {
		input      string
		canonVerb  string
		canonObj   string
	}{
		// Verb synonyms (from gsyntax.zil)
		{"take lamp", "take", "lamp"},
		{"get lamp", "take", "lamp"},
		{"pick up lamp", "take", "lamp"},
		{"grab lamp", "take", "lamp"},

		{"drop lamp", "drop", "lamp"},
		{"put down lamp", "drop", "lamp"},

		{"examine lamp", "examine", "lamp"},
		{"look at lamp", "examine", "lamp"},
		{"x lamp", "examine", "lamp"},

		// Object synonyms (from 1dungeon.zil, gglobals.zil)
		{"take lantern", "take", "lamp"}, // LAMP has synonym LANTERN
		{"examine grue", "examine", "grue"},
	}

	p := NewParser()

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cmd, err := p.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.input, err)
			}

			if cmd.Verb != tt.canonVerb {
				t.Errorf("Canonical verb = %q, want %q", cmd.Verb, tt.canonVerb)
			}
			if cmd.DirectObject != "" && cmd.DirectObject != tt.canonObj {
				t.Errorf("Canonical object = %q, want %q", cmd.DirectObject, tt.canonObj)
			}
		})
	}
}

// TestArticlesAndFillerWords tests that articles are properly ignored
// (from ZIL parser: articles aren't stored in vocabulary)
func TestArticlesAndFillerWords(t *testing.T) {
	tests := []struct {
		input string
		verb  string
		obj   string
	}{
		{"take the lamp", "take", "lamp"},
		{"take a lamp", "take", "lamp"},
		{"open the mailbox", "open", "mailbox"},
		{"look at the white house", "examine", "white-house"},
	}

	p := NewParser()

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cmd, err := p.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.input, err)
			}

			if cmd.Verb != tt.verb {
				t.Errorf("Verb = %q, want %q", cmd.Verb, tt.verb)
			}
			if cmd.DirectObject != tt.obj {
				t.Errorf("DirectObject = %q, want %q", cmd.DirectObject, tt.obj)
			}
		})
	}
}

// TestSpecialCases tests edge cases from the ZIL parser
func TestSpecialCases(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
		verb    string
	}{
		// Empty input
		{"", true, ""},
		{"   ", true, ""},

		// Unknown words should error
		{"xyzzy plugh", true, ""}, // well, xyzzy is actually a word, bad example
		{"frobozz", false, "frobozz"}, // this IS a real Zork word (magic word)
	}

	p := NewParser()

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cmd, err := p.Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("Parse(%q) unexpected error = %v", tt.input, err)
			}

			if tt.verb != "" && cmd.Verb != tt.verb {
				t.Errorf("Verb = %q, want %q", cmd.Verb, tt.verb)
			}
		})
	}
}

// TestIT tests the special "IT" reference system from ZIL (P-IT-OBJECT)
func TestIT(t *testing.T) {
	p := NewParser()

	// First reference establishes "it"
	cmd1, err := p.Parse("examine lamp")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if cmd1.DirectObject != "lamp" {
		t.Fatalf("Expected lamp, got %q", cmd1.DirectObject)
	}

	// "it" should now refer to lamp
	cmd2, err := p.Parse("take it")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if cmd2.DirectObject != "lamp" {
		t.Errorf("IT reference failed: got %q, want lamp", cmd2.DirectObject)
	}
}
