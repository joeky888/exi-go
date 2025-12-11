package exi

import (
	"fmt"

	"example.com/exi-go/pkg/v2g/generated"
)

// Phase 3 complex message type implementations.
// These message types involve complex nested structures, union types, and certificate chains.
//
// PowerDeliveryReq (event 21): Header + EVProcessing + ChargeProgress + optional EVPowerProfile + optional BPT_ChannelSelection
// AuthorizationReq (event 0): Header + SelectedAuthorizationService + union of EIM/PnC modes
// ServiceDetailRes (event 30): Header + ResponseCode + ServiceID + ServiceParameterList

// --------------------------- PowerDeliveryReq ------------------------------

// EncodeTopLevelPowerDeliveryReq writes EXI header + event code (21) then
// delegates to EncodePowerDeliveryReq.
func EncodeTopLevelPowerDeliveryReq(bs *BitStream, v *generated.PowerDeliveryReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelPowerDeliveryReq: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for PowerDeliveryReq is 21
	if err := bs.WriteBits(6, 21); err != nil {
		return err
	}
	return EncodePowerDeliveryReq(bs, v)
}

// EncodePowerDeliveryReq encodes the PowerDeliveryReq body:
// Header, EVProcessing, ChargeProgress, optional EVPowerProfile, and optional BPT_ChannelSelection.
func EncodePowerDeliveryReq(bs *BitStream, v *generated.PowerDeliveryReq) error {
	if v == nil {
		return fmt.Errorf("EncodePowerDeliveryReq: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// START EVProcessing (required string/enum)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// EVProcessing encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// 2-bit enum (Finished=0, Ongoing=1)
	evProcessing := mapEVSEProcessingToEnum(v.EVProcessing)
	if err := bs.WriteBits(2, uint32(evProcessing)); err != nil {
		return err
	}
	// END EVProcessing
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// START ChargeProgress (required string)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := writeString(bs, v.ChargeProgress); err != nil {
		return err
	}
	// END ChargeProgress
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Optional EVPowerProfile and BPT_ChannelSelection
	// 3-bit choice: 0=EVPowerProfile only, 1=BPT_ChannelSelection only, 2=both, 3=END
	if v.EVPowerProfile != nil && v.BPT_ChannelSelection != nil {
		// Both present
		if err := bs.WriteBits(3, 0); err != nil {
			return err
		}
		// START EVPowerProfile
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodeEVPowerProfile(bs, v.EVPowerProfile); err != nil {
			return err
		}
		// END EVPowerProfile
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		// START BPT_ChannelSelection
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := writeString(bs, *v.BPT_ChannelSelection); err != nil {
			return err
		}
		// END BPT_ChannelSelection
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	} else if v.EVPowerProfile != nil {
		// EVPowerProfile only
		if err := bs.WriteBits(3, 1); err != nil {
			return err
		}
		// START EVPowerProfile
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodeEVPowerProfile(bs, v.EVPowerProfile); err != nil {
			return err
		}
		// END EVPowerProfile
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	} else if v.BPT_ChannelSelection != nil {
		// BPT_ChannelSelection only
		if err := bs.WriteBits(3, 2); err != nil {
			return err
		}
		// START BPT_ChannelSelection
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := writeString(bs, *v.BPT_ChannelSelection); err != nil {
			return err
		}
		// END BPT_ChannelSelection
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	} else {
		// Neither present
		if err := bs.WriteBits(3, 3); err != nil {
			return err
		}
	}

	// END PowerDeliveryReq
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// encodeEVPowerProfile encodes an EVPowerProfile.
func encodeEVPowerProfile(bs *BitStream, v *generated.EVPowerProfile) error {
	// START TimeAnchor (required uint64)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// TimeAnchor encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write timestamp as unsigned-var
	if err := bs.WriteUnsignedVar(v.TimeAnchor); err != nil {
		return err
	}
	// END TimeAnchor
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// START Entries (required array)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write array length
	if err := bs.WriteUnsignedVar(uint64(len(v.Entries))); err != nil {
		return err
	}
	// Encode each entry
	for _, entry := range v.Entries {
		// START entry
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodeEVPowerProfileEntry(bs, &entry); err != nil {
			return err
		}
		// END entry
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	// END Entries
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// encodeEVPowerProfileEntry encodes an EVPowerProfileEntry.
func encodeEVPowerProfileEntry(bs *BitStream, v *generated.EVPowerProfileEntry) error {
	// START Duration (required uint32)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Duration encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write as unsigned-var
	if err := bs.WriteUnsignedVar(uint64(v.Duration)); err != nil {
		return err
	}
	// END Duration
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// START Power (required RationalNumber)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeRationalNumber(bs, &v.Power); err != nil {
		return err
	}
	// END Power
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// encodeRationalNumber encodes a RationalNumber.
func encodeRationalNumber(bs *BitStream, v *generated.RationalNumber) error {
	// START Exponent (required int8)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Exponent encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write as signed byte (8 bits)
	if err := bs.WriteBits(8, uint32(uint8(v.Exponent))); err != nil {
		return err
	}
	// END Exponent
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// START Value (required int16)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Value encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write as signed 16-bit (using writeUint16 for unsigned representation)
	if err := writeUint16(bs, uint16(v.Value)); err != nil {
		return err
	}
	// END Value
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodePowerDeliveryReq decodes the PowerDeliveryReq body.
func DecodePowerDeliveryReq(bs *BitStream) (*generated.PowerDeliveryReq, error) {
	// START Header
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// START EVProcessing
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// EVProcessing encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	evProcessingEnum, err := bs.ReadBits(2)
	if err != nil {
		return nil, err
	}
	evProcessing := mapEnumToEVSEProcessing(uint8(evProcessingEnum))
	// END EVProcessing
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// START ChargeProgress
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	chargeProgress, err := readString(bs)
	if err != nil {
		return nil, err
	}
	// END ChargeProgress
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Optional choice (3 bits)
	choice, err := bs.ReadBits(3)
	if err != nil {
		return nil, err
	}

	var evPowerProfile *generated.EVPowerProfile
	var bptChannelSelection *string

	if choice == 0 {
		// Both present
		// START EVPowerProfile
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		profile, err := decodeEVPowerProfile(bs)
		if err != nil {
			return nil, err
		}
		evPowerProfile = profile
		// END EVPowerProfile
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		// START BPT_ChannelSelection
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		bpt, err := readString(bs)
		if err != nil {
			return nil, err
		}
		bptChannelSelection = &bpt
		// END BPT_ChannelSelection
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	} else if choice == 1 {
		// EVPowerProfile only
		// START EVPowerProfile
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		profile, err := decodeEVPowerProfile(bs)
		if err != nil {
			return nil, err
		}
		evPowerProfile = profile
		// END EVPowerProfile
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	} else if choice == 2 {
		// BPT_ChannelSelection only
		// START BPT_ChannelSelection
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		bpt, err := readString(bs)
		if err != nil {
			return nil, err
		}
		bptChannelSelection = &bpt
		// END BPT_ChannelSelection
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}
	// else choice == 3, neither present

	// END PowerDeliveryReq
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.PowerDeliveryReq{
		Header:               *header,
		EVProcessing:         evProcessing,
		ChargeProgress:       chargeProgress,
		EVPowerProfile:       evPowerProfile,
		BPT_ChannelSelection: bptChannelSelection,
	}, nil
}

// decodeEVPowerProfile decodes an EVPowerProfile.
func decodeEVPowerProfile(bs *BitStream) (*generated.EVPowerProfile, error) {
	// START TimeAnchor
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// TimeAnchor encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	timeAnchor, err := bs.ReadUnsignedVar()
	if err != nil {
		return nil, err
	}
	// END TimeAnchor
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// START Entries
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Read array length
	count, err := bs.ReadUnsignedVar()
	if err != nil {
		return nil, err
	}
	entries := make([]generated.EVPowerProfileEntry, count)
	for i := uint64(0); i < count; i++ {
		// START entry
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		entry, err := decodeEVPowerProfileEntry(bs)
		if err != nil {
			return nil, err
		}
		entries[i] = *entry
		// END entry
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}
	// END Entries
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.EVPowerProfile{
		TimeAnchor: timeAnchor,
		Entries:    entries,
	}, nil
}

// decodeEVPowerProfileEntry decodes an EVPowerProfileEntry.
func decodeEVPowerProfileEntry(bs *BitStream) (*generated.EVPowerProfileEntry, error) {
	// START Duration
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Duration encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	duration, err := bs.ReadUnsignedVar()
	if err != nil {
		return nil, err
	}
	// END Duration
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// START Power
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	power, err := decodeRationalNumber(bs)
	if err != nil {
		return nil, err
	}
	// END Power
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.EVPowerProfileEntry{
		Duration: uint32(duration),
		Power:    *power,
	}, nil
}

// decodeRationalNumber decodes a RationalNumber.
func decodeRationalNumber(bs *BitStream) (*generated.RationalNumber, error) {
	// START Exponent
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Exponent encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	exponentBits, err := bs.ReadBits(8)
	if err != nil {
		return nil, err
	}
	exponent := int8(exponentBits)
	// END Exponent
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// START Value
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Value encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	valueUint, err := readUint16(bs)
	if err != nil {
		return nil, err
	}
	value := int16(valueUint)
	// END Value
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.RationalNumber{
		Exponent: exponent,
		Value:    value,
	}, nil
}

// --------------------------- AuthorizationReq ------------------------------

// EncodeTopLevelAuthorizationReq writes EXI header + event code (0) then
// delegates to EncodeAuthorizationReq.
func EncodeTopLevelAuthorizationReq(bs *BitStream, v *generated.AuthorizationReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelAuthorizationReq: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for AuthorizationReq is 0
	if err := bs.WriteBits(6, 0); err != nil {
		return err
	}
	return EncodeAuthorizationReq(bs, v)
}

// EncodeAuthorizationReq encodes the AuthorizationReq body:
// Header, SelectedAuthorizationService, and union of EIM/PnC modes.
func EncodeAuthorizationReq(bs *BitStream, v *generated.AuthorizationReq) error {
	if v == nil {
		return fmt.Errorf("EncodeAuthorizationReq: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// START SelectedAuthorizationService (required string)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := writeString(bs, v.SelectedAuthorizationService); err != nil {
		return err
	}
	// END SelectedAuthorizationService
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Union choice: EIM_AReqAuthorizationMode or PnC_AReqAuthorizationMode
	// 2-bit choice: 0 = EIM mode, 1 = PnC mode, 2 = END (no mode)
	if v.EIM_AReqAuthorizationMode != nil {
		if err := bs.WriteBits(2, 0); err != nil {
			return err
		}
		// START EIM_AReqAuthorizationMode
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodeEIM_AReqAuthorizationMode(bs, v.EIM_AReqAuthorizationMode); err != nil {
			return err
		}
		// END EIM_AReqAuthorizationMode
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		// END message
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	} else if v.PnC_AReqAuthorizationMode != nil {
		if err := bs.WriteBits(2, 1); err != nil {
			return err
		}
		// START PnC_AReqAuthorizationMode
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodePnC_AReqAuthorizationMode(bs, v.PnC_AReqAuthorizationMode); err != nil {
			return err
		}
		// END PnC_AReqAuthorizationMode
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

// encodeEIM_AReqAuthorizationMode encodes an EIM authorization mode.
// Currently a placeholder as the type has no fields defined.
func encodeEIM_AReqAuthorizationMode(bs *BitStream, v *generated.EIM_AReqAuthorizationMode) error {
	// Placeholder - no fields to encode
	return nil
}

// encodePnC_AReqAuthorizationMode encodes a PnC authorization mode.
func encodePnC_AReqAuthorizationMode(bs *BitStream, v *generated.PnC_AReqAuthorizationMode) error {
	// Optional GenChallenge (1 bit presence)
	if v.GenChallenge != nil && len(v.GenChallenge) > 0 {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
		// START GenChallenge
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
	} else {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}

	// Optional ContractCertificateChain (1 bit presence)
	if v.ContractCertificateChain != nil {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
		// START ContractCertificateChain
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodeCertificateChain(bs, v.ContractCertificateChain); err != nil {
			return err
		}
		// END ContractCertificateChain
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	} else {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}

	return nil
}

// encodeCertificateChain encodes a CertificateChain.
func encodeCertificateChain(bs *BitStream, v *generated.CertificateChain) error {
	// START Certificates (array)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write array length
	if err := bs.WriteUnsignedVar(uint64(len(v.Certificates))); err != nil {
		return err
	}
	// Encode each certificate
	for _, cert := range v.Certificates {
		// START certificate
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		// hexBinary encoding flag
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		// Write length
		if err := writeUint16(bs, uint16(len(cert))); err != nil {
			return err
		}
		// Write bytes
		if err := writeRawBytes(bs, cert); err != nil {
			return err
		}
		// END certificate
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	// END Certificates
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeAuthorizationReq decodes the AuthorizationReq body.
func DecodeAuthorizationReq(bs *BitStream) (*generated.AuthorizationReq, error) {
	// START Header
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// START SelectedAuthorizationService
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	selectedAuthService, err := readString(bs)
	if err != nil {
		return nil, err
	}
	// END SelectedAuthorizationService
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Union choice (2 bits)
	choice, err := bs.ReadBits(2)
	if err != nil {
		return nil, err
	}

	var eimMode *generated.EIM_AReqAuthorizationMode
	var pncMode *generated.PnC_AReqAuthorizationMode

	if choice == 0 {
		// EIM mode
		// START EIM_AReqAuthorizationMode
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		eim, err := decodeEIM_AReqAuthorizationMode(bs)
		if err != nil {
			return nil, err
		}
		eimMode = eim
		// END EIM_AReqAuthorizationMode
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		// END message
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	} else if choice == 1 {
		// PnC mode
		// START PnC_AReqAuthorizationMode
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		pnc, err := decodePnC_AReqAuthorizationMode(bs)
		if err != nil {
			return nil, err
		}
		pncMode = pnc
		// END PnC_AReqAuthorizationMode
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		// END message
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}
	// else choice == 2, no mode, message already ended

	return &generated.AuthorizationReq{
		Header:                       *header,
		SelectedAuthorizationService: selectedAuthService,
		EIM_AReqAuthorizationMode:    eimMode,
		PnC_AReqAuthorizationMode:    pncMode,
	}, nil
}

// decodeEIM_AReqAuthorizationMode decodes an EIM authorization mode.
func decodeEIM_AReqAuthorizationMode(bs *BitStream) (*generated.EIM_AReqAuthorizationMode, error) {
	// Placeholder - no fields to decode
	return &generated.EIM_AReqAuthorizationMode{}, nil
}

// decodePnC_AReqAuthorizationMode decodes a PnC authorization mode.
func decodePnC_AReqAuthorizationMode(bs *BitStream) (*generated.PnC_AReqAuthorizationMode, error) {
	// Optional GenChallenge presence
	p1, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	var genChallenge []byte
	if p1 == 1 {
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
		genChallenge = make([]byte, genChallengeLen)
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
	}

	// Optional ContractCertificateChain presence
	p2, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	var contractCertChain *generated.CertificateChain
	if p2 == 1 {
		// START ContractCertificateChain
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		chain, err := decodeCertificateChain(bs)
		if err != nil {
			return nil, err
		}
		contractCertChain = chain
		// END ContractCertificateChain
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}

	return &generated.PnC_AReqAuthorizationMode{
		GenChallenge:             genChallenge,
		ContractCertificateChain: contractCertChain,
	}, nil
}

// decodeCertificateChain decodes a CertificateChain.
func decodeCertificateChain(bs *BitStream) (*generated.CertificateChain, error) {
	// START Certificates
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Read array length
	count, err := bs.ReadUnsignedVar()
	if err != nil {
		return nil, err
	}
	certificates := make([][]byte, count)
	for i := uint64(0); i < count; i++ {
		// START certificate
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		// hexBinary encoding flag
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		// Read length
		certLen, err := readUint16(bs)
		if err != nil {
			return nil, err
		}
		// Read bytes
		cert := make([]byte, certLen)
		for j := 0; j < int(certLen); j++ {
			b, err := bs.ReadOctet()
			if err != nil {
				return nil, err
			}
			cert[j] = b
		}
		certificates[i] = cert
		// END certificate
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}
	// END Certificates
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.CertificateChain{
		Certificates: certificates,
	}, nil
}

// --------------------------- ServiceDetailRes ------------------------------

// EncodeTopLevelServiceDetailRes writes EXI header + event code (30) then
// delegates to EncodeServiceDetailRes.
func EncodeTopLevelServiceDetailRes(bs *BitStream, v *generated.ServiceDetailRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelServiceDetailRes: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for ServiceDetailRes is 30
	if err := bs.WriteBits(6, 30); err != nil {
		return err
	}
	return EncodeServiceDetailRes(bs, v)
}

// EncodeServiceDetailRes encodes the ServiceDetailRes body:
// Header, ResponseCode, ServiceID, and ServiceParameterList.
func EncodeServiceDetailRes(bs *BitStream, v *generated.ServiceDetailRes) error {
	if v == nil {
		return fmt.Errorf("EncodeServiceDetailRes: nil value")
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

	// START ServiceID (required uint16)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ServiceID encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := writeUint16(bs, v.ServiceID); err != nil {
		return err
	}
	// END ServiceID
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// START ServiceParameterList (required)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeServiceParameterList(bs, &v.ServiceParameterList); err != nil {
		return err
	}
	// END ServiceParameterList
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// END ServiceDetailRes
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// encodeServiceParameterList encodes a ServiceParameterList.
func encodeServiceParameterList(bs *BitStream, v *generated.ServiceParameterList) error {
	// START ParameterSets (array)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write array length
	if err := bs.WriteUnsignedVar(uint64(len(v.ParameterSets))); err != nil {
		return err
	}
	// Encode each parameter set
	for _, paramSet := range v.ParameterSets {
		// START parameter set
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodeParameterSet(bs, &paramSet); err != nil {
			return err
		}
		// END parameter set
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	// END ParameterSets
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// encodeParameterSet encodes a ParameterSet.
func encodeParameterSet(bs *BitStream, v *generated.ParameterSet) error {
	// START ParameterSetID (required uint16)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ParameterSetID encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := writeUint16(bs, v.ParameterSetID); err != nil {
		return err
	}
	// END ParameterSetID
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// START Parameters (array)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write array length
	if err := bs.WriteUnsignedVar(uint64(len(v.Parameters))); err != nil {
		return err
	}
	// Encode each parameter
	for _, param := range v.Parameters {
		// START parameter
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodeParameter(bs, &param); err != nil {
			return err
		}
		// END parameter
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	// END Parameters
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// encodeParameter encodes a Parameter.
func encodeParameter(bs *BitStream, v *generated.Parameter) error {
	// START Name (required string)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := writeString(bs, v.Name); err != nil {
		return err
	}
	// END Name
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Optional value fields (union choice)
	// 3-bit choice: 0=IntValue, 1=StrValue, 2=BoolValue, 3=END
	if v.IntValue != nil {
		if err := bs.WriteBits(3, 0); err != nil {
			return err
		}
		// START IntValue
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		// IntValue encoding flag
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		// Write as signed 16-bit
		if err := writeUint16(bs, uint16(*v.IntValue)); err != nil {
			return err
		}
		// END IntValue
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	} else if v.StrValue != nil {
		if err := bs.WriteBits(3, 1); err != nil {
			return err
		}
		// START StrValue
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := writeString(bs, *v.StrValue); err != nil {
			return err
		}
		// END StrValue
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	} else if v.BoolValue != nil {
		if err := bs.WriteBits(3, 2); err != nil {
			return err
		}
		// START BoolValue
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		// BoolValue encoding flag
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		// Boolean value (1 bit)
		if *v.BoolValue {
			if err := bs.WriteBits(1, 1); err != nil {
				return err
			}
		} else {
			if err := bs.WriteBits(1, 0); err != nil {
				return err
			}
		}
		// END BoolValue
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	} else {
		// No value
		if err := bs.WriteBits(3, 3); err != nil {
			return err
		}
	}

	return nil
}

// DecodeServiceDetailRes decodes the ServiceDetailRes body.
func DecodeServiceDetailRes(bs *BitStream) (*generated.ServiceDetailRes, error) {
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

	// START ServiceID
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ServiceID encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	serviceID, err := readUint16(bs)
	if err != nil {
		return nil, err
	}
	// END ServiceID
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// START ServiceParameterList
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	paramList, err := decodeServiceParameterList(bs)
	if err != nil {
		return nil, err
	}
	// END ServiceParameterList
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// END ServiceDetailRes
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.ServiceDetailRes{
		Header:               *header,
		ResponseCode:         responseCode,
		ServiceID:            serviceID,
		ServiceParameterList: *paramList,
	}, nil
}

// decodeServiceParameterList decodes a ServiceParameterList.
func decodeServiceParameterList(bs *BitStream) (*generated.ServiceParameterList, error) {
	// START ParameterSets
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Read array length
	count, err := bs.ReadUnsignedVar()
	if err != nil {
		return nil, err
	}
	parameterSets := make([]generated.ParameterSet, count)
	for i := uint64(0); i < count; i++ {
		// START parameter set
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		paramSet, err := decodeParameterSet(bs)
		if err != nil {
			return nil, err
		}
		parameterSets[i] = *paramSet
		// END parameter set
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}
	// END ParameterSets
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.ServiceParameterList{
		ParameterSets: parameterSets,
	}, nil
}

// decodeParameterSet decodes a ParameterSet.
func decodeParameterSet(bs *BitStream) (*generated.ParameterSet, error) {
	// START ParameterSetID
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// ParameterSetID encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	parameterSetID, err := readUint16(bs)
	if err != nil {
		return nil, err
	}
	// END ParameterSetID
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// START Parameters
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Read array length
	count, err := bs.ReadUnsignedVar()
	if err != nil {
		return nil, err
	}
	parameters := make([]generated.Parameter, count)
	for i := uint64(0); i < count; i++ {
		// START parameter
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		param, err := decodeParameter(bs)
		if err != nil {
			return nil, err
		}
		parameters[i] = *param
		// END parameter
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}
	// END Parameters
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.ParameterSet{
		ParameterSetID: parameterSetID,
		Parameters:     parameters,
	}, nil
}

// decodeParameter decodes a Parameter.
func decodeParameter(bs *BitStream) (*generated.Parameter, error) {
	// START Name
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	name, err := readString(bs)
	if err != nil {
		return nil, err
	}
	// END Name
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Optional value choice (3 bits)
	choice, err := bs.ReadBits(3)
	if err != nil {
		return nil, err
	}

	var intValue *int16
	var strValue *string
	var boolValue *bool

	if choice == 0 {
		// IntValue
		// START IntValue
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		// IntValue encoding flag
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		intVal, err := readUint16(bs)
		if err != nil {
			return nil, err
		}
		val := int16(intVal)
		intValue = &val
		// END IntValue
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	} else if choice == 1 {
		// StrValue
		// START StrValue
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		str, err := readString(bs)
		if err != nil {
			return nil, err
		}
		strValue = &str
		// END StrValue
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	} else if choice == 2 {
		// BoolValue
		// START BoolValue
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		// BoolValue encoding flag
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		boolBit, err := bs.ReadBits(1)
		if err != nil {
			return nil, err
		}
		val := boolBit == 1
		boolValue = &val
		// END BoolValue
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}
	// else choice == 3, no value

	return &generated.Parameter{
		Name:      name,
		IntValue:  intValue,
		StrValue:  strValue,
		BoolValue: boolValue,
	}, nil
}
