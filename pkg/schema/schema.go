// Package schema provides helper types and placeholder implementation for
// parsing XML Schema (XSD) files and driving a code generator that produces
// Go types and EXI codec stubs.
//
// This file is an initial scaffold for the generator approach (Option A).
// It intentionally implements safe, minimal behavior so the rest of the
// system can be integrated and tested while the full XSD parsing and
// code-generation logic is implemented iteratively.
//
// Goals and responsibilities of this package (future):
//   - discover and load XSD files (individual files or directories)
//   - parse relevant XSD constructs needed for ISO 15118 (types, elements,
//     namespaces, substitution groups, enums)
//   - provide an intermediate representation (IR) of types and elements
//     suitable for code generation
//   - offer a stable Generator interface that accepts the IR and emits
//     Go code and glue for the EXI runtime
//
// NOTE on ISO XSDs and redistribution:
//   - Some ISO 15118 schema files are subject to ISO distribution terms.
//     Do not embed or redistribute ISO-proprietary XSDs unless you have the
//     right to do so. The generator should accept external XSD paths or
//     optionally perform downloads when explicitly requested by the user.
package schema

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ErrNotImplemented is returned by functions that are intentionally left
// as placeholders in this scaffold.
var ErrNotImplemented = errors.New("schema: not implemented")

// Schema represents a loaded XSD file and minimal metadata extracted from it.
type Schema struct {
	// Path is the original filesystem path to the XSD file.
	Path string

	// Namespace is the target namespace declared by the XSD (may be empty).
	Namespace string

	// Raw contains the bytes of the XSD file as read from disk. Parsers
	// can consume these bytes to build the IR.
	Raw []byte
}

// Field describes a single field/member inside a type definition.
type Field struct {
	Name       string // element or attribute name
	Type       string // raw or resolved type name
	IsOptional bool   // true if minOccurs=0 or optional attribute
	IsArray    bool   // true if maxOccurs>1 or unbounded
}

// TypeDefinition is an intermediate representation of a schema type
// (simple or complex) that the code generator will use as input.
type TypeDefinition struct {
	Name       string  // e.g. "SessionSetupReqType"
	Namespace  string  // originating namespace
	Fields     []Field // ordered fields for the type
	IsEnum     bool    // whether this type is an enumeration
	EnumValues []string
	Docs       string // captured documentation/comments (if any)
}

// Generator is the interface implemented by components that turn the
// schema IR into concrete Go source files and EXI codec stubs.
type Generator interface {
	// GenerateFromSchemas receives a list of loaded Schema structures and
	// writes Go source files into outDir. Implementations should create
	// files under outDir and return an error on failure.
	GenerateFromSchemas(schemas []*Schema, outDir string) error
}

// ValidateSchemaPaths accepts a list of paths (files or directories) and
// expands them to a list of candidate XSD file paths. It returns the
// discovered XSD file paths or an error.
func ValidateSchemaPaths(paths []string) ([]string, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("no schema paths provided")
	}
	var result []string
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		info, err := os.Stat(p)
		if err != nil {
			return nil, fmt.Errorf("stat %s: %w", p, err)
		}
		if info.IsDir() {
			// find .xsd files in this directory (non-recursive)
			entries, err := os.ReadDir(p)
			if err != nil {
				return nil, fmt.Errorf("read dir %s: %w", p, err)
			}
			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				name := e.Name()
				if strings.HasSuffix(strings.ToLower(name), ".xsd") {
					result = append(result, filepath.Join(p, name))
				}
			}
		} else {
			// single file - ensure it has .xsd extension
			if strings.HasSuffix(strings.ToLower(p), ".xsd") {
				result = append(result, p)
			} else {
				return nil, fmt.Errorf("path %s is not an XSD file", p)
			}
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no XSD files found in provided paths")
	}
	return result, nil
}

// LoadSchemas reads XSD files from the provided paths (files or directories)
// and returns a slice of Schema instances with file contents loaded.
// This is a thin helper: the full parsing of XSD to IR is done by ParseSchemas.
func LoadSchemas(paths []string) ([]*Schema, error) {
	expanded, err := ValidateSchemaPaths(paths)
	if err != nil {
		return nil, err
	}
	var schemas []*Schema
	for _, p := range expanded {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", p, err)
		}
		// Minimal namespace sniff: look for targetNamespace="..." in the raw bytes.
		// This is a convenience; the real parser should derive namespace properly.
		ns := sniffNamespace(string(data))
		schemas = append(schemas, &Schema{
			Path:      p,
			Namespace: ns,
			Raw:       data,
		})
	}
	return schemas, nil
}

// sniffNamespace performs a very small heuristic extraction of the
// targetNamespace attribute from the XSD contents. This is intentionally
// simplistic and only used for metadata in the scaffold.
func sniffNamespace(content string) string {
	lc := strings.ToLower(content)
	idx := strings.Index(lc, "targetnamespace")
	if idx == -1 {
		return ""
	}
	// find first '=' after the keyword
	rest := lc[idx:]
	eq := strings.Index(rest, "=")
	if eq == -1 {
		return ""
	}
	rest = rest[eq+1:]
	rest = strings.TrimLeft(rest, " \t\r\n")
	if len(rest) == 0 {
		return ""
	}
	quote := rest[0]
	rest = rest[1:]
	qidx := strings.IndexByte(rest, quote)
	if qidx == -1 {
		return ""
	}
	return rest[:qidx]
}

// ParseSchemas is implemented in parse_xsd.go.
//
// The full, incremental XSD parser was moved into a dedicated file
// (`parse_xsd.go`) to keep the implementation focused and maintainable.
// That implementation provides a pragmatic parser for common XSD constructs
// used by the generator (named complexType with sequence children,
// simpleType enumerations, top-level elements and basic inline complexTypes).
//
// Callers should use ParseSchemas as declared in this package; the concrete
// implementation is provided in `parse_xsd.go`.

// DownloadSchemas downloads the provided URLs and writes them into destDir.
// It will create destDir if necessary. Each URL is saved under its base filename.
// The caller must ensure that it has the right to download and redistribute
// the requested files; for ISO schemas this typically requires acceptance of
// ISO's license/terms. This function performs a simple HTTP GET for each URL.
func DownloadSchemas(urls []string, destDir string, client *http.Client) ([]string, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("no URLs provided")
	}
	if destDir == "" {
		return nil, fmt.Errorf("destDir must be provided")
	}
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return nil, fmt.Errorf("create dest dir: %w", err)
	}
	var saved []string
	for _, u := range urls {
		// Basic validation
		if strings.TrimSpace(u) == "" {
			continue
		}
		resp, err := client.Get(u)
		if err != nil {
			return nil, fmt.Errorf("download %s: %w", u, err)
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			return nil, fmt.Errorf("download %s: unexpected status %s", u, resp.Status)
		}
		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("read %s body: %w", u, err)
		}
		// derive filename from URL path
		base := filepath.Base(resp.Request.URL.Path)
		if base == "" || base == "/" {
			// fallback filename
			base = "schema.xsd"
		}
		dst := filepath.Join(destDir, base)
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return nil, fmt.Errorf("write %s: %w", dst, err)
		}
		saved = append(saved, dst)
	}
	return saved, nil
}

// DownloadISO15118Schemas is a convenience helper that downloads a set of
// known public ISO 15118 XSDs into destDir. Because ISO schemas are often
// distributed under ISO terms, the caller MUST set acceptISO to true to
// indicate explicit acceptance of ISO's usage terms before the download
// proceeds. The returned slice contains local file paths of the downloaded
// files.
//
// Note: The list below is intentionally empty as a placeholder. Replace the
// entries with authoritative URLs from a permitted source (or provide a list
// of local paths). Do NOT include proprietary ISO-distributed files here.
func DownloadISO15118Schemas(destDir string, acceptISO bool) ([]string, error) {
	if !acceptISO {
		return nil, fmt.Errorf("must accept ISO schema terms to download ISO 15118 schemas")
	}
	// Example placeholder URLs. Replace or extend with real, permitted URLs.
	urls := []string{
		// "https://example.org/iso15118/ISO15118-20.xsd",
		// "https://example.org/iso15118/ISO15118-20-Common.xsd",
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("no public ISO schema URLs configured; please provide URLs or use local schema files")
	}
	return DownloadSchemas(urls, destDir, nil)
}

// GenerateTypes is a high-level convenience function that runs the common
// flow: discover/load schemas, parse them into the IR and run the supplied
// Generator to write code to disk. The function performs basic validation
// and writes a minimal placeholder if the generator is nil.
func GenerateTypes(schemaPaths []string, outDir string, gen Generator) error {
	// Load schemas from provided paths (files or directories)
	schemas, err := LoadSchemas(schemaPaths)
	if err != nil {
		return fmt.Errorf("load schemas: %w", err)
	}

	// Parse into intermediate representation (IR)
	types, err := ParseSchemas(schemas)
	if err != nil {
		// During early development it's acceptable to let consumers proceed
		// with a placeholder generated module. Return the not implemented
		// error so callers may decide how to proceed.
		return fmt.Errorf("parse schemas: %w", err)
	}
	_ = types // currently unused - kept to show intended flow

	// If no generator provided, produce a placeholder generated file to show
	// the output flow. Otherwise, call the generator.
	if gen == nil {
		// Ensure output dir exists
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			return fmt.Errorf("create outdir: %w", err)
		}
		placeholder := filepath.Join(outDir, "zz_generated_placeholder.go")
		content := []byte(`// Code generated by exi-go (placeholder).
// This file is created by GenerateTypes when no Generator is supplied.
// Replace with a real generator implementation.
package generated
`)
		if err := os.WriteFile(placeholder, content, 0o644); err != nil {
			return fmt.Errorf("write placeholder: %w", err)
		}
		return nil
	}

	// Call the user-supplied generator to produce files.
	if err := gen.GenerateFromSchemas(schemas, outDir); err != nil {
		return fmt.Errorf("generator failed: %w", err)
	}
	return nil
}

// PackageReadme is a short README describing the responsibilities of this
// package and how to extend it.
const PackageReadme = `schema package - README

This package is a scaffold that will evolve into the XSD parsing and
code generation core of the exi-go project.

How to extend:
  1. Implement ParseSchemas to build a reliable IR (TypeDefinition) from
     Schema.Raw bytes. Start small: parse elements, simple/complex
     types, sequences and attributes.
  2. Create a Generator implementation that consumes the IR and emits
     Go types and EXI grammar glue. Keep code generation templates simple
     and test roundtrips for a small set of ISO 15118 message samples.
  3. Add unit tests and sample XSDs under an internal testdata folder and
     iterate until the generated code covers the necessary message set.

Licensing and schemas:
  - Do not include ISO-distributed XSD files in the repository unless you
    have permission to redistribute them. The generator should accept
    external XSD paths or offer an explicit download mechanism that
    clearly documents ISO's usage terms.
`
