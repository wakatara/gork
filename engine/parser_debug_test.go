package engine

import (
	"testing"
)

func TestDebugLookAtWhiteHouse(t *testing.T) {
	p := NewParser()
	input := "look at the white house"

	t.Logf("\n=== Debugging: %q ===", input)

	// Step 1: Tokenize
	tokens := p.tokenize(input)
	t.Logf("After tokenize: %v", tokens)

	// Step 2: Resolve synonyms
	tokens = p.resolveSynonyms(tokens)
	t.Logf("After synonyms: %v", tokens)

	// Step 3: Check multi-word verb
	if len(tokens) >= 2 {
		twoWord := tokens[0] + " " + tokens[1]
		verb := p.vocabulary.GetVerb(twoWord)
		t.Logf("Two-word verb check: %q -> %q", twoWord, verb)
	}

	// Full parse
	cmd, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	t.Logf("Result: verb=%q, obj=%q", cmd.Verb, cmd.DirectObject)

	if cmd.Verb != "examine" {
		t.Errorf("Verb = %q, want examine", cmd.Verb)
	}
	if cmd.DirectObject != "white-house" {
		t.Errorf("DirectObject = %q, want white-house", cmd.DirectObject)
	}
}
