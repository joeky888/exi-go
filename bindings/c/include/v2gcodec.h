/*
 * v2gcodec.h - C API for exi-go EXI codec
 *
 * Copyright 2025
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *
 * This header declares a small, stable C API for the v2g EXI codec
 * implemented in Go (built as a C-shared library). The API is intentionally
 * minimal and focused on encoding/decoding operations and runtime lifecycle.
 *
 * Notes:
 *  - All output buffers returned from encode/decode functions are allocated
 *    by the library. Call `v2g_free_buffer` to release them.
 *  - Functions return 0 (V2G_OK) on success, non-zero error codes otherwise.
 *  - The library maintains a thread-local last error string accessible via
 *    `v2g_last_error()`; copy the string if you need it to survive further
 *    library calls.
 *  - The API is designed to be simple and portable. Future versions may add
 *    more advanced configuration hooks.
 */

#ifndef EXIGO_V2GCODEC_H
#define EXIGO_V2GCODEC_H

#include <stddef.h> /* for size_t */
#include <stdint.h> /* for uint8_t */

#ifdef __cplusplus
extern "C" {
#endif

/* Error / return codes */
enum v2g_status {
  V2G_OK = 0,              /* success */
  V2G_ERR_INIT = 1,        /* initialization error */
  V2G_ERR_SHUTDOWN = 2,    /* shutdown error */
  V2G_ERR_INVALID_ARG = 3, /* invalid parameter */
  V2G_ERR_ENCODE = 4,      /* encode failure */
  V2G_ERR_DECODE = 5,      /* decode failure */
  V2G_ERR_SCHEMA = 6,      /* schema / grammar error */
  V2G_ERR_OOM = 7,         /* out of memory */
  V2G_ERR_INTERNAL = 254   /* internal/unclassified error */
};

/*
 * v2g_init
 *
 * Initialize global runtime state. Must be called before other API calls
 * except `v2g_version` and `v2g_last_error`.
 *
 * Returns:
 *   V2G_OK on success, or an error code on failure.
 */
int v2g_init(void);

/*
 * v2g_shutdown
 *
 * Release global resources and perform an orderly shutdown. After a successful
 * shutdown, v2g_init may be called again to reinitialize.
 *
 * Returns:
 *   V2G_OK on success, or an error code on failure.
 */
int v2g_shutdown(void);

/*
 * v2g_load_schemas
 *
 * Load and register one or more XSD schema files into the runtime. The runtime
 * will use loaded schemas to build schema-informed EXI grammars used during
 * encode/decode operations.
 *
 * Parameters:
 *   paths - pointer to an array of NUL-terminated UTF-8 file path strings
 *   count - number of entries in `paths`
 *
 * Ownership:
 *   The callee does not take ownership of the `paths` strings; they must remain
 *   valid for the duration of the call.
 *
 * Returns:
 *   V2G_OK on success, or V2G_ERR_SCHEMA / other error code on failure.
 *
 * Notes:
 *   - Some XSD files (e.g. ISO-distributed schemas) may be subject to external
 *     licensing; the library does not attempt to redistribute schemas.
 */
int v2g_load_schemas(const char **paths, size_t count);

/*
 * v2g_encode_xml
 *
 * Encode the provided XML document (UTF-8 bytes) into EXI bytes.
 *
 * Parameters:
 *   xml       - pointer to XML bytes
 *   xml_len   - length of xml in bytes
 *   out_exi   - pointer to a uint8_t* variable that will receive the address
 *               of the allocated EXI buffer on success
 *   out_len   - pointer to a size_t that will receive the length of the buffer
 *
 * On success:
 *   - returns V2G_OK
 *   - *out_exi points to a heap-allocated buffer owned by the library
 *   - *out_len contains the buffer length
 *
 * On failure:
 *   - returns a non-zero error code
 *   - *out_exi and *out_len are unspecified
 *
 * Memory ownership:
 *   Caller must call `v2g_free_buffer(*out_exi)` to release the buffer when
 *   it is no longer needed.
 */
int v2g_encode_xml(const uint8_t *xml, size_t xml_len, uint8_t **out_exi,
                   size_t *out_len);

/*
 * v2g_decode_exi
 *
 * Decode the provided EXI bytes into an XML document (UTF-8).
 *
 * Parameters:
 *   exi       - pointer to EXI bytes
 *   exi_len   - length of exi in bytes
 *   out_xml   - pointer to a char* variable that will receive the address of
 *               the allocated NUL-terminated XML string on success
 *   out_len   - pointer to a size_t that will receive the byte-length (not
 *               counting the trailing NUL) of the returned XML
 *
 * On success:
 *   - returns V2G_OK
 *   - *out_xml points to a NUL-terminated heap-allocated C string owned by the
 *     library
 *   - *out_len contains the length of the XML payload (bytes, excluding NUL)
 *
 * Memory ownership:
 *   Caller must call `v2g_free_buffer(*out_xml)` to release the returned
 * string.
 */
int v2g_decode_exi(const uint8_t *exi, size_t exi_len, char **out_xml,
                   size_t *out_len);

/*
 * v2g_free_buffer
 *
 * Free a buffer previously allocated and returned by the library. This is the
 * single function callers must use to release memory allocated by the
 * library in encode/decode operations.
 *
 * Parameters:
 *   buf - pointer returned by the library (may be NULL)
 *
 * Behavior:
 *   - If buf is NULL, the function does nothing.
 *   - After freeing, the caller must not use buf again.
 */
void v2g_free_buffer(void *buf);

/*
 * v2g_last_error
 *
 * Return a human-readable error message string describing the most recent
 * error on the calling thread. The returned pointer is valid until the next
 * API call on the same thread (or until shutdown). The pointer is owned by
 * the library and must not be freed by the caller.
 *
 * Returns:
 *   NUL-terminated C string or NULL if no error message is available.
 */
const char *v2g_last_error(void);

/*
 * v2g_version
 *
 * Return a statically allocated NUL-terminated string describing the library
 * version (e.g. "vX.Y.Z" or "dev"). This pointer is valid for the lifetime of
 * the process and must not be freed by the caller.
 */
const char *v2g_version(void);

/*
 * v2g_set_option (optional / experimental)
 *
 * Set a runtime option by name. Intended for experimental and advanced
 * configuration of the runtime (e.g. toggling stub mode or enabling debug
 * logging). Options and values are interpreted by the implementation.
 *
 * Parameters:
 *   name  - NUL-terminated option name
 *   value - NUL-terminated option value
 *
 * Returns:
 *   V2G_OK on success or an appropriate error code. Unknown option names
 *   should return V2G_ERR_INVALID_ARG.
 */
int v2g_set_option(const char *name, const char *value);

/* Thread-safety:
 *
 * - The library supports concurrent calls from multiple threads for encode /
 *   decode operations provided that v2g_init has been called and v2g_shutdown
 *   is not concurrently executing. Certain implementations may require
 *   calling v2g_init from the main thread; consult implementation notes.
 *
 * - The last-error string accessed via v2g_last_error() is thread-local when
 *   possible. Do not assume it is global across threads.
 */

/*
 * v2g_encode_struct
 *
 * Encode a message structure directly to EXI bytes using JSON as the
 * intermediate representation. This bypasses XML encoding and provides
 * a more efficient path for applications working with structured data.
 *
 * Parameters:
 *   msg_type  - message type identifier (see V2G_MSG_* constants below)
 *   json_data - pointer to JSON-encoded message structure
 *   json_len  - length of json_data in bytes
 *   out_exi   - pointer to a uint8_t* that will receive the EXI buffer
 *   out_len   - pointer to a size_t that will receive the EXI buffer length
 *
 * Returns:
 *   V2G_OK on success, or an error code on failure.
 *
 * Memory ownership:
 *   Caller must call v2g_free_buffer(*out_exi) to release the buffer.
 */
int v2g_encode_struct(int msg_type, const char *json_data, size_t json_len,
                      uint8_t **out_exi, size_t *out_len);

/*
 * v2g_decode_struct
 *
 * Decode EXI bytes directly to a message structure in JSON format.
 * This bypasses XML decoding for more efficient structured data access.
 *
 * Parameters:
 *   msg_type  - message type identifier (see V2G_MSG_* constants below)
 *   exi_data  - pointer to EXI bytes
 *   exi_len   - length of exi_data in bytes
 *   out_json  - pointer to a char* that will receive the JSON string
 *   out_len   - pointer to a size_t that will receive the JSON length
 *
 * Returns:
 *   V2G_OK on success, or an error code on failure.
 *
 * Memory ownership:
 *   Caller must call v2g_free_buffer(*out_json) to release the string.
 */
int v2g_decode_struct(int msg_type, const uint8_t *exi_data, size_t exi_len,
                      char **out_json, size_t *out_len);

/*
 * v2g_message_type_name
 *
 * Get the human-readable name for a message type constant.
 *
 * Parameters:
 *   msg_type - message type identifier
 *
 * Returns:
 *   NUL-terminated string with the message type name. The string is
 *   dynamically allocated and must be freed with v2g_free_buffer().
 */
const char *v2g_message_type_name(int msg_type);

/* Message type constants for ISO 15118-20 CommonMessages */
#define V2G_MSG_AuthorizationReq 0
#define V2G_MSG_AuthorizationRes 1
#define V2G_MSG_AuthorizationSetupReq 2
#define V2G_MSG_AuthorizationSetupRes 3
#define V2G_MSG_CLReqControlMode 4
#define V2G_MSG_CLResControlMode 5
#define V2G_MSG_CertificateInstallationReq 7
#define V2G_MSG_CertificateInstallationRes 8
#define V2G_MSG_MeteringConfirmationReq 16
#define V2G_MSG_MeteringConfirmationRes 17
#define V2G_MSG_PowerDeliveryReq 21
#define V2G_MSG_PowerDeliveryRes 22
#define V2G_MSG_ScheduleExchangeReq 27
#define V2G_MSG_ScheduleExchangeRes 28
#define V2G_MSG_ServiceDetailReq 29
#define V2G_MSG_ServiceDetailRes 30
#define V2G_MSG_ServiceDiscoveryReq 31
#define V2G_MSG_ServiceDiscoveryRes 32
#define V2G_MSG_ServiceSelectionReq 33
#define V2G_MSG_ServiceSelectionRes 34
#define V2G_MSG_SessionSetupReq 35
#define V2G_MSG_SessionSetupRes 36
#define V2G_MSG_SessionStopReq 37
#define V2G_MSG_SessionStopRes 38
#define V2G_MSG_VehicleCheckInReq 49
#define V2G_MSG_VehicleCheckInRes 50
#define V2G_MSG_VehicleCheckOutReq 51
#define V2G_MSG_VehicleCheckOutRes 52

#ifdef __cplusplus
} /* extern "C" */
#endif

#endif /* EXIGO_V2GCODEC_H */
