package kml

import (
	"encoding/xml"
	"io"
	"strconv"
	"strings"
)

// GxNamespace is the default namespace for Google Earth extensions.
const GxNamespace = "http://www.google.com/kml/ext/2.2"

// A GxOptionName is a gx:option name.
type GxOptionName string

func (e GxOptionName) String() string { return string(e) }

// GxOptionNames.
const (
	GxOptionNameHistoricalImagery GxOptionName = "historicalimagery"
	GxOptionNameStreetView        GxOptionName = "streetview"
	GxOptionNameSunlight          GxOptionName = "sunlight"
)

// A GxAnglesElement is a gx:angles element.
type GxAnglesElement struct {
	Heading float64
	Tilt    float64
	Roll    float64
}

// GxAngles returns a new GxAnglesElement.
func GxAngles(heading, tilt, roll float64) *GxAnglesElement {
	return &GxAnglesElement{
		Heading: heading,
		Tilt:    tilt,
		Roll:    roll,
	}
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *GxAnglesElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{Name: xml.Name{Local: "gx:angles"}}
	var builder strings.Builder
	builder.Grow(3 * float64StringSize)
	builder.WriteString(strconv.FormatFloat(e.Heading, 'f', -1, 64))
	builder.WriteByte(' ')
	builder.WriteString(strconv.FormatFloat(e.Tilt, 'f', -1, 64))
	builder.WriteByte(' ')
	builder.WriteString(strconv.FormatFloat(e.Roll, 'f', -1, 64))
	charData := xml.CharData(builder.String())
	return encodeElementWithCharData(encoder, startElement, charData)
}

// A GxCoordElement is a gx:coord element.
type GxCoordElement Coordinate

// GxCoord returns a new GxCoordElement.
func GxCoord(coordinate Coordinate) GxCoordElement {
	return GxCoordElement(coordinate)
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e GxCoordElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{Name: xml.Name{Local: "gx:coord"}}
	var builder strings.Builder
	builder.Grow(3 * float64StringSize)
	builder.WriteString(strconv.FormatFloat(e.Lon, 'f', -1, 64))
	builder.WriteByte(' ')
	builder.WriteString(strconv.FormatFloat(e.Lat, 'f', -1, 64))
	builder.WriteByte(' ')
	builder.WriteString(strconv.FormatFloat(e.Alt, 'f', -1, 64))
	charData := xml.CharData(builder.String())
	return encodeElementWithCharData(encoder, startElement, charData)
}

// A GxKMLElement is a kml element with gx: extensions.
type GxKMLElement struct {
	Child Element
}

// GxKML returns a new GxKMLElement.
func GxKML(child Element) *GxKMLElement {
	return &GxKMLElement{
		Child: child,
	}
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *GxKMLElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{
		Name: xml.Name{Space: Namespace, Local: "kml"},
		Attr: []xml.Attr{
			{
				Name:  xml.Name{Local: "xmlns:gx"},
				Value: GxNamespace,
			},
		},
	}
	return encodeElementWithChild(encoder, startElement, e.Child)
}

// Write writes e to w.
func (e *GxKMLElement) Write(w io.Writer) error {
	return write(w, e)
}

// WriteIndent writes e to w with the given prefix and indent.
func (e *GxKMLElement) WriteIndent(w io.Writer, prefix, indent string) error {
	return writeIndent(w, e, prefix, indent)
}

// A GxOptionElement is a gx:option element.
type GxOptionElement struct {
	Name    GxOptionName
	Enabled bool
}

// GxOption returns a new gx:option element.
func GxOption(name GxOptionName, enabled bool) *GxOptionElement {
	return &GxOptionElement{
		Name:    name,
		Enabled: enabled,
	}
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *GxOptionElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{
		Name: xml.Name{Local: "gx:option"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "name"}, Value: string(e.Name)},
			{Name: xml.Name{Local: "enabled"}, Value: strconv.FormatBool(e.Enabled)},
		},
	}
	return encodeElement(encoder, startElement)
}

// A GxSimpleArrayDataElement is a SimpleArrayData element.
type GxSimpleArrayDataElement struct {
	Name     string
	Children []Element
}

// GxSimpleArrayData returns a new GxSimpleArrayDataElement.
func GxSimpleArrayData(name string, children ...Element) *GxSimpleArrayDataElement {
	return &GxSimpleArrayDataElement{
		Name:     name,
		Children: children,
	}
}

// Add appends children to e and returns e as a ParentElement.
func (e *GxSimpleArrayDataElement) Add(children ...Element) ParentElement {
	return e.Append(children...)
}

// Append appends children to e and returns e.
func (e *GxSimpleArrayDataElement) Append(children ...Element) *GxSimpleArrayDataElement {
	e.Children = append(e.Children, children...)
	return e
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *GxSimpleArrayDataElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{
		Name: xml.Name{Local: "gx:SimpleArrayData"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "name"}, Value: e.Name},
		},
	}
	return encodeElementWithChildren(encoder, startElement, e.Children)
}

// A GxSimpleArrayFieldElement is a gx:SimpleArrayField element.
type GxSimpleArrayFieldElement struct {
	Name     string
	Type     string
	Children []Element
}

// GxSimpleArrayField returns a new GxSimpleArrayFieldElement.
func GxSimpleArrayField(name, _type string, children ...Element) *GxSimpleArrayFieldElement {
	return &GxSimpleArrayFieldElement{
		Name:     name,
		Type:     _type,
		Children: children,
	}
}

// Add appends children to e and returns e as a ParentElement.
func (e *GxSimpleArrayFieldElement) Add(children ...Element) ParentElement {
	return e.Append(children...)
}

// Append appends children to e and returns e.
func (e *GxSimpleArrayFieldElement) Append(children ...Element) *GxSimpleArrayFieldElement {
	e.Children = append(e.Children, children...)
	return e
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *GxSimpleArrayFieldElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{
		Name: xml.Name{Local: "gx:SimpleArrayField"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "name"}, Value: e.Name},
			{Name: xml.Name{Local: "type"}, Value: e.Type},
		},
	}
	return encodeElementWithChildren(encoder, startElement, e.Children)
}

// GxFloat64Value returns a new GxValueElement with the given float64 value.
func GxFloat64Value(value float64) *GxValueElement {
	return GxValue(strconv.FormatFloat(value, 'f', -1, 64))
}

// GxIntValue returns a new GxValueElement with the given float64 value.
func GxIntValue(value int) *GxValueElement {
	return GxValue(strconv.Itoa(value))
}
