package sphere

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/twpayne/go-kml/v2"
)

func TestSphereHaversineDistance(t *testing.T) {
	for i, tc := range []struct {
		sphere   T
		c1       kml.Coordinate
		c2       kml.Coordinate
		expected float64
		delta    float64
	}{
		{
			sphere:   FAI,
			c1:       kml.Coordinate{Lon: -108.6180554, Lat: 35.4325002},
			c2:       kml.Coordinate{Lon: -108.61, Lat: 35.43},
			expected: 781,
			delta:    1e-3,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.InDeltaf(t, tc.expected, tc.sphere.HaversineDistance(tc.c1, tc.c2), tc.delta, "")
		})
	}
}

func TestInitialBearingTo(t *testing.T) {
	for i, tc := range []struct {
		sphere   T
		c1       kml.Coordinate
		c2       kml.Coordinate
		expected float64
	}{
		{
			sphere:   FAI,
			c1:       kml.Coordinate{Lon: 0, Lat: 0},
			c2:       kml.Coordinate{Lon: 0, Lat: 1},
			expected: 0,
		},
		{
			sphere:   FAI,
			c1:       kml.Coordinate{Lon: 0, Lat: 0},
			c2:       kml.Coordinate{Lon: 1, Lat: 0},
			expected: 90,
		},
		{
			sphere:   FAI,
			c1:       kml.Coordinate{Lon: 0, Lat: 0},
			c2:       kml.Coordinate{Lon: 0, Lat: -1},
			expected: 180,
		},
		{
			sphere:   FAI,
			c1:       kml.Coordinate{Lon: 0, Lat: 0},
			c2:       kml.Coordinate{Lon: -1, Lat: 0},
			expected: -90,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.sphere.InitialBearingTo(tc.c1, tc.c2))
		})
	}
}
