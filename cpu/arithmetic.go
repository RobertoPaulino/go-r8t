package cpu

import "fmt"

// PerformArithmeticOperation handles all arithmetic and bitwise operations
// specified by the 0x8XXX instructions.
// Parameters:
//   - x: The first register index
//   - y: The second register index
//   - sel: The selector nibble that determines which operation to perform
func (cpu *CPU) PerformArithmeticOperation(x, y, sel uint16) {
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
