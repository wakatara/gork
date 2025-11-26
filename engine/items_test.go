package engine

import (
	"testing"
)

func TestAllItemsCreated(t *testing.T) {
	g := NewGameV2("test")

	// Count items
	itemCount := len(g.Items)
	if itemCount < 80 {
		t.Errorf("Expected at least 80 items, got %d", itemCount)
	}
	t.Logf("Created %d items total", itemCount)
}

func TestTreasuresExist(t *testing.T) {
	g := NewGameV2("test")

	treasures := []string{
		"diamond",
		"emerald",
		"chalice",
		"jade",
		"coins",
		"painting",
		"bracelet",
		"bauble",
		"scarab",
		"pot-of-gold",
		"trident",
		"sceptre",
		"egg",
	}

	for _, treasureID := range treasures {
		item := g.Items[treasureID]
		if item == nil {
			t.Errorf("Treasure %s does not exist", treasureID)
			continue
		}

		if !item.Flags.IsTreasure {
			t.Errorf("Item %s should have IsTreasure flag", treasureID)
		}

		if item.Value <= 0 {
			t.Errorf("Treasure %s should have Value > 0, got %d", treasureID, item.Value)
		}
	}
}

func TestWeaponsExist(t *testing.T) {
	g := NewGameV2("test")

	weapons := []string{"sword", "knife", "rusty-knife", "stiletto", "axe"}

	for _, weaponID := range weapons {
		item := g.Items[weaponID]
		if item == nil {
			t.Errorf("Weapon %s does not exist", weaponID)
			continue
		}

		if !item.Flags.IsWeapon {
			t.Errorf("Item %s should have IsWeapon flag", weaponID)
		}
	}
}

func TestContainersExist(t *testing.T) {
	g := NewGameV2("test")

	containers := []string{
		"mailbox",
		"trophy-case",
		"bottle",
		"coffin",
		"sandwich-bag",
		"large-bag",
		"nest",
	}

	for _, containerID := range containers {
		item := g.Items[containerID]
		if item == nil {
			t.Errorf("Container %s does not exist", containerID)
			continue
		}

		if !item.Flags.IsContainer {
			t.Errorf("Item %s should have IsContainer flag", containerID)
		}
	}
}

func TestStartingItems(t *testing.T) {
	g := NewGameV2("test")

	// Test mailbox is in west-of-house
	mailbox := g.Items["mailbox"]
	if mailbox == nil {
		t.Fatal("Mailbox does not exist")
	}
	if mailbox.Location != "west-of-house" {
		t.Errorf("Mailbox should be in west-of-house, got %s", mailbox.Location)
	}

	// Test leaflet is in mailbox
	leaflet := g.Items["leaflet"]
	if leaflet == nil {
		t.Fatal("Leaflet does not exist")
	}
	if leaflet.Location != "mailbox" {
		t.Errorf("Leaflet should be in mailbox, got %s", leaflet.Location)
	}

	// Test lamp is in living-room
	lamp := g.Items["lamp"]
	if lamp == nil {
		t.Fatal("Lamp does not exist")
	}
	if lamp.Location != "living-room" {
		t.Errorf("Lamp should be in living-room, got %s", lamp.Location)
	}

	// Test window is in behind-house
	window := g.Items["kitchen-window"]
	if window == nil {
		t.Fatal("Kitchen window does not exist")
	}
	if window.Location != "behind-house" {
		t.Errorf("Kitchen window should be in behind-house, got %s", window.Location)
	}
}

func TestItemAliases(t *testing.T) {
	g := NewGameV2("test")

	tests := []struct {
		itemID string
		alias  string
	}{
		{"mailbox", "box"},
		{"lamp", "lantern"},
		{"lamp", "light"},
		{"sword", "blade"},
		{"chalice", "cup"},
		{"kitchen-window", "window"},
	}

	for _, tt := range tests {
		item := g.Items[tt.itemID]
		if item == nil {
			t.Errorf("Item %s does not exist", tt.itemID)
			continue
		}

		found := false
		for _, alias := range item.Aliases {
			if alias == tt.alias {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Item %s should have alias %q, aliases: %v", tt.itemID, tt.alias, item.Aliases)
		}
	}
}

func TestLightSources(t *testing.T) {
	g := NewGameV2("test")

	lightSources := []struct {
		id    string
		isLit bool
	}{
		{"lamp", true},
		{"torch", true},
		{"candles", false},
	}

	for _, ls := range lightSources {
		item := g.Items[ls.id]
		if item == nil {
			t.Errorf("Light source %s does not exist", ls.id)
			continue
		}

		if !item.Flags.IsLightSource {
			t.Errorf("Item %s should have IsLightSource flag", ls.id)
		}

		if item.Flags.IsLit != ls.isLit {
			t.Errorf("Item %s IsLit should be %v, got %v", ls.id, ls.isLit, item.Flags.IsLit)
		}
	}
}

func TestReadableItems(t *testing.T) {
	g := NewGameV2("test")

	readableItems := []string{
		"leaflet",
		"book",
		"advertisement",
		"guide",
		"map",
		"match",
		"owners-manual",
	}

	for _, itemID := range readableItems {
		item := g.Items[itemID]
		if item == nil {
			t.Errorf("Readable item %s does not exist", itemID)
			continue
		}

		if !item.Flags.IsReadable {
			t.Errorf("Item %s should have IsReadable flag", itemID)
		}
	}
}

func TestItemsInRooms(t *testing.T) {
	g := NewGameV2("test")

	// Check that items are added to their room's Contents
	westOfHouse := g.Rooms["west-of-house"]
	if westOfHouse == nil {
		t.Fatal("West of house does not exist")
	}

	foundMailbox := false
	for _, itemID := range westOfHouse.Contents {
		if itemID == "mailbox" {
			foundMailbox = true
			break
		}
	}

	if !foundMailbox {
		t.Errorf("Mailbox should be in west-of-house Contents, got: %v", westOfHouse.Contents)
	}
}
