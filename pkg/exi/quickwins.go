package exi

import (
	"fmt"

	"example.com/exi-go/pkg/v2g/generated"
)

// Quick-win encoders/decoders for three simple message types:
// - PowerDeliveryRes (event 22)           : Header + ResponseCode + optional EVSEStatus
// - VehicleCheckInRes (event 50)          : Header + ResponseCode + optional VehicleCheckInResult
// - VehicleCheckOutRes (event 52)         : Header + ResponseCode + EVSECheckOutStatus
//
// These follow established patterns in this package: write an EXI simple header
// (0x80), a 6-bit event code, then encode the message body using the shared
// MessageHeaderType helpers and small-field helpers (writeString/readString,
// writeUint16/readUint16, etc).

// --------------------------- PowerDeliveryRes ------------------------------

// EncodeTopLevelPowerDeliveryRes writes EXI header + event code (22) then
// delegates to EncodePowerDeliveryRes.
func EncodeTopLevelPowerDeliveryRes(bs *BitStream, v *generated.PowerDeliveryRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelPowerDeliveryRes: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for PowerDeliveryRes is 22
	if err := bs.WriteBits(6, 22); err != nil {
		return err
	}
	return EncodePowerDeliveryRes(bs, v)
}

// EncodePowerDeliveryRes encodes the PowerDeliveryRes body:
// Header, ResponseCode, optional EVSEStatus.
func EncodePowerDeliveryRes(bs *BitStream, v *generated.PowerDeliveryRes) error {
	if v == nil {
		return fmt.Errorf("EncodePowerDeliveryRes: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// START ResponseCode
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ResponseCode encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// 6-bit enum
	responseCode := mapResponseCodeToEnum(v.ResponseCode)
	if err := bs.WriteBits(6, uint32(responseCode)); err != nil {
		return err
	}
	// END ResponseCode
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Optional EVSEStatus presence (1 bit: 1 = present, 0 = absent)
	if v.EVSEStatus != nil {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
		// START EVSEStatus (struct)
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		// NotificationMaxDelay (uint16)
		if err := writeUint16(bs, v.EVSEStatus.NotificationMaxDelay); err != nil {
			return err
		}
		// EVSENotification (string)
		// Use the string helper which writes a 16-bit length + bytes
		if err := writeString(bs, v.EVSEStatus.EVSENotification); err != nil {
			return err
		}
		// END EVSEStatus
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	} else {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}

	// END PowerDeliveryRes
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodePowerDeliveryRes decodes the PowerDeliveryRes body.
func DecodePowerDeliveryRes(bs *BitStream) (*generated.PowerDeliveryRes, error) {
	// START Header
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// START ResponseCode
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ResponseCode encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	rcEnum, err := bs.ReadBits(6)
	if err != nil {
		return nil, err
	}
	responseCode := mapEnumToResponseCode(uint8(rcEnum))
	// END ResponseCode
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Optional EVSEStatus presence
	p, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	var evseStatus *generated.EVSEStatus
	if p == 1 {
		// START EVSEStatus
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		// NotificationMaxDelay
		nmd, err := readUint16(bs)
		if err != nil {
			return nil, err
		}
		// EVSENotification
		evsen, err := readString(bs)
		if err != nil {
			return nil, err
		}
		// END EVSEStatus
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		evseStatus = &generated.EVSEStatus{
			NotificationMaxDelay: nmd,
			EVSENotification:     evsen,
		}
	}

	// END PowerDeliveryRes
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.PowerDeliveryRes{
		Header:       *header,
		ResponseCode: responseCode,
		EVSEStatus:   evseStatus,
	}, nil
}

// ------------------------- VehicleCheckInRes -------------------------------

// EncodeTopLevelVehicleCheckInRes writes EXI header + event code (43) then
// delegates to EncodeVehicleCheckInRes.
func EncodeTopLevelVehicleCheckInRes(bs *BitStream, v *generated.VehicleCheckInRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelVehicleCheckInRes: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for VehicleCheckInRes is 50
	if err := bs.WriteBits(6, 50); err != nil {
		return err
	}
	return EncodeVehicleCheckInRes(bs, v)
}

// EncodeVehicleCheckInRes encodes Header, ResponseCode and optional
// VehicleCheckInResult.
func EncodeVehicleCheckInRes(bs *BitStream, v *generated.VehicleCheckInRes) error {
	if v == nil {
		return fmt.Errorf("EncodeVehicleCheckInRes: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// START ResponseCode
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	code := mapResponseCodeToEnum(v.ResponseCode)
	if err := bs.WriteBits(6, uint32(code)); err != nil {
		return err
	}
	// END ResponseCode
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Optional VehicleCheckInResult (1 bit presence)
	if v.VehicleCheckInResult != nil {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
		// Write the string value
		if err := writeString(bs, *v.VehicleCheckInResult); err != nil {
			return err
		}
	} else {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}

	// END VehicleCheckInRes
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	return nil
}

// DecodeVehicleCheckInRes decodes the VehicleCheckInRes body.
func DecodeVehicleCheckInRes(bs *BitStream) (*generated.VehicleCheckInRes, error) {
	// START Header
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// START ResponseCode
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	rcEnum, err := bs.ReadBits(6)
	if err != nil {
		return nil, err
	}
	responseCode := mapEnumToResponseCode(uint8(rcEnum))
	// END ResponseCode
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Optional VehicleCheckInResult presence
	p, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	var result *string
	if p == 1 {
		r, err := readString(bs)
		if err != nil {
			return nil, err
		}
		result = &r
	}

	// END VehicleCheckInRes
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.VehicleCheckInRes{
		Header:               *header,
		ResponseCode:         responseCode,
		VehicleCheckInResult: result,
	}, nil
}

// ------------------------ VehicleCheckOutRes -------------------------------

// EncodeTopLevelVehicleCheckOutRes writes EXI header + event code (45) then
// delegates to EncodeVehicleCheckOutRes.
func EncodeTopLevelVehicleCheckOutRes(bs *BitStream, v *generated.VehicleCheckOutRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelVehicleCheckOutRes: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for VehicleCheckOutRes is 52
	if err := bs.WriteBits(6, 52); err != nil {
		return err
	}
	return EncodeVehicleCheckOutRes(bs, v)
}

// EncodeVehicleCheckOutRes encodes Header, ResponseCode and required
// EVSECheckOutStatus.
func EncodeVehicleCheckOutRes(bs *BitStream, v *generated.VehicleCheckOutRes) error {
	if v == nil {
		return fmt.Errorf("EncodeVehicleCheckOutRes: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// START ResponseCode
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	code := mapResponseCodeToEnum(v.ResponseCode)
	if err := bs.WriteBits(6, uint32(code)); err != nil {
		return err
	}
	// END ResponseCode
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// EVSECheckOutStatus (required string)
	// Use START marker for structural symmetry
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := writeString(bs, v.EVSECheckOutStatus); err != nil {
		return err
	}
	// END EVSECheckOutStatus
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// END VehicleCheckOutRes
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	return nil
}

// DecodeVehicleCheckOutRes decodes the VehicleCheckOutRes body.
func DecodeVehicleCheckOutRes(bs *BitStream) (*generated.VehicleCheckOutRes, error) {
	// START Header
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// START ResponseCode
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	rcEnum, err := bs.ReadBits(6)
	if err != nil {
		return nil, err
	}
	responseCode := mapEnumToResponseCode(uint8(rcEnum))
	// END ResponseCode
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// EVSECheckOutStatus START
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	status, err := readString(bs)
	if err != nil {
		return nil, err
	}
	// EVSECheckOutStatus END
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// END VehicleCheckOutRes
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.VehicleCheckOutRes{
		Header:             *header,
		ResponseCode:       responseCode,
		EVSECheckOutStatus: status,
	}, nil
}
