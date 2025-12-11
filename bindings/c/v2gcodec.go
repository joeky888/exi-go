/*
cgo bridge for exi-go

This file implements the C-exported API surface for the v2g EXI codec.
It provides lightweight, safe bridging between C callers and the Go runtime.

Build note:
  - This package is intended to be built as a C-shared library:
    go build -buildmode=c-shared -o libv2gcodec.so
  - The header corresponding to this implementation is provided at:
    exi-go/bindings/c/include/v2gcodec.h

License: Apache-2.0 (match repository)
*/
package main

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"

	"example.com/exi-go/pkg/exi"
)

// Global runtime state -----------------------------------------------------

var (
	// codec is the shared EXI codec instance used by C calls.
	// It is protected by stateMu when being replaced.
	codec   *exi.Codec
	stateMu sync.Mutex

	// lastErr stores the last human-readable error message. Protected by errMu.
	errMu      sync.Mutex
	lastErrStr string
	lastErrC   *C.char

	// versionC is a statically allocated C string for v2g_version.
	versionC *C.char
)

const (
	// Return codes corresponding to v2g_status in the header.
	_v2g_ok           = 0
	_v2g_err_init     = 1
	_v2g_err_shutdown = 2
	_v2g_err_invalid  = 3
	_v2g_err_encode   = 4
	_v2g_err_decode   = 5
	_v2g_err_schema   = 6
	_v2g_err_oom      = 7
	_v2g_err_internal = 254
)

// helper: set last error string (thread-safe). Keeps a C copy in lastErrC.
func setLastError(format string, a ...interface{}) {
	errMu.Lock()
	defer errMu.Unlock()
	lastErrStr = fmt.Sprintf(format, a...)
	// free previous C allocation if any
	if lastErrC != nil {
		C.free(unsafe.Pointer(lastErrC))
		lastErrC = nil
	}
	lastErrC = C.CString(lastErrStr)
}

// getLastErrorC returns the last error C string (may be nil).
func getLastErrorC() *C.char {
	errMu.Lock()
	defer errMu.Unlock()
	return lastErrC
}

// init prepares a version string and ensures no previous allocation leak.
func init() {
	version := "dev"
	versionC = C.CString(version)
	// lastErrC starts nil
}

// cleanup releases C allocations created by this package.
// Called on program exit or when needed.
func cleanup() {
	errMu.Lock()
	if lastErrC != nil {
		C.free(unsafe.Pointer(lastErrC))
		lastErrC = nil
	}
	errMu.Unlock()
	if versionC != nil {
		C.free(unsafe.Pointer(versionC))
		versionC = nil
	}
}

// EXPORTS ------------------------------------------------------------------
// Each exported function is declared immediately after a comment of the form:
// //export <Name>
// The functions use C-compatible types and minimal pointer arithmetic.

//export v2g_init
func v2g_init() C.int {
	stateMu.Lock()
	defer stateMu.Unlock()

	// If already initialized, return success.
	if codec != nil {
		return C.int(_v2g_ok)
	}

	// Create codec with default config (using stub by default).
	c := exi.NewCodec(nil)
	if err := c.Init(); err != nil {
		setLastError("init: %v", err)
		return C.int(_v2g_err_init)
	}
	codec = c
	return C.int(_v2g_ok)
}

//export v2g_shutdown
func v2g_shutdown() C.int {
	stateMu.Lock()
	defer stateMu.Unlock()

	if codec == nil {
		// Nothing to do
		return C.int(_v2g_ok)
	}
	if err := codec.Shutdown(); err != nil {
		setLastError("shutdown: %v", err)
		return C.int(_v2g_err_shutdown)
	}
	codec = nil
	// cleanup any C-allocated strings
	cleanup()
	return C.int(_v2g_ok)
}

// cStringsToGoStrings converts a C array of *C.char (paths) into a Go []string.
// pathsPtr is **C.char and count is number of entries.
func cStringsToGoStrings(paths **C.char, count C.size_t) []string {
	if paths == nil || count == 0 {
		return nil
	}
	var out []string
	// Size of pointer on this platform
	ptrSize := unsafe.Sizeof(uintptr(0))
	base := uintptr(unsafe.Pointer(paths))
	for i := 0; i < int(count); i++ {
		// compute pointer to C.char*
		elemPtr := *(**C.char)(unsafe.Pointer(base + uintptr(i)*ptrSize))
		if elemPtr == nil {
			out = append(out, "")
		} else {
			out = append(out, C.GoString(elemPtr))
		}
	}
	return out
}

//export v2g_load_schemas
func v2g_load_schemas(paths **C.char, count C.size_t) C.int {
	// Convert paths to Go slice of strings
	goPaths := cStringsToGoStrings(paths, count)
	if len(goPaths) == 0 {
		setLastError("v2g_load_schemas: no schema paths provided")
		return C.int(_v2g_err_invalid)
	}

	// For this scaffold, store the schema paths by creating a codec configured
	// with the provided paths and initialize it.
	stateMu.Lock()
	defer stateMu.Unlock()

	// Shutdown existing codec if any
	if codec != nil {
		_ = codec.Shutdown()
		codec = nil
	}

	cfg := &exi.Config{
		SchemaPaths: goPaths,
		UseStub:     true, // generator will later switch as needed
	}
	c := exi.NewCodec(cfg)
	if err := c.Init(); err != nil {
		setLastError("load_schemas: init failed: %v", err)
		return C.int(_v2g_err_schema)
	}
	codec = c
	return C.int(_v2g_ok)
}

//export v2g_encode_xml
func v2g_encode_xml(xml *C.uint8_t, xml_len C.size_t, out_exi **C.uint8_t, out_len *C.size_t) C.int {
	if xml == nil || xml_len == 0 || out_exi == nil || out_len == nil {
		setLastError("v2g_encode_xml: invalid arguments")
		return C.int(_v2g_err_invalid)
	}

	// Ensure codec initialized
	stateMu.Lock()
	c := codec
	stateMu.Unlock()
	if c == nil {
		setLastError("v2g_encode_xml: codec not initialized")
		return C.int(_v2g_err_init)
	}

	// Copy input bytes from C to Go
	input := C.GoBytes(unsafe.Pointer(xml), C.int(xml_len))

	// Perform encoding
	result, err := c.EncodeXML(input)
	if err != nil {
		setLastError("encode failed: %v", err)
		// map errors to encode error code
		return C.int(_v2g_err_encode)
	}

	// Allocate C memory and copy result into it. Use C.CBytes which allocates via malloc.
	cbuf := C.CBytes(result) // returns void*
	if cbuf == nil {
		setLastError("encode: out of memory")
		return C.int(_v2g_err_oom)
	}

	// Assign output parameters
	*out_exi = (*C.uint8_t)(cbuf)
	*out_len = C.size_t(len(result))
	return C.int(_v2g_ok)
}

//export v2g_decode_exi
func v2g_decode_exi(exiBuf *C.uint8_t, exi_len C.size_t, out_xml **C.char, out_len *C.size_t) C.int {
	if exiBuf == nil || exi_len == 0 || out_xml == nil || out_len == nil {
		setLastError("v2g_decode_exi: invalid arguments")
		return C.int(_v2g_err_invalid)
	}

	// Ensure codec initialized
	stateMu.Lock()
	c := codec
	stateMu.Unlock()
	if c == nil {
		setLastError("v2g_decode_exi: codec not initialized")
		return C.int(_v2g_err_init)
	}

	input := C.GoBytes(unsafe.Pointer(exiBuf), C.int(exi_len))

	// Perform decode
	xmlBytes, err := c.DecodeEXI(input)
	if err != nil {
		setLastError("decode failed: %v", err)
		return C.int(_v2g_err_decode)
	}

	// Allocate C memory for NUL-terminated string
	// Use malloc to allocate len+1 bytes and copy contents, set trailing NUL.
	clen := C.size_t(len(xmlBytes) + 1)
	cptr := C.malloc(clen)
	if cptr == nil {
		setLastError("decode: out of memory")
		return C.int(_v2g_err_oom)
	}

	// Copy bytes into allocated memory
	C.memcpy(cptr, unsafe.Pointer(&xmlBytes[0]), C.size_t(len(xmlBytes)))
	// Set trailing NUL (write zero at the final byte)
	lastBytePtr := unsafe.Pointer(uintptr(cptr) + uintptr(len(xmlBytes)))
	*(*byte)(lastBytePtr) = 0

	*out_xml = (*C.char)(cptr)
	*out_len = C.size_t(len(xmlBytes))
	return C.int(_v2g_ok)
}

//export v2g_free_buffer
func v2g_free_buffer(buf unsafe.Pointer) {
	if buf == nil {
		return
	}
	C.free(buf)
}

//export v2g_last_error
func v2g_last_error() *C.char {
	// Return the last error C string (may be nil)
	return getLastErrorC()
}

//export v2g_version
func v2g_version() *C.char {
	return versionC
}

//export v2g_set_option
func v2g_set_option(name *C.char, value *C.char) C.int {
	if name == nil {
		setLastError("v2g_set_option: name is nil")
		return C.int(_v2g_err_invalid)
	}
	n := C.GoString(name)
	v := ""
	if value != nil {
		v = C.GoString(value)
	}

	// support a small set of options for now
	switch n {
	case "use-stub":
		// expect "true"/"false"
		stateMu.Lock()
		if codec == nil {
			// no codec yet; create one with specified option
			cfg := &exi.Config{UseStub: v != "false"}
			c := exi.NewCodec(cfg)
			if err := c.Init(); err != nil {
				stateMu.Unlock()
				setLastError("set_option(use-stub): init failed: %v", err)
				return C.int(_v2g_err_internal)
			}
			codec = c
			stateMu.Unlock()
			return C.int(_v2g_ok)
		}
		// codec exists: try to update config via re-create (safe approach)
		_ = codec.Shutdown()
		cfg := &exi.Config{UseStub: v != "false"}
		c := exi.NewCodec(cfg)
		if err := c.Init(); err != nil {
			// attempt to restore previous codec by re-init default
			codec = exi.NewCodec(nil)
			_ = codec.Init()
			stateMu.Unlock()
			setLastError("set_option(use-stub): reinit failed: %v", err)
			return C.int(_v2g_err_internal)
		}
		codec = c
		stateMu.Unlock()
		return C.int(_v2g_ok)
	default:
		setLastError("unknown option: %s", n)
		return C.int(_v2g_err_invalid)
	}
}

// Keep a small main so the package builds as a C-shared library.
// The main function will never be executed when the library is loaded by C,
// but is required by the Go toolchain when building a c-shared package.
func main() {}
