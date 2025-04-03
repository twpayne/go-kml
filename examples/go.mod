module github.com/twpayne/go-kml/examples

go 1.24.0

toolchain go1.24.2

require (
	github.com/twpayne/go-gpx v1.4.1
	github.com/twpayne/go-kml/v3 v3.3.0
	github.com/twpayne/go-polyline v1.1.1
	github.com/twpayne/go-waypoint v0.1.0
)

require (
	github.com/twpayne/go-geom v1.6.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)

replace github.com/twpayne/go-kml/v3 => ..
