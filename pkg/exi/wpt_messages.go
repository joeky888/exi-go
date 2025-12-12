package exi

import (
	"fmt"

	"example.com/exi-go/pkg/v2g/generated"
)

// WPT (Wireless Power Transfer) message type implementations for ISO 15118-20.
// These messages support wireless charging operations including alignment,
// positioning, and charge loop control.

// ========================== WPT_AlignmentCheckReq ==========================

// EncodeTopLevelWPT_AlignmentCheckReq writes EXI header + event code (53) then
// delegates to EncodeWPT_AlignmentCheckReq.
func EncodeTopLevelWPT_AlignmentCheckReq(bs *BitStream, v *generated.WPT_AlignmentCheckReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelWPT_AlignmentCheckReq: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for WPT_AlignmentCheckReq is 53
	if err := bs.WriteBits(6, 53); err != nil {
		return err
	}
	return EncodeWPT_AlignmentCheckReq(bs, v)
}

// EncodeWPT_AlignmentCheckReq encodes the WPT_AlignmentCheckReq body.
func EncodeWPT_AlignmentCheckReq(bs *BitStream, v *generated.WPT_AlignmentCheckReq) error {
	if v == nil {
		return fmt.Errorf("EncodeWPT_AlignmentCheckReq: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// END WPT_AlignmentCheckReq (no additional fields)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeWPT_AlignmentCheckReq decodes the WPT_AlignmentCheckReq body.
func DecodeWPT_AlignmentCheckReq(bs *BitStream) (*generated.WPT_AlignmentCheckReq, error) {
	// START Header
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// END WPT_AlignmentCheckReq
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.WPT_AlignmentCheckReq{
		Header: *header,
	}, nil
}

// ========================== WPT_AlignmentCheckRes ==========================

// EncodeTopLevelWPT_AlignmentCheckRes writes EXI header + event code (54) then
// delegates to EncodeWPT_AlignmentCheckRes.
func EncodeTopLevelWPT_AlignmentCheckRes(bs *BitStream, v *generated.WPT_AlignmentCheckRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelWPT_AlignmentCheckRes: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for WPT_AlignmentCheckRes is 54
	if err := bs.WriteBits(6, 54); err != nil {
		return err
	}
	return EncodeWPT_AlignmentCheckRes(bs, v)
}

// EncodeWPT_AlignmentCheckRes encodes the WPT_AlignmentCheckRes body.
func EncodeWPT_AlignmentCheckRes(bs *BitStream, v *generated.WPT_AlignmentCheckRes) error {
	if v == nil {
		return fmt.Errorf("EncodeWPT_AlignmentCheckRes: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// ResponseCode (required, encoded as enum)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	responseCode := mapResponseCodeToEnum(v.ResponseCode)
	if err := bs.WriteBits(6, uint32(responseCode)); err != nil {
		return err
	}
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// AlignmentStatus (required, encoded as enum - simplified as 2-bit)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Map alignment status to simple enum (0=Aligned, 1=NotAligned, 2=InProgress)
	var alignStatus uint32 = 0
	switch v.AlignmentStatus {
	case "Aligned":
		alignStatus = 0
	case "NotAligned":
		alignStatus = 1
	case "InProgress":
		alignStatus = 2
	}
	if err := bs.WriteBits(2, alignStatus); err != nil {
		return err
	}
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// AlignmentOffset_X (optional)
	if v.AlignmentOffset_X != nil {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodeRationalNumber(bs, v.AlignmentOffset_X); err != nil {
			return err
		}
	}

	// AlignmentOffset_Y (optional)
	if v.AlignmentOffset_Y != nil {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodeRationalNumber(bs, v.AlignmentOffset_Y); err != nil {
			return err
		}
	}

	// AlignmentOffset_Z (optional)
	if v.AlignmentOffset_Z != nil {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodeRationalNumber(bs, v.AlignmentOffset_Z); err != nil {
			return err
		}
	}

	// END WPT_AlignmentCheckRes
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeWPT_AlignmentCheckRes decodes the WPT_AlignmentCheckRes body.
func DecodeWPT_AlignmentCheckRes(bs *BitStream) (*generated.WPT_AlignmentCheckRes, error) {
	// START Header
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// ResponseCode
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	responseCodeVal, err := bs.ReadBits(6)
	if err != nil {
		return nil, err
	}
	responseCode := mapEnumToResponseCode(uint8(responseCodeVal))
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// AlignmentStatus
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	alignStatusVal, err := bs.ReadBits(2)
	if err != nil {
		return nil, err
	}
	var alignmentStatus string
	switch alignStatusVal {
	case 0:
		alignmentStatus = "Aligned"
	case 1:
		alignmentStatus = "NotAligned"
	case 2:
		alignmentStatus = "InProgress"
	default:
		alignmentStatus = "NotAligned"
	}
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	res := &generated.WPT_AlignmentCheckRes{
		Header:          *header,
		ResponseCode:    responseCode,
		AlignmentStatus: alignmentStatus,
	}

	// Check for optional AlignmentOffset_X
	eventCode, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	if eventCode == 0 {
		offsetX, err := decodeRationalNumber(bs)
		if err != nil {
			return nil, err
		}
		res.AlignmentOffset_X = offsetX

		// Check for optional AlignmentOffset_Y
		eventCode, err = bs.ReadBits(1)
		if err != nil {
			return nil, err
		}
		if eventCode == 0 {
			offsetY, err := decodeRationalNumber(bs)
			if err != nil {
				return nil, err
			}
			res.AlignmentOffset_Y = offsetY

			// Check for optional AlignmentOffset_Z
			eventCode, err = bs.ReadBits(1)
			if err != nil {
				return nil, err
			}
			if eventCode == 0 {
				offsetZ, err := decodeRationalNumber(bs)
				if err != nil {
					return nil, err
				}
				res.AlignmentOffset_Z = offsetZ

				// END marker
				if _, err := bs.ReadBits(1); err != nil {
					return nil, err
				}
			}
		}
	}

	return res, nil
}

// ========================== WPT_FinePositioningReq ==========================

// EncodeTopLevelWPT_FinePositioningReq writes EXI header + event code (55).
func EncodeTopLevelWPT_FinePositioningReq(bs *BitStream, v *generated.WPT_FinePositioningReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelWPT_FinePositioningReq: nil value")
	}
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	if err := bs.WriteBits(6, 55); err != nil {
		return err
	}
	return EncodeWPT_FinePositioningReq(bs, v)
}

func EncodeWPT_FinePositioningReq(bs *BitStream, v *generated.WPT_FinePositioningReq) error {
	if v == nil {
		return fmt.Errorf("EncodeWPT_FinePositioningReq: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// END (no additional fields)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

func DecodeWPT_FinePositioningReq(bs *BitStream) (*generated.WPT_FinePositioningReq, error) {
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.WPT_FinePositioningReq{
		Header: *header,
	}, nil
}

// ========================== WPT_FinePositioningRes ==========================

// EncodeTopLevelWPT_FinePositioningRes writes EXI header + event code (56).
func EncodeTopLevelWPT_FinePositioningRes(bs *BitStream, v *generated.WPT_FinePositioningRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelWPT_FinePositioningRes: nil value")
	}
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	if err := bs.WriteBits(6, 56); err != nil {
		return err
	}
	return EncodeWPT_FinePositioningRes(bs, v)
}

func EncodeWPT_FinePositioningRes(bs *BitStream, v *generated.WPT_FinePositioningRes) error {
	if v == nil {
		return fmt.Errorf("EncodeWPT_FinePositioningRes: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// ResponseCode (required)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	responseCode := mapResponseCodeToEnum(v.ResponseCode)
	if err := bs.WriteBits(6, uint32(responseCode)); err != nil {
		return err
	}
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// PositioningStatus (required) - encoded as 2-bit enum
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	var posStatus uint32 = 0
	switch v.PositioningStatus {
	case "Complete":
		posStatus = 0
	case "InProgress":
		posStatus = 1
	case "Failed":
		posStatus = 2
	}
	if err := bs.WriteBits(2, posStatus); err != nil {
		return err
	}
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Optional fields - simplified encoding (skip for stub)
	// if v.GuidanceDirection != nil { ... }

	// END
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

func DecodeWPT_FinePositioningRes(bs *BitStream) (*generated.WPT_FinePositioningRes, error) {
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// ResponseCode
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	responseCodeVal, err := bs.ReadBits(6)
	if err != nil {
		return nil, err
	}
	responseCode := mapEnumToResponseCode(uint8(responseCodeVal))
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// PositioningStatus
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	posStatusVal, err := bs.ReadBits(2)
	if err != nil {
		return nil, err
	}
	var positioningStatus string
	switch posStatusVal {
	case 0:
		positioningStatus = "Complete"
	case 1:
		positioningStatus = "InProgress"
	case 2:
		positioningStatus = "Failed"
	default:
		positioningStatus = "InProgress"
	}
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Skip optional fields and END marker
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.WPT_FinePositioningRes{
		Header:            *header,
		ResponseCode:      responseCode,
		PositioningStatus: positioningStatus,
	}, nil
}

// ========================== WPT_ChargeLoopReq ==========================

// EncodeTopLevelWPT_ChargeLoopReq writes EXI header + event code (57).
func EncodeTopLevelWPT_ChargeLoopReq(bs *BitStream, v *generated.WPT_ChargeLoopReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelWPT_ChargeLoopReq: nil value")
	}
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	if err := bs.WriteBits(6, 57); err != nil {
		return err
	}
	return EncodeWPT_ChargeLoopReq(bs, v)
}

func EncodeWPT_ChargeLoopReq(bs *BitStream, v *generated.WPT_ChargeLoopReq) error {
	if v == nil {
		return fmt.Errorf("EncodeWPT_ChargeLoopReq: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// EVProcessing (required) - use existing EVProcessing mapping
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	evProc := mapEVProcessingToEnum(v.EVProcessing)
	if err := bs.WriteBits(2, uint32(evProc)); err != nil {
		return err
	}
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Simplified - skip optional fields for now
	// END
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

func DecodeWPT_ChargeLoopReq(bs *BitStream) (*generated.WPT_ChargeLoopReq, error) {
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// EVProcessing
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	evProcVal, err := bs.ReadBits(2)
	if err != nil {
		return nil, err
	}
	evProcessing := mapEnumToEVProcessing(uint8(evProcVal))
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Skip to END marker
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.WPT_ChargeLoopReq{
		Header:       *header,
		EVProcessing: evProcessing,
	}, nil
}

// ========================== WPT_ChargeLoopRes ==========================

// EncodeTopLevelWPT_ChargeLoopRes writes EXI header + event code (58).
func EncodeTopLevelWPT_ChargeLoopRes(bs *BitStream, v *generated.WPT_ChargeLoopRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelWPT_ChargeLoopRes: nil value")
	}
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	if err := bs.WriteBits(6, 58); err != nil {
		return err
	}
	return EncodeWPT_ChargeLoopRes(bs, v)
}

func EncodeWPT_ChargeLoopRes(bs *BitStream, v *generated.WPT_ChargeLoopRes) error {
	if v == nil {
		return fmt.Errorf("EncodeWPT_ChargeLoopRes: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// ResponseCode (required)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	responseCode := mapResponseCodeToEnum(v.ResponseCode)
	if err := bs.WriteBits(6, uint32(responseCode)); err != nil {
		return err
	}
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// EVSEProcessing (required)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	evseProc := mapEVSEProcessingToEnum(v.EVSEProcessing)
	if err := bs.WriteBits(2, uint32(evseProc)); err != nil {
		return err
	}
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Simplified - skip optional fields for now
	// END
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

func DecodeWPT_ChargeLoopRes(bs *BitStream) (*generated.WPT_ChargeLoopRes, error) {
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// ResponseCode
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	responseCodeVal, err := bs.ReadBits(6)
	if err != nil {
		return nil, err
	}
	responseCode := mapEnumToResponseCode(uint8(responseCodeVal))
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// EVSEProcessing
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	evseProcVal, err := bs.ReadBits(2)
	if err != nil {
		return nil, err
	}
	evseProcessing := mapEnumToEVSEProcessing(uint8(evseProcVal))
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Skip to END marker
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.WPT_ChargeLoopRes{
		Header:         *header,
		ResponseCode:   responseCode,
		EVSEProcessing: evseProcessing,
	}, nil
}
