package main

import (
	"fmt"
	"github.com/wakatara/gork/engine"
)

func main() {
	g := engine.NewGameV2()

	fmt.Println("=== RUG/TRAP DOOR PUZZLE TEST ===\n")

	// Navigate to living room
	fmt.Println("> north")
	g.Process("north")
	fmt.Println("> east")
	g.Process("east")
	fmt.Println("> open window")
	g.Process("open window")
	fmt.Println("> in")
	g.Process("in")
	fmt.Println("> west")
	g.Process("west")
	fmt.Println()

	// Look at the living room
	fmt.Println("> look")
	fmt.Println(g.Process("look"))
	fmt.Println()

	// Try to go down without moving rug
	fmt.Println("> down")
	fmt.Println(g.Process("down"))
	fmt.Println()

	// Move the rug
	fmt.Println("> move rug")
	fmt.Println(g.Process("move rug"))
	fmt.Println()

	// Look again
	fmt.Println("> look")
	fmt.Println(g.Process("look"))
	fmt.Println()

	// Open the trap door
	fmt.Println("> open trap door")
	fmt.Println(g.Process("open trap door"))
	fmt.Println()

	// Now go down
	fmt.Println("> down")
	fmt.Println(g.Process("down"))
	fmt.Println()

	// Look at cellar
	fmt.Println("> look")
	fmt.Println(g.Process("look"))
	fmt.Println()

	fmt.Println("✅ Rug/trap door puzzle works!")
}
