package exi

import (
	"fmt"

	"example.com/exi-go/pkg/v2g/generated"
)

// Phase 2 moderate complexity message type implementations.
// These message types involve nested structures, arrays, and optional union types.
//
// AuthorizationSetupRes (event 3): Header + ResponseCode + AuthorizationServices array
//                                   + CertificateInstallationService bool
//                                   + optional EIM/PnC authorization modes

// ------------------------ AuthorizationSetupRes ----------------------------

// EncodeTopLevelAuthorizationSetupRes writes EXI header + event code (3) then
// delegates to EncodeAuthorizationSetupRes.
func EncodeTopLevelAuthorizationSetupRes(bs *BitStream, v *generated.AuthorizationSetupRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelAuthorizationSetupRes: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for AuthorizationSetupRes is 3
	if err := bs.WriteBits(6, 3); err != nil {
		return err
	}
	return EncodeAuthorizationSetupRes(bs, v)
}

// EncodeAuthorizationSetupRes encodes the AuthorizationSetupRes body:
// Header, ResponseCode, AuthorizationServices, CertificateInstallationService,
// and optional EIM_ASResAuthorizationMode or PnC_ASResAuthorizationMode.
func EncodeAuthorizationSetupRes(bs *BitStream, v *generated.AuthorizationSetupRes) error {
	if v == nil {
		return fmt.Errorf("EncodeAuthorizationSetupRes: nil value")
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

	// START AuthorizationServices (required array of strings)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write array length
	if err := bs.WriteUnsignedVar(uint64(len(v.AuthorizationServices))); err != nil {
		return err
	}
	// Encode each service string
	for _, service := range v.AuthorizationServices {
		// START service
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := writeString(bs, service); err != nil {
			return err
		}
		// END service
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	// END AuthorizationServices
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// START CertificateInstallationService (required boolean)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Boolean encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Boolean value (1 bit)
	if v.CertificateInstallationService {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
	} else {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	// END CertificateInstallationService
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Optional EIM_ASResAuthorizationMode or PnC_ASResAuthorizationMode
	// This is a choice between two optional modes
	// 2-bit choice: 0 = EIM mode, 1 = PnC mode, 2 = END (no mode)
	if v.EIM_ASResAuthorizationMode != nil {
		if err := bs.WriteBits(2, 0); err != nil {
			return err
		}
		// START EIM_ASResAuthorizationMode
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodeEIM_ASResAuthorizationMode(bs, v.EIM_ASResAuthorizationMode); err != nil {
			return err
		}
		// END EIM_ASResAuthorizationMode
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		// END message
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	} else if v.PnC_ASResAuthorizationMode != nil {
		if err := bs.WriteBits(2, 1); err != nil {
			return err
		}
		// START PnC_ASResAuthorizationMode
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodePnC_ASResAuthorizationMode(bs, v.PnC_ASResAuthorizationMode); err != nil {
			return err
		}
		// END PnC_ASResAuthorizationMode
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		// END message
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	} else {
		// No optional mode, just END
		if err := bs.WriteBits(2, 2); err != nil {
			return err
		}
	}

	return nil
}

// encodeEIM_ASResAuthorizationMode encodes an EIM authorization mode.
// Currently a placeholder as the type has no fields defined.
func encodeEIM_ASResAuthorizationMode(bs *BitStream, v *generated.EIM_ASResAuthorizationMode) error {
	// Placeholder - no fields to encode
	// If fields are added to the type, encode them here
	return nil
}

// encodePnC_ASResAuthorizationMode encodes a PnC authorization mode.
func encodePnC_ASResAuthorizationMode(bs *BitStream, v *generated.PnC_ASResAuthorizationMode) error {
	// START GenChallenge (required binary data)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// hexBinary encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write length
	if err := writeUint16(bs, uint16(len(v.GenChallenge))); err != nil {
		return err
	}
	// Write bytes
	if err := writeRawBytes(bs, v.GenChallenge); err != nil {
		return err
	}
	// END GenChallenge
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// START SupportedProviders (required array of strings)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write array length
	if err := bs.WriteUnsignedVar(uint64(len(v.SupportedProviders))); err != nil {
		return err
	}
	// Encode each provider string
	for _, provider := range v.SupportedProviders {
		// START provider
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := writeString(bs, provider); err != nil {
			return err
		}
		// END provider
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	// END SupportedProviders
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeAuthorizationSetupRes decodes the AuthorizationSetupRes body.
func DecodeAuthorizationSetupRes(bs *BitStream) (*generated.AuthorizationSetupRes, error) {
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

	// START AuthorizationServices
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Read array length
	count, err := bs.ReadUnsignedVar()
	if err != nil {
		return nil, err
	}
	authServices := make([]string, count)
	for i := uint64(0); i < count; i++ {
		// START service
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		service, err := readString(bs)
		if err != nil {
			return nil, err
		}
		authServices[i] = service
		// END service
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}
	// END AuthorizationServices
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// START CertificateInstallationService
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Boolean encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Boolean value
	certInstallBit, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	certInstallService := certInstallBit == 1
	// END CertificateInstallationService
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Optional mode choice (2 bits)
	choice, err := bs.ReadBits(2)
	if err != nil {
		return nil, err
	}

	var eimMode *generated.EIM_ASResAuthorizationMode
	var pncMode *generated.PnC_ASResAuthorizationMode

	if choice == 0 {
		// EIM mode
		// START EIM_ASResAuthorizationMode
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		eim, err := decodeEIM_ASResAuthorizationMode(bs)
		if err != nil {
			return nil, err
		}
		eimMode = eim
		// END EIM_ASResAuthorizationMode
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		// END message
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	} else if choice == 1 {
		// PnC mode
		// START PnC_ASResAuthorizationMode
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		pnc, err := decodePnC_ASResAuthorizationMode(bs)
		if err != nil {
			return nil, err
		}
		pncMode = pnc
		// END PnC_ASResAuthorizationMode
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		// END message
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}
	// else choice == 2, no mode, message already ended

	return &generated.AuthorizationSetupRes{
		Header:                         *header,
		ResponseCode:                   responseCode,
		AuthorizationServices:          authServices,
		CertificateInstallationService: certInstallService,
		EIM_ASResAuthorizationMode:     eimMode,
		PnC_ASResAuthorizationMode:     pncMode,
	}, nil
}

// decodeEIM_ASResAuthorizationMode decodes an EIM authorization mode.
func decodeEIM_ASResAuthorizationMode(bs *BitStream) (*generated.EIM_ASResAuthorizationMode, error) {
	// Placeholder - no fields to decode
	return &generated.EIM_ASResAuthorizationMode{}, nil
}

// decodePnC_ASResAuthorizationMode decodes a PnC authorization mode.
func decodePnC_ASResAuthorizationMode(bs *BitStream) (*generated.PnC_ASResAuthorizationMode, error) {
	// START GenChallenge
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// hexBinary encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Read length
	genChallengeLen, err := readUint16(bs)
	if err != nil {
		return nil, err
	}
	// Read bytes
	genChallenge := make([]byte, genChallengeLen)
	for i := 0; i < int(genChallengeLen); i++ {
		b, err := bs.ReadOctet()
		if err != nil {
			return nil, err
		}
		genChallenge[i] = b
	}
	// END GenChallenge
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// START SupportedProviders
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Read array length
	count, err := bs.ReadUnsignedVar()
	if err != nil {
		return nil, err
	}
	providers := make([]string, count)
	for i := uint64(0); i < count; i++ {
		// START provider
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		provider, err := readString(bs)
		if err != nil {
			return nil, err
		}
		providers[i] = provider
		// END provider
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}
	// END SupportedProviders
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.PnC_ASResAuthorizationMode{
		GenChallenge:       genChallenge,
		SupportedProviders: providers,
	}, nil
}

// ------------------------ ScheduleExchangeReq ------------------------------

// EncodeTopLevelScheduleExchangeReq writes EXI header + event code (27) then
// delegates to EncodeScheduleExchangeReq.
func EncodeTopLevelScheduleExchangeReq(bs *BitStream, v *generated.ScheduleExchangeReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelScheduleExchangeReq: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for ScheduleExchangeReq is 27
	if err := bs.WriteBits(6, 27); err != nil {
		return err
	}
	return EncodeScheduleExchangeReq(bs, v)
}

// EncodeScheduleExchangeReq encodes the ScheduleExchangeReq body:
// Header and MaximumSupportingPoints.
func EncodeScheduleExchangeReq(bs *BitStream, v *generated.ScheduleExchangeReq) error {
	if v == nil {
		return fmt.Errorf("EncodeScheduleExchangeReq: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// START MaximumSupportingPoints (required uint16)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// MaximumSupportingPoints encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write as uint16
	if err := writeUint16(bs, v.MaximumSupportingPoints); err != nil {
		return err
	}
	// END MaximumSupportingPoints
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// END ScheduleExchangeReq
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeScheduleExchangeReq decodes the ScheduleExchangeReq body.
func DecodeScheduleExchangeReq(bs *BitStream) (*generated.ScheduleExchangeReq, error) {
	// START Header
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// START MaximumSupportingPoints
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// MaximumSupportingPoints encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	maxSupportingPoints, err := readUint16(bs)
	if err != nil {
		return nil, err
	}
	// END MaximumSupportingPoints
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// END ScheduleExchangeReq
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.ScheduleExchangeReq{
		Header:                  *header,
		MaximumSupportingPoints: maxSupportingPoints,
	}, nil
}

// ------------------------ ScheduleExchangeRes ------------------------------

// EncodeTopLevelScheduleExchangeRes writes EXI header + event code (28) then
// delegates to EncodeScheduleExchangeRes.
func EncodeTopLevelScheduleExchangeRes(bs *BitStream, v *generated.ScheduleExchangeRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelScheduleExchangeRes: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for ScheduleExchangeRes is 28
	if err := bs.WriteBits(6, 28); err != nil {
		return err
	}
	return EncodeScheduleExchangeRes(bs, v)
}

// EncodeScheduleExchangeRes encodes the ScheduleExchangeRes body:
// Header, ResponseCode, and EVSEProcessing.
func EncodeScheduleExchangeRes(bs *BitStream, v *generated.ScheduleExchangeRes) error {
	if v == nil {
		return fmt.Errorf("EncodeScheduleExchangeRes: nil value")
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

	// START EVSEProcessing (required string/enum)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// EVSEProcessing encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// 2-bit enum (Finished=0, Ongoing=1)
	evseProcessing := mapEVSEProcessingToEnum(v.EVSEProcessing)
	if err := bs.WriteBits(2, uint32(evseProcessing)); err != nil {
		return err
	}
	// END EVSEProcessing
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// END ScheduleExchangeRes
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeScheduleExchangeRes decodes the ScheduleExchangeRes body.
func DecodeScheduleExchangeRes(bs *BitStream) (*generated.ScheduleExchangeRes, error) {
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

	// START EVSEProcessing
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// EVSEProcessing encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	evseProcessingEnum, err := bs.ReadBits(2)
	if err != nil {
		return nil, err
	}
	evseProcessing := mapEnumToEVSEProcessing(uint8(evseProcessingEnum))
	// END EVSEProcessing
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// END ScheduleExchangeRes
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.ScheduleExchangeRes{
		Header:         *header,
		ResponseCode:   responseCode,
		EVSEProcessing: evseProcessing,
	}, nil
}
