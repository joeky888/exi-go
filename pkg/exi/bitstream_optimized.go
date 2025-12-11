package exi

// Optimized ReadBits implementation that reads multiple bits at once
// This version significantly reduces function call overhead by processing
// bits in larger chunks instead of calling readBit() for each bit.

// ReadBitsOptimized is an optimized version of ReadBits that processes
// multiple bits at once to reduce function call overhead.
func (bs *BitStream) ReadBitsOptimized(bitCount int) (uint32, error) {
	if bitCount <= 0 || bitCount > 32 {
		return 0, ErrBitCountTooLarge
	}
	if bs.data == nil {
		return 0, ErrBitstreamNotInitial
	}

	var value uint32
	bitsRemaining := bitCount

	for bitsRemaining > 0 {
		// Check if we need to advance to next byte
		if bs.bitCount == exiBitstreamMaxBitCount {
			if err := bs.ensureBitCapacity(); err != nil {
				return 0, err
			}
		}

		// Check bounds
		if bs.bytePos >= bs.dataSize {
			return 0, ErrBitstreamOverflow
		}

		// Calculate how many bits we can read from current byte
		bitsAvailableInByte := int(exiBitstreamMaxBitCount - bs.bitCount)
		bitsToRead := bitsRemaining
		if bitsToRead > bitsAvailableInByte {
			bitsToRead = bitsAvailableInByte
		}

		// Read multiple bits at once from current byte
		currentByte := bs.data[bs.bytePos]
		shift := exiBitstreamMaxBitCount - (bs.bitCount + uint8(bitsToRead))
		mask := uint8((1 << bitsToRead) - 1)
		bits := uint32((currentByte >> shift) & mask)

		// Append to result
		value = (value << bitsToRead) | bits

		// Update position
		bs.bitCount += uint8(bitsToRead)
		bitsRemaining -= bitsToRead
	}

	return value, nil
}

// WriteBitsOptimized is an optimized version of WriteBits that processes
// multiple bits at once.
func (bs *BitStream) WriteBitsOptimized(bitCount int, value uint32) error {
	if bitCount <= 0 || bitCount > 32 {
		return ErrBitCountTooLarge
	}
	if bs.data == nil {
		return ErrBitstreamNotInitial
	}

	bitsRemaining := bitCount

	for bitsRemaining > 0 {
		// Ensure we have a byte allocated
		if bs.data == nil {
			return ErrBitstreamNotInitial
		}
		if bs.bitCount == exiBitstreamMaxBitCount {
			if err := bs.ensureBitCapacity(); err != nil {
				return err
			}
		}
		if bs.bytePos >= bs.dataSize {
			return ErrBitstreamOverflow
		}

		// Calculate how many bits we can write to current byte
		bitsAvailableInByte := int(exiBitstreamMaxBitCount - bs.bitCount)
		bitsToWrite := bitsRemaining
		if bitsToWrite > bitsAvailableInByte {
			bitsToWrite = bitsAvailableInByte
		}

		// Extract the bits to write from value (from MSB side)
		shift := uint(bitsRemaining - bitsToWrite)
		mask := uint32((1 << bitsToWrite) - 1)
		bits := uint8((value >> shift) & mask)

		// Write to current byte
		current := &bs.data[bs.bytePos]
		bitShift := exiBitstreamMaxBitCount - (bs.bitCount + uint8(bitsToWrite))
		*current = *current | (bits << bitShift)

		// Update position
		bs.bitCount += uint8(bitsToWrite)
		bitsRemaining -= bitsToWrite
	}

	return nil
}
