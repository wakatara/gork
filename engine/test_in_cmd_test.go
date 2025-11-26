package engine

import (
	"testing"
)

func TestInDirection(t *testing.T) {
	p := NewParser()
	
	// Test parsing "in" as a direction
	cmd, err := p.Parse("in")
	if err != nil {
		t.Fatalf("Failed to parse 'in': %v", err)
	}
	
	t.Logf("Command: %+v", cmd)
	
	if cmd.Verb != "walk" {
		t.Errorf("Expected verb 'walk', got %q", cmd.Verb)
	}
	if cmd.Direction != "in" {
		t.Errorf("Expected direction 'in', got %q", cmd.Direction)
	}
}
