package kml

import (
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
