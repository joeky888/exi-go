package exi

import (
	"testing"

	"example.com/exi-go/pkg/v2g/generated"
)

// BenchmarkSessionSetupReq benchmarks SessionSetupReq encoding/decoding
func BenchmarkSessionSetupReq(b *testing.B) {
	msg := &generated.SessionSetupReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x01, 0x02, 0x03, 0x04},
			TimeStamp: 1234567890,
		},
		EVCCID: []byte{0x0A, 0x1B, 0x2C, 0x3D, 0x4E, 0x5F},
	}

	b.Run("Encode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := EncodeStruct(msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	encoded, _ := EncodeStruct(msg)
	b.Run("Decode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := DecodeStruct(encoded, (*generated.SessionSetupReq)(nil))
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("RoundTrip", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			enc, err := EncodeStruct(msg)
			if err != nil {
				b.Fatal(err)
			}
			_, err = DecodeStruct(enc, (*generated.SessionSetupReq)(nil))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkSessionSetupRes benchmarks SessionSetupRes encoding/decoding
func BenchmarkSessionSetupRes(b *testing.B) {
	msg := &generated.SessionSetupRes{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x01, 0x02, 0x03, 0x04},
			TimeStamp: 1234567890,
		},
		ResponseCode: "OK",
		EVSEID:       []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11},
	}

	b.Run("Encode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := EncodeStruct(msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	encoded, _ := EncodeStruct(msg)
	b.Run("Decode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := DecodeStruct(encoded, (*generated.SessionSetupRes)(nil))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkServiceDiscoveryRes benchmarks ServiceDiscoveryRes encoding/decoding
func BenchmarkServiceDiscoveryRes(b *testing.B) {
	msg := &generated.ServiceDiscoveryRes{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
			TimeStamp: 1672531200,
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

	b.Run("Encode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := EncodeStruct(msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	encoded, _ := EncodeStruct(msg)
	b.Run("Decode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := DecodeStruct(encoded, (*generated.ServiceDiscoveryRes)(nil))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkPowerDeliveryReq benchmarks complex message with nested structures
func BenchmarkPowerDeliveryReq(b *testing.B) {
	msg := &generated.PowerDeliveryReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x01, 0x02, 0x03, 0x04},
			TimeStamp: 1234567890,
		},
		EVProcessing:         "Ongoing",
		ChargeProgress:       "Start",
		EVPowerProfile:       &generated.EVPowerProfile{},
		BPT_ChannelSelection: nil,
	}

	b.Run("Encode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := EncodeStruct(msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	encoded, _ := EncodeStruct(msg)
	b.Run("Decode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := DecodeStruct(encoded, (*generated.PowerDeliveryReq)(nil))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkVehicleCheckInReq benchmarks large message type
func BenchmarkVehicleCheckInReq(b *testing.B) {
	msg := &generated.VehicleCheckInReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x01, 0x02, 0x03, 0x04},
			TimeStamp: 1234567890,
		},
		EVCheckInStatus: "CheckIn",
		ParkingMethod:   stringPtr("AutoParking"),
	}

	b.Run("Encode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := EncodeStruct(msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	encoded, _ := EncodeStruct(msg)
	b.Run("Decode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := DecodeStruct(encoded, (*generated.VehicleCheckInReq)(nil))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkCertificateInstallationReq benchmarks certificate message with chains
func BenchmarkCertificateInstallationReq(b *testing.B) {
	msg := &generated.CertificateInstallationReq{
		Header: generated.MessageHeaderType{
			SessionID: []byte{0x01, 0x02, 0x03, 0x04},
			TimeStamp: 1234567890,
		},
		OEMProvisioningCertChain: generated.CertificateChain{
			Certificates: [][]byte{{0x30, 0x82, 0x01, 0x00}},
		},
		ListOfRootCertificateIDs: []string{"cert1", "cert2"},
	}

	b.Run("Encode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := EncodeStruct(msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	encoded, _ := EncodeStruct(msg)
	b.Run("Decode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := DecodeStruct(encoded, (*generated.CertificateInstallationReq)(nil))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkAllMessageTypes runs benchmarks for all 26 message types
func BenchmarkAllMessageTypes(b *testing.B) {
	messages := []struct {
		name      string
		msg       interface{}
		prototype interface{}
	}{
		{"SessionSetupReq", &generated.SessionSetupReq{
			Header: generated.MessageHeaderType{SessionID: []byte{1, 2, 3, 4}, TimeStamp: 1234567890},
			EVCCID: []byte{10, 27, 44, 61, 78, 95},
		}, (*generated.SessionSetupReq)(nil)},
		{"SessionSetupRes", &generated.SessionSetupRes{
			Header:       generated.MessageHeaderType{SessionID: []byte{1, 2, 3, 4}, TimeStamp: 1234567890},
			ResponseCode: "OK",
			EVSEID:       []byte{0xAA, 0xBB, 0xCC},
		}, (*generated.SessionSetupRes)(nil)},
		{"ServiceDiscoveryReq", &generated.ServiceDiscoveryReq{
			Header: generated.MessageHeaderType{SessionID: []byte{1, 2, 3, 4}, TimeStamp: 1234567890},
		}, (*generated.ServiceDiscoveryReq)(nil)},
		{"ServiceDetailReq", &generated.ServiceDetailReq{
			Header:    generated.MessageHeaderType{SessionID: []byte{1, 2, 3, 4}, TimeStamp: 1234567890},
			ServiceID: 1,
		}, (*generated.ServiceDetailReq)(nil)},
		{"SessionStopReq", &generated.SessionStopReq{
			Header:          generated.MessageHeaderType{SessionID: []byte{1, 2, 3, 4}, TimeStamp: 1234567890},
			ChargingSession: "Terminate",
		}, (*generated.SessionStopReq)(nil)},
		{"AuthorizationSetupReq", &generated.AuthorizationSetupReq{
			Header: generated.MessageHeaderType{SessionID: []byte{1, 2, 3, 4}, TimeStamp: 1234567890},
		}, (*generated.AuthorizationSetupReq)(nil)},
		{"MeteringConfirmationReq", &generated.MeteringConfirmationReq{
			Header: generated.MessageHeaderType{SessionID: []byte{1, 2, 3, 4}, TimeStamp: 1234567890},
		}, (*generated.MeteringConfirmationReq)(nil)},
		{"CLReqControlMode", &generated.CLReqControlMode{
			Header: generated.MessageHeaderType{SessionID: []byte{1, 2, 3, 4}, TimeStamp: 1234567890},
		}, (*generated.CLReqControlMode)(nil)},
		{"CLResControlMode", &generated.CLResControlMode{
			Header: generated.MessageHeaderType{SessionID: []byte{1, 2, 3, 4}, TimeStamp: 1234567890},
		}, (*generated.CLResControlMode)(nil)},
	}

	for _, msg := range messages {
		b.Run(msg.name+"/Encode", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := EncodeStruct(msg.msg)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		encoded, _ := EncodeStruct(msg.msg)
		b.Run(msg.name+"/Decode", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := DecodeStruct(encoded, msg.prototype)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
