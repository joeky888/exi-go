package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"example.com/exi-go/pkg/schema"
)

// This small tool regenerates Go types from XSD schemas using the in-repo
// CodeGenerator. It is intended to be run from the go module directory
// (exi-go) during development to refresh generated types.
//
// Usage examples:
//
//	# regenerate using the community schema folders (default)
//	go run tools_regen.go
//
//	# regenerate using explicit schema directories and output directory
//	go run tools_regen.go -schemas ./schemas/switchev,./schemas/tux-evse -out ./pkg/v2g/generated -pkg generated -author "exi-go"
//
// Notes:
//   - The tool will load all .xsd files found in the provided schema paths
//     (each schema path may be a file or a directory).
//   - If -force is set the output directory will be removed before generation.
func main() {
	var (
		schemaPathsFlag string
		outDir          string
		pkgName         string
		author          string
		force           bool
	)

	flag.StringVar(&schemaPathsFlag, "schemas", "./schemas/switchev", "Comma-separated list of schema files or directories to use")
	flag.StringVar(&outDir, "out", "./pkg/v2g/generated", "Output directory for generated Go code")
	flag.StringVar(&pkgName, "pkg", "generated", "Package name for generated code")
	flag.StringVar(&author, "author", "exi-go", "Author string to embed in generated headers")
	flag.BoolVar(&force, "force", false, "Remove output directory before generation")
	flag.Parse()

	// Expand the schema path list
	var rawPaths []string
	for _, p := range strings.Split(schemaPathsFlag, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		rawPaths = append(rawPaths, p)
	}

	if len(rawPaths) == 0 {
		log.Fatalf("no schema paths provided; use -schemas flag")
	}

	// Validate and expand directories to XSD file paths using the helper.
	var allSchemaFiles []string
	for _, rp := range rawPaths {
		// If it's a directory, add the directory path (LoadSchemas will discover .xsd files)
		if info, err := os.Stat(rp); err == nil && info.IsDir() {
			allSchemaFiles = append(allSchemaFiles, rp)
			continue
		}
		// Otherwise, accept the path as a file (will be validated by LoadSchemas)
		allSchemaFiles = append(allSchemaFiles, rp)
	}

	absOut, err := filepath.Abs(outDir)
	if err != nil {
		log.Fatalf("failed to resolve out dir: %v", err)
	}

	// Handle force removal of output directory if requested
	if force {
		if err := os.RemoveAll(absOut); err != nil {
			log.Fatalf("failed to remove output dir %s: %v", absOut, err)
		}
	}

	// Ensure parent directories exist
	if err := os.MkdirAll(absOut, 0o755); err != nil {
		log.Fatalf("failed to create output dir %s: %v", absOut, err)
	}

	// Load schemas (the loader accepts file or directory paths)
	schemas, err := schema.LoadSchemas(allSchemaFiles)
	if err != nil {
		log.Fatalf("failed to load schemas: %v", err)
	}
	if len(schemas) == 0 {
		log.Fatalf("no schemas found in provided paths: %v", allSchemaFiles)
	}

	// Use the CodeGenerator defined in the schema package
	gen := &schema.CodeGenerator{
		PackageName: pkgName,
		Author:      author,
	}

	log.Printf("generating Go types into %s (package %s) from %d schema files", absOut, pkgName, len(schemas))
	if err := gen.GenerateFromSchemas(schemas, absOut); err != nil {
		log.Fatalf("generation failed: %v", err)
	}

	fmt.Printf("Generation completed successfully. Output written to: %s\n", absOut)
}
