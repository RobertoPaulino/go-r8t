package main

import (
	"go-r8t/cpu"

	"github.com/nsf/termbox-go"
)

// Display buffer with fade-out state to reduce flickering
var displayBuffer [64 * 32]int

// TerminalDisplay is responsible for rendering the CHIP-8 display in the terminal
func TerminalDisplay(cpu *cpu.CPU) {
	// Clear screen
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Build the display string
	display := cpu.GetDisplay()

	// Calculate display dimensions
	termWidth, _ := termbox.Size()
	startX := (termWidth - 128) / 2
	if startX < 0 {
		startX = 0
	}

	// Render border and display
	renderBorder(startX, 0, 128+2, 32+2)

	// Render the CHIP-8 display with phosphor effect
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			index := y*64 + x

			// Update buffer with new pixel state
			if display[index] == 1 {
				displayBuffer[index] = 3 // Full brightness
			} else if displayBuffer[index] > 0 {
				displayBuffer[index]-- // Fade out
			}

			// Set character and color based on buffer state
			if displayBuffer[index] > 0 {
				var color termbox.Attribute
				switch displayBuffer[index] {
				case 3:
					color = termbox.ColorWhite // Full brightness
				case 2:
					color = termbox.ColorWhite // Medium brightness
				case 1:
					color = termbox.ColorDarkGray // Dim
				}

				// Set two characters for each pixel (for better aspect ratio)
				termbox.SetCell(startX+1+x*2, y+1, '█', color, termbox.ColorDefault)
				termbox.SetCell(startX+1+x*2+1, y+1, '█', color, termbox.ColorDefault)
			}
		}
	}

	// Render current ROM information
	renderROMInfo(startX, 34)

	// Force update
	termbox.Flush()
}

// InitializeTerminal initializes the terminal UI
func InitializeTerminal() error {
	return termbox.Init()
}

// CloseTerminal closes the terminal UI
func CloseTerminal() {
	termbox.Close()
}

// renderBorder draws a border around the display
func renderBorder(x, y, width, height int) {
	// Draw top and bottom borders
	for i := 0; i < width; i++ {
		termbox.SetCell(x+i, y, '─', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x+i, y+height-1, '─', termbox.ColorWhite, termbox.ColorDefault)
	}

	// Draw left and right borders
	for i := 0; i < height; i++ {
		termbox.SetCell(x, y+i, '│', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x+width-1, y+i, '│', termbox.ColorWhite, termbox.ColorDefault)
	}

	// Draw corners
	termbox.SetCell(x, y, '┌', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(x+width-1, y, '┐', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(x, y+height-1, '└', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(x+width-1, y+height-1, '┘', termbox.ColorWhite, termbox.ColorDefault)
}

// Current ROM name
var currentROM = "No ROM loaded"

// SetCurrentROM sets the current ROM name
func SetCurrentROM(name string) {
	currentROM = name
}

// renderROMInfo renders information about the current ROM
func renderROMInfo(x, y int) {
	drawString(x, y, "ROM: "+currentROM+" (Press ESC to exit)", termbox.ColorWhite)
}

// drawString draws a string at the specified position
func drawString(x, y int, str string, color termbox.Attribute) {
	for i, c := range str {
		termbox.SetCell(x+i, y, c, color, termbox.ColorDefault)
	}
}
