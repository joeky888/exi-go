"""
cffi-based Python wrapper for the exi-go v2g EXI codec C API.

This file provides a small, ergonomic wrapper around the C API exposed by the
Go-built shared library (libv2gcodec). It is intended as an example and
starting point; in production you should add robust error handling, logging,
and packaging (wheel) logic.

Usage (example):
    from v2gcodec_cffi import V2GCodec

    codec = V2GCodec()            # loads the shared library
    codec.init()
    codec.load_schemas(['/path/to/iso15118.xsd'])
    exi = codec.encode_xml(b'<SomeXml>...</SomeXml>')
    xml = codec.decode_exi(exi)
    codec.shutdown()

The wrapper looks for the shared library in the following order:
  - Path in environment variable V2G_CODEC_LIBRARY
  - libv2gcodec.so (Linux)
  - libv2gcodec.dylib (macOS)
  - v2gcodec.dll (Windows)
You can pass an explicit path when constructing V2GCodec.
"""

import os
import sys

from cffi import FFI

ffi = FFI()

# C function declarations (keep in sync with exi-go/bindings/c/include/v2gcodec.h)
ffi.cdef("""
    typedef unsigned long size_t; /* ensure size_t is available to cffi */

    int v2g_init(void);
    int v2g_shutdown(void);
    int v2g_load_schemas(const char** paths, size_t count);
    int v2g_encode_xml(const unsigned char* xml, size_t xml_len,
                       unsigned char** out_exi, size_t* out_len);
    int v2g_decode_exi(const unsigned char* exi, size_t exi_len,
                       char** out_xml, size_t* out_len);
    int v2g_encode_struct(int msg_type, const char* json_data, size_t json_len,
                          unsigned char** out_exi, size_t* out_len);
    int v2g_decode_struct(int msg_type, const unsigned char* exi_data, size_t exi_len,
                          char** out_json, size_t* out_len);
    const char* v2g_message_type_name(int msg_type);
    void v2g_free_buffer(void* buf);
    const char* v2g_last_error(void);
    const char* v2g_version(void);
    int v2g_set_option(const char* name, const char* value);
""")

# Default library names to try if none provided.
_DEFAULT_LIB_CANDIDATES = [
    "libv2gcodec.so",  # Linux
    "libv2gcodec.dylib",  # macOS
    "v2gcodec.dll",  # Windows
]


# ISO 15118-20 message type constants
class MessageType:
    """Message type constants for ISO 15118-20."""

    # Common messages (all services)
    AuthorizationReq = 0
    AuthorizationRes = 1
    AuthorizationSetupReq = 2
    AuthorizationSetupRes = 3
    CLReqControlMode = 4
    CLResControlMode = 5
    CertificateInstallationReq = 7
    CertificateInstallationRes = 8
    MeteringConfirmationReq = 16
    MeteringConfirmationRes = 17
    PowerDeliveryReq = 21
    PowerDeliveryRes = 22
    ScheduleExchangeReq = 27
    ScheduleExchangeRes = 28
    ServiceDetailReq = 29
    ServiceDetailRes = 30
    ServiceDiscoveryReq = 31
    ServiceDiscoveryRes = 32
    ServiceSelectionReq = 33
    ServiceSelectionRes = 34
    SessionSetupReq = 35
    SessionSetupRes = 36
    SessionStopReq = 37
    SessionStopRes = 38
    VehicleCheckInReq = 49
    VehicleCheckInRes = 50
    VehicleCheckOutReq = 51
    VehicleCheckOutRes = 52

    # WPT (Wireless Power Transfer) messages
    WPT_AlignmentCheckReq = 53
    WPT_AlignmentCheckRes = 54
    WPT_FinePositioningReq = 55
    WPT_FinePositioningRes = 56
    WPT_ChargeLoopReq = 57
    WPT_ChargeLoopRes = 58

    # ACDP (AC Dynamic Power) messages
    DC_ACDPReq = 59
    DC_ACDPRes = 60
    DC_ACDP_BPTReq = 61
    DC_ACDP_BPTRes = 62


class V2GError(RuntimeError):
    """Base exception for v2g codec errors."""


class V2GCodec:
    """
    Thin Python wrapper around the C API.

    Methods raise V2GError on failure and return Python-native types on success.
    """

    def __init__(self, lib_path: str = None):
        """
        Load the shared library.

        :param lib_path: Optional explicit path to the shared library. If None,
                         the wrapper searches environment variables and defaults.
        """
        lib_path = lib_path or os.environ.get("V2G_CODEC_LIBRARY")
        self._lib = None

        tried = []
        if lib_path:
            tried.append(lib_path)
            try:
                self._lib = ffi.dlopen(lib_path)
            except OSError as e:
                raise V2GError(f"failed to load library at {lib_path}: {e}") from e
        else:
            # try default names
            for name in _DEFAULT_LIB_CANDIDATES:
                tried.append(name)
                try:
                    self._lib = ffi.dlopen(name)
                    break
                except OSError:
                    continue
        if self._lib is None:
            raise V2GError(f"could not load v2g codec library; tried: {tried}")

    # ---- low-level helpers ----
    def _last_error(self) -> str:
        p = self._lib.v2g_last_error()
        if p == ffi.NULL:
            return "unknown error"
        try:
            return ffi.string(p).decode("utf-8", errors="replace")
        except Exception:
            return "<invalid error string>"

    def _check_status(self, code: int, context: str = ""):
        if code == 0:
            return
        # obtain last error from library if available
        msg = self._last_error()
        if context:
            raise V2GError(f"{context}: {msg} (code {code})")
        raise V2GError(f"{msg} (code {code})")

    # ---- lifecycle ----
    def init(self):
        """Initialize the codec runtime."""
        rc = self._lib.v2g_init()
        self._check_status(rc, "v2g_init")

    def shutdown(self):
        """Shutdown the codec runtime and free resources."""
        rc = self._lib.v2g_shutdown()
        self._check_status(rc, "v2g_shutdown")

    def version(self) -> str:
        """Return the library version string."""
        p = self._lib.v2g_version()
        if p == ffi.NULL:
            return "unknown"
        return ffi.string(p).decode("utf-8", errors="replace")

    def set_option(self, name: str, value: str):
        """Set a runtime option by name (experimental)."""
        rc = self._lib.v2g_set_option(name.encode("utf-8"), value.encode("utf-8"))
        self._check_status(rc, "v2g_set_option")

    # ---- schema loading ----
    def load_schemas(self, paths):
        """
        Load schema files into the runtime.

        :param paths: iterable of file path strings (UTF-8)
        """
        if not paths:
            raise ValueError("paths must be a non-empty iterable of file paths")
        arr = []
        for p in paths:
            if p is None:
                arr.append(ffi.NULL)
            else:
                arr.append(ffi.new("char[]", p.encode("utf-8")))
        c_array = ffi.new("char*[]", arr)
        rc = self._lib.v2g_load_schemas(c_array, len(arr))
        self._check_status(rc, "v2g_load_schemas")

    # ---- encode / decode ----
    def encode_xml(self, xml_bytes: bytes) -> bytes:
        """
        Encode XML bytes into EXI bytes.

        :param xml_bytes: bytes containing UTF-8 XML
        :returns: EXI payload as bytes
        :raises V2GError on failure
        """
        if not isinstance(xml_bytes, (bytes, bytearray)):
            raise TypeError("xml_bytes must be bytes")
        in_buf = ffi.from_buffer(xml_bytes)
        out_exi_ptr = ffi.new("unsigned char**")
        out_len = ffi.new("size_t*")
        rc = self._lib.v2g_encode_xml(in_buf, len(xml_bytes), out_exi_ptr, out_len)
        if rc != 0:
            self._check_status(rc, "v2g_encode_xml")
        # out_exi_ptr[0] is a pointer to heap memory allocated by library
        out_ptr = out_exi_ptr[0]
        length = int(out_len[0])
        if out_ptr == ffi.NULL or length == 0:
            # still consider it success but return empty bytes
            return b""
        try:
            result = bytes(ffi.buffer(out_ptr, length))
            return result
        finally:
            # free the buffer allocated by the C library
            self._lib.v2g_free_buffer(out_ptr)

    def decode_exi(self, exi_bytes: bytes) -> str:
        """
        Decode EXI bytes into an XML UTF-8 string.

        :param exi_bytes: EXI payload bytes
        :returns: decoded XML string (UTF-8)
        """
        if not isinstance(exi_bytes, (bytes, bytearray)):
            raise TypeError("exi_bytes must be bytes")
        in_buf = ffi.from_buffer(exi_bytes)
        out_xml_ptr = ffi.new("char**")
        out_len = ffi.new("size_t*")
        rc = self._lib.v2g_decode_exi(in_buf, len(exi_bytes), out_xml_ptr, out_len)
        if rc != 0:
            self._check_status(rc, "v2g_decode_exi")
        out_ptr = out_xml_ptr[0]
        length = int(out_len[0])
        if out_ptr == ffi.NULL or length == 0:
            return ""
        try:
            # out_ptr is a NUL-terminated C string (but length is also provided)
            b = bytes(ffi.buffer(out_ptr, length))
            return b.decode("utf-8", errors="replace")
        finally:
            # free the returned string pointer
            self._lib.v2g_free_buffer(out_ptr)

    # ---- native struct encode/decode ----
    def encode_struct(self, msg_type: int, data: dict) -> bytes:
        """
        Encode a message structure directly to EXI bytes.

        :param msg_type: Message type constant from MessageType class
        :param data: Dictionary containing the message structure
        :returns: EXI payload as bytes
        :raises V2GError on failure
        """
        import json

        if not isinstance(data, dict):
            raise TypeError("data must be a dictionary")

        json_str = json.dumps(data)
        json_bytes = json_str.encode("utf-8")

        in_buf = ffi.from_buffer(json_bytes)
        out_exi_ptr = ffi.new("unsigned char**")
        out_len = ffi.new("size_t*")

        rc = self._lib.v2g_encode_struct(
            msg_type, in_buf, len(json_bytes), out_exi_ptr, out_len
        )
        if rc != 0:
            self._check_status(rc, "v2g_encode_struct")

        out_ptr = out_exi_ptr[0]
        length = int(out_len[0])
        if out_ptr == ffi.NULL or length == 0:
            return b""

        try:
            result = bytes(ffi.buffer(out_ptr, length))
            return result
        finally:
            self._lib.v2g_free_buffer(out_ptr)

    def decode_struct(self, msg_type: int, exi_bytes: bytes) -> dict:
        """
        Decode EXI bytes directly to a message structure.

        :param msg_type: Message type constant from MessageType class
        :param exi_bytes: EXI payload bytes
        :returns: Dictionary containing the decoded message structure
        :raises V2GError on failure
        """
        import json

        if not isinstance(exi_bytes, (bytes, bytearray)):
            raise TypeError("exi_bytes must be bytes")

        in_buf = ffi.from_buffer(exi_bytes)
        out_json_ptr = ffi.new("char**")
        out_len = ffi.new("size_t*")

        rc = self._lib.v2g_decode_struct(
            msg_type, in_buf, len(exi_bytes), out_json_ptr, out_len
        )
        if rc != 0:
            self._check_status(rc, "v2g_decode_struct")

        out_ptr = out_json_ptr[0]
        length = int(out_len[0])
        if out_ptr == ffi.NULL or length == 0:
            return {}

        try:
            json_bytes = bytes(ffi.buffer(out_ptr, length))
            json_str = json_bytes.decode("utf-8", errors="replace")
            return json.loads(json_str)
        finally:
            self._lib.v2g_free_buffer(out_ptr)

    def message_type_name(self, msg_type: int) -> str:
        """
        Get the human-readable name for a message type.

        :param msg_type: Message type constant
        :returns: Message type name string
        """
        p = self._lib.v2g_message_type_name(msg_type)
        if p == ffi.NULL:
            return f"Unknown({msg_type})"
        try:
            name = ffi.string(p).decode("utf-8", errors="replace")
            return name
        finally:
            self._lib.v2g_free_buffer(p)


# Simple CLI demonstration when executed directly
def _demo():
    import argparse

    parser = argparse.ArgumentParser(description="v2gcodec cffi demo")
    parser.add_argument("--lib", help="path to v2g codec shared library")
    parser.add_argument(
        "--schema", action="append", help="schema XSD path (can be used multiple times)"
    )
    args = parser.parse_args()

    codec = V2GCodec(lib_path=args.lib) if args.lib else V2GCodec()
    print("Loaded codec library:", codec.version())
    codec.init()
    try:
        if args.schema:
            print("Loading schemas:", args.schema)
            codec.load_schemas(args.schema)
        sample_xml = b"<Sample>hello</Sample>"
        print("Encoding sample XML:", sample_xml)
        exi = codec.encode_xml(sample_xml)
        print("Encoded EXI (len=%d) hex: %s" % (len(exi), exi.hex()[:256]))
        decoded = codec.decode_exi(exi)
        print("Decoded XML:", decoded)

        # Demo native struct encoding
        print("\n--- Native Struct Encoding Demo ---")
        session_setup_req = {
            "Header": {"SessionID": [0x01, 0x02, 0x03, 0x04], "TimeStamp": 1234567890},
            "EVCCID": [0x0A, 0x1B, 0x2C, 0x3D, 0x4E, 0x5F],
        }
        print("Encoding SessionSetupReq:", session_setup_req)
        struct_exi = codec.encode_struct(MessageType.SessionSetupReq, session_setup_req)
        print("Struct EXI (len=%d) hex: %s" % (len(struct_exi), struct_exi.hex()))

        decoded_struct = codec.decode_struct(MessageType.SessionSetupReq, struct_exi)
        print("Decoded struct:", decoded_struct)

        print(
            "\nMessage type name for event 35:",
            codec.message_type_name(MessageType.SessionSetupReq),
        )
    finally:
        codec.shutdown()


if __name__ == "__main__":
    # Provide a friendly message if the library cannot be loaded.
    try:
        _demo()
    except V2GError as e:
        print("v2gcodec error:", e, file=sys.stderr)
        sys.exit(2)
    except Exception as e:
        print("unexpected error:", e, file=sys.stderr)
        sys.exit(3)
