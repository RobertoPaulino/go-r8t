package main

import (
	"fmt"
	"math/rand"
)

// CPU represents the CHIP-8 virtual machine state.
// It contains all the registers, memory, and state needed to execute CHIP-8 programs.
type CPU struct {
	PC            uint16        // Program Counter - points to the current instruction in memory
	memory        [4096]byte    // 4KB of memory (0x000-0x1FF: System memory, 0x200-0xFFF: Program memory)
	V             [16]byte      // 16 general-purpose registers (V0-VF)
	stack         [16]uint16    // Stack for subroutine calls (16 levels deep)
	I             uint16        // Index register - used for memory operations and sprite drawing
	SP            uint8         // Stack Pointer - points to the current stack level
	delayTimer    uint8         // Delay timer - decrements at 60Hz when non-zero
	soundTimer    uint8         // Sound timer - decrements at 60Hz when non-zero, beeps when non-zero
	keys          [16]bool      // State of the 16-key hexadecimal keypad (0x0-0xF)
	display       [64 * 32]byte // Display memory (64x32 pixels, 1 bit per pixel)
	currentOpcode uint16        // The current instruction being executed
}

func main() {
	cpu := &CPU{}

	testProgram := []byte{0x00, 0xE0, 0x12, 0x00}

	cpu.loadProgram(testProgram)

}

// loadProgram loads a CHIP-8 program into memory starting at address 0x200.
// This is the standard starting address for CHIP-8 programs.
// Parameters:
//   - program: The byte slice containing the CHIP-8 program to load
func (cpu *CPU) loadProgram(program []byte) {
	for i := range program {
		cpu.memory[0x200+i] = program[i]
	}
}

// executeInstruction decodes and executes a single CHIP-8 instruction.
// The instruction is processed based on its opcode pattern, and the appropriate
// operation is performed on the CPU state.
// Parameters:
//   - instruction: The 16-bit instruction to execute
func (cpu *CPU) executeInstruction(instruction uint16) {
	switch instruction & 0xF000 {
	case 0x0000:
		if instruction == 0x00E0 {
			// 00E0 - CLS
			// Clear the display
			cpu.clearScreen()
			return
		} else if instruction == 0x00EE {
			// 00EE - RET
			// Return from a subroutine
			// The interpreter sets the program counter to the address at the top of the stack,
			// then subtracts 1 from the stack pointer.
			cpu.returnFromSubroutine()
			return
		}
	case 0x1000:
		// 1NNN - JP addr
		// Jump to location NNN
		// The interpreter sets the program counter to NNN.
		nextAddress := instruction & 0x0FFF
		cpu.PC = nextAddress
	case 0x2000:
		// 2NNN - CALL addr
		// Call subroutine at NNN
		// The interpreter increments the stack pointer, then puts the current PC on the top of the stack.
		// The PC is then set to NNN.
		target := instruction & 0x0FFF
		cpu.stack[cpu.SP] = cpu.PC
		cpu.SP++
		cpu.PC = target
	case 0x3000:
		// 3XNN - SE Vx, byte
		// Skip next instruction if Vx = NN
		// The interpreter compares register Vx to NN, and if they are equal,
		// increments the program counter by 2.
		x := (instruction & 0x0F00) >> 8
		nn := byte(instruction & 0x00FF)
		if cpu.V[x] == nn {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	case 0x4000:
		// 4XNN - SNE Vx, byte
		// Skip next instruction if Vx != NN
		// The interpreter compares register Vx to NN, and if they are not equal,
		// increments the program counter by 2.
		x := (instruction & 0x0F00) >> 8
		nn := byte(instruction & 0x00FF)
		if cpu.V[x] != nn {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	case 0x5000:
		// 5XY0 - SE Vx, Vy
		// Skip next instruction if Vx = Vy
		// The interpreter compares register Vx to register Vy, and if they are equal,
		// increments the program counter by 2.
		x := (instruction & 0x0F00) >> 8
		y := (instruction & 0x00F0) >> 4
		if cpu.V[x] == cpu.V[y] {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	case 0x6000:
		// 6XNN - LD Vx, byte
		// Set Vx = NN
		// The interpreter puts the value NN into register Vx.
		x := (instruction & 0x0F00) >> 8
		cpu.V[x] = byte(instruction & 0x00FF)
	case 0x7000:
		// 7XNN - ADD Vx, byte
		// Set Vx = Vx + NN
		// Adds the value NN to the value of register Vx, then stores the result in Vx.
		x := (instruction & 0x0F00) >> 8
		cpu.V[x] += byte(instruction & 0x00FF)
	case 0x8000:
		// 8XYN - Various arithmetic and logical operations
		// The interpreter performs the specified operation on Vx and Vy, then stores the result in Vx.
		x := (instruction & 0x0F00) >> 8
		y := (instruction & 0x00F0) >> 4
		sel := instruction & 0x000F
		cpu.performArithmeticOperation(x, y, sel)
	case 0x9000:
		// 9XY0 - SNE Vx, Vy
		// Skip next instruction if Vx != Vy
		// The interpreter compares register Vx to register Vy, and if they are not equal,
		// increments the program counter by 2.
		x := (instruction & 0x0F00) >> 8
		y := (instruction & 0x00F0) >> 4
		if cpu.V[x] != cpu.V[y] {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	case 0xA000:
		// ANNN - LD I, addr
		// Set I = NNN
		// The value of register I is set to NNN.
		cpu.I = instruction & 0x0FFF
	case 0xB000:
		// BNNN - JP V0, addr
		// Jump to location NNN + V0
		// The program counter is set to NNN plus the value of V0.
		nextAddress := (instruction & 0x0FFF) + uint16(cpu.V[0])
		cpu.PC = nextAddress
	case 0xC000:
		// CXNN - RND Vx, byte
		// Set Vx = random byte AND NN
		// The interpreter generates a random number from 0 to 255, which is then ANDed with NN.
		// The results are stored in Vx.
		x := (instruction & 0x0F00) >> 8
		nn := byte(instruction & 0x00FF)
		random := byte(rand.Intn(256))
		cpu.V[x] = random & nn
	case 0xD000:
		// DXYN - DRW Vx, Vy, nibble
		// Display N-byte sprite starting at memory location I at (Vx, Vy)
		// The interpreter reads N bytes from memory, starting at the address stored in I.
		// These bytes are then displayed as sprites on screen at coordinates (Vx, Vy).
		// Sprites are XORed onto the existing screen. If this causes any pixels to be erased,
		// VF is set to 1, otherwise it is set to 0.
		x := (instruction & 0x0F00) >> 8
		y := (instruction & 0x00F0) >> 4
		size := instruction & 0x000F
		cpu.drawSprite(x, y, size)
		cpu.PC += 2
	case 0xE000:
		// EX9E - SKP Vx
		// Skip next instruction if key with the value of Vx is pressed
		// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the down position,
		// PC is increased by 2.
		// EXA1 - SKNP Vx
		// Skip next instruction if key with the value of Vx is not pressed
		// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position,
		// PC is increased by 2.
		x := (instruction & 0x0F00) >> 8
		opcodeLow := byte(instruction & 0x00FF)
		keyState := cpu.keys[cpu.V[x]]
		if (keyState == true) && (opcodeLow == 0x9E) {
			cpu.PC += 4
		} else if (keyState == false) && (opcodeLow == 0xA1) {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	case 0xF000:
		// Various operations starting with F
		// These are handled in perform0xFOperation
		x := (instruction & 0x0F00) >> 8
		sel := instruction & 0x00FF
		cpu.perform0xFOperation(x, sel)
	default:
		// Unknown opcode
		fmt.Printf("Unknown opcode: 0x%04X\n", instruction)
	}
}

// drawSprite draws a sprite at coordinates (VX, VY) with N bytes of sprite data.
// The sprite is drawn using XOR logic, and VF is set to 1 if any pixels are flipped from set to unset.
// Parameters:
//   - x: The register index containing the X coordinate
//   - y: The register index containing the Y coordinate
//   - size: The number of bytes of sprite data to draw (height of sprite)
func (cpu *CPU) drawSprite(x, y, size uint16) {
	xCord := uint16(cpu.V[x]) // Cast to uint16 for correct math
	yCord := uint16(cpu.V[y]) // Cast to uint16 for correct math

	cpu.V[0xF] = 0 // Reset collision flag

	// Loop over the sprite rows (from 0 to size-1)
	for row := uint16(0); row < size; row++ {
		spriteByte := cpu.memory[cpu.I+row] // Get the byte for the current row of the sprite

		// Loop through each bit (8 bits = 1 byte)
		for bit := uint16(0); bit < 8; bit++ {
			// Calculate screen x, y position
			screenX := (xCord + bit) % 64 // Wrap around horizontally if needed
			screenY := (yCord + row) % 32 // Wrap around vertically if needed

			// Get the index in the display array
			displayIndex := screenY*64 + screenX

			// Check if the pixel should be turned off (collision)
			if (spriteByte & (0x80 >> bit)) != 0 { // Check if the current bit is 1
				if cpu.display[displayIndex] == 1 {
					cpu.V[0xF] = 1 // Set VF flag if there's a collision (pixel was already 1)
				}
				// Flip the pixel on the screen using XOR
				cpu.display[displayIndex] ^= 1
			}
		}
	}
}

// perform0xFOperation handles all instructions starting with 0xF.
// These instructions are typically used for memory and timer operations.
// Parameters:
//   - x: The register index specified in the instruction
//   - sel: The selector byte that determines which 0xF operation to perform
func (cpu *CPU) perform0xFOperation(x, sel uint16) {
	switch sel {
	case 0x07:
		// FX07 - LD Vx, DT
		// Set Vx = delay timer value
		// The value of DT is placed into Vx.
		cpu.V[x] = cpu.delayTimer
	case 0x0A:
		// FX0A - LD Vx, K
		// Wait for a key press, store the value of the key in Vx
		// All execution stops until a key is pressed, then the value of that key is stored in Vx.
		keyPressed := false
		for i, key := range cpu.keys {
			if key == true {
				cpu.V[x] = byte(i)
				cpu.PC += 2
				keyPressed = true
				break
			}
		}
		if !keyPressed {
			return
		}
	case 0x15:
		// FX15 - LD DT, Vx
		// Set delay timer = Vx
		// DT is set equal to the value of Vx.
		cpu.delayTimer = cpu.V[x]
	case 0x18:
		// FX18 - LD ST, Vx
		// Set sound timer = Vx
		// ST is set equal to the value of Vx.
		cpu.soundTimer = cpu.V[x]
	case 0x1E:
		// FX1E - ADD I, Vx
		// Set I = I + Vx
		// The values of I and Vx are added, and the results are stored in I.
		cpu.I += uint16(cpu.V[x])
	case 0x29:
		// FX29 - LD F, Vx
		// Set I = location of sprite for digit Vx
		// The value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx.
		cpu.I = uint16(cpu.V[x]) * 5
	case 0x33:
		// FX33 - LD B, Vx
		// Store BCD representation of Vx in memory locations I, I+1, and I+2
		// The interpreter takes the decimal value of Vx, and places the hundreds digit in memory at location in I,
		// the tens digit at location I+1, and the ones digit at location I+2.
		cpu.memory[cpu.I] = (cpu.V[x] / 100)
		cpu.memory[cpu.I+1] = (cpu.V[x] % 100) / 10
		cpu.memory[cpu.I+2] = (cpu.V[x] % 10)
	default:
		// Unknown opcode
		fmt.Printf("Unknown 0xF opcode: 0x%02X\n", sel)
	}
}

// performArithmeticOperation handles all arithmetic and bitwise operations
// specified by the 0x8XXX instructions.
// Parameters:
//   - x: The first register index
//   - y: The second register index
//   - sel: The selector nibble that determines which operation to perform
func (cpu *CPU) performArithmeticOperation(x, y, sel uint16) {
	switch sel {
	case 0x0:
		// 8XY0 - LD Vx, Vy
		// Set Vx = Vy
		// Stores the value of register Vy in register Vx.
		cpu.V[x] = cpu.V[y]
	case 0x1:
		// 8XY1 - OR Vx, Vy
		// Set Vx = Vx OR Vy
		// Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx.
		cpu.V[x] |= cpu.V[y]
	case 0x2:
		// 8XY2 - AND Vx, Vy
		// Set Vx = Vx AND Vy
		// Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx.
		cpu.V[x] &= cpu.V[y]
	case 0x3:
		// 8XY3 - XOR Vx, Vy
		// Set Vx = Vx XOR Vy
		// Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx.
		cpu.V[x] ^= cpu.V[y]
	case 0x4:
		// 8XY4 - ADD Vx, Vy
		// Set Vx = Vx + Vy, set VF = carry
		// The values of Vx and Vy are added together. If the result is greater than 8 bits (i.e., > 255),
		// VF is set to 1, otherwise 0. Only the lowest 8 bits of the result are kept, and stored in Vx.
		result := uint16(cpu.V[x]) + uint16(cpu.V[y])
		if result > 0xFF {
			cpu.V[0xF] = 1
		} else {
			cpu.V[0xF] = 0
		}
		cpu.V[x] = byte(result)
	case 0x5:
		// 8XY5 - SUB Vx, Vy
		// Set Vx = Vx - Vy, set VF = NOT borrow
		// If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.
		if uint16(cpu.V[x]) > uint16(cpu.V[y]) {
			cpu.V[0xF] = 1
		} else {
			cpu.V[0xF] = 0
		}
		cpu.V[x] -= cpu.V[y]
	case 0x6:
		// 8XY6 - SHR Vx {, Vy}
		// Set Vx = Vx SHR 1
		// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
		cpu.V[0xF] = cpu.V[x] & 0x1
		cpu.V[x] >>= 1
	case 0x7:
		// 8XY7 - SUBN Vx, Vy
		// Set Vx = Vy - Vx, set VF = NOT borrow
		// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
		if uint16(cpu.V[y]) > uint16(cpu.V[x]) {
			cpu.V[0xF] = 1
		} else {
			cpu.V[0xF] = 0
		}
		cpu.V[x] = cpu.V[y] - cpu.V[x]
	case 0xE:
		// 8XYE - SHL Vx {, Vy}
		// Set Vx = Vx SHL 1
		// If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.
		cpu.V[0xF] = (cpu.V[x] >> 7) & 0x1
		cpu.V[x] <<= 1
	default:
		// Unknown arithmetic operation
		fmt.Printf("Unknown arithmetic operation: 0x%X\n", sel)
	}
}

func (cpu *CPU) clearScreen() {
	for i := range cpu.display {
		cpu.display[i] = 0
	}
}

func (cpu *CPU) returnFromSubroutine() {
	cpu.SP -= 1
	cpu.PC = cpu.stack[cpu.SP]
}
