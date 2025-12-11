package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"example.com/exi-go/pkg/exi"
	"example.com/exi-go/pkg/schema"
	"example.com/exi-go/pkg/v2g/generated"
)

// v is populated at build time with -ldflags "-X main.version=..."
var version = "dev"

// Exit codes
const (
	ExitOK = iota
	ExitErrUsage
	ExitErrIO
	ExitErrProcessing
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		usage()
		os.Exit(ExitErrUsage)
	}

	switch os.Args[1] {
	case "help", "-h", "--help":
		usage()
		os.Exit(ExitOK)
	case "version", "--version", "-v":
		fmt.Println("v2gcodec", version)
		os.Exit(ExitOK)
	case "encode":
		if err := runEncode(os.Args[2:]); err != nil {
			log.Println("encode:", err)
			os.Exit(mapErrorToCode(err))
		}
	case "decode":
		if err := runDecode(os.Args[2:]); err != nil {
			log.Println("decode:", err)
			os.Exit(mapErrorToCode(err))
		}
	case "generate":
		if err := runGenerate(os.Args[2:]); err != nil {
			log.Println("generate:", err)
			os.Exit(mapErrorToCode(err))
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n\n", os.Args[1])
		usage()
		os.Exit(ExitErrUsage)
	}
}

// usage prints the top-level usage text.
func usage() {
	prog := filepath.Base(os.Args[0])
	fmt.Printf(`%s - ISO 15118-20 EXI codec (Go) - encoder/decoder CLI

Usage:
  %s <command> [options] [args]

Commands:
  encode    Encode JSON message into EXI binary (hex format)
  decode    Decode EXI binary (hex format) into JSON
  generate  Generate Go types and encoder/decoder stubs from XSDs
  version   Print version
  help      Print this help

Decode Examples:
  # decode hex string (direct argument)
  %s decode "808c02050d961e8809ac39d06204050d961ea72f80"

  # decode from file
  %s decode -in message.exi

  # decode from stdin
  cat message.hex | %s decode

Encode Examples:
  # encode JSON to hex (direct argument)
  %s encode -type SessionSetupReq '{"Header":{"SessionID":"ChssPQ==","TimeStamp":1672531200},"EVCCID":"ChssPU5f"}'

  # encode from file
  %s encode -type SessionSetupReq -in message.json

  # encode from stdin
  echo '{"Header":{...}}' | %s encode -type SessionSetupReq

  # encode to binary file (not hex)
  %s encode -type SessionSetupReq -hex=false -in message.json -out message.exi

Round-trip Example:
  # encode then decode
  %s encode -type SessionSetupReq '{"Header":{"SessionID":"ChssPQ==","TimeStamp":1672531200},"EVCCID":"ChssPU5f"}' | %s decode

Supported Message Types:
  SessionSetupReq, SessionSetupRes, AuthorizationSetupReq, AuthorizationSetupRes,
  AuthorizationReq, AuthorizationRes, ServiceDiscoveryReq, ServiceDiscoveryRes,
  ServiceDetailReq, ServiceDetailRes, ServiceSelectionReq, ServiceSelectionRes,
  ScheduleExchangeReq, ScheduleExchangeRes, PowerDeliveryReq, PowerDeliveryRes,
  MeteringConfirmationReq, MeteringConfirmationRes, SessionStopReq, SessionStopRes,
  CertificateInstallationReq, CertificateInstallationRes, VehicleCheckInReq,
  VehicleCheckInRes, VehicleCheckOutReq, VehicleCheckOutRes

`, prog, prog, prog, prog, prog, prog, prog, prog, prog, prog, prog)
}

// mapErrorToCode converts an error into an exit code.
func mapErrorToCode(err error) int {
	if err == nil {
		return ExitOK
	}
	// Simple heuristics: treat os.PathError and io errors as IO problems.
	if _, ok := err.(*os.PathError); ok {
		return ExitErrIO
	}
	if err == io.EOF {
		return ExitErrIO
	}
	return ExitErrProcessing
}

// runEncode handles the `encode` subcommand.
func runEncode(args []string) error {
	fs := flag.NewFlagSet("encode", flag.ContinueOnError)
	inPath := fs.String("in", "-", "Input file path (JSON). Use '-' for stdin")
	outPath := fs.String("out", "-", "Output file path (EXI binary). Use '-' for stdout")
	asHex := fs.Bool("hex", true, "Write hex-encoded EXI instead of raw binary (default)")
	msgType := fs.String("type", "", "Message type (e.g., SessionSetupReq, SessionSetupRes)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Support direct JSON as positional argument
	var inputBytes []byte
	if fs.NArg() > 0 {
		// First positional arg is JSON string
		inputBytes = []byte(fs.Arg(0))
	} else {
		// Read from file or stdin
		in, err := openInput(*inPath)
		if err != nil {
			return err
		}
		defer closeIfFile(in)

		inputBytes, err = io.ReadAll(in)
		if err != nil {
			return err
		}
	}

	// Parse JSON into appropriate message type
	var msg interface{}
	if *msgType != "" {
		msg = getMessageTypeByName(*msgType)
		if msg == nil {
			return fmt.Errorf("unknown message type: %s", *msgType)
		}
	} else {
		// Try to auto-detect message type from JSON
		var generic map[string]interface{}
		if err := json.Unmarshal(inputBytes, &generic); err != nil {
			return fmt.Errorf("invalid JSON input: %w", err)
		}
		msg = getMessageTypeByName("SessionSetupReq") // default for now
	}

	// Unmarshal JSON into the message struct
	if err := json.Unmarshal(inputBytes, msg); err != nil {
		return fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	// Encode using real EXI encoder
	exiBytes, err := exi.EncodeStruct(msg)
	if err != nil {
		return fmt.Errorf("EXI encode failed: %w", err)
	}

	// Output result
	out, outIsFile, err := openOutput(*outPath)
	if err != nil {
		return err
	}
	if outIsFile {
		defer closeIfFile(out)
	}

	if *asHex {
		encoded := hex.EncodeToString(exiBytes)
		_, err = io.WriteString(out, encoded+"\n")
	} else {
		_, err = out.Write(exiBytes)
	}
	return err
}

// runDecode handles the `decode` subcommand.
func runDecode(args []string) error {
	fs := flag.NewFlagSet("decode", flag.ContinueOnError)
	inPath := fs.String("in", "-", "Input file path (EXI binary or hex). Use '-' for stdin")
	outPath := fs.String("out", "-", "Output file path (JSON). Use '-' for stdout")
	asJSON := fs.Bool("json", true, "Output as JSON (default)")
	// future flags: schema, pretty, xml
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Support direct hex string as positional argument
	var exiBytes []byte
	if fs.NArg() > 0 {
		// First positional arg is hex string
		hexStr := strings.TrimSpace(fs.Arg(0))
		decoded, err := hex.DecodeString(hexStr)
		if err != nil {
			return fmt.Errorf("invalid hex string: %w", err)
		}
		exiBytes = decoded
	} else {
		// Read from file or stdin
		in, err := openInput(*inPath)
		if err != nil {
			return err
		}
		defer closeIfFile(in)

		inputBytes, err := io.ReadAll(in)
		if err != nil {
			return err
		}

		// Auto-detect if input looks like hex (all printable hex chars)
		trimmed := strings.TrimSpace(string(inputBytes))
		if isHexString(trimmed) {
			decoded, err := hex.DecodeString(trimmed)
			if err != nil {
				return fmt.Errorf("invalid hex input: %w", err)
			}
			exiBytes = decoded
		} else {
			exiBytes = inputBytes
		}
	}

	// Decode using real EXI decoder
	msg, err := exi.DecodeStruct(exiBytes, nil)
	if err != nil {
		return fmt.Errorf("EXI decode failed: %w", err)
	}

	// Output result
	out, outIsFile, err := openOutput(*outPath)
	if err != nil {
		return err
	}
	if outIsFile {
		defer closeIfFile(out)
	}

	if *asJSON {
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(msg); err != nil {
			return fmt.Errorf("JSON encode failed: %w", err)
		}
	} else {
		// For now, output JSON even without -json flag
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(msg); err != nil {
			return fmt.Errorf("JSON encode failed: %w", err)
		}
	}

	return nil
}

// runGenerate handles the `generate` subcommand.
// This is a placeholder which will later implement XSD parsing and codegen.
func runGenerate(args []string) error {
	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	schemaDir := fs.String("schema", "", "Directory containing XSD files (optional if --auto-download)")
	autoDownload := fs.Bool("auto-download", false, "Download public ISO 15118 XSDs into download-dir (must accept ISO terms)")
	acceptISO := fs.Bool("accept-iso", false, "Accept ISO schema terms when using --auto-download")
	downloadDir := fs.String("download-dir", "./schemas", "Destination directory for downloaded schemas")
	outDir := fs.String("out", "./generated", "Output directory for generated Go code")
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Determine schema source(s): either provided schema dir(s) or auto-download.
	var schemaPaths []string
	if *schemaDir != "" {
		schemaPaths = append(schemaPaths, *schemaDir)
	}

	if *autoDownload {
		// User explicitly requested downloading ISO schemas. They must accept ISO terms.
		downloaded, err := schema.DownloadISO15118Schemas(*downloadDir, *acceptISO)
		if err != nil {
			return fmt.Errorf("failed to download ISO schemas: %w", err)
		}
		// Append all downloaded file paths to schemaPaths
		schemaPaths = append(schemaPaths, downloaded...)
	}

	if len(schemaPaths) == 0 {
		return fmt.Errorf("no schema source provided; specify -schema or use --auto-download")
	}

	// Validate schema paths exist (files or directories)
	for _, p := range schemaPaths {
		if _, err := os.Stat(p); err != nil {
			return fmt.Errorf("schema path error for %s: %w", p, err)
		}
	}

	// Ensure output directory exists
	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	// Defer to schema.GenerateTypes to perform generation. For now the package
	// will create a placeholder if no Generator is supplied.
	if err := schema.RunSimpleGeneration(schemaPaths, *outDir, "exi-go"); err != nil {
		return fmt.Errorf("generate types: %w", err)
	}

	// created placeholder or real generated output by GenerateTypes
	fmt.Fprintf(os.Stderr, "Generation completed, output in %s\n", *outDir)
	return nil
}

// Helpers ---------------------------------------------------------

// openInput opens input path or returns stdin if "-" is passed.
func openInput(path string) (io.ReadCloser, error) {
	if path == "-" || path == "" {
		// wrap os.Stdin with NopCloser
		return io.NopCloser(os.Stdin), nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// openOutput opens output path or returns stdout if "-" is passed.
// Returns (writer, isFile, error)
func openOutput(path string) (io.Writer, bool, error) {
	if path == "-" || path == "" {
		return os.Stdout, false, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, false, err
	}
	return f, true, nil
}

// closeIfFile closes the reader/writer if it's an *os.File (and not stdin/stdout).
// Accepts any value and only closes when the value implements io.Closer. It
// avoids closing standard streams (stdin/stdout/stderr).
func closeIfFile(v interface{}) {
	if v == nil {
		return
	}
	// Only close values that implement io.Closer, but avoid closing standard files.
	if c, ok := v.(io.Closer); ok {
		if f, ok := c.(*os.File); ok {
			// don't close standard streams
			if f == os.Stdin || f == os.Stdout || f == os.Stderr {
				return
			}
			_ = f.Close()
			return
		}
		// For other closers, attempt to close.
		_ = c.Close()
	}
}

// isHexString checks if a string contains only valid hex characters
func isHexString(s string) bool {
	if len(s) == 0 || len(s)%2 != 0 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// getMessageTypeByName returns a pointer to a new instance of the message type
func getMessageTypeByName(typeName string) interface{} {
	switch typeName {
	case "SessionSetupReq":
		return &generated.SessionSetupReq{}
	case "SessionSetupRes":
		return &generated.SessionSetupRes{}
	case "AuthorizationSetupReq":
		return &generated.AuthorizationSetupReq{}
	case "AuthorizationSetupRes":
		return &generated.AuthorizationSetupRes{}
	case "AuthorizationReq":
		return &generated.AuthorizationReq{}
	case "AuthorizationRes":
		return &generated.AuthorizationRes{}
	case "ServiceDiscoveryReq":
		return &generated.ServiceDiscoveryReq{}
	case "ServiceDiscoveryRes":
		return &generated.ServiceDiscoveryRes{}
	case "ServiceDetailReq":
		return &generated.ServiceDetailReq{}
	case "ServiceDetailRes":
		return &generated.ServiceDetailRes{}
	case "ServiceSelectionReq":
		return &generated.ServiceSelectionReq{}
	case "ServiceSelectionRes":
		return &generated.ServiceSelectionRes{}
	case "ScheduleExchangeReq":
		return &generated.ScheduleExchangeReq{}
	case "ScheduleExchangeRes":
		return &generated.ScheduleExchangeRes{}
	case "PowerDeliveryReq":
		return &generated.PowerDeliveryReq{}
	case "PowerDeliveryRes":
		return &generated.PowerDeliveryRes{}
	case "MeteringConfirmationReq":
		return &generated.MeteringConfirmationReq{}
	case "MeteringConfirmationRes":
		return &generated.MeteringConfirmationRes{}
	case "SessionStopReq":
		return &generated.SessionStopReq{}
	case "SessionStopRes":
		return &generated.SessionStopRes{}
	case "CertificateInstallationReq":
		return &generated.CertificateInstallationReq{}
	case "CertificateInstallationRes":
		return &generated.CertificateInstallationRes{}
	case "VehicleCheckInReq":
		return &generated.VehicleCheckInReq{}
	case "VehicleCheckInRes":
		return &generated.VehicleCheckInRes{}
	case "VehicleCheckOutReq":
		return &generated.VehicleCheckOutReq{}
	case "VehicleCheckOutRes":
		return &generated.VehicleCheckOutRes{}
	default:
		return nil
	}
}
