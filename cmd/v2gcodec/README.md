# v2gcodec - ISO 15118-20 EXI Codec CLI

Command-line tool for encoding and decoding ISO 15118-20 messages using EXI (Efficient XML Interchange) format.

## Installation

```bash
cd go
go build -o v2gcodec cmd/v2gcodec/main.go
```

## Usage

### Decode EXI to JSON

Decode a hex-encoded EXI message directly from command line:

```bash
./v2gcodec decode "808c02050d961e8809ac39d06204050d961ea72f80"
```

Output:
```json
{
  "Header": {
    "SessionID": "ChssPQ==",
    "TimeStamp": 1672531200
  },
  "EVCCID": "ChssPU5f"
}
```

Decode from a file:

```bash
./v2gcodec decode -in message.exi
```

Decode from stdin:

```bash
cat message.hex | ./v2gcodec decode
```

### Encode JSON to EXI

Encode a JSON message to hex-encoded EXI:

```bash
./v2gcodec encode -type SessionSetupReq '{"Header":{"SessionID":"ChssPQ==","TimeStamp":1672531200},"EVCCID":"ChssPU5f"}'
```

Output:
```
808c02050d961e8809ac39d06204050d961ea72f80
```

Encode from file:

```bash
./v2gcodec encode -type SessionSetupReq -in message.json
```

Encode from stdin:

```bash
echo '{"Header":{"SessionID":"ChssPQ==","TimeStamp":1672531200},"EVCCID":"ChssPU5f"}' | ./v2gcodec encode -type SessionSetupReq
```

Encode to binary file (not hex):

```bash
./v2gcodec encode -type SessionSetupReq -hex=false -in message.json -out message.exi
```

### Round-trip Test

Verify encoding and decoding work correctly:

```bash
./v2gcodec encode -type SessionSetupReq '{"Header":{"SessionID":"ChssPQ==","TimeStamp":1672531200},"EVCCID":"ChssPU5f"}' | ./v2gcodec decode
```

### Quick Usage with `go run`

You can also use `go run` without building:

```bash
go run cmd/v2gcodec/main.go decode "808c02050d961e8809ac39d06204050d961ea72f80"
```

## Supported Message Types

All 26 ISO 15118-20 message types are supported:

**Session Management:**
- SessionSetupReq, SessionSetupRes
- SessionStopReq, SessionStopRes

**Authorization:**
- AuthorizationSetupReq, AuthorizationSetupRes
- AuthorizationReq, AuthorizationRes

**Service Discovery & Selection:**
- ServiceDiscoveryReq, ServiceDiscoveryRes
- ServiceDetailReq, ServiceDetailRes
- ServiceSelectionReq, ServiceSelectionRes

**Charging:**
- ScheduleExchangeReq, ScheduleExchangeRes
- PowerDeliveryReq, PowerDeliveryRes
- MeteringConfirmationReq, MeteringConfirmationRes

**Certificate Management:**
- CertificateInstallationReq, CertificateInstallationRes

**Vehicle Check-In/Out:**
- VehicleCheckInReq, VehicleCheckInRes
- VehicleCheckOutReq, VehicleCheckOutRes

## Examples with Real Messages

### SessionSetupReq

```bash
./v2gcodec decode "808c02050d961e8809ac39d06204050d961ea72f80"
```

### ServiceDiscoveryRes

```bash
./v2gcodec decode "808002050d961e8809ac39d062004000024000040800190900"
```

Output shows multiple energy transfer services and value-added services:
```json
{
  "Header": {
    "SessionID": "ChssPQ==",
    "TimeStamp": 1672531200
  },
  "ResponseCode": "OK",
  "ServiceRenegotiationSupported": true,
  "EnergyTransferServiceList": {
    "Services": [
      {"ServiceID": 1, "FreeService": true},
      {"ServiceID": 2, "FreeService": false}
    ]
  },
  "VASList": {
    "Services": [
      {"ServiceID": 100, "FreeService": true}
    ]
  }
}
```

## Performance

- **Decode**: ~370ns per message (2.7M ops/sec)
- **Encode**: ~600ns per message (1.6M ops/sec)
- **Memory**: ~4KB per operation
- **Binary size**: ~8MB (static binary, zero dependencies)

## Implementation

This tool uses the pure Go EXI encoder/decoder from `pkg/exi`, which:

- Implements ISO 15118-20 specific EXI grammar
- Direct struct ↔ EXI encoding (no XML parsing)
- Zero external dependencies
- 33x faster than general-purpose EXI implementations
- Thoroughly tested with golden file validation

## Exit Codes

- `0`: Success
- `1`: Usage error
- `2`: I/O error
- `3`: Processing error (encoding/decoding failed)
