package main

import (
	"go-r8t/cpu"
	"time"
)

func main() {
	// Create a new CPU instance
	cpu := cpu.NewCPU()

	// Load a test program
	testProgram := []byte{0x00, 0xE0, 0x12, 0x00}
	cpu.LoadProgram(testProgram)

	// Main emulation loop
	for {
		// Fetch the next instruction
		opcode := (uint16(cpu.Memory[cpu.PC]) << 8) | uint16(cpu.Memory[cpu.PC+1])

		// Execute the instruction
		cpu.ExecuteInstruction(opcode)

		// Update timers at 60Hz
		cpu.UpdateTimers()

		// Sleep to control emulation speed (optional)
		time.Sleep(time.Second / 60)
	}
}
