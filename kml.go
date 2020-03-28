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

const (
	// Namespace is the default namespace.
	Namespace = "http://www.opengis.net/kml/2.2"
	// GxNamespace is the default namespace for Google Earth extensions.
	GxNamespace = "http://www.google.com/kml/ext/2.2"
)

var (
	coordinatesStartElement = xml.StartElement{
		Name: xml.Name{
			Local: "coordinates",
		},
	}
	coordinatesEndElement = coordinatesStartElement.End()
)

// A GxAngle represents an angle.
type GxAngle struct {
	Heading, Tilt, Roll float64
}

// A Coordinate represents a single geographical coordinate.
// Lon and Lat are in degrees, Alt is in meters.
type Coordinate struct {
	Lon, Lat, Alt float64
}

// A Vec2 represents a screen position.
type Vec2 struct {
	X, Y           float64
	XUnits, YUnits UnitsEnum
}

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

// CoordinatesElement is a coordinates element.
type CoordinatesElement struct {
	coordinates []Coordinate
}

// CoordinatesArrayElement is a coordinates element.
type CoordinatesArrayElement struct {
	coordinates [][]float64
}

// CoordinatesFlatElement is a coordinates element.
type CoordinatesFlatElement struct {
	flatCoords               []float64
	offset, end, stride, dim int
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

// MarshalXML marshals ce to e. start is ignored.
func (ce *CoordinatesElement) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(coordinatesStartElement); err != nil {
		return err
	}
	for i, c := range ce.coordinates {
		s := ""
		if i != 0 {
			s = " "
		}
		s += strconv.FormatFloat(c.Lon, 'f', -1, 64) + "," + strconv.FormatFloat(c.Lat, 'f', -1, 64)
		if c.Alt != 0 {
			s += "," + strconv.FormatFloat(c.Alt, 'f', -1, 64)
		}
		if err := e.EncodeToken(xml.CharData([]byte(s))); err != nil {
			return err
		}
	}
	return e.EncodeToken(coordinatesEndElement)
}

// Write writes an XML header and ce to w.
func (ce *CoordinatesElement) Write(w io.Writer) error {
	return write(w, "", "  ", ce)
}

// WriteIndent writes an XML header and ce to w.
func (ce *CoordinatesElement) WriteIndent(w io.Writer, prefix, indent string) error {
	return write(w, prefix, indent, ce)
}

// MarshalXML marshals cae to e. start is ignored.
func (cae *CoordinatesArrayElement) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(coordinatesStartElement); err != nil {
		return err
	}
	for i, c := range cae.coordinates {
		s := ""
		if i != 0 {
			s = " "
		}
		s += strconv.FormatFloat(c[0], 'f', -1, 64) + "," + strconv.FormatFloat(c[1], 'f', -1, 64)
		if len(c) > 2 && c[2] != 0 {
			s += "," + strconv.FormatFloat(c[2], 'f', -1, 64)
		}
		if err := e.EncodeToken(xml.CharData([]byte(s))); err != nil {
			return err
		}
	}
	return e.EncodeToken(coordinatesEndElement)
}

// Write writes an XML header and cae to w.
func (cae *CoordinatesArrayElement) Write(w io.Writer) error {
	return write(w, "", "  ", cae)
}

// WriteIndent writes an XML header and cae to w.
func (cae *CoordinatesArrayElement) WriteIndent(w io.Writer, prefix, indent string) error {
	return write(w, prefix, indent, cae)
}

// MarshalXML marshals cfe to e. start is ignored.
func (cfe *CoordinatesFlatElement) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(coordinatesStartElement); err != nil {
		return err
	}
	for i := cfe.offset; i < cfe.end; i += cfe.stride {
		s := ""
		if i != cfe.offset {
			s = " "
		}
		s += strconv.FormatFloat(cfe.flatCoords[i], 'f', -1, 64) + "," + strconv.FormatFloat(cfe.flatCoords[i+1], 'f', -1, 64)
		if cfe.dim > 2 && cfe.flatCoords[i+2] != 0 {
			s += "," + strconv.FormatFloat(cfe.flatCoords[i+2], 'f', -1, 64)
		}
		if err := e.EncodeToken(xml.CharData([]byte(s))); err != nil {
			return err
		}
	}
	return e.EncodeToken(coordinatesEndElement)
}

// Write writes an XML header and cfe to w.
func (cfe *CoordinatesFlatElement) Write(w io.Writer) error {
	return write(w, "", "  ", cfe)
}

// WriteIndent writes an XML header and cfe to w.
func (cfe *CoordinatesFlatElement) WriteIndent(w io.Writer, prefix, indent string) error {
	return write(w, prefix, indent, cfe)
}

// Coordinates returns a new CoordinatesElement.
func Coordinates(value ...Coordinate) *CoordinatesElement {
	return &CoordinatesElement{coordinates: value}
}

// CoordinatesArray returns a new CoordinatesArrayElement.
func CoordinatesArray(value ...[]float64) *CoordinatesArrayElement {
	return &CoordinatesArrayElement{coordinates: value}
}

// CoordinatesFlat returns a new Coordinates element from flat coordinates.
func CoordinatesFlat(flatCoords []float64, offset, end, stride, dim int) *CoordinatesFlatElement {
	return &CoordinatesFlatElement{
		flatCoords: flatCoords,
		offset:     offset,
		end:        end,
		stride:     stride,
		dim:        dim,
	}
}

// GxAngles returns a new gx:angles element.
func GxAngles(value GxAngle) *SimpleElement {
	return &SimpleElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: "gx:angles"},
		},
		value: strconv.FormatFloat(value.Heading, 'f', -1, 64) + " " +
			strconv.FormatFloat(value.Tilt, 'f', -1, 64) + " " +
			strconv.FormatFloat(value.Roll, 'f', -1, 64),
	}
}

// GxCoord returns a new gx:coord element.
func GxCoord(value Coordinate) *SimpleElement {
	return &SimpleElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: "gx:coord"},
		},
		value: strconv.FormatFloat(value.Lon, 'f', -1, 64) + " " +
			strconv.FormatFloat(value.Lat, 'f', -1, 64) + " " +
			strconv.FormatFloat(value.Alt, 'f', -1, 64),
	}
}

// GxSimpleArrayField returns a new gx:SimpleArrayField element.
func GxSimpleArrayField(name, _type string) *CompoundElement {
	return &CompoundElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: "gx:SimpleArrayField"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "name"}, Value: name},
				{Name: xml.Name{Local: "type"}, Value: _type},
			},
		},
	}
}

// LinkSnippet returns a new linkSnippet element.
func LinkSnippet(maxLines int, value string) *SimpleElement {
	return &SimpleElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: "linkSnippet"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "maxLines"}, Value: strconv.Itoa(maxLines)},
			},
		},
		value: value,
	}
}

// Schema returns a new Schema element.
func Schema(id, name string, children ...Element) *SharedElement {
	return &SharedElement{
		CompoundElement: CompoundElement{
			StartElement: xml.StartElement{
				Name: xml.Name{Local: "Schema"},
				Attr: []xml.Attr{
					{Name: xml.Name{Local: "id"}, Value: id},
					{Name: xml.Name{Local: "name"}, Value: name},
				},
			},
			children: children,
		},
		id: id,
	}
}

// SchemaData returns a new SchemaData element.
func SchemaData(schemaURL string, children ...Element) *CompoundElement {
	return &CompoundElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: "SchemaData"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "schemaUrl"}, Value: schemaURL},
			},
		},
		children: children,
	}
}

// SharedStyle returns a new shared Style element.
func SharedStyle(id string, children ...Element) *SharedElement {
	return newSharedE("Style", id, children)
}

// SharedStyleMap returns a new shared StyleMap element.
func SharedStyleMap(id string, children ...Element) *SharedElement {
	return newSharedE("StyleMap", id, children)
}

// SimpleData returns a new SimpleData element.
func SimpleData(name, value string) *SimpleElement {
	return &SimpleElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: "SimpleData"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "name"}, Value: name},
			},
		},
		value: value,
	}
}

// SimpleField returns a new SimpleField element.
func SimpleField(name, _type string, children ...Element) *CompoundElement {
	return &CompoundElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: "SimpleField"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "name"}, Value: name},
				{Name: xml.Name{Local: "type"}, Value: _type},
			},
		},
		children: children,
	}
}

// KML returns a new kml element.
func KML(child Element) *CompoundElement {
	return &CompoundElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Space: Namespace, Local: "kml"},
		},
		children: []Element{child},
	}
}

// GxKML returns a new kml element with Google Earth extensions.
func GxKML(child Element) *CompoundElement {
	kml := KML(child)
	// FIXME find a more correct way to do this
	kml.Attr = append(kml.Attr, xml.Attr{Name: xml.Name{Local: "xmlns:gx"}, Value: GxNamespace})
	return kml
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

func newSEString(name string, value string) *SimpleElement {
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
