package kml

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Namespace is the default namespace.
const Namespace = "http://www.opengis.net/kml/2.2"

const defaultLinkSnippetMaxLines = 2

// A Coordinate is a single geographical coordinate.
type Coordinate struct {
	Lon float64 // Longitude in degrees.
	Lat float64 // Latitude in degrees.
	Alt float64 // Altitude in meters.
}

// CoordinatesElement is a coordinates element composed of Coordinates.
type CoordinatesElement []Coordinate

// Coordinates returns a new CoordinatesElement.
func Coordinates(value ...Coordinate) CoordinatesElement {
	return CoordinatesElement(value)
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e CoordinatesElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{Name: xml.Name{Local: "coordinates"}}
	var builder strings.Builder
	builder.Grow(3 * float64StringSize * len(e))
	for i, c := range e {
		if i != 0 {
			builder.WriteByte(' ')
		}
		builder.WriteString(strconv.FormatFloat(c.Lon, 'f', -1, 64))
		builder.WriteByte(',')
		builder.WriteString(strconv.FormatFloat(c.Lat, 'f', -1, 64))
		if c.Alt != 0 {
			builder.WriteByte(',')
			builder.WriteString(strconv.FormatFloat(c.Alt, 'f', -1, 64))
		}
	}
	charData := xml.CharData(builder.String())
	return encodeElementWithCharData(encoder, startElement, charData)
}

// CoordinatesFlatElement is a coordinates element composed of flat coordinates.
type CoordinatesFlatElement struct {
	FlatCoords []float64
	Offset     int
	End        int
	Stride     int
	Dim        int
}

// CoordinatesFlat returns a new Coordinates element from flat coordinates.
func CoordinatesFlat(flatCoords []float64, offset, end, stride, dim int) *CoordinatesFlatElement {
	return &CoordinatesFlatElement{
		FlatCoords: flatCoords,
		Offset:     offset,
		End:        end,
		Stride:     stride,
		Dim:        dim,
	}
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *CoordinatesFlatElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{Name: xml.Name{Local: "coordinates"}}
	var builder strings.Builder
	builder.Grow(3 * float64StringSize * (e.End - e.Offset) / e.Stride)
	for i := e.Offset; i < e.End; i += e.Stride {
		if i != e.Offset {
			builder.WriteByte(' ')
		}
		builder.WriteString(strconv.FormatFloat(e.FlatCoords[i], 'f', -1, 64))
		builder.WriteByte(',')
		builder.WriteString(strconv.FormatFloat(e.FlatCoords[i+1], 'f', -1, 64))
		if e.Dim > 2 && e.FlatCoords[i+2] != 0 {
			builder.WriteByte(',')
			builder.WriteString(strconv.FormatFloat(e.FlatCoords[i+2], 'f', -1, 64))
		}
	}
	charData := xml.CharData(builder.String())
	return encodeElementWithCharData(encoder, startElement, charData)
}

// CoordinatesSliceElement is a coordinates element composed of a slice of []float64s.
type CoordinatesSliceElement [][]float64

// CoordinatesSlice returns a new CoordinatesArrayElement.
func CoordinatesSlice(value ...[]float64) CoordinatesSliceElement {
	return CoordinatesSliceElement(value)
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e CoordinatesSliceElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{Name: xml.Name{Local: "coordinates"}}
	var builder strings.Builder
	builder.Grow(3 * float64StringSize * len(e))
	for i, c := range e {
		if i != 0 {
			builder.WriteByte(' ')
		}
		builder.WriteString(strconv.FormatFloat(c[0], 'f', -1, 64))
		builder.WriteByte(',')
		builder.WriteString(strconv.FormatFloat(c[1], 'f', -1, 64))
		if len(c) > 2 && c[2] != 0 {
			builder.WriteByte(',')
			builder.WriteString(strconv.FormatFloat(c[2], 'f', -1, 64))
		}
	}
	charData := xml.CharData(builder.String())
	return encodeElementWithCharData(encoder, startElement, charData)
}

// A DataElement is a Data element.
type DataElement struct {
	Name     string
	Children []Element
}

// Data returns a new DataElement.
func Data(name string, children ...Element) *DataElement {
	return &DataElement{
		Name:     name,
		Children: children,
	}
}

// Append appends children to e.
func (e *DataElement) Append(children ...Element) *DataElement {
	e.Children = append(e.Children, children...)
	return e
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *DataElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{
		Name: xml.Name{Local: "Data"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "name"}, Value: e.Name},
		},
	}
	return encodeElementWithChildren(encoder, startElement, e.Children)
}

// A KMLElement is a kml element.
type KMLElement struct { //nolint:revive
	Child Element
}

// KML returns a new KMLElement.
func KML(child Element) *KMLElement {
	return &KMLElement{
		Child: child,
	}
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *KMLElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{
		Name: xml.Name{Space: Namespace, Local: "kml"},
	}
	return encodeElementWithChild(encoder, startElement, e.Child)
}

// Write writes e to w.
func (e *KMLElement) Write(w io.Writer) error {
	return write(w, e)
}

// WriteIndent writes e to w with the given prefix and indent.
func (e *KMLElement) WriteIndent(w io.Writer, prefix, indent string) error {
	return writeIndent(w, e, prefix, indent)
}

// A LinkSnippetElement is a LinkSnippet element.
type LinkSnippetElement struct {
	MaxLines int
	Value    string
}

// LinkSnippet returns a new LinkSnippetElement.
func LinkSnippet(value string) *LinkSnippetElement {
	return &LinkSnippetElement{
		MaxLines: defaultLinkSnippetMaxLines,
		Value:    value,
	}
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *LinkSnippetElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{Name: xml.Name{Local: "linkSnippet"}}
	if e.MaxLines != defaultLinkSnippetMaxLines {
		startElement.Attr = []xml.Attr{
			{Name: xml.Name{Local: "maxLines"}, Value: strconv.Itoa(e.MaxLines)},
		}
	}
	charData := xml.CharData(e.Value)
	return encodeElementWithCharData(encoder, startElement, charData)
}

// WithMaxLines sets e's maxLines attribute.
func (e *LinkSnippetElement) WithMaxLines(maxLines int) *LinkSnippetElement {
	e.MaxLines = maxLines
	return e
}

// A ModelScaleElement is a Scale element.
type ModelScaleElement struct {
	Children []Element
}

// ModelScale returns a new ModelScaleElement.
func ModelScale(children ...Element) *ModelScaleElement {
	return &ModelScaleElement{
		Children: children,
	}
}

// Append appends children to e.
func (e *ModelScaleElement) Append(children ...Element) *ModelScaleElement {
	e.Children = append(e.Children, children...)
	return e
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *ModelScaleElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{Name: xml.Name{Local: "Scale"}}
	return encodeElementWithChildren(encoder, startElement, e.Children)
}

// A SchemaElement is a Schema element.
type SchemaElement struct {
	ID       string
	Name     string
	Children []Element
}

// NamedSchema returns a new SchemaElement with the given name.
func NamedSchema(id, name string, children ...Element) *SchemaElement {
	return &SchemaElement{
		ID:       id,
		Name:     name,
		Children: children,
	}
}

// Schema returns a new SchemaElement.
func Schema(id string, children ...Element) *SchemaElement {
	return &SchemaElement{
		ID:       id,
		Children: children,
	}
}

// Append appends children to e.
func (e *SchemaElement) Append(children ...Element) *SchemaElement {
	e.Children = append(e.Children, children...)
	return e
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *SchemaElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{
		Name: xml.Name{Local: "Schema"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "id"}, Value: e.ID},
		},
	}
	if e.Name != "" {
		startElement.Attr = append(startElement.Attr,
			xml.Attr{Name: xml.Name{Local: "name"}, Value: e.Name},
		)
	}
	return encodeElementWithChildren(encoder, startElement, e.Children)
}

// WithName sets e's name.
func (e *SchemaElement) WithName(name string) *SchemaElement {
	e.Name = name
	return e
}

// URL return e's URL.
func (e *SchemaElement) URL() string {
	if e.ID == "" {
		return ""
	}
	return "#" + e.ID
}

// A SchemaDataElement is a SchemaData element.
type SchemaDataElement struct {
	SchemaURL string
	Children  []Element
}

// SchemaData returns a new SchemaDataElement.
func SchemaData(schemaURL string, children ...Element) *SchemaDataElement {
	return &SchemaDataElement{
		SchemaURL: schemaURL,
		Children:  children,
	}
}

// Append appends children to e.
func (e *SchemaDataElement) Append(children ...Element) *SchemaDataElement {
	e.Children = append(e.Children, children...)
	return e
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *SchemaDataElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{
		Name: xml.Name{Local: "SchemaData"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "schemaUrl"}, Value: e.SchemaURL},
		},
	}
	return encodeElementWithChildren(encoder, startElement, e.Children)
}

// A SimpleDataElement is a SimpleData element.
type SimpleDataElement struct {
	Name  string
	Value string
}

// SimpleData returns a new SimpleDataElement.
func SimpleData(name, value string) *SimpleDataElement {
	return &SimpleDataElement{
		Name:  name,
		Value: value,
	}
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *SimpleDataElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{
		Name: xml.Name{Local: "SimpleData"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "name"}, Value: e.Name},
		},
	}
	charData := xml.CharData(e.Value)
	return encodeElementWithCharData(encoder, startElement, charData)
}

// A SimpleFieldElement is a SimpleField element.
type SimpleFieldElement struct {
	Name     string
	Type     string
	Children []Element
}

// SimpleField returns a new SimpleFieldElement.
func SimpleField(name, _type string, children ...Element) *SimpleFieldElement {
	return &SimpleFieldElement{
		Name:     name,
		Type:     _type,
		Children: children,
	}
}

// Append appends children to e.
func (e *SimpleFieldElement) Append(children ...Element) *SimpleFieldElement {
	e.Children = append(e.Children, children...)
	return e
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *SimpleFieldElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{
		Name: xml.Name{Local: "SimpleField"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "name"}, Value: e.Name},
			{Name: xml.Name{Local: "type"}, Value: e.Type},
		},
	}
	return encodeElementWithChildren(encoder, startElement, e.Children)
}

// A SnippetElement is a snippet element.
type SnippetElement struct {
	MaxLines int
	Value    string
}

// Snippet returns a new SnippetElement.
func Snippet(value string) *SnippetElement {
	return &SnippetElement{
		MaxLines: defaultLinkSnippetMaxLines,
		Value:    value,
	}
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *SnippetElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{Name: xml.Name{Local: "Snippet"}}
	if e.MaxLines != defaultLinkSnippetMaxLines {
		startElement.Attr = []xml.Attr{
			{Name: xml.Name{Local: "maxLines"}, Value: strconv.Itoa(e.MaxLines)},
		}
	}
	charData := xml.CharData(e.Value)
	return encodeElementWithCharData(encoder, startElement, charData)
}

// WithMaxLines sets e's maxLines attribute.
func (e *SnippetElement) WithMaxLines(maxLines int) *SnippetElement {
	e.MaxLines = maxLines
	return e
}

// A StyleElement is a Style element.
type StyleElement struct {
	ID       string
	Children []Element
}

// SharedStyle returns a new StyleElement with the given id.
func SharedStyle(id string, children ...Element) *StyleElement {
	return &StyleElement{
		ID:       id,
		Children: children,
	}
}

// Style returns a new StyleElement.
func Style(children ...Element) *StyleElement {
	return &StyleElement{
		Children: children,
	}
}

// Append appends children to e.
func (e *StyleElement) Append(children ...Element) *StyleElement {
	e.Children = append(e.Children, children...)
	return e
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *StyleElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{Name: xml.Name{Local: "Style"}}
	if e.ID != "" {
		startElement.Attr = []xml.Attr{
			{Name: xml.Name{Local: "id"}, Value: e.ID},
		}
	}
	return encodeElementWithChildren(encoder, startElement, e.Children)
}

// URL return e's URL.
func (e *StyleElement) URL() string {
	if e.ID == "" {
		return ""
	}
	return "#" + e.ID
}

// WithID sets e's ID.
func (e *StyleElement) WithID(id string) *StyleElement {
	e.ID = id
	return e
}

// A StyleMapElement is a StyleMap element.
type StyleMapElement struct {
	ID       string
	Children []Element
}

// SharedStyleMap returns a new StyleMapElement with the given id.
func SharedStyleMap(id string, children ...Element) *StyleMapElement {
	return &StyleMapElement{
		ID:       id,
		Children: children,
	}
}

// StyleMap returns a new StyleMapElement.
func StyleMap(children ...Element) *StyleMapElement {
	return &StyleMapElement{
		Children: children,
	}
}

// Append appends children to e.
func (e *StyleMapElement) Append(children ...Element) *StyleMapElement {
	e.Children = append(e.Children, children...)
	return e
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *StyleMapElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{Name: xml.Name{Local: "StyleMap"}}
	if e.ID != "" {
		startElement.Attr = []xml.Attr{
			{Name: xml.Name{Local: "id"}, Value: e.ID},
		}
	}
	return encodeElementWithChildren(encoder, startElement, e.Children)
}

// URL return e's URL.
func (e *StyleMapElement) URL() string {
	if e.ID == "" {
		return ""
	}
	return "#" + e.ID
}

// WithID sets e's ID.
func (e *StyleMapElement) WithID(id string) *StyleMapElement {
	e.ID = id
	return e
}

// A ValueElement is a value element.
type ValueElement struct {
	Value any
}

// Value returns a new ValueElement.
func Value(value any) *ValueElement {
	return &ValueElement{
		Value: value,
	}
}

// MarshalXML implements encoding/xml.Marshaler.MarshalXML.
func (e *ValueElement) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	startElement := xml.StartElement{Name: xml.Name{Local: "value"}}
	charData, err := charData(e.Value)
	if err != nil {
		return err
	}
	return encodeElementWithCharData(encoder, startElement, charData)
}

// A Vec2 is a vec2.
type Vec2 struct {
	X      float64
	Y      float64
	XUnits UnitsEnum
	YUnits UnitsEnum
}

// attr returns a slice of attributes populated with v's values.
func (v *Vec2) attr() []xml.Attr {
	return []xml.Attr{
		{Name: xml.Name{Local: "x"}, Value: strconv.FormatFloat(v.X, 'f', -1, 64)},
		{Name: xml.Name{Local: "y"}, Value: strconv.FormatFloat(v.Y, 'f', -1, 64)},
		{Name: xml.Name{Local: "xunits"}, Value: string(v.XUnits)},
		{Name: xml.Name{Local: "yunits"}, Value: string(v.YUnits)},
	}
}

func charData(value any) (xml.CharData, error) {
	switch value := value.(type) {
	case nil:
		return nil, nil
	case xml.CharData:
		return value, nil
	case fmt.Stringer:
		return xml.CharData(value.String()), nil
	case []byte:
		return xml.CharData(value), nil
	case bool:
		return xml.CharData(strconv.FormatBool(value)), nil
	case complex64:
		return xml.CharData(strconv.FormatComplex(complex128(value), 'f', -1, 64)), nil
	case complex128:
		return xml.CharData(strconv.FormatComplex(value, 'f', -1, 128)), nil
	case float32:
		return xml.CharData(strconv.FormatFloat(float64(value), 'f', -1, 64)), nil
	case float64:
		return xml.CharData(strconv.FormatFloat(value, 'f', -1, 64)), nil
	case int:
		return xml.CharData(strconv.Itoa(value)), nil
	case int8:
		return xml.CharData(strconv.FormatInt(int64(value), 10)), nil
	case int16:
		return xml.CharData(strconv.FormatInt(int64(value), 10)), nil
	case int32:
		return xml.CharData(strconv.FormatInt(int64(value), 10)), nil
	case int64:
		return xml.CharData(strconv.FormatInt(value, 10)), nil
	case string:
		return xml.CharData(value), nil
	case uint:
		return xml.CharData(strconv.FormatUint(uint64(value), 10)), nil
	case uint8:
		return xml.CharData(strconv.FormatUint(uint64(value), 10)), nil
	case uint16:
		return xml.CharData(strconv.FormatUint(uint64(value), 10)), nil
	case uint32:
		return xml.CharData(strconv.FormatUint(uint64(value), 10)), nil
	case uint64:
		return xml.CharData(strconv.FormatUint(value, 10)), nil
	default:
		return nil, &unsupportedTypeError{value: value}
	}
}
