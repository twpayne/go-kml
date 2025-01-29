module github.com/twpayne/go-kml/examples

go 1.22

require (
	github.com/twpayne/go-gpx v1.4.1
	github.com/twpayne/go-kml/v3 v3.2.1
	github.com/twpayne/go-polyline v1.1.1
	github.com/twpayne/go-waypoint v0.1.0
)

require (
	github.com/twpayne/go-geom v1.5.7 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)

replace github.com/twpayne/go-kml/v3 => ..
