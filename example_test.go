package kml_test

import (
	"image/color"
	"log"
	"os"

	"github.com/twpayne/go-kml/v3"
)

func ExamplePlacemark() {
	k := kml.KML(
		kml.Placemark(
			kml.Name("Simple placemark"),
			kml.Description("Attached to the ground. Intelligently places itself at the height of the underlying terrain."),
			kml.Point(
				kml.Coordinates(kml.Coordinate{Lon: -122.0822035425683, Lat: 37.42228990140251}),
			),
		),
	)
	if err := k.WriteIndent(os.Stdout, "", "  "); err != nil {
		log.Fatal(err)
	}
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <kml xmlns="http://www.opengis.net/kml/2.2">
	//   <Placemark>
	//     <name>Simple placemark</name>
	//     <description>Attached to the ground. Intelligently places itself at the height of the underlying terrain.</description>
	//     <Point>
	//       <coordinates>-122.0822035425683,37.42228990140251</coordinates>
	//     </Point>
	//   </Placemark>
	// </kml>
}

func ExampleDescription() {
	k := kml.KML(
		kml.Document(
			kml.Placemark(
				kml.Name("CDATA example"),
				kml.Description(`<h1>CDATA Tags are useful!</h1> <p><font color="red">Text is <i>more readable</i> and <b>easier to write</b> when you can avoid using entity references.</font></p>`),
				kml.Point(
					kml.Coordinates(kml.Coordinate{Lon: 102.595626, Lat: 14.996729}),
				),
			),
		),
	)
	if err := k.WriteIndent(os.Stdout, "", "  "); err != nil {
		log.Fatal(err)
	}
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <kml xmlns="http://www.opengis.net/kml/2.2">
	//   <Document>
	//     <Placemark>
	//       <name>CDATA example</name>
	//       <description>&lt;h1&gt;CDATA Tags are useful!&lt;/h1&gt; &lt;p&gt;&lt;font color=&#34;red&#34;&gt;Text is &lt;i&gt;more readable&lt;/i&gt; and &lt;b&gt;easier to write&lt;/b&gt; when you can avoid using entity references.&lt;/font&gt;&lt;/p&gt;</description>
	//       <Point>
	//         <coordinates>102.595626,14.996729</coordinates>
	//       </Point>
	//     </Placemark>
	//   </Document>
	// </kml>
}

func ExampleGroundOverlay() {
	k := kml.KML(
		kml.Folder(
			kml.Name("Ground Overlays"),
			kml.Description("Examples of ground overlays"),
			kml.GroundOverlay(
				kml.Name("Large-scale overlay on terrain"),
				kml.Description("Overlay shows Mount Etna erupting on July 13th, 2001."),
				kml.Icon(
					kml.Href("https://developers.google.com/kml/documentation/images/etna.jpg"),
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
	)
	if err := k.WriteIndent(os.Stdout, "", "  "); err != nil {
		log.Fatal(err)
	}
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <kml xmlns="http://www.opengis.net/kml/2.2">
	//   <Folder>
	//     <name>Ground Overlays</name>
	//     <description>Examples of ground overlays</description>
	//     <GroundOverlay>
	//       <name>Large-scale overlay on terrain</name>
	//       <description>Overlay shows Mount Etna erupting on July 13th, 2001.</description>
	//       <Icon>
	//         <href>https://developers.google.com/kml/documentation/images/etna.jpg</href>
	//       </Icon>
	//       <LatLonBox>
	//         <north>37.91904192681665</north>
	//         <south>37.46543388598137</south>
	//         <east>15.35832653742206</east>
	//         <west>14.60128369746704</west>
	//         <rotation>-0.1556640799496235</rotation>
	//       </LatLonBox>
	//     </GroundOverlay>
	//   </Folder>
	// </kml>
}

func ExampleLineString() {
	k := kml.KML(
		kml.Document(
			kml.Name("Paths"),
			kml.Description("Examples of paths. Note that the tessellate tag is by default set to 0. If you want to create tessellated lines, they must be authored (or edited) directly in KML."),
			kml.SharedStyle(
				"yellowLineGreenPoly",
				kml.LineStyle(
					kml.Color(color.RGBA{R: 255, G: 255, B: 0, A: 127}),
					kml.Width(4),
				),
				kml.PolyStyle(
					kml.Color(color.RGBA{R: 0, G: 255, B: 0, A: 127}),
				),
			),
			kml.Placemark(
				kml.Name("Absolute Extruded"),
				kml.Description("Transparent green wall with yellow outlines"),
				kml.StyleURL("#yellowLineGreenPoly"),
				kml.LineString(
					kml.Extrude(true),
					kml.Tessellate(true),
					kml.AltitudeMode(kml.AltitudeModeAbsolute),
					kml.Coordinates([]kml.Coordinate{
						{Lon: -112.2550785337791, Lat: 36.07954952145647, Alt: 2357},
						{Lon: -112.2549277039738, Lat: 36.08117083492122, Alt: 2357},
						{Lon: -112.2552505069063, Lat: 36.08260761307279, Alt: 2357},
						{Lon: -112.2564540158376, Lat: 36.08395660588506, Alt: 2357},
						{Lon: -112.2580238976449, Lat: 36.08511401044813, Alt: 2357},
						{Lon: -112.2595218489022, Lat: 36.08584355239394, Alt: 2357},
						{Lon: -112.2608216347552, Lat: 36.08612634548589, Alt: 2357},
						{Lon: -112.262073428656, Lat: 36.08626019085147, Alt: 2357},
						{Lon: -112.2633204928495, Lat: 36.08621519860091, Alt: 2357},
						{Lon: -112.2644963846444, Lat: 36.08627897945274, Alt: 2357},
						{Lon: -112.2656969554589, Lat: 36.08649599090644, Alt: 2357},
					}...),
				),
			),
		),
	)
	if err := k.WriteIndent(os.Stdout, "", "  "); err != nil {
		log.Fatal(err)
	}
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <kml xmlns="http://www.opengis.net/kml/2.2">
	//   <Document>
	//     <name>Paths</name>
	//     <description>Examples of paths. Note that the tessellate tag is by default set to 0. If you want to create tessellated lines, they must be authored (or edited) directly in KML.</description>
	//     <Style id="yellowLineGreenPoly">
	//       <LineStyle>
	//         <color>7f00ffff</color>
	//         <width>4</width>
	//       </LineStyle>
	//       <PolyStyle>
	//         <color>7f00ff00</color>
	//       </PolyStyle>
	//     </Style>
	//     <Placemark>
	//       <name>Absolute Extruded</name>
	//       <description>Transparent green wall with yellow outlines</description>
	//       <styleUrl>#yellowLineGreenPoly</styleUrl>
	//       <LineString>
	//         <extrude>1</extrude>
	//         <tessellate>1</tessellate>
	//         <altitudeMode>absolute</altitudeMode>
	//         <coordinates>-112.2550785337791,36.07954952145647,2357 -112.2549277039738,36.08117083492122,2357 -112.2552505069063,36.08260761307279,2357 -112.2564540158376,36.08395660588506,2357 -112.2580238976449,36.08511401044813,2357 -112.2595218489022,36.08584355239394,2357 -112.2608216347552,36.08612634548589,2357 -112.262073428656,36.08626019085147,2357 -112.2633204928495,36.08621519860091,2357 -112.2644963846444,36.08627897945274,2357 -112.2656969554589,36.08649599090644,2357</coordinates>
	//       </LineString>
	//     </Placemark>
	//   </Document>
	// </kml>
}

func ExamplePolygon() {
	k := kml.KML(
		kml.Placemark(
			kml.Name("The Pentagon"),
			kml.Polygon(
				kml.Extrude(true),
				kml.AltitudeMode(kml.AltitudeModeRelativeToGround),
				kml.OuterBoundaryIs(
					kml.LinearRing(
						kml.Coordinates([]kml.Coordinate{
							{Lon: -77.05788457660967, Lat: 38.87253259892824, Alt: 100},
							{Lon: -77.05465973756702, Lat: 38.87291016281703, Alt: 100},
							{Lon: -77.0531553685479, Lat: 38.87053267794386, Alt: 100},
							{Lon: -77.05552622493516, Lat: 38.868757801256, Alt: 100},
							{Lon: -77.05844056290393, Lat: 38.86996206506943, Alt: 100},
							{Lon: -77.05788457660967, Lat: 38.87253259892824, Alt: 100},
						}...),
					),
				),
				kml.InnerBoundaryIs(
					kml.LinearRing(
						kml.Coordinates([]kml.Coordinate{
							{Lon: -77.05668055019126, Lat: 38.87154239798456, Alt: 100},
							{Lon: -77.05542625960818, Lat: 38.87167890344077, Alt: 100},
							{Lon: -77.05485125901023, Lat: 38.87076535397792, Alt: 100},
							{Lon: -77.05577677433152, Lat: 38.87008686581446, Alt: 100},
							{Lon: -77.05691162017543, Lat: 38.87054446963351, Alt: 100},
							{Lon: -77.05668055019126, Lat: 38.87154239798456, Alt: 100},
						}...),
					),
				),
			),
		),
	)
	if err := k.WriteIndent(os.Stdout, "", "  "); err != nil {
		log.Fatal(err)
	}
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <kml xmlns="http://www.opengis.net/kml/2.2">
	//   <Placemark>
	//     <name>The Pentagon</name>
	//     <Polygon>
	//       <extrude>1</extrude>
	//       <altitudeMode>relativeToGround</altitudeMode>
	//       <outerBoundaryIs>
	//         <LinearRing>
	//           <coordinates>-77.05788457660967,38.87253259892824,100 -77.05465973756702,38.87291016281703,100 -77.0531553685479,38.87053267794386,100 -77.05552622493516,38.868757801256,100 -77.05844056290393,38.86996206506943,100 -77.05788457660967,38.87253259892824,100</coordinates>
	//         </LinearRing>
	//       </outerBoundaryIs>
	//       <innerBoundaryIs>
	//         <LinearRing>
	//           <coordinates>-77.05668055019126,38.87154239798456,100 -77.05542625960818,38.87167890344077,100 -77.05485125901023,38.87076535397792,100 -77.05577677433152,38.87008686581446,100 -77.05691162017543,38.87054446963351,100 -77.05668055019126,38.87154239798456,100</coordinates>
	//         </LinearRing>
	//       </innerBoundaryIs>
	//     </Polygon>
	//   </Placemark>
	// </kml>
}

func ExampleStyle() {
	k := kml.KML(
		kml.Document(
			kml.SharedStyle(
				"transBluePoly",
				kml.LineStyle(
					kml.Width(1.5),
				),
				kml.PolyStyle(
					kml.Color(color.RGBA{R: 0, G: 0, B: 255, A: 125}),
				),
			),
			kml.Placemark(
				kml.Name("Building 41"),
				kml.StyleURL("#transBluePoly"),
				kml.Polygon(
					kml.Extrude(true),
					kml.AltitudeMode(kml.AltitudeModeRelativeToGround),
					kml.OuterBoundaryIs(
						kml.LinearRing(
							kml.Coordinates([]kml.Coordinate{
								{Lon: -122.0857412771483, Lat: 37.42227033155257, Alt: 17},
								{Lon: -122.0858169768481, Lat: 37.42231408832346, Alt: 17},
								{Lon: -122.085852582875, Lat: 37.42230337469744, Alt: 17},
								{Lon: -122.0858799945639, Lat: 37.42225686138789, Alt: 17},
								{Lon: -122.0858860101409, Lat: 37.4222311076138, Alt: 17},
								{Lon: -122.0858069157288, Lat: 37.42220250173855, Alt: 17},
								{Lon: -122.0858379542653, Lat: 37.42214027058678, Alt: 17},
								{Lon: -122.0856732640519, Lat: 37.42208690214408, Alt: 17},
								{Lon: -122.0856022926407, Lat: 37.42214885429042, Alt: 17},
								{Lon: -122.0855902778436, Lat: 37.422128290487, Alt: 17},
								{Lon: -122.0855841672237, Lat: 37.42208171967246, Alt: 17},
								{Lon: -122.0854852065741, Lat: 37.42210455874995, Alt: 17},
								{Lon: -122.0855067264352, Lat: 37.42214267949824, Alt: 17},
								{Lon: -122.0854430712915, Lat: 37.42212783846172, Alt: 17},
								{Lon: -122.0850990714904, Lat: 37.42251282407603, Alt: 17},
								{Lon: -122.0856769818632, Lat: 37.42281815323651, Alt: 17},
								{Lon: -122.0860162273783, Lat: 37.42244918858722, Alt: 17},
								{Lon: -122.0857260327004, Lat: 37.42229239604253, Alt: 17},
								{Lon: -122.0857412771483, Lat: 37.42227033155257, Alt: 17},
							}...),
						),
					),
				),
			),
		),
	)
	if err := k.WriteIndent(os.Stdout, "", "  "); err != nil {
		log.Fatal(err)
	}
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <kml xmlns="http://www.opengis.net/kml/2.2">
	//   <Document>
	//     <Style id="transBluePoly">
	//       <LineStyle>
	//         <width>1.5</width>
	//       </LineStyle>
	//       <PolyStyle>
	//         <color>7dff0000</color>
	//       </PolyStyle>
	//     </Style>
	//     <Placemark>
	//       <name>Building 41</name>
	//       <styleUrl>#transBluePoly</styleUrl>
	//       <Polygon>
	//         <extrude>1</extrude>
	//         <altitudeMode>relativeToGround</altitudeMode>
	//         <outerBoundaryIs>
	//           <LinearRing>
	//             <coordinates>-122.0857412771483,37.42227033155257,17 -122.0858169768481,37.42231408832346,17 -122.085852582875,37.42230337469744,17 -122.0858799945639,37.42225686138789,17 -122.0858860101409,37.4222311076138,17 -122.0858069157288,37.42220250173855,17 -122.0858379542653,37.42214027058678,17 -122.0856732640519,37.42208690214408,17 -122.0856022926407,37.42214885429042,17 -122.0855902778436,37.422128290487,17 -122.0855841672237,37.42208171967246,17 -122.0854852065741,37.42210455874995,17 -122.0855067264352,37.42214267949824,17 -122.0854430712915,37.42212783846172,17 -122.0850990714904,37.42251282407603,17 -122.0856769818632,37.42281815323651,17 -122.0860162273783,37.42244918858722,17 -122.0857260327004,37.42229239604253,17 -122.0857412771483,37.42227033155257,17</coordinates>
	//           </LinearRing>
	//         </outerBoundaryIs>
	//       </Polygon>
	//     </Placemark>
	//   </Document>
	// </kml>
}

func ExampleSharedStyleMap() {
	k := kml.KML(
		kml.Document(
			kml.Name("Highlighted Icon"),
			kml.Description("Place your mouse over the icon to see it display the new icon"),
			kml.SharedStyle(
				"highlightPlacemark",
				kml.IconStyle(
					kml.Icon(
						kml.Href("http://maps.google.com/mapfiles/kml/paddle/red-stars.png"),
					),
				),
			),
			kml.SharedStyle(
				"normalPlacemark",
				kml.IconStyle(
					kml.Icon(
						kml.Href("http://maps.google.com/mapfiles/kml/paddle/wht-blank.png"),
					),
				),
			),
			kml.SharedStyleMap(
				"exampleStyleMap",
				kml.Pair(
					kml.Key(kml.StyleStateNormal),
					kml.StyleURL("#normalPlacemark"),
				),
				kml.Pair(
					kml.Key(kml.StyleStateHighlight),
					kml.StyleURL("#highlightPlacemark"),
				),
			),
			kml.Placemark(
				kml.Name("Roll over this icon"),
				kml.StyleURL("#exampleStyleMap"),
				kml.Point(
					kml.Coordinates(kml.Coordinate{Lon: -122.0856545755255, Lat: 37.42243077405461}),
				),
			),
		),
	)
	if err := k.WriteIndent(os.Stdout, "", "  "); err != nil {
		log.Fatal(err)
	}
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <kml xmlns="http://www.opengis.net/kml/2.2">
	//   <Document>
	//     <name>Highlighted Icon</name>
	//     <description>Place your mouse over the icon to see it display the new icon</description>
	//     <Style id="highlightPlacemark">
	//       <IconStyle>
	//         <Icon>
	//           <href>http://maps.google.com/mapfiles/kml/paddle/red-stars.png</href>
	//         </Icon>
	//       </IconStyle>
	//     </Style>
	//     <Style id="normalPlacemark">
	//       <IconStyle>
	//         <Icon>
	//           <href>http://maps.google.com/mapfiles/kml/paddle/wht-blank.png</href>
	//         </Icon>
	//       </IconStyle>
	//     </Style>
	//     <StyleMap id="exampleStyleMap">
	//       <Pair>
	//         <key>normal</key>
	//         <styleUrl>#normalPlacemark</styleUrl>
	//       </Pair>
	//       <Pair>
	//         <key>highlight</key>
	//         <styleUrl>#highlightPlacemark</styleUrl>
	//       </Pair>
	//     </StyleMap>
	//     <Placemark>
	//       <name>Roll over this icon</name>
	//       <styleUrl>#exampleStyleMap</styleUrl>
	//       <Point>
	//         <coordinates>-122.0856545755255,37.42243077405461</coordinates>
	//       </Point>
	//     </Placemark>
	//   </Document>
	// </kml>
}

func ExampleScreenOverlay() {
	k := kml.KML(
		kml.ScreenOverlay(
			kml.Name("Absolute Positioning: Top left"),
			kml.Icon(
				kml.Href("http://developers.google.com/kml/documentation/images/top_left.jpg"),
			),
			kml.OverlayXY(kml.Vec2{X: 0, Y: 1, XUnits: kml.UnitsFraction, YUnits: kml.UnitsFraction}),
			kml.ScreenXY(kml.Vec2{X: 0, Y: 1, XUnits: kml.UnitsFraction, YUnits: kml.UnitsFraction}),
			kml.RotationXY(kml.Vec2{X: 0, Y: 0, XUnits: kml.UnitsFraction, YUnits: kml.UnitsFraction}),
			kml.Size(kml.Vec2{X: 0, Y: 0, XUnits: kml.UnitsFraction, YUnits: kml.UnitsFraction}),
		),
	)
	if err := k.WriteIndent(os.Stdout, "", "  "); err != nil {
		log.Fatal(err)
	}
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <kml xmlns="http://www.opengis.net/kml/2.2">
	//   <ScreenOverlay>
	//     <name>Absolute Positioning: Top left</name>
	//     <Icon>
	//       <href>http://developers.google.com/kml/documentation/images/top_left.jpg</href>
	//     </Icon>
	//     <overlayXY x="0" y="1" xunits="fraction" yunits="fraction"></overlayXY>
	//     <screenXY x="0" y="1" xunits="fraction" yunits="fraction"></screenXY>
	//     <rotationXY x="0" y="0" xunits="fraction" yunits="fraction"></rotationXY>
	//     <size x="0" y="0" xunits="fraction" yunits="fraction"></size>
	//   </ScreenOverlay>
	// </kml>
}

func ExampleNetworkLink() {
	k := kml.KML(
		kml.Folder(
			kml.Name("Network Links"),
			kml.Visibility(false),
			kml.Open(false),
			kml.Description("Network link example 1"),
			kml.NetworkLink(
				kml.Name("Random Placemark"),
				kml.Visibility(false),
				kml.Open(false),
				kml.Description("A simple server-side script that generates a new random placemark on each call"),
				kml.RefreshVisibility(false),
				kml.FlyToView(false),
				kml.Link(
					kml.Href("http://yourserver.com/cgi-bin/randomPlacemark.py"),
				),
			),
		),
	)
	if err := k.WriteIndent(os.Stdout, "", "  "); err != nil {
		log.Fatal(err)
	}
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <kml xmlns="http://www.opengis.net/kml/2.2">
	//   <Folder>
	//     <name>Network Links</name>
	//     <visibility>0</visibility>
	//     <open>0</open>
	//     <description>Network link example 1</description>
	//     <NetworkLink>
	//       <name>Random Placemark</name>
	//       <visibility>0</visibility>
	//       <open>0</open>
	//       <description>A simple server-side script that generates a new random placemark on each call</description>
	//       <refreshVisibility>0</refreshVisibility>
	//       <flyToView>0</flyToView>
	//       <Link>
	//         <href>http://yourserver.com/cgi-bin/randomPlacemark.py</href>
	//       </Link>
	//     </NetworkLink>
	//   </Folder>
	// </kml>
}
