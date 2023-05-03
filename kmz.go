package kml

import (
	"archive/zip"
	"fmt"
	"io"
	"sort"
)

// WriteKMZ writes a KMZ file containing files to w. The values of the files map
// can be []bytes, strings, *KMLElements, *GxKMLElements, Elements, or
// io.Readers.
func WriteKMZ(w io.Writer, files map[string]any) error {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)

	zipWriter := zip.NewWriter(w)
	for _, filename := range names {
		zipFileWriter, err := zipWriter.Create(filename)
		if err != nil {
			return err
		}
		switch value := files[filename].(type) {
		case []byte:
			if _, err := zipFileWriter.Write(value); err != nil {
				return err
			}
		case string:
			if _, err := zipFileWriter.Write([]byte(value)); err != nil {
				return err
			}
		case *KMLElement:
			if err := value.Write(zipFileWriter); err != nil {
				return err
			}
		case *GxKMLElement:
			if err := value.Write(zipFileWriter); err != nil {
				return err
			}
		case Element:
			if err := KML(value).Write(zipFileWriter); err != nil {
				return err
			}
		case io.Reader:
			if _, err := io.Copy(zipFileWriter, value); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%T: unsupported type", value)
		}
	}
	return zipWriter.Close()
}
