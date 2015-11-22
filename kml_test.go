package kml

import (
	"bytes"
	"image/color"
	"testing"
	"time"
)

type testCase struct {
	e    Element
	want string
}

func (tc testCase) testWrite(t *testing.T) {
	b := &bytes.Buffer{}
	if err := tc.e.Write(b); err != nil {
		t.Errorf("%#v.Write(b) == %#v, want nil", tc.e, err)
		return
	}
	if got := b.String(); got != tc.want {
		t.Errorf("%#v.Write(b) wrote %#v, want %#v", tc.e, got, tc.want)
	}
}

func (tc testCase) testStringXML(t *testing.T) {
	got, err := tc.e.StringXML()
	if err != nil {
		t.Errorf("%#v.StringXML() == %#v, %#v, want ..., nil", tc.e, got, err)
		return
	}
	if got != tc.want {
		t.Errorf("%#v.StringXML() == ..., nil,\n got %#v, nil,\nwant %#v, nil", tc.e, got, tc.want)
	}
}

func TestSimpleElements(t *testing.T) {
	for _, tc := range []testCase{
		{
			Altitude(0),
			`<altitude>0</altitude>`,
		},
		{
			AltitudeMode("absolute"),
			`<altitudeMode>absolute</altitudeMode>`,
		},
		{
			Begin(time.Date(1876, 8, 1, 0, 0, 0, 0, time.UTC)),
			`<begin>1876-08-01T00:00:00Z</begin>`,
		},
		{
			BgColor(color.Black),
			`<bgColor>ff000000</bgColor>`,
		},
		{
			Color(color.White),
			`<color>ffffffff</color>`,
		},
		{
			Coordinates(Coordinate{Lon: 1.23, Lat: 4.56, Alt: 7.89}),
			`<coordinates>1.23,4.56,7.89</coordinates>`,
		},
		{
			Description("text"),
			`<description>text</description>`,
		},
		{
			End(time.Date(2015, 12, 31, 23, 59, 59, 0, time.UTC)),
			`<end>2015-12-31T23:59:59Z</end>`,
		},
		{
			Extrude(false),
			`<extrude>0</extrude>`,
		},
		{
			Folder(),
			`<Folder></Folder>`,
		},
		{
			GxCoord(Coordinate{1.23, 4.56, 7.89}),
			`<gx:coord>1.23 4.56 7.89</gx:coord>`,
		},
		{
			Heading(0),
			`<heading>0</heading>`,
		},
		{
			HotSpot(Vec2{X: 0.5, Y: 0.5, XUnits: "pixels", YUnits: "pixels"}),
			`<hotSpot x="0.5" y="0.5" xunits="pixels" yunits="pixels"></hotSpot>`,
		},
		{
			HrefMustParse("https://www.google.com/"),
			`<href>https://www.google.com/</href>`,
		},
		{
			Latitude(0),
			`<latitude>0</latitude>`,
		},
		{
			ListItemType("check"),
			`<listItemType>check</listItemType>`,
		},
		{
			OverlayXY(Vec2{X: 0, Y: 0, XUnits: "fraction", YUnits: "fraction"}),
			`<overlayXY x="0" y="0" xunits="fraction" yunits="fraction"></overlayXY>`,
		},
		// FIXME More simple elements
	} {
		tc.testStringXML(t)
	}
}

func TestSharedStyles(t *testing.T) {
	style0 := Style("0")
	highlightPlacemarkStyle := Style(
		"highlightPlacemark",
		IconStyle(
			Icon(
				HrefMustParse("http://maps.google.com/mapfiles/kml/paddle/red-stars.png"),
			),
		),
	)
	normalPlacemarkStyle := Style(
		"normalPlacemark",
		IconStyle(
			Icon(
				HrefMustParse("http://maps.google.com/mapfiles/kml/paddle/wht-blank.png"),
			),
		),
	)
	exampleStyleMap := StyleMap(
		"exampleStyleMap",
		Pair(
			Key("normal"),
			StyleURL(normalPlacemarkStyle),
		),
		Pair(
			Key("highlight"),
			StyleURL(highlightPlacemarkStyle),
		),
	)
	for _, tc := range []testCase{
		{
			e: Folder(
				style0,
				Placemark(
					StyleURL(style0),
				),
			),
			want: `<Folder>` +
				`<Style id="0">` +
				`</Style>` +
				`<Placemark>` +
				`<styleUrl>#0</styleUrl>` +
				`</Placemark>` +
				`</Folder>`,
		},
		{
			e: KML(
				Document(
					Name("Highlighted Icon"),
					Description("Place your mouse over the icon to see it display the new icon"),
					highlightPlacemarkStyle,
					normalPlacemarkStyle,
					exampleStyleMap,
					Placemark(
						Name("Roll over this icon"),
						StyleURL(exampleStyleMap),
						Point(
							Coordinates(Coordinate{Lon: -122.0856545755255, Lat: 37.42243077405461}),
						),
					),
				),
			),
			want: `<kml xmlns="http://www.opengis.net/kml/2.2">` +
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
			e: KML(
				Document(
					Schema("TrailHeadTypeId", "TrailHeadType",
						SimpleField("TrailHeadName", "string",
							DisplayName("<b>Trail Head Name</b>"),
						),
						SimpleField("TrailLength", "double",
							DisplayName("<i>The length in miles</i>"),
						),
						SimpleField("ElevationGain", "int",
							DisplayName("<i>change in altitude</i>"),
						),
					),
				),
			),
			want: `<kml xmlns="http://www.opengis.net/kml/2.2">` +
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
		tc.testStringXML(t)
	}
}

func TestWrite(t *testing.T) {
	for _, tc := range []testCase{
		{
			e: KML(),
			want: `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
				`<kml xmlns="http://www.opengis.net/kml/2.2"></kml>`,
		},
		{
			e: KML(
				Placemark(
					Name("Simple placemark"),
					Description("Attached to the ground. Intelligently places itself at the height of the underlying terrain."),
					Point(
						Coordinates(Coordinate{Lon: -122.0822035425683, Lat: 37.42228990140251}),
					),
				),
			),
			want: `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
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
			e: KML(
				Document(
					Placemark(
						Name("Entity references example"),
						Description(
							`<h1>Entity references are hard to type!</h1>`+
								`<p><font color="red">Text is <i>more readable</i> and `+
								`<b>easier to write</b> when you can avoid using entity `+
								`references.</font></p>`,
						),
						Point(
							Coordinates(Coordinate{Lon: 102.594411, Lat: 14.998518}),
						),
					),
				),
			),
			want: `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
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
			e: KML(
				Folder(
					Name("Ground Overlays"),
					Description("Examples of ground overlays"),
					GroundOverlay(
						Name("Large-scale overlay on terrain"),
						Description("Overlay shows Mount Etna erupting on July 13th, 2001."),
						Icon(
							HrefMustParse("http://developers.google.com/kml/documentation/images/etna.jpg"),
						),
						LatLonBox(
							North(37.91904192681665),
							South(37.46543388598137),
							East(15.35832653742206),
							West(14.60128369746704),
							Rotation(-0.1556640799496235),
						),
					),
				),
			),
			want: `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
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
			e: KML(
				Placemark(
					Name("The Pentagon"),
					Polygon(
						Extrude(true),
						AltitudeMode("relativeToGround"),
						OuterBoundaryIs(
							LinearRing(
								Coordinates([]Coordinate{
									{-77.05788457660967, 38.87253259892824, 100},
									{-77.05465973756702, 38.87291016281703, 100},
									{-77.05315536854791, 38.87053267794386, 100},
									{-77.05552622493516, 38.868757801256, 100},
									{-77.05844056290393, 38.86996206506943, 100},
									{-77.05788457660967, 38.87253259892824, 100},
								}...),
							),
						),
						InnerBoundaryIs(
							LinearRing(
								Coordinates([]Coordinate{
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
			want: `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
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
			e: GxKML(),
			want: `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
				`<kml xmlns="http://www.opengis.net/kml/2.2" xmlns:gx="http://www.google.com/kml/ext/2.2"></kml>`,
		},
		{
			e: GxKML(
				Folder(
					Placemark(
						GxTrack(
							When(time.Date(2010, 5, 28, 2, 2, 9, 0, time.UTC)),
							When(time.Date(2010, 5, 28, 2, 2, 35, 0, time.UTC)),
							When(time.Date(2010, 5, 28, 2, 2, 44, 0, time.UTC)),
							When(time.Date(2010, 5, 28, 2, 2, 53, 0, time.UTC)),
							When(time.Date(2010, 5, 28, 2, 2, 54, 0, time.UTC)),
							When(time.Date(2010, 5, 28, 2, 2, 55, 0, time.UTC)),
							When(time.Date(2010, 5, 28, 2, 2, 56, 0, time.UTC)),
							GxCoord(Coordinate{-122.207881, 37.371915, 156.000000}),
							GxCoord(Coordinate{-122.205712, 37.373288, 152.000000}),
							GxCoord(Coordinate{-122.204678, 37.373939, 147.000000}),
							GxCoord(Coordinate{-122.203572, 37.374630, 142.199997}),
							GxCoord(Coordinate{-122.203451, 37.374706, 141.800003}),
							GxCoord(Coordinate{-122.203329, 37.374780, 141.199997}),
							GxCoord(Coordinate{-122.203207, 37.374857, 140.199997}),
						),
					),
				),
			),
			want: `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
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
		tc.testWrite(t)
	}
}
