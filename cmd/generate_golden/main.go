package main

import (
	"fmt"
	"os"
	"path/filepath"

	"example.com/exi-go/pkg/exi"
	"example.com/exi-go/pkg/v2g/generated"
)

func main() {
	// Create output directory
	outDir := filepath.Join("..", "..", "..", "iso15118-encoders", "testvectors")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	codec := exi.NewCodec(nil)
	if err := codec.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize codec: %v\n", err)
		os.Exit(1)
	}
	defer codec.Shutdown()

	// Fixed timestamp matching the Makefile XML templates
	timestamp := uint64(1672531200)
	sessionID := []byte{0x0A, 0x1B, 0x2C, 0x3D}

	// Track success/failure
	successCount := 0
	failCount := 0

	// Helper to encode and write
	encode := func(name string, msg interface{}) {
		data, err := codec.EncodeStruct(msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to encode %s: %v\n", name, err)
			failCount++
			return
		}

		outPath := filepath.Join(outDir, name+".exi")
		if err := os.WriteFile(outPath, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s.exi: %v\n", name, err)
			failCount++
			return
		}
		fmt.Printf("Generated %s (%d bytes)\n", outPath, len(data))
		successCount++
	}

	// 1. SessionSetupReq (event 35)
	encode("SessionSetupReq", &generated.SessionSetupReq{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		EVCCID: []byte{0x0A, 0x1B, 0x2C, 0x3D, 0x4E, 0x5F},
	})

	// 2. SessionSetupRes (event 36)
	encode("SessionSetupRes", &generated.SessionSetupRes{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ResponseCode: "OK",
		EVSEID:       []byte("ABCDEF123456"),
		DateTimeNow:  nil,
	})

	// 3. ServiceDiscoveryReq (event 31)
	encode("ServiceDiscoveryReq", &generated.ServiceDiscoveryReq{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ServiceScope:    nil,
		ServiceCategory: nil,
	})

	// 4. ServiceDiscoveryRes (event 32)
	encode("ServiceDiscoveryRes", &generated.ServiceDiscoveryRes{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
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
	})

	// 5. ServiceDetailReq (event 29)
	encode("ServiceDetailReq", &generated.ServiceDetailReq{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ServiceID: 1,
	})

	// 6. ServiceDetailRes (event 30)
	encode("ServiceDetailRes", &generated.ServiceDetailRes{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ResponseCode: "OK",
		ServiceID:    1,
		ServiceParameterList: generated.ServiceParameterList{
			ParameterSets: []generated.ParameterSet{},
		},
	})

	// 7. ServiceSelectionReq (event 33)
	encode("ServiceSelectionReq", &generated.ServiceSelectionReq{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		SelectedEnergyTransferService: generated.SelectedService{
			ServiceID:      1,
			ParameterSetID: nil,
		},
		SelectedVASList: nil,
	})

	// 8. ServiceSelectionRes (event 34)
	encode("ServiceSelectionRes", &generated.ServiceSelectionRes{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ResponseCode: "OK",
	})

	// 9. SessionStopReq (event 37)
	encode("SessionStopReq", &generated.SessionStopReq{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ChargingSession:          "Terminate",
		EVTerminationCode:        nil,
		EVTerminationExplanation: nil,
	})

	// 10. SessionStopRes (event 38)
	encode("SessionStopRes", &generated.SessionStopRes{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ResponseCode: "OK",
	})

	// 11. AuthorizationSetupReq (event 2)
	encode("AuthorizationSetupReq", &generated.AuthorizationSetupReq{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
	})

	// 12. AuthorizationSetupRes (event 3)
	encode("AuthorizationSetupRes", &generated.AuthorizationSetupRes{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ResponseCode: "OK",
	})

	// 13. AuthorizationReq (event 0)
	encode("AuthorizationReq", &generated.AuthorizationReq{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
	})

	// 14. AuthorizationRes (event 1)
	encode("AuthorizationRes", &generated.AuthorizationRes{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ResponseCode:   "OK",
		EVSEProcessing: "Finished",
	})

	// 15. PowerDeliveryReq (event 21)
	encode("PowerDeliveryReq", &generated.PowerDeliveryReq{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ChargeProgress: "Start",
		EVPowerProfile: nil,
	})

	// 16. PowerDeliveryRes (event 22)
	encode("PowerDeliveryRes", &generated.PowerDeliveryRes{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ResponseCode: "OK",
	})

	// 17. ScheduleExchangeReq (event 27)
	encode("ScheduleExchangeReq", &generated.ScheduleExchangeReq{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		MaximumSupportingPoints: 1024,
	})

	// 18. ScheduleExchangeRes (event 28)
	encode("ScheduleExchangeRes", &generated.ScheduleExchangeRes{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ResponseCode:   "OK",
		EVSEProcessing: "Finished",
	})

	// 19. MeteringConfirmationReq (event 16)
	encode("MeteringConfirmationReq", &generated.MeteringConfirmationReq{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
	})

	// 20. MeteringConfirmationRes (event 17)
	encode("MeteringConfirmationRes", &generated.MeteringConfirmationRes{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ResponseCode: "OK",
	})

	// 21. CertificateInstallationReq (event 7)
	encode("CertificateInstallationReq", &generated.CertificateInstallationReq{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		OEMProvisioningCertChain: generated.CertificateChain{
			Certificates: [][]byte{
				[]byte("CERT123"),
			},
		},
		ListOfRootCertificateIDs: []string{
			"root-cert-1",
		},
	})

	// 22. CertificateInstallationRes (event 8)
	encode("CertificateInstallationRes", &generated.CertificateInstallationRes{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ResponseCode: "OK",
	})

	// 23. VehicleCheckInReq (event 49)
	parkingMethod := "AutomaticParking"
	encode("VehicleCheckInReq", &generated.VehicleCheckInReq{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		EVCheckInStatus: "CheckIn",
		ParkingMethod:   &parkingMethod,
	})

	// 24. VehicleCheckInRes (event 50)
	vehicleCheckInResult := "CheckIn"
	encode("VehicleCheckInRes", &generated.VehicleCheckInRes{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ResponseCode:         "OK",
		VehicleCheckInResult: &vehicleCheckInResult,
	})

	// 25. VehicleCheckOutReq (event 51)
	encode("VehicleCheckOutReq", &generated.VehicleCheckOutReq{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		EVCheckOutStatus: "CheckOut",
	})

	// 26. VehicleCheckOutRes (event 52)
	encode("VehicleCheckOutRes", &generated.VehicleCheckOutRes{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
		ResponseCode: "OK",
	})

	// 27. CLReqControlMode (event 4)
	encode("CLReqControlMode", &generated.CLReqControlMode{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
	})

	// 28. CLResControlMode (event 5)
	encode("CLResControlMode", &generated.CLResControlMode{
		Header: generated.MessageHeaderType{
			SessionID: sessionID,
			TimeStamp: timestamp,
		},
	})

	// Summary
	fmt.Println()
	fmt.Printf("Golden file generation complete!\n")
	fmt.Printf("Success: %d/%d\n", successCount, successCount+failCount)
	if failCount > 0 {
		fmt.Printf("Failed: %d\n", failCount)
		os.Exit(1)
	}
}
