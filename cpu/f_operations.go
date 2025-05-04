package cpu

import "fmt"

// Perform0xFOperation handles all instructions starting with 0xF.
// These instructions are typically used for memory and timer operations.
// Parameters:
//   - x: The register index specified in the instruction
//   - sel: The selector byte that determines which 0xF operation to perform
func (cpu *CPU) Perform0xFOperation(x, sel uint16) {
	switch sel {
	case 0x07:
		// FX07 - LD Vx, DT
		// Set Vx = delay timer value
		// The value of DT is placed into Vx.
		cpu.V[x] = cpu.DelayTimer
	case 0x0A:
		// FX0A - LD Vx, K
		// Wait for a key press, store the value of the key in Vx
		// All execution stops until a key is pressed, then the value of that key is stored in Vx.
		pressed := false
		for i := 0; i < len(cpu.Keys); i++ {
			if cpu.Keys[i] {
				cpu.V[x] = byte(i)
				pressed = true
				break
			}
		}
		if !pressed {
			// If no key is pressed, decrement PC to run this instruction again
			cpu.PC -= 2
			return // Important: Skip the automatic PC increment
		}
	case 0x15:
		// FX15 - LD DT, Vx
		// Set delay timer = Vx
		// DT is set equal to the value of Vx.
		cpu.DelayTimer = cpu.V[x]
	case 0x18:
		// FX18 - LD ST, Vx
		// Set sound timer = Vx
		// ST is set equal to the value of Vx.
		cpu.SoundTimer = cpu.V[x]
	case 0x1E:
		// FX1E - ADD I, Vx
		// Set I = I + Vx
		// The values of I and Vx are added, and the results are stored in I.
		// Set VF to 1 if there's a range overflow (I+Vx > 0xFFF), 0 otherwise
		if cpu.I+uint16(cpu.V[x]) > 0xFFF {
			cpu.V[0xF] = 1
		} else {
			cpu.V[0xF] = 0
		}
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
		cpu.Memory[cpu.I] = cpu.V[x] / 100
		cpu.Memory[cpu.I+1] = (cpu.V[x] / 10) % 10
		cpu.Memory[cpu.I+2] = cpu.V[x] % 10
	case 0x55:
		// FX55 - LD [I], Vx
		// Store registers V0 through Vx in memory starting at location I
		// The interpreter copies the values of registers V0 through Vx into memory, starting at the address in I.
		for i := uint16(0); i <= x; i++ {
			cpu.Memory[cpu.I+i] = cpu.V[i]
		}
		// Original CHIP-8 behavior increments I
		cpu.I += x + 1
	case 0x65:
		// FX65 - LD Vx, [I]
		// Read registers V0 through Vx from memory starting at location I
		// The interpreter reads values from memory starting at location I into registers V0 through Vx.
		for i := uint16(0); i <= x; i++ {
			cpu.V[i] = cpu.Memory[cpu.I+i]
		}
		// Original CHIP-8 behavior increments I
		cpu.I += x + 1
	default:
		// Unknown opcode
		fmt.Printf("Unknown 0xF opcode: 0x%02X\n", sel)
	}
}
