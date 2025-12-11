package exi

import (
	"fmt"

	"example.com/exi-go/pkg/v2g/generated"
)

// Phase 4 certificate-related message type implementations.
// These message types involve certificate chains and complex nested structures.
//
// CertificateInstallationReq (event 7): Header + OEMProvisioningCertChain + ListOfRootCertificateIDs
// CertificateInstallationRes (event 8): Header + ResponseCode + EVSEProcessing + multiple certificate chains

// --------------------- CertificateInstallationReq --------------------------

// EncodeTopLevelCertificateInstallationReq writes EXI header + event code (7) then
// delegates to EncodeCertificateInstallationReq.
func EncodeTopLevelCertificateInstallationReq(bs *BitStream, v *generated.CertificateInstallationReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelCertificateInstallationReq: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for CertificateInstallationReq is 7
	if err := bs.WriteBits(6, 7); err != nil {
		return err
	}
	return EncodeCertificateInstallationReq(bs, v)
}

// EncodeCertificateInstallationReq encodes the CertificateInstallationReq body:
// Header, OEMProvisioningCertChain, and ListOfRootCertificateIDs.
func EncodeCertificateInstallationReq(bs *BitStream, v *generated.CertificateInstallationReq) error {
	if v == nil {
		return fmt.Errorf("EncodeCertificateInstallationReq: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// START OEMProvisioningCertChain (required)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeCertificateChain(bs, &v.OEMProvisioningCertChain); err != nil {
		return err
	}
	// END OEMProvisioningCertChain
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// START ListOfRootCertificateIDs (required array of strings)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write array length
	if err := bs.WriteUnsignedVar(uint64(len(v.ListOfRootCertificateIDs))); err != nil {
		return err
	}
	// Encode each certificate ID string
	for _, certID := range v.ListOfRootCertificateIDs {
		// START certID
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := writeString(bs, certID); err != nil {
			return err
		}
		// END certID
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}
	// END ListOfRootCertificateIDs
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// END CertificateInstallationReq
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeCertificateInstallationReq decodes the CertificateInstallationReq body.
func DecodeCertificateInstallationReq(bs *BitStream) (*generated.CertificateInstallationReq, error) {
	// START Header
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// START OEMProvisioningCertChain
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	oemProvisioningCertChain, err := decodeCertificateChain(bs)
	if err != nil {
		return nil, err
	}
	// END OEMProvisioningCertChain
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// START ListOfRootCertificateIDs
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Read array length
	count, err := bs.ReadUnsignedVar()
	if err != nil {
		return nil, err
	}
	rootCertIDs := make([]string, count)
	for i := uint64(0); i < count; i++ {
		// START certID
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		certID, err := readString(bs)
		if err != nil {
			return nil, err
		}
		rootCertIDs[i] = certID
		// END certID
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}
	// END ListOfRootCertificateIDs
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// END CertificateInstallationReq
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.CertificateInstallationReq{
		Header:                   *header,
		OEMProvisioningCertChain: *oemProvisioningCertChain,
		ListOfRootCertificateIDs: rootCertIDs,
	}, nil
}

// --------------------- CertificateInstallationRes --------------------------

// EncodeTopLevelCertificateInstallationRes writes EXI header + event code (8) then
// delegates to EncodeCertificateInstallationRes.
func EncodeTopLevelCertificateInstallationRes(bs *BitStream, v *generated.CertificateInstallationRes) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelCertificateInstallationRes: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for CertificateInstallationRes is 8
	if err := bs.WriteBits(6, 8); err != nil {
		return err
	}
	return EncodeCertificateInstallationRes(bs, v)
}

// EncodeCertificateInstallationRes encodes the CertificateInstallationRes body:
// Header, ResponseCode, EVSEProcessing, CPSCertificateChain, ContractSignatureEncryptedPrivateKey,
// DHPublicKey, and ContractCertificateChain.
func EncodeCertificateInstallationRes(bs *BitStream, v *generated.CertificateInstallationRes) error {
	if v == nil {
		return fmt.Errorf("EncodeCertificateInstallationRes: nil value")
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

	// START CPSCertificateChain (required)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeCertificateChain(bs, &v.CPSCertificateChain); err != nil {
		return err
	}
	// END CPSCertificateChain
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// START ContractSignatureEncryptedPrivateKey (required string)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := writeString(bs, v.ContractSignatureEncryptedPrivateKey); err != nil {
		return err
	}
	// END ContractSignatureEncryptedPrivateKey
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// START DHPublicKey (required binary data)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// hexBinary encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write length
	if err := writeUint16(bs, uint16(len(v.DHPublicKey))); err != nil {
		return err
	}
	// Write bytes
	if err := writeRawBytes(bs, v.DHPublicKey); err != nil {
		return err
	}
	// END DHPublicKey
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// START ContractCertificateChain (required)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeCertificateChain(bs, &v.ContractCertificateChain); err != nil {
		return err
	}
	// END ContractCertificateChain
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// END CertificateInstallationRes
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeCertificateInstallationRes decodes the CertificateInstallationRes body.
func DecodeCertificateInstallationRes(bs *BitStream) (*generated.CertificateInstallationRes, error) {
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

	// START CPSCertificateChain
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	cpsCertChain, err := decodeCertificateChain(bs)
	if err != nil {
		return nil, err
	}
	// END CPSCertificateChain
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// START ContractSignatureEncryptedPrivateKey
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	contractSigEncPrivKey, err := readString(bs)
	if err != nil {
		return nil, err
	}
	// END ContractSignatureEncryptedPrivateKey
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// START DHPublicKey
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// hexBinary encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// Read length
	dhPubKeyLen, err := readUint16(bs)
	if err != nil {
		return nil, err
	}
	// Read bytes
	dhPublicKey := make([]byte, dhPubKeyLen)
	for i := 0; i < int(dhPubKeyLen); i++ {
		b, err := bs.ReadOctet()
		if err != nil {
			return nil, err
		}
		dhPublicKey[i] = b
	}
	// END DHPublicKey
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// START ContractCertificateChain
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	contractCertChain, err := decodeCertificateChain(bs)
	if err != nil {
		return nil, err
	}
	// END ContractCertificateChain
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// END CertificateInstallationRes
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.CertificateInstallationRes{
		Header:                               *header,
		ResponseCode:                         responseCode,
		EVSEProcessing:                       evseProcessing,
		CPSCertificateChain:                  *cpsCertChain,
		ContractSignatureEncryptedPrivateKey: contractSigEncPrivKey,
		DHPublicKey:                          dhPublicKey,
		ContractCertificateChain:             *contractCertChain,
	}, nil
}
