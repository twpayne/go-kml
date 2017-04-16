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
