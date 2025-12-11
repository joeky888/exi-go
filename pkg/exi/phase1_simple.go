package exi

import (
	"fmt"

	"example.com/exi-go/pkg/v2g/generated"
)

// Phase 1 simple message type implementations: VehicleCheckInReq and VehicleCheckOutReq.
// These are straightforward message types that follow established patterns in this package.
//
// VehicleCheckInReq (event 49):  Header + EVCheckInStatus + optional ParkingMethod
// VehicleCheckOutReq (event 51): Header + EVCheckOutStatus + CheckOutTime

// --------------------------- VehicleCheckInReq -----------------------------

// EncodeTopLevelVehicleCheckInReq writes EXI header + event code (42) then
// delegates to EncodeVehicleCheckInReq.
func EncodeTopLevelVehicleCheckInReq(bs *BitStream, v *generated.VehicleCheckInReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelVehicleCheckInReq: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for VehicleCheckInReq is 49
	if err := bs.WriteBits(6, 49); err != nil {
		return err
	}
	return EncodeVehicleCheckInReq(bs, v)
}

// EncodeVehicleCheckInReq encodes the VehicleCheckInReq body:
// Header, EVCheckInStatus, and optional ParkingMethod.
func EncodeVehicleCheckInReq(bs *BitStream, v *generated.VehicleCheckInReq) error {
	if v == nil {
		return fmt.Errorf("EncodeVehicleCheckInReq: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// START EVCheckInStatus (required string)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := writeString(bs, v.EVCheckInStatus); err != nil {
		return err
	}
	// END EVCheckInStatus
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Optional ParkingMethod (1 bit presence flag)
	if v.ParkingMethod != nil {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
		// START ParkingMethod
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := writeString(bs, *v.ParkingMethod); err != nil {
			return err
		}
		// END ParkingMethod
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	} else {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}

	// END VehicleCheckInReq
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeVehicleCheckInReq decodes the VehicleCheckInReq body.
func DecodeVehicleCheckInReq(bs *BitStream) (*generated.VehicleCheckInReq, error) {
	// START Header
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// START EVCheckInStatus
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	evCheckInStatus, err := readString(bs)
	if err != nil {
		return nil, err
	}
	// END EVCheckInStatus
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Optional ParkingMethod presence
	p, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	var parkingMethod *string
	if p == 1 {
		// START ParkingMethod
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		pm, err := readString(bs)
		if err != nil {
			return nil, err
		}
		parkingMethod = &pm
		// END ParkingMethod
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}

	// END VehicleCheckInReq
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.VehicleCheckInReq{
		Header:          *header,
		EVCheckInStatus: evCheckInStatus,
		ParkingMethod:   parkingMethod,
	}, nil
}

// -------------------------- VehicleCheckOutReq -----------------------------

// EncodeTopLevelVehicleCheckOutReq writes EXI header + event code (44) then
// delegates to EncodeVehicleCheckOutReq.
func EncodeTopLevelVehicleCheckOutReq(bs *BitStream, v *generated.VehicleCheckOutReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelVehicleCheckOutReq: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for VehicleCheckOutReq is 51
	if err := bs.WriteBits(6, 51); err != nil {
		return err
	}
	return EncodeVehicleCheckOutReq(bs, v)
}

// EncodeVehicleCheckOutReq encodes the VehicleCheckOutReq body:
// Header, EVCheckOutStatus, and CheckOutTime.
func EncodeVehicleCheckOutReq(bs *BitStream, v *generated.VehicleCheckOutReq) error {
	if v == nil {
		return fmt.Errorf("EncodeVehicleCheckOutReq: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// START EVCheckOutStatus (required string)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := writeString(bs, v.EVCheckOutStatus); err != nil {
		return err
	}
	// END EVCheckOutStatus
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// START CheckOutTime (required uint64)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// CheckOutTime encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write timestamp as unsigned-var
	if err := bs.WriteUnsignedVar(v.CheckOutTime); err != nil {
		return err
	}
	// END CheckOutTime
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// END VehicleCheckOutReq
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// DecodeVehicleCheckOutReq decodes the VehicleCheckOutReq body.
func DecodeVehicleCheckOutReq(bs *BitStream) (*generated.VehicleCheckOutReq, error) {
	// START Header
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// START EVCheckOutStatus
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	evCheckOutStatus, err := readString(bs)
	if err != nil {
		return nil, err
	}
	// END EVCheckOutStatus
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// START CheckOutTime
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	// CheckOutTime encoding flag
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	checkOutTime, err := bs.ReadUnsignedVar()
	if err != nil {
		return nil, err
	}
	// END CheckOutTime
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// END VehicleCheckOutReq
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.VehicleCheckOutReq{
		Header:           *header,
		EVCheckOutStatus: evCheckOutStatus,
		CheckOutTime:     checkOutTime,
	}, nil
}

// ------------------------- ServiceSelectionReq -----------------------------

// EncodeTopLevelServiceSelectionReq writes EXI header + event code (33) then
// delegates to EncodeServiceSelectionReq.
func EncodeTopLevelServiceSelectionReq(bs *BitStream, v *generated.ServiceSelectionReq) error {
	if v == nil {
		return fmt.Errorf("EncodeTopLevelServiceSelectionReq: nil value")
	}
	// EXI simple header
	if err := bs.WriteBits(8, 0x80); err != nil {
		return err
	}
	// Event code for ServiceSelectionReq is 33
	if err := bs.WriteBits(6, 33); err != nil {
		return err
	}
	return EncodeServiceSelectionReq(bs, v)
}

// EncodeServiceSelectionReq encodes the ServiceSelectionReq body:
// Header, SelectedEnergyTransferService, and optional SelectedVASList.
func EncodeServiceSelectionReq(bs *BitStream, v *generated.ServiceSelectionReq) error {
	if v == nil {
		return fmt.Errorf("EncodeServiceSelectionReq: nil value")
	}

	// START Header
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeMessageHeaderType(bs, &v.Header); err != nil {
		return err
	}

	// START SelectedEnergyTransferService (required)
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	if err := encodeSelectedService(bs, &v.SelectedEnergyTransferService); err != nil {
		return err
	}
	// END SelectedEnergyTransferService
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Optional SelectedVASList (1 bit presence flag)
	if v.SelectedVASList != nil {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
		// START SelectedVASList
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodeSelectedServiceList(bs, v.SelectedVASList); err != nil {
			return err
		}
		// END SelectedVASList
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	} else {
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}

	// END ServiceSelectionReq
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	return nil
}

// encodeSelectedService encodes a SelectedService (ServiceID + optional ParameterSetID).
func encodeSelectedService(bs *BitStream, v *generated.SelectedService) error {
	// START ServiceID
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// ServiceID encoding flag
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}
	// Write ServiceID as uint16
	if err := writeUint16(bs, v.ServiceID); err != nil {
		return err
	}
	// END ServiceID
	if err := bs.WriteBits(1, 0); err != nil {
		return err
	}

	// Optional ParameterSetID (1 bit presence flag)
	if v.ParameterSetID != nil {
		if err := bs.WriteBits(1, 1); err != nil {
			return err
		}
		// START ParameterSetID
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		// ParameterSetID encoding flag
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := writeUint16(bs, *v.ParameterSetID); err != nil {
			return err
		}
		// END ParameterSetID
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

// encodeSelectedServiceList encodes a SelectedServiceList (array of SelectedService).
func encodeSelectedServiceList(bs *BitStream, v *generated.SelectedServiceList) error {
	// Write array length
	if err := bs.WriteUnsignedVar(uint64(len(v.SelectedServices))); err != nil {
		return err
	}

	// Encode each SelectedService
	for _, service := range v.SelectedServices {
		// START SelectedService
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
		if err := encodeSelectedService(bs, &service); err != nil {
			return err
		}
		// END SelectedService
		if err := bs.WriteBits(1, 0); err != nil {
			return err
		}
	}

	return nil
}

// DecodeServiceSelectionReq decodes the ServiceSelectionReq body.
func DecodeServiceSelectionReq(bs *BitStream) (*generated.ServiceSelectionReq, error) {
	// START Header
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	header, err := decodeMessageHeaderType(bs)
	if err != nil {
		return nil, err
	}

	// START SelectedEnergyTransferService
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}
	selectedEnergyTransferService, err := decodeSelectedService(bs)
	if err != nil {
		return nil, err
	}
	// END SelectedEnergyTransferService
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	// Optional SelectedVASList presence
	p, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	var selectedVASList *generated.SelectedServiceList
	if p == 1 {
		// START SelectedVASList
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		vasList, err := decodeSelectedServiceList(bs)
		if err != nil {
			return nil, err
		}
		selectedVASList = vasList
		// END SelectedVASList
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}

	// END ServiceSelectionReq
	if _, err := bs.ReadBits(1); err != nil {
		return nil, err
	}

	return &generated.ServiceSelectionReq{
		Header:                        *header,
		SelectedEnergyTransferService: *selectedEnergyTransferService,
		SelectedVASList:               selectedVASList,
	}, nil
}

// decodeSelectedService decodes a SelectedService.
func decodeSelectedService(bs *BitStream) (*generated.SelectedService, error) {
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

	// Optional ParameterSetID presence
	p, err := bs.ReadBits(1)
	if err != nil {
		return nil, err
	}
	var parameterSetID *uint16
	if p == 1 {
		// START ParameterSetID
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		// ParameterSetID encoding flag
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		psID, err := readUint16(bs)
		if err != nil {
			return nil, err
		}
		parameterSetID = &psID
		// END ParameterSetID
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}

	return &generated.SelectedService{
		ServiceID:      serviceID,
		ParameterSetID: parameterSetID,
	}, nil
}

// decodeSelectedServiceList decodes a SelectedServiceList.
func decodeSelectedServiceList(bs *BitStream) (*generated.SelectedServiceList, error) {
	// Read array length
	count, err := bs.ReadUnsignedVar()
	if err != nil {
		return nil, err
	}

	services := make([]generated.SelectedService, count)
	for i := uint64(0); i < count; i++ {
		// START SelectedService
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
		service, err := decodeSelectedService(bs)
		if err != nil {
			return nil, err
		}
		services[i] = *service
		// END SelectedService
		if _, err := bs.ReadBits(1); err != nil {
			return nil, err
		}
	}

	return &generated.SelectedServiceList{
		SelectedServices: services,
	}, nil
}
