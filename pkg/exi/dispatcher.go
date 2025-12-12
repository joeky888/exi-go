package exi

import (
	"fmt"

	"example.com/exi-go/pkg/v2g/generated"
)

// EncodeStruct encodes any supported ISO 15118-20 message struct to EXI bytes.
// This is a standalone function suitable for CGO bindings and external use.
func EncodeStruct(v interface{}) ([]byte, error) {
	buf := make([]byte, defaultEncodeBufferSize)
	bs := &BitStream{}
	bs.Init(buf, 0)

	switch val := v.(type) {
	// Session management messages
	case *generated.SessionSetupReq:
		if err := EncodeTopLevelSessionSetupReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.SessionSetupRes:
		if err := EncodeTopLevelSessionSetupRes(bs, val); err != nil {
			return nil, err
		}
	case *generated.SessionStopReq:
		if err := EncodeTopLevelSessionStopReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.SessionStopRes:
		if err := EncodeTopLevelSessionStopRes(bs, val); err != nil {
			return nil, err
		}

	// Service discovery and selection messages
	case *generated.ServiceDiscoveryReq:
		if err := EncodeTopLevelServiceDiscoveryReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.ServiceDiscoveryRes:
		if err := EncodeTopLevelServiceDiscoveryRes(bs, val); err != nil {
			return nil, err
		}
	case *generated.ServiceDetailReq:
		if err := EncodeTopLevelServiceDetailReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.ServiceDetailRes:
		if err := EncodeTopLevelServiceDetailRes(bs, val); err != nil {
			return nil, err
		}
	case *generated.ServiceSelectionReq:
		if err := EncodeTopLevelServiceSelectionReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.ServiceSelectionRes:
		if err := EncodeTopLevelServiceSelectionRes(bs, val); err != nil {
			return nil, err
		}

	// Authorization messages
	case *generated.AuthorizationReq:
		if err := EncodeTopLevelAuthorizationReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.AuthorizationRes:
		if err := EncodeTopLevelAuthorizationRes(bs, val); err != nil {
			return nil, err
		}
	case *generated.AuthorizationSetupReq:
		if err := EncodeTopLevelAuthorizationSetupReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.AuthorizationSetupRes:
		if err := EncodeTopLevelAuthorizationSetupRes(bs, val); err != nil {
			return nil, err
		}

	// Power delivery messages
	case *generated.PowerDeliveryReq:
		if err := EncodeTopLevelPowerDeliveryReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.PowerDeliveryRes:
		if err := EncodeTopLevelPowerDeliveryRes(bs, val); err != nil {
			return nil, err
		}

	// Schedule exchange messages
	case *generated.ScheduleExchangeReq:
		if err := EncodeTopLevelScheduleExchangeReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.ScheduleExchangeRes:
		if err := EncodeTopLevelScheduleExchangeRes(bs, val); err != nil {
			return nil, err
		}

	// Metering confirmation messages
	case *generated.MeteringConfirmationReq:
		if err := EncodeTopLevelMeteringConfirmationReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.MeteringConfirmationRes:
		if err := EncodeTopLevelMeteringConfirmationRes(bs, val); err != nil {
			return nil, err
		}

	// Certificate installation messages
	case *generated.CertificateInstallationReq:
		if err := EncodeTopLevelCertificateInstallationReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.CertificateInstallationRes:
		if err := EncodeTopLevelCertificateInstallationRes(bs, val); err != nil {
			return nil, err
		}

	// Vehicle check in/out messages
	case *generated.VehicleCheckInReq:
		if err := EncodeTopLevelVehicleCheckInReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.VehicleCheckInRes:
		if err := EncodeTopLevelVehicleCheckInRes(bs, val); err != nil {
			return nil, err
		}
	case *generated.VehicleCheckOutReq:
		if err := EncodeTopLevelVehicleCheckOutReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.VehicleCheckOutRes:
		if err := EncodeTopLevelVehicleCheckOutRes(bs, val); err != nil {
			return nil, err
		}

	// Control loop messages
	case *generated.CLReqControlMode:
		if err := EncodeTopLevelCLReqControlMode(bs, val); err != nil {
			return nil, err
		}
	case *generated.CLResControlMode:
		if err := EncodeTopLevelCLResControlMode(bs, val); err != nil {
			return nil, err
		}

	// Certificate update messages (from original implementation)
	case *generated.CertificateUpdateReq:
		if err := EncodeCertificateUpdateReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.CertificateUpdateRes:
		if err := EncodeCertificateUpdateRes(bs, val); err != nil {
			return nil, err
		}

	// WPT (Wireless Power Transfer) messages
	case *generated.WPT_AlignmentCheckReq:
		if err := EncodeTopLevelWPT_AlignmentCheckReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.WPT_AlignmentCheckRes:
		if err := EncodeTopLevelWPT_AlignmentCheckRes(bs, val); err != nil {
			return nil, err
		}
	case *generated.WPT_FinePositioningReq:
		if err := EncodeTopLevelWPT_FinePositioningReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.WPT_FinePositioningRes:
		if err := EncodeTopLevelWPT_FinePositioningRes(bs, val); err != nil {
			return nil, err
		}
	case *generated.WPT_ChargeLoopReq:
		if err := EncodeTopLevelWPT_ChargeLoopReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.WPT_ChargeLoopRes:
		if err := EncodeTopLevelWPT_ChargeLoopRes(bs, val); err != nil {
			return nil, err
		}

	// ACDP (AC Dynamic Power) messages
	case *generated.DC_ACDPReq:
		if err := EncodeTopLevelDC_ACDPReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.DC_ACDPRes:
		if err := EncodeTopLevelDC_ACDPRes(bs, val); err != nil {
			return nil, err
		}
	case *generated.DC_ACDP_BPTReq:
		if err := EncodeTopLevelDC_ACDP_BPTReq(bs, val); err != nil {
			return nil, err
		}
	case *generated.DC_ACDP_BPTRes:
		if err := EncodeTopLevelDC_ACDP_BPTRes(bs, val); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("EncodeStruct: unsupported type %T", v)
	}

	// Return written bytes
	outLen := bs.Length()
	if outLen > len(buf) {
		outLen = len(buf)
	}
	out := make([]byte, outLen)
	copy(out, buf[:outLen])
	return out, nil
}

// DecodeStruct decodes EXI bytes into the appropriate message struct type.
// The prototypeMsg parameter is used to determine the target type - pass an
// empty instance of the desired type (e.g., &generated.SessionSetupReq{}).
func DecodeStruct(data []byte, prototypeMsg interface{}) (interface{}, error) {
	bs := &BitStream{}
	bs.Init(data, 0)

	// Decode EXI header (8 bits)
	exiHeader, err := bs.ReadBits(8)
	if err != nil {
		return nil, fmt.Errorf("DecodeStruct: failed to read EXI header: %w", err)
	}
	if exiHeader != 0x80 {
		return nil, fmt.Errorf("DecodeStruct: invalid EXI header: expected 0x80, got 0x%02x", exiHeader)
	}

	// Decode event code (6 bits) to determine message type
	eventCode, err := bs.ReadBits(6)
	if err != nil {
		return nil, fmt.Errorf("DecodeStruct: failed to read event code: %w", err)
	}

	// Decode based on event code (the event codes match the ISO 15118-20 spec)
	switch eventCode {
	case 0: // AuthorizationReq
		return DecodeAuthorizationReq(bs)
	case 1: // AuthorizationRes
		return DecodeAuthorizationRes(bs)
	case 2: // AuthorizationSetupReq
		return DecodeAuthorizationSetupReq(bs)
	case 3: // AuthorizationSetupRes
		return DecodeAuthorizationSetupRes(bs)
	case 4: // CLReqControlMode
		return DecodeCLReqControlMode(bs)
	case 5: // CLResControlMode
		return DecodeCLResControlMode(bs)
	case 7: // CertificateInstallationReq
		return DecodeCertificateInstallationReq(bs)
	case 8: // CertificateInstallationRes
		return DecodeCertificateInstallationRes(bs)
	case 16: // MeteringConfirmationReq
		return DecodeMeteringConfirmationReq(bs)
	case 17: // MeteringConfirmationRes
		return DecodeMeteringConfirmationRes(bs)
	case 21: // PowerDeliveryReq
		return DecodePowerDeliveryReq(bs)
	case 22: // PowerDeliveryRes
		return DecodePowerDeliveryRes(bs)
	case 27: // ScheduleExchangeReq
		return DecodeScheduleExchangeReq(bs)
	case 28: // ScheduleExchangeRes
		return DecodeScheduleExchangeRes(bs)
	case 29: // ServiceDetailReq
		return DecodeServiceDetailReq(bs)
	case 30: // ServiceDetailRes
		return DecodeServiceDetailRes(bs)
	case 31: // ServiceDiscoveryReq
		return DecodeServiceDiscoveryReq(bs)
	case 32: // ServiceDiscoveryRes
		return DecodeServiceDiscoveryRes(bs)
	case 33: // ServiceSelectionReq
		return DecodeServiceSelectionReq(bs)
	case 34: // ServiceSelectionRes
		return DecodeServiceSelectionRes(bs)
	case 35: // SessionSetupReq
		return DecodeSessionSetupReq(bs)
	case 36: // SessionSetupRes
		return DecodeSessionSetupRes(bs)
	case 37: // SessionStopReq
		return DecodeSessionStopReq(bs)
	case 38: // SessionStopRes
		return DecodeSessionStopRes(bs)
	case 50: // VehicleCheckInRes
		return DecodeVehicleCheckInRes(bs)
	case 45: // VehicleCheckOutRes (legacy code, should be 52)
		return DecodeVehicleCheckOutRes(bs)
	case 49: // VehicleCheckInReq
		return DecodeVehicleCheckInReq(bs)
	case 51: // VehicleCheckOutReq
		return DecodeVehicleCheckOutReq(bs)
	case 52: // VehicleCheckOutRes
		return DecodeVehicleCheckOutRes(bs)

	// WPT (Wireless Power Transfer) messages
	case 53: // WPT_AlignmentCheckReq
		return DecodeWPT_AlignmentCheckReq(bs)
	case 54: // WPT_AlignmentCheckRes
		return DecodeWPT_AlignmentCheckRes(bs)
	case 55: // WPT_FinePositioningReq
		return DecodeWPT_FinePositioningReq(bs)
	case 56: // WPT_FinePositioningRes
		return DecodeWPT_FinePositioningRes(bs)
	case 57: // WPT_ChargeLoopReq
		return DecodeWPT_ChargeLoopReq(bs)
	case 58: // WPT_ChargeLoopRes
		return DecodeWPT_ChargeLoopRes(bs)

	// ACDP (AC Dynamic Power) messages
	case 59: // DC_ACDPReq
		return DecodeDC_ACDPReq(bs)
	case 60: // DC_ACDPRes
		return DecodeDC_ACDPRes(bs)
	case 61: // DC_ACDP_BPTReq
		return DecodeDC_ACDP_BPTReq(bs)
	case 62: // DC_ACDP_BPTRes
		return DecodeDC_ACDP_BPTRes(bs)

	default:
		return nil, fmt.Errorf("DecodeStruct: unsupported event code %d", eventCode)
	}
}
