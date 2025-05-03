# CHIP-8 Emulator

A CHIP-8 emulator written in Go, with both desktop and web browser support.

## Features

- Full CHIP-8 instruction set implementation
- 64x32 pixel display
- 16-key hexadecimal keypad input
- Sound support
- ROM loading from file
- WebAssembly support for browser-based emulation

## Running the Desktop Version

```bash
go run main.go [rom_file]
```

## Running the Web Version

1. Build the WebAssembly version:
```bash
./build.sh
```

2. Start a local web server:
```bash
python3 -m http.server
```

3. Open your browser and navigate to:
```
http://localhost:8000
```

4. Use the file input to load a CHIP-8 ROM file.

## Keyboard Controls

The CHIP-8 keypad is mapped to the following keyboard keys:

```
1 2 3 C
4 5 6 D
7 8 9 E
A 0 B F
```

## Building from Source

```bash
go build
```

## License

MIT License - see LICENSE file for details.
