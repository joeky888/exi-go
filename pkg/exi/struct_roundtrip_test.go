package exi_test

import (
	"bytes"
	"testing"

	"example.com/exi-go/pkg/exi"
	"example.com/exi-go/pkg/v2g/generated"
)

// TestStructRoundTrip tests encode/decode round-trips for all implemented message types
// using the direct struct encoding API (not XML-based).

func TestStructRoundTripSessionSetupReq(t *testing.T) {
	orig := &generated.SessionSetupReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		EVCCID: []byte{0x0A, 0x1B, 0x2C, 0x3D, 0x4E, 0x5F},
	}

	// Encode
	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	// Decode
	decoded, err := exi.DecodeStruct(encoded, (*generated.SessionSetupReq)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.SessionSetupReq)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	// Compare
	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch: expected %x, got %x", orig.Header.SessionID, got.Header.SessionID)
	}
	if orig.Header.TimeStamp != got.Header.TimeStamp {
		t.Errorf("TimeStamp mismatch: expected %d, got %d", orig.Header.TimeStamp, got.Header.TimeStamp)
	}
	if !bytes.Equal(orig.EVCCID, got.EVCCID) {
		t.Errorf("EVCCID mismatch: expected %x, got %x", orig.EVCCID, got.EVCCID)
	}
}

func TestStructRoundTripSessionSetupRes(t *testing.T) {
	orig := &generated.SessionSetupRes{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ResponseCode: "OK",
		EVSEID:       []byte("ABCDEF123456"),
		DateTimeNow:  nil, // DateTimeNow is not currently encoded/decoded
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.SessionSetupRes)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.SessionSetupRes)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch: expected %x, got %x", orig.Header.SessionID, got.Header.SessionID)
	}
	if orig.Header.TimeStamp != got.Header.TimeStamp {
		t.Errorf("TimeStamp mismatch: expected %d, got %d", orig.Header.TimeStamp, got.Header.TimeStamp)
	}
	if orig.ResponseCode != got.ResponseCode {
		t.Errorf("ResponseCode mismatch: expected %s, got %s", orig.ResponseCode, got.ResponseCode)
	}
	if !bytes.Equal(orig.EVSEID, got.EVSEID) {
		t.Errorf("EVSEID mismatch: expected %x, got %x", orig.EVSEID, got.EVSEID)
	}
	// DateTimeNow is optional and not currently encoded in the implementation
	if got.DateTimeNow != nil {
		t.Logf("Note: DateTimeNow was decoded as %v (optional field)", *got.DateTimeNow)
	}
}

func TestStructRoundTripServiceDiscoveryReq(t *testing.T) {
	orig := &generated.ServiceDiscoveryReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ServiceScope:    nil,
		ServiceCategory: nil,
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.ServiceDiscoveryReq)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.ServiceDiscoveryReq)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch: expected %x, got %x", orig.Header.SessionID, got.Header.SessionID)
	}
	if orig.Header.TimeStamp != got.Header.TimeStamp {
		t.Errorf("TimeStamp mismatch: expected %d, got %d", orig.Header.TimeStamp, got.Header.TimeStamp)
	}
}

func TestStructRoundTripServiceDetailReq(t *testing.T) {
	serviceID := uint16(42)
	orig := &generated.ServiceDetailReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ServiceID: serviceID,
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.ServiceDetailReq)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.ServiceDetailReq)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch: expected %x, got %x", orig.Header.SessionID, got.Header.SessionID)
	}
	if orig.Header.TimeStamp != got.Header.TimeStamp {
		t.Errorf("TimeStamp mismatch: expected %d, got %d", orig.Header.TimeStamp, got.Header.TimeStamp)
	}
	if orig.ServiceID != got.ServiceID {
		t.Errorf("ServiceID mismatch: expected %d, got %d", orig.ServiceID, got.ServiceID)
	}
}

func TestStructRoundTripSessionStopReq(t *testing.T) {
	orig := &generated.SessionStopReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ChargingSession:          "Terminate",
		EVTerminationCode:        nil,
		EVTerminationExplanation: nil,
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.SessionStopReq)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.SessionStopReq)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch: expected %x, got %x", orig.Header.SessionID, got.Header.SessionID)
	}
	if orig.Header.TimeStamp != got.Header.TimeStamp {
		t.Errorf("TimeStamp mismatch: expected %d, got %d", orig.Header.TimeStamp, got.Header.TimeStamp)
	}
	if orig.ChargingSession != got.ChargingSession {
		t.Errorf("ChargingSession mismatch: expected %s, got %s", orig.ChargingSession, got.ChargingSession)
	}
}

func TestStructRoundTripSessionStopRes(t *testing.T) {
	orig := &generated.SessionStopRes{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ResponseCode: "OK",
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.SessionStopRes)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.SessionStopRes)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch: expected %x, got %x", orig.Header.SessionID, got.Header.SessionID)
	}
	if orig.Header.TimeStamp != got.Header.TimeStamp {
		t.Errorf("TimeStamp mismatch: expected %d, got %d", orig.Header.TimeStamp, got.Header.TimeStamp)
	}
	if orig.ResponseCode != got.ResponseCode {
		t.Errorf("ResponseCode mismatch: expected %s, got %s", orig.ResponseCode, got.ResponseCode)
	}
}

func TestStructRoundTripAuthorizationSetupReq(t *testing.T) {
	orig := &generated.AuthorizationSetupReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.AuthorizationSetupReq)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.AuthorizationSetupReq)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch: expected %x, got %x", orig.Header.SessionID, got.Header.SessionID)
	}
	if orig.Header.TimeStamp != got.Header.TimeStamp {
		t.Errorf("TimeStamp mismatch: expected %d, got %d", orig.Header.TimeStamp, got.Header.TimeStamp)
	}
}

func TestStructRoundTripMeteringConfirmationReq(t *testing.T) {
	orig := &generated.MeteringConfirmationReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.MeteringConfirmationReq)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.MeteringConfirmationReq)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch: expected %x, got %x", orig.Header.SessionID, got.Header.SessionID)
	}
	if orig.Header.TimeStamp != got.Header.TimeStamp {
		t.Errorf("TimeStamp mismatch: expected %d, got %d", orig.Header.TimeStamp, got.Header.TimeStamp)
	}
}

func TestStructRoundTripMeteringConfirmationRes(t *testing.T) {
	orig := &generated.MeteringConfirmationRes{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ResponseCode: "OK",
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.MeteringConfirmationRes)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.MeteringConfirmationRes)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch: expected %x, got %x", orig.Header.SessionID, got.Header.SessionID)
	}
	if orig.Header.TimeStamp != got.Header.TimeStamp {
		t.Errorf("TimeStamp mismatch: expected %d, got %d", orig.Header.TimeStamp, got.Header.TimeStamp)
	}
	if orig.ResponseCode != got.ResponseCode {
		t.Errorf("ResponseCode mismatch: expected %s, got %s", orig.ResponseCode, got.ResponseCode)
	}
}

func TestStructRoundTripAuthorizationRes(t *testing.T) {
	orig := &generated.AuthorizationRes{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ResponseCode:   "OK",
		EVSEProcessing: "Finished",
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.AuthorizationRes)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.AuthorizationRes)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch: expected %x, got %x", orig.Header.SessionID, got.Header.SessionID)
	}
	if orig.Header.TimeStamp != got.Header.TimeStamp {
		t.Errorf("TimeStamp mismatch: expected %d, got %d", orig.Header.TimeStamp, got.Header.TimeStamp)
	}
	if orig.ResponseCode != got.ResponseCode {
		t.Errorf("ResponseCode mismatch: expected %s, got %s", orig.ResponseCode, got.ResponseCode)
	}
	if orig.EVSEProcessing != got.EVSEProcessing {
		t.Errorf("EVSEProcessing mismatch: expected %s, got %s", orig.EVSEProcessing, got.EVSEProcessing)
	}
}

func TestStructRoundTripVehicleCheckInReq(t *testing.T) {
	parkingMethod := "AutomaticParking"
	orig := &generated.VehicleCheckInReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		EVCheckInStatus: "CheckIn",
		ParkingMethod:   &parkingMethod,
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.VehicleCheckInReq)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.VehicleCheckInReq)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch: expected %x, got %x", orig.Header.SessionID, got.Header.SessionID)
	}
	if orig.Header.TimeStamp != got.Header.TimeStamp {
		t.Errorf("TimeStamp mismatch: expected %d, got %d", orig.Header.TimeStamp, got.Header.TimeStamp)
	}
	if orig.EVCheckInStatus != got.EVCheckInStatus {
		t.Errorf("EVCheckInStatus mismatch: expected %s, got %s", orig.EVCheckInStatus, got.EVCheckInStatus)
	}
}

func TestStructRoundTripVehicleCheckInRes(t *testing.T) {
	vehicleCheckInResult := "CheckIn"
	orig := &generated.VehicleCheckInRes{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ResponseCode:         "OK",
		VehicleCheckInResult: &vehicleCheckInResult,
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.VehicleCheckInRes)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.VehicleCheckInRes)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch: expected %x, got %x", orig.Header.SessionID, got.Header.SessionID)
	}
	if orig.Header.TimeStamp != got.Header.TimeStamp {
		t.Errorf("TimeStamp mismatch: expected %d, got %d", orig.Header.TimeStamp, got.Header.TimeStamp)
	}
	if orig.ResponseCode != got.ResponseCode {
		t.Errorf("ResponseCode mismatch: expected %s, got %s", orig.ResponseCode, got.ResponseCode)
	}
	// Compare optional VehicleCheckInResult
	if orig.VehicleCheckInResult == nil && got.VehicleCheckInResult != nil {
		t.Errorf("VehicleCheckInResult mismatch: expected nil, got %v", *got.VehicleCheckInResult)
	}
	if orig.VehicleCheckInResult != nil && got.VehicleCheckInResult == nil {
		t.Errorf("VehicleCheckInResult mismatch: expected %v, got nil", *orig.VehicleCheckInResult)
	}
	if orig.VehicleCheckInResult != nil && got.VehicleCheckInResult != nil && *orig.VehicleCheckInResult != *got.VehicleCheckInResult {
		t.Errorf("VehicleCheckInResult mismatch: expected %v, got %v", *orig.VehicleCheckInResult, *got.VehicleCheckInResult)
	}
}

func TestStructRoundTripServiceDiscoveryRes(t *testing.T) {
	orig := &generated.ServiceDiscoveryRes{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ResponseCode:                  "OK",
		ServiceRenegotiationSupported: true,
		EnergyTransferServiceList: generated.ServiceList{
			Services: []generated.ServiceType{
				{ServiceID: 1, FreeService: true},
				{ServiceID: 2, FreeService: false},
			},
		},
		VASList: &generated.ServiceList{
			Services: []generated.ServiceType{
				{ServiceID: 100, FreeService: true},
			},
		},
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.ServiceDiscoveryRes)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.ServiceDiscoveryRes)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	// Compare Header
	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch: expected %x, got %x", orig.Header.SessionID, got.Header.SessionID)
	}
	if orig.Header.TimeStamp != got.Header.TimeStamp {
		t.Errorf("TimeStamp mismatch: expected %d, got %d", orig.Header.TimeStamp, got.Header.TimeStamp)
	}

	// Compare ResponseCode
	if orig.ResponseCode != got.ResponseCode {
		t.Errorf("ResponseCode mismatch: expected %s, got %s", orig.ResponseCode, got.ResponseCode)
	}

	// Compare ServiceRenegotiationSupported
	if orig.ServiceRenegotiationSupported != got.ServiceRenegotiationSupported {
		t.Errorf("ServiceRenegotiationSupported mismatch: expected %v, got %v",
			orig.ServiceRenegotiationSupported, got.ServiceRenegotiationSupported)
	}

	// Compare EnergyTransferServiceList
	if len(orig.EnergyTransferServiceList.Services) != len(got.EnergyTransferServiceList.Services) {
		t.Errorf("EnergyTransferServiceList length mismatch: expected %d, got %d",
			len(orig.EnergyTransferServiceList.Services), len(got.EnergyTransferServiceList.Services))
	} else {
		for i, origSvc := range orig.EnergyTransferServiceList.Services {
			gotSvc := got.EnergyTransferServiceList.Services[i]
			if origSvc.ServiceID != gotSvc.ServiceID {
				t.Errorf("EnergyTransferServiceList[%d].ServiceID mismatch: expected %d, got %d",
					i, origSvc.ServiceID, gotSvc.ServiceID)
			}
			if origSvc.FreeService != gotSvc.FreeService {
				t.Errorf("EnergyTransferServiceList[%d].FreeService mismatch: expected %v, got %v",
					i, origSvc.FreeService, gotSvc.FreeService)
			}
		}
	}

	// Compare VASList
	if orig.VASList == nil && got.VASList != nil {
		t.Errorf("VASList mismatch: expected nil, got %d services", len(got.VASList.Services))
	}
	if orig.VASList != nil && got.VASList == nil {
		t.Errorf("VASList mismatch: expected %d services, got nil", len(orig.VASList.Services))
	}
	if orig.VASList != nil && got.VASList != nil {
		if len(orig.VASList.Services) != len(got.VASList.Services) {
			t.Errorf("VASList length mismatch: expected %d, got %d",
				len(orig.VASList.Services), len(got.VASList.Services))
		} else {
			for i, origSvc := range orig.VASList.Services {
				gotSvc := got.VASList.Services[i]
				if origSvc.ServiceID != gotSvc.ServiceID {
					t.Errorf("VASList[%d].ServiceID mismatch: expected %d, got %d",
						i, origSvc.ServiceID, gotSvc.ServiceID)
				}
				if origSvc.FreeService != gotSvc.FreeService {
					t.Errorf("VASList[%d].FreeService mismatch: expected %v, got %v",
						i, origSvc.FreeService, gotSvc.FreeService)
				}
			}
		}
	}
}

func TestStructRoundTripServiceDetailRes(t *testing.T) {
	orig := &generated.ServiceDetailRes{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ResponseCode: "OK",
		ServiceID:    1,
		ServiceParameterList: generated.ServiceParameterList{
			ParameterSets: []generated.ParameterSet{},
		},
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.ServiceDetailRes)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.ServiceDetailRes)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch")
	}
	if orig.ResponseCode != got.ResponseCode {
		t.Errorf("ResponseCode mismatch")
	}
	if orig.ServiceID != got.ServiceID {
		t.Errorf("ServiceID mismatch")
	}
}

func TestStructRoundTripServiceSelectionReq(t *testing.T) {
	orig := &generated.ServiceSelectionReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		SelectedEnergyTransferService: generated.SelectedService{
			ServiceID:      1,
			ParameterSetID: nil,
		},
		SelectedVASList: nil,
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.ServiceSelectionReq)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.ServiceSelectionReq)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch")
	}
	if orig.SelectedEnergyTransferService.ServiceID != got.SelectedEnergyTransferService.ServiceID {
		t.Errorf("ServiceID mismatch")
	}
}

func TestStructRoundTripServiceSelectionRes(t *testing.T) {
	orig := &generated.ServiceSelectionRes{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ResponseCode: "OK",
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.ServiceSelectionRes)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.ServiceSelectionRes)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch")
	}
	if orig.ResponseCode != got.ResponseCode {
		t.Errorf("ResponseCode mismatch")
	}
}

func TestStructRoundTripAuthorizationSetupRes(t *testing.T) {
	orig := &generated.AuthorizationSetupRes{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ResponseCode: "OK",
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.AuthorizationSetupRes)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.AuthorizationSetupRes)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch")
	}
	if orig.ResponseCode != got.ResponseCode {
		t.Errorf("ResponseCode mismatch")
	}
}

func TestStructRoundTripAuthorizationReq(t *testing.T) {
	orig := &generated.AuthorizationReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.AuthorizationReq)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.AuthorizationReq)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch")
	}
}

func TestStructRoundTripPowerDeliveryReq(t *testing.T) {
	orig := &generated.PowerDeliveryReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ChargeProgress: "Start",
		EVPowerProfile: nil,
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.PowerDeliveryReq)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.PowerDeliveryReq)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch")
	}
	if orig.ChargeProgress != got.ChargeProgress {
		t.Errorf("ChargeProgress mismatch")
	}
}

func TestStructRoundTripPowerDeliveryRes(t *testing.T) {
	orig := &generated.PowerDeliveryRes{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ResponseCode: "OK",
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.PowerDeliveryRes)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.PowerDeliveryRes)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch")
	}
	if orig.ResponseCode != got.ResponseCode {
		t.Errorf("ResponseCode mismatch")
	}
}

func TestStructRoundTripScheduleExchangeReq(t *testing.T) {
	orig := &generated.ScheduleExchangeReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		MaximumSupportingPoints: 1024,
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.ScheduleExchangeReq)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.ScheduleExchangeReq)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch")
	}
	if orig.MaximumSupportingPoints != got.MaximumSupportingPoints {
		t.Errorf("MaximumSupportingPoints mismatch")
	}
}

func TestStructRoundTripScheduleExchangeRes(t *testing.T) {
	orig := &generated.ScheduleExchangeRes{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ResponseCode:   "OK",
		EVSEProcessing: "Finished",
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.ScheduleExchangeRes)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.ScheduleExchangeRes)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch")
	}
	if orig.ResponseCode != got.ResponseCode {
		t.Errorf("ResponseCode mismatch")
	}
	if orig.EVSEProcessing != got.EVSEProcessing {
		t.Errorf("EVSEProcessing mismatch")
	}
}

func TestStructRoundTripCertificateInstallationReq(t *testing.T) {
	orig := &generated.CertificateInstallationReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		OEMProvisioningCertChain: generated.CertificateChain{
			Certificates: [][]byte{
				[]byte("CERT123"),
			},
		},
		ListOfRootCertificateIDs: []string{
			"root-cert-1",
		},
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.CertificateInstallationReq)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.CertificateInstallationReq)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch")
	}
}

func TestStructRoundTripCertificateInstallationRes(t *testing.T) {
	orig := &generated.CertificateInstallationRes{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ResponseCode: "OK",
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.CertificateInstallationRes)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.CertificateInstallationRes)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch")
	}
	if orig.ResponseCode != got.ResponseCode {
		t.Errorf("ResponseCode mismatch")
	}
}

func TestStructRoundTripVehicleCheckOutReq(t *testing.T) {
	orig := &generated.VehicleCheckOutReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		EVCheckOutStatus: "CheckOut",
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.VehicleCheckOutReq)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.VehicleCheckOutReq)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch")
	}
	if orig.EVCheckOutStatus != got.EVCheckOutStatus {
		t.Errorf("EVCheckOutStatus mismatch")
	}
}

func TestStructRoundTripVehicleCheckOutRes(t *testing.T) {
	orig := &generated.VehicleCheckOutRes{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
		ResponseCode: "OK",
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.VehicleCheckOutRes)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.VehicleCheckOutRes)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch")
	}
	if orig.ResponseCode != got.ResponseCode {
		t.Errorf("ResponseCode mismatch")
	}
}

func TestStructRoundTripCLReqControlMode(t *testing.T) {
	orig := &generated.CLReqControlMode{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.CLReqControlMode)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.CLReqControlMode)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch")
	}
}

func TestStructRoundTripCLResControlMode(t *testing.T) {
	orig := &generated.CLResControlMode{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: uint64(1672531200),
		},
	}

	encoded, err := exi.EncodeStruct(orig)
	if err != nil {
		t.Fatalf("EncodeStruct failed: %v", err)
	}
	t.Logf("Encoded %d bytes", len(encoded))

	decoded, err := exi.DecodeStruct(encoded, (*generated.CLResControlMode)(nil))
	if err != nil {
		t.Fatalf("DecodeStruct failed: %v", err)
	}

	got, ok := decoded.(*generated.CLResControlMode)
	if !ok {
		t.Fatalf("decoded type assertion failed: got %T", decoded)
	}

	if !bytes.Equal(orig.Header.SessionID, got.Header.SessionID) {
		t.Errorf("SessionID mismatch")
	}
}
