#!/bin/bash
# Build script for exi-go C shared library
# This script compiles the Go code into a C-shared library that can be
# used from C, C++, Python, and other languages.

set -e

# Detect OS and set library extension
OS=$(uname -s)
case "$OS" in
    Linux*)
        LIB_EXT="so"
        ;;
    Darwin*)
        LIB_EXT="dylib"
        ;;
    CYGWIN*|MINGW*|MSYS*)
        LIB_EXT="dll"
        ;;
    *)
        echo "Unknown OS: $OS"
        exit 1
        ;;
esac

# Build configuration
BUILD_DIR="$(cd "$(dirname "$0")" && pwd)"
OUTPUT_DIR="${BUILD_DIR}/lib"
LIB_NAME="libv2gcodec.${LIB_EXT}"
OUTPUT_PATH="${OUTPUT_DIR}/${LIB_NAME}"

echo "Building exi-go C shared library..."
echo "OS: $OS"
echo "Library extension: $LIB_EXT"
echo "Output: $OUTPUT_PATH"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Build the shared library
cd "$BUILD_DIR"
go build -buildmode=c-shared -o "$OUTPUT_PATH" \
    -ldflags="-s -w" \
    v2gcodec.go v2gcodec_native.go

echo ""
echo "Build complete!"
echo "Library: $OUTPUT_PATH"
echo "Header: ${OUTPUT_PATH%.${LIB_EXT}}.h"
echo ""
echo "To use from C/C++:"
echo "  #include \"include/v2gcodec.h\""
echo "  gcc your_app.c -L${OUTPUT_DIR} -lv2gcodec -o your_app"
echo ""
echo "To use from Python:"
echo "  export V2G_CODEC_LIBRARY=${OUTPUT_PATH}"
echo "  python3 ../python/cffi/v2gcodec_cffi.py"
echo ""
echo "To install system-wide (Linux/macOS):"
echo "  sudo cp ${OUTPUT_PATH} /usr/local/lib/"
echo "  sudo cp include/v2gcodec.h /usr/local/include/"
echo "  sudo ldconfig  # Linux only"
echo ""
