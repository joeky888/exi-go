package exi

import (
	"fmt"

	"example.com/exi-go/pkg/v2g/generated"
)

// Control loop message type implementations for ISO 15118-20.
// These are the simplest message types - they contain only the header with no additional fields.
//
// CLReqControlMode (event 4): Header only (empty message body)
// CLResControlMode (event 5): Header only (empty message body)

// --------------------------- CLReqControlMode ------------------------------

// EncodeTopLevelCLReqControlMode writes EXI header + event code (4) then
// delegates to EncodeCLReqControlMode.
func EncodeTopLevelCLReqControlMode(bs *BitStream, v *generated.CLReqControlMode) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelCLReqControlMode: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for CLReqControlMode is 4
	if err := bs.WriteBits(6, 4); err != nil {
		return err
	}
	return EncodeCLReqControlMode(bs, v)
}

// EncodeCLReqControlMode encodes the CLReqControlMode body.
// This is an empty message type that only contains the header - no additional fields.
func EncodeCLReqControlMode(bs *BitStream, v *generated.CLReqControlMode) error {
	if v == nil {
		return fmt.Errorf("EncodeCLReqControlMode: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// END CLReqControlMode (no additional fields)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeCLReqControlMode decodes the CLReqControlMode body.
func DecodeCLReqControlMode(bs *BitStream) (*generated.CLReqControlMode, error) {
	// START Header
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// END CLReqControlMode
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.CLReqControlMode{
		Header: *header,
	}, nil
}

// --------------------------- CLResControlMode ------------------------------

// EncodeTopLevelCLResControlMode writes EXI header + event code (5) then
// delegates to EncodeCLResControlMode.
func EncodeTopLevelCLResControlMode(bs *BitStream, v *generated.CLResControlMode) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelCLResControlMode: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for CLResControlMode is 5
	if err := bs.WriteBits(6, 5); err != nil {
		return err
	}
	return EncodeCLResControlMode(bs, v)
}

// EncodeCLResControlMode encodes the CLResControlMode body.
// This is an empty message type that only contains the header - no additional fields.
func EncodeCLResControlMode(bs *BitStream, v *generated.CLResControlMode) error {
	if v == nil {
		return fmt.Errorf("EncodeCLResControlMode: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// END CLResControlMode (no additional fields)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeCLResControlMode decodes the CLResControlMode body.
func DecodeCLResControlMode(bs *BitStream) (*generated.CLResControlMode, error) {
	// START Header
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// END CLResControlMode
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.CLResControlMode{
		Header: *header,
	}, nil
}
