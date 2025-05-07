# Go-R8T: A CHIP-8 Emulator in Go

Go-R8T is a simple yet powerful CHIP-8 emulator written in Go. It provides a terminal-based interface to run classic CHIP-8 games and applications.

## What is CHIP-8?

CHIP-8 is an interpreted programming language developed in the mid-1970s, designed to allow video games to be more easily programmed for early microcomputers. It was initially used on the COSMAC VIP and Telmac 1800 8-bit microcomputers. CHIP-8 programs are run on a virtual machine which interprets the instructions.

## Features

- **Complete CHIP-8 instruction set**: Supports all original CHIP-8 instructions
- **Terminal display**: Play games in your terminal with a phosphor effect for reduced flickering
- **Simple ROM loading**: Easily load and run CHIP-8 ROMs
- **Full keyboard input**: Maps standard keyboard keys to the CHIP-8 hexadecimal keypad

## Installation

### Prerequisites

- Go 1.18 or higher
- termbox-go (for terminal interface)

### Building from Source

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/go-r8t.git
   cd go-r8t
   ```

2. Build the emulator:
   ```
   go build
   ```

## Usage

Run the emulator with a ROM file:

```
./go-r8t path/to/rom.ch8
```

### Controls

The CHIP-8 uses a 16-key hexadecimal keypad mapped to the following keys:

```
CHIP-8 Keypad    Keyboard
+-+-+-+-+        +-+-+-+-+
|1|2|3|C|        |1|2|3|4|
+-+-+-+-+        +-+-+-+-+
|4|5|6|D|        |Q|W|E|R|
+-+-+-+-+   =>   +-+-+-+-+
|7|8|9|E|        |A|S|D|F|
+-+-+-+-+        +-+-+-+-+
|A|0|B|F|        |Z|X|C|V|
+-+-+-+-+        +-+-+-+-+
```

Function keys are also mapped as alternatives:

```
F1-F4:   Keys 1, 2, 3, C
F5-F8:   Keys 4, 5, 6, D
F9-F12:  Keys 7, 8, 9, E
```

Press `ESC` to exit the emulator.

## Architecture

The emulator consists of several key components:

- **CPU**: Implements the CHIP-8 instruction set and manages system state
- **Display**: Renders the 64×32 pixel monochrome display with phosphor effect to reduce flickering
- **Keyboard**: Handles input from the 16-key hexadecimal keypad
- **Memory**: Manages the 4KB of RAM

## Technical Details

- 4KB of memory
- 16 8-bit registers (V0-VF)
- 16-level stack for subroutine calls
- 64×32 pixel monochrome display
- 16-key hexadecimal keypad
- Two timers (delay and sound) that decrement at 60Hz

## Finding ROMs

CHIP-8 ROMs are widely available online. Some classic games include:
- Pong
- Space Invaders
- Tetris
- Breakout

Make sure the files have a `.ch8` extension or are in binary format.

## Troubleshooting

### Display Issues
- If the display appears flickery, try running in a larger terminal window
- The phosphor effect is designed to reduce flickering, but may not work perfectly in all terminal emulators

### ROM Compatibility
- Some ROMs may require specific timing adjustments
- Super CHIP-8 or CHIP-48 ROMs might have limited compatibility

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the GNU General Public License (GPL) Version 2, June 1991 - see the LICENSE file for details.

## Acknowledgments

- Thanks to the CHIP-8 developer community for documentation and test ROMs
- Special thanks to the authors of termbox-go
