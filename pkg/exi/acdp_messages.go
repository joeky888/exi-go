package exi

import (
	"fmt"

	"example.com/exi-go/pkg/v2g/generated"
)

// ACDP (AC Dynamic Power) message type implementations for ISO 15118-20.
// These messages support DC charging with AC-side dynamic power control
// and bidirectional power transfer (BPT) capabilities.

// mapEVProcessingToEnum converts a string EVProcessing to its enum value (2-bit)
func mapEVProcessingToEnum(processing string) uint8 {
	switch processing {
	case "Finished":
		return 0
	case "Ongoing":
		return 1
	default:
		return 1 // Default to Ongoing
	}
}

// mapEnumToEVProcessing converts an enum value to its string EVProcessing
func mapEnumToEVProcessing(processing uint8) string {
	switch processing {
	case 0:
		return "Finished"
	case 1:
		return "Ongoing"
	default:
		return "Ongoing"
	}
}

// ========================== DC_ACDPReq ==========================

// EncodeTopLevelDC_ACDPReq writes EXI header + event code (59) then
// delegates to EncodeDC_ACDPReq.
func EncodeTopLevelDC_ACDPReq(bs *BitStream, v *generated.DC_ACDPReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelDC_ACDPReq: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for DC_ACDPReq is 59
	if err := bs.WriteBits(6, 59); err != nil {
		return err
	}
	return EncodeDC_ACDPReq(bs, v)
}

// EncodeDC_ACDPReq encodes the DC_ACDPReq body.
func EncodeDC_ACDPReq(bs *BitStream, v *generated.DC_ACDPReq) error {
	if v == nil {
		return fmt.Errorf("EncodeDC_ACDPReq: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// EVProcessing (required)
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

	// EVTargetEnergyRequest (required)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeRationalNumber(bs, &v.EVTargetEnergyRequest); err != nil {
		return err
	}

	// Optional fields - simplified for now
	// END
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeDC_ACDPReq decodes the DC_ACDPReq body.
func DecodeDC_ACDPReq(bs *BitStream) (*generated.DC_ACDPReq, error) {
	// START Header
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

	// EVTargetEnergyRequest
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	targetEnergy, err := decodeRationalNumber(bs)
	if err != nil {
		return nil, err
	}

	// Skip to END marker
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.DC_ACDPReq{
		Header:                *header,
		EVProcessing:          evProcessing,
		EVTargetEnergyRequest: *targetEnergy,
	}, nil
}

// ========================== DC_ACDPRes ==========================

// EncodeTopLevelDC_ACDPRes writes EXI header + event code (60).
func EncodeTopLevelDC_ACDPRes(bs *BitStream, v *generated.DC_ACDPRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelDC_ACDPRes: nil value")
	}
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	if err := bs.WriteBits(6, 60); err != nil {
		return err
	}
	return EncodeDC_ACDPRes(bs, v)
}

func EncodeDC_ACDPRes(bs *BitStream, v *generated.DC_ACDPRes) error {
	if v == nil {
		return fmt.Errorf("EncodeDC_ACDPRes: nil value")
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

	// Optional fields - simplified
	// END
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

func DecodeDC_ACDPRes(bs *BitStream) (*generated.DC_ACDPRes, error) {
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

	return &generated.DC_ACDPRes{
		Header:         *header,
		ResponseCode:   responseCode,
		EVSEProcessing: evseProcessing,
	}, nil
}

// ========================== DC_ACDP_BPTReq ==========================

// EncodeTopLevelDC_ACDP_BPTReq writes EXI header + event code (61).
func EncodeTopLevelDC_ACDP_BPTReq(bs *BitStream, v *generated.DC_ACDP_BPTReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelDC_ACDP_BPTReq: nil value")
	}
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	if err := bs.WriteBits(6, 61); err != nil {
		return err
	}
	return EncodeDC_ACDP_BPTReq(bs, v)
}

func EncodeDC_ACDP_BPTReq(bs *BitStream, v *generated.DC_ACDP_BPTReq) error {
	if v == nil {
		return fmt.Errorf("EncodeDC_ACDP_BPTReq: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// EVProcessing (required)
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

	// EVTargetEnergyRequest (required)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeRationalNumber(bs, &v.EVTargetEnergyRequest); err != nil {
		return err
	}

	// Optional fields - simplified
	// END
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

func DecodeDC_ACDP_BPTReq(bs *BitStream) (*generated.DC_ACDP_BPTReq, error) {
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

	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	targetEnergy, err := decodeRationalNumber(bs)
	if err != nil {
		return nil, err
	}

	// Skip to END marker
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.DC_ACDP_BPTReq{
		Header:                *header,
		EVProcessing:          evProcessing,
		EVTargetEnergyRequest: *targetEnergy,
	}, nil
}

// ========================== DC_ACDP_BPTRes ==========================

// EncodeTopLevelDC_ACDP_BPTRes writes EXI header + event code (62).
func EncodeTopLevelDC_ACDP_BPTRes(bs *BitStream, v *generated.DC_ACDP_BPTRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelDC_ACDP_BPTRes: nil value")
	}
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	if err := bs.WriteBits(6, 62); err != nil {
		return err
	}
	return EncodeDC_ACDP_BPTRes(bs, v)
}

func EncodeDC_ACDP_BPTRes(bs *BitStream, v *generated.DC_ACDP_BPTRes) error {
	if v == nil {
		return fmt.Errorf("EncodeDC_ACDP_BPTRes: nil value")
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

	// Optional fields - simplified
	// END
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

func DecodeDC_ACDP_BPTRes(bs *BitStream) (*generated.DC_ACDP_BPTRes, error) {
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

	return &generated.DC_ACDP_BPTRes{
		Header:         *header,
		ResponseCode:   responseCode,
		EVSEProcessing: evseProcessing,
	}, nil
}
