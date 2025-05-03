//go:build js && wasm
// +build js,wasm

package main

import (
	"go-r8t/cpu"
	"image/color"
	"log"
	"syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game represents the main game state
type Game struct {
	cpu *cpu.CPU
}

// NewGame creates a new game instance
func NewGame() *Game {
	g := &Game{
		cpu: cpu.NewCPU(),
	}
	return g
}

// Update updates the game state
func (g *Game) Update() error {
	// Handle input
	g.handleInput()

	// Run one CPU cycle
	opcode := (uint16(g.cpu.Memory[g.cpu.PC]) << 8) | uint16(g.cpu.Memory[g.cpu.PC+1])
	g.cpu.ExecuteInstruction(opcode)

	// Update timers at 60Hz
	g.cpu.UpdateTimers()

	return nil
}

// Draw draws the game screen
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Get screen dimensions
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Calculate pixel size to maintain aspect ratio
	pixelWidth := float64(screenWidth) / 64
	pixelHeight := float64(screenHeight) / 32

	// Draw CHIP-8 display
	pixelsDrawn := 0
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if g.cpu.Display[y*64+x] == 1 {
				ebitenutil.DrawRect(
					screen,
					float64(x)*pixelWidth,
					float64(y)*pixelHeight,
					pixelWidth,
					pixelHeight,
					color.RGBA{255, 255, 255, 255},
				)
				pixelsDrawn++
			}
		}
	}

	// Log display state for debugging
	if pixelsDrawn > 0 {
		log.Printf("Display: %d pixels drawn, screen size: %dx%d, pixel size: %.2fx%.2f",
			pixelsDrawn, screenWidth, screenHeight, pixelWidth, pixelHeight)
	}
}

// Layout implements ebiten.Game's Layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 256, 128 // Very small CHIP-8 display
}

// handleInput updates the CPU's key state based on keyboard input
func (g *Game) handleInput() {
	// CHIP-8 keypad mapping to keyboard
	keyMap := map[ebiten.Key]uint8{
		ebiten.Key1: 0x1,
		ebiten.Key2: 0x2,
		ebiten.Key3: 0x3,
		ebiten.Key4: 0xC,
		ebiten.KeyQ: 0x4,
		ebiten.KeyW: 0x5,
		ebiten.KeyE: 0x6,
		ebiten.KeyR: 0xD,
		ebiten.KeyA: 0x7,
		ebiten.KeyS: 0x8,
		ebiten.KeyD: 0x9,
		ebiten.KeyF: 0xE,
		ebiten.KeyZ: 0xA,
		ebiten.KeyX: 0x0,
		ebiten.KeyC: 0xB,
		ebiten.KeyV: 0xF,
	}

	// Update key states
	for key, value := range keyMap {
		g.cpu.SetKey(value, ebiten.IsKeyPressed(key))
	}
}

var game *Game

func main() {
	game = NewGame()

	// Set up WebAssembly-specific initialization
	js.Global().Set("loadROM", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Error: Expected 1 argument"
		}
		romData := args[0]
		if romData.Type() != js.TypeObject {
			return "Error: Expected Uint8Array"
		}

		// Copy ROM data to memory
		romBytes := make([]byte, romData.Length())
		js.CopyBytesToGo(romBytes, romData)
		game.cpu.LoadProgram(romBytes)
		return "ROM loaded successfully"
	}))

	// Configure the window
	ebiten.SetWindowTitle("CHIP-8 Emulator")
	ebiten.SetWindowSize(640, 320)
	ebiten.SetWindowResizable(true)
	ebiten.SetMaxTPS(60) // Limit to 60 frames per second

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
