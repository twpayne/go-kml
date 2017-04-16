# go-kml

[![Build Status](https://travis-ci.org/twpayne/go-kml.svg?branch=master)](https://travis-ci.org/twpayne/go-kml)
[![GoDoc](https://godoc.org/github.com/twpayne/go-kml?status.svg)](https://godoc.org/github.com/twpayne/go-kml)
[![Report Card](https://goreportcard.com/badge/github.com/twpayne/go-kml)](https://goreportcard.com/report/github.com/twpayne/go-kml)

Package go-kml provides convenience methods for creating and writing KML documents.

Example:

```go
func ExampleKML() {
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
```

[License](LICENSE)
