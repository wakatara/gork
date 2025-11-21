package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/wakatara/gork/engine"
	"github.com/wakatara/gork/ui"
)

func main() {
	// Display title
	ui.PrintTitle()

	// Create new game with refactored types
	game := engine.NewGameV2()

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
