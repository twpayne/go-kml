package kml

import (
	"image/color"
	"log"
	"os"
)

func ExamplePlacemark() {
	k := KML(
		Placemark(
			Name("Simple placemark"),
			Description("Attached to the ground. Intelligently places itself at the height of the underlying terrain."),
			Point(
				Coordinates(Coordinate{Lon: -122.0822035425683, Lat: 37.42228990140251}),
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
	k := KML(
		Document(
			Placemark(
				Name("CDATA example"),
				Description(`<h1>CDATA Tags are useful!</h1> <p><font color="red">Text is <i>more readable</i> and <b>easier to write</b> when you can avoid using entity references.</font></p>`),
				Point(
					Coordinates(Coordinate{Lon: 102.595626, Lat: 14.996729}),
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
	k := KML(
		Folder(
			Name("Ground Overlays"),
			Description("Examples of ground overlays"),
			GroundOverlay(
				Name("Large-scale overlay on terrain"),
				Description("Overlay shows Mount Etna erupting on July 13th, 2001."),
				Icon(
					Href("https://developers.google.com/kml/documentation/images/etna.jpg"),
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
	k := KML(
		Document(
			Name("Paths"),
			Description("Examples of paths. Note that the tessellate tag is by default set to 0. If you want to create tessellated lines, they must be authored (or edited) directly in KML."),
			SharedStyle(
				"yellowLineGreenPoly",
				LineStyle(
					Color(color.RGBA{R: 255, G: 255, B: 0, A: 127}),
					Width(4),
				),
				PolyStyle(
					Color(color.RGBA{R: 0, G: 255, B: 0, A: 127}),
				),
			),
			Placemark(
				Name("Absolute Extruded"),
				Description("Transparent green wall with yellow outlines"),
				StyleURL("#yellowLineGreenPoly"),
				LineString(
					Extrude(true),
					Tessellate(true),
					AltitudeMode("absolute"),
					Coordinates([]Coordinate{
						{-112.2550785337791, 36.07954952145647, 2357},
						{-112.2549277039738, 36.08117083492122, 2357},
						{-112.2552505069063, 36.08260761307279, 2357},
						{-112.2564540158376, 36.08395660588506, 2357},
						{-112.2580238976449, 36.08511401044813, 2357},
						{-112.2595218489022, 36.08584355239394, 2357},
						{-112.2608216347552, 36.08612634548589, 2357},
						{-112.262073428656, 36.08626019085147, 2357},
						{-112.2633204928495, 36.08621519860091, 2357},
						{-112.2644963846444, 36.08627897945274, 2357},
						{-112.2656969554589, 36.08649599090644, 2357},
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
	k := KML(
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
							{-77.0531553685479, 38.87053267794386, 100},
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
							{-77.05485125901023, 38.87076535397792, 100},
							{-77.05577677433152, 38.87008686581446, 100},
							{-77.05691162017543, 38.87054446963351, 100},
							{-77.05668055019126, 38.87154239798456, 100},
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
	k := KML(
		Document(
			SharedStyle(
				"transBluePoly",
				LineStyle(
					Width(1.5),
				),
				PolyStyle(
					Color(color.RGBA{R: 0, G: 0, B: 255, A: 125}),
				),
			),
			Placemark(
				Name("Building 41"),
				StyleURL("#transBluePoly"),
				Polygon(
					Extrude(true),
					AltitudeMode("relativeToGround"),
					OuterBoundaryIs(
						LinearRing(
							Coordinates([]Coordinate{
								{-122.0857412771483, 37.42227033155257, 17},
								{-122.0858169768481, 37.42231408832346, 17},
								{-122.085852582875, 37.42230337469744, 17},
								{-122.0858799945639, 37.42225686138789, 17},
								{-122.0858860101409, 37.4222311076138, 17},
								{-122.0858069157288, 37.42220250173855, 17},
								{-122.0858379542653, 37.42214027058678, 17},
								{-122.0856732640519, 37.42208690214408, 17},
								{-122.0856022926407, 37.42214885429042, 17},
								{-122.0855902778436, 37.422128290487, 17},
								{-122.0855841672237, 37.42208171967246, 17},
								{-122.0854852065741, 37.42210455874995, 17},
								{-122.0855067264352, 37.42214267949824, 17},
								{-122.0854430712915, 37.42212783846172, 17},
								{-122.0850990714904, 37.42251282407603, 17},
								{-122.0856769818632, 37.42281815323651, 17},
								{-122.0860162273783, 37.42244918858722, 17},
								{-122.0857260327004, 37.42229239604253, 17},
								{-122.0857412771483, 37.42227033155257, 17},
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
	k := KML(
		Document(
			Name("Highlighted Icon"),
			Description("Place your mouse over the icon to see it display the new icon"),
			SharedStyle(
				"highlightPlacemark",
				IconStyle(
					Icon(
						Href("http://maps.google.com/mapfiles/kml/paddle/red-stars.png"),
					),
				),
			),
			SharedStyle(
				"normalPlacemark",
				IconStyle(
					Icon(
						Href("http://maps.google.com/mapfiles/kml/paddle/wht-blank.png"),
					),
				),
			),
			SharedStyleMap(
				"exampleStyleMap",
				Pair(
					Key("normal"),
					StyleURL("#normalPlacemark"),
				),
				Pair(
					Key("highlight"),
					StyleURL("#highlightPlacemark"),
				),
			),
			Placemark(
				Name("Roll over this icon"),
				StyleURL("#exampleStyleMap"),
				Point(
					Coordinates(Coordinate{Lon: -122.0856545755255, Lat: 37.42243077405461}),
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
