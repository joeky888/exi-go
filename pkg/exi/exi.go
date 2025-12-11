package exi

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"errors"
	"io"
	"regexp"
)

// A simple EXI-like encoder/decoder implemented for development and tests.
// This is NOT a full EXI implementation. It provides a compact, deterministic
// binary representation used for integration tests: XML is minified then
// compressed using gzip. The resulting payload can be round-tripped by the
// corresponding DecodeEXI implementation.
//
// This lightweight approach gives us a meaningful binary format for the
// generator/runtime integration while we later implement a real EXI engine.

var (
	// ErrNotImplemented used when a non-stubbed branch is invoked
	ErrNotImplemented = errors.New("exi: not implemented")

	// ErrInvalidEXI indicates that the input bytes are not a recognized EXI payload
	ErrInvalidEXI = errors.New("exi: invalid exi data")
)

// Config configures the EXI runtime behavior.
type Config struct {
	// SchemaPaths is an ordered list of XSD file paths that the generator
	// or runtime can use to build schema-informed grammars. If empty, the
	// runtime will operate in a schema-less mode (more generic).
	SchemaPaths []string

	// UseStub controls whether the runtime uses the lightweight stub encoder.
	// For development the minify+gzip encoder is used regardless; UseStub is
	// kept for compatibility and future feature toggling.
	UseStub bool
}

// Codec is a minimal EXI codec instance.
type Codec struct {
	cfg *Config
}

// NewCodec creates a new Codec instance using the provided Config.
func NewCodec(cfg *Config) *Codec {
	if cfg == nil {
		cfg = &Config{UseStub: true}
	}
	return &Codec{cfg: cfg}
}

// Init prepares the codec for use. No-op for this simple implementation.
func (c *Codec) Init() error {
	if c == nil {
		return errors.New("exi: codec is nil")
	}
	// real implementation: parse grammars here
	return nil
}

// Shutdown releases resources (no-op here).
func (c *Codec) Shutdown() error {
	return nil
}

// minifyXML removes unnecessary whitespace between XML tags and trims the input.
// It's a conservative minifier that preserves textual content while removing
// inter-tag and leading/trailing whitespace.
func minifyXML(in []byte) []byte {
	s := string(in)
	// collapse sequences of whitespace between tags: >   <  ->  ><
	re := regexp.MustCompile(`>\s+<`)
	s = re.ReplaceAllString(s, "><")
	// trim
	return []byte(s)
}

// gzipCompress compresses bytes using gzip and returns the result.
func gzipCompress(in []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	// write raw bytes
	if _, err := zw.Write(in); err != nil {
		_ = zw.Close()
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// gzipDecompress decompresses gzip data into raw bytes.
func gzipDecompress(in []byte) ([]byte, error) {
	r := bytes.NewReader(in)
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	out, err := io.ReadAll(zr)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EncodeXML encodes an XML document (as bytes) into a gzipped, minified
// binary representation. This provides a compact, deterministic payload
// for testing. It also validates that the input is well-formed XML.
func (c *Codec) EncodeXML(xmlBytes []byte) ([]byte, error) {
	if c == nil {
		return nil, errors.New("exi: codec is nil")
	}
	// Ensure input is valid XML to catch user errors early.
	var tmp any
	if err := xml.Unmarshal(xmlBytes, &tmp); err != nil {
		return nil, err
	}
	// Minify the XML
	min := minifyXML(xmlBytes)
	// Compress
	out, err := gzipCompress(min)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DecodeEXI decodes bytes produced by EncodeXML: it attempts to gunzip and
// returns the minified XML bytes. The caller can unmarshal or pretty-print
// as needed.
func (c *Codec) DecodeEXI(exiBytes []byte) ([]byte, error) {
	if c == nil {
		return nil, errors.New("exi: codec is nil")
	}
	// Try to decompress gzip payload
	xmlOut, err := gzipDecompress(exiBytes)
	if err == nil {
		// Optionally validate XML
		var tmp any
		if err := xml.Unmarshal(xmlOut, &tmp); err != nil {
			return nil, err
		}
		return xmlOut, nil
	}
	// If not gzip, consider it invalid for this codec
	return nil, ErrInvalidEXI
}

// Convenience wrappers -----------------------------------------------------

// EncodeXMLToEXI is a helper that creates a temporary codec, initializes it,
// encodes the XML and shuts the codec down.
func EncodeXMLToEXI(xmlBytes []byte) ([]byte, error) {
	c := NewCodec(nil)
	if err := c.Init(); err != nil {
		return nil, err
	}
	defer c.Shutdown()
	return c.EncodeXML(xmlBytes)
}

// DecodeEXIToXML is a helper to decode exi bytes into XML bytes using a
// temporary codec.
func DecodeEXIToXML(exiBytes []byte) ([]byte, error) {
	c := NewCodec(nil)
	if err := c.Init(); err != nil {
		return nil, err
	}
	defer c.Shutdown()
	return c.DecodeEXI(exiBytes)
}
