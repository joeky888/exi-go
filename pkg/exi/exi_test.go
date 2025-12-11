package exi_test

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"

	"example.com/exi-go/pkg/exi"
	"example.com/exi-go/pkg/v2g/generated"
)

// helper that encodes the provided XML bytes with the exi codec and decodes them back.
func roundTripXMLWithCodec(t *testing.T, codec *exi.Codec, xmlIn []byte) []byte {
	t.Helper()
	// encode
	exiBytes, err := codec.EncodeXML(xmlIn)
	if err != nil {
		t.Fatalf("EncodeXML failed: %v", err)
	}
	// decode
	outXML, err := codec.DecodeEXI(exiBytes)
	if err != nil {
		t.Fatalf("DecodeEXI failed: %v", err)
	}
	return outXML
}

func TestSessionSetupReqRoundTrip(t *testing.T) {
	orig := &generated.SessionSetupReq{
		EVCCID: []byte{0x0A, 0x1B, 0x2C, 0x3D, 0x4E, 0x5F},
	}

	xmlIn, err := xml.Marshal(orig)
	if err != nil {
		t.Fatalf("xml.Marshal failed: %v", err)
	}

	c := exi.NewCodec(nil)
	if err := c.Init(); err != nil {
		t.Fatalf("codec Init failed: %v", err)
	}
	defer c.Shutdown()

	outXML := roundTripXMLWithCodec(t, c, xmlIn)

	var got generated.SessionSetupReq
	if err := xml.Unmarshal(outXML, &got); err != nil {
		t.Fatalf("xml.Unmarshal failed: %v\noutput: %s", err, string(outXML))
	}

	if !reflect.DeepEqual(orig.EVCCID, got.EVCCID) {
		t.Fatalf("SessionSetupReq mismatch:\nexpected: %#v\nactual:   %#v", orig, got)
	}
}

func TestServiceDiscoveryReqRoundTrip(t *testing.T) {
	scope := "public"
	category := "charging"
	orig := &generated.ServiceDiscoveryReq{
		ServiceScope:    &scope,
		ServiceCategory: &category,
	}

	xmlIn, err := xml.Marshal(orig)
	if err != nil {
		t.Fatalf("xml.Marshal failed: %v", err)
	}

	c := exi.NewCodec(nil)
	if err := c.Init(); err != nil {
		t.Fatalf("codec Init failed: %v", err)
	}
	defer c.Shutdown()

	outXML := roundTripXMLWithCodec(t, c, xmlIn)

	var got generated.ServiceDiscoveryReq
	if err := xml.Unmarshal(outXML, &got); err != nil {
		t.Fatalf("xml.Unmarshal failed: %v\noutput: %s", err, string(outXML))
	}

	// Compare pointers' values safely
	if orig.ServiceScope == nil && got.ServiceScope != nil {
		t.Fatalf("ServiceScope mismatch: expected nil got %v", *got.ServiceScope)
	}
	if orig.ServiceScope != nil && got.ServiceScope == nil {
		t.Fatalf("ServiceScope mismatch: expected %v got nil", *orig.ServiceScope)
	}
	if orig.ServiceScope != nil && got.ServiceScope != nil && *orig.ServiceScope != *got.ServiceScope {
		t.Fatalf("ServiceScope mismatch: expected %v got %v", *orig.ServiceScope, *got.ServiceScope)
	}

	if orig.ServiceCategory == nil && got.ServiceCategory != nil {
		t.Fatalf("ServiceCategory mismatch: expected nil got %v", *got.ServiceCategory)
	}
	if orig.ServiceCategory != nil && got.ServiceCategory == nil {
		t.Fatalf("ServiceCategory mismatch: expected %v got nil", *orig.ServiceCategory)
	}
	if orig.ServiceCategory != nil && got.ServiceCategory != nil && *orig.ServiceCategory != *got.ServiceCategory {
		t.Fatalf("ServiceCategory mismatch: expected %v got %v", *orig.ServiceCategory, *got.ServiceCategory)
	}
}

func TestCertificateUpdateReqRoundTrip(t *testing.T) {
	orig := &generated.CertificateUpdateReq{
		Id: "cert-req-1",
		ContractSignatureCertChain: generated.CertificateChain{
			Certificates: [][]byte{
				[]byte("MIIBIjAN..."),
				[]byte("MIICJz..."),
			},
		},
		ContractID:               "contract-123",
		ListOfRootCertificateIDs: []string{"root-1", "root-2"},
		DHParams:                 []byte("c2FtcGxlLWRoLXBhcmFtcz0="),
	}

	xmlIn, err := xml.Marshal(orig)
	if err != nil {
		t.Fatalf("xml.Marshal failed: %v", err)
	}

	c := exi.NewCodec(nil)
	if err := c.Init(); err != nil {
		t.Fatalf("codec Init failed: %v", err)
	}
	defer c.Shutdown()

	outXML := roundTripXMLWithCodec(t, c, xmlIn)

	var got generated.CertificateUpdateReq
	if err := xml.Unmarshal(outXML, &got); err != nil {
		t.Fatalf("xml.Unmarshal failed: %v\noutput: %s", err, string(outXML))
	}

	// compare primitive fields
	if orig.Id != got.Id {
		t.Fatalf("Id mismatch: expected %s got %s", orig.Id, got.Id)
	}
	if orig.ContractID != got.ContractID {
		t.Fatalf("ContractID mismatch: expected %s got %s", orig.ContractID, got.ContractID)
	}
	if !bytes.Equal(orig.DHParams, got.DHParams) {
		t.Fatalf("DHParams mismatch: expected %v got %v", orig.DHParams, got.DHParams)
	}
	// compare slices
	if !reflect.DeepEqual(orig.ListOfRootCertificateIDs, got.ListOfRootCertificateIDs) {
		t.Fatalf("ListOfRootCertificateIDs mismatch: expected %#v got %#v", orig.ListOfRootCertificateIDs, got.ListOfRootCertificateIDs)
	}
	if !reflect.DeepEqual(orig.ContractSignatureCertChain, got.ContractSignatureCertChain) {
		t.Fatalf("ContractSignatureCertChain mismatch: expected %#v got %#v", orig.ContractSignatureCertChain, got.ContractSignatureCertChain)
	}
}
