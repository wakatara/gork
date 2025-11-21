package ui

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// ANSI color codes for retro terminal look
const (
	// Classic amber monitor color
	ColorAmber = "\033[38;5;214m"
	// Classic green monitor color
	ColorGreen = "\033[38;5;46m"
	// White/gray for normal text
	ColorDefault = "\033[0m"
	ColorBold    = "\033[1m"
	ColorDim     = "\033[2m"

	// Current theme (change to ColorGreen for green CRT feel)
	ThemeColor = ColorAmber
)

var (
	// EnableTypingEffect controls whether text appears character-by-character
	EnableTypingEffect = false // Set to true for classic typing effect

	// TypingSpeed is the delay between characters (in milliseconds)
	TypingSpeed = 20
)

// PrintTitle displays the game title with ASCII art
func PrintTitle() {
	title := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║   ▄████  ▒█████   ██▀███   ██ ▄█▀                          ║
║  ██▒ ▀█▒▒██▒  ██▒▓██ ▒ ██▒ ██▄█▒                           ║
║ ▒██░▄▄▄░▒██░  ██▒▓██ ░▄█ ▒▓███▄░                           ║
║ ░▓█  ██▓▒██   ██░▒██▀▀█▄  ▓██ █▄                           ║
║ ░▒▓███▀▒░ ████▓▒░░██▓ ▒██▒▒██▒ █▄                          ║
║  ░▒   ▒ ░ ▒░▒░▒░ ░ ▒▓ ░▒▓░▒ ▒▒ ▓▒                          ║
║   ░   ░   ░ ▒ ▒░   ░▒ ░ ▒░░ ░▒ ▒░                          ║
║ ░ ░   ░ ░ ░ ░ ▒    ░░   ░ ░ ░░ ░                           ║
║       ░     ░ ░     ░     ░  ░                              ║
║                                                              ║
║           The Great Underground Empire (Go Edition)         ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	// Print title in theme color
	fmt.Print(ThemeColor + title + ColorDefault)
	time.Sleep(500 * time.Millisecond)
	fmt.Println()
}

// PrintPrompt displays the input prompt
func PrintPrompt() {
	fmt.Print(ThemeColor + "> " + ColorDefault)
}

// PrintSlow prints text with optional typing effect
func PrintSlow(text string) {
	if !EnableTypingEffect {
		fmt.Print(ThemeColor + text + ColorDefault)
		return
	}

	// Character-by-character typing effect
	for _, char := range text {
		fmt.Print(ThemeColor + string(char) + ColorDefault)

		// Add slight random variation to typing speed for authenticity
		delay := time.Duration(TypingSpeed + rand.Intn(10))
		time.Sleep(delay * time.Millisecond)

		// Slightly longer pause after punctuation
		if char == '.' || char == '!' || char == '?' {
			time.Sleep(200 * time.Millisecond)
		} else if char == ',' || char == ';' {
			time.Sleep(100 * time.Millisecond)
		} else if char == '\n' {
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// PrintBox prints text in a bordered box
func PrintBox(text string) {
	lines := strings.Split(text, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Top border
	fmt.Print(ThemeColor + "┌")
	for i := 0; i < maxLen+2; i++ {
		fmt.Print("─")
	}
	fmt.Println("┐" + ColorDefault)

	// Content lines
	for _, line := range lines {
		padding := maxLen - len(line)
		fmt.Print(ThemeColor + "│ " + ColorDefault)
		fmt.Print(line)
		for i := 0; i < padding; i++ {
			fmt.Print(" ")
		}
		fmt.Println(ThemeColor + " │" + ColorDefault)
	}

	// Bottom border
	fmt.Print(ThemeColor + "└")
	for i := 0; i < maxLen+2; i++ {
		fmt.Print("─")
	}
	fmt.Println("┘" + ColorDefault)
}

// EnableRetroMode turns on all retro features
func EnableRetroMode() {
	EnableTypingEffect = true
	TypingSpeed = 30 // Slightly slower for that 1980s modem feel
}

// PrintDeath displays a dramatic death message
func PrintDeath(message string) {
	fmt.Println()
	fmt.Print(ColorBold)
	PrintBox("☠ " + message + " ☠")
	fmt.Print(ColorDefault)
	fmt.Println()
}

// PrintTreasure displays a treasure acquisition message
func PrintTreasure(message string) {
	fmt.Println()
	fmt.Print(ColorBold + ThemeColor)
	PrintBox("✦ " + message + " ✦")
	fmt.Print(ColorDefault)
	fmt.Println()
}

// CheckTerminalSupport checks if the terminal supports ANSI colors
func CheckTerminalSupport() bool {
	term := os.Getenv("TERM")
	return term != "" && term != "dumb"
}

func init() {
	// Seed random for typing effect variation
	rand.Seed(time.Now().UnixNano())

	// Check if we should disable colors
	if !CheckTerminalSupport() {
		// Disable colors on unsupported terminals
		// (In a full implementation, we'd set all colors to empty strings)
	}
}
