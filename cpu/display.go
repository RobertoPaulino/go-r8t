package cpu

// DrawSprite draws a sprite at coordinates (VX, VY) with N bytes of sprite data.
// The sprite is drawn using XOR logic, and VF is set to 1 if any pixels are flipped from set to unset.
// Parameters:
//   - x: The register index containing the X coordinate
//   - y: The register index containing the Y coordinate
//   - size: The number of bytes of sprite data to draw (height of sprite)
func (cpu *CPU) DrawSprite(x, y, size uint16) {
	xPos := uint16(cpu.V[x]) // X coordinate from register VX
	yPos := uint16(cpu.V[y]) // Y coordinate from register VY
	cpu.V[0xF] = 0           // Reset collision flag

	// Loop through each row of the sprite
	for j := uint16(0); j < size; j++ {
		// Get the sprite data for this row
		pixel := cpu.Memory[cpu.I+j]

		// Loop through each bit in the sprite data
		for i := uint16(0); i < 8; i++ {
			// Check if the current pixel is set in the sprite data (1)
			if (pixel & (0x80 >> i)) != 0 {
				// Calculate the x and y position with wrapping
				posX := (xPos + i) % 64
				posY := (yPos + j) % 32
				idx := posY*64 + posX

				// Check for collision and set VF if a pixel is flipped
				if cpu.Display[idx] == 1 {
					cpu.V[0xF] = 1 // Set collision flag
				}

				// XOR the pixel in the display
				cpu.Display[idx] ^= 1
			}
		}
	}
}
