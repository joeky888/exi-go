package exi_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"example.com/exi-go/pkg/exi"
	"example.com/exi-go/pkg/v2g/generated"
)

// helper to open golden file path relative to this test package
func goldenPath(name string) string {
	// current package directory: exi-go/pkg/exi
	// golden files are at iso15118-encoders/testvectors/
	// relative path: ../../../iso15118-encoders/testvectors/
	return filepath.Join("..", "..", "..", "iso15118-encoders", "testvectors", name)
}

func readGolden(t *testing.T, fname string) []byte {
	t.Helper()
	path := goldenPath(fname)
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skipf("golden file not found, skipping: %s", path)
		}
		t.Fatalf("failed to read golden file %s: %v", path, err)
	}
	return b
}

func diffReport(a, b []byte) string {
	if bytes.Equal(a, b) {
		return "identical"
	}
	min := len(a)
	if len(b) < min {
		min = len(b)
	}
	idx := -1
	for i := 0; i < min; i++ {
		if a[i] != b[i] {
			idx = i
			break
		}
	}
	if idx == -1 {
		// all common bytes equal; difference in length
		return fmt.Sprintf("same prefix (%d bytes), different lengths (got %d, want %d)", min, len(a), len(b))
	}
	// show a small hex context around the first difference
	const ctx = 24
	start := idx
	if start > ctx {
		start = idx - ctx
	}
	endA := idx + ctx
	if endA > len(a) {
		endA = len(a)
	}
	endB := idx + ctx
	if endB > len(b) {
		endB = len(b)
	}
	return fmt.Sprintf("first diff at byte %d\n got  (%d bytes): %s\n want (%d bytes): %s\n",
		idx,
		len(a), hex.EncodeToString(a[start:endA]),
		len(b), hex.EncodeToString(b[start:endB]),
	)
}

func TestCompareGeneratedToGolden(t *testing.T) {
	cases := []struct {
		name      string
		golden    string
		construct func() interface{}
		typename  string // for decoder path if needed
	}{
		{
			name:   "SessionSetupReq",
			golden: "SessionSetupReq.exi",
			construct: func() interface{} {
				// Match the exact test data used to generate the golden file.
				// SessionID: 4 bytes, TimeStamp: 1672531200, EVCCID: 6 bytes
				h := generated.MessageHeaderType{
					SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
					TimeStamp: uint64(1672531200),
				}
				return &generated.SessionSetupReq{
					Header: h,
					EVCCID: []byte{0x0A, 0x1B, 0x2C, 0x3D, 0x4E, 0x5F},
				}
			},
			typename: "SessionSetupReq",
		},
		{
			name:   "ServiceDiscoveryReq",
			golden: "ServiceDiscoveryReq.exi",
			construct: func() interface{} {
				// Match the exact test data used to generate the golden file.
				// The C golden uses only Header (no SupportedServiceIDs).
				// Our Go type has ServiceScope/ServiceCategory but those don't exist in ISO-20 spec.
				// Set them to nil to match the empty golden encoding.
				h := generated.MessageHeaderType{
					SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
					TimeStamp: uint64(1672531200),
				}
				return &generated.ServiceDiscoveryReq{
					Header:          h,
					ServiceScope:    nil,
					ServiceCategory: nil,
				}
			},
			typename: "ServiceDiscoveryReq",
		},
		{
			name:   "SessionSetupRes",
			golden: "SessionSetupRes.exi",
			construct: func() interface{} {
				h := generated.MessageHeaderType{
					SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
					TimeStamp: uint64(1672531200),
				}
				return &generated.SessionSetupRes{
					Header:       h,
					ResponseCode: "OK",
					EVSEID:       []byte("ABCDEF123456"),
					DateTimeNow:  nil, // DateTimeNow is not currently encoded
				}
			},
			typename: "SessionSetupRes",
		},
		{
			name:   "CertificateUpdateReq",
			golden: "CertificateUpdateReq.exi",
			construct: func() interface{} {
				h := generated.MessageHeaderType{
					SessionID: []byte{0x0A, 0x1B, 0x2C},
					TimeStamp: uint64(1672531200),
				}
				chain := generated.CertificateChain{
					Certificates: [][]byte{
						[]byte("MIIBIjAN..."),
						[]byte("MIICJz..."),
					},
				}
				return &generated.CertificateUpdateReq{
					Header:                     h,
					Id:                         "cert-req-1",
					ContractSignatureCertChain: chain,
					ContractID:                 "contract-123",
					ListOfRootCertificateIDs:   []string{"root-1", "root-2"},
					DHParams:                   []byte("dhparams"),
				}
			},
			typename: "CertificateUpdateReq",
		},
		// New test cases for all generated golden files
		{
			name:   "ServiceDiscoveryRes",
			golden: "ServiceDiscoveryRes.exi",
			construct: func() interface{} {
				return &generated.ServiceDiscoveryRes{
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
			},
			typename: "ServiceDiscoveryRes",
		},
		{
			name:   "ServiceDetailReq",
			golden: "ServiceDetailReq.exi",
			construct: func() interface{} {
				return &generated.ServiceDetailReq{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					ServiceID: 1,
				}
			},
			typename: "ServiceDetailReq",
		},
		{
			name:   "ServiceDetailRes",
			golden: "ServiceDetailRes.exi",
			construct: func() interface{} {
				return &generated.ServiceDetailRes{
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
			},
			typename: "ServiceDetailRes",
		},
		{
			name:   "ServiceSelectionReq",
			golden: "ServiceSelectionReq.exi",
			construct: func() interface{} {
				return &generated.ServiceSelectionReq{
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
			},
			typename: "ServiceSelectionReq",
		},
		{
			name:   "ServiceSelectionRes",
			golden: "ServiceSelectionRes.exi",
			construct: func() interface{} {
				return &generated.ServiceSelectionRes{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					ResponseCode: "OK",
				}
			},
			typename: "ServiceSelectionRes",
		},
		{
			name:   "SessionStopReq",
			golden: "SessionStopReq.exi",
			construct: func() interface{} {
				return &generated.SessionStopReq{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					ChargingSession:          "Terminate",
					EVTerminationCode:        nil,
					EVTerminationExplanation: nil,
				}
			},
			typename: "SessionStopReq",
		},
		{
			name:   "SessionStopRes",
			golden: "SessionStopRes.exi",
			construct: func() interface{} {
				return &generated.SessionStopRes{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					ResponseCode: "OK",
				}
			},
			typename: "SessionStopRes",
		},
		{
			name:   "AuthorizationSetupReq",
			golden: "AuthorizationSetupReq.exi",
			construct: func() interface{} {
				return &generated.AuthorizationSetupReq{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
				}
			},
			typename: "AuthorizationSetupReq",
		},
		{
			name:   "AuthorizationSetupRes",
			golden: "AuthorizationSetupRes.exi",
			construct: func() interface{} {
				return &generated.AuthorizationSetupRes{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					ResponseCode: "OK",
				}
			},
			typename: "AuthorizationSetupRes",
		},
		{
			name:   "AuthorizationReq",
			golden: "AuthorizationReq.exi",
			construct: func() interface{} {
				return &generated.AuthorizationReq{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
				}
			},
			typename: "AuthorizationReq",
		},
		{
			name:   "AuthorizationRes",
			golden: "AuthorizationRes.exi",
			construct: func() interface{} {
				return &generated.AuthorizationRes{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					ResponseCode:   "OK",
					EVSEProcessing: "Finished",
				}
			},
			typename: "AuthorizationRes",
		},
		{
			name:   "PowerDeliveryReq",
			golden: "PowerDeliveryReq.exi",
			construct: func() interface{} {
				return &generated.PowerDeliveryReq{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					ChargeProgress: "Start",
					EVPowerProfile: nil,
				}
			},
			typename: "PowerDeliveryReq",
		},
		{
			name:   "PowerDeliveryRes",
			golden: "PowerDeliveryRes.exi",
			construct: func() interface{} {
				return &generated.PowerDeliveryRes{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					ResponseCode: "OK",
				}
			},
			typename: "PowerDeliveryRes",
		},
		{
			name:   "ScheduleExchangeReq",
			golden: "ScheduleExchangeReq.exi",
			construct: func() interface{} {
				return &generated.ScheduleExchangeReq{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					MaximumSupportingPoints: 1024,
				}
			},
			typename: "ScheduleExchangeReq",
		},
		{
			name:   "ScheduleExchangeRes",
			golden: "ScheduleExchangeRes.exi",
			construct: func() interface{} {
				return &generated.ScheduleExchangeRes{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					ResponseCode:   "OK",
					EVSEProcessing: "Finished",
				}
			},
			typename: "ScheduleExchangeRes",
		},
		{
			name:   "MeteringConfirmationReq",
			golden: "MeteringConfirmationReq.exi",
			construct: func() interface{} {
				return &generated.MeteringConfirmationReq{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
				}
			},
			typename: "MeteringConfirmationReq",
		},
		{
			name:   "MeteringConfirmationRes",
			golden: "MeteringConfirmationRes.exi",
			construct: func() interface{} {
				return &generated.MeteringConfirmationRes{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					ResponseCode: "OK",
				}
			},
			typename: "MeteringConfirmationRes",
		},
		{
			name:   "CertificateInstallationReq",
			golden: "CertificateInstallationReq.exi",
			construct: func() interface{} {
				return &generated.CertificateInstallationReq{
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
			},
			typename: "CertificateInstallationReq",
		},
		{
			name:   "CertificateInstallationRes",
			golden: "CertificateInstallationRes.exi",
			construct: func() interface{} {
				return &generated.CertificateInstallationRes{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					ResponseCode: "OK",
				}
			},
			typename: "CertificateInstallationRes",
		},
		{
			name:   "VehicleCheckInReq",
			golden: "VehicleCheckInReq.exi",
			construct: func() interface{} {
				parkingMethod := "AutomaticParking"
				return &generated.VehicleCheckInReq{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					EVCheckInStatus: "CheckIn",
					ParkingMethod:   &parkingMethod,
				}
			},
			typename: "VehicleCheckInReq",
		},
		{
			name:   "VehicleCheckInRes",
			golden: "VehicleCheckInRes.exi",
			construct: func() interface{} {
				vehicleCheckInResult := "CheckIn"
				return &generated.VehicleCheckInRes{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					ResponseCode:         "OK",
					VehicleCheckInResult: &vehicleCheckInResult,
				}
			},
			typename: "VehicleCheckInRes",
		},
		{
			name:   "VehicleCheckOutReq",
			golden: "VehicleCheckOutReq.exi",
			construct: func() interface{} {
				return &generated.VehicleCheckOutReq{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					EVCheckOutStatus: "CheckOut",
				}
			},
			typename: "VehicleCheckOutReq",
		},
		{
			name:   "VehicleCheckOutRes",
			golden: "VehicleCheckOutRes.exi",
			construct: func() interface{} {
				return &generated.VehicleCheckOutRes{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
					ResponseCode: "OK",
				}
			},
			typename: "VehicleCheckOutRes",
		},
		{
			name:   "CLReqControlMode",
			golden: "CLReqControlMode.exi",
			construct: func() interface{} {
				return &generated.CLReqControlMode{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
				}
			},
			typename: "CLReqControlMode",
		},
		{
			name:   "CLResControlMode",
			golden: "CLResControlMode.exi",
			construct: func() interface{} {
				return &generated.CLResControlMode{
					Header: generated.MessageHeaderType{
						SessionID: []byte{0x0A, 0x1B, 0x2C, 0x3D},
						TimeStamp: uint64(1672531200),
					},
				}
			},
			typename: "CLResControlMode",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// read golden
			golden := readGolden(t, tc.golden)

			// initialize codec
			c := exi.NewCodec(nil)
			if err := c.Init(); err != nil {
				t.Fatalf("codec init failed: %v", err)
			}
			defer c.Shutdown()

			// encode using schema-informed encoder
			enc, err := c.EncodeStruct(tc.construct())
			if err != nil {
				t.Fatalf("EncodeStruct failed: %v", err)
			}

			// Compare bytes
			if bytes.Equal(enc, golden) {
				// success
				return
			}

			// If mismatched, write debug files to TMP for inspection
			tmpGot := filepath.Join(os.TempDir(), fmt.Sprintf("got_%s.exi", tc.name))
			tmpWant := filepath.Join(os.TempDir(), fmt.Sprintf("want_%s.exi", tc.name))
			_ = os.WriteFile(tmpGot, enc, 0o644)
			_ = os.WriteFile(tmpWant, golden, 0o644)

			t.Fatalf("encoded output differs from golden:\n%s\nLocal copies: got=%s want=%s",
				diffReport(enc, golden), tmpGot, tmpWant)
		})
	}
}
