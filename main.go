package main

import "math/rand"

type CPU struct {
	PC            uint16 // Program Counter
	memory        [4096]byte
	V             [16]byte // 16 registers (V0-VF)
	stack         [16]uint16
	I             uint16
	SP            uint8         // Stack Pointer
	delayTimer    uint8         // Delay Timer
	soundTimer    uint8         // Sound Timer
	keys          [16]bool      // Keypad state
	display       [64 * 32]byte // Display (64x32 screen)
	currentOpcode uint16        // The current instruction (opcode)
}

func main() {
	cpu := &CPU{}

	testProgram := []byte{0x00, 0xE0, 0x12, 0x00}

	cpu.loadProgram(testProgram)

}

func (cpu *CPU) loadProgram(program []byte) {
	for i := range program {
		cpu.memory[0x200+i] = program[i]
	}
}

func (cpu *CPU) executeInstruction(instruction uint16) {
	switch instruction & 0xF000 {
	case 0x0000:
		if instruction == 0x00E0 {
			cpu.clearScreen() // 0x00E0: Clear the display
			return
		} else if instruction == 0x00EE {
			cpu.returnFromSubroutine() // 0x00EE: Return from subroutine (sets PC to address on top of stack and decrements SP)
			return
		}
	case 0x1000: // 1NNN: Jump to address NNN
		nextAddress := instruction & 0x0FFF
		cpu.PC = nextAddress
	case 0x2000: // 2NNN: Call subroutine at address NNN (pushes current PC to stack, sets PC to NNN)
		target := instruction & 0x0FFF
		cpu.stack[cpu.SP] = cpu.PC
		cpu.SP++
		cpu.PC = target
	case 0x6000: // Set VX (register) to NN
		x := (instruction & 0x0F00) >> 8
		cpu.V[x] = byte(instruction & 0x00FF)
	case 0x7000: // Add NN to VX (without carry)
		x := (instruction & 0x0F00) >> 8
		cpu.V[x] += byte(instruction & 0x00FF)
	case 0xA000:
		cpu.I = instruction & 0x0FFF
	case 0x3000: // skip if register V[x] is equal to nn
		x := (instruction & 0x0F00) >> 8
		nn := byte(instruction & 0x00FF)
		if cpu.V[x] == nn {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	case 0x4000: // skip if register V[x] is not equal to nn
		x := (instruction & 0x0F00) >> 8
		nn := byte(instruction & 0x00FF)
		if cpu.V[x] != nn {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	case 0x5000: //compare registers x and y
		x := (instruction & 0x0F00) >> 8
		y := (instruction & 0x00F0) >> 4

		if cpu.V[x] == cpu.V[y] {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	case 0x8000: //Arithmetic and bitwise operations on registers

		x := (instruction & 0x0F00) >> 8
		y := (instruction & 0x00F0) >> 4
		sel := instruction & 0x000F // Get the last nibble
		cpu.performArithmeticOperation(x, y, sel)

	case 0x9000:
		x := (instruction & 0x0F00) >> 8
		y := (instruction & 0x00F0) >> 4

		if cpu.V[x] != cpu.V[y] {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	case 0xB000:
		nextAddress := (instruction & 0x0FFF) + uint16(cpu.V[0])
		cpu.PC = nextAddress
	case 0xC000:
		x := (instruction & 0x0F00) >> 8
		nn := byte(instruction & 0x00FF)
		random := byte(rand.Intn(256))
		cpu.V[x] = random & nn
	case 0xE000:
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
		x := (instruction & 0x0F00) >> 8
		sel := instruction & 0x00FF // Get the last nibble
		cpu.perform0xFOperation(x, sel)
	}
}

func (cpu *CPU) perform0xFOperation(x, sel uint16) {
	switch sel {
	case 0x07:
		cpu.V[x] = cpu.delayTimer
	case 0x0A:
		for i, key := range cpu.keys {
			if key == true {
				cpu.V[x] = byte(i)
				cpu.PC += 2
				break
			}
		}
	case 0x01:
		cpu.V[x] = cpu.soundTimer
	case 0x02:
		cpu.I = uint16(cpu.V[x]) * 5
	case 0x03:

	}
}

func (cpu *CPU) performArithmeticOperation(x, y, sel uint16) {
	switch sel {

	case 0x0:
		cpu.V[x] = cpu.V[y]

	//bitwise operations
	case 0x1:
		cpu.V[x] |= cpu.V[y]
	case 0x2:
		cpu.V[x] &= cpu.V[y]
	case 0x3:
		cpu.V[x] ^= cpu.V[y]

	case 0x4: // 0x8xy4: Vx = Vx + Vy, set VF to 1 if carry, 0 otherwise
		result := uint16(cpu.V[x]) + uint16(cpu.V[y])
		if result > 0xFF {
			cpu.V[0xF] = 1
		} else {
			cpu.V[0xF] = 0
		}
		cpu.V[x] = byte(result)
	case 0x5: // 0x8xy5: Vx = Vx - Vy, set VF to 0 if borrow, 1 otherwise
		if uint16(cpu.V[x]) > uint16(cpu.V[y]) {
			cpu.V[0xF] = 1
		} else {
			cpu.V[0xF] = 0
		}
		cpu.V[x] -= cpu.V[y]

	case 0x6: // 0x8xy6: Vx = Vx >> 1, set VF to least significant bit of Vx before shift
		cpu.V[0xF] = cpu.V[x] & 0x1
		cpu.V[x] >>= 1

	case 0x7: // 0x8xy7: Vx = Vy - Vx, set VF to 0 if borrow, 1 otherwise
		if uint16(cpu.V[y]) > uint16(cpu.V[x]) {
			cpu.V[0xF] = 1 // Set VF to 1 if no borrow
		} else {
			cpu.V[0xF] = 0 // Set VF to 0 if there is a borrow
		}
		cpu.V[x] = cpu.V[y] - cpu.V[x]
	case 0xE: // 0x8xyE: Vx = Vx << 1, set VF to most significant bit of Vx before shift
		cpu.V[0xF] = (cpu.V[x] >> 7) & 0x1 // Set VF to the MSB of Vx
		cpu.V[x] <<= 1
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
