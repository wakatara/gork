package engine

import (
	"fmt"
	"testing"
)

func TestDebugLookAtWhiteHouse(t *testing.T) {
	p := NewParser()
	input := "look at the white house"

	fmt.Printf("\n=== Debugging: %q ===\n", input)

	// Step 1: Tokenize
	tokens := p.tokenize(input)
	fmt.Printf("After tokenize: %v\n", tokens)

	// Step 2: Resolve synonyms
	tokens = p.resolveSynonyms(tokens)
	fmt.Printf("After synonyms: %v\n", tokens)

	// Step 3: Check multi-word verb
	if len(tokens) >= 2 {
		twoWord := tokens[0] + " " + tokens[1]
		verb := p.vocabulary.GetVerb(twoWord)
		fmt.Printf("Two-word verb check: %q -> %q\n", twoWord, verb)
	}

	// Full parse
	cmd, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	fmt.Printf("Result: verb=%q, obj=%q\n", cmd.Verb, cmd.DirectObject)

	if cmd.Verb != "examine" {
		t.Errorf("Verb = %q, want examine", cmd.Verb)
	}
	if cmd.DirectObject != "white-house" {
		t.Errorf("DirectObject = %q, want white-house", cmd.DirectObject)
	}
}
