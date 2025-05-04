package cpu

import (
	"fmt"
	"math/rand"
)

// Font data for hexadecimal digits 0-F
var fontSet = []byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

// CPU represents the CHIP-8 virtual machine state.
// It contains all the registers, memory, and state needed to execute CHIP-8 programs.
type CPU struct {
	PC            uint16        // Program Counter - points to the current instruction in memory
	Memory        [4096]byte    // 4KB of memory (0x000-0x1FF: System memory, 0x200-0xFFF: Program memory)
	V             [16]byte      // 16 general-purpose registers (V0-VF)
	Stack         [16]uint16    // Stack for subroutine calls (16 levels deep)
	I             uint16        // Index register - used for memory operations and sprite drawing
	SP            uint8         // Stack Pointer - points to the current stack level
	DelayTimer    uint8         // Delay timer - decrements at 60Hz when non-zero
	SoundTimer    uint8         // Sound timer - decrements at 60Hz when non-zero, beeps when non-zero
	Keys          [16]bool      // State of the 16-key hexadecimal keypad (0x0-0xF)
	Display       [64 * 32]byte // Display memory (64x32 pixels, 1 bit per pixel)
	CurrentOpcode uint16        // The current instruction being executed
}

// NewCPU creates and returns a new CPU instance.
func NewCPU() *CPU {
	cpu := &CPU{
		PC: 0x200, // Program counter starts at 0x200
	}
<<<<<<< HEAD
	cpu.ClearScreen() // Clear the display on initialization
=======

	// Load font data into memory at address 0
	for i, fontByte := range fontSet {
		cpu.Memory[i] = fontByte
	}

>>>>>>> 71beedc (Terminal display finished, fixed bug with I register not updating correctly)
	return cpu
}

// LoadProgram loads a CHIP-8 program into memory starting at address 0x200.
// This is the standard starting address for CHIP-8 programs.
// Parameters:
//   - program: The byte slice containing the CHIP-8 program to load
func (cpu *CPU) LoadProgram(program []byte) {
	for i := range program {
		cpu.Memory[0x200+i] = program[i]
	}
}

// ExecuteInstruction decodes and executes a single CHIP-8 instruction.
// The instruction is processed based on its opcode pattern, and the appropriate
// operation is performed on the CPU state.
// Parameters:
//   - instruction: The 16-bit instruction to execute
func (cpu *CPU) ExecuteInstruction(instruction uint16) {
	cpu.CurrentOpcode = instruction
	switch instruction & 0xF000 {
	case 0x0000:
		if instruction == 0x00E0 {
			// 00E0 - CLS
			// Clear the display
			cpu.ClearScreen()
			cpu.PC += 2
		} else if instruction == 0x00EE {
			// 00EE - RET
			// Return from a subroutine
			cpu.ReturnFromSubroutine()
		}
	case 0x1000:
		// 1NNN - JP addr
		// Jump to location NNN
		cpu.PC = instruction & 0x0FFF
	case 0x2000:
		// 2NNN - CALL addr
		// Call subroutine at NNN
		cpu.Stack[cpu.SP] = cpu.PC
		cpu.SP++
		cpu.PC = instruction & 0x0FFF
	case 0x3000:
		// 3XNN - SE Vx, byte
		// Skip next instruction if Vx = NN
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
		x := (instruction & 0x0F00) >> 8
		cpu.V[x] = byte(instruction & 0x00FF)
		cpu.PC += 2
	case 0x7000:
		// 7XNN - ADD Vx, byte
		// Set Vx = Vx + NN
		x := (instruction & 0x0F00) >> 8
		cpu.V[x] += byte(instruction & 0x00FF)
		cpu.PC += 2
	case 0x8000:
		// 8XYN - Various arithmetic and logical operations
		x := (instruction & 0x0F00) >> 8
		y := (instruction & 0x00F0) >> 4
		sel := instruction & 0x000F
		cpu.PerformArithmeticOperation(x, y, sel)
		cpu.PC += 2
	case 0x9000:
		// 9XY0 - SNE Vx, Vy
		// Skip next instruction if Vx != Vy
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
		cpu.I = instruction & 0x0FFF
		cpu.PC += 2
	case 0xB000:
		// BNNN - JP V0, addr
		// Jump to location NNN + V0
		cpu.PC = (instruction & 0x0FFF) + uint16(cpu.V[0])
	case 0xC000:
		// CXNN - RND Vx, byte
		// Set Vx = random byte AND NN
		x := (instruction & 0x0F00) >> 8
		nn := byte(instruction & 0x00FF)
		random := byte(rand.Intn(256))
		cpu.V[x] = random & nn
		cpu.PC += 2
	case 0xD000:
		// DXYN - DRW Vx, Vy, nibble
		// Display N-byte sprite starting at memory location I at (Vx, Vy)
		x := (instruction & 0x0F00) >> 8
		y := (instruction & 0x00F0) >> 4
		size := instruction & 0x000F
		cpu.DrawSprite(x, y, size)
		cpu.PC += 2
	case 0xE000:
		// EX9E - SKP Vx
		// Skip next instruction if key with the value of Vx is pressed
		// EXA1 - SKNP Vx
		// Skip next instruction if key with the value of Vx is not pressed
		x := (instruction & 0x0F00) >> 8
		opcodeLow := byte(instruction & 0x00FF)
		keyState := cpu.Keys[cpu.V[x]]
		if (keyState == true) && (opcodeLow == 0x9E) {
			cpu.PC += 4
		} else if (keyState == false) && (opcodeLow == 0xA1) {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	case 0xF000:
		// Various operations starting with F
		x := (instruction & 0x0F00) >> 8
		sel := instruction & 0x00FF
		cpu.Perform0xFOperation(x, sel)
		cpu.PC += 2
	default:
		// Unknown opcode
		fmt.Printf("Unknown opcode: 0x%04X\n", instruction)
		cpu.PC += 2
	}
}

// ClearScreen clears the display memory.
func (cpu *CPU) ClearScreen() {
	for i := range cpu.Display {
		cpu.Display[i] = 0
	}
}

// ReturnFromSubroutine returns from a subroutine by popping the return address from the stack.
func (cpu *CPU) ReturnFromSubroutine() {
	cpu.SP--
	cpu.PC = cpu.Stack[cpu.SP]
	cpu.PC += 2 // Increment PC after return to avoid infinite loop
}

// GetDisplay returns the current state of the display.
func (cpu *CPU) GetDisplay() [64 * 32]byte {
	return cpu.Display
}

// SetKey sets the state of a key in the keypad.
func (cpu *CPU) SetKey(key uint8, pressed bool) {
	if key < 16 {
		cpu.Keys[key] = pressed
	}
}

// GetKeys returns the current state of all keys.
func (cpu *CPU) GetKeys() [16]bool {
	return cpu.Keys
}

// GetDelayTimer returns the current value of the delay timer.
func (cpu *CPU) GetDelayTimer() uint8 {
	return cpu.DelayTimer
}

// GetSoundTimer returns the current value of the sound timer.
func (cpu *CPU) GetSoundTimer() uint8 {
	return cpu.SoundTimer
}

// UpdateTimers updates the delay and sound timers at 60Hz.
func (cpu *CPU) UpdateTimers() {
	if cpu.DelayTimer > 0 {
		cpu.DelayTimer--
	}
	if cpu.SoundTimer > 0 {
		cpu.SoundTimer--
		// TODO: Implement sound
	}
}
