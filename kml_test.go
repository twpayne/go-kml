package kml_test

import (
	"encoding/xml"
	"image/color"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	kml "github.com/twpayne/go-kml/v3"
)

var (
	_ kml.TopLevelElement = &kml.GxKMLElement{}
	_ kml.TopLevelElement = &kml.KMLElement{}
)

func TestSimpleElements(t *testing.T) {
	for _, tc := range []struct {
		name        string
		element     kml.Element
		expected    string
		expectedErr string
	}{
		{
			name:     "altitude",
			element:  kml.Altitude(0),
			expected: `<altitude>0</altitude>`,
		},
		{
			name:     "altitudeMode",
			element:  kml.AltitudeMode(kml.AltitudeModeAbsolute),
			expected: `<altitudeMode>absolute</altitudeMode>`,
		},
		{
			name:     "begin",
			element:  kml.Begin(time.Date(1876, 8, 1, 0, 0, 0, 0, time.UTC)),
			expected: `<begin>1876-08-01T00:00:00Z</begin>`,
		},
		{
			name:     "bgColor",
			element:  kml.BgColor(color.Black),
			expected: `<bgColor>ff000000</bgColor>`,
		},
		{
			name:     "color",
			element:  kml.Color(color.White),
			expected: `<color>ffffffff</color>`,
		},
		{
			name:     "coordinates",
			element:  kml.Coordinates(kml.Coordinate{Lon: 1.23, Lat: 4.56, Alt: 7.89}),
			expected: `<coordinates>1.23,4.56,7.89</coordinates>`,
		},
		{
			name:     "coordinatesFlat0",
			element:  kml.CoordinatesFlat([]float64{1.23, 4.56, 7.89, 0.12}, 0, 4, 2, 2),
			expected: `<coordinates>1.23,4.56 7.89,0.12</coordinates>`,
		},
		{
			name:     "coordinatesFlat1",
			element:  kml.CoordinatesFlat([]float64{1.23, 4.56, 0, 7.89, 0.12, 0}, 0, 6, 3, 3),
			expected: `<coordinates>1.23,4.56 7.89,0.12</coordinates>`,
		},
		{
			name:     "coordinatesFlat2",
			element:  kml.CoordinatesFlat([]float64{1.23, 4.56, 7.89, 0.12, 3.45, 6.78}, 0, 6, 3, 3),
			expected: `<coordinates>1.23,4.56,7.89 0.12,3.45,6.78</coordinates>`,
		},
		{
			name:     "coordinatesSlice0",
			element:  kml.CoordinatesSlice([]float64{1.23, 4.56}),
			expected: `<coordinates>1.23,4.56</coordinates>`,
		},
		{
			name:     "coordinatesSlice1",
			element:  kml.CoordinatesSlice([]float64{1.23, 4.56, 7.89}),
			expected: `<coordinates>1.23,4.56,7.89</coordinates>`,
		},
		{
			name:     "coordinatesSlice2",
			element:  kml.CoordinatesSlice([][]float64{{1.23, 4.56}, {7.89, 0.12}}...),
			expected: `<coordinates>1.23,4.56 7.89,0.12</coordinates>`,
		},
		{
			name:     "description",
			element:  kml.Description("text"),
			expected: `<description>text</description>`,
		},
		{
			name:     "end",
			element:  kml.End(time.Date(2015, 12, 31, 23, 59, 59, 0, time.UTC)),
			expected: `<end>2015-12-31T23:59:59Z</end>`,
		},
		{
			name:     "extrude",
			element:  kml.Extrude(false),
			expected: `<extrude>0</extrude>`,
		},
		{
			name:     "Folder",
			element:  kml.Folder(),
			expected: `<Folder></Folder>`,
		},
		{
			name:     "gx:angles",
			element:  kml.GxAngles(1.23, 4.56, 7.89),
			expected: `<gx:angles>1.23 4.56 7.89</gx:angles>`,
		},
		{
			name:     "gx:coord",
			element:  kml.GxCoord(kml.Coordinate{1.23, 4.56, 7.89}),
			expected: `<gx:coord>1.23 4.56 7.89</gx:coord>`,
		},
		{
			name:     "gx:option",
			element:  kml.GxOption(kml.GxOptionNameStreetView, true),
			expected: `<gx:option name="streetview" enabled="true"></gx:option>`,
		},
		{
			name:     "heading",
			element:  kml.Heading(0),
			expected: `<heading>0</heading>`,
		},
		{
			name:     "hotSpot",
			element:  kml.HotSpot(kml.Vec2{X: 0.5, Y: 0.5, XUnits: kml.UnitsPixels, YUnits: kml.UnitsPixels}),
			expected: `<hotSpot x="0.5" y="0.5" xunits="pixels" yunits="pixels"></hotSpot>`,
		},
		{
			name:     "href",
			element:  kml.Href("https://www.google.com/"),
			expected: `<href>https://www.google.com/</href>`,
		},
		{
			name:     "latitude",
			element:  kml.Latitude(0),
			expected: `<latitude>0</latitude>`,
		},
		{
			name:     "linkSnippet0",
			element:  kml.LinkSnippet("snippet"),
			expected: `<linkSnippet>snippet</linkSnippet>`,
		},
		{
			name:     "linkSnippet1",
			element:  kml.LinkSnippet("snippet").WithMaxLines(1),
			expected: `<linkSnippet maxLines="1">snippet</linkSnippet>`,
		},
		{
			name:     "listItemType",
			element:  kml.ListItemType(kml.ListItemTypeCheck),
			expected: `<listItemType>check</listItemType>`,
		},
		{
			name:     "name",
			element:  kml.Name("value"),
			expected: "<name>value</name>",
		},
		{
			name:     "overlayXY",
			element:  kml.OverlayXY(kml.Vec2{X: 0, Y: 0, XUnits: kml.UnitsFraction, YUnits: kml.UnitsFraction}),
			expected: `<overlayXY x="0" y="0" xunits="fraction" yunits="fraction"></overlayXY>`,
		},
		{
			name:     "Snippet",
			element:  kml.Snippet("snippet").WithMaxLines(1),
			expected: `<Snippet maxLines="1">snippet</Snippet>`,
		},
		{
			name:     "value_charData",
			element:  kml.Value(xml.CharData("<>")),
			expected: "<value>&lt;&gt;</value>",
		},
		{
			name:     "value_stringer",
			element:  kml.Value(kml.AltitudeModeAbsolute),
			expected: "<value>absolute</value>",
		},
		{
			name:     "value_byte_slice",
			element:  kml.Value([]byte("&")),
			expected: "<value>&amp;</value>",
		},
		{
			name:     "value_bool",
			element:  kml.Value(true),
			expected: "<value>true</value>",
		},
		{
			name:     "value_complex64",
			element:  kml.Value(complex64(1 + 2i)),
			expected: "<value>(1+2i)</value>",
		},
		{
			name:     "value_complex128",
			element:  kml.Value(1 + 2i),
			expected: "<value>(1+2i)</value>",
		},
		{
			name:     "value_float32",
			element:  kml.Value(float32(1.25)),
			expected: "<value>1.25</value>",
		},
		{
			name:     "value_float64",
			element:  kml.Value(1.2),
			expected: "<value>1.2</value>",
		},
		{
			name:     "value_int",
			element:  kml.Value(1),
			expected: "<value>1</value>",
		},
		{
			name:     "value_int8",
			element:  kml.Value(int8(-8)),
			expected: "<value>-8</value>",
		},
		{
			name:     "value_int16",
			element:  kml.Value(int16(-16)),
			expected: "<value>-16</value>",
		},
		{
			name:     "value_int32",
			element:  kml.Value(int32(-32)),
			expected: "<value>-32</value>",
		},
		{
			name:     "value_int64",
			element:  kml.Value(int64(-64)),
			expected: "<value>-64</value>",
		},
		{
			name:     "value_nil",
			element:  kml.Value(nil),
			expected: "<value></value>",
		},
		{
			name:     "value_string",
			element:  kml.Value("<>"),
			expected: "<value>&lt;&gt;</value>",
		},
		{
			name:     "value_uint",
			element:  kml.Value(uint(1)),
			expected: "<value>1</value>",
		},
		{
			name:     "value_uint8",
			element:  kml.Value(uint8(8)),
			expected: "<value>8</value>",
		},
		{
			name:     "value_uint16",
			element:  kml.Value(uint16(16)),
			expected: "<value>16</value>",
		},
		{
			name:     "value_uint32",
			element:  kml.Value(uint32(32)),
			expected: "<value>32</value>",
		},
		{
			name:     "value_uint64",
			element:  kml.Value(uint64(64)),
			expected: "<value>64</value>",
		},
		{
			name:     "value_nil",
			element:  kml.Value(nil),
			expected: "<value></value>",
		},
		{
			name:        "value_unsupported",
			element:     kml.Value(kml.Value(nil)),
			expectedErr: "*kml.ValueElement: unsupported type",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var builder strings.Builder
			err := xml.NewEncoder(&builder).Encode(tc.element)
			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, builder.String())
			}
		})
	}
}

func TestParentElements(t *testing.T) {
	for _, tc := range []struct {
		name     string
		element  kml.ParentElement
		expected string
	}{
		{
			name: "easy_trail",
			element: kml.Placemark(
				kml.Name("Easy trail"),
				kml.ExtendedData(
					kml.SchemaData("#TrailHeadTypeId",
						kml.SimpleData("TrailHeadName", "Pi in the sky"),
						kml.SimpleData("TrailLength", "3.14159"),
						kml.SimpleData("ElevationGain", "10"),
					),
				),
				kml.Point(
					kml.Coordinates(kml.Coordinate{Lon: -122.000, Lat: 37.002}),
				),
			),
			expected: `` +
				`<Placemark>` +
				`<name>Easy trail</name>` +
				`<ExtendedData>` +
				`<SchemaData schemaUrl="#TrailHeadTypeId">` +
				`<SimpleData name="TrailHeadName">Pi in the sky</SimpleData>` +
				`<SimpleData name="TrailLength">3.14159</SimpleData>` +
				`<SimpleData name="ElevationGain">10</SimpleData>` +
				`</SchemaData>` +
				`</ExtendedData>` +
				`<Point>` +
				`<coordinates>-122,37.002</coordinates>` +
				`</Point>` +
				`</Placemark>`,
		},
		{
			name: "simple_crosshairs",
			element: kml.ScreenOverlay(
				kml.Name("Simple crosshairs"),
				kml.Description("This screen overlay uses fractional positioning to put the image in the exact center of the screen"),
				kml.Icon(
					kml.Href("http://myserver/myimage.jpg"),
				),
				kml.OverlayXY(kml.Vec2{X: 0.5, Y: 0.5, XUnits: kml.UnitsFraction, YUnits: kml.UnitsFraction}),
				kml.ScreenXY(kml.Vec2{X: 0.5, Y: 0.5, XUnits: kml.UnitsFraction, YUnits: kml.UnitsFraction}),
				kml.Rotation(39.37878630116985),
				kml.Size(kml.Vec2{X: 0, Y: 0, XUnits: kml.UnitsPixels, YUnits: kml.UnitsPixels}),
			),
			expected: `` +
				`<ScreenOverlay>` +
				`<name>Simple crosshairs</name>` +
				`<description>This screen overlay uses fractional positioning to put the image in the exact center of the screen</description>` +
				`<Icon>` +
				`<href>http://myserver/myimage.jpg</href>` +
				`</Icon>` +
				`<overlayXY x="0.5" y="0.5" xunits="fraction" yunits="fraction"></overlayXY>` +
				`<screenXY x="0.5" y="0.5" xunits="fraction" yunits="fraction"></screenXY>` +
				`<rotation>39.37878630116985</rotation>` +
				`<size x="0" y="0" xunits="pixels" yunits="pixels"></size>` +
				`</ScreenOverlay>`,
		},
		{
			name: "extended_data",
			element: kml.Placemark(
				kml.Name("Club house"),
				kml.ExtendedData(
					kml.Data("holeNumber", kml.Value(1)),
					kml.Data("holeYardage", kml.Value(234)),
					kml.Data("holePar", kml.Value(4)),
				),
			),
			expected: `` +
				`<Placemark>` +
				`<name>Club house</name>` +
				`<ExtendedData>` +
				`<Data name="holeNumber">` +
				`<value>1</value>` +
				`</Data>` +
				`<Data name="holeYardage">` +
				`<value>234</value>` +
				`</Data>` +
				`<Data name="holePar">` +
				`<value>4</value>` +
				`</Data>` +
				`</ExtendedData>` +
				`</Placemark>`,
		},
		{
			name: "Schema",
			element: kml.Schema("schema",
				kml.GxSimpleArrayField("heartrate", "int", kml.DisplayName("Heart Rate")),
				kml.GxSimpleArrayField("cadence", "int", kml.DisplayName("Cadence")),
				kml.GxSimpleArrayField("power", "float", kml.DisplayName("Power")),
			),
			expected: `` +
				`<Schema id="schema">` +
				`<gx:SimpleArrayField name="heartrate" type="int">` +
				`<displayName>Heart Rate</displayName>` +
				`</gx:SimpleArrayField>` +
				`<gx:SimpleArrayField name="cadence" type="int">` +
				`<displayName>Cadence</displayName>` +
				`</gx:SimpleArrayField>` +
				`<gx:SimpleArrayField name="power" type="float">` +
				`<displayName>Power</displayName>` +
				`</gx:SimpleArrayField>` +
				`</Schema>`,
		},
		{
			name: "gx:Wait",
			element: kml.GxWait(
				kml.GxDuration(2500 * time.Millisecond),
			),
			expected: `` +
				`<gx:Wait>` +
				`<gx:duration>2.5</gx:duration>` +
				`</gx:Wait>`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var builder strings.Builder
			encoder := xml.NewEncoder(&builder)
			assert.NoError(t, encoder.Encode(tc.element))
			assert.Equal(t, tc.expected, builder.String())
		})
	}
}

func TestModel(t *testing.T) {
	k := kml.KML(
		kml.Placemark(
			kml.Name("SketchUp Model of Macky Auditorium"),
			kml.Description("University of Colorado, Boulder; model created by Noël Nemcik."),
			kml.LookAt(
				kml.Longitude(-105.2727379358738),
				kml.Latitude(40.01000594412381),
				kml.Altitude(0),
				kml.Range(127.2393107680517),
				kml.Tilt(65.74454495876547),
				kml.Heading(-27.70337734057933),
			),
			kml.Model(
				kml.AltitudeMode(kml.AltitudeModeRelativeToGround),
				kml.Location(
					kml.Longitude(-105.272774533734),
					kml.Latitude(40.009993372683),
					kml.Altitude(0),
				),
				kml.Orientation(
					kml.Heading(0),
					kml.Tilt(0),
					kml.Roll(0),
				),
				kml.ModelScale(
					kml.X(1),
					kml.Y(1),
					kml.Z(1),
				),
			),
		),
	)
	expected := `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
		`<kml xmlns="http://www.opengis.net/kml/2.2">` +
		`<Placemark>` +
		`<name>SketchUp Model of Macky Auditorium</name>` +
		`<description>University of Colorado, Boulder; model created by Noël Nemcik.</description>` +
		`<LookAt>` +
		`<longitude>-105.2727379358738</longitude>` +
		`<latitude>40.01000594412381</latitude>` +
		`<altitude>0</altitude>` +
		`<range>127.2393107680517</range>` +
		`<tilt>65.74454495876547</tilt>` +
		`<heading>-27.70337734057933</heading>` +
		`</LookAt>` +
		`<Model>` +
		`<altitudeMode>relativeToGround</altitudeMode>` +
		`<Location>` +
		`<longitude>-105.272774533734</longitude>` +
		`<latitude>40.009993372683</latitude>` +
		`<altitude>0</altitude>` +
		`</Location>` +
		`<Orientation>` +
		`<heading>0</heading>` +
		`<tilt>0</tilt>` +
		`<roll>0</roll>` +
		`</Orientation>` +
		`<Scale>` +
		`<x>1</x>` +
		`<y>1</y>` +
		`<z>1</z>` +
		`</Scale>` +
		`</Model>` +
		`</Placemark>` +
		`</kml>`
	var builder strings.Builder
	assert.NoError(t, k.Write(&builder))
	assert.Equal(t, expected, builder.String())
}

func TestSharedStyles(t *testing.T) {
	style0 := kml.SharedStyle("0")
	highlightPlacemarkStyle := kml.SharedStyle(
		"highlightPlacemark",
		kml.IconStyle(
			kml.Icon(
				kml.Href("http://maps.google.com/mapfiles/kml/paddle/red-stars.png"),
			),
		),
	)
	normalPlacemarkStyle := kml.SharedStyle(
		"normalPlacemark",
		kml.IconStyle(
			kml.Icon(
				kml.Href("http://maps.google.com/mapfiles/kml/paddle/wht-blank.png"),
			),
		),
	)
	exampleStyleMap := kml.SharedStyleMap(
		"exampleStyleMap",
		kml.Pair(
			kml.Key(kml.StyleStateNormal),
			kml.StyleURL(normalPlacemarkStyle.URL()),
		),
		kml.Pair(
			kml.Key(kml.StyleStateHighlight),
			kml.StyleURL(highlightPlacemarkStyle.URL()),
		),
	)
	for _, tc := range []struct {
		name     string
		element  kml.Element
		expected string
	}{
		{
			name: "folder",
			element: kml.Folder(
				style0,
				kml.Placemark(
					kml.StyleURL(style0.URL()),
				),
			),
			expected: `<Folder>` +
				`<Style id="0">` +
				`</Style>` +
				`<Placemark>` +
				`<styleUrl>#0</styleUrl>` +
				`</Placemark>` +
				`</Folder>`,
		},
		{
			name: "highlighted_icon",
			element: kml.KML(
				kml.Document(
					kml.Name("Highlighted Icon"),
					kml.Description("Place your mouse over the icon to see it display the new icon"),
					highlightPlacemarkStyle,
					normalPlacemarkStyle,
					exampleStyleMap,
					kml.Placemark(
						kml.Name("Roll over this icon"),
						kml.StyleURL(exampleStyleMap.URL()),
						kml.Point(
							kml.Coordinates(kml.Coordinate{Lon: -122.0856545755255, Lat: 37.42243077405461}),
						),
					),
				),
			),
			expected: `<kml xmlns="http://www.opengis.net/kml/2.2">` +
				`<Document>` +
				`<name>Highlighted Icon</name>` +
				`<description>Place your mouse over the icon to see it display the new icon</description>` +
				`<Style id="highlightPlacemark">` +
				`<IconStyle>` +
				`<Icon>` +
				`<href>http://maps.google.com/mapfiles/kml/paddle/red-stars.png</href>` +
				`</Icon>` +
				`</IconStyle>` +
				`</Style>` +
				`<Style id="normalPlacemark">` +
				`<IconStyle>` +
				`<Icon>` +
				`<href>http://maps.google.com/mapfiles/kml/paddle/wht-blank.png</href>` +
				`</Icon>` +
				`</IconStyle>` +
				`</Style>` +
				`<StyleMap id="exampleStyleMap">` +
				`<Pair>` +
				`<key>normal</key>` +
				`<styleUrl>#normalPlacemark</styleUrl>` +
				`</Pair>` +
				`<Pair>` +
				`<key>highlight</key>` +
				`<styleUrl>#highlightPlacemark</styleUrl>` +
				`</Pair>` +
				`</StyleMap>` +
				`<Placemark>` +
				`<name>Roll over this icon</name>` +
				`<styleUrl>#exampleStyleMap</styleUrl>` +
				`<Point>` +
				`<coordinates>-122.0856545755255,37.42243077405461</coordinates>` +
				`</Point>` +
				`</Placemark>` +
				`</Document>` +
				`</kml>`,
		},
		{
			name: "trail_head_type",
			element: kml.KML(
				kml.Document(
					kml.NamedSchema("TrailHeadTypeId", "TrailHeadType",
						kml.SimpleField("TrailHeadName", "string",
							kml.DisplayName("<b>Trail Head Name</b>"),
						),
						kml.SimpleField("TrailLength", "double",
							kml.DisplayName("<i>The length in miles</i>"),
						),
						kml.SimpleField("ElevationGain", "int",
							kml.DisplayName("<i>change in altitude</i>"),
						),
					),
				),
			),
			expected: `<kml xmlns="http://www.opengis.net/kml/2.2">` +
				`<Document>` +
				`<Schema id="TrailHeadTypeId" name="TrailHeadType">` +
				`<SimpleField name="TrailHeadName" type="string">` +
				`<displayName>&lt;b&gt;Trail Head Name&lt;/b&gt;</displayName>` +
				`</SimpleField>` +
				`<SimpleField name="TrailLength" type="double">` +
				`<displayName>&lt;i&gt;The length in miles&lt;/i&gt;</displayName>` +
				`</SimpleField>` +
				`<SimpleField name="ElevationGain" type="int">` +
				`<displayName>&lt;i&gt;change in altitude&lt;/i&gt;</displayName>` +
				`</SimpleField>` +
				`</Schema>` +
				`</Document>` +
				`</kml>`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var builder strings.Builder
			encoder := xml.NewEncoder(&builder)
			assert.NoError(t, encoder.Encode(tc.element))
			assert.Equal(t, tc.expected, builder.String())
		})
	}
}

func TestWrite(t *testing.T) {
	for _, tc := range []struct {
		name     string
		element  kml.TopLevelElement
		expected string
	}{
		{
			name:    "placemark",
			element: kml.KML(kml.Placemark()),
			expected: `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
				`<kml xmlns="http://www.opengis.net/kml/2.2">` +
				`<Placemark>` +
				`</Placemark>` +
				`</kml>`,
		},
		{
			name: "simple_placemark",
			element: kml.KML(
				kml.Placemark(
					kml.Name("Simple placemark"),
					kml.Description("Attached to the ground. Intelligently places itself at the height of the underlying terrain."),
					kml.Point(
						kml.Coordinates(kml.Coordinate{Lon: -122.0822035425683, Lat: 37.42228990140251}),
					),
				),
			),
			expected: `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
				`<kml xmlns="http://www.opengis.net/kml/2.2">` +
				`<Placemark>` +
				`<name>Simple placemark</name>` +
				`<description>Attached to the ground. Intelligently places itself at the height of the underlying terrain.</description>` +
				`<Point>` +
				`<coordinates>-122.0822035425683,37.42228990140251</coordinates>` +
				`</Point>` +
				`</Placemark>` +
				`</kml>`,
		},
		{
			name: "entity_references_example",
			element: kml.KML(
				kml.Document(
					kml.Placemark(
						kml.Name("Entity references example"),
						kml.Description(
							`<h1>Entity references are hard to type!</h1>`+
								`<p><font color="red">Text is <i>more readable</i> and `+
								`<b>easier to write</b> when you can avoid using entity `+
								`references.</font></p>`,
						),
						kml.Point(
							kml.Coordinates(kml.Coordinate{Lon: 102.594411, Lat: 14.998518}),
						),
					),
				),
			),
			expected: `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
				`<kml xmlns="http://www.opengis.net/kml/2.2">` +
				`<Document>` +
				`<Placemark>` +
				`<name>Entity references example</name>` +
				`<description>` +
				`&lt;h1&gt;Entity references are hard to type!&lt;/h1&gt;` +
				`&lt;p&gt;&lt;font color=&#34;red&#34;&gt;Text is ` +
				`&lt;i&gt;more readable&lt;/i&gt; ` +
				`and &lt;b&gt;easier to write&lt;/b&gt; ` +
				`when you can avoid using entity references.&lt;/font&gt;&lt;/p&gt;` +
				`</description>` +
				`<Point>` +
				`<coordinates>102.594411,14.998518</coordinates>` +
				`</Point>` +
				`</Placemark>` +
				`</Document>` +
				`</kml>`,
		},
		{
			name: "ground_overlays",
			element: kml.KML(
				kml.Folder(
					kml.Name("Ground Overlays"),
					kml.Description("Examples of ground overlays"),
					kml.GroundOverlay(
						kml.Name("Large-scale overlay on terrain"),
						kml.Description("Overlay shows Mount Etna erupting on July 13th, 2001."),
						kml.Icon(
							kml.Href("http://developers.google.com/kml/documentation/images/etna.jpg"),
						),
						kml.LatLonBox(
							kml.North(37.91904192681665),
							kml.South(37.46543388598137),
							kml.East(15.35832653742206),
							kml.West(14.60128369746704),
							kml.Rotation(-0.1556640799496235),
						),
					),
				),
			),
			expected: `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
				`<kml xmlns="http://www.opengis.net/kml/2.2">` +
				`<Folder>` +
				`<name>Ground Overlays</name>` +
				`<description>Examples of ground overlays</description>` +
				`<GroundOverlay>` +
				`<name>Large-scale overlay on terrain</name>` +
				`<description>Overlay shows Mount Etna erupting on July 13th, 2001.</description>` +
				`<Icon>` +
				`<href>http://developers.google.com/kml/documentation/images/etna.jpg</href>` +
				`</Icon>` +
				`<LatLonBox>` +
				`<north>37.91904192681665</north>` +
				`<south>37.46543388598137</south>` +
				`<east>15.35832653742206</east>` +
				`<west>14.60128369746704</west>` +
				`<rotation>-0.1556640799496235</rotation>` +
				`</LatLonBox>` +
				`</GroundOverlay>` +
				`</Folder>` +
				`</kml>`,
		},
		{
			name: "the_pentagon",
			element: kml.KML(
				kml.Placemark(
					kml.Name("The Pentagon"),
					kml.Polygon(
						kml.Extrude(true),
						kml.AltitudeMode(kml.AltitudeModeRelativeToGround),
						kml.OuterBoundaryIs(
							kml.LinearRing(
								kml.Coordinates([]kml.Coordinate{
									{-77.05788457660967, 38.87253259892824, 100},
									{-77.05465973756702, 38.87291016281703, 100},
									{-77.05315536854791, 38.87053267794386, 100},
									{-77.05552622493516, 38.868757801256, 100},
									{-77.05844056290393, 38.86996206506943, 100},
									{-77.05788457660967, 38.87253259892824, 100},
								}...),
							),
						),
						kml.InnerBoundaryIs(
							kml.LinearRing(
								kml.Coordinates([]kml.Coordinate{
									{-77.05668055019126, 38.87154239798456, 100},
									{-77.05542625960818, 38.87167890344077, 100},
									{-77.05485125901024, 38.87076535397792, 100},
									{-77.05577677433152, 38.87008686581446, 100},
									{-77.05691162017543, 38.87054446963351, 100},
									{-77.05668055019126, 38.87154239798456, 100},
								}...),
							),
						),
					),
				),
			),
			expected: `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
				`<kml xmlns="http://www.opengis.net/kml/2.2">` +
				`<Placemark>` +
				`<name>The Pentagon</name>` +
				`<Polygon>` +
				`<extrude>1</extrude>` +
				`<altitudeMode>relativeToGround</altitudeMode>` +
				`<outerBoundaryIs>` +
				`<LinearRing>` +
				`<coordinates>` +
				`-77.05788457660967,38.87253259892824,100 ` +
				`-77.05465973756702,38.87291016281703,100 ` +
				`-77.0531553685479,38.87053267794386,100 ` +
				`-77.05552622493516,38.868757801256,100 ` +
				`-77.05844056290393,38.86996206506943,100 ` +
				`-77.05788457660967,38.87253259892824,100` +
				`</coordinates>` +
				`</LinearRing>` +
				`</outerBoundaryIs>` +
				`<innerBoundaryIs>` +
				`<LinearRing>` +
				`<coordinates>` +
				`-77.05668055019126,38.87154239798456,100 ` +
				`-77.05542625960818,38.87167890344077,100 ` +
				`-77.05485125901023,38.87076535397792,100 ` +
				`-77.05577677433152,38.87008686581446,100 ` +
				`-77.05691162017543,38.87054446963351,100 ` +
				`-77.05668055019126,38.87154239798456,100` +
				`</coordinates>` +
				`</LinearRing>` +
				`</innerBoundaryIs>` +
				`</Polygon>` +
				`</Placemark>` +
				`</kml>`,
		},
		{
			name:    "gx_placemark",
			element: kml.GxKML(kml.Placemark()),
			expected: `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
				`<kml xmlns="http://www.opengis.net/kml/2.2" xmlns:gx="http://www.google.com/kml/ext/2.2">` +
				`<Placemark>` +
				`</Placemark>` +
				`</kml>`,
		},
		{
			name: "gx_track",
			element: kml.GxKML(
				kml.Folder(
					kml.Placemark(
						kml.GxTrack(
							kml.When(time.Date(2010, 5, 28, 2, 2, 9, 0, time.UTC)),
							kml.When(time.Date(2010, 5, 28, 2, 2, 35, 0, time.UTC)),
							kml.When(time.Date(2010, 5, 28, 2, 2, 44, 0, time.UTC)),
							kml.When(time.Date(2010, 5, 28, 2, 2, 53, 0, time.UTC)),
							kml.When(time.Date(2010, 5, 28, 2, 2, 54, 0, time.UTC)),
							kml.When(time.Date(2010, 5, 28, 2, 2, 55, 0, time.UTC)),
							kml.When(time.Date(2010, 5, 28, 2, 2, 56, 0, time.UTC)),
							kml.GxCoord(kml.Coordinate{-122.207881, 37.371915, 156.000000}),
							kml.GxCoord(kml.Coordinate{-122.205712, 37.373288, 152.000000}),
							kml.GxCoord(kml.Coordinate{-122.204678, 37.373939, 147.000000}),
							kml.GxCoord(kml.Coordinate{-122.203572, 37.374630, 142.199997}),
							kml.GxCoord(kml.Coordinate{-122.203451, 37.374706, 141.800003}),
							kml.GxCoord(kml.Coordinate{-122.203329, 37.374780, 141.199997}),
							kml.GxCoord(kml.Coordinate{-122.203207, 37.374857, 140.199997}),
						),
					),
				),
			),
			expected: `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
				`<kml xmlns="http://www.opengis.net/kml/2.2" xmlns:gx="http://www.google.com/kml/ext/2.2">` +
				`<Folder>` +
				`<Placemark>` +
				`<gx:Track>` +
				`<when>2010-05-28T02:02:09Z</when>` +
				`<when>2010-05-28T02:02:35Z</when>` +
				`<when>2010-05-28T02:02:44Z</when>` +
				`<when>2010-05-28T02:02:53Z</when>` +
				`<when>2010-05-28T02:02:54Z</when>` +
				`<when>2010-05-28T02:02:55Z</when>` +
				`<when>2010-05-28T02:02:56Z</when>` +
				`<gx:coord>-122.207881 37.371915 156</gx:coord>` +
				`<gx:coord>-122.205712 37.373288 152</gx:coord>` +
				`<gx:coord>-122.204678 37.373939 147</gx:coord>` +
				`<gx:coord>-122.203572 37.37463 142.199997</gx:coord>` +
				`<gx:coord>-122.203451 37.374706 141.800003</gx:coord>` +
				`<gx:coord>-122.203329 37.37478 141.199997</gx:coord>` +
				`<gx:coord>-122.203207 37.374857 140.199997</gx:coord>` +
				`</gx:Track>` +
				`</Placemark>` +
				`</Folder>` +
				`</kml>`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var builder strings.Builder
			assert.NoError(t, tc.element.Write(&builder))
			assert.Equal(t, tc.expected, builder.String())
		})
	}
}
