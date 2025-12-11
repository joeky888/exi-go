package schema

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ParseSchemas parses loaded XSD Schema structures and produces a set of
// TypeDefinition values representing the schema's types.
//
// Implemented features (minimal, pragmatic):
//   - top-level <xs:complexType name="..."> with an inner <xs:sequence>
//     containing <xs:element name="..." type="..."> children -> produces
//     a TypeDefinition with Fields for those elements.
//   - top-level <xs:simpleType name="..."> with a <xs:restriction> that
//     contains <xs:enumeration value="..."> children -> produces an enum
//     TypeDefinition with EnumValues.
//   - top-level <xs:element name="..." type="..."> declarations are recorded
//     as TypeDefinitions with Docs pointing to the referenced type (if any).
//   - inline (anonymous) complexType under an element is handled by creating
//     a TypeDefinition named <ElementName>Type and associating it with the
//     element name.
//
// Limitations:
//   - does not support all XSD constructs (choices, attributes, substitution
//     groups, complexContent/extension, unions, lists, etc.).
//   - only handles sequence->element children for complex types.
//   - referenced type names are recorded as-is; no attempt is made to
//     resolve namespace prefixes to schema files here.
//
// This implementation is intentionally small so the rest of the generator can
// iterate on top of a working XSD parser. Extend as needed.
func ParseSchemas(schemas []*Schema) ([]*TypeDefinition, error) {
	if len(schemas) == 0 {
		return nil, fmt.Errorf("no schemas to parse")
	}

	// map[typeName] -> TypeDefinition to avoid duplicates
	typeMap := map[string]*TypeDefinition{}

	for _, s := range schemas {
		if s == nil || len(s.Raw) == 0 {
			continue
		}
		decoder := xml.NewDecoder(bytes.NewReader(s.Raw))
		// Process tokens
		for {
			tok, err := decoder.Token()
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, fmt.Errorf("xml parse error in %s: %w", s.Path, err)
			}
			switch se := tok.(type) {
			case xml.StartElement:
				ln := localName(se.Name)
				switch ln {
				case "complexType":
					name := getAttr(se.Attr, "name")
					if name == "" {
						// anonymous complexType - caller likely handles inline under element
						// skip here
						// consume until matching end to avoid streaming confusion
						if err := skipElement(decoder, se); err != nil {
							return nil, err
						}
						continue
					}
					td, err := parseComplexType(decoder, se, s)
					if err != nil {
						return nil, err
					}
					// prefer existing if present, but merge minimally
					if existing, ok := typeMap[name]; ok {
						// merge fields if existing has none
						if len(existing.Fields) == 0 && len(td.Fields) > 0 {
							existing.Fields = td.Fields
						}
						if existing.Docs == "" {
							existing.Docs = td.Docs
						}
					} else {
						typeMap[name] = td
					}
				case "simpleType":
					name := getAttr(se.Attr, "name")
					td, err := parseSimpleType(decoder, se, s)
					if err != nil {
						return nil, err
					}
					if name == "" {
						// anonymous simpleType - skip
						continue
					}
					td.Name = name
					if existing, ok := typeMap[name]; ok {
						if !existing.IsEnum && td.IsEnum {
							existing.IsEnum = true
							existing.EnumValues = td.EnumValues
						}
					} else {
						typeMap[name] = td
					}
				case "element":
					// top-level element declaration
					name := getAttr(se.Attr, "name")
					typ := getAttr(se.Attr, "type")
					if name == "" {
						// skip unnamed
						if err := skipElement(decoder, se); err != nil {
							return nil, err
						}
						continue
					}
					// if element has an inline complexType, parse it and create a type
					inlineTd, hasInline, err := parseElementMaybeInlineComplex(decoder, se, s, name)
					if err != nil {
						return nil, err
					}
					// If inline type created, add it and also add an element-level TypeDefinition
					if hasInline && inlineTd != nil {
						typeMap[inlineTd.Name] = inlineTd
						// Also create an element wrapper type referencing inline type
						elemTd := &TypeDefinition{
							Name:      name,
							Namespace: s.Namespace,
							Docs:      fmt.Sprintf("element -> inline complexType %s", inlineTd.Name),
						}
						typeMap[name] = elemTd
						continue
					}
					// No inline type: create a TypeDefinition that references typ (if any)
					td := &TypeDefinition{
						Name:      name,
						Namespace: s.Namespace,
						Docs:      fmt.Sprintf("element references type %s", typ),
					}
					typeMap[name] = td
				default:
					// ignore other start elements
				}
			default:
				// skip other tokens
			}
		}
	}

	// convert map to slice
	var defs []*TypeDefinition
	for _, v := range typeMap {
		defs = append(defs, v)
	}
	return defs, nil
}

// parseComplexType parses a named complexType start element. It expects the
// StartElement for <complexType name="..."> to have already been read.
// It returns a TypeDefinition that may contain Fields parsed from an inner
// <sequence> of <element> children.
func parseComplexType(decoder *xml.Decoder, start xml.StartElement, s *Schema) (*TypeDefinition, error) {
	name := getAttr(start.Attr, "name")
	td := &TypeDefinition{
		Name:      name,
		Namespace: s.Namespace,
		Fields:    nil,
	}

	// Depth-first scan until matching end for complexType
	for {
		tok, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			ln := localName(t.Name)
			if ln == "sequence" {
				// parse elements inside sequence
				fields, err := parseSequenceElements(decoder, t)
				if err != nil {
					return nil, err
				}
				td.Fields = append(td.Fields, fields...)
			} else if ln == "attribute" {
				// parse attribute as field (optional)
				attrName := getAttr(t.Attr, "name")
				attrType := getAttr(t.Attr, "type")
				if attrName != "" {
					f := Field{
						Name:       attrName,
						Type:       attrType,
						IsOptional: true,
						IsArray:    false,
					}
					td.Fields = append(td.Fields, f)
				}
				// consume attribute element end
				if err := skipElement(decoder, t); err != nil {
					return nil, err
				}
			} else {
				// skip other start elements inside complexType
				if err := skipElement(decoder, t); err != nil {
					return nil, err
				}
			}
		case xml.EndElement:
			if localName(t.Name) == "complexType" {
				return td, nil
			}
		default:
			// continue scanning
		}
	}
}

// parseSequenceElements parses a <sequence> start element and returns a slice
// of Fields corresponding to <element> children. The sequence StartElement is
// already consumed.
func parseSequenceElements(decoder *xml.Decoder, start xml.StartElement) ([]Field, error) {
	var fields []Field
	for {
		tok, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if localName(t.Name) == "element" {
				name := getAttr(t.Attr, "name")
				typ := getAttr(t.Attr, "type")
				minOccurs := getAttr(t.Attr, "minOccurs")
				maxOccurs := getAttr(t.Attr, "maxOccurs")
				isOptional := false
				isArray := false
				if minOccurs != "" {
					if mv, err := strconv.Atoi(minOccurs); err == nil && mv == 0 {
						isOptional = true
					}
				}
				if maxOccurs != "" {
					if maxOccurs == "unbounded" {
						isArray = true
					} else if mv, err := strconv.Atoi(maxOccurs); err == nil && mv > 1 {
						isArray = true
					}
				}
				// consume any nested inline complexType or simpleType for this element
				// but do not attempt to fully parse it here (higher-level parser will)
				if hasInlineChild(decoder, t) {
					// hasInlineChild consumed nested content already
				}
				if name == "" {
					// element without name - skip / ignore
					continue
				}
				fields = append(fields, Field{
					Name:       name,
					Type:       typ,
					IsOptional: isOptional,
					IsArray:    isArray,
				})
			} else {
				// skip nested start elements we don't handle
				if err := skipElement(decoder, t); err != nil {
					return nil, err
				}
			}
		case xml.EndElement:
			if localName(t.Name) == "sequence" {
				return fields, nil
			}
		}
	}
}

// parseSimpleType parses a simpleType start element. It looks for a restriction
// with enumeration children to produce an enum TypeDefinition. The caller is
// responsible for assigning the TypeDefinition.Name (if available).
func parseSimpleType(decoder *xml.Decoder, start xml.StartElement, s *Schema) (*TypeDefinition, error) {
	td := &TypeDefinition{
		Namespace: s.Namespace,
	}
	var enumVals []string
	for {
		tok, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if localName(t.Name) == "restriction" {
				// parse enumeration children until end of restriction
				for {
					rtok, err := decoder.Token()
					if err != nil {
						return nil, err
					}
					switch r := rtok.(type) {
					case xml.StartElement:
						if localName(r.Name) == "enumeration" {
							val := getAttr(r.Attr, "value")
							if val != "" {
								enumVals = append(enumVals, val)
							}
							// consume until end of enumeration
							if err := skipElement(decoder, r); err != nil {
								return nil, err
							}
						} else {
							// skip other nested nodes
							if err := skipElement(decoder, r); err != nil {
								return nil, err
							}
						}
					case xml.EndElement:
						if localName(r.Name) == "restriction" {
							break
						}
					}
					// continue until restriction end is found in the outer loop
				}
			} else {
				// skip other nested nodes
				if err := skipElement(decoder, t); err != nil {
					return nil, err
				}
			}
		case xml.EndElement:
			if localName(t.Name) == "simpleType" {
				td.IsEnum = len(enumVals) > 0
				td.EnumValues = enumVals
				return td, nil
			}
		}
	}
}

// parseElementMaybeInlineComplex handles an <element ...> StartElement and
// determines whether it contains an inline complexType. If so, it parses it
// and returns a created TypeDefinition with a generated name "<ElementName>Type".
// It returns (td, true, nil) when an inline complexType was parsed and td!=nil.
// If there is no inline type it returns (nil, false, nil).
//
// Note: the provided start element is the element start token that was seen by
// the caller. This function will consume until the element end.
func parseElementMaybeInlineComplex(decoder *xml.Decoder, start xml.StartElement, s *Schema, elemName string) (*TypeDefinition, bool, error) {
	// Look ahead for an inline complexType by reading tokens up to the end of the element.
	depth := 1
	var inlineStarted xml.StartElement
	for {
		tok, err := decoder.Token()
		if err != nil {
			return nil, false, err
		}
		switch tk := tok.(type) {
		case xml.StartElement:
			if localName(tk.Name) == "complexType" {
				// parse the inline complexType into a generated type named elemName + "Type"
				inlName := elemName + "Type"
				inlineStarted = tk
				td, err := parseComplexType(decoder, inlineStarted, s)
				if err != nil {
					return nil, false, err
				}
				// Ensure the generated type has an explicit name
				td.Name = inlName
				return td, true, nil
			} else {
				// skip nested content we don't care about
				if err := skipElement(decoder, tk); err != nil {
					return nil, false, err
				}
			}
		case xml.EndElement:
			if localName(tk.Name) == "element" {
				// reached end of element and no inline complexType found
				return nil, false, nil
			}
		}
		_ = depth
	}
}

// getAttr returns the value of the named attribute in attrs or empty string.
func getAttr(attrs []xml.Attr, name string) string {
	for _, a := range attrs {
		if a.Name.Local == name {
			return strings.TrimSpace(a.Value)
		}
	}
	return ""
}

// localName returns the local (unprefixed) name of an xml.Name.
func localName(n xml.Name) string {
	return n.Local
}

// skipElement consumes tokens from decoder until the matching EndElement for
// the provided start element is found.
func skipElement(decoder *xml.Decoder, start xml.StartElement) error {
	depth := 1
	for depth > 0 {
		tok, err := decoder.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			depth++
		case xml.EndElement:
			depth--
		default:
			// ignore other tokens
			_ = t
		}
	}
	return nil
}

// hasInlineChild checks whether the start element has an inline child (like
// an inline complexType or simpleType). If such child exists, it consumes
// that nested element content (so callers don't need to re-skip it).
//
// This helper reads one token to detect the nested start element and delegates
// skipElement to consume any nested content; it returns true if a nested start
// element was found and consumed.
func hasInlineChild(decoder *xml.Decoder, start xml.StartElement) bool {
	// Peek by reading tokens until we see either a StartElement or EndElement
	for {
		tok, err := decoder.Token()
		if err != nil {
			return false
		}
		switch t := tok.(type) {
		case xml.StartElement:
			// consume nested element entirely
			_ = skipElement(decoder, t)
			return true
		case xml.EndElement:
			// no inline children
			return false
		default:
			// continue scanning (text, comments, etc.)
		}
	}
}
