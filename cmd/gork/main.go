package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/wakatara/gork/engine"
	"github.com/wakatara/gork/ui"
)

// version is set by GoReleaser at build time
var version = "dev"

func main() {
	// Handle flags
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-version", "-v":
			fmt.Printf("gork version %s\n", version)
			fmt.Println("ZORK I: The Great Underground Empire")
			fmt.Println("Go Edition - https://github.com/wakatara/gork")
			return
		case "--help", "-help", "-h":
			fmt.Println("GORK - ZORK I: The Great Underground Empire (Go Edition)")
			fmt.Println()
			fmt.Println("Usage:")
			fmt.Println("  gork              Start the game")
			fmt.Println("  gork --version    Show version information")
			fmt.Println("  gork --help       Show this help message")
			fmt.Println()
			fmt.Println("In-game commands:")
			fmt.Println("  Type 'help' in the game for available commands")
			fmt.Println("  Type 'quit' to exit the game")
			return
		}
	}

	// Display title
	ui.PrintTitle()

	// Create new game with refactored types
	game := engine.NewGameV2(version)

	// Display initial message
	ui.PrintSlow(game.GetInitialMessage())
	fmt.Println()

	// Main game loop (REPL)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Print prompt
		ui.PrintPrompt()

		// Read input
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if input == "" {
			continue
		}

		// Process command
		output := game.Process(input)

		// Check for special clear screen command
		if output == "<<CLEAR_SCREEN>>" {
			// Clear screen and reprint title + current location
			locationDesc := game.Process("look")
			ui.PrintTitleAndLocation(locationDesc)
			continue
		}

		// Display result
		if output != "" {
			ui.PrintSlow(output)
			fmt.Println()
		}

		// Check if game is over
		if game.GameOver {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}
