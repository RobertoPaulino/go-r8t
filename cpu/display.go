package cpu

// DrawSprite draws a sprite at coordinates (VX, VY) with N bytes of sprite data.
// The sprite is drawn using XOR logic, and VF is set to 1 if any pixels are flipped from set to unset.
// Parameters:
//   - x: The register index containing the X coordinate
//   - y: The register index containing the Y coordinate
//   - size: The number of bytes of sprite data to draw (height of sprite)
func (cpu *CPU) DrawSprite(x, y, size uint16) {
	xCord := uint16(cpu.V[x]) // Cast to uint16 for correct math
	yCord := uint16(cpu.V[y]) // Cast to uint16 for correct math

	cpu.V[0xF] = 0 // Reset collision flag

	// Loop over the sprite rows (from 0 to size-1)
	for row := uint16(0); row < size; row++ {
		spriteByte := cpu.Memory[cpu.I+row] // Get the byte for the current row of the sprite

		// Loop through each bit (8 bits = 1 byte)
		for bit := uint16(0); bit < 8; bit++ {
			// Calculate screen x, y position
			screenX := (xCord + bit) % 64 // Wrap around horizontally if needed
			screenY := (yCord + row) % 32 // Wrap around vertically if needed

			// Get the index in the display array
			displayIndex := screenY*64 + screenX

			// Check if the pixel should be turned off (collision)
			if (spriteByte & (0x80 >> bit)) != 0 { // Check if the current bit is 1
				if cpu.Display[displayIndex] == 1 {
					cpu.V[0xF] = 1 // Set VF flag if there's a collision (pixel was already 1)
				}
				// Flip the pixel on the screen using XOR
				cpu.Display[displayIndex] ^= 1
			}
		}
	}
}
