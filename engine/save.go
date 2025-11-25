package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SaveGame represents a serializable game state
type SaveGame struct {
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	GameState GameState `json:"game_state"`
}

// GameState holds all the dynamic game state that needs to be saved
type GameState struct {
	Location      string            `json:"location"`
	Score         int               `json:"score"`
	Moves         int               `json:"moves"`
	Flags         map[string]bool   `json:"flags"`
	GameOver      bool              `json:"game_over"`
	Won           bool              `json:"won"`
	PlayerState   PlayerState       `json:"player"`
	ItemStates    map[string]ItemState `json:"items"`
	NPCStates     map[string]NPCState  `json:"npcs"`
}

// PlayerState holds serializable player data
type PlayerState struct {
	Inventory []string `json:"inventory"`
	Health    int      `json:"health"`
	MaxWeight int      `json:"max_weight"`
}

// ItemState holds serializable item data
type ItemState struct {
	Location  string    `json:"location"`
	Flags     ItemFlags `json:"flags"`
	Fuel      int       `json:"fuel,omitempty"`
	GlowLevel int       `json:"glow_level,omitempty"`
}

// NPCState holds serializable NPC data
type NPCState struct {
	Location  string   `json:"location"`
	Flags     NPCFlags `json:"flags"`
	Strength  int      `json:"strength"`
	Weapon    string   `json:"weapon,omitempty"`
	Inventory []string `json:"inventory,omitempty"`
	Hostile   bool     `json:"hostile"`
}

// getSaveDir returns the platform-specific save directory
func getSaveDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}

	saveDir := filepath.Join(configDir, "gork", "saves")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create save directory: %w", err)
	}

	return saveDir, nil
}

// Save serializes the current game state to a JSON file
func (g *GameV2) Save(filename string) error {
	// If no filename provided, generate one with timestamp
	if filename == "" {
		filename = fmt.Sprintf("gork_save_%s.json", time.Now().Format("20060102_150405"))
	}

	// Ensure .json extension
	if filepath.Ext(filename) != ".json" {
		filename += ".json"
	}

	// Get save directory
	saveDir, err := getSaveDir()
	if err != nil {
		return err
	}

	savePath := filepath.Join(saveDir, filename)

	// Create save game structure
	save := SaveGame{
		Version:   "1.0",
		Timestamp: time.Now(),
		GameState: g.serializeState(),
	}

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(save, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal save data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(savePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	return nil
}

// serializeState extracts the current game state into a serializable format
func (g *GameV2) serializeState() GameState {
	state := GameState{
		Location: g.Location,
		Score:    g.Score,
		Moves:    g.Moves,
		Flags:    make(map[string]bool),
		GameOver: g.GameOver,
		Won:      g.Won,
		PlayerState: PlayerState{
			Inventory: g.Player.Inventory,
			Health:    g.Player.Health,
			MaxWeight: g.Player.MaxWeight,
		},
		ItemStates: make(map[string]ItemState),
		NPCStates:  make(map[string]NPCState),
	}

	// Copy flags
	for k, v := range g.Flags {
		state.Flags[k] = v
	}

	// Serialize items (only dynamic state, not static definitions)
	for id, item := range g.Items {
		state.ItemStates[id] = ItemState{
			Location:  item.Location,
			Flags:     item.Flags,
			Fuel:      item.Fuel,
			GlowLevel: item.GlowLevel,
		}
	}

	// Serialize NPCs (only dynamic state)
	for id, npc := range g.NPCs {
		state.NPCStates[id] = NPCState{
			Location:  npc.Location,
			Flags:     npc.Flags,
			Strength:  npc.Strength,
			Weapon:    npc.Weapon,
			Inventory: npc.Inventory,
			Hostile:   npc.Hostile,
		}
	}

	return state
}

// Restore loads a saved game state from a JSON file
func (g *GameV2) Restore(filename string) error {
	// Ensure .json extension
	if filepath.Ext(filename) != ".json" {
		filename += ".json"
	}

	// Get save directory
	saveDir, err := getSaveDir()
	if err != nil {
		return err
	}

	savePath := filepath.Join(saveDir, filename)

	// Read file
	data, err := os.ReadFile(savePath)
	if err != nil {
		return fmt.Errorf("failed to read save file: %w", err)
	}

	// Unmarshal JSON
	var save SaveGame
	if err := json.Unmarshal(data, &save); err != nil {
		return fmt.Errorf("failed to parse save file: %w", err)
	}

	// Version check (for future compatibility)
	if save.Version != "1.0" {
		return fmt.Errorf("incompatible save file version: %s (expected 1.0)", save.Version)
	}

	// Apply state to current game
	g.deserializeState(save.GameState)

	return nil
}

// deserializeState applies a saved game state to the current game
func (g *GameV2) deserializeState(state GameState) {
	g.Location = state.Location
	g.Score = state.Score
	g.Moves = state.Moves
	g.GameOver = state.GameOver
	g.Won = state.Won

	// Restore player state
	g.Player.Inventory = state.PlayerState.Inventory
	g.Player.Health = state.PlayerState.Health
	g.Player.MaxWeight = state.PlayerState.MaxWeight

	// Restore flags
	g.Flags = make(map[string]bool)
	for k, v := range state.Flags {
		g.Flags[k] = v
	}

	// Restore item states
	for id, itemState := range state.ItemStates {
		if item, ok := g.Items[id]; ok {
			item.Location = itemState.Location
			item.Flags = itemState.Flags
			item.Fuel = itemState.Fuel
			item.GlowLevel = itemState.GlowLevel
		}
	}

	// Restore NPC states
	for id, npcState := range state.NPCStates {
		if npc, ok := g.NPCs[id]; ok {
			npc.Location = npcState.Location
			npc.Flags = npcState.Flags
			npc.Strength = npcState.Strength
			npc.Weapon = npcState.Weapon
			npc.Inventory = npcState.Inventory
			npc.Hostile = npcState.Hostile
		}
	}
}

// ListSaves returns a list of available save files
func ListSaves() ([]string, error) {
	saveDir, err := getSaveDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(saveDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read save directory: %w", err)
	}

	var saves []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			saves = append(saves, entry.Name())
		}
	}

	return saves, nil
}

// GetSavePath returns the full path to a save file
func GetSavePath(filename string) (string, error) {
	if filepath.Ext(filename) != ".json" {
		filename += ".json"
	}

	saveDir, err := getSaveDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(saveDir, filename), nil
}
