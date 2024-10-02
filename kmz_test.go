package kml_test

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
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

	var buffer bytes.Buffer
	buffer.Grow(512)
	if err := kml.WriteKMZ(&buffer, map[string]any{
		"doc.kml": doc,
	}); err != nil {
		panic(err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buffer.Bytes()), int64(buffer.Len()))
	if err != nil {
		panic(err)
	}
	for _, zipFile := range zipReader.File {
		fmt.Println(zipFile.Name + ":")
		file, err := zipFile.Open()
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(os.Stdout, file); err != nil { //nolint:gosec
			panic(err)
		}
		file.Close()
	}

	// Output:
	// doc.kml:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <kml xmlns="http://www.opengis.net/kml/2.2"><Placemark><name>Zürich</name><Point><coordinates>8.541111,47.374444</coordinates></Point></Placemark></kml>
}
