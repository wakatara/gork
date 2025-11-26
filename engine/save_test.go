package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveAndRestore(t *testing.T) {
	// Create a new game
	g := NewGameV2("test")

	// Modify game state
	g.Location = "kitchen"
	g.Score = 50
	g.Moves = 25
	g.Flags["grate-open"] = true
	g.Flags["lamp-on"] = true
	g.Player.Inventory = append(g.Player.Inventory, "lamp", "sword")
	g.Player.Health = 80

	// Modify some item states
	if lamp, ok := g.Items["lamp"]; ok {
		lamp.Location = "inventory"
		lamp.Flags.IsLit = true
		lamp.Fuel = 200
	}

	// Modify NPC state
	if thief, ok := g.NPCs["thief"]; ok {
		thief.Location = "treasure-room"
		thief.Hostile = true
	}

	// Save the game
	testFilename := "test_save"
	err := g.Save(testFilename)
	if err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}

	// Verify save file exists
	savePath, err := GetSavePath(testFilename)
	if err != nil {
		t.Fatalf("Failed to get save path: %v", err)
	}

	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		t.Fatalf("Save file does not exist: %s", savePath)
	}

	// Create a new game (clean state)
	g2 := NewGameV2("test")

	// Verify it has default state
	if g2.Location != "west-of-house" {
		t.Errorf("New game should start at west-of-house, got %s", g2.Location)
	}
	if g2.Score != 0 {
		t.Errorf("New game should have score 0, got %d", g2.Score)
	}

	// Restore the saved game
	err = g2.Restore(testFilename)
	if err != nil {
		t.Fatalf("Failed to restore game: %v", err)
	}

	// Verify all state was restored correctly
	if g2.Location != "kitchen" {
		t.Errorf("Expected location kitchen, got %s", g2.Location)
	}
	if g2.Score != 50 {
		t.Errorf("Expected score 50, got %d", g2.Score)
	}
	if g2.Moves != 25 {
		t.Errorf("Expected moves 25, got %d", g2.Moves)
	}
	if !g2.Flags["grate-open"] {
		t.Error("Expected grate-open flag to be true")
	}
	if !g2.Flags["lamp-on"] {
		t.Error("Expected lamp-on flag to be true")
	}

	// Verify player state
	if g2.Player.Health != 80 {
		t.Errorf("Expected player health 80, got %d", g2.Player.Health)
	}
	if len(g2.Player.Inventory) != 2 {
		t.Errorf("Expected 2 items in inventory, got %d", len(g2.Player.Inventory))
	}

	// Verify item state
	if lamp, ok := g2.Items["lamp"]; ok {
		if lamp.Location != "inventory" {
			t.Errorf("Expected lamp in inventory, got %s", lamp.Location)
		}
		if !lamp.Flags.IsLit {
			t.Error("Expected lamp to be lit")
		}
		if lamp.Fuel != 200 {
			t.Errorf("Expected lamp fuel 200, got %d", lamp.Fuel)
		}
	}

	// Verify NPC state
	if thief, ok := g2.NPCs["thief"]; ok {
		if thief.Location != "treasure-room" {
			t.Errorf("Expected thief in treasure-room, got %s", thief.Location)
		}
		if !thief.Hostile {
			t.Error("Expected thief to be hostile")
		}
	}

	// Clean up test file
	os.Remove(savePath)
}

func TestSaveCommand(t *testing.T) {
	g := NewGameV2("test")
	g.Location = "kitchen"
	g.Score = 100

	// Test save command with filename
	result := g.Process("save testsave")
	if !strings.Contains(result, "Game saved") {
		t.Errorf("Expected save success message, got: %s", result)
	}

	// Verify file exists
	savePath, _ := GetSavePath("testsave")
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		t.Error("Save file was not created")
	}

	// Clean up
	os.Remove(savePath)
}

func TestRestoreCommand(t *testing.T) {
	// Create and save a game
	g := NewGameV2("test")
	g.Location = "kitchen"
	g.Score = 75
	g.Moves = 10

	testFilename := "test_restore"
	err := g.Save(testFilename)
	if err != nil {
		t.Fatalf("Failed to save test game: %v", err)
	}

	// Create new game and restore
	g2 := NewGameV2("test")
	result := g2.Process("restore " + testFilename)

	if !strings.Contains(result, "Game restored") {
		t.Errorf("Expected restore success message, got: %s", result)
	}

	// Verify state was restored
	if g2.Location != "kitchen" {
		t.Errorf("Expected location kitchen, got %s", g2.Location)
	}
	if g2.Score != 75 {
		t.Errorf("Expected score 75, got %d", g2.Score)
	}

	// Clean up
	savePath, _ := GetSavePath(testFilename)
	os.Remove(savePath)
}

func TestRestoreCommandListsSaves(t *testing.T) {
	// Create a test save
	g := NewGameV2("test")
	testFilename := "test_list"
	g.Save(testFilename)

	// Test restore without filename (should list saves)
	g2 := NewGameV2("test")
	result := g2.Process("restore")

	if !strings.Contains(result, "Available saved games") {
		t.Errorf("Expected save list, got: %s", result)
	}

	if !strings.Contains(result, testFilename) {
		t.Errorf("Expected test save in list, got: %s", result)
	}

	// Clean up
	savePath, _ := GetSavePath(testFilename)
	os.Remove(savePath)
}

func TestGetSaveDir(t *testing.T) {
	saveDir, err := getSaveDir()
	if err != nil {
		t.Fatalf("Failed to get save directory: %v", err)
	}

	// Verify directory exists
	info, err := os.Stat(saveDir)
	if err != nil {
		t.Fatalf("Save directory does not exist: %v", err)
	}

	if !info.IsDir() {
		t.Error("Save path is not a directory")
	}

	// Verify it's in a reasonable location
	if !strings.Contains(saveDir, "gork") {
		t.Errorf("Expected 'gork' in save path, got: %s", saveDir)
	}
}

func TestListSaves(t *testing.T) {
	// Create a few test saves
	g := NewGameV2("test")
	saves := []string{"test1", "test2", "test3"}

	for _, name := range saves {
		err := g.Save(name)
		if err != nil {
			t.Fatalf("Failed to create test save %s: %v", name, err)
		}
	}

	// List saves
	foundSaves, err := ListSaves()
	if err != nil {
		t.Fatalf("Failed to list saves: %v", err)
	}

	// Verify all test saves are listed
	for _, name := range saves {
		found := false
		expectedName := name + ".json"
		for _, save := range foundSaves {
			if save == expectedName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find save %s in list", expectedName)
		}
	}

	// Clean up
	for _, name := range saves {
		savePath, _ := GetSavePath(name)
		os.Remove(savePath)
	}
}

func TestSaveFileFormat(t *testing.T) {
	g := NewGameV2("test")
	g.Location = "kitchen"
	g.Score = 42

	testFilename := "test_format"
	err := g.Save(testFilename)
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Read the save file
	savePath, _ := GetSavePath(testFilename)
	data, err := os.ReadFile(savePath)
	if err != nil {
		t.Fatalf("Failed to read save file: %v", err)
	}

	content := string(data)

	// Verify JSON structure
	if !strings.Contains(content, `"version"`) {
		t.Error("Save file missing version field")
	}
	if !strings.Contains(content, `"timestamp"`) {
		t.Error("Save file missing timestamp field")
	}
	if !strings.Contains(content, `"game_state"`) {
		t.Error("Save file missing game_state field")
	}
	if !strings.Contains(content, `"kitchen"`) {
		t.Error("Save file missing location data")
	}
	if !strings.Contains(content, `42`) {
		t.Error("Save file missing score data")
	}

	// Clean up
	os.Remove(savePath)
}

func TestSaveFilenameExtension(t *testing.T) {
	g := NewGameV2("test")

	// Test with extension
	err := g.Save("test.json")
	if err != nil {
		t.Fatalf("Failed to save with .json extension: %v", err)
	}

	// Test without extension
	err = g.Save("test2")
	if err != nil {
		t.Fatalf("Failed to save without extension: %v", err)
	}

	// Verify both files exist
	path1, _ := GetSavePath("test.json")
	path2, _ := GetSavePath("test2.json")

	if _, err := os.Stat(path1); os.IsNotExist(err) {
		t.Error("Save with extension failed")
	}
	if _, err := os.Stat(path2); os.IsNotExist(err) {
		t.Error("Save without extension failed")
	}

	// Clean up
	os.Remove(path1)
	os.Remove(path2)
}

func TestAutoGeneratedFilename(t *testing.T) {
	g := NewGameV2("test")

	// Save without filename
	err := g.Save("")
	if err != nil {
		t.Fatalf("Failed to save with auto-generated filename: %v", err)
	}

	// List saves to find the auto-generated file
	saves, err := ListSaves()
	if err != nil {
		t.Fatalf("Failed to list saves: %v", err)
	}

	// Should have at least one save with gork_save_ prefix
	found := false
	var autoSave string
	for _, save := range saves {
		if strings.HasPrefix(save, "gork_save_") {
			found = true
			autoSave = save
			break
		}
	}

	if !found {
		t.Error("Auto-generated save file not found")
	}

	// Clean up
	if autoSave != "" {
		savePath, _ := GetSavePath(strings.TrimSuffix(autoSave, filepath.Ext(autoSave)))
		os.Remove(savePath)
	}
}
