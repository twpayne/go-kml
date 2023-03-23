//go:generate go run ./internal/generate -o kml22gx.gen.go -n gx: xsd/kml22gx.xsd
//go:generate go run ./internal/generate -o ogckml22.gen.go xsd/ogckml22.xsd

// Package kml provides convenience methods for creating and writing KML documents.
//
// See https://developers.google.com/kml/.
//
// Goals
//
//   - Convenient API for creating both simple and complex KML documents.
//   - 1:1 mapping between functions and KML elements.
//
// Non-goals
//
//   - Protection against generating invalid documents.
//   - Concealment of KML complexity.
//   - Fine-grained control over generated XML.
package kml

import (
	"encoding/xml"
	"fmt"
	"io"
)

const float64StringSize = 16

type unsupportedTypeError struct {
	value any
}

func (e unsupportedTypeError) Error() string {
	return fmt.Sprintf("%T: unsupported type", e.value)
}

// An Element is a KML element.
type Element interface {
	xml.Marshaler
}

// A TopLevelElement is a top level KML element.
type TopLevelElement interface {
	Element
	Write(w io.Writer) error
	WriteIndent(w io.Writer, prefix, indent string) error
}

func encodeElement(encoder *xml.Encoder, startElement xml.StartElement) error {
	if err := encoder.EncodeToken(startElement); err != nil {
		return err
	}
	return encoder.EncodeToken(startElement.End())
}

func encodeElementWithCharData(encoder *xml.Encoder, startElement xml.StartElement, charData xml.CharData) error {
	if err := encoder.EncodeToken(startElement); err != nil {
		return err
	}
	if charData != nil {
		if err := encoder.EncodeToken(charData); err != nil {
			return err
		}
	}
	return encoder.EncodeToken(startElement.End())
}

func encodeElementWithChild(encoder *xml.Encoder, startElement xml.StartElement, child Element) error {
	if err := encoder.EncodeToken(startElement); err != nil {
		return err
	}
	if child != nil {
		if err := child.MarshalXML(encoder, xml.StartElement{}); err != nil {
			return err
		}
	}
	return encoder.EncodeToken(startElement.End())
}

func encodeElementWithChildren(encoder *xml.Encoder, startElement xml.StartElement, children []Element) error {
	if err := encoder.EncodeToken(startElement); err != nil {
		return err
	}
	for _, child := range children {
		if child == nil {
			continue
		}
		if err := child.MarshalXML(encoder, xml.StartElement{}); err != nil {
			return err
		}
	}
	return encoder.EncodeToken(startElement.End())
}

func write(w io.Writer, e Element) error {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	return xml.NewEncoder(w).Encode(e)
}

func writeIndent(w io.Writer, e Element, prefix, indent string) error {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent(prefix, indent)
	return encoder.Encode(e)
}
