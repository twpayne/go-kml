package kml

import (
	"encoding/xml"
	"fmt"
	"image/color"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	Header = xml.Header
	NS     = "http://www.opengis.net/kml/2.2"
	NS_GX  = "http://www.google.com/kml/ext/2.2"
)

var (
	id      int
	idMutex sync.Mutex
)

type Coordinate struct {
	Lon, Lat, Alt float64
}

type Vec2 struct {
	X, Y           float64
	XUnits, YUnits string
}

type SimpleElement struct {
	xml.StartElement
	value string
}

type CompoundElement struct {
	xml.StartElement
	id       int
	children []xml.Token
}

func getId() int {
	idMutex.Lock()
	result := id
	id++
	idMutex.Unlock()
	return result
}

func (v2 Vec2) Attr() []xml.Attr {
	return []xml.Attr{
		{Name: xml.Name{Local: "x"}, Value: strconv.FormatFloat(v2.X, 'f', -1, 64)},
		{Name: xml.Name{Local: "y"}, Value: strconv.FormatFloat(v2.Y, 'f', -1, 64)},
		{Name: xml.Name{Local: "xunits"}, Value: v2.XUnits},
		{Name: xml.Name{Local: "yunits"}, Value: v2.YUnits},
	}
}

func (se *SimpleElement) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(se.StartElement); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.CharData(se.value)); err != nil {
		return err
	}
	endElement := xml.EndElement{Name: se.Name}
	if err := e.EncodeToken(endElement); err != nil {
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

func newSEPosition(name string, value Vec2) *SimpleElement {
	return &SimpleElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: name},
			Attr: value.Attr(),
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

func newCE(name string, children ...xml.Token) *CompoundElement {
	return &CompoundElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: name},
		},
		children: children,
		id:       getId(),
	}
}

func (ce *CompoundElement) Add(children ...xml.Token) *CompoundElement {
	ce.children = append(ce.children, children...)
	return ce
}

func (ce *CompoundElement) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(ce.StartElement); err != nil {
		return err
	}
	for _, c := range ce.children {
		if err := e.EncodeElement(c, ce.StartElement); err != nil {
			return err
		}
	}
	endElement := xml.EndElement{Name: ce.Name}
	if err := e.EncodeToken(endElement); err != nil {
		return err
	}
	return nil
}

func Altitude(value int) *SimpleElement                    { return newSEInt("altitude", value) }
func AltitudeMode(value string) *SimpleElement             { return newSEString("altitudeMode", value) }
func BalloonStyle(children ...xml.Token) *CompoundElement  { return newCE("BalloonStyle", children) }
func Begin(value time.Time) *SimpleElement                 { return newSETime("begin", value) }
func BgColor(value color.Color) *SimpleElement             { return newSEColor("bgColor", value) }
func Camera(children ...xml.Token) *CompoundElement        { return newCE("Camera", children) }
func Color(value color.Color) *SimpleElement               { return newSEColor("color", value) }
func Data(children ...xml.Token) *CompoundElement          { return newCE("Data", children) }
func Description(value string) *SimpleElement              { return newSEString("description", value) }
func Document(children ...xml.Token) *CompoundElement      { return newCE("Document", children) }
func East(value float64) *SimpleElement                    { return newSEFloat("east", value) }
func End(value time.Time) *SimpleElement                   { return newSETime("end", value) }
func Extrude(value bool) *SimpleElement                    { return newSEBool("extrude", value) }
func Folder(children ...xml.Token) *CompoundElement        { return newCE("Folder", children) }
func GroundOverlay(children ...xml.Token) *CompoundElement { return newCE("GroundOverlay", children) }
func Heading(value float64) *SimpleElement                 { return newSEFloat("heading", value) }
func Href(value *url.URL) *SimpleElement                   { return newSEString("href", value.String()) }
func HotSpot(value Vec2) *SimpleElement                    { return newSEPosition("hotSpot", value) }
func Icon(children ...xml.Token) *CompoundElement          { return newCE("Icon", children) }
func IconStyle(children ...xml.Token) *CompoundElement     { return newCE("IconStyle", children) }
func LabelStyle(children ...xml.Token) *CompoundElement    { return newCE("LabelStyle", children) }
func LatLonBox(children ...xml.Token) *CompoundElement     { return newCE("LatLonBox", children) }
func Latitude(value float64) *SimpleElement                { return newSEFloat("latitude", value) }
func LineString(children ...xml.Token) *CompoundElement    { return newCE("LineString", children) }
func LineStyle(children ...xml.Token) *CompoundElement     { return newCE("LineStyle", children) }
func ListItemType(value string) *SimpleElement             { return newSEString("listItemType", value) }
func ListStyle(children ...xml.Token) *CompoundElement     { return newCE("ListStyle", children) }
func Longitude(value float64) *SimpleElement               { return newSEFloat("longitude", value) }
func MultiGeometry(children ...xml.Token) *CompoundElement { return newCE("MultiGeometry", children) }
func Name(value string) *SimpleElement                     { return newSEString("name", value) }
func North(value float64) *SimpleElement                   { return newSEFloat("north", value) }
func Open(value bool) *SimpleElement                       { return newSEBool("open", value) }
func OverlayXY(value Vec2) *SimpleElement                  { return newSEPosition("overlayXY", value) }
func Placemark(children ...xml.Token) *CompoundElement     { return newCE("Placemark", children) }
func Point(children ...xml.Token) *CompoundElement         { return newCE("Point", children) }
func PolyStyle(children ...xml.Token) *CompoundElement     { return newCE("PolyStyle", children) }
func Roll(value float64) *SimpleElement                    { return newSEFloat("roll", value) }
func Rotation(value float64) *SimpleElement                { return newSEFloat("rotation", value) }
func Scale(value float64) *SimpleElement                   { return newSEFloat("scale", value) }
func ScreenOverlay(children ...xml.Token) *CompoundElement { return newCE("ScreenOverlay", children) }
func ScreenXY(value Vec2) *SimpleElement                   { return newSEPosition("screenXY", value) }
func Snippet(value string) *SimpleElement                  { return newSEString("snippet", value) }
func South(value float64) *SimpleElement                   { return newSEFloat("south", value) }
func Style(children ...xml.Token) *CompoundElement         { return newCE("Style", children) }
func Tesselate(value bool) *SimpleElement                  { return newSEBool("tesselate", value) }
func Text(value string) *SimpleElement                     { return newSEString("text", value) }
func Tilt(value float64) *SimpleElement                    { return newSEFloat("tilt", value) }
func TimeSpan(children ...xml.Token) *CompoundElement      { return newCE("TimeSpan", children) }
func Value(value string) *SimpleElement                    { return newSEString("value", value) }
func Visibility(value bool) *SimpleElement                 { return newSEBool("visibility", value) }
func West(value float64) *SimpleElement                    { return newSEFloat("west", value) }
func When(value time.Time) *SimpleElement                  { return newSETime("time", value) }
func Width(value float64) *SimpleElement                   { return newSEFloat("width", value) }

func Coordinates(value ...Coordinate) *SimpleElement {
	cs := make([]string, len(value))
	for i, c := range value {
		cs[i] = strconv.FormatFloat(c.Lon, 'f', -1, 64) + "," +
			strconv.FormatFloat(c.Lat, 'f', -1, 64) + "," +
			strconv.FormatFloat(c.Alt, 'f', -1, 64)
	}
	return &SimpleElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: "coordinates"},
		},
		value: strings.Join(cs, " "),
	}
}

func HrefMustParse(value string) *SimpleElement {
	url, err := url.Parse(value)
	if err != nil {
		panic(err)
	}
	return Href(url)
}

func KML(children ...xml.Token) *CompoundElement {
	return &CompoundElement{
		StartElement: xml.StartElement{
			Name: xml.Name{Space: NS, Local: "kml"},
		},
		children: children,
		id:       getId(),
	}
}
