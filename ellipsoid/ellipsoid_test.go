package ellipsoid_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-kml/v3"
	"github.com/twpayne/go-kml/v3/ellipsoid"
	"github.com/twpayne/go-kml/v3/sphere"
)

func TestEllipsoid_Distance(t *testing.T) {
	t.Parallel()

	sphere := sphere.WGS84
	ellipsoid := ellipsoid.WGS84

	const (
		deltaDegrees       = 0.02 // A delta of +0.02 degrees in both lat and lon at null island about 3145m
		sphericalTolerance = 11
	)

	for lat := -89; lat <= 89; lat++ {
		t.Run("lat"+strconv.Itoa(lat), func(t *testing.T) {
			t.Parallel()
			for lon := -179; lon <= 179; lon++ {
				t.Run("lon"+strconv.Itoa(lon), func(t *testing.T) {
					for _, deltaLat := range []float64{-deltaDegrees, deltaDegrees} {
						for _, deltaLon := range []float64{-deltaDegrees, deltaDegrees} {
							c1 := kml.Coordinate{Lat: float64(lat), Lon: float64(lon)}
							c2 := kml.Coordinate{Lat: float64(lat) + deltaLat, Lon: float64(lon) + deltaLon}
							sphericalDistance := sphere.HaversineDistance(c1, c2)
							ellipsoidDistance := ellipsoid.Distance(c1, c2)
							assert.True(t, math.Abs(ellipsoidDistance-sphericalDistance) < sphericalTolerance)
						}
					}
				})
			}
		})
	}
}
