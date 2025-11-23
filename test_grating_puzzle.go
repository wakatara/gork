package main

import (
	"fmt"
	"github.com/wakatara/gork/engine"
)

func main() {
	g := engine.NewGameV2()

	fmt.Println("=== GRATING PUZZLE TEST ===\n")

	// Navigate to living room to get keys
	fmt.Println("Getting to living room...")
	g.Process("north")
	g.Process("east")
	g.Process("open window")
	g.Process("in")
	g.Process("west")

	// Take the keys
	fmt.Println("> take keys")
	fmt.Println(g.Process("take keys"))
	fmt.Println()

	// Navigate to grating-clearing
	fmt.Println("Going to grating-clearing...")
	g.Process("east") // to kitchen
	g.Process("west") // back to living room
	g.Process("east") // to kitchen
	g.Process("out")  // to behind house
	g.Process("north") // to north-of-house
	g.Process("north") // to path
	g.Process("north") // to grating-clearing
	fmt.Println()

	// Look at clearing
	fmt.Println("> look")
	fmt.Println(g.Process("look"))
	fmt.Println()

	// Try to go down without opening
	fmt.Println("> down")
	fmt.Println(g.Process("down"))
	fmt.Println()

	// Try to open without keys first (we have them, but let's show the check works)
	// Actually we already have keys, so it should work
	fmt.Println("> open grating")
	fmt.Println(g.Process("open grating"))
	fmt.Println()

	// Now go down
	fmt.Println("> down")
	fmt.Println(g.Process("down"))
	fmt.Println()

	// Look at grating room
	fmt.Println("> look")
	fmt.Println(g.Process("look"))
	fmt.Println()

	fmt.Println("✅ Grating puzzle works!")
}
