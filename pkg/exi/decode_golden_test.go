package exi_test

import (
	"encoding/hex"
	"os"
	"testing"

	"example.com/exi-go/pkg/exi"
)

// TestDecodeGoldenTimestamps decodes the golden EXI files to extract and print
// the timestamps they contain. This helps us understand what timestamp values
// were used when the golden files were generated.
func TestDecodeGoldenTimestamps(t *testing.T) {
	goldenFiles := []struct {
		name string
		path string
	}{
		{"SessionSetupReq", goldenPath("SessionSetupReq.exi")},
		{"SessionSetupRes", goldenPath("SessionSetupRes.exi")},
		{"ServiceDiscoveryReq", goldenPath("ServiceDiscoveryReq.exi")},
	}

	for _, gf := range goldenFiles {
		t.Run(gf.name, func(t *testing.T) {
			data, err := os.ReadFile(gf.path)
			if err != nil {
				if os.IsNotExist(err) {
					t.Skipf("golden file not found: %s", gf.path)
				}
				t.Fatalf("failed to read golden file: %v", err)
			}

			t.Logf("Golden file: %s", gf.path)
			t.Logf("Size: %d bytes", len(data))
			t.Logf("Hex: %s", hex.EncodeToString(data))

			// Initialize a bitstream to decode
			bs := &exi.BitStream{}
			bs.Init(data, 0)

			// Decode EXI header (8 bits)
			header, err := bs.ReadBits(8)
			if err != nil {
				t.Fatalf("failed to read EXI header: %v", err)
			}
			t.Logf("EXI header: 0x%02x", header)

			// Decode event code (6 bits)
			eventCode, err := bs.ReadBits(6)
			if err != nil {
				t.Fatalf("failed to read event code: %v", err)
			}
			t.Logf("Event code: %d", eventCode)

			// Decode message START bit (1 bit) - Grammar ID=404/etc: START Header
			startBit, err := bs.ReadBits(1)
			if err != nil {
				t.Fatalf("failed to read start bit: %v", err)
			}
			t.Logf("Header START bit: %d", startBit)

			// Decode SessionID START bit (1 bit) - Grammar ID=277: START SessionID
			sidStartBit, err := bs.ReadBits(1)
			if err != nil {
				t.Fatalf("failed to read SessionID START bit: %v", err)
			}
			t.Logf("SessionID START bit: %d", sidStartBit)

			// Decode SessionID encoding flag (1 bit)
			sidFlag, err := bs.ReadBits(1)
			if err != nil {
				t.Fatalf("failed to read SessionID flag: %v", err)
			}
			t.Logf("SessionID encoding flag: %d", sidFlag)

			// Decode SessionID length (unsigned-var)
			sidLen, err := bs.ReadUnsignedVar()
			if err != nil {
				t.Fatalf("failed to read SessionID length: %v", err)
			}
			t.Logf("SessionID length: %d", sidLen)

			// Read SessionID bytes
			sidBytes := make([]byte, sidLen)
			for i := uint64(0); i < sidLen; i++ {
				b, err := bs.ReadOctet()
				if err != nil {
					t.Fatalf("failed to read SessionID byte %d: %v", i, err)
				}
				sidBytes[i] = b
			}
			t.Logf("SessionID: %s", hex.EncodeToString(sidBytes))

			// Decode SessionID END bit (1 bit)
			sidEndBit, err := bs.ReadBits(1)
			if err != nil {
				t.Fatalf("failed to read SessionID END bit: %v", err)
			}
			t.Logf("SessionID END bit: %d", sidEndBit)

			// Decode TimeStamp START bit (1 bit) - Grammar ID=278: START TimeStamp
			tsStartBit, err := bs.ReadBits(1)
			if err != nil {
				t.Fatalf("failed to read TimeStamp START bit: %v", err)
			}
			t.Logf("TimeStamp START bit: %d", tsStartBit)

			// Decode TimeStamp encoding flag (1 bit)
			tsFlag, err := bs.ReadBits(1)
			if err != nil {
				t.Fatalf("failed to read TimeStamp flag: %v", err)
			}
			t.Logf("TimeStamp encoding flag: %d", tsFlag)

			// Decode TimeStamp value (unsigned-var uint64)
			timestamp, err := bs.ReadUnsignedVar()
			if err != nil {
				t.Fatalf("failed to read TimeStamp value: %v", err)
			}
			t.Logf("TimeStamp: %d (0x%x)", timestamp, timestamp)

			// Decode TimeStamp END bit (1 bit)
			tsEndBit, err := bs.ReadBits(1)
			if err != nil {
				t.Fatalf("failed to read TimeStamp END bit: %v", err)
			}
			t.Logf("TimeStamp END bit: %d", tsEndBit)

			t.Logf("Successfully decoded timestamp from golden file: %d", timestamp)
		})
	}
}
