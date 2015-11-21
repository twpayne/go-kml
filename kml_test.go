package kml

import (
	"bytes"
	"image/color"
	"testing"
	"time"
)

func TestTutorial(t *testing.T) {
	for _, tc := range []struct {
		e    Element
		want string
	}{
		{
			e: KML(
				Placemark(
					Name("Simple placemark"),
					Description("Attached to the ground. Intelligently places itself at the height of the underlying terrain."),
					Point(
						Coordinates(Coordinate{Lon: -122.0822035425683, Lat: 37.42228990140251, Alt: 0}),
					),
				),
			),
			want: `<kml xmlns="http://www.opengis.net/kml/2.2">` +
				`<Placemark>` +
				`<name>Simple placemark</name>` +
				`<description>Attached to the ground. Intelligently places itself at the height of the underlying terrain.</description>` +
				`<Point>` +
				`<coordinates>-122.0822035425683,37.42228990140251,0</coordinates>` +
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
			want: `<kml xmlns="http://www.opengis.net/kml/2.2">` +
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
				`<coordinates>102.594411,14.998518,0</coordinates>` +
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
			want: `<kml xmlns="http://www.opengis.net/kml/2.2">` +
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
	} {
		if got, err := tc.e.StringXML(); err != nil || got != tc.want {
			t.Errorf("%#v.StringXML() == %#v, %#v, want %#v, nil", tc.e, got, err, tc.want)
		}
	}
}

func TestSimpleElements(t *testing.T) {
	for _, tc := range []struct {
		e    Element
		want string
	}{
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
		if got, err := tc.e.StringXML(); err != nil || got != tc.want {
			t.Errorf("%#v.StringXML() == %#v, %#v, want %#v, nil", tc.e, got, err, tc.want)
		}
	}
}

func TestWrite(t *testing.T) {
	for _, tc := range []struct {
		e    Element
		want string
	}{
		{
			KML(),
			`<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
				`<kml xmlns="http://www.opengis.net/kml/2.2"></kml>`,
		},
	} {
		b := &bytes.Buffer{}
		if err := tc.e.Write(b); err != nil {
			t.Errorf("%#v.Write(b) == %#v, want nil", tc.e, err)
			continue
		}
		if got := b.String(); got != tc.want {
			t.Errorf("%#v.Write(b) wrote %#v, want %#v", tc.e, got, tc.want)
		}
	}
}
