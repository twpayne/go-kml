// Package kml provides convenince methods for creating and writing KML documents.
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
	"strings"
	"sync"
	"time"
)

const (
	NS    = "http://www.opengis.net/kml/2.2"
	NS_GX = "http://www.google.com/kml/ext/2.2"
)

var (
	lastId      int
	lastIdMutex sync.Mutex
)

func GetId() string {
	lastIdMutex.Lock()
	lastId++
	id := lastId
	lastIdMutex.Unlock()
	return strconv.Itoa(id)
}

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
	XUnits, YUnits string
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

// MarshalXML marshals se to e. start is ignored.
func (se *SimpleElement) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(xml.CharData(se.value), se.StartElement)
}

// Write writes an XML header and se to w.
func (se *SimpleElement) Write(w io.Writer) error {
	return write(w, "", "", se)
}

// WriteIndent writes an XML and se to w.
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
	if err := e.EncodeToken(ce.End()); err != nil {
		return err
	}
	return nil
}

// Write writes an XML header and ce to w.
func (ce *CompoundElement) Write(w io.Writer) error {
	return write(w, "", "", ce)
}

// WriteIndent writes an XML and se to w.
func (ce *CompoundElement) WriteIndent(w io.Writer, prefix, indent string) error {
	return write(w, prefix, indent, ce)
}

// Id returns se's id.
func (se *SharedElement) Id() string {
	return se.id
}

func Address(value string) *SimpleElement                  { return newSEString("address", value) }
func Altitude(value float64) *SimpleElement                { return newSEFloat("altitude", value) }
func AltitudeMode(value string) *SimpleElement             { return newSEString("altitudeMode", value) }
func BalloonStyle(children ...Element) *CompoundElement    { return newCE("BalloonStyle", children) }
func Begin(value time.Time) *SimpleElement                 { return newSETime("begin", value) }
func BgColor(value color.Color) *SimpleElement             { return newSEColor("bgColor", value) }
func Camera(children ...Element) *CompoundElement          { return newCE("Camera", children) }
func Change(children ...Element) *CompoundElement          { return newCE("Change", children) }
func Color(value color.Color) *SimpleElement               { return newSEColor("color", value) }
func ColorMode(value string) *SimpleElement                { return newSEString("colorMode", value) }
func Cookie(value string) *SimpleElement                   { return newSEString("cookie", value) }
func Create(children ...Element) *CompoundElement          { return newCE("Create", children) }
func Data(children ...Element) *CompoundElement            { return newCE("Data", children) }
func Delete(children ...Element) *CompoundElement          { return newCE("Delete", children) }
func Description(value string) *SimpleElement              { return newSEString("description", value) }
func DisplayName(value string) *SimpleElement              { return newSEString("displayName", value) }
func Document(children ...Element) *CompoundElement        { return newCE("Document", children) }
func DrawOrder(value int) *SimpleElement                   { return newSEInt("drawOrder", value) }
func East(value float64) *SimpleElement                    { return newSEFloat("east", value) }
func End(value time.Time) *SimpleElement                   { return newSETime("end", value) }
func Expires(value time.Time) *SimpleElement               { return newSETime("expires", value) }
func ExtendedData(children ...Element) *CompoundElement    { return newCE("ExtendedData", children) }
func Extrude(value bool) *SimpleElement                    { return newSEBool("extrude", value) }
func Fill(value bool) *SimpleElement                       { return newSEBool("fill", value) }
func Folder(children ...Element) *CompoundElement          { return newCE("Folder", children) }
func GroundOverlay(children ...Element) *CompoundElement   { return newCE("GroundOverlay", children) }
func GxAltitudeMode(value string) *SimpleElement           { return newSEString("gx:altitudeMode", value) }
func GxAltitudeOffset(value float64) *SimpleElement        { return newSEFloat("gx:altitudeOffset", value) }
func GxBalloonVisibility(value bool) *SimpleElement        { return newSEBool("gx:balloonVisibility", value) }
func GxDelayedStart(value float64) *SimpleElement          { return newSEFloat("gx:delayedStart", value) }
func GxDuration(value float64) *SimpleElement              { return newSEFloat("gx:duration", value) }
func GxFlyTo(children ...Element) *CompoundElement         { return newCE("gx:FlyTo", children) }
func GxLabelVisibility(value bool) *SimpleElement          { return newSEBool("gx:labelVisibility", value) }
func GxLatLonQuad(children ...Element) *CompoundElement    { return newCE("gx:LatLonQuad", children) }
func GxMultiTrack(children ...Element) *CompoundElement    { return newCE("gx:MultiTrack", children) }
func GxOuterColor(value color.Color) *SimpleElement        { return newSEColor("gx:outerColor", value) }
func GxOuterWidth(value float64) *SimpleElement            { return newSEFloat("gx:outerWidth", value) }
func GxPhysicalWidth(value float64) *SimpleElement         { return newSEFloat("gx:physicalWidth", value) }
func GxPlaylist(children ...Element) *CompoundElement      { return newCE("gx:Playlist", children) }
func GxSoundCue(children ...Element) *CompoundElement      { return newCE("gx:SoundCue", children) }
func GxTour(children ...Element) *CompoundElement          { return newCE("gx:Tour", children) }
func GxTourControl(children ...Element) *CompoundElement   { return newCE("gx:TourControl", children) }
func GxTourPrimitive(children ...Element) *CompoundElement { return newCE("gx:TourPrimitive", children) }
func GxTrack(children ...Element) *CompoundElement         { return newCE("gx:Track", children) }
func GxWait(children ...Element) *CompoundElement          { return newCE("gx:Wait", children) }
func HTTPQuery(value string) *SimpleElement                { return newSEString("httpQuery", value) }
func Heading(value float64) *SimpleElement                 { return newSEFloat("heading", value) }
func HotSpot(value Vec2) *SimpleElement                    { return newSEVec2("hotSpot", value) }
func Href(value string) *SimpleElement                     { return newSEString("href", value) }
func Icon(children ...Element) *CompoundElement            { return newCE("Icon", children) }
func IconStyle(children ...Element) *CompoundElement       { return newCE("IconStyle", children) }
func InnerBoundaryIs(value Element) *CompoundElement       { return newCEElement("innerBoundaryIs", value) }
func Key(value string) *SimpleElement                      { return newSEString("key", value) }
func LabelStyle(children ...Element) *CompoundElement      { return newCE("LabelStyle", children) }
func LatLonBox(children ...Element) *CompoundElement       { return newCE("LatLonBox", children) }
func Latitude(value float64) *SimpleElement                { return newSEFloat("latitude", value) }
func LineString(children ...Element) *CompoundElement      { return newCE("LineString", children) }
func LineStyle(children ...Element) *CompoundElement       { return newCE("LineStyle", children) }
func LinearRing(children ...Element) *CompoundElement      { return newCE("LinearRing", children) }
func Link(children ...Element) *CompoundElement            { return newCE("Link", children) }
func LinkDescription(value string) *SimpleElement          { return newSEString("linkDescription", value) }
func LinkName(value string) *SimpleElement                 { return newSEString("linkName", value) }
func ListItemType(value string) *SimpleElement             { return newSEString("listItemType", value) }
func ListStyle(children ...Element) *CompoundElement       { return newCE("ListStyle", children) }
func Longitude(value float64) *SimpleElement               { return newSEFloat("longitude", value) }
func LookAt(children ...Element) *CompoundElement          { return newCE("LookAt", children) }
func MaxAltitude(value float64) *SimpleElement             { return newSEFloat("maxAltitude", value) }
func MaxFadeExtent(value int) *SimpleElement               { return newSEInt("maxFadeExtent", value) }
func MaxLodPixel(value int) *SimpleElement                 { return newSEInt("maxLodPixels", value) }
func Message(value string) *SimpleElement                  { return newSEString("message", value) }
func MinAltitude(value float64) *SimpleElement             { return newSEFloat("minAltitude", value) }
func MinFadeExtent(value int) *SimpleElement               { return newSEInt("minFadeExtent", value) }
func MinLodPixel(value int) *SimpleElement                 { return newSEInt("minLodPixels", value) }
func Model(children ...Element) *CompoundElement           { return newCE("Model", children) }
func MultiGeometry(children ...Element) *CompoundElement   { return newCE("MultiGeometry", children) }
func Name(value string) *SimpleElement                     { return newSEString("name", value) }
func NetworkLink(children ...Element) *CompoundElement     { return newCE("gx:NetworkLink", children) }
func North(value float64) *SimpleElement                   { return newSEFloat("north", value) }
func Open(value bool) *SimpleElement                       { return newSEBool("open", value) }
func OuterBoundaryIs(value Element) *CompoundElement       { return newCEElement("outerBoundaryIs", value) }
func Outline(value bool) *SimpleElement                    { return newSEBool("outline", value) }
func OverlayXY(value Vec2) *SimpleElement                  { return newSEVec2("overlayXY", value) }
func Pair(children ...Element) *CompoundElement            { return newCE("Pair", children) }
func PhoneNumber(value string) *SimpleElement              { return newSEString("phoneNumber", value) }
func Placemark(children ...Element) *CompoundElement       { return newCE("Placemark", children) }
func Point(children ...Element) *CompoundElement           { return newCE("Point", children) }
func PolyStyle(children ...Element) *CompoundElement       { return newCE("PolyStyle", children) }
func Polygon(children ...Element) *CompoundElement         { return newCE("Polygon", children) }
func Range(value float64) *SimpleElement                   { return newSEFloat("range", value) }
func RefreshInterval(value float64) *SimpleElement         { return newSEFloat("refreshInterval", value) }
func RefreshMode(value string) *SimpleElement              { return newSEString("refreshMode", value) }
func Region(children ...Element) *CompoundElement          { return newCE("Region", children) }
func Roll(value float64) *SimpleElement                    { return newSEFloat("roll", value) }
func Rotation(value float64) *SimpleElement                { return newSEFloat("rotation", value) }
func RotationXY(value Vec2) *SimpleElement                 { return newSEVec2("rotationXY", value) }
func Scale(value float64) *SimpleElement                   { return newSEFloat("scale", value) }
func ScreenOverlay(children ...Element) *CompoundElement   { return newCE("ScreenOverlay", children) }
func ScreenXY(value Vec2) *SimpleElement                   { return newSEVec2("screenXY", value) }
func Size(value Vec2) *SimpleElement                       { return newSEVec2("size", value) }
func Snippet(value string) *SimpleElement                  { return newSEString("snippet", value) }
func South(value float64) *SimpleElement                   { return newSEFloat("south", value) }
func Style(id string, children ...Element) *SharedElement  { return newSharedE("Style", id, children) }
func StyleURL(style *SharedElement) *SimpleElement         { return newSEString("styleUrl", "#"+style.Id()) }
func TargetHref(value string) *SimpleElement               { return newSEString("targetHref", value) }
func Tesselate(value bool) *SimpleElement                  { return newSEBool("tesselate", value) }
func Text(value string) *SimpleElement                     { return newSEString("text", value) }
func Tilt(value float64) *SimpleElement                    { return newSEFloat("tilt", value) }
func TimeSpan(children ...Element) *CompoundElement        { return newCE("TimeSpan", children) }
func TimeStamp(children ...Element) *CompoundElement       { return newCE("TimeStamp", children) }
func Value(value string) *SimpleElement                    { return newSEString("value", value) }
func ViewBoundScale(value float64) *SimpleElement          { return newSEFloat("viewBoundScale", value) }
func ViewFormat(value string) *SimpleElement               { return newSEString("viewFormat", value) }
func ViewRefreshMode(value string) *SimpleElement          { return newSEString("viewRefreshMode", value) }
func ViewRefreshTime(value float64) *SimpleElement         { return newSEFloat("viewRefreshTime", value) }
func Visibility(value bool) *SimpleElement                 { return newSEBool("visibility", value) }
func West(value float64) *SimpleElement                    { return newSEFloat("west", value) }
func When(value time.Time) *SimpleElement                  { return newSETime("when", value) }
func Width(value float64) *SimpleElement                   { return newSEFloat("width", value) }

func coordinates(value string) *SimpleElement {
	return &SimpleElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: "coordinates"},
		},
		value: value,
	}
}

func Coordinates(value ...Coordinate) *SimpleElement {
	cs := make([]string, len(value))
	for i, c := range value {
		cs[i] = strconv.FormatFloat(c.Lon, 'f', -1, 64) + "," + strconv.FormatFloat(c.Lat, 'f', -1, 64)
		if c.Alt != 0 {
			cs[i] += "," + strconv.FormatFloat(c.Alt, 'f', -1, 64)
		}
	}
	return coordinates(strings.Join(cs, " "))
}

func CoordinatesArray(value ...[]float64) *SimpleElement {
	cs := make([]string, len(value))
	for i, c := range value {
		if len(c) < 2 {
			continue
		}
		cs[i] = strconv.FormatFloat(c[0], 'f', -1, 64) + "," + strconv.FormatFloat(c[1], 'f', -1, 64)
		if len(c) >= 3 && c[2] != 0 {
			cs[i] += "," + strconv.FormatFloat(c[2], 'f', -1, 64)
		}
	}
	return coordinates(strings.Join(cs, " "))
}

func CoordinatesFlat(flatCoords []float64, offset, end, stride, dim int) *SimpleElement {
	cs := make([]string, (end-offset)/stride)
	src := offset
	for dst := range cs {
		cs[dst] = strconv.FormatFloat(flatCoords[src], 'f', -1, 64) + "," + strconv.FormatFloat(flatCoords[src+1], 'f', -1, 64)
		if dim > 2 && flatCoords[src+2] != 0 {
			cs[dst] += "," + strconv.FormatFloat(flatCoords[src+2], 'f', -1, 64)
		}
		src += stride
	}
	return coordinates(strings.Join(cs, " "))
}

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

func GxAnimatedUpdate(children ...Element) *CompoundElement {
	return newCE("gx:AnimatedUpdate", children)
}

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

func GxSimpleArrayField(name, type_ string) *CompoundElement {
	return &CompoundElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: "gx:SimpleArrayField"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "name"}, Value: name},
				{Name: xml.Name{Local: "type"}, Value: type_},
			},
		},
	}
}

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

func NetworkLinkControl(children ...Element) *CompoundElement {
	return newCE("gx:NetworkLinkControl", children)
}

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

func SchemaData(schemaUrl string, children ...Element) *CompoundElement {
	return &CompoundElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: "SchemaData"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "schemaUrl"}, Value: schemaUrl},
			},
		},
		children: children,
	}
}

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

func SimpleField(name, type_ string, children ...Element) *CompoundElement {
	return &CompoundElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: "SimpleField"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "name"}, Value: name},
				{Name: xml.Name{Local: "type"}, Value: type_},
			},
		},
		children: children,
	}
}

func StyleMap(id string, children ...Element) *SharedElement {
	return newSharedE("StyleMap", id, children)
}

func KML(children ...Element) *CompoundElement {
	return &CompoundElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Space: NS, Local: "kml"},
		},
		children: children,
	}
}

func GxKML(children ...Element) *CompoundElement {
	kml := KML(children...)
	// FIXME find a more correct way to do this
	kml.Attr = append(kml.Attr, xml.Attr{Name: xml.Name{Local: "xmlns:gx"}, Value: NS_GX})
	return kml
}

func write(w io.Writer, prefix, indent string, m xml.Marshaler) error {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	e := xml.NewEncoder(w)
	e.Indent(prefix, indent)
	if err := e.Encode(m); err != nil {
		return err
	}
	return nil
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
				{Name: xml.Name{Local: "xunits"}, Value: value.XUnits},
				{Name: xml.Name{Local: "yunits"}, Value: value.YUnits},
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

func newSEVoid(name string) *SimpleElement {
	return &SimpleElement{
		StartElement: xml.StartElement{Name: xml.Name{Local: name}},
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

func newCEElement(name string, child Element) *CompoundElement {
	return newCE(name, []Element{child})
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
