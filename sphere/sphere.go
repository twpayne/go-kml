// Package sphere contains convenience methods for generating coordinates on
// a sphere. All angles are measured in degrees.
package sphere

import (
	"math"

	"github.com/twpayne/go-kml"
)

const (
	degrees = 180 / math.Pi
	radians = math.Pi / 180
)

// A T is a sphere of radius R.
type T struct {
	R float64
}

var (
	// Unit is the unit sphere.
	Unit = T{R: 1}

	// FAI is the FAI sphere, measured in meters.
	FAI = T{R: 6371000}

	// WGS84 is a sphere whose radius is equal to the the semi-major axis of
	// the WGS84 ellipsoid, measured in meters.
	WGS84 = T{R: 6378137}
)

// Offset returns the coordinate at distance from origin in direction bearing.
func (t T) Offset(origin kml.Coordinate, distance, bearing float64) kml.Coordinate {
	lat := math.Asin(math.Sin(origin.Lat*radians)*math.Cos(distance/t.R) + math.Cos(origin.Lat*radians)*math.Sin(distance/t.R)*math.Cos(bearing*radians))
	lon := origin.Lon*radians + math.Atan2(math.Sin(bearing*radians)*math.Sin(distance/t.R)*math.Cos(origin.Lat*radians), math.Cos(distance/t.R)-math.Sin(origin.Lat*radians)*math.Sin(lat))
	return kml.Coordinate{
		Lon: lon * degrees,
		Lat: lat * degrees,
		Alt: origin.Alt,
	}
}

// Circle returns an array of kml.Coordinates that approximate a circle of
// radius radius centered on center with a maximum error of maxErr.
func (t T) Circle(center kml.Coordinate, radius, maxErr float64) []kml.Coordinate {
	numVertices := int(math.Ceil(math.Pi / math.Acos((radius-maxErr)/(radius+maxErr))))
	cs := make([]kml.Coordinate, numVertices+1)
	for i := 0; i < numVertices; i++ {
		cs[i] = t.Offset(center, radius, 360*float64(i)/float64(numVertices))
	}
	cs[numVertices] = cs[0]
	return cs
}
