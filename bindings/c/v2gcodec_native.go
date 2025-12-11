/*
cgo bridge for exi-go - Native struct encoding/decoding

This file implements C-exported functions for encoding/decoding individual
ISO 15118-20 message types directly using the native Go structs, bypassing XML.
This provides a more efficient API for applications that work directly with
message structures.

Build note:
  - Build together with v2gcodec.go as a C-shared library:
    go build -buildmode=c-shared -o libv2gcodec.so

License: Apache-2.0 (match repository)
*/
package main

/*
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"example.com/exi-go/pkg/exi"
	"example.com/exi-go/pkg/v2g/generated"
)

// Message type enum matching ISO 15118-20 CommonMessages event codes
const (
	// Event codes for ISO 15118-20 CommonMessages
	V2G_MSG_AuthorizationReq           = 0
	V2G_MSG_AuthorizationRes           = 1
	V2G_MSG_AuthorizationSetupReq      = 2
	V2G_MSG_AuthorizationSetupRes      = 3
	V2G_MSG_CLReqControlMode           = 4
	V2G_MSG_CLResControlMode           = 5
	V2G_MSG_CertificateInstallationReq = 7
	V2G_MSG_CertificateInstallationRes = 8
	V2G_MSG_MeteringConfirmationReq    = 16
	V2G_MSG_MeteringConfirmationRes    = 17
	V2G_MSG_PowerDeliveryReq           = 21
	V2G_MSG_PowerDeliveryRes           = 22
	V2G_MSG_ScheduleExchangeReq        = 27
	V2G_MSG_ScheduleExchangeRes        = 28
	V2G_MSG_ServiceDetailReq           = 29
	V2G_MSG_ServiceDetailRes           = 30
	V2G_MSG_ServiceDiscoveryReq        = 31
	V2G_MSG_ServiceDiscoveryRes        = 32
	V2G_MSG_ServiceSelectionReq        = 33
	V2G_MSG_ServiceSelectionRes        = 34
	V2G_MSG_SessionSetupReq            = 35
	V2G_MSG_SessionSetupRes            = 36
	V2G_MSG_SessionStopReq             = 37
	V2G_MSG_SessionStopRes             = 38
	V2G_MSG_VehicleCheckInReq          = 49
	V2G_MSG_VehicleCheckInRes          = 50
	V2G_MSG_VehicleCheckOutReq         = 51
	V2G_MSG_VehicleCheckOutRes         = 52
)

//export v2g_encode_struct
func v2g_encode_struct(msg_type C.int, json_data *C.char, json_len C.size_t, out_exi **C.uint8_t, out_len *C.size_t) C.int {
	if json_data == nil || json_len == 0 || out_exi == nil || out_len == nil {
		setLastError("v2g_encode_struct: invalid arguments")
		return C.int(_v2g_err_invalid)
	}

	// Convert JSON to Go bytes
	jsonBytes := C.GoBytes(unsafe.Pointer(json_data), C.int(json_len))

	// Decode JSON into appropriate struct type based on msg_type
	var result []byte
	var err error

	switch int(msg_type) {
	case V2G_MSG_SessionSetupReq:
		var msg generated.SessionSetupReq
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal SessionSetupReq: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_SessionSetupRes:
		var msg generated.SessionSetupRes
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal SessionSetupRes: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_ServiceDiscoveryReq:
		var msg generated.ServiceDiscoveryReq
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal ServiceDiscoveryReq: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_ServiceDiscoveryRes:
		var msg generated.ServiceDiscoveryRes
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal ServiceDiscoveryRes: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_ServiceDetailReq:
		var msg generated.ServiceDetailReq
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal ServiceDetailReq: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_ServiceDetailRes:
		var msg generated.ServiceDetailRes
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal ServiceDetailRes: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_AuthorizationReq:
		var msg generated.AuthorizationReq
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal AuthorizationReq: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_AuthorizationRes:
		var msg generated.AuthorizationRes
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal AuthorizationRes: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_AuthorizationSetupReq:
		var msg generated.AuthorizationSetupReq
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal AuthorizationSetupReq: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_AuthorizationSetupRes:
		var msg generated.AuthorizationSetupRes
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal AuthorizationSetupRes: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_ServiceSelectionReq:
		var msg generated.ServiceSelectionReq
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal ServiceSelectionReq: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_ServiceSelectionRes:
		var msg generated.ServiceSelectionRes
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal ServiceSelectionRes: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_PowerDeliveryReq:
		var msg generated.PowerDeliveryReq
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal PowerDeliveryReq: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_PowerDeliveryRes:
		var msg generated.PowerDeliveryRes
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal PowerDeliveryRes: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_SessionStopReq:
		var msg generated.SessionStopReq
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal SessionStopReq: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_SessionStopRes:
		var msg generated.SessionStopRes
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal SessionStopRes: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_ScheduleExchangeReq:
		var msg generated.ScheduleExchangeReq
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal ScheduleExchangeReq: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_ScheduleExchangeRes:
		var msg generated.ScheduleExchangeRes
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal ScheduleExchangeRes: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_MeteringConfirmationReq:
		var msg generated.MeteringConfirmationReq
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal MeteringConfirmationReq: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_MeteringConfirmationRes:
		var msg generated.MeteringConfirmationRes
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal MeteringConfirmationRes: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_CertificateInstallationReq:
		var msg generated.CertificateInstallationReq
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal CertificateInstallationReq: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_CertificateInstallationRes:
		var msg generated.CertificateInstallationRes
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal CertificateInstallationRes: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_VehicleCheckInReq:
		var msg generated.VehicleCheckInReq
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal VehicleCheckInReq: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_VehicleCheckInRes:
		var msg generated.VehicleCheckInRes
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal VehicleCheckInRes: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_VehicleCheckOutReq:
		var msg generated.VehicleCheckOutReq
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal VehicleCheckOutReq: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_VehicleCheckOutRes:
		var msg generated.VehicleCheckOutRes
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal VehicleCheckOutRes: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_CLReqControlMode:
		var msg generated.CLReqControlMode
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal CLReqControlMode: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	case V2G_MSG_CLResControlMode:
		var msg generated.CLResControlMode
		if err := json.Unmarshal(jsonBytes, &msg); err != nil {
			setLastError("unmarshal CLResControlMode: %v", err)
			return C.int(_v2g_err_invalid)
		}
		result, err = exi.EncodeStruct(&msg)

	default:
		setLastError("v2g_encode_struct: unsupported message type %d", msg_type)
		return C.int(_v2g_err_invalid)
	}

	if err != nil {
		setLastError("encode failed: %v", err)
		return C.int(_v2g_err_encode)
	}

	// Allocate C memory and copy result
	cbuf := C.CBytes(result)
	if cbuf == nil {
		setLastError("encode: out of memory")
		return C.int(_v2g_err_oom)
	}

	*out_exi = (*C.uint8_t)(cbuf)
	*out_len = C.size_t(len(result))
	return C.int(_v2g_ok)
}

//export v2g_decode_struct
func v2g_decode_struct(msg_type C.int, exi_data *C.uint8_t, exi_len C.size_t, out_json **C.char, out_len *C.size_t) C.int {
	if exi_data == nil || exi_len == 0 || out_json == nil || out_len == nil {
		setLastError("v2g_decode_struct: invalid arguments")
		return C.int(_v2g_err_invalid)
	}

	// Convert EXI to Go bytes
	exiBytes := C.GoBytes(unsafe.Pointer(exi_data), C.int(exi_len))

	// Decode EXI into appropriate struct type
	var jsonBytes []byte
	var err error

	switch int(msg_type) {
	case V2G_MSG_SessionSetupReq:
		msg, err := exi.DecodeStruct(exiBytes, &generated.SessionSetupReq{})
		if err != nil {
			setLastError("decode SessionSetupReq: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_SessionSetupRes:
		msg, err := exi.DecodeStruct(exiBytes, &generated.SessionSetupRes{})
		if err != nil {
			setLastError("decode SessionSetupRes: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_ServiceDiscoveryReq:
		msg, err := exi.DecodeStruct(exiBytes, &generated.ServiceDiscoveryReq{})
		if err != nil {
			setLastError("decode ServiceDiscoveryReq: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_ServiceDiscoveryRes:
		msg, err := exi.DecodeStruct(exiBytes, &generated.ServiceDiscoveryRes{})
		if err != nil {
			setLastError("decode ServiceDiscoveryRes: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_ServiceDetailReq:
		msg, err := exi.DecodeStruct(exiBytes, &generated.ServiceDetailReq{})
		if err != nil {
			setLastError("decode ServiceDetailReq: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_ServiceDetailRes:
		msg, err := exi.DecodeStruct(exiBytes, &generated.ServiceDetailRes{})
		if err != nil {
			setLastError("decode ServiceDetailRes: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_AuthorizationReq:
		msg, err := exi.DecodeStruct(exiBytes, &generated.AuthorizationReq{})
		if err != nil {
			setLastError("decode AuthorizationReq: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_AuthorizationRes:
		msg, err := exi.DecodeStruct(exiBytes, &generated.AuthorizationRes{})
		if err != nil {
			setLastError("decode AuthorizationRes: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_AuthorizationSetupReq:
		msg, err := exi.DecodeStruct(exiBytes, &generated.AuthorizationSetupReq{})
		if err != nil {
			setLastError("decode AuthorizationSetupReq: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_AuthorizationSetupRes:
		msg, err := exi.DecodeStruct(exiBytes, &generated.AuthorizationSetupRes{})
		if err != nil {
			setLastError("decode AuthorizationSetupRes: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_ServiceSelectionReq:
		msg, err := exi.DecodeStruct(exiBytes, &generated.ServiceSelectionReq{})
		if err != nil {
			setLastError("decode ServiceSelectionReq: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_ServiceSelectionRes:
		msg, err := exi.DecodeStruct(exiBytes, &generated.ServiceSelectionRes{})
		if err != nil {
			setLastError("decode ServiceSelectionRes: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_PowerDeliveryReq:
		msg, err := exi.DecodeStruct(exiBytes, &generated.PowerDeliveryReq{})
		if err != nil {
			setLastError("decode PowerDeliveryReq: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_PowerDeliveryRes:
		msg, err := exi.DecodeStruct(exiBytes, &generated.PowerDeliveryRes{})
		if err != nil {
			setLastError("decode PowerDeliveryRes: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_SessionStopReq:
		msg, err := exi.DecodeStruct(exiBytes, &generated.SessionStopReq{})
		if err != nil {
			setLastError("decode SessionStopReq: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_SessionStopRes:
		msg, err := exi.DecodeStruct(exiBytes, &generated.SessionStopRes{})
		if err != nil {
			setLastError("decode SessionStopRes: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_ScheduleExchangeReq:
		msg, err := exi.DecodeStruct(exiBytes, &generated.ScheduleExchangeReq{})
		if err != nil {
			setLastError("decode ScheduleExchangeReq: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_ScheduleExchangeRes:
		msg, err := exi.DecodeStruct(exiBytes, &generated.ScheduleExchangeRes{})
		if err != nil {
			setLastError("decode ScheduleExchangeRes: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_MeteringConfirmationReq:
		msg, err := exi.DecodeStruct(exiBytes, &generated.MeteringConfirmationReq{})
		if err != nil {
			setLastError("decode MeteringConfirmationReq: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_MeteringConfirmationRes:
		msg, err := exi.DecodeStruct(exiBytes, &generated.MeteringConfirmationRes{})
		if err != nil {
			setLastError("decode MeteringConfirmationRes: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_CertificateInstallationReq:
		msg, err := exi.DecodeStruct(exiBytes, &generated.CertificateInstallationReq{})
		if err != nil {
			setLastError("decode CertificateInstallationReq: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_CertificateInstallationRes:
		msg, err := exi.DecodeStruct(exiBytes, &generated.CertificateInstallationRes{})
		if err != nil {
			setLastError("decode CertificateInstallationRes: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_VehicleCheckInReq:
		msg, err := exi.DecodeStruct(exiBytes, &generated.VehicleCheckInReq{})
		if err != nil {
			setLastError("decode VehicleCheckInReq: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_VehicleCheckInRes:
		msg, err := exi.DecodeStruct(exiBytes, &generated.VehicleCheckInRes{})
		if err != nil {
			setLastError("decode VehicleCheckInRes: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_VehicleCheckOutReq:
		msg, err := exi.DecodeStruct(exiBytes, &generated.VehicleCheckOutReq{})
		if err != nil {
			setLastError("decode VehicleCheckOutReq: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_VehicleCheckOutRes:
		msg, err := exi.DecodeStruct(exiBytes, &generated.VehicleCheckOutRes{})
		if err != nil {
			setLastError("decode VehicleCheckOutRes: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_CLReqControlMode:
		msg, err := exi.DecodeStruct(exiBytes, &generated.CLReqControlMode{})
		if err != nil {
			setLastError("decode CLReqControlMode: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	case V2G_MSG_CLResControlMode:
		msg, err := exi.DecodeStruct(exiBytes, &generated.CLResControlMode{})
		if err != nil {
			setLastError("decode CLResControlMode: %v", err)
			return C.int(_v2g_err_decode)
		}
		jsonBytes, err = json.Marshal(msg)

	default:
		setLastError("v2g_decode_struct: unsupported message type %d", msg_type)
		return C.int(_v2g_err_invalid)
	}

	if err != nil {
		setLastError("json marshal failed: %v", err)
		return C.int(_v2g_err_internal)
	}

	// Allocate C memory for NUL-terminated JSON string
	clen := C.size_t(len(jsonBytes) + 1)
	cptr := C.malloc(clen)
	if cptr == nil {
		setLastError("decode: out of memory")
		return C.int(_v2g_err_oom)
	}

	// Copy JSON bytes and add NUL terminator
	C.memcpy(cptr, unsafe.Pointer(&jsonBytes[0]), C.size_t(len(jsonBytes)))
	lastBytePtr := unsafe.Pointer(uintptr(cptr) + uintptr(len(jsonBytes)))
	*(*byte)(lastBytePtr) = 0

	*out_json = (*C.char)(cptr)
	*out_len = C.size_t(len(jsonBytes))
	return C.int(_v2g_ok)
}

//export v2g_message_type_name
func v2g_message_type_name(msg_type C.int) *C.char {
	var name string
	switch int(msg_type) {
	case V2G_MSG_AuthorizationReq:
		name = "AuthorizationReq"
	case V2G_MSG_AuthorizationRes:
		name = "AuthorizationRes"
	case V2G_MSG_AuthorizationSetupReq:
		name = "AuthorizationSetupReq"
	case V2G_MSG_AuthorizationSetupRes:
		name = "AuthorizationSetupRes"
	case V2G_MSG_CLReqControlMode:
		name = "CLReqControlMode"
	case V2G_MSG_CLResControlMode:
		name = "CLResControlMode"
	case V2G_MSG_CertificateInstallationReq:
		name = "CertificateInstallationReq"
	case V2G_MSG_CertificateInstallationRes:
		name = "CertificateInstallationRes"
	case V2G_MSG_MeteringConfirmationReq:
		name = "MeteringConfirmationReq"
	case V2G_MSG_MeteringConfirmationRes:
		name = "MeteringConfirmationRes"
	case V2G_MSG_PowerDeliveryReq:
		name = "PowerDeliveryReq"
	case V2G_MSG_PowerDeliveryRes:
		name = "PowerDeliveryRes"
	case V2G_MSG_ScheduleExchangeReq:
		name = "ScheduleExchangeReq"
	case V2G_MSG_ScheduleExchangeRes:
		name = "ScheduleExchangeRes"
	case V2G_MSG_ServiceDetailReq:
		name = "ServiceDetailReq"
	case V2G_MSG_ServiceDetailRes:
		name = "ServiceDetailRes"
	case V2G_MSG_ServiceDiscoveryReq:
		name = "ServiceDiscoveryReq"
	case V2G_MSG_ServiceDiscoveryRes:
		name = "ServiceDiscoveryRes"
	case V2G_MSG_ServiceSelectionReq:
		name = "ServiceSelectionReq"
	case V2G_MSG_ServiceSelectionRes:
		name = "ServiceSelectionRes"
	case V2G_MSG_SessionSetupReq:
		name = "SessionSetupReq"
	case V2G_MSG_SessionSetupRes:
		name = "SessionSetupRes"
	case V2G_MSG_SessionStopReq:
		name = "SessionStopReq"
	case V2G_MSG_SessionStopRes:
		name = "SessionStopRes"
	case V2G_MSG_VehicleCheckInReq:
		name = "VehicleCheckInReq"
	case V2G_MSG_VehicleCheckInRes:
		name = "VehicleCheckInRes"
	case V2G_MSG_VehicleCheckOutReq:
		name = "VehicleCheckOutReq"
	case V2G_MSG_VehicleCheckOutRes:
		name = "VehicleCheckOutRes"
	default:
		name = fmt.Sprintf("Unknown(%d)", msg_type)
	}
	return C.CString(name)
}
