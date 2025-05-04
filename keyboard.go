package main

import (
	"go-r8t/cpu"

	"time"

	"github.com/nsf/termbox-go"
)

// Map of terminal keys to CHIP-8 keypad
var keyMap = map[termbox.Key]byte{
	termbox.KeyF1:  0x1, // 1
	termbox.KeyF2:  0x2, // 2
	termbox.KeyF3:  0x3, // 3
	termbox.KeyF4:  0xC, // 4 -> C
	termbox.KeyF5:  0x4, // Q -> 4
	termbox.KeyF6:  0x5, // W -> 5
	termbox.KeyF7:  0x6, // E -> 6
	termbox.KeyF8:  0xD, // R -> D
	termbox.KeyF9:  0x7, // A -> 7
	termbox.KeyF10: 0x8, // S -> 8
	termbox.KeyF11: 0x9, // D -> 9
	termbox.KeyF12: 0xE, // F -> E
}

// Map of character runes to CHIP-8 keypad
var runeKeyMap = map[rune]byte{
	'1': 0x1, // 1
	'2': 0x2, // 2
	'3': 0x3, // 3
	'4': 0xC, // 4 -> C
	'q': 0x4, // q -> 4
	'w': 0x5, // w -> 5
	'e': 0x6, // e -> 6
	'r': 0xD, // r -> D
	'a': 0x7, // a -> 7
	's': 0x8, // s -> 8
	'd': 0x9, // d -> 9
	'f': 0xE, // f -> E
	'z': 0xA, // z -> A
	'x': 0x0, // x -> 0
	'c': 0xB, // c -> B
	'v': 0xF, // v -> F
}

// StartKeyboardInput initializes keyboard input handling in a separate goroutine
func StartKeyboardInput(chip8 *cpu.CPU, quit chan struct{}) {
	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				ev := termbox.PollEvent()

				switch ev.Type {
				case termbox.EventKey:
					if ev.Key == termbox.KeyEsc {
						// Send signal to quit
						close(quit)
						return
					}

					// Handle key press
					if chipKey, ok := keyMap[ev.Key]; ok {
						chip8.SetKey(chipKey, true)

						// Create a goroutine to handle auto-release after a short delay
						go func(key byte) {
							time.Sleep(100 * time.Millisecond)
							chip8.SetKey(key, false)
						}(chipKey)
					} else if chipKey, ok := runeKeyMap[ev.Ch]; ok {
						chip8.SetKey(chipKey, true)

						// Create a goroutine to handle auto-release after a short delay
						go func(key byte) {
							time.Sleep(100 * time.Millisecond)
							chip8.SetKey(key, false)
						}(chipKey)
					}
				case termbox.EventError:
					panic(ev.Err)
				}
			}
		}
	}()
}
