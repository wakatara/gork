package engine

import (
	"strings"
	"testing"
)

// TestAllTextProperties verifies all readable items have .Text set
func TestAllTextProperties(t *testing.T) {
	g := NewGameV2("test")

	// Items that should have .Text properties (from audit)
	itemsWithText := []struct {
		id           string
		name         string
		textContains string
	}{
		{"advertisement", "advertisement", "WELCOME TO ZORK"},
		{"match", "matchbook", "PAPER SHUFFLING"},
		{"prayer", "prayer", "ancient script"},
		{"map", "map", "three clearings"},
		{"boat-label", "boat label", "FROBOZZ MAGIC BOAT"},
		{"guide", "tour guide", "Flood Control Dam #3"},
		{"tube", "tube", "Frobozz Magic Gunk"},
		{"engravings", "engravings", "bas reliefs"},
		{"owners-manual", "owner's manual", "Congratulations"},
		{"book", "black book", "Commandment"},
		{"wooden-door", "wooden door", "This space intentionally left blank"},
	}

	for _, item := range itemsWithText {
		t.Run(item.id, func(t *testing.T) {
			obj := g.Items[item.id]
			if obj == nil {
				t.Fatalf("Item %s not found", item.id)
			}

			if obj.Text == "" {
				t.Errorf("Item %s (%s) has empty .Text property", item.id, item.name)
			}

			if !strings.Contains(obj.Text, item.textContains) {
				t.Errorf("Item %s text does not contain expected substring %q.\nGot: %s",
					item.id, item.textContains, obj.Text)
			}

			if !obj.Flags.IsReadable {
				t.Errorf("Item %s should be marked IsReadable", item.id)
			}
		})
	}
}
