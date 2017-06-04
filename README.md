# go-kml

[![Build Status](https://travis-ci.org/twpayne/go-kml.svg?branch=master)](https://travis-ci.org/twpayne/go-kml)
[![GoDoc](https://godoc.org/github.com/twpayne/go-kml?status.svg)](https://godoc.org/github.com/twpayne/go-kml)
[![Report Card](https://goreportcard.com/badge/github.com/twpayne/go-kml)](https://goreportcard.com/report/github.com/twpayne/go-kml)

Package kml provides convenience methods for creating and writing KML documents.

## Key Features

 * Simple API for building arbitrarily complex KML documents.
 * Support for all KML elements, including Google Earth `gx:` extensions.
 * Compatibilty with the standard library [`encoding/xml`](https://godoc.org/encoding/xml) package.
 * Pretty (neatly indented) and compact (minimum size) output formats.
 * Support for shared `Style` and `StyleMap` elements.
 * Simple mapping between functions and KML elements.


## Example

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
