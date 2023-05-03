package kml_test

import (
	"os"

	"github.com/twpayne/go-kml/v3"
)

func ExampleWriteKMZ() {
	doc := kml.KML(
		kml.Placemark(
			kml.Name("Zürich"),
			kml.Point(
				kml.Coordinates(
					kml.Coordinate{Lat: 47.374444, Lon: 8.541111},
				),
			),
		),
	)

	if err := kml.WriteKMZ(os.Stdout, map[string]any{
		"doc.kml": doc,
	}); err != nil {
		panic(err)
	}
}
