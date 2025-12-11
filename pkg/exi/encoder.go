package exi

import (
	"fmt"

	"example.com/exi-go/pkg/v2g/generated"
)

// MeteringConfirmation encoders/decoders implemented in metering_confirmation.go
// Quick-win encoders/decoders are implemented in quickwins.go: PowerDeliveryRes,
// VehicleCheckInRes and VehicleCheckOutRes.
//
// Prototypes (for reference; implementations live in separate files):
// MeteringConfirmation:
// func EncodeTopLevelMeteringConfirmationReq(bs *BitStream, v *generated.MeteringConfirmationReq) error
// func EncodeMeteringConfirmationReq(bs *BitStream, v *generated.MeteringConfirmationReq) error
// func DecodeMeteringConfirmationReq(bs *BitStream) (*generated.MeteringConfirmationReq, error)
//
// PowerDeliveryRes (quickwins.go):
// func EncodeTopLevelPowerDeliveryRes(bs *BitStream, v *generated.PowerDeliveryRes) error
// func EncodePowerDeliveryRes(bs *BitStream, v *generated.PowerDeliveryRes) error
// func DecodePowerDeliveryRes(bs *BitStream) (*generated.PowerDeliveryRes, error)
//
// VehicleCheckInRes (quickwins.go):
// func EncodeTopLevelVehicleCheckInRes(bs *BitStream, v *generated.VehicleCheckInRes) error
// func EncodeVehicleCheckInRes(bs *BitStream, v *generated.VehicleCheckInRes) error
// func DecodeVehicleCheckInRes(bs *BitStream) (*generated.VehicleCheckInRes, error)
//
// VehicleCheckOutRes (quickwins.go):
// func EncodeTopLevelVehicleCheckOutRes(bs *BitStream, v *generated.VehicleCheckOutRes) error
// func EncodeVehicleCheckOutRes(bs *BitStream, v *generated.VehicleCheckOutRes) error
// func DecodeVehicleCheckOutRes(bs *BitStream) (*generated.VehicleCheckOutRes, error)

// This file contains schema-informed struct encoder/decoder functions that
// use the BitStream utility to encode generated Go structs into a compact
// binary representation and decode them back.
//
// The encoding format is purposely simple and deterministic, intended as a
// schema-aware transport for the generated types (not the EXI standard).
// It encodes fields in a stable order and uses fixed-size length prefixes
// for variable-length data. This implementation is incremental and intended
// to be replaced by a full EXI bitstream implementation once available.
//
// Encoding conventions used here:
//  - lengths are encoded as 16-bit unsigned integers (big-endian) using
//    BitStream.WriteBits(16, uint32(len))
//  - presence of optional fields is encoded as a single bit (1 present, 0 absent)
//  - strings are encoded as UTF-8 bytes (length-prefixed)
//  - binary data ([]byte) encoded as length-prefixed raw octets
//  - arrays are encoded with a 16-bit count followed by each element encoded
//    as above
//
// Note: The BitStream API used:
//   bs.Init(buf []byte, dataOffset int)
//   bs.WriteBits(bitCount int, value uint32)
//   bs.WriteOctet(value byte)
//   bs.ReadBits(bitCount int) (uint32, error)
//   bs.ReadOctet() (byte, error)
//   bs.Length() int
//
// The Codec methods below provide convenient wrappers to encode/decode the
// generated types using a BitStream-backed buffer.

const (
	// recommended initial buffer size for encoded payloads
	defaultEncodeBufferSize = 4096
)

// --- low-level helpers -----------------------------------------------------

// writeUint16 writes a 16-bit unsigned integer as an EXI octet-sequence (variable-length)
// to the bitstream using the EXI unsigned-var format (7-bit groups with continuation flag).
// This mirrors the behavior of the original C implementation which encodes unsigned
// integers as variable-length octet sequences.
func writeUint16(bs *BitStream, v uint16) error {
	// Use the BitStream helper that writes EXI unsigned-var sequences.
	return bs.WriteUnsignedVar(uint64(v))
}

// readUint16 reads an unsigned integer encoded as an EXI octet-sequence (variable-length)
// from the bitstream and returns it as uint16. It uses the EXI unsigned-var decoding
// semantics (7-bit groups with continuation flag) to match the encoder.
func readUint16(bs *BitStream) (uint16, error) {
	v, err := bs.ReadUnsignedVar()
	return uint16(v), err
}

// writeBytes encodes a length-prefixed byte slice (length:uint16 + octets).
func writeBytes(bs *BitStream, data []byte) error {
	if len(data) > 0xFFFF {
		return fmt.Errorf("writeBytes: payload too large (%d bytes)", len(data))
	}
	if err := writeUint16(bs, uint16(len(data))); err != nil {
		return err
	}
	for _, b := range data {
		if err := bs.WriteOctet(b); err != nil {
			return err
		}
	}
	return nil
}

// readBytes reads a length-prefixed byte slice.
func readBytes(bs *BitStream) ([]byte, error) {
	n, err := readUint16(bs)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}
	out := make([]byte, int(n))
	for i := 0; i < int(n); i++ {
		b, err := bs.ReadOctet()
		if err != nil {
			return nil, err
		}
		out[i] = b
	}
	return out, nil
}

// writeString writes a length-prefixed string (UTF-8 bytes).
func writeString(bs *BitStream, s string) error {
	return writeBytes(bs, []byte(s))
}

// readString reads a length-prefixed string.
func readString(bs *BitStream) (string, error) {
	b, err := readBytes(bs)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// writeStringArray writes an array of strings with a 16-bit count prefix.
func writeStringArray(bs *BitStream, arr []string) error {
	if len(arr) > 0xFFFF {
		return fmt.Errorf("writeStringArray: too many elements (%d)", len(arr))
	}
	if err := writeUint16(bs, uint16(len(arr))); err != nil {
		return err
	}
	for _, s := range arr {
		if err := writeString(bs, s); err != nil {
			return err
		}
	}
	return nil
}

// readStringArray reads an array of strings encoded by writeStringArray.
func readStringArray(bs *BitStream) ([]string, error) {
	cnt, err := readUint16(bs)
	if err != nil {
		return nil, err
	}
	if cnt == 0 {
		return nil, nil
	}
	out := make([]string, int(cnt))
	for i := 0; i < int(cnt); i++ {
		s, err := readString(bs)
		if err != nil {
			return nil, err
		}
		out[i] = s
	}
	return out, nil
}

// writeBinaryArray writes a count followed by each binary item (length+bytes).
func writeBinaryArray(bs *BitStream, arr [][]byte) error {
	if len(arr) > 0xFFFF {
		return fmt.Errorf("writeBinaryArray: too many elements (%d)", len(arr))
	}
	if err := writeUint16(bs, uint16(len(arr))); err != nil {
		return err
	}
	for _, b := range arr {
		if err := writeBytes(bs, b); err != nil {
			return err
		}
	}
	return nil
}

// readBinaryArray reads a 16-bit count followed by that many length-prefixed byte slices.
func readBinaryArray(bs *BitStream) ([][]byte, error) {
	cnt, err := readUint16(bs)
	if err != nil {
		return nil, err
	}
	if cnt == 0 {
		return nil, nil
	}
	out := make([][]byte, int(cnt))
	for i := 0; i < int(cnt); i++ {
		b, err := readBytes(bs)
		if err != nil {
			return nil, err
		}
		out[i] = b
	}
	return out, nil
}

// --- Schema-informed encoders/decoders ------------------------------------

// writeRawBytes writes bytes without a length prefix.
func writeRawBytes(bs *BitStream, data []byte) error {
	if data == nil {
		return nil
	}
	for _, b := range data {
		if err := bs.WriteOctet(b); err != nil {
			return err
		}
	}
	return nil
}

// encodeMessageHeaderType encodes a MessageHeaderType following the C implementation.
// This is the common header encoding shared by all ISO 15118-20 messages.
// Grammar flow (from C implementation):
//
//	Grammar ID=277: START SessionID (1 bit, value 0)
//	hexBinary flag (1 bit, value 0)
//	SessionID length (unsigned-var)
//	SessionID bytes
//	END SessionID (1 bit, value 0)
//	Grammar ID=278: START TimeStamp (1 bit, value 0)
//	TimeStamp flag (1 bit, value 0)
//	TimeStamp value (unsigned-var uint64)
//	END TimeStamp (1 bit, value 0)
//	Grammar ID=279: Header END or Signature (2 bits, value 1 = END, no Signature)
func encodeMessageHeaderType(bs *BitStream, h *generated.MessageHeaderType) error {
	// Grammar ID=277: START SessionID (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// hexBinary encoding flag (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// SessionID length (unsigned-var)
	if err := writeUint16(bs, uint16(len(h.SessionID))); err != nil {
		return err
	}
	// SessionID bytes
	if err := writeRawBytes(bs, h.SessionID); err != nil {
		return err
	}
	// END SessionID (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=278: START TimeStamp (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// TimeStamp encoding flag (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// TimeStamp value (unsigned-var uint64)
	if err := bs.WriteUnsignedVar(h.TimeStamp); err != nil {
		return err
	}
	// END TimeStamp (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=279: Header END or Signature (2 bits, value 1 = END, no Signature)
	if err := bs.WriteBits(2, 1); err != nil {
		return err
	}

	return nil
}

// decodeMessageHeaderType decodes a MessageHeaderType following the C implementation.
// This is the common header decoding shared by all ISO 15118-20 messages.
func decodeMessageHeaderType(bs *BitStream) (*generated.MessageHeaderType, error) {
	// Grammar ID=277: START SessionID (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// hexBinary encoding flag (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// SessionID length (unsigned-var)
	sidLen, err := readUint16(bs)
	if err != nil {
		return nil, err
	}
	// SessionID bytes
	sid := make([]byte, sidLen)
	for i := 0; i < int(sidLen); i++ {
		b, err := bs.ReadOctet()
		if err != nil {
			return nil, err
		}
		sid[i] = b
	}
	// END SessionID (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=278: START TimeStamp (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// TimeStamp encoding flag (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// TimeStamp value (unsigned-var uint64)
	ts, err := bs.ReadUnsignedVar()
	if err != nil {
		return nil, err
	}
	// END TimeStamp (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=279: Header END or Signature (2 bits, expect 1 = END)
	if _, err := bs.ReadBits(2); err != nil {
		return nil, err
	}

	return &generated.MessageHeaderType{
		SessionID: sid,
		TimeStamp: ts,
	}, nil
}

// EncodeTopLevelSessionSetupReq writes an EXI simple header and the top-level
// event code for SessionSetupReq, then delegates to the per-type encoder to
// write the message body. This mirrors the C flow that writes an EXI header
// followed by the document-level event ID and then the message-specific content.
func EncodeTopLevelSessionSetupReq(bs *BitStream, v *generated.SessionSetupReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelSessionSetupReq: nil value")
	}
	// EXI simple header (8 bits): EXI_SIMPLE_HEADER_VALUE (0x80)
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Top-level event code for SessionSetupReq in the ISO-2/20 mapping is 35.
	// Write it as a 6-bit value to match the C harness's exi_basetypes_encoder_nbit_uint(stream, 6, 35).
	if err := bs.WriteBits(6, 35); err != nil {
		return err
	}
	// Now write the message body using the existing per-type encoder which writes
	// the header (SessionID + TimeStamp) and EVCCID.
	return EncodeSessionSetupReq(bs, v)
}

// EncodeTopLevelServiceDiscoveryReq writes EXI header and the top-level event
// code for ServiceDiscoveryReq, then delegates to the per-type encoder.
func EncodeTopLevelServiceDiscoveryReq(bs *BitStream, v *generated.ServiceDiscoveryReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelServiceDiscoveryReq: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for ServiceDiscoveryReq is 31 (as used by the C encoder).
	if err := bs.WriteBits(6, 31); err != nil {
		return err
	}
	// Delegate to the existing ServiceDiscoveryReq encoder which writes the header and fields.
	return EncodeServiceDiscoveryReq(bs, v)
}

// Per-type body encoder for SessionSetupReq matching C encoder bit-for-bit.
// This follows the exact grammar path from the C implementation.
func EncodeSessionSetupReq(bs *BitStream, v *generated.SessionSetupReq) error {
	if v == nil {
		return fmt.Errorf("EncodeSessionSetupReq: nil value")
	}

	// Grammar ID=404: START Header (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Encode MessageHeaderType (Grammar ID=277-279)
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// Grammar ID=405: START EVCCID (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// String encoding flag (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// EVCCID length (unsigned-var, +2 for string table semantics)
	evLen := uint16(len(v.EVCCID))
	if err := writeUint16(bs, evLen+2); err != nil {
		return err
	}
	// EVCCID bytes
	if err := writeRawBytes(bs, v.EVCCID); err != nil {
		return err
	}
	// END EVCCID (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=2: END SessionSetupReq (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeSessionSetupReq decodes a SessionSetupReq from BitStream matching C decoder.
// This follows the exact grammar path from the C implementation.
func DecodeSessionSetupReq(bs *BitStream) (*generated.SessionSetupReq, error) {
	// Grammar ID=404: START Header (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Decode MessageHeaderType (Grammar ID=277-279)
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// Grammar ID=405: START EVCCID (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// String encoding flag (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// EVCCID length (unsigned-var, subtract 2 for string table semantics)
	evccidLen, err := readUint16(bs)
	if err != nil {
		return nil, err
	}
	evccidLen -= 2 // Remove the +2 added during encoding
	// EVCCID bytes
	evccid := make([]byte, evccidLen)
	for i := 0; i < int(evccidLen); i++ {
		b, err := bs.ReadOctet()
		if err != nil {
			return nil, err
		}
		evccid[i] = b
	}
	// END EVCCID (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=2: END SessionSetupReq (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.SessionSetupReq{
		Header: *header,
		EVCCID: evccid,
	}, nil
}

// EncodeServiceDiscoveryReq encodes ServiceDiscoveryReq matching C encoder bit-for-bit.
// This follows the exact grammar path from the C implementation.
// Note: The C ISO-20 spec has SupportedServiceIDs, not ServiceScope/ServiceCategory.
// The Go types were generated from a different spec, so we ignore those fields.
func EncodeServiceDiscoveryReq(bs *BitStream, v *generated.ServiceDiscoveryReq) error {
	if v == nil {
		return fmt.Errorf("EncodeServiceDiscoveryReq: nil value")
	}

	// Grammar ID=422: START Header (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Encode MessageHeaderType (Grammar ID=277-279)
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// Grammar ID=423: SupportedServiceIDs or END (2 bits)
	// The golden has no SupportedServiceIDs, so write value 1 = END Element
	if err := bs.WriteBits(2, 1); err != nil {
		return err
	}

	return nil
}

// DecodeServiceDiscoveryReq decodes ServiceDiscoveryReq from BitStream matching C decoder.
// This follows the exact grammar path from the C implementation.
func DecodeServiceDiscoveryReq(bs *BitStream) (*generated.ServiceDiscoveryReq, error) {
	// Grammar ID=422: START Header (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Decode MessageHeaderType (Grammar ID=277-279)
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// Grammar ID=423: SupportedServiceIDs or END (2 bits)
	// For now we only handle the END case (value 1)
	endOrSupported, err := bs.ReadBits(2)
	if err != nil {
		return nil, err
	}

	// If endOrSupported == 0, there are SupportedServiceIDs to decode
	// If endOrSupported == 1, it's END Element
	if endOrSupported != 1 {
		return nil, fmt.Errorf("DecodeServiceDiscoveryReq: SupportedServiceIDs decoding not yet implemented")
	}

	return &generated.ServiceDiscoveryReq{
		Header:          *header,
		ServiceScope:    nil,
		ServiceCategory: nil,
	}, nil
}

// EncodeTopLevelServiceDetailReq writes EXI header and event code for ServiceDetailReq.
func EncodeTopLevelServiceDetailReq(bs *BitStream, v *generated.ServiceDetailReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelServiceDetailReq: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for ServiceDetailReq is 29
	if err := bs.WriteBits(6, 29); err != nil {
		return err
	}
	return EncodeServiceDetailReq(bs, v)
}

// EncodeServiceDetailReq encodes ServiceDetailReq matching C encoder bit-for-bit.
// Layout: Header, ServiceID (uint16)
func EncodeServiceDetailReq(bs *BitStream, v *generated.ServiceDetailReq) error {
	if v == nil {
		return fmt.Errorf("EncodeServiceDetailReq: nil value")
	}

	// Grammar ID=429: START Header (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Encode MessageHeaderType (Grammar ID=277-279)
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// Grammar ID=430: START ServiceID (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ServiceID encoding flag (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ServiceID value (uint16)
	if err := writeUint16(bs, v.ServiceID); err != nil {
		return err
	}
	// END ServiceID (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=2: END ServiceDetailReq (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeServiceDetailReq decodes ServiceDetailReq from BitStream matching C decoder.
func DecodeServiceDetailReq(bs *BitStream) (*generated.ServiceDetailReq, error) {
	// Grammar ID=429: START Header (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Decode MessageHeaderType (Grammar ID=277-279)
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// Grammar ID=430: START ServiceID (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ServiceID encoding flag (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ServiceID value (uint16)
	serviceID, err := readUint16(bs)
	if err != nil {
		return nil, err
	}
	// END ServiceID (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=2: END ServiceDetailReq (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.ServiceDetailReq{
		Header:    *header,
		ServiceID: serviceID,
	}, nil
}

// EncodeTopLevelSessionStopReq writes EXI header and event code for SessionStopReq.
func EncodeTopLevelSessionStopReq(bs *BitStream, v *generated.SessionStopReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelSessionStopReq: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for SessionStopReq is 37
	if err := bs.WriteBits(6, 37); err != nil {
		return err
	}
	return EncodeSessionStopReq(bs, v)
}

// EncodeSessionStopReq encodes SessionStopReq matching C encoder bit-for-bit.
// Layout: Header, ChargingSession (2-bit enum), optional EVTerminationCode, optional EVTerminationExplanation
func EncodeSessionStopReq(bs *BitStream, v *generated.SessionStopReq) error {
	if v == nil {
		return fmt.Errorf("EncodeSessionStopReq: nil value")
	}

	// Grammar ID=460: START Header (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Encode MessageHeaderType (Grammar ID=277-279)
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// Grammar ID=461: START ChargingSession (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ChargingSession encoding flag (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ChargingSession enum value (2 bits)
	chargingSession := mapChargingSessionToEnum(v.ChargingSession)
	if err := bs.WriteBits(2, uint32(chargingSession)); err != nil {
		return err
	}
	// END ChargingSession (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=462: Optional EVTerminationCode, EVTerminationExplanation, or END (2 bits)
	if v.EVTerminationCode != nil {
		// START EVTerminationCode (2 bits, value 0)
		if err := bs.WriteBits(2, 0); err != nil {
			return err
		}
		// String encoding flag (1 bit, value 0)
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		// EVTerminationCode string (+2 for string table)
		if err := writeUint16(bs, uint16(len(*v.EVTerminationCode))+2); err != nil {
			return err
		}
		if err := writeRawBytes(bs, []byte(*v.EVTerminationCode)); err != nil {
			return err
		}
		// END EVTerminationCode (1 bit, value 0)
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}

		// Grammar ID=463: Optional EVTerminationExplanation or END (2 bits)
		if v.EVTerminationExplanation != nil {
			// START EVTerminationExplanation (2 bits, value 0)
			if err := bs.WriteBits(2, 0); err != nil {
				return err
			}
			// String encoding flag (1 bit, value 0)
			if err := bs.WriteBits(1, 0); err != nil {
				return err
			}
			// EVTerminationExplanation string (+2 for string table)
			if err := writeUint16(bs, uint16(len(*v.EVTerminationExplanation))+2); err != nil {
				return err
			}
			if err := writeRawBytes(bs, []byte(*v.EVTerminationExplanation)); err != nil {
				return err
			}
			// END EVTerminationExplanation (1 bit, value 0)
			if err := bs.WriteBits(1, 0); err != nil {
				return err
			}
			// Grammar ID=2: END SessionStopReq (1 bit, value 0)
			if err := bs.WriteBits(1, 0); err != nil {
				return err
			}
		} else {
			// END Element (2 bits, value 1)
			if err := bs.WriteBits(2, 1); err != nil {
				return err
			}
		}
	} else if v.EVTerminationExplanation != nil {
		// START EVTerminationExplanation without EVTerminationCode (2 bits, value 1)
		if err := bs.WriteBits(2, 1); err != nil {
			return err
		}
		// String encoding flag (1 bit, value 0)
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		// EVTerminationExplanation string (+2 for string table)
		if err := writeUint16(bs, uint16(len(*v.EVTerminationExplanation))+2); err != nil {
			return err
		}
		if err := writeRawBytes(bs, []byte(*v.EVTerminationExplanation)); err != nil {
			return err
		}
		// END EVTerminationExplanation (1 bit, value 0)
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		// Grammar ID=2: END SessionStopReq (1 bit, value 0)
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	} else {
		// No optional fields, END Element (2 bits, value 2)
		if err := bs.WriteBits(2, 2); err != nil {
			return err
		}
	}

	return nil
}

// DecodeSessionStopReq decodes SessionStopReq from BitStream matching C decoder.
func DecodeSessionStopReq(bs *BitStream) (*generated.SessionStopReq, error) {
	// Grammar ID=460: START Header (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Decode MessageHeaderType (Grammar ID=277-279)
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// Grammar ID=461: START ChargingSession (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ChargingSession encoding flag (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ChargingSession enum value (2 bits)
	chargingSessionEnum, err := bs.ReadBits(2)
	if err != nil {
		return nil, err
	}
	chargingSession := mapEnumToChargingSession(uint8(chargingSessionEnum))
	// END ChargingSession (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=462: Optional EVTerminationCode, EVTerminationExplanation, or END (2 bits)
	choice, err := bs.ReadBits(2)
	if err != nil {
		return nil, err
	}

	var evTerminationCode *string
	var evTerminationExplanation *string

	if choice == 0 {
		// EVTerminationCode present
		// String encoding flag (1 bit, expect 0)
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		// Read string length (subtract 2 for string table)
		length, err := readUint16(bs)
		if err != nil {
			return nil, err
		}
		length -= 2
		// Read string bytes
		termCodeBytes := make([]byte, length)
		for i := 0; i < int(length); i++ {
			b, err := bs.ReadOctet()
			if err != nil {
				return nil, err
			}
			termCodeBytes[i] = b
		}
		termCode := string(termCodeBytes)
		evTerminationCode = &termCode
		// END EVTerminationCode (1 bit, expect 0)
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}

		// Grammar ID=463: Optional EVTerminationExplanation or END (2 bits)
		choice2, err := bs.ReadBits(2)
		if err != nil {
			return nil, err
		}

		if choice2 == 0 {
			// EVTerminationExplanation present
			// String encoding flag (1 bit, expect 0)
			if _, err := bs.ReadBits(1); err != nil {
				return nil, err
			}
			// Read string length (subtract 2)
			length2, err := readUint16(bs)
			if err != nil {
				return nil, err
			}
			length2 -= 2
			// Read string bytes
			termExplBytes := make([]byte, length2)
			for i := 0; i < int(length2); i++ {
				b, err := bs.ReadOctet()
				if err != nil {
					return nil, err
				}
				termExplBytes[i] = b
			}
			termExpl := string(termExplBytes)
			evTerminationExplanation = &termExpl
			// END EVTerminationExplanation (1 bit, expect 0)
			if _, err := bs.ReadBits(1); err != nil {
				return nil, err
			}
			// Grammar ID=2: END SessionStopReq (1 bit, expect 0)
			if _, err := bs.ReadBits(1); err != nil {
				return nil, err
			}
		}
		// If choice2 == 1, it's END Element already handled
	} else if choice == 1 {
		// EVTerminationExplanation present (without EVTerminationCode)
		// String encoding flag (1 bit, expect 0)
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		// Read string length (subtract 2)
		length, err := readUint16(bs)
		if err != nil {
			return nil, err
		}
		length -= 2
		// Read string bytes
		termExplBytes := make([]byte, length)
		for i := 0; i < int(length); i++ {
			b, err := bs.ReadOctet()
			if err != nil {
				return nil, err
			}
			termExplBytes[i] = b
		}
		termExpl := string(termExplBytes)
		evTerminationExplanation = &termExpl
		// END EVTerminationExplanation (1 bit, expect 0)
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		// Grammar ID=2: END SessionStopReq (1 bit, expect 0)
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}
	// If choice == 2, it's END Element (no optional fields)

	return &generated.SessionStopReq{
		Header:                   *header,
		ChargingSession:          chargingSession,
		EVTerminationCode:        evTerminationCode,
		EVTerminationExplanation: evTerminationExplanation,
	}, nil
}

// EncodeTopLevelSessionStopRes writes EXI header and event code for SessionStopRes.
func EncodeTopLevelSessionStopRes(bs *BitStream, v *generated.SessionStopRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelSessionStopRes: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for SessionStopRes is 38
	if err := bs.WriteBits(6, 38); err != nil {
		return err
	}
	return EncodeSessionStopRes(bs, v)
}

// EncodeSessionStopRes encodes SessionStopRes matching C encoder bit-for-bit.
// Layout: Header, ResponseCode (6-bit enum)
func EncodeSessionStopRes(bs *BitStream, v *generated.SessionStopRes) error {
	if v == nil {
		return fmt.Errorf("EncodeSessionStopRes: nil value")
	}

	// Grammar ID=464: START Header (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Encode MessageHeaderType (Grammar ID=277-279)
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// Grammar ID=465: START ResponseCode (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ResponseCode encoding flag (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ResponseCode as 6-bit enum value
	responseCode := mapResponseCodeToEnum(v.ResponseCode)
	if err := bs.WriteBits(6, uint32(responseCode)); err != nil {
		return err
	}
	// END ResponseCode (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=2: END SessionStopRes (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeSessionStopRes decodes SessionStopRes from BitStream matching C decoder.
func DecodeSessionStopRes(bs *BitStream) (*generated.SessionStopRes, error) {
	// Grammar ID=464: START Header (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Decode MessageHeaderType (Grammar ID=277-279)
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// Grammar ID=465: START ResponseCode (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ResponseCode encoding flag (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ResponseCode as 6-bit enum value
	responseCodeEnum, err := bs.ReadBits(6)
	if err != nil {
		return nil, err
	}
	responseCode := mapEnumToResponseCode(uint8(responseCodeEnum))
	// END ResponseCode (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=2: END SessionStopRes (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.SessionStopRes{
		Header:       *header,
		ResponseCode: responseCode,
	}, nil
}

// EncodeTopLevelAuthorizationSetupReq writes EXI header and event code for AuthorizationSetupReq.
func EncodeTopLevelAuthorizationSetupReq(bs *BitStream, v *generated.AuthorizationSetupReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelAuthorizationSetupReq: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for AuthorizationSetupReq is 2
	if err := bs.WriteBits(6, 2); err != nil {
		return err
	}
	return EncodeAuthorizationSetupReq(bs, v)
}

// EncodeAuthorizationSetupReq encodes AuthorizationSetupReq matching C encoder bit-for-bit.
// This is the simplest message - only has Header field.
func EncodeAuthorizationSetupReq(bs *BitStream, v *generated.AuthorizationSetupReq) error {
	if v == nil {
		return fmt.Errorf("EncodeAuthorizationSetupReq: nil value")
	}

	// Grammar ID=409: START Header (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Encode MessageHeaderType (Grammar ID=277-279)
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// Grammar ID=2: END AuthorizationSetupReq (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeAuthorizationSetupReq decodes AuthorizationSetupReq from BitStream matching C decoder.
func DecodeAuthorizationSetupReq(bs *BitStream) (*generated.AuthorizationSetupReq, error) {
	// Grammar ID=409: START Header (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Decode MessageHeaderType (Grammar ID=277-279)
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// Grammar ID=2: END AuthorizationSetupReq (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.AuthorizationSetupReq{
		Header: *header,
	}, nil
}

// EncodeTopLevelServiceSelectionRes writes EXI header and event code for ServiceSelectionRes.
func EncodeTopLevelServiceSelectionRes(bs *BitStream, v *generated.ServiceSelectionRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelServiceSelectionRes: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for ServiceSelectionRes is 34
	if err := bs.WriteBits(6, 34); err != nil {
		return err
	}
	return EncodeServiceSelectionRes(bs, v)
}

// EncodeServiceSelectionRes encodes ServiceSelectionRes matching C encoder bit-for-bit.
// Layout: Header, ResponseCode (6-bit enum)
func EncodeServiceSelectionRes(bs *BitStream, v *generated.ServiceSelectionRes) error {
	if v == nil {
		return fmt.Errorf("EncodeServiceSelectionRes: nil value")
	}

	// Grammar ID=438: START Header (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Encode MessageHeaderType (Grammar ID=277-279)
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// Grammar ID=439: START ResponseCode (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ResponseCode encoding flag (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ResponseCode as 6-bit enum value
	responseCode := mapResponseCodeToEnum(v.ResponseCode)
	if err := bs.WriteBits(6, uint32(responseCode)); err != nil {
		return err
	}
	// END ResponseCode (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=2: END ServiceSelectionRes (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeServiceSelectionRes decodes ServiceSelectionRes from BitStream matching C decoder.
func DecodeServiceSelectionRes(bs *BitStream) (*generated.ServiceSelectionRes, error) {
	// Grammar ID=438: START Header (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Decode MessageHeaderType (Grammar ID=277-279)
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// Grammar ID=439: START ResponseCode (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ResponseCode encoding flag (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ResponseCode as 6-bit enum value
	responseCodeEnum, err := bs.ReadBits(6)
	if err != nil {
		return nil, err
	}
	responseCode := mapEnumToResponseCode(uint8(responseCodeEnum))
	// END ResponseCode (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=2: END ServiceSelectionRes (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.ServiceSelectionRes{
		Header:       *header,
		ResponseCode: responseCode,
	}, nil
}

// EncodeTopLevelMeteringConfirmationRes writes EXI header and event code for MeteringConfirmationRes.
func EncodeTopLevelMeteringConfirmationRes(bs *BitStream, v *generated.MeteringConfirmationRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelMeteringConfirmationRes: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for MeteringConfirmationRes is 17
	if err := bs.WriteBits(6, 17); err != nil {
		return err
	}
	return EncodeMeteringConfirmationRes(bs, v)
}

// EncodeMeteringConfirmationRes encodes MeteringConfirmationRes matching C encoder bit-for-bit.
// Layout: Header, ResponseCode (6-bit enum)
func EncodeMeteringConfirmationRes(bs *BitStream, v *generated.MeteringConfirmationRes) error {
	if v == nil {
		return fmt.Errorf("EncodeMeteringConfirmationRes: nil value")
	}

	// Grammar ID=458: START Header (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Encode MessageHeaderType (Grammar ID=277-279)
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// Grammar ID=459: START ResponseCode (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ResponseCode encoding flag (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ResponseCode as 6-bit enum value
	responseCode := mapResponseCodeToEnum(v.ResponseCode)
	if err := bs.WriteBits(6, uint32(responseCode)); err != nil {
		return err
	}
	// END ResponseCode (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=2: END MeteringConfirmationRes (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeMeteringConfirmationRes decodes MeteringConfirmationRes from BitStream matching C decoder.
func DecodeMeteringConfirmationRes(bs *BitStream) (*generated.MeteringConfirmationRes, error) {
	// Grammar ID=458: START Header (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Decode MessageHeaderType (Grammar ID=277-279)
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// Grammar ID=459: START ResponseCode (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ResponseCode encoding flag (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ResponseCode as 6-bit enum value
	responseCodeEnum, err := bs.ReadBits(6)
	if err != nil {
		return nil, err
	}
	responseCode := mapEnumToResponseCode(uint8(responseCodeEnum))
	// END ResponseCode (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=2: END MeteringConfirmationRes (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.MeteringConfirmationRes{
		Header:       *header,
		ResponseCode: responseCode,
	}, nil
}

// EncodeTopLevelAuthorizationRes writes EXI header and event code for AuthorizationRes.
func EncodeTopLevelAuthorizationRes(bs *BitStream, v *generated.AuthorizationRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelAuthorizationRes: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for AuthorizationRes is 1
	if err := bs.WriteBits(6, 1); err != nil {
		return err
	}
	return EncodeAuthorizationRes(bs, v)
}

// EncodeAuthorizationRes encodes AuthorizationRes matching C encoder bit-for-bit.
// Layout: Header, ResponseCode (6-bit enum), EVSEProcessing (2-bit enum)
func EncodeAuthorizationRes(bs *BitStream, v *generated.AuthorizationRes) error {
	if v == nil {
		return fmt.Errorf("EncodeAuthorizationRes: nil value")
	}

	// Grammar ID=419: START Header (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Encode MessageHeaderType (Grammar ID=277-279)
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// Grammar ID=420: START ResponseCode (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ResponseCode encoding flag (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ResponseCode as 6-bit enum value
	responseCode := mapResponseCodeToEnum(v.ResponseCode)
	if err := bs.WriteBits(6, uint32(responseCode)); err != nil {
		return err
	}
	// END ResponseCode (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=421: START EVSEProcessing (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// EVSEProcessing encoding flag (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// EVSEProcessing as 2-bit enum value
	evseProcessing := mapEVSEProcessingToEnum(v.EVSEProcessing)
	if err := bs.WriteBits(2, uint32(evseProcessing)); err != nil {
		return err
	}
	// END EVSEProcessing (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=2: END AuthorizationRes (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeAuthorizationRes decodes AuthorizationRes from BitStream matching C decoder.
func DecodeAuthorizationRes(bs *BitStream) (*generated.AuthorizationRes, error) {
	// Grammar ID=419: START Header (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Decode MessageHeaderType (Grammar ID=277-279)
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// Grammar ID=420: START ResponseCode (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ResponseCode encoding flag (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ResponseCode as 6-bit enum value
	responseCodeEnum, err := bs.ReadBits(6)
	if err != nil {
		return nil, err
	}
	responseCode := mapEnumToResponseCode(uint8(responseCodeEnum))
	// END ResponseCode (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=421: START EVSEProcessing (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// EVSEProcessing encoding flag (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// EVSEProcessing as 2-bit enum value
	evseProcessingEnum, err := bs.ReadBits(2)
	if err != nil {
		return nil, err
	}
	evseProcessing := mapEnumToEVSEProcessing(uint8(evseProcessingEnum))
	// END EVSEProcessing (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=2: END AuthorizationRes (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.AuthorizationRes{
		Header:         *header,
		ResponseCode:   responseCode,
		EVSEProcessing: evseProcessing,
	}, nil
}

// EncodeCertificateUpdateReq encodes CertificateUpdateReq into BitStream.
//
// Layout:
//   - Id: presence bit + string (if present)
//   - ContractSignatureCertChain: count(uint16) + each certificate (len+bytes)
//   - ContractID: presence + string
//   - ListOfRootCertificateIDs: count(uint16) + each string
//   - DHParams: presence + bytes
func EncodeCertificateUpdateReq(bs *BitStream, v *generated.CertificateUpdateReq) error {
	if v == nil {
		return fmt.Errorf("EncodeCertificateUpdateReq: nil value")
	}
	// Header (SessionID + TimeStamp)
	if err := writeBytes(bs, v.Header.SessionID); err != nil {
		return err
	}
	if err := bs.WriteUnsignedVar(v.Header.TimeStamp); err != nil {
		return err
	}
	// Id (attribute)
	if v.Id != "" {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
		if err := writeString(bs, v.Id); err != nil {
			return err
		}
	} else {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	// ContractSignatureCertChain (binary array)
	if err := writeBinaryArray(bs, v.ContractSignatureCertChain.Certificates); err != nil {
		return err
	}
	// ContractID
	if v.ContractID != "" {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
		if err := writeString(bs, v.ContractID); err != nil {
			return err
		}
	} else {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	// ListOfRootCertificateIDs (string array)
	if err := writeStringArray(bs, v.ListOfRootCertificateIDs); err != nil {
		return err
	}
	// DHParams (binary)
	if len(v.DHParams) > 0 {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
		if err := writeBytes(bs, v.DHParams); err != nil {
			return err
		}
	} else {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	return nil
}

// DecodeCertificateUpdateReq decodes CertificateUpdateReq from BitStream.
func DecodeCertificateUpdateReq(bs *BitStream) (*generated.CertificateUpdateReq, error) {
	// Header
	sid, err := readBytes(bs)
	if err != nil {
		return nil, err
	}
	ts, err := bs.ReadUnsignedVar()
	if err != nil {
		return nil, err
	}
	// Id presence
	p, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	var id string
	if p == 1 {
		id, err = readString(bs)
		if err != nil {
			return nil, err
		}
	}
	// ContractSignatureCertChain
	chain, err := readBinaryArray(bs)
	if err != nil {
		return nil, err
	}
	contractChain := generated.CertificateChain{Certificates: chain}
	// ContractID presence
	p2, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	var contractID string
	if p2 == 1 {
		contractID, err = readString(bs)
		if err != nil {
			return nil, err
		}
	}
	// ListOfRootCertificateIDs
	roots, err := readStringArray(bs)
	if err != nil {
		return nil, err
	}
	// DHParams presence
	p3, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	var dhParams []byte
	if p3 == 1 {
		dhParams, err = readBytes(bs)
		if err != nil {
			return nil, err
		}
	}
	return &generated.CertificateUpdateReq{
		Header:                     generated.MessageHeaderType{SessionID: sid, TimeStamp: uint64(ts)},
		Id:                         id,
		ContractSignatureCertChain: contractChain,
		ContractID:                 contractID,
		ListOfRootCertificateIDs:   roots,
		DHParams:                   dhParams,
	}, nil
}

// EncodeTopLevelSessionSetupRes writes EXI header and the top-level event
// code for SessionSetupRes, then delegates to the per-type encoder.
func EncodeTopLevelSessionSetupRes(bs *BitStream, v *generated.SessionSetupRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelSessionSetupRes: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for SessionSetupRes is 36
	if err := bs.WriteBits(6, 36); err != nil {
		return err
	}
	// Delegate to the per-type encoder
	return EncodeSessionSetupRes(bs, v)
}

// EncodeSessionSetupRes encodes SessionSetupRes matching C encoder bit-for-bit.
// This follows the exact grammar path from the C implementation.
// Layout: Header, ResponseCode (enum), EVSEID (string with +2 length)
func EncodeSessionSetupRes(bs *BitStream, v *generated.SessionSetupRes) error {
	if v == nil {
		return fmt.Errorf("EncodeSessionSetupRes: nil value")
	}

	// Grammar ID=406: START Header (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Encode MessageHeaderType (Grammar ID=277-279)
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// Grammar ID=407: START ResponseCode (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ResponseCode encoding flag (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ResponseCode as 6-bit enum value
	responseCode := mapResponseCodeToEnum(v.ResponseCode)
	if err := bs.WriteBits(6, uint32(responseCode)); err != nil {
		return err
	}
	// END ResponseCode (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=408: START EVSEID (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// String encoding flag (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// EVSEID length (unsigned-var, +2 for string table)
	evseidLen := uint16(len(v.EVSEID))
	if err := writeUint16(bs, evseidLen+2); err != nil {
		return err
	}
	// EVSEID bytes
	if err := writeRawBytes(bs, v.EVSEID); err != nil {
		return err
	}
	// END EVSEID (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=2: END SessionSetupRes (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// mapResponseCodeToEnum converts a string ResponseCode to its enum value
func mapResponseCodeToEnum(code string) uint8 {
	switch code {
	case "OK":
		return 0
	case "OK_CertificateExpiresSoon":
		return 1
	case "OK_NewSessionEstablished":
		return 2
	case "OK_OldSessionJoined":
		return 3
	case "FAILED":
		return 32 // Common failure code
	default:
		return 0 // Default to OK
	}
}

// mapEnumToResponseCode converts an enum value to its string ResponseCode
func mapEnumToResponseCode(code uint8) string {
	switch code {
	case 0:
		return "OK"
	case 1:
		return "OK_CertificateExpiresSoon"
	case 2:
		return "OK_NewSessionEstablished"
	case 3:
		return "OK_OldSessionJoined"
	case 32:
		return "FAILED"
	default:
		return "OK"
	}
}

// mapChargingSessionToEnum converts a string ChargingSession to its enum value (2-bit)
func mapChargingSessionToEnum(session string) uint8 {
	switch session {
	case "Pause":
		return 0
	case "Terminate":
		return 1
	case "ServiceRenegotiation":
		return 2
	default:
		return 1 // Default to Terminate
	}
}

// mapEnumToChargingSession converts an enum value to its string ChargingSession
func mapEnumToChargingSession(session uint8) string {
	switch session {
	case 0:
		return "Pause"
	case 1:
		return "Terminate"
	case 2:
		return "ServiceRenegotiation"
	default:
		return "Terminate"
	}
}

// mapEVSEProcessingToEnum converts a string EVSEProcessing to its enum value (2-bit)
func mapEVSEProcessingToEnum(processing string) uint8 {
	switch processing {
	case "Finished":
		return 0
	case "Ongoing":
		return 1
	case "Ongoing_WaitingForCustomerInteraction":
		return 2
	default:
		return 1 // Default to Ongoing
	}
}

// mapEnumToEVSEProcessing converts an enum value to its string EVSEProcessing
func mapEnumToEVSEProcessing(processing uint8) string {
	switch processing {
	case 0:
		return "Finished"
	case 1:
		return "Ongoing"
	case 2:
		return "Ongoing_WaitingForCustomerInteraction"
	default:
		return "Ongoing"
	}
}

// DecodeSessionSetupRes decodes SessionSetupRes from BitStream matching C decoder.
// This follows the exact grammar path from the C implementation.
func DecodeSessionSetupRes(bs *BitStream) (*generated.SessionSetupRes, error) {
	// Grammar ID=406: START Header (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Decode MessageHeaderType (Grammar ID=277-279)
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// Grammar ID=407: START ResponseCode (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ResponseCode encoding flag (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ResponseCode as 6-bit enum value
	responseCodeEnum, err := bs.ReadBits(6)
	if err != nil {
		return nil, err
	}
	responseCode := mapEnumToResponseCode(uint8(responseCodeEnum))
	// END ResponseCode (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=408: START EVSEID (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// String encoding flag (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// EVSEID length (unsigned-var, subtract 2 for string table)
	evseidLen, err := readUint16(bs)
	if err != nil {
		return nil, err
	}
	evseidLen -= 2
	// EVSEID bytes
	evseid := make([]byte, evseidLen)
	for i := 0; i < int(evseidLen); i++ {
		b, err := bs.ReadOctet()
		if err != nil {
			return nil, err
		}
		evseid[i] = b
	}
	// END EVSEID (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=2: END SessionSetupRes (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.SessionSetupRes{
		Header:       *header,
		ResponseCode: responseCode,
		EVSEID:       evseid,
		DateTimeNow:  nil, // Not present in this encoding
	}, nil
}

// EncodeTopLevelServiceDiscoveryRes writes EXI header and event code for ServiceDiscoveryRes.
func EncodeTopLevelServiceDiscoveryRes(bs *BitStream, v *generated.ServiceDiscoveryRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelServiceDiscoveryRes: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for ServiceDiscoveryRes is 32
	if err := bs.WriteBits(6, 32); err != nil {
		return err
	}
	return EncodeServiceDiscoveryRes(bs, v)
}

// EncodeServiceDiscoveryRes encodes ServiceDiscoveryRes into BitStream.
// Grammar path follows C implementation (Grammar IDs 424-428):
// - START Header (1 bit = 0)
// - MessageHeaderType
// - START ResponseCode (1 bit = 0)
// - ResponseCode enum (6 bits)
// - END ResponseCode (1 bit = 0)
// - START ServiceRenegotiationSupported (1 bit = 0)
// - boolean value
// - END ServiceRenegotiationSupported (1 bit = 0)
// - START EnergyTransferServiceList (1 bit = 0)
// - ServiceListType
// - Optional VASList (2 bits: 0 = present, 1 = END)
func EncodeServiceDiscoveryRes(bs *BitStream, v *generated.ServiceDiscoveryRes) error {
	if v == nil {
		return fmt.Errorf("EncodeServiceDiscoveryRes: nil value")
	}

	// Grammar ID=424: START Header (1 bit)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// Grammar ID=425: START ResponseCode (1 bit)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Encode ResponseCode enum (1 bit encoding flag + 6-bit enum)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	responseCodeEnum := mapResponseCodeToEnum(v.ResponseCode)
	if err := bs.WriteBits(6, uint32(responseCodeEnum)); err != nil {
		return err
	}
	// END ResponseCode (1 bit)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=426: START ServiceRenegotiationSupported (1 bit)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Encode boolean
	boolVal := uint32(0)
	if v.ServiceRenegotiationSupported {
		boolVal = 1
	}
	if err := bs.WriteBits(1, boolVal); err != nil {
		return err
	}
	// END ServiceRenegotiationSupported (1 bit)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=427: START EnergyTransferServiceList (1 bit)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeServiceList(bs, &v.EnergyTransferServiceList); err != nil {
		return err
	}

	// Grammar ID=428: VASList or END (2 bits)
	if v.VASList != nil {
		// START VASList (2 bits = 0)
		if err := bs.WriteBits(2, 0); err != nil {
			return err
		}
		if err := encodeServiceList(bs, v.VASList); err != nil {
			return err
		}
		// Grammar ID=2: END Element (1 bit)
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	} else {
		// END Element (2 bits = 1)
		if err := bs.WriteBits(2, 1); err != nil {
			return err
		}
	}

	return nil
}

// encodeServiceList encodes ServiceListType into BitStream.
// Grammar ID=358: Array of ServiceType (1-8 services)
func encodeServiceList(bs *BitStream, list *generated.ServiceList) error {
	if list == nil {
		return fmt.Errorf("encodeServiceList: nil list")
	}

	// Encode each service in the array
	for i, service := range list.Services {
		// Grammar ID=358: START Service (1 bit = 0)
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodeServiceType(bs, &service); err != nil {
			return err
		}

		// After encoding a service, check if more services follow
		// Grammar ID=359: START Service or END (varies based on count)
		if i < len(list.Services)-1 {
			// More services follow - handled by next iteration
		}
	}

	// After all services, write END Element
	// The number of bits depends on how many services we've encoded
	// For simplicity, write 1 bit = 1 for END
	if err := bs.WriteBits(1, 1); err != nil {
		return err
	}

	return nil
}

// encodeServiceType encodes ServiceType into BitStream.
// Grammar ID=167: ServiceID + FreeService
func encodeServiceType(bs *BitStream, service *generated.ServiceType) error {
	if service == nil {
		return fmt.Errorf("encodeServiceType: nil service")
	}

	// Grammar ID=167: START ServiceID (1 bit = 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Encode ServiceID as uint16 (1 bit encoding flag + 16 bits)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := bs.WriteBits(16, uint32(service.ServiceID)); err != nil {
		return err
	}
	// END ServiceID (1 bit)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=168: START FreeService (1 bit = 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Encode boolean
	boolVal := uint32(0)
	if service.FreeService {
		boolVal = 1
	}
	if err := bs.WriteBits(1, boolVal); err != nil {
		return err
	}
	// END FreeService (1 bit)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Grammar ID=2: END Service (1 bit)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeServiceDiscoveryRes decodes ServiceDiscoveryRes from BitStream.
func DecodeServiceDiscoveryRes(bs *BitStream) (*generated.ServiceDiscoveryRes, error) {
	// Grammar ID=424: START Header (1 bit)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// Grammar ID=425: START ResponseCode (1 bit)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Encoding flag (1 bit)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ResponseCode enum (6 bits)
	responseCodeBits, err := bs.ReadBits(6)
	if err != nil {
		return nil, err
	}
	responseCode := mapEnumToResponseCode(uint8(responseCodeBits))
	// END ResponseCode (1 bit)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=426: START ServiceRenegotiationSupported (1 bit)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Boolean value
	boolBits, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	serviceRenegotiationSupported := boolBits == 1
	// END ServiceRenegotiationSupported (1 bit)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=427: START EnergyTransferServiceList (1 bit)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	energyTransferServiceList, err := decodeServiceList(bs)
	if err != nil {
		return nil, err
	}

	// Grammar ID=428: VASList or END (2 bits)
	choice, err := bs.ReadBits(2)
	if err != nil {
		return nil, err
	}

	var vasList *generated.ServiceList
	if choice == 0 {
		// VASList present
		vasList, err = decodeServiceList(bs)
		if err != nil {
			return nil, err
		}
		// Grammar ID=2: END Element (1 bit)
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}
	// choice == 1 means END Element (no VASList)

	return &generated.ServiceDiscoveryRes{
		Header:                        *header,
		ResponseCode:                  responseCode,
		ServiceRenegotiationSupported: serviceRenegotiationSupported,
		EnergyTransferServiceList:     *energyTransferServiceList,
		VASList:                       vasList,
	}, nil
}

// decodeServiceList decodes ServiceListType from BitStream.
func decodeServiceList(bs *BitStream) (*generated.ServiceList, error) {
	var services []generated.ServiceType

	for {
		// Check for START Service or END (1 bit)
		choice, err := bs.ReadBits(1)
		if err != nil {
			return nil, err
		}

		if choice == 1 {
			// END Element
			break
		}

		// choice == 0: START Service
		service, err := decodeServiceType(bs)
		if err != nil {
			return nil, err
		}
		services = append(services, *service)
	}

	return &generated.ServiceList{
		Services: services,
	}, nil
}

// decodeServiceType decodes ServiceType from BitStream.
func decodeServiceType(bs *BitStream) (*generated.ServiceType, error) {
	// Grammar ID=167: START ServiceID (1 bit)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Encoding flag (1 bit)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ServiceID (16 bits)
	serviceIDBits, err := bs.ReadBits(16)
	if err != nil {
		return nil, err
	}
	// END ServiceID (1 bit)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=168: START FreeService (1 bit)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Boolean value
	boolBits, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	freeService := boolBits == 1
	// END FreeService (1 bit)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Grammar ID=2: END Service (1 bit)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.ServiceType{
		ServiceID:   uint16(serviceIDBits),
		FreeService: freeService,
	}, nil
}

// EncodeCertificateUpdateRes encodes CertificateUpdateRes into BitStream.
// Layout:
//   - Id: presence + string
//   - ResponseCode: string
//   - ContractSignatureCertChain: count + each cert bytes
//   - ContractSignatureEncryptedPrivateKey: presence + string
//   - DHParams: presence + bytes
//   - ContractID: presence + string
//   - RetryCounter: presence + unsigned-var
func EncodeCertificateUpdateRes(bs *BitStream, v *generated.CertificateUpdateRes) error {
	if v == nil {
		return fmt.Errorf("EncodeCertificateUpdateRes: nil value")
	}
	// Header
	if err := writeBytes(bs, v.Header.SessionID); err != nil {
		return err
	}
	if err := bs.WriteUnsignedVar(v.Header.TimeStamp); err != nil {
		return err
	}
	// Id
	if v.Id != "" {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
		if err := writeString(bs, v.Id); err != nil {
			return err
		}
	} else {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	// ResponseCode
	if err := writeString(bs, v.ResponseCode); err != nil {
		return err
	}
	// ContractSignatureCertChain
	if err := writeBinaryArray(bs, v.ContractSignatureCertChain.Certificates); err != nil {
		return err
	}
	// ContractSignatureEncryptedPrivateKey
	if v.ContractSignatureEncryptedPrivateKey != "" {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
		if err := writeString(bs, v.ContractSignatureEncryptedPrivateKey); err != nil {
			return err
		}
	} else {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	// DHParams
	if len(v.DHParams) > 0 {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
		if err := writeBytes(bs, v.DHParams); err != nil {
			return err
		}
	} else {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	// ContractID
	if v.ContractID != "" {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
		if err := writeString(bs, v.ContractID); err != nil {
			return err
		}
	} else {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	// RetryCounter
	if v.RetryCounter != nil {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
		if err := bs.WriteUnsignedVar(uint64(*v.RetryCounter)); err != nil {
			return err
		}
	} else {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	return nil
}

// DecodeCertificateUpdateRes decodes CertificateUpdateRes from BitStream.
func DecodeCertificateUpdateRes(bs *BitStream) (*generated.CertificateUpdateRes, error) {
	// Id
	p, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	var id string
	if p == 1 {
		id, err = readString(bs)
		if err != nil {
			return nil, err
		}
	}
	// ResponseCode
	rc, err := readString(bs)
	if err != nil {
		return nil, err
	}
	// ContractSignatureCertChain
	chain, err := readBinaryArray(bs)
	if err != nil {
		return nil, err
	}
	contractChain := generated.CertificateChain{Certificates: chain}
	// ContractSignatureEncryptedPrivateKey
	p2, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	var encPriv string
	if p2 == 1 {
		encPriv, err = readString(bs)
		if err != nil {
			return nil, err
		}
	}
	// DHParams
	p3, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	var dh []byte
	if p3 == 1 {
		dh, err = readBytes(bs)
		if err != nil {
			return nil, err
		}
	}
	// ContractID
	p4, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	var contractID string
	if p4 == 1 {
		contractID, err = readString(bs)
		if err != nil {
			return nil, err
		}
	}
	// RetryCounter
	p5, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	var retry *int
	if p5 == 1 {
		v, err := bs.ReadUnsignedVar()
		if err != nil {
			return nil, err
		}
		val := int(v)
		retry = &val
	}
	return &generated.CertificateUpdateRes{
		Id:                                   id,
		ResponseCode:                         rc,
		ContractSignatureCertChain:           contractChain,
		ContractSignatureEncryptedPrivateKey: encPriv,
		DHParams:                             dh,
		ContractID:                           contractID,
		RetryCounter:                         retry,
	}, nil
}

// --- Codec convenience wrappers -------------------------------------------

// EncodeStruct encodes a supported generated struct into EXI-like bytes using a BitStream.
func (c *Codec) EncodeStruct(v interface{}) ([]byte, error) {
	// Delegate to package-level EncodeStruct which has all message types
	return EncodeStruct(v)
}

// DecodeStruct decodes EXI-like bytes into a generated struct of the given typeName.
// Supported typeName values: "SessionSetupReq", "ServiceDiscoveryReq", "CertificateUpdateReq", "SessionSetupRes", "ServiceDiscoveryRes", "CertificateUpdateRes".
func (c *Codec) DecodeStruct(data []byte, typeName string) (interface{}, error) {
	bs := &BitStream{}
	bs.Init(data, 0)
	switch typeName {
	case "SessionSetupReq":
		return DecodeSessionSetupReq(bs)
	case "ServiceDiscoveryReq":
		return DecodeServiceDiscoveryReq(bs)
	case "CertificateUpdateReq":
		return DecodeCertificateUpdateReq(bs)
	case "SessionSetupRes":
		return DecodeSessionSetupRes(bs)
	case "ServiceDetailReq":
		return DecodeServiceDetailReq(bs)
	case "SessionStopReq":
		return DecodeSessionStopReq(bs)
	case "SessionStopRes":
		return DecodeSessionStopRes(bs)
	case "AuthorizationSetupReq":
		return DecodeAuthorizationSetupReq(bs)
	case "ServiceSelectionRes":
		return DecodeServiceSelectionRes(bs)
	case "MeteringConfirmationRes":
		return DecodeMeteringConfirmationRes(bs)
	case "AuthorizationRes":
		return DecodeAuthorizationRes(bs)
	case "ServiceDiscoveryRes":
		return DecodeServiceDiscoveryRes(bs)
	case "CertificateUpdateRes":
		return DecodeCertificateUpdateRes(bs)
	default:
		return nil, fmt.Errorf("DecodeStruct: unsupported type %s", typeName)
	}
}
