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
	"strings"
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

// CoordinatesElement is a coordinates element.
type CoordinatesElement struct {
	coordinates []Coordinate
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
	return e.EncodeToken(ce.End())
}

// Write writes an XML header and ce to w.
func (ce *CompoundElement) Write(w io.Writer) error {
	return write(w, "", "", ce)
}

// WriteIndent writes an XML and se to w.
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

// MarshalXML marshals ee to e. start is ignored.
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

// WriteIndent writes an XML and se to w.
func (ce *CoordinatesElement) WriteIndent(w io.Writer, prefix, indent string) error {
	return write(w, prefix, indent, ce)
}

// Address returns a new Address element.
func Address(value string) *SimpleElement { return newSEString("address", value) }

// Altitude returns a new Altitude element.
func Altitude(value float64) *SimpleElement { return newSEFloat("altitude", value) }

// AltitudeMode returns a new AltitudeMode element.
func AltitudeMode(value string) *SimpleElement { return newSEString("altitudeMode", value) }

// BalloonStyle returns a new BalloonStyle element.
func BalloonStyle(children ...Element) *CompoundElement { return newCE("BalloonStyle", children) }

// Begin returns a new Begin element.
func Begin(value time.Time) *SimpleElement { return newSETime("begin", value) }

// BgColor returns a new BgColor element.
func BgColor(value color.Color) *SimpleElement { return newSEColor("bgColor", value) }

// Camera returns a new Camera element.
func Camera(children ...Element) *CompoundElement { return newCE("Camera", children) }

// Change returns a new Change element.
func Change(children ...Element) *CompoundElement { return newCE("Change", children) }

// Color returns a new Color element.
func Color(value color.Color) *SimpleElement { return newSEColor("color", value) }

// ColorMode returns a new ColorMode element.
func ColorMode(value string) *SimpleElement { return newSEString("colorMode", value) }

// Cookie returns a new Cookie element.
func Cookie(value string) *SimpleElement { return newSEString("cookie", value) }

// Create returns a new Create element.
func Create(children ...Element) *CompoundElement { return newCE("Create", children) }

// Data returns a new Data element.
func Data(children ...Element) *CompoundElement { return newCE("Data", children) }

// Delete returns a new Delete element.
func Delete(children ...Element) *CompoundElement { return newCE("Delete", children) }

// Description returns a new Description element.
func Description(value string) *SimpleElement { return newSEString("description", value) }

// DisplayName returns a new DisplayName element.
func DisplayName(value string) *SimpleElement { return newSEString("displayName", value) }

// Document returns a new Document element.
func Document(children ...Element) *CompoundElement { return newCE("Document", children) }

// DrawOrder returns a new DrawOrder element.
func DrawOrder(value int) *SimpleElement { return newSEInt("drawOrder", value) }

// East returns a new East element.
func East(value float64) *SimpleElement { return newSEFloat("east", value) }

// End returns a new End element.
func End(value time.Time) *SimpleElement { return newSETime("end", value) }

// Expires returns a new Expires element.
func Expires(value time.Time) *SimpleElement { return newSETime("expires", value) }

// ExtendedData returns a new ExtendedData element.
func ExtendedData(children ...Element) *CompoundElement { return newCE("ExtendedData", children) }

// Extrude returns a new Extrude element.
func Extrude(value bool) *SimpleElement { return newSEBool("extrude", value) }

// Fill returns a new Fill element.
func Fill(value bool) *SimpleElement { return newSEBool("fill", value) }

// FlyToView returns a new FlyToView element.
func FlyToView(value bool) *SimpleElement { return newSEBool("flyToView", value) }

// Folder returns a new Folder element.
func Folder(children ...Element) *CompoundElement { return newCE("Folder", children) }

// GroundOverlay returns a new GroundOverlay element.
func GroundOverlay(children ...Element) *CompoundElement { return newCE("GroundOverlay", children) }

// GxAltitudeMode returns a new gx:AltitudeMode element.
func GxAltitudeMode(value string) *SimpleElement { return newSEString("gx:altitudeMode", value) }

// GxAltitudeOffset returns a new gx:AltitudeOffset element.
func GxAltitudeOffset(value float64) *SimpleElement { return newSEFloat("gx:altitudeOffset", value) }

// GxBalloonVisibility returns a new gx:BalloonVisibility element.
func GxBalloonVisibility(value bool) *SimpleElement { return newSEBool("gx:balloonVisibility", value) }

// GxDelayedStart returns a new gx:DelayedStart element.
func GxDelayedStart(value float64) *SimpleElement { return newSEFloat("gx:delayedStart", value) }

// GxDuration returns a new gx:Duration element.
func GxDuration(value float64) *SimpleElement { return newSEFloat("gx:duration", value) }

// GxFlyTo returns a new gx:FlyTo element.
func GxFlyTo(children ...Element) *CompoundElement { return newCE("gx:FlyTo", children) }

// GxLabelVisibility returns a new gx:LabelVisibility element.
func GxLabelVisibility(value bool) *SimpleElement { return newSEBool("gx:labelVisibility", value) }

// GxLatLonQuad returns a new gx:LatLonQuad element.
func GxLatLonQuad(children ...Element) *CompoundElement { return newCE("gx:LatLonQuad", children) }

// GxMultiTrack returns a new gx:MultiTrack element.
func GxMultiTrack(children ...Element) *CompoundElement { return newCE("gx:MultiTrack", children) }

// GxNetworkLink returns a new GxNetworkLink element.
func GxNetworkLink(children ...Element) *CompoundElement { return newCE("gx:NetworkLink", children) }

// GxOuterColor returns a new gx:OuterColor element.
func GxOuterColor(value color.Color) *SimpleElement { return newSEColor("gx:outerColor", value) }

// GxOuterWidth returns a new gx:OuterWidth element.
func GxOuterWidth(value float64) *SimpleElement { return newSEFloat("gx:outerWidth", value) }

// GxPhysicalWidth returns a new gx:PhysicalWidth element.
func GxPhysicalWidth(value float64) *SimpleElement { return newSEFloat("gx:physicalWidth", value) }

// GxPlaylist returns a new gx:Playlist element.
func GxPlaylist(children ...Element) *CompoundElement { return newCE("gx:Playlist", children) }

// GxSoundCue returns a new gx:SoundCue element.
func GxSoundCue(children ...Element) *CompoundElement { return newCE("gx:SoundCue", children) }

// GxTour returns a new gx:Tour element.
func GxTour(children ...Element) *CompoundElement { return newCE("gx:Tour", children) }

// GxTourControl returns a new gx:TourControl element.
func GxTourControl(children ...Element) *CompoundElement { return newCE("gx:TourControl", children) }

// GxTourPrimitive returns a new gx:TourPrimitive element.
func GxTourPrimitive(children ...Element) *CompoundElement { return newCE("gx:TourPrimitive", children) }

// GxTrack returns a new gx:Track element.
func GxTrack(children ...Element) *CompoundElement { return newCE("gx:Track", children) }

// GxWait returns a new gx:Wait element.
func GxWait(children ...Element) *CompoundElement { return newCE("gx:Wait", children) }

// HTTPQuery returns a new HTTPQuery element.
func HTTPQuery(value string) *SimpleElement { return newSEString("httpQuery", value) }

// Heading returns a new Heading element.
func Heading(value float64) *SimpleElement { return newSEFloat("heading", value) }

// HotSpot returns a new HotSpot element.
func HotSpot(value Vec2) *SimpleElement { return newSEVec2("hotSpot", value) }

// Href returns a new Href element.
func Href(value string) *SimpleElement { return newSEString("href", value) }

// Icon returns a new Icon element.
func Icon(children ...Element) *CompoundElement { return newCE("Icon", children) }

// IconStyle returns a new IconStyle element.
func IconStyle(children ...Element) *CompoundElement { return newCE("IconStyle", children) }

// InnerBoundaryIs returns a new InnerBoundaryIs element.
func InnerBoundaryIs(value Element) *CompoundElement { return newCEElement("innerBoundaryIs", value) }

// Key returns a new Key element.
func Key(value string) *SimpleElement { return newSEString("key", value) }

// LabelStyle returns a new LabelStyle element.
func LabelStyle(children ...Element) *CompoundElement { return newCE("LabelStyle", children) }

// LatLonBox returns a new LatLonBox element.
func LatLonBox(children ...Element) *CompoundElement { return newCE("LatLonBox", children) }

// Latitude returns a new Latitude element.
func Latitude(value float64) *SimpleElement { return newSEFloat("latitude", value) }

// LineString returns a new LineString element.
func LineString(children ...Element) *CompoundElement { return newCE("LineString", children) }

// LineStyle returns a new LineStyle element.
func LineStyle(children ...Element) *CompoundElement { return newCE("LineStyle", children) }

// LinearRing returns a new LinearRing element.
func LinearRing(children ...Element) *CompoundElement { return newCE("LinearRing", children) }

// Link returns a new Link element.
func Link(children ...Element) *CompoundElement { return newCE("Link", children) }

// LinkDescription returns a new LinkDescription element.
func LinkDescription(value string) *SimpleElement { return newSEString("linkDescription", value) }

// LinkName returns a new LinkName element.
func LinkName(value string) *SimpleElement { return newSEString("linkName", value) }

// ListItemType returns a new ListItemType element.
func ListItemType(value string) *SimpleElement { return newSEString("listItemType", value) }

// ListStyle returns a new ListStyle element.
func ListStyle(children ...Element) *CompoundElement { return newCE("ListStyle", children) }

// Longitude returns a new Longitude element.
func Longitude(value float64) *SimpleElement { return newSEFloat("longitude", value) }

// LookAt returns a new LookAt element.
func LookAt(children ...Element) *CompoundElement { return newCE("LookAt", children) }

// MaxAltitude returns a new MaxAltitude element.
func MaxAltitude(value float64) *SimpleElement { return newSEFloat("maxAltitude", value) }

// MaxFadeExtent returns a new MaxFadeExtent element.
func MaxFadeExtent(value int) *SimpleElement { return newSEInt("maxFadeExtent", value) }

// MaxLodPixel returns a new MaxLodPixel element.
func MaxLodPixel(value int) *SimpleElement { return newSEInt("maxLodPixels", value) }

// Message returns a new Message element.
func Message(value string) *SimpleElement { return newSEString("message", value) }

// MinAltitude returns a new MinAltitude element.
func MinAltitude(value float64) *SimpleElement { return newSEFloat("minAltitude", value) }

// MinFadeExtent returns a new MinFadeExtent element.
func MinFadeExtent(value int) *SimpleElement { return newSEInt("minFadeExtent", value) }

// MinLodPixel returns a new MinLodPixel element.
func MinLodPixel(value int) *SimpleElement { return newSEInt("minLodPixels", value) }

// Model returns a new Model element.
func Model(children ...Element) *CompoundElement { return newCE("Model", children) }

// MultiGeometry returns a new MultiGeometry element.
func MultiGeometry(children ...Element) *CompoundElement { return newCE("MultiGeometry", children) }

// Name returns a new Name element.
func Name(value string) *SimpleElement { return newSEString("name", value) }

// NetworkLink returns a new NetworkLink element.
func NetworkLink(children ...Element) *CompoundElement { return newCE("NetworkLink", children) }

// North returns a new North element.
func North(value float64) *SimpleElement { return newSEFloat("north", value) }

// Open returns a new Open element.
func Open(value bool) *SimpleElement { return newSEBool("open", value) }

// OuterBoundaryIs returns a new OuterBoundaryIs element.
func OuterBoundaryIs(value Element) *CompoundElement { return newCEElement("outerBoundaryIs", value) }

// Outline returns a new Outline element.
func Outline(value bool) *SimpleElement { return newSEBool("outline", value) }

// OverlayXY returns a new OverlayXY element.
func OverlayXY(value Vec2) *SimpleElement { return newSEVec2("overlayXY", value) }

// Pair returns a new Pair element.
func Pair(children ...Element) *CompoundElement { return newCE("Pair", children) }

// PhoneNumber returns a new PhoneNumber element.
func PhoneNumber(value string) *SimpleElement { return newSEString("phoneNumber", value) }

// Placemark returns a new Placemark element.
func Placemark(children ...Element) *CompoundElement { return newCE("Placemark", children) }

// Point returns a new Point element.
func Point(children ...Element) *CompoundElement { return newCE("Point", children) }

// PolyStyle returns a new PolyStyle element.
func PolyStyle(children ...Element) *CompoundElement { return newCE("PolyStyle", children) }

// Polygon returns a new Polygon element.
func Polygon(children ...Element) *CompoundElement { return newCE("Polygon", children) }

// Range returns a new Range element.
func Range(value float64) *SimpleElement { return newSEFloat("range", value) }

// RefreshInterval returns a new RefreshInterval element.
func RefreshInterval(value float64) *SimpleElement { return newSEFloat("refreshInterval", value) }

// RefreshMode returns a new RefreshMode element.
func RefreshMode(value string) *SimpleElement { return newSEString("refreshMode", value) }

// RefreshVisibility returns a new RefreshVisibility element.
func RefreshVisibility(value bool) *SimpleElement { return newSEBool("refreshVisibility", value) }

// Region returns a new Region element.
func Region(children ...Element) *CompoundElement { return newCE("Region", children) }

// Roll returns a new Roll element.
func Roll(value float64) *SimpleElement { return newSEFloat("roll", value) }

// Rotation returns a new Rotation element.
func Rotation(value float64) *SimpleElement { return newSEFloat("rotation", value) }

// RotationXY returns a new RotationXY element.
func RotationXY(value Vec2) *SimpleElement { return newSEVec2("rotationXY", value) }

// Scale returns a new Scale element.
func Scale(value float64) *SimpleElement { return newSEFloat("scale", value) }

// ScreenOverlay returns a new ScreenOverlay element.
func ScreenOverlay(children ...Element) *CompoundElement { return newCE("ScreenOverlay", children) }

// ScreenXY returns a new ScreenXY element.
func ScreenXY(value Vec2) *SimpleElement { return newSEVec2("screenXY", value) }

// Size returns a new Size element.
func Size(value Vec2) *SimpleElement { return newSEVec2("size", value) }

// Snippet returns a new Snippet element.
func Snippet(value string) *SimpleElement { return newSEString("snippet", value) }

// South returns a new South element.
func South(value float64) *SimpleElement { return newSEFloat("south", value) }

// Style returns a new Style element.
func Style(children ...Element) *CompoundElement { return newCE("Style", children) }

// StyleMap returns a new StyleMap element.
func StyleMap(children ...Element) *CompoundElement { return newCE("StyleMap", children) }

// StyleURL returns a new StyleURL element.
func StyleURL(value string) *SimpleElement { return newSEString("styleUrl", value) }

// TargetHref returns a new TargetHref element.
func TargetHref(value string) *SimpleElement { return newSEString("targetHref", value) }

// Tessellate returns a new Tessellate element.
func Tessellate(value bool) *SimpleElement { return newSEBool("tessellate", value) }

// Text returns a new Text element.
func Text(value string) *SimpleElement { return newSEString("text", value) }

// Tilt returns a new Tilt element.
func Tilt(value float64) *SimpleElement { return newSEFloat("tilt", value) }

// TimeSpan returns a new TimeSpan element.
func TimeSpan(children ...Element) *CompoundElement { return newCE("TimeSpan", children) }

// TimeStamp returns a new TimeStamp element.
func TimeStamp(children ...Element) *CompoundElement { return newCE("TimeStamp", children) }

// Value returns a new Value element.
func Value(value string) *SimpleElement { return newSEString("value", value) }

// ViewBoundScale returns a new ViewBoundScale element.
func ViewBoundScale(value float64) *SimpleElement { return newSEFloat("viewBoundScale", value) }

// ViewFormat returns a new ViewFormat element.
func ViewFormat(value string) *SimpleElement { return newSEString("viewFormat", value) }

// ViewRefreshMode returns a new ViewRefreshMode element.
func ViewRefreshMode(value string) *SimpleElement { return newSEString("viewRefreshMode", value) }

// ViewRefreshTime returns a new ViewRefreshTime element.
func ViewRefreshTime(value float64) *SimpleElement { return newSEFloat("viewRefreshTime", value) }

// Visibility returns a new Visibility element.
func Visibility(value bool) *SimpleElement { return newSEBool("visibility", value) }

// West returns a new West element.
func West(value float64) *SimpleElement { return newSEFloat("west", value) }

// When returns a new When element.
func When(value time.Time) *SimpleElement { return newSETime("when", value) }

// Width returns a new Width element.
func Width(value float64) *SimpleElement { return newSEFloat("width", value) }

func coordinates(value string) *SimpleElement {
	return &SimpleElement{
		StartElement: coordinatesStartElement,
		value:        value,
	}
}

// Coordinates returns a new CoordinatesElement.
func Coordinates(value ...Coordinate) *CoordinatesElement {
	return &CoordinatesElement{coordinates: value}
}

// CoordinatesArray returns a new Coordinates element from an array of coordinates.
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

// CoordinatesFlat returns a new Coordinates element from flat coordinates.
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

// GxAngles returns a new gx:Angles element.
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

// GxAnimatedUpdate returns a new gx:AnimatedUpdate element.
func GxAnimatedUpdate(children ...Element) *CompoundElement {
	return newCE("gx:AnimatedUpdate", children)
}

// GxCoord returns a new gx:Coord element.
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
func GxSimpleArrayField(name, typ string) *CompoundElement {
	return &CompoundElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: "gx:SimpleArrayField"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "name"}, Value: name},
				{Name: xml.Name{Local: "type"}, Value: typ},
			},
		},
	}
}

// LinkSnippet returns a new LinkSnippet element.
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

// NetworkLinkControl returns a new NetworkLinkControl element.
func NetworkLinkControl(children ...Element) *CompoundElement {
	return newCE("gx:NetworkLinkControl", children)
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
func SimpleField(name, typ string, children ...Element) *CompoundElement {
	return &CompoundElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: "SimpleField"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "name"}, Value: name},
				{Name: xml.Name{Local: "type"}, Value: typ},
			},
		},
		children: children,
	}
}

// KML returns a new kml element.
func KML(children ...Element) *CompoundElement {
	return &CompoundElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Space: Namespace, Local: "kml"},
		},
		children: children,
	}
}

// GxKML returns a new kml element with Google Earth extensions.
func GxKML(children ...Element) *CompoundElement {
	kml := KML(children...)
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
