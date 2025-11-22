package engine

import (
	"strings"
	"testing"
)

func TestOpenCloseCommands(t *testing.T) {
	tests := []struct {
		name        string
		commands    []string
		lastContains string
	}{
		{
			"open already open container",
			[]string{"open mailbox"},
			"already open",
		},
		{
			"close container",
			[]string{"close mailbox"},
			"Closed",
		},
		{
			"close already closed",
			[]string{"close mailbox", "close mailbox"},
			"already closed",
		},
		{
			"open non-container",
			[]string{"take leaflet", "open leaflet"},
			"can't open",
		},
		{
			"close non-container",
			[]string{"take leaflet", "close leaflet"},
			"can't close",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGameV2()
			var result string
			for _, cmd := range tt.commands {
				result = g.Process(cmd)
			}
			if !strings.Contains(strings.ToLower(result), strings.ToLower(tt.lastContains)) {
				t.Errorf("Expected output to contain %q, got %q", tt.lastContains, result)
			}
		})
	}
}

func TestReadCommand(t *testing.T) {
	tests := []struct {
		name         string
		commands     []string
		lastContains string
	}{
		{
			"read leaflet from mailbox",
			[]string{"take leaflet", "read leaflet"},
			"WELCOME TO ZORK",
		},
		{
			"read non-readable item",
			[]string{"read mailbox"},
			"How does one read",
		},
		{
			"read nothing",
			[]string{"read"},
			"What do you want to read",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGameV2()
			var result string
			for _, cmd := range tt.commands {
				result = g.Process(cmd)
			}
			if !strings.Contains(result, tt.lastContains) {
				t.Errorf("Expected output to contain %q, got %q", tt.lastContains, result)
			}
		})
	}
}

func TestLookInCommand(t *testing.T) {
	tests := []struct {
		name         string
		commands     []string
		lastContains string
	}{
		{
			"look in mailbox",
			[]string{"look in mailbox"},
			"leaflet",
		},
		{
			"look in empty container",
			[]string{"take leaflet", "look in mailbox"},
			"empty",
		},
		{
			"look in non-container",
			[]string{"take leaflet", "look in leaflet"},
			"can't look inside",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGameV2()
			var result string
			for _, cmd := range tt.commands {
				result = g.Process(cmd)
			}
			if !strings.Contains(strings.ToLower(result), strings.ToLower(tt.lastContains)) {
				t.Errorf("Expected output to contain %q, got %q", tt.lastContains, result)
			}
		})
	}
}

func TestHelpCommand(t *testing.T) {
	g := NewGameV2()

	result := g.Process("help")

	expectedPhrases := []string{
		"Available commands",
		"Movement:",
		"Actions:",
		"TAKE",
		"DROP",
		"OPEN",
		"CLOSE",
		"READ",
		"Obvious exits",
		"NORTH",
		"SOUTH",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(result, phrase) {
			t.Errorf("Help output should contain %q, got:\n%s", phrase, result)
		}
	}
}

func TestTurnOnOff(t *testing.T) {
	g := NewGameV2()

	// Note: lamp is in living-room, so we need to navigate there first
	// For this test, let's temporarily move the lamp to west-of-house
	lamp := g.Items["lamp"]
	lamp.Location = "west-of-house"
	g.Rooms["west-of-house"].AddItem("lamp")
	g.Rooms["living-room"].RemoveItem("lamp")

	tests := []struct {
		name         string
		commands     []string
		lastContains string
	}{
		{
			"turn off lit lamp",
			[]string{"take lamp", "turn off lamp"},
			"now off",
		},
		{
			"turn on unlit lamp",
			[]string{"take lamp", "turn off lamp", "turn on lamp"},
			"now on",
		},
		{
			"turn off already off",
			[]string{"take lamp", "turn off lamp", "turn off lamp"},
			"already off",
		},
		{
			"turn on already on",
			[]string{"take lamp", "turn on lamp"},
			"already on",
		},
		{
			"turn on non-light-source",
			[]string{"take leaflet", "turn on leaflet"},
			"can't turn that on",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset lamp state for each test
			g := NewGameV2()
			lamp := g.Items["lamp"]
			lamp.Location = "west-of-house"
			g.Rooms["west-of-house"].AddItem("lamp")
			g.Rooms["living-room"].RemoveItem("lamp")
			lamp.Flags.IsLit = true

			var result string
			for _, cmd := range tt.commands {
				result = g.Process(cmd)
			}
			if !strings.Contains(strings.ToLower(result), strings.ToLower(tt.lastContains)) {
				t.Errorf("Expected output to contain %q, got %q", tt.lastContains, result)
			}
		})
	}
}

func TestExamineContainer(t *testing.T) {
	g := NewGameV2()

	// Examine open container
	result := g.Process("examine mailbox")
	if !strings.Contains(result, "It is open") {
		t.Errorf("Examining open container should show it's open, got: %s", result)
	}
	if !strings.Contains(result, "leaflet") {
		t.Errorf("Examining open container should show contents, got: %s", result)
	}

	// Close and examine
	g.Process("close mailbox")
	result = g.Process("examine mailbox")
	if !strings.Contains(result, "It is closed") {
		t.Errorf("Examining closed container should show it's closed, got: %s", result)
	}

	// Transparent container still shows contents when closed
	if !strings.Contains(result, "leaflet") {
		t.Errorf("Examining transparent closed container should still show contents, got: %s", result)
	}
}

func TestTakeFromContainer(t *testing.T) {
	g := NewGameV2()

	// Should be able to take from open/transparent container
	result := g.Process("take leaflet")
	if !strings.Contains(result, "Taken") {
		t.Errorf("Should be able to take item from open container, got: %s", result)
	}

	// Verify it's in inventory
	result = g.Process("inventory")
	if !strings.Contains(result, "leaflet") {
		t.Errorf("Taken item should be in inventory, got: %s", result)
	}
}
