package engine

import (
	"fmt"
	"strings"
)

// Command represents a parsed player command (equivalent to ZIL's PRSA/PRSO/PRSI)
type Command struct {
	Verb           string // PRSA in ZIL (action)
	DirectObject   string // PRSO in ZIL (direct object)
	Preposition    string // Preposition (IN, ON, WITH, etc.)
	IndirectObject string // PRSI in ZIL (indirect object)
	Direction      string // Special case for movement (P-WALK-DIR in ZIL)
	Raw            string // Original input
}

// Parser handles natural language parsing (equivalent to ZIL's PARSER routine)
type Parser struct {
	vocabulary  *Vocabulary
	lastObject  string // P-IT-OBJECT in ZIL - tracks "it" references
}

// NewParser creates a new parser with initialized vocabulary
func NewParser() *Parser {
	return &Parser{
		vocabulary: NewVocabulary(),
		lastObject: "",
	}
}

// Parse converts a natural language input into a Command
// This implements the core logic from PARSER routine in gparser.zil
func (p *Parser) Parse(input string) (*Command, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("please enter a command")
	}

	// Store raw input
	cmd := &Command{Raw: input}

	// Tokenize input (equivalent to READ and LEXV in ZIL)
	tokens := p.tokenize(input)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("I don't understand that")
	}

	// Note: We DON'T resolve object synonyms here because we need to identify
	// multi-word phrases first (like "white house"). Object resolution happens
	// in resolveObject() after we've collected the noun phrase tokens.

	// Handle "it" references (P-IT-OBJECT in ZIL)
	tokens = p.resolveIt(tokens)

	// Check if this is a direction command
	if len(tokens) == 1 {
		if dir := p.vocabulary.GetDirection(tokens[0]); dir != "" {
			cmd.Verb = "walk"
			cmd.Direction = dir
			return cmd, nil
		}
	}

	// Check for "go <direction>" pattern
	if len(tokens) == 2 && (tokens[0] == "go" || tokens[0] == "walk") {
		if dir := p.vocabulary.GetDirection(tokens[1]); dir != "" {
			cmd.Verb = "walk"
			cmd.Direction = dir
			return cmd, nil
		}
	}

	// Check for multi-word verbs like "look at" -> "examine", "put down" -> "drop"
	verb := ""
	verbTokens := 1
	if len(tokens) >= 2 {
		twoWord := tokens[0] + " " + tokens[1]
		if v := p.vocabulary.GetVerb(twoWord); v != "" {
			verb = v
			verbTokens = 2
		}
	}

	// Try single word verb
	if verb == "" {
		verb = p.vocabulary.GetVerb(tokens[0])
		if verb == "" {
			return nil, fmt.Errorf("I don't know the word %q", tokens[0])
		}
	}
	cmd.Verb = verb

	// Handle verb-only commands (INVENTORY, LOOK, etc.)
	if len(tokens) == verbTokens {
		return cmd, nil
	}

	// Parse remaining tokens for objects and prepositions
	// This implements the noun clause parsing from ZIL
	pos := verbTokens
	objTokens := []string{}

	// Collect tokens until we hit a preposition or run out
	for pos < len(tokens) {
		if p.vocabulary.IsPreposition(tokens[pos]) {
			break
		}
		objTokens = append(objTokens, tokens[pos])
		pos++
	}

	// Resolve direct object
	if len(objTokens) > 0 {
		// Special case for save/restore commands - allow arbitrary filenames
		if verb == "save" || verb == "restore" {
			cmd.DirectObject = strings.Join(objTokens, "_")
		} else {
			obj := p.resolveObject(objTokens)
			if obj == "" {
				return nil, fmt.Errorf("I don't know the word %q", objTokens[0])
			}
			cmd.DirectObject = obj
			// Track for "it" references (P-IT-OBJECT)
			p.lastObject = obj
		}
	}

	// Check for preposition and indirect object
	if pos < len(tokens) {
		prep := tokens[pos]
		if !p.vocabulary.IsPreposition(prep) {
			return nil, fmt.Errorf("I don't understand how to use %q here", prep)
		}
		cmd.Preposition = prep
		pos++

		// Collect indirect object tokens
		indirectTokens := []string{}
		for pos < len(tokens) {
			indirectTokens = append(indirectTokens, tokens[pos])
			pos++
		}

		if len(indirectTokens) > 0 {
			indirect := p.resolveObject(indirectTokens)
			if indirect == "" {
				return nil, fmt.Errorf("I don't know the word %q", indirectTokens[0])
			}
			cmd.IndirectObject = indirect
		}
	}

	return cmd, nil
}

// tokenize breaks input into words (equivalent to ZIL's LEXV processing)
func (p *Parser) tokenize(input string) []string {
	// Convert to lowercase and split on whitespace
	input = strings.ToLower(input)
	words := strings.Fields(input)

	// Remove punctuation and articles
	filtered := []string{}
	for _, word := range words {
		// Remove trailing punctuation
		word = strings.TrimRight(word, ".,!?;:")

		// Skip articles (not in ZIL vocabulary)
		if word == "the" || word == "a" || word == "an" {
			continue
		}

		if word != "" {
			filtered = append(filtered, word)
		}
	}

	return filtered
}

// resolveSynonyms is no longer used - we resolve objects in resolveObject()
// after identifying multi-word phrases
func (p *Parser) resolveSynonyms(tokens []string) []string {
	// Just return tokens unchanged - resolution happens later
	return tokens
}

// resolveIt replaces "it" with the last referenced object (P-IT-OBJECT in ZIL)
func (p *Parser) resolveIt(tokens []string) []string {
	if p.lastObject == "" {
		return tokens
	}

	result := make([]string, len(tokens))
	for i, token := range tokens {
		if token == "it" || token == "them" {
			result[i] = p.lastObject
		} else {
			result[i] = token
		}
	}
	return result
}

// resolveObject attempts to match multi-word objects (like "kitchen window")
func (p *Parser) resolveObject(tokens []string) string {
	// Try multi-word match first (longest match wins)
	for length := len(tokens); length > 0; length-- {
		phrase := strings.Join(tokens[:length], " ")
		if obj := p.vocabulary.GetObject(phrase); obj != "" {
			return obj
		}
	}

	// Try hyphenated version (ZIL uses hyphens: KITCHEN-WINDOW)
	if len(tokens) > 1 {
		hyphenated := strings.Join(tokens, "-")
		if obj := p.vocabulary.GetObject(hyphenated); obj != "" {
			return obj
		}
	}

	// Single word
	if len(tokens) == 1 {
		return p.vocabulary.GetObject(tokens[0])
	}

	// If we still haven't found it, try each token individually
	// and return the first match (this handles cases like "white house"
	// where "house" alone might be valid)
	for _, token := range tokens {
		if obj := p.vocabulary.GetObject(token); obj != "" {
			return obj
		}
	}

	return ""
}
