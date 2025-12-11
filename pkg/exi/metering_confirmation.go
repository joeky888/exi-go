package exi

import (
	"fmt"

	"example.com/exi-go/pkg/v2g/generated"
)

// EncodeTopLevelMeteringConfirmationReq writes an EXI simple header and the
// top-level event code for MeteringConfirmationReq (event 16) and delegates
// to the per-type encoder.
//
// Mirrors the pattern used by other top-level encoders in this package.
func EncodeTopLevelMeteringConfirmationReq(bs *BitStream, v *generated.MeteringConfirmationReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelMeteringConfirmationReq: nil value")
	}
	// EXI simple header (8 bits): EXI_SIMPLE_HEADER_VALUE (0x80)
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Top-level event code for MeteringConfirmationReq is 16 (6-bit)
	if err := bs.WriteBits(6, 16); err != nil {
		return err
	}
	return EncodeMeteringConfirmationReq(bs, v)
}

// EncodeMeteringConfirmationReq encodes the MeteringConfirmationReq body.
// Current generated type only contains the common MessageHeaderType.
// If SignedMeteringData is added later, expand this function following
// established patterns.
func EncodeMeteringConfirmationReq(bs *BitStream, v *generated.MeteringConfirmationReq) error {
	if v == nil {
		return fmt.Errorf("EncodeMeteringConfirmationReq: nil value")
	}

	// START Header (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Encode shared MessageHeaderType (SessionID + TimeStamp)
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// END MeteringConfirmationReq (1 bit, value 0)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeMeteringConfirmationReq decodes a MeteringConfirmationReq from the BitStream.
// It currently decodes only the shared MessageHeaderType. If SignedMeteringData
// becomes part of the generated types, extend the decoder to parse it.
func DecodeMeteringConfirmationReq(bs *BitStream) (*generated.MeteringConfirmationReq, error) {
	// START Header (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Decode shared MessageHeaderType
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// END MeteringConfirmationReq (1 bit, expect 0)
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.MeteringConfirmationReq{
		Header: *header,
	}, nil
}
