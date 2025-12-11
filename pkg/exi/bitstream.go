package exi

import (
	"errors"
)

// Bitstream error definitions
var (
	ErrBitstreamOverflow   = errors.New("exi: bitstream overflow")
	ErrBitCountTooLarge    = errors.New("exi: bit count larger than 32")
	ErrInvalidBitCount     = errors.New("exi: invalid bit count")
	ErrBitstreamNotInitial = errors.New("exi: bitstream not initialized")
)

// Maximum number of bits per byte
const exiBitstreamMaxBitCount = 8

// BitStream is a lightweight bit/byte stream abstraction intended to mirror
// the functionality of the reference C exi_bitstream_t structure. It supports
// writing/reading individual bits and octets and tracking position/length
// within a backing buffer.
type BitStream struct {
	// backing data buffer (mutable)
	data     []byte
	dataSize int
	// current position within the buffer
	bytePos int
	// number of bits already used in current byte (0..7)
	bitCount uint8
	// initialization flag and saved offset for reset
	initCalled  bool
	flagBytePos int
	// optional status callback (not used here, placeholder)
	StatusCallback func(messageID int, statusCode int, value1 int, value2 int)
}

// Init sets up the bitstream to operate on the given buffer. dataOffset is
// the starting byte offset within data where payload begins. The stream will
// write/read starting at that offset. The BitStream keeps a reference to the
// provided slice (no copy) so callers must ensure it remains valid.
func (bs *BitStream) Init(data []byte, dataOffset int) {
	bs.data = data
	if data == nil {
		bs.dataSize = 0
	} else {
		bs.dataSize = len(data)
	}
	if dataOffset < 0 {
		dataOffset = 0
	}
	if dataOffset > bs.dataSize {
		dataOffset = bs.dataSize
	}
	bs.bytePos = dataOffset
	bs.bitCount = 0
	bs.initCalled = true
	bs.flagBytePos = dataOffset
}

// Reset resets the stream to the last saved init state (i.e., rewinds to the
// offset passed in Init). If Init hasn't been called, this resets to zero.
func (bs *BitStream) Reset() {
	if bs.initCalled {
		bs.bytePos = bs.flagBytePos
	} else {
		bs.bytePos = 0
	}
	bs.bitCount = 0
}

// DataSize returns the total capacity (size) of the backing buffer.
func (bs *BitStream) DataSize() int {
	return bs.dataSize
}

// Length returns the number of bytes that have been written/read so far,
// taking into account the initial offset used at Init.
func (bs *BitStream) Length() int {
	length := bs.bytePos
	if bs.initCalled && bs.flagBytePos > 0 {
		length -= bs.flagBytePos
	}
	if bs.bitCount > 0 {
		length += 1
	}
	return length
}

// internal helper: ensure there is space (or advance to next byte) before writing/reading a bit.
func (bs *BitStream) ensureBitCapacity() error {
	if bs.bitCount == exiBitstreamMaxBitCount {
		// We've filled a byte; advance to next byte if possible.
		if bs.bytePos < bs.dataSize {
			bs.bytePos++
			bs.bitCount = 0
			// If bytePos equals dataSize, subsequent writes will overflow.
			if bs.bytePos >= bs.dataSize {
				// If we've consumed the last available byte, further writes will overflow.
				// The caller will detect when attempting to write a bit.
			}
			return nil
		}
		return ErrBitstreamOverflow
	}
	return nil
}

// writeBit writes a single bit (0 or 1) to the stream.
func (bs *BitStream) writeBit(bit uint8) error {
	if bs.data == nil {
		return ErrBitstreamNotInitial
	}
	// Ensure we have capacity (advance byte if needed)
	if bs.bitCount == exiBitstreamMaxBitCount {
		if err := bs.ensureBitCapacity(); err != nil {
			return err
		}
	}
	// If byte position is beyond buffer, overflow
	if bs.bytePos >= bs.dataSize {
		return ErrBitstreamOverflow
	}
	// Get pointer to current byte
	current := &bs.data[bs.bytePos]
	if bs.bitCount == 0 {
		// Starting a fresh byte - clear it
		*current = 0
	}
	if bit != 0 {
		// Set the appropriate bit: place bits from MSB to LSB
		*current = *current | (1 << (exiBitstreamMaxBitCount - (bs.bitCount + 1)))
	}
	bs.bitCount++
	return nil
}

// readBit reads a single bit from the stream and returns it (0 or 1).
func (bs *BitStream) readBit() (uint8, error) {
	if bs.data == nil {
		return 0, ErrBitstreamNotInitial
	}
	// If we've consumed full byte, advance
	if bs.bitCount == exiBitstreamMaxBitCount {
		if err := bs.ensureBitCapacity(); err != nil {
			return 0, err
		}
	}
	// If byte position is beyond buffer, overflow
	if bs.bytePos >= bs.dataSize {
		return 0, ErrBitstreamOverflow
	}
	current := bs.data[bs.bytePos]
	shift := exiBitstreamMaxBitCount - (bs.bitCount + 1)
	b := (current >> shift) & 1
	bs.bitCount++
	return b, nil
}

// WriteBits writes bitCount bits of value (most-significant-bit first) to the stream.
// bitCount must be between 1 and 32 inclusive.
func (bs *BitStream) WriteBits(bitCount int, value uint32) error {
	if bitCount <= 0 || bitCount > 32 {
		return ErrBitCountTooLarge
	}
	// Write each bit from MSB to LSB
	for n := 0; n < bitCount; n++ {
		shift := uint(bitCount - n - 1)
		bit := uint8((value >> shift) & 1)
		if err := bs.writeBit(bit); err != nil {
			return err
		}
	}
	return nil
}

// WriteOctet writes one byte to the stream using WriteBits(8, value).
func (bs *BitStream) WriteOctet(value byte) error {
	return bs.WriteBits(8, uint32(value))
}

// ReadBits reads bitCount bits from the stream (MSB-first) and returns the
// resulting uint32. bitCount must be between 1 and 32 inclusive.
func (bs *BitStream) ReadBits(bitCount int) (uint32, error) {
	if bitCount <= 0 || bitCount > 32 {
		return 0, ErrBitCountTooLarge
	}
	var value uint32
	for i := 0; i < bitCount; i++ {
		b, err := bs.readBit()
		if err != nil {
			return 0, err
		}
		value = (value << 1) | uint32(b)
	}
	return value, nil
}

// ReadOctet reads a single byte (8 bits) from the stream.
func (bs *BitStream) ReadOctet() (byte, error) {
	v, err := bs.ReadBits(8)
	if err != nil {
		return 0, err
	}
	return byte(v & 0xFF), nil
}

// WriteUnsignedVar writes a variable-length unsigned integer using the EXI
// octet-sequence format (7-bit groups with continuation flag in high bit).
// This mirrors the behavior of the C exi_basetypes_convert_to_unsigned and
// exi_basetypes_encoder_write_unsigned helpers.
func (bs *BitStream) WriteUnsignedVar(value uint64) error {
	// Collect 7-bit groups, low-order first
	var octs []byte
	v := value
	for {
		b := byte(v & 0x7F)
		v >>= 7
		if v != 0 {
			b |= 0x80
		}
		octs = append(octs, b)
		if v == 0 {
			break
		}
	}

	// Write each octet as a full octet in the same order produced (LSB-first)
	for _, o := range octs {
		if err := bs.WriteOctet(o); err != nil {
			return err
		}
	}

	return nil
}

// ReadUnsignedVar reads a variable-length unsigned integer encoded as an
// EXI octet-sequence (7-bit groups with continuation flag) and returns the value.
func (bs *BitStream) ReadUnsignedVar() (uint64, error) {
	var shift uint
	var result uint64
	for {
		b, err := bs.ReadOctet()
		if err != nil {
			return 0, err
		}
		value := uint64(b & 0x7F)
		result |= value << shift
		if (b & 0x80) == 0 {
			break
		}
		shift += 7
		if shift >= 64 {
			return 0, ErrBitCountTooLarge
		}
	}
	return result, nil
}
