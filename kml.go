//go:generate go run ./internal/generate -f -o kml22gx.gen.go -n gx: xsd/kml22gx.xsd
//go:generate go run ./internal/generate -f -o ogckml22.gen.go xsd/ogckml22.xsd

// Package kml provides convenience methods for creating and writing KML documents.
//
// See https://developers.google.com/kml/
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
	"image/color"
	"io"
	"strconv"
	"time"
)

// An Element represents an abstract KML element.
type Element interface {
	xml.Marshaler
	Write(io.Writer) error
	WriteIndent(io.Writer, string, string) error
}

// A SimpleElement is an Element with a single value.
type SimpleElement struct {
	xml.StartElement
	value string
}

// A CompoundElement is an Element with children.
type CompoundElement struct {
	xml.StartElement
	children []Element
}

// A SharedElement is an element with an id.
type SharedElement struct {
	CompoundElement
	id string
}

// MarshalXML marshals se to e. start is ignored.
func (se *SimpleElement) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(xml.CharData(se.value), se.StartElement)
}

// Write writes an XML header and se to w.
func (se *SimpleElement) Write(w io.Writer) error {
	return write(w, "", "", se)
}

// WriteIndent writes an XML header and se to w.
func (se *SimpleElement) WriteIndent(w io.Writer, prefix, indent string) error {
	return write(w, prefix, indent, se)
}

// Add adds children to ce.
func (ce *CompoundElement) Add(children ...Element) *CompoundElement {
	ce.children = append(ce.children, children...)
	return ce
}

// MarshalXML marshals ce to e. start is ignored.
func (ce *CompoundElement) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(ce.StartElement); err != nil {
		return err
	}
	for _, c := range ce.children {
		if err := e.EncodeElement(c, ce.StartElement); err != nil {
			return err
		}
	}
	return e.EncodeToken(ce.End())
}

// Write writes an XML header and ce to w.
func (ce *CompoundElement) Write(w io.Writer) error {
	return write(w, "", "", ce)
}

// WriteIndent writes an XML header and ce to w.
func (ce *CompoundElement) WriteIndent(w io.Writer, prefix, indent string) error {
	return write(w, prefix, indent, ce)
}

// ID returns se's id.
func (se *SharedElement) ID() string {
	return se.id
}

// URL returns se's URL.
func (se *SharedElement) URL() string {
	return "#" + se.ID()
}

func write(w io.Writer, prefix, indent string, m xml.Marshaler) error {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	e := xml.NewEncoder(w)
	e.Indent(prefix, indent)
	return e.Encode(m)
}

func newSEBool(name string, value bool) *SimpleElement {
	var v string
	if value {
		v = "1"
	} else {
		v = "0"
	}
	return &SimpleElement{
		StartElement: xml.StartElement{Name: xml.Name{Local: name}},
		value:        v,
	}
}

func newSEColor(name string, value color.Color) *SimpleElement {
	r, g, b, a := value.RGBA()
	return &SimpleElement{
		StartElement: xml.StartElement{Name: xml.Name{Local: name}},
		value:        fmt.Sprintf("%02x%02x%02x%02x", a/256, b/256, g/256, r/256),
	}
}

func newSEElement(name string, value Element) *CompoundElement {
	return &CompoundElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: name},
		},
		children: []Element{value},
	}
}

func newSEFloat(name string, value float64) *SimpleElement {
	return &SimpleElement{
		StartElement: xml.StartElement{Name: xml.Name{Local: name}},
		value:        strconv.FormatFloat(value, 'f', -1, 64),
	}
}

func newSEInt(name string, value int) *SimpleElement {
	return &SimpleElement{
		StartElement: xml.StartElement{Name: xml.Name{Local: name}},
		value:        strconv.Itoa(value),
	}
}

func newSEVec2(name string, value Vec2) *SimpleElement {
	return &SimpleElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: name},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "x"}, Value: strconv.FormatFloat(value.X, 'f', -1, 64)},
				{Name: xml.Name{Local: "y"}, Value: strconv.FormatFloat(value.Y, 'f', -1, 64)},
				{Name: xml.Name{Local: "xunits"}, Value: string(value.XUnits)},
				{Name: xml.Name{Local: "yunits"}, Value: string(value.YUnits)},
			},
		},
	}
}

func newSEString(name, value string) *SimpleElement {
	return &SimpleElement{
		StartElement: xml.StartElement{Name: xml.Name{Local: name}},
		value:        value,
	}
}

func newSETime(name string, value time.Time) *SimpleElement {
	return &SimpleElement{
		StartElement: xml.StartElement{Name: xml.Name{Local: name}},
		value:        value.Format(time.RFC3339),
	}
}

func newCE(name string, children []Element) *CompoundElement {
	return &CompoundElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: name},
		},
		children: children,
	}
}

func newSharedE(name, id string, children []Element) *SharedElement {
	var attr []xml.Attr
	if id != "" {
		attr = append(attr, xml.Attr{Name: xml.Name{Local: "id"}, Value: id})
	}
	return &SharedElement{
		CompoundElement: CompoundElement{
			StartElement: xml.StartElement{
				Name: xml.Name{Local: name},
				Attr: attr,
			},
			children: children,
		},
		id: id,
	}
}
