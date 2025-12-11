//go:build go1.18
// +build go1.18

package exi

import (
	"testing"

	"example.com/exi-go/pkg/v2g/generated"
)

// FuzzSessionSetupReq fuzzes SessionSetupReq encoding/decoding
func FuzzSessionSetupReq(f *testing.F) {
	// Seed corpus
	f.Add([]byte{0x01, 0x02, 0x03, 0x04}, uint64(1234567890), []byte{0x0A, 0x1B, 0x2C, 0x3D, 0x4E, 0x5F})

	f.Fuzz(func(t *testing.T, sessionID []byte, timestamp uint64, evccid []byte) {
		// Limit input sizes to reasonable bounds
		if len(sessionID) > 8 || len(evccid) > 6 {
			return
		}

		msg := &generated.SessionSetupReq{
			Header: generated.MessageHeaderType{
				SessionID: sessionID,
				TimeStamp: timestamp,
			},
			EVCCID: evccid,
		}

		// Encode
		encoded, err := EncodeStruct(msg)
		if err != nil {
			return // Invalid input, skip
		}

		// Decode
		decoded, err := DecodeStruct(encoded, (*generated.SessionSetupReq)(nil))
		if err != nil {
			t.Fatalf("Decode failed after successful encode: %v", err)
		}

		// Verify
		decodedMsg := decoded.(*generated.SessionSetupReq)
		if string(decodedMsg.EVCCID) != string(msg.EVCCID) {
			t.Errorf("EVCCID mismatch: got %v, want %v", decodedMsg.EVCCID, msg.EVCCID)
		}
	})
}

// FuzzServiceDiscoveryRes fuzzes ServiceDiscoveryRes encoding/decoding
func FuzzServiceDiscoveryRes(f *testing.F) {
	// Seed corpus
	f.Add([]byte{0x0A, 0x1B, 0x2C, 0x3D}, uint64(1672531200), uint16(1), uint16(2))

	f.Fuzz(func(t *testing.T, sessionID []byte, timestamp uint64, serviceID1 uint16, serviceID2 uint16) {
		if len(sessionID) > 8 {
			return
		}

		msg := &generated.ServiceDiscoveryRes{
			Header: generated.MessageHeaderType{
				SessionID: sessionID,
				TimeStamp: timestamp,
			},
			ResponseCode:                  "OK",
			ServiceRenegotiationSupported: true,
			EnergyTransferServiceList: generated.ServiceList{
				Services: []generated.ServiceType{
					{ServiceID: serviceID1, FreeService: true},
					{ServiceID: serviceID2, FreeService: false},
				},
			},
		}

		// Encode
		encoded, err := EncodeStruct(msg)
		if err != nil {
			return
		}

		// Decode
		decoded, err := DecodeStruct(encoded, (*generated.ServiceDiscoveryRes)(nil))
		if err != nil {
			t.Fatalf("Decode failed after successful encode: %v", err)
		}

		// Verify
		decodedMsg := decoded.(*generated.ServiceDiscoveryRes)
		if decodedMsg.ResponseCode != msg.ResponseCode {
			t.Errorf("ResponseCode mismatch")
		}
	})
}

// FuzzBitStream fuzzes the BitStream operations directly
func FuzzBitStream(f *testing.F) {
	// Seed corpus with various bit counts
	f.Add([]byte{0xFF, 0x00, 0xAA, 0x55}, 1)
	f.Add([]byte{0xFF, 0x00, 0xAA, 0x55}, 6)
	f.Add([]byte{0xFF, 0x00, 0xAA, 0x55}, 8)
	f.Add([]byte{0xFF, 0x00, 0xAA, 0x55}, 16)

	f.Fuzz(func(t *testing.T, data []byte, bitCount int) {
		// Limit to valid bit counts
		if bitCount <= 0 || bitCount > 32 {
			return
		}
		if len(data) == 0 {
			return
		}

		// Test ReadBits
		bs := &BitStream{}
		bs.Init(data, 0)

		value, err := bs.ReadBits(bitCount)
		if err != nil {
			return // Not enough data, acceptable
		}

		// Test WriteBits
		buf := make([]byte, len(data)+10)
		wbs := &BitStream{}
		wbs.Init(buf, 0)

		err = wbs.WriteBits(bitCount, value)
		if err != nil {
			t.Fatalf("WriteBits failed: %v", err)
		}

		// Read back and verify
		wbs.Reset()
		readValue, err := wbs.ReadBits(bitCount)
		if err != nil {
			t.Fatalf("ReadBits after WriteBits failed: %v", err)
		}

		if readValue != value {
			t.Errorf("Value mismatch: wrote %v, read %v (bitCount=%d)", value, readValue, bitCount)
		}
	})
}

// FuzzBitStreamOptimized fuzzes the optimized BitStream operations
func FuzzBitStreamOptimized(f *testing.F) {
	// Seed corpus
	f.Add([]byte{0xFF, 0x00, 0xAA, 0x55}, 6)
	f.Add([]byte{0xFF, 0x00, 0xAA, 0x55}, 16)

	f.Fuzz(func(t *testing.T, data []byte, bitCount int) {
		if bitCount <= 0 || bitCount > 32 || len(data) == 0 {
			return
		}

		// Test original vs optimized consistency
		bs1 := &BitStream{}
		bs1.Init(data, 0)
		value1, err1 := bs1.ReadBits(bitCount)

		bs2 := &BitStream{}
		bs2.Init(data, 0)
		value2, err2 := bs2.ReadBitsOptimized(bitCount)

		// Both should succeed or both should fail
		if (err1 == nil) != (err2 == nil) {
			t.Fatalf("Inconsistent errors: original=%v, optimized=%v", err1, err2)
		}

		// If both succeeded, values should match
		if err1 == nil && value1 != value2 {
			t.Errorf("Value mismatch: original=%v, optimized=%v (bitCount=%d)", value1, value2, bitCount)
		}
	})
}
